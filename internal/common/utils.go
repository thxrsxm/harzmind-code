package common

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"

	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
)

// CreateFileIfNotExists creates a file if it does not exist.
func CreateFileIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// CreateDirIfNotExists creates a directory if it does not exist.
func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

// FileExists checks if a file exists.
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

// IsValidURL checks if a URL is valid.
func IsValidURL(s string) bool {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return false
	}
	return true
}

// PrintTitle displays the HarzMind Code title and help message.
func PrintTitle() {
	fmt.Printf("\n\nWelcome to %s!\n\n\n", rnbw.String(rnbw.Green, "HarzMind Code"))
	rnbw.ForgroundColor(rnbw.Green)
	fmt.Print(TITLE)
	rnbw.ResetColor()
	fmt.Print("\n\n\n")
}

// PrintBrocken displays the ASCII art for the Brocken mountain.
func PrintBrocken() {
	output.SetWriteMode(output.STDOUT)
	output.Println(BROCKEN)
	output.SetWriteMode(output.ALL)
}

func getBinaryDataPath() string {
	goos := runtime.GOOS
	homeDir, err := os.UserHomeDir()
	if err != nil && goos != "windows" {
		// Fallback to current directory if home dir cannot be determined
		return "."
	}
	switch runtime.GOOS {
	case "darwin":
		// macOS: ~/Library/Application Support/hzmind/
		return filepath.Join(homeDir, "Library", "Application Support", DIR_MAIN)
	case "linux":
		// Linux: ~/.config/hzmind/
		return filepath.Join(homeDir, ".config", DIR_MAIN)
	case "windows":
		// Windows: binary path
		exePath, err := os.Executable()
		if err != nil {
			// Fallback to current directory if home dir cannot be determined
			return "."
		}
		return filepath.Dir(exePath)
	default:
		// Fallback: ~/.hzmind/
		return filepath.Join(homeDir, "."+DIR_MAIN)
	}
}
