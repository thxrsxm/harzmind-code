// Package repl implements the Read-Eval-Print Loop for interactive user interaction,
// defining and executing slash commands for account management, codebase operations,
// external tool execution, and handling user messages sent to the LLM.
package repl

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/briandowns/spinner"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/config"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
	"golang.org/x/term"
)

// REPL represents a Read-Eval-Print Loop for the HarzMind Code application.
type REPL struct {
	running  bool
	out      *output.Output
	config   *config.Config
	reader   *bufio.Reader
	commands []CMD
	messages []api.Message
}

// NewREPL creates a new REPL instance.
// If outputFile is true, it will write output to a file.
func NewREPL(outputFile bool) (*REPL, error) {
	r := &REPL{
		running:  false,
		reader:   bufio.NewReader(os.Stdin),
		commands: []CMD{},
		messages: []api.Message{},
	}
	// Load config
	cnfg, err := config.LoadConfig(common.PATH_FILE_CONFIG)
	if err != nil {
		return nil, err
	}
	r.config = cnfg
	// Create output
	o, err := output.NewOutput(common.PATH_DIR_OUT, outputFile)
	if err != nil {
		return nil, err
	}
	r.out = o
	// Add commands
	addAllCommands(r)
	return r, nil
}

// handleUserMessage processes a user's message.
// It retrieves the codebase, constructs a system prompt,
// adds the user's message, and sends it to the API for processing.
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
	data, err := os.ReadFile(common.PATH_FILE_README)
	if err != nil {
		r.out.PrintfWarning("no %s file\n\n", common.FILE_IGNORE)
		data = []byte{}
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
	// Get current account
	account, err := r.config.GetCurrentAccount()
	if err != nil {
		return "", err
	}
	// Initialize and start the spinner for visual feedback
	// Use a dot spinner style
	s := spinner.New(spinner.CharSets[14], 100*time.Millisecond)
	// Start spinning in a goroutine
	s.Start()
	s.Suffix = " Sending codebase and querying LLM..."
	resp, err := api.SendMessage(account.ApiUrl, account.Model, account.ApiKey, r.messages)
	// Stop the spinner after the call completes
	s.Stop()
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

// printTitle displays the HarzMind Code title and help message.
func (r *REPL) printTitle() {
	fmt.Printf("\n\nWelcome to %s!\n\n\n", rnbw.String(rnbw.Green, "HarzMind Code"))
	rnbw.ForgroundColor(rnbw.Green)
	fmt.Print(TITLE)
	rnbw.ResetColor()
	fmt.Print("\n\n\n")
	helpCMD(r, nil)
}

// readInput reads a line of input from the user.
// It also writes the input to the output file if available.
func (r *REPL) readInput(writeToFile bool) (string, error) {
	input, err := r.reader.ReadString('\n')
	// Write user input to output file
	if writeToFile && r.out.File != nil {
		r.out.File.Print(input)
	}
	// Handle input error
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)
	return input, nil
}

// readPassword reads a password from the user securely.
func (r *REPL) readPassword() (string, error) {
	bytePassword, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		return "", err
	}
	r.out.Println()
	return string(bytePassword), nil
}

// AddCommand adds a new command to the REPL.
func (r *REPL) AddCommand(command *CMD) {
	r.commands = append(r.commands, *command)
}

// HandleCommand handles a slash command.
// It searches for a matching command and executes it.
func (r *REPL) HandleCommand(command string, args []string) error {
	for _, v := range r.commands {
		if command == v.name {
			return v.command(r, args)
		}
	}
	return fmt.Errorf("unknown command")
}

// Run starts the REPL event loop.
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
	// Graceful cleanup
	defer func() {
		r.out.CloseOutput()
	}()
	for r.running {
		rnbw.ResetColor()
		r.out.Println()
		// Show no account warning
		if _, err := r.config.GetCurrentAccount(); err != nil {
			r.out.PrintfWarning("no account\n\n")
		}
		// Show no ignore file warning
		if !codebase.IgnoreFileExists() {
			r.out.PrintfWarning("no %s file\n\n", common.FILE_IGNORE)
		}
		rnbw.ForgroundColor(rnbw.Green)
		r.out.Print("> ")
		rnbw.ResetColor()
		input, err := r.readInput(true)
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
			r.out.Println()
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
