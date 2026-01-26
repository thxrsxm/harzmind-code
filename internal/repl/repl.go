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
	"github.com/pkoukk/tiktoken-go"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/config"
	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
	"golang.org/x/term"
)

// REPL represents a Read-Eval-Print Loop for the HarzMind Code application.
type REPL struct {
	running  bool
	config   *config.Config
	reader   *bufio.Reader
	tokens   int
	commands []CMD
	messages []api.Message
}

// NewREPL creates a new REPL instance.
// If outputFile is true, it will write output to a file.
func NewREPL() (*REPL, error) {
	r := &REPL{
		running:  false,
		reader:   bufio.NewReader(os.Stdin),
		tokens:   0,
		commands: []CMD{},
		messages: []api.Message{},
	}
	// Load config
	cnfg, err := config.LoadConfig(common.PATH_FILE_CONFIG)
	if err != nil {
		return nil, err
	}
	r.config = cnfg
	// Add commands
	addAllCommands(r)
	return r, nil
}

// updateTokens calculates the current token count in the conversation context.
func (r *REPL) updateTokens() {
	account, err := r.config.GetCurrentAccount()
	model := ""
	if err == nil {
		model = account.Model
	}
	encoding, err := tiktoken.EncodingForModel(model)
	if err != nil {
		// Fallback to cl100k_base (GPT-4 encoding)
		encoding, _ = tiktoken.GetEncoding("cl100k_base")
	}
	count := 0
	for _, v := range r.messages {
		count += len(encoding.Encode(v.Content, nil, nil))
	}
	r.tokens = count
}

// createSystemPrompt builds the system prompt by combining HZMIND.md and codebase data.
func (r *REPL) createSystemPrompt() (string, error) {
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
		output.PrintfWarning("no %s file\n\n", common.FILE_IGNORE)
		logger.Log(logger.ERROR, "%v", err)
		data = []byte{}
	}
	// Create System Prompt message
	return string(data) + "\n\n## Codebase\n\n" + string(jsonCodeBase), nil
}

// handleUserMessage processes a user's message.
// It retrieves the codebase, constructs a system prompt,
// adds the user's message, and sends it to the API for processing.
func (r *REPL) handleUserMessage(msg string) (string, error) {
	logger.Log(logger.INFO, "handling user message (length: %d chars)", len(msg))
	sysPrompt, err := r.createSystemPrompt()
	if err != nil {
		return "", err
	}
	if len(r.messages) > 0 {
		r.messages[0].Content = sysPrompt
	} else {
		r.messages = append(r.messages, api.Message{Role: "system", Content: sysPrompt})
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
	logger.Log(logger.INFO, "%s", "sending codebase and querying LLM")
	// Stop the spinner after the call completes
	s.Stop()
	if err != nil {
		logger.Log(logger.ERROR, "API call failed for user message: %v", err)
		// Remove last message from messages (user message)
		if len(r.messages) >= 1 {
			r.messages = r.messages[:len(r.messages)-1]
		}
		return "", err
	}
	logger.Log(logger.INFO, "received response from API for user message")
	// Add AI message to messages
	r.messages = append(r.messages, api.Message{
		Role:    "assistant",
		Content: resp,
	})
	// Update tokens amount
	r.updateTokens()
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
	if writeToFile {
		output.SetWriteMode(output.FILE)
		output.Print(input)
		output.SetWriteMode(output.ALL)
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
	output.Println()
	return string(bytePassword), nil
}

func (r *REPL) GetContext() string {
	var sb strings.Builder
	for _, v := range r.messages {
		sb.WriteString(v.Content)
	}
	return sb.String()
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
			logger.Log(logger.INFO, "command '/%s' was entered", command)
			return v.command(r, args)
		}
	}
	logger.Log(logger.ERROR, "unknown command was entered: /%s", command)
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
		output.Printf("\nSuccessfully logged in to %s\n", account.Name)
		rnbw.ResetColor()
		logger.Log(logger.INFO, "logged in to '%s'", account.Name)
	}
	// Graceful cleanup
	defer func() {
		logger.Log(logger.INFO, "%s", "graceful cleanup")
		logger.Close()
		output.Close()
	}()
	logger.Log(logger.INFO, "%s", "REPL started")
	for r.running {
		rnbw.ResetColor()
		output.Println()
		// Show no account warning
		if _, err := r.config.GetCurrentAccount(); err != nil {
			output.PrintfWarning("no account\n\n")
		}
		// Show no ignore file warning
		if !codebase.IgnoreFileExists() {
			output.PrintfWarning("no %s file\n\n", common.FILE_IGNORE)
		}
		rnbw.ForgroundColor(rnbw.Green)
		output.Print("> ")
		rnbw.ResetColor()
		input, err := r.readInput(true)
		if err != nil {
			output.PrintfError("%v\n", err)
			logger.Log(logger.ERROR, "%v", err)
			continue
		}
		// Handle input is empty
		if len(input) == 0 {
			continue
		}
		// Handle slash command
		if input[0] == '/' && len(input) > 1 {
			output.Println()
			args := strings.Split(input[1:], " ")
			if len(args) >= 1 {
				err := r.HandleCommand(strings.ToLower(args[0]), args[1:])
				if err != nil {
					output.PrintfError("%v\n", err)
					logger.Log(logger.ERROR, "%v", err)
				}
			} else {
				output.PrintlnError("unknown command")
			}
			continue
		}
		// Handle user message
		resp, err := r.handleUserMessage(input)
		if err != nil {
			output.PrintfError("%v\n", err)
			logger.Log(logger.ERROR, "%v", err)
			continue
		}
		output.Printf("\n%s\n", resp)
	}
}
