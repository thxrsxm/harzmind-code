// Package repl implements the Read-Eval-Print Loop for interactive user interaction,
// defining and executing slash commands for account management, codebase operations,
// external tool execution, and handling user messages sent to the LLM.
package repl

import (
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/input"
	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
)

// REPL represents a Read-Eval-Print Loop for the HarzMind Code application.
type REPL struct {
	running  bool
	commands []CMD
	main     func(arg string) error
}

// NewREPL creates a new REPL instance.
// If outputFile is true, it will write output to a file.
func NewREPL(main func(arg string) error) (*REPL, error) {
	r := &REPL{
		running:  false,
		commands: []CMD{},
		main:     main,
	}
	// Add commands
	// /help
	r.AddCommand(NewCMD(
		"help",
		"List all commands",
		func(arg string) error { return r.PrintHelp() },
	))
	// /exit
	r.AddCommand(NewCMD(
		"exit",
		"End the conversation",
		func(arg string) error { return r.ExitREPL() },
	))
	return r, nil
}

func (r *REPL) PrintHelp() error {
	output.SetWriteMode(output.STDOUT)
	for _, v := range r.commands {
		output.Printf("'/%s' ", v.name)
		rnbw.ForgroundColor(rnbw.Gray)
		output.Printf("- %s\n", v.info)
		rnbw.ResetColor()
	}
	output.SetWriteMode(output.ALL)
	return nil
}

func (r *REPL) ExitREPL() error {
	r.running = false
	logger.Log(logger.INFO, "%s", "repl exit")
	return nil
}

// Run starts the REPL event loop.
func (r *REPL) Run() {
	r.running = true
	logger.Log(logger.INFO, "%s", "REPL started")
	for r.running {
		rnbw.ResetColor()
		output.Println()
		rnbw.ForgroundColor(rnbw.Green)
		output.Print("> ")
		rnbw.ResetColor()
		input, err := input.ReadInput(true)
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
			parts := strings.SplitN(input[1:], " ", 2)
			if len(parts) >= 1 {
				arg := ""
				if len(parts) == 2 {
					arg = parts[1]
				}
				err := r.HandleCommand(strings.ToLower(parts[0]), arg)
				if err != nil {
					output.PrintfError("%v\n", err)
					logger.Log(logger.ERROR, "%v", err)
				}
			} else {
				output.PrintlnError("unknown command")
			}
			continue
		}
		// Handle main
		if err := r.main(input); err != nil {
			output.Println()
			output.PrintfError("%v\n", err)
			logger.Log(logger.ERROR, "%v", err)
			continue
		}
	}
}
