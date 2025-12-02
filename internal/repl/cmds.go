package repl

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/api"
	"github.com/thxrsxm/harzmind-code/internal/config"
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
	// /config
	r.AddCommand(NewCMD(
		"config",
		"Show config",
		configCMD,
	))
	// /editor
	r.AddCommand(NewCMD(
		"editor",
		"Open CLI editor",
		editorCMD,
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
	return nil
}

func initCMD(r *REPL, args []string) error {
	err := setup.SetupWorkingDir()
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
	config, err := config.LoadConfig(internal.PATH_FILE_CONFIG)
	if err != nil {
		return err
	}
	models, err := api.GetModels(config.API, r.token)
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

func configCMD(r *REPL, args []string) error {
	if r.config == nil {
		return fmt.Errorf("config is missing")
	}
	if len(args) == 0 {
		r.out.Println(r.config.String())
	} else if len(args) == 1 {
		switch args[0] {
		case "model":
			r.out.Println(r.config.Model)
		case "api":
			r.out.Println(r.config.API)
		case "outfile":
			r.out.Println(r.config.Outfile)
		case "reload":
			config, err := config.LoadConfig(internal.PATH_FILE_CONFIG)
			if err != nil {
				return err
			}
			r.config = config
			rnbw.ForgroundColor(rnbw.Green)
			r.out.Println("Config successfully reloaded")
			rnbw.ResetColor()
		}
	} else if len(args) == 2 {
		switch args[0] {
		case "model":
			r.config.Model = args[1]
		case "api":
			r.config.API = args[1]
		case "outfile":
			switch args[1] {
			case "true":
				r.config.Outfile = true
			case "false":
				r.config.Outfile = false
			default:
				return fmt.Errorf("format is wrong (only true/false)")
			}
		}
		// Save config changes
		err := r.config.Save(internal.PATH_FILE_CONFIG)
		if err == nil {
			rnbw.ForgroundColor(rnbw.Green)
			r.out.Println("Config successfully updated")
			rnbw.ResetColor()
		}
		return err
	}
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
