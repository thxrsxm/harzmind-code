package main

import (
	"fmt"
	"os"

	"github.com/thxrsxm/harzmind-code/internal/repl"
)

func main() {
	repl, err := repl.NewREPL()
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	repl.Run()
}
