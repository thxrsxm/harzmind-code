package main

import (
	"fmt"
	"os"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/args"
	"github.com/thxrsxm/harzmind-code/internal/repl"
	"github.com/thxrsxm/harzmind-code/internal/setup"
	"github.com/thxrsxm/rnbw"
)

func main() {
	// Setup config file
	if err := setup.SetupConfigFile(); err != nil {
		rnbw.ForgroundColor(rnbw.Red)
		fmt.Fprintf(os.Stderr, "[ERROR] setting up config file: %v\n", err)
		rnbw.ResetColor()
		os.Exit(1)
	}
	// Parse command line flags
	args.Parse()
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
		err := setup.SetupProjectDir()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		fmt.Fprintf(os.Stdout, "Project initiated")
	}
	// Create new REPL
	repl, err := repl.NewREPL(*args.OutputFlag)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	repl.Run()
}
