package main

import (
	"fmt"
	"os"

	"github.com/thxrsxm/harzmind-code/internal/args"
	"github.com/thxrsxm/harzmind-code/internal/repl"
	"github.com/thxrsxm/harzmind-code/internal/setup"
)

func main() {
	if args.Parse() {
		os.Exit(0)
	}
	apiToken, err := setup.Setup()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	repl, err := repl.NewREPL(apiToken)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	repl.Run()
}
