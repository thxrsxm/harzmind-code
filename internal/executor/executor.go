package executor

import (
	"os"
	"os/exec"
)

func ExecuteBash(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func OpenEditor(editor, fileName string) error {
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
