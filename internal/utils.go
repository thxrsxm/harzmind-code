// Package internal provides various internal functions and constants.
package internal

import (
	"os"
	"path/filepath"
	"regexp"
)

// GetBinaryPath returns the path of the binary executable.
func GetBinaryPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

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
func IsValidURL(url string) bool {
	// Regular expression to validate URL
	re := regexp.MustCompile(`^https?://`)
	return re.MatchString(url)
}
