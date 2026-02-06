// Package app orchestrates the application's lifecycle, processing command-line arguments,
// setting up configuration files, initializing the project if needed, and launching the REPL.
package app

import (
	"fmt"
	"os"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/args"
	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/input"
	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/harzmind-code/internal/repl"
	"github.com/thxrsxm/harzmind-code/internal/setup"
	"github.com/thxrsxm/rnbw"
)

func Run() {
	// Parse command line flags
	args.Parse()
	// Initialize logger
	if *args.LogFlag {
		if err := logger.Init(common.PATH_FILE_LOG); err != nil {
			rnbw.ForgroundColor(rnbw.Red)
			fmt.Fprintf(os.Stderr, "[ERROR] failed to initialize logger: %v\n", err)
			rnbw.ResetColor()
			os.Exit(1)
		}
		defer logger.Close()
	}
	logger.Log(logger.INFO, "%s", "HarzMind Code started")
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
	if err := setup.SetupConfigFile(); err != nil {
		msg := fmt.Sprintf("setting up config file: %v", err)
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "[ERROR] %s\n", msg)
		rnbw.ResetColor()
		logger.Log(logger.ERROR, "%s", msg)
		os.Exit(1)
	}
	// Show help
	if *args.HelpFlag {
		fmt.Println("HarzMind Code")
		fmt.Println("Usage: hzmind [flags]")
		fmt.Println("Flags:")
		args.PrintDefaults()
		os.Exit(0)
	}
	// Show version
	if *args.VersionFlag {
		fmt.Fprintf(os.Stdout, "v%s\n", internal.VERSION_DATE)
		os.Exit(0)
	}
	// Init project
	if *args.InitFlag {
		if err := setup.SetupProjectDir(); err != nil {
			rnbw.ForgroundColor(rnbw.Red)
			fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
			rnbw.ResetColor()
			logger.Log(logger.ERROR, "%v", err)
		}
		rnbw.ForgroundColor(rnbw.Green)
		fmt.Fprint(os.Stdout, "Project initiated\n")
		rnbw.ResetColor()
		logger.Log(logger.INFO, "%s", "project initiated")
	}
	// Create new REPL
	repl, err := repl.NewREPL()
	if err != nil {
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stdout, "%v\n", err)
		rnbw.ResetColor()
		logger.Log(logger.ERROR, "%v", err)
		os.Exit(1)
	}
	repl.Run()
}
