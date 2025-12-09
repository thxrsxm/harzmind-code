package repl

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/executor"
	"github.com/thxrsxm/harzmind-code/internal/setup"
	"github.com/thxrsxm/rnbw"
)

type CMD struct {
	name    string
	info    string
	command func(r *REPL, args []string) error
}

func NewCMD(name, info string, command func(r *REPL, args []string) error) *CMD {
	return &CMD{
		name:    name,
		info:    info,
		command: command,
	}
}

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
	// /forget
	r.AddCommand(NewCMD(
		"forget",
		"Clear session context",
		forgetCMD,
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
}

func helpCMD(r *REPL, args []string) error {
	for _, v := range r.commands {
		r.out.Printf("'/%s' - %s\n", v.name, v.info)
	}
	return nil
}

func exitCMD(r *REPL, args []string) error {
	r.running = false
	r.out.CloseOutput()
	return nil
}

func initCMD(r *REPL, args []string) error {
	err := setup.SetupProjectDir()
	if err != nil {
		return err
	}
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Println("Project initiated")
	rnbw.ResetColor()
	return nil
}

func forgetCMD(r *REPL, args []string) error {
	r.messages = []api.Message{}
	r.messages = append(r.messages, api.Message{Role: "system", Content: ""})
	return nil
}

func modelsCMD(r *REPL, args []string) error {
	account, err := r.config.GetCurrentAccount()
	if err != nil {
		return err
	}
	models, err := api.GetModels(account.ApiUrl, account.ApiKey)
	if err != nil {
		return err
	}
	for _, v := range models {
		r.out.Println(v)
	}
	return nil
}

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
	r.out.Println(out)
	return nil
}

func infoCMD(r *REPL, args []string) error {
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Print("HarzMind Code")
	rnbw.ResetColor()
	r.out.Printf(" v%s\n", internal.VERSION_DATE)
	r.out.Println("Created by Erik Andrè Thürsam")
	return nil
}

func accCMD(r *REPL, args []string) error {
	if len(args) == 0 {
		// Show all accounts
		for _, v := range r.config.Accounts {
			r.out.Println(v)
			r.out.Println()
		}
		return nil
	} else if len(args) == 1 {
		switch args[0] {
		case "new":
			// Create a new account
			if err := r.handleAccountCreation(); err != nil {
				return err
			}
			return r.config.SaveConfig(internal.PATH_FILE_CONFIG)
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
			r.out.Printf("Successfully logged in to %s\n", args[1])
			rnbw.ResetColor()
			return r.config.SaveConfig(internal.PATH_FILE_CONFIG)
		case "logout":
			// Logout
			r.config.CurrentAccountName = ""
			return r.config.SaveConfig(internal.PATH_FILE_CONFIG)
		case "remove":
			// Remove account
			r.config.RemoveAccount(args[1])
			return r.config.SaveConfig(internal.PATH_FILE_CONFIG)
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

func modelCMD(r *REPL, args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("wrong format")
	}
	account, err := r.config.GetCurrentAccount()
	if err != nil {
		return err
	}
	account.Model = args[0]
	r.config.SaveConfig(internal.PATH_FILE_CONFIG)
	rnbw.ForgroundColor(rnbw.Green)
	r.out.Printf("Successfully changed model to %s for account %s\n", args[0], r.config.CurrentAccountName)
	rnbw.ResetColor()
	return nil
}

func editorCMD(r *REPL, args []string) error {
	if len(args) == 1 {
		return executor.OpenEditor(args[0], "")
	} else if len(args) >= 2 {
		return executor.OpenEditor(args[0], args[1])
	}
	return fmt.Errorf("wrong format")
}
