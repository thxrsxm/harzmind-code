package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/config"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
	"golang.org/x/term"
)

type REPL struct {
	running  bool
	out      *output.Output
	config   *config.Config
	reader   *bufio.Reader
	commands []CMD
	messages []api.Message
}

func NewREPL(outputFile bool) (*REPL, error) {
	r := &REPL{
		running:  false,
		reader:   bufio.NewReader(os.Stdin),
		commands: []CMD{},
		messages: []api.Message{},
	}
	// Load config
	cnfg, err := config.LoadConfig(internal.PATH_FILE_CONFIG)
	if err != nil {
		return nil, err
	}
	r.config = cnfg
	// Create output
	o, err := output.NewOutput(internal.PATH_DIR_OUT, outputFile)
	if err != nil {
		return nil, err
	}
	r.out = o
	// Add commands
	addAllCommands(r)
	// DEBUG
	// Load API token
	//apiToken := os.Getenv("HARZMIND_API_TOKEN")
	//r.account = config.NewAccount("test", "https://api.mammouth.ai/v1/chat/completions", apiToken, "grok-3-mini")
	return r, nil
}

func (r *REPL) AddCommand(command *CMD) {
	r.commands = append(r.commands, *command)
}

func (r *REPL) handleUserMessage(msg string) (string, error) {
	// Get code base
	files, err := codebase.GetCodeBase(".")
	if err != nil {
		return "", err
	}
	jsonCodeBase, err := json.Marshal(files)
	if err != nil {
		return "", err
	}
	// Load HZMIND.md
	data, err := os.ReadFile(internal.PATH_FILE_README)
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
	account, err := r.config.GetCurrentAccount()
	if err != nil {
		return "", err
	}
	// TODO check account members before send the request
	resp, err := api.SendMessage(account.ApiUrl, account.Model, account.ApiKey, r.messages)
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

func (r *REPL) readInput() (string, error) {
	input, err := r.reader.ReadString('\n')
	// Write user input to output file
	if r.out.File != nil {
		r.out.File.Print(input)
	}
	// Handle input error
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, nil
}

func (r *REPL) readPassword() (string, error) {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	r.out.Println()
	return string(bytePassword), nil
}

func (r *REPL) HandleCommand(command string, args []string) error {
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
	r.messages = append(r.messages, api.Message{Role: "system", Content: ""})
	// Login to current account
	if account, err := r.config.GetCurrentAccount(); err == nil {
		rnbw.ForgroundColor(rnbw.Green)
		r.out.Printf("\nSuccessfully logged in to %s\n", account.Name)
		rnbw.ResetColor()
	}
	for r.running {
		rnbw.ResetColor()
		r.out.Println()
		if _, err := r.config.GetCurrentAccount(); err != nil {
			r.out.PrintlnWarning("no account")
			r.out.Println()
		}
		rnbw.ForgroundColor(rnbw.Yellow)
		r.out.Print("> ")
		rnbw.ResetColor()
		input, err := r.readInput()
		if err != nil {
			r.out.PrintfError("%v\n", err)
			continue
		}
		// Handle input is empty
		if len(input) == 0 {
			continue
		}
		// Handle slash command
		if input[0] == '/' && len(input) > 1 {
			args := strings.Split(input[1:], " ")
			if len(args) >= 1 {
				err := r.HandleCommand(strings.ToLower(args[0]), args[1:])
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
