// Package executor provides utilities for executing external commands,
// including running bash scripts and opening files in terminal-based editors.
package executor

import (
	"fmt"
	"os"
	"os/exec"
)

// ExecuteBash executes a bash command and returns the output and error.
func ExecuteBash(command string) (string, error) {
	// Create a new bash command with the given command string
	cmd := exec.Command("bash", "-c", command)
	// Run the command and capture the output
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// OpenEditor opens a file in the specified editor.
func OpenEditor(editor, fileName string) error {
	// Check if editor binary exists
	if _, err := exec.LookPath(editor); err != nil {
		return fmt.Errorf("editor %s not found: %w", editor, err)
	}
	cmd := exec.Command(editor, fileName)
	// Direct the standard input, output, and error streams of the editor process to those of the current Go process.
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	// Run the command
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
