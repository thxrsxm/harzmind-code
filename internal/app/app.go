// Package app orchestrates the application's lifecycle, processing command-line arguments,
// setting up configuration files, initializing the project if needed, and launching the REPL.
package app

import (
	"fmt"
	"os"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/args"
	"github.com/thxrsxm/harzmind-code/internal/codebase"
	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/executor"
	"github.com/thxrsxm/harzmind-code/internal/input"
	"github.com/thxrsxm/harzmind-code/internal/llmx"
	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/harzmind-code/internal/repl"
	"github.com/thxrsxm/harzmind-code/internal/setup"
	"github.com/thxrsxm/rnbw"
)

func Run() {
	// Parse command line flags
	args.Parse()
	// Initialize binary data directory
	if err := setup.SetupBinaryDataDir(); err != nil {
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "%v\n", err)
		rnbw.ResetColor()
		os.Exit(1)
	}
	// Init project
	if *args.InitFlag {
		if err := setup.SetupProjectDir(); err != nil {
			rnbw.ForgroundColor(rnbw.Red)
			fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
			rnbw.ResetColor()
		}
		rnbw.ForgroundColor(rnbw.Green)
		fmt.Fprint(os.Stdout, "Project initiated :)\n")
		rnbw.ResetColor()
	}
	// Initialize logger
	if *args.LogFlag {
		if err := logger.Init(common.PATH_FILE_LOG); err != nil {
			rnbw.ForgroundColor(rnbw.Red)
			fmt.Fprintf(os.Stderr, "[ERROR] failed to initialize logger: %v\n\n", err)
			rnbw.ResetColor()
			args.PrintUsage()
			os.Exit(1)
		}
		defer logger.Close()
	}
	logger.Log(logger.INFO, "HarzMind Code started v%s", internal.VERSION_DATE)
	logger.Log(logger.INFO, "Config directory: %s", common.PATH_DIR_BINARY_DATA)
	// Log project initialization
	if *args.InitFlag {
		logger.Log(logger.INFO, "%s", "project initiated")
	}
	// Initialize ouput
	if err := output.Init(common.PATH_DIR_OUT, *args.OutputFlag); err != nil {
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "%v\n", err)
		rnbw.ResetColor()
		logger.Log(logger.ERROR, "%v", err)
		os.Exit(1)
	}
	// Initialize input
	if err := input.Init(); err != nil {
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "%v\n", err)
		rnbw.ResetColor()
		logger.Log(logger.ERROR, "%v", err)
		os.Exit(1)
	}
	// Setup config file
	config, err := setup.SetupConfigFile()
	if err != nil {
		msg := fmt.Sprintf("setting up config file: %v", err)
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "[ERROR] %s\n", msg)
		rnbw.ResetColor()
		logger.Log(logger.ERROR, "%s", msg)
		os.Exit(1)
	}
	// Show help
	if *args.HelpFlag {
		args.PrintUsage()
		os.Exit(0)
	}
	// Show version
	if *args.VersionFlag {
		fmt.Fprintf(os.Stdout, "v%s\n", internal.VERSION_DATE)
		os.Exit(0)
	}
	// Create new llm client
	llmClient := llmx.NewLLMx()
	// Create new REPL
	r, err := repl.NewREPL(func(input string) error {
		// Get current account
		account, err := config.GetAccountManager().GetCurrentAccount()
		if err != nil {
			return err
		}
		// Handle user message
		resp, err := llmClient.HandleUserMessage(input, account.ApiUrl, account.Model, account.ApiKey)
		if err != nil {
			return err
		}
		output.Printf("\n%s\n", resp)
		return nil
	})
	if err != nil {
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "%v\n", err)
		rnbw.ResetColor()
		logger.Log(logger.ERROR, "%v", err)
		os.Exit(1)
	}
	// -----------------------
	// DEFINE AND ADD COMMANDS
	// -----------------------
	// /brocken
	r.AddCommand(repl.NewCMD(
		"brocken",
		"Shows the Brocken",
		func(arg string) error {
			common.PrintBrocken()
			return nil
		},
	))
	// /tree
	r.AddCommand(repl.NewCMD(
		"tree",
		"Codebase tree visualization",
		func(arg string) error {
			files, err := codebase.GetCodeBase(".")
			if err != nil {
				return err
			}
			output.Print(codebase.Tree(files))
			return nil
		},
	))
	// /info
	r.AddCommand(repl.NewCMD(
		"info",
		"Show info",
		func(arg string) error {
			rnbw.ForgroundColor(rnbw.Green)
			output.Print("HarzMind Code")
			rnbw.ResetColor()
			output.Printf(" v%s\n", internal.VERSION_DATE)
			output.Println("Created by Erik Andrè Thürsam")
			return nil
		},
	))
	// /session
	r.AddCommand(repl.NewCMD(
		"session",
		"Shows current session info",
		func(arg string) error {
			accountName := "-"
			model := "-"
			// Get current account
			account, err := config.GetAccountManager().GetCurrentAccount()
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
			// Print session
			output.Printf("Account:	'%s'\n", accountName)
			output.Printf("Model:		'%s'\n", model)
			output.Printf("Directory:	'%s'\n", dir)
			output.Printf("Context:	%d tokens\n", llmClient.GetTokens())
			return nil
		},
	))
	// /bash
	r.AddCommand(repl.NewCMD(
		"bash",
		"Run bash",
		func(arg string) error {
			out, err := executor.ExecuteBash(arg)
			if err != nil {
				rnbw.ForgroundColor(rnbw.Red)
			}
			output.Print(out)
			if len(out) >= 1 && out[len(out)-1] != '\n' {
				output.Println()
			}
			return nil
		},
	))
	// /editor
	r.AddCommand(repl.NewCMD(
		"editor",
		"Open CLI editor",
		func(arg string) error {
			if len(arg) == 0 {
				return fmt.Errorf("wrong format")
			}
			args := strings.Split(arg, " ")
			if len(args) == 1 {
				return executor.OpenEditor(args[0], "")
			} else if len(args) >= 2 {
				return executor.OpenEditor(args[0], args[1])
			}
			return fmt.Errorf("wrong format")
		},
	))
	// /clear
	r.AddCommand(repl.NewCMD(
		"clear",
		"Clear session context",
		func(arg string) error {
			llmClient.ClearMessages()
			rnbw.ForgroundColor(rnbw.Green)
			output.Println("Context was successfully deleted")
			rnbw.ResetColor()
			logger.Log(logger.INFO, "%s", "completed context clearing")
			return nil
		},
	))
	// /acc
	r.AddCommand(repl.NewCMD(
		"acc",
		"Account management",
		func(arg string) error {
			return config.GetAccountManager().HandleCommands(arg)
		},
	))
	// /model
	r.AddCommand(repl.NewCMD(
		"model",
		"Change model",
		func(arg string) error {
			if len(arg) == 0 {
				return fmt.Errorf("wrong format")
			}
			// Get current account
			account, err := config.GetAccountManager().GetCurrentAccount()
			if err != nil {
				return err
			}
			// Set model
			account.Model = arg
			// Save config file
			err = config.SaveConfig()
			if err != nil {
				return err
			}
			// Show success message
			rnbw.ForgroundColor(rnbw.Green)
			output.Printf("Successfully changed model to '%s' for account '%s'\n", arg, account.Name)
			rnbw.ResetColor()
			logger.Log(logger.INFO, "changed model to '%s' for account '%s'", arg, account.Name)
			return nil
		},
	))
	// /models
	r.AddCommand(repl.NewCMD(
		"models",
		"List all models",
		func(arg string) error {
			// Get current account
			account, err := config.GetAccountManager().GetCurrentAccount()
			if err != nil {
				return err
			}
			// Get available models
			models, err := api.GetModels(account.ApiUrl, account.ApiKey)
			logger.Log(logger.INFO, "%s", "fetching available models")
			if err != nil {
				return err
			}
			for i := range models {
				output.Println(models[i])
			}
			return nil
		},
	))
	// /init
	r.AddCommand(repl.NewCMD(
		"init",
		"Initialize project",
		func(arg string) error {
			err := setup.SetupProjectDir()
			if err != nil {
				return err
			}
			rnbw.ForgroundColor(rnbw.Green)
			output.Println("Project initiated")
			rnbw.ResetColor()
			logger.Log(logger.INFO, "%s", "project initiated")
			return nil
		},
	))
	// ----------
	// START REPL
	// ----------
	// Print title
	common.PrintTitle()
	// Print help
	fmt.Printf("Type %s to list all commands\n", rnbw.String(rnbw.Gray, "'/help'"))
	// Login to current account
	output.Println()
	if account, err := config.GetAccountManager().GetCurrentAccount(); err == nil {
		err = config.GetAccountManager().Login(account.Name)
		if err != nil {
			output.PrintWarning("no account\n")
		} else {
			rnbw.ForgroundColor(rnbw.Green)
			output.Printf("Successfully logged in to %s\n", account.Name)
			rnbw.ResetColor()
			logger.Log(logger.INFO, "logged in to '%s'", account.Name)
		}
	} else {
		output.PrintWarning("no account\n")
	}
	// Run REPL
	r.Run()
}
