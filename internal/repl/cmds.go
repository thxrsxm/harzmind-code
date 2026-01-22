package repl

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/executor"
	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/setup"
	"github.com/thxrsxm/rnbw"
)

// CMD represents a command that can be executed in the REPL.
type CMD struct {
	name    string
	info    string
	command func(r *REPL, args []string) error
}

// NewCMD creates a new command instance.
func NewCMD(name, info string, command func(r *REPL, args []string) error) *CMD {
	return &CMD{
		name:    name,
		info:    info,
		command: command,
	}
}

// addAllCommands adds all available commands to the REPL.
func addAllCommands(r *REPL) {
	// /help
	r.AddCommand(NewCMD(
		"help",
		"List all commands",
		helpCMD,
	))
	// /exit
	r.AddCommand(NewCMD(
		"exit",
		"End the conversation",
		exitCMD,
	))
	// /init
	r.AddCommand(NewCMD(
		"init",
		"Initialize project",
		initCMD,
	))
	// /clear
	r.AddCommand(NewCMD(
		"clear",
		"Clear session context",
		clearCMD,
	))
	// /models
	r.AddCommand(NewCMD(
		"models",
		"List all models",
		modelsCMD,
	))
	// /bash
	r.AddCommand(NewCMD(
		"bash",
		"Run bash",
		bashCMD,
	))
	// /info
	r.AddCommand(NewCMD(
		"info",
		"Show info",
		infoCMD,
	))
	// /acc
	r.AddCommand(NewCMD(
		"acc",
		"Account management",
		accCMD,
	))
	// /editor
	r.AddCommand(NewCMD(
		"editor",
		"Open CLI editor",
		editorCMD,
	))
	// model
	r.AddCommand(NewCMD(
		"model",
		"Change model",
		modelCMD,
	))
	// tree
	r.AddCommand(NewCMD(
		"tree",
		"Codebase tree visualization",
		treeCMD,
	))
	// brocken
	r.AddCommand(NewCMD(
		"brocken",
		"Shows the Brocken",
		brockenCMD,
	))
	// status
	r.AddCommand(NewCMD(
		"status",
		"Shows current session info",
		statusCMD,
	))
	// Sort commands
	sort.Slice(r.commands, func(i, j int) bool {
		return r.commands[i].name < r.commands[j].name
	})
}

// helpCMD displays help information for all commands.
func helpCMD(r *REPL, args []string) error {
	for _, v := range r.commands {
		r.out.Stdout.Printf("'/%s' ", v.name)
		rnbw.ForgroundColor(rnbw.Gray)
		r.out.Stdout.Printf("- %s\n", v.info)
		rnbw.ResetColor()
	}
	return nil
}

// exitCMD exits the REPL.
func exitCMD(r *REPL, args []string) error {
	r.running = false
	r.out.CloseOutput()
	logger.Log(logger.INFO, "%s", "exit")
	return nil
}

// initCMD initializes a new project.
func initCMD(r *REPL, args []string) error {
	err := setup.SetupProjectDir()
	if err != nil {
		return err
	}
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Println("Project initiated")
	rnbw.ResetColor()
	logger.Log(logger.INFO, "%s", "project initiated")
	return nil
}

// clearCMD clears the session context.
func clearCMD(r *REPL, args []string) error {
	r.messages = []api.Message{}
	r.messages = append(r.messages, api.Message{Role: "system", Content: ""})
	r.updateTokens()
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Println("Context was successfully deleted")
	rnbw.ResetColor()
	logger.Log(logger.INFO, "%s", "completed context clear")
	return nil
}

// modelsCMD lists all available models.
func modelsCMD(r *REPL, args []string) error {
	account, err := r.config.GetCurrentAccount()
	if err != nil {
		return err
	}
	models, err := api.GetModels(account.ApiUrl, account.ApiKey)
	logger.Log(logger.INFO, "%s", "fetching available models")
	if err != nil {
		return err
	}
	for _, v := range models {
		r.out.Println(v)
	}
	return nil
}

// bashCMD runs a bash command.
func bashCMD(r *REPL, args []string) error {
	var sb strings.Builder
	for i, v := range args {
		sb.WriteString(v)
		if i < len(args)-1 {
			sb.WriteRune(' ')
		}
	}
	out, err := executor.ExecuteBash(sb.String())
	if err != nil {
		rnbw.ForgroundColor(rnbw.Red)
	}
	r.out.Print(out)
	if len(out) >= 1 && out[len(out)-1] != '\n' {
		r.out.Println()
	}
	return nil
}

// infoCMD displays information about the HarzMind Code application.
func infoCMD(r *REPL, args []string) error {
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Print("HarzMind Code")
	rnbw.ResetColor()
	r.out.Printf(" v%s\n", internal.VERSION_DATE)
	r.out.Println("Created by Erik Andrè Thürsam")
	return nil
}

// accCMD manages accounts.
func accCMD(r *REPL, args []string) error {
	if len(args) == 0 {
		// Show all accounts
		for i, v := range r.config.Accounts {
			r.out.Println(v)
			if i < len(r.config.Accounts)-1 {
				r.out.Println()
			}
		}
		return nil
	} else if len(args) == 1 {
		switch args[0] {
		case "new":
			// Create a new account
			account, err := r.handleAccountCreation()
			if err != nil {
				return err
			}
			if err := r.config.SaveConfig(common.PATH_FILE_CONFIG); err != nil {
				rnbw.ForgroundColor(rnbw.Green)
				r.out.Printf("Successfully created the account '%s'\n", account.Name)
				rnbw.ResetColor()
				logger.Log(logger.INFO, "created account '%s'", account.Name)
			}
			return nil
		case "logout":
			// Logout
			account := r.config.CurrentAccountName
			r.config.CurrentAccountName = ""
			if err := r.config.SaveConfig(common.PATH_FILE_CONFIG); err != nil {
				return err
			}
			rnbw.ForgroundColor(rnbw.Green)
			r.out.Printf("Successfully logged out from '%s'\n", account)
			rnbw.ResetColor()
			logger.Log(logger.INFO, "logged out from '%s'", account)
			return nil
		default:
			return fmt.Errorf("command not found")
		}
	} else if len(args) == 2 {
		if len(args[1]) == 0 {
			return fmt.Errorf("argument is missing")
		}
		switch args[0] {
		case "login":
			// Login
			if _, err := r.config.GetAccount(args[1]); err != nil {
				return err
			}
			r.config.CurrentAccountName = args[1]
			rnbw.ForgroundColor(rnbw.Green)
			r.out.Printf("Successfully logged in to '%s'\n", args[1])
			rnbw.ResetColor()
			logger.Log(logger.INFO, "logged in to '%s'", args[1])
			return r.config.SaveConfig(common.PATH_FILE_CONFIG)
		case "remove":
			// Remove account
			r.config.RemoveAccount(args[1])
			if err := r.config.SaveConfig(common.PATH_FILE_CONFIG); err != nil {
				return err
			}
			r.out.Printf("Successfully removed account '%s'\n", args[1])
			logger.Log(logger.WARNING, "removed account '%s'", args[1])
			return nil
		case "info":
			// Show account info
			account, err := r.config.GetAccount(args[1])
			if err != nil {
				return err
			}
			r.out.Println(account)
		default:
			return fmt.Errorf("command not found")
		}
	}
	return fmt.Errorf("command not found")
}

// modelCMD changes the model for the current account.
func modelCMD(r *REPL, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong format")
	}
	account, err := r.config.GetCurrentAccount()
	if err != nil {
		return err
	}
	account.Model = args[0]
	err = r.config.SaveConfig(common.PATH_FILE_CONFIG)
	if err != nil {
		return err
	}
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Printf("Successfully changed model to '%s' for account '%s'\n", args[0], r.config.CurrentAccountName)
	rnbw.ResetColor()
	logger.Log(logger.INFO, "changed model to '%s' for account '%s'", args[0], r.config.CurrentAccountName)
	return nil
}

// editorCMD opens the CLI editor.
func editorCMD(r *REPL, args []string) error {
	if len(args) == 1 {
		return executor.OpenEditor(args[0], "")
	} else if len(args) >= 2 {
		return executor.OpenEditor(args[0], args[1])
	}
	return fmt.Errorf("wrong format")
}

// treeCMD generates a hierarchical tree view of all files and directories included in the codebase.
func treeCMD(r *REPL, args []string) error {
	files, err := codebase.GetCodeBase(".")
	if err != nil {
		return err
	}
	r.out.Print(codebase.Tree(files))
	return nil
}

// brockenCMD displays the ASCII art for the Brocken mountain.
func brockenCMD(r *REPL, args []string) error {
	r.out.Stdout.Println(BROCKEN)
	return nil
}

// statusCMD shows the current session status, including account, model, working directory and token usage.
func statusCMD(r *REPL, args []string) error {
	accountName := "-"
	model := "-"
	// Get current account
	account, err := r.config.GetCurrentAccount()
	if err == nil {
		accountName = account.Name
		model = account.Model
		if len(model) == 0 {
			model = "-"
		}
	}
	// Get working directory
	dir, err := os.Getwd()
	if err != nil {
		dir = "-"
		logger.Log(logger.ERROR, "%v", err)
	}
	// Print status
	r.out.Printf("Account:	'%s'\n", accountName)
	r.out.Printf("Model:		'%s'\n", model)
	r.out.Printf("Directory:	'%s'\n", dir)
	r.out.Printf("Context:	%d tokens\n", r.tokens)
	return nil
}
