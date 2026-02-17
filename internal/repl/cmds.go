package repl

import (
	"fmt"
	"sort"

	"github.com/thxrsxm/harzmind-code/internal/logger"
)

// CMD represents a command that can be executed in the REPL.
type CMD struct {
	name    string
	info    string
	command func(arg string) error
}

// NewCMD creates a new command instance.
func NewCMD(name, info string, command func(arg string) error) *CMD {
	return &CMD{
		name:    name,
		info:    info,
		command: command,
	}
}

// AddCommand adds a new command to the REPL.
func (r *REPL) AddCommand(command *CMD) {
	r.commands = append(r.commands, *command)
	r.sortCommands()
}

// HandleCommand looks up and executes a registered slash command.
func (r *REPL) HandleCommand(command, arg string) error {
	for i := range r.commands {
		if command == r.commands[i].name {
			logger.Log(logger.INFO, "command '/%s' was entered", command)
			return r.commands[i].command(arg)
		}
	}
	logger.Log(logger.ERROR, "unknown command was entered: /%s", command)
	return fmt.Errorf("unknown command")
}

// sortCommands sorts the registered commands alphabetically by name (case-sensitive).
func (r *REPL) sortCommands() {
	sort.Slice(r.commands, func(i, j int) bool {
		return r.commands[i].name < r.commands[j].name
	})
}
