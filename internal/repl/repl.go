package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/config"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
)

type REPL struct {
	running  bool
	token    string
	out      *output.Output
	config   *config.Config
	commands []CMD
	messages []api.Message
}

func NewREPL() (*REPL, error) {
	apiToken := os.Getenv(API_TOKEN_NAME)
	if len(apiToken) == 0 {
		return nil, fmt.Errorf("API token is missing (%s)", API_TOKEN_NAME)
	}
	r := &REPL{
		running:  false,
		token:    apiToken,
		commands: []CMD{},
		messages: []api.Message{},
	}
	initCMD(nil, nil)
	// Load config
	config, err := config.LoadConfig(PATH_FILE_CONFIG)
	if err != nil {
		return nil, err
	}
	r.config = config
	// Create output
	o, err := output.NewOutput(PATH_DIR_OUT, config.Outfile)
	if err != nil {
		return nil, err
	}
	r.out = o
	// Add commands
	addAllCommands(r)
	return r, nil
}

func (r *REPL) AddCommand(command *CMD) {
	r.commands = append(r.commands, *command)
}

func (r *REPL) handleUserMessage(msg string) (string, error) {
	// Get code base
	files := codebase.GetCodeBase(".", r.config.Ignore)
	jsonCodeBase, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	// Load HZMIND.md
	data, err := os.ReadFile(PATH_FILE_README)
	if err != nil {
		r.out.PrintWarning("no HZMIND.md file!")
	}
	// Overwrite System Prompt message
	sysprompt := string(data) + "\n\n## Codebase\n\n" + string(jsonCodeBase)
	if len(r.messages) > 0 {
		r.messages[0].Content = sysprompt
	} else {
		r.messages = append(r.messages, api.Message{Role: "system", Content: sysprompt})
	}
	// Add user message to messages
	userMsg := api.Message{
		Role:    "user",
		Content: msg,
	}
	r.messages = append(r.messages, userMsg)
	// Send API-Request
	resp, err := api.SendMessage(r.config.API, r.config.Model, r.token, r.messages)
	if err != nil {
		// Remove last message from messages (user message)
		if len(r.messages) >= 1 {
			r.messages = r.messages[:len(r.messages)-1]
		}
		return "", err
	}
	// Add AI message to messages
	r.messages = append(r.messages, api.Message{
		Role:    "assistant",
		Content: resp,
	})
	return resp, nil
}

func (r *REPL) printTitle() {
	fmt.Printf("\n\nWelcome to %s!\n\n\n", rnbw.String(rnbw.Green, "HarzMind Code"))
	rnbw.ForgroundColor(rnbw.Green)
	fmt.Print(TITLE)
	rnbw.ResetColor()
	fmt.Print("\n\n\n")
	helpCMD(r, nil)
}

func (r *REPL) handleCommand(command string, args []string) error {
	for _, v := range r.commands {
		if command == v.name {
			return v.command(r, args)
		}
	}
	return fmt.Errorf("unknown command")
}

func (r *REPL) Run() {
	r.running = true
	r.printTitle()
	rnbw.ResetColor()
	reader := bufio.NewReader(os.Stdin)
	r.messages = append(r.messages, api.Message{Role: "system", Content: ""})
	for r.running {
		rnbw.ResetColor()
		r.out.Println()
		rnbw.ForgroundColor(rnbw.Yellow)
		r.out.Print("> ")
		rnbw.ResetColor()
		input, err := reader.ReadString('\n')
		// Write user input to output file
		if r.config.Outfile && r.out.File != nil {
			r.out.File.Print(input)
		}
		// Handle input error
		if err != nil {
			r.out.PrintfError("reading input: %v\n", err)
			continue
		}
		input = strings.TrimSpace(input)
		// Handle input is empty
		if len(input) == 0 {
			continue
		}
		// Handle slash command
		if input[0] == '/' && len(input) > 1 {
			args := strings.Split(input[1:], " ")
			if len(args) >= 1 {
				err := r.handleCommand(strings.ToLower(args[0]), args[1:])
				if err != nil {
					r.out.PrintfError("%v\n", err)
				}
			} else {
				r.out.PrintlnError("unknown command")
			}
			continue
		}
		// Handle user message
		resp, err := r.handleUserMessage(input)
		if err != nil {
			r.out.PrintfError("%v\n", err)
			continue
		}
		r.out.Printf("\n%s\n", resp)
	}
}
