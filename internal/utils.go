package internal

import (
	"os"
	"path/filepath"
	"regexp"
)

func GetBinaryPath() (string, error) {
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Dir(exePath), nil
}

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

func CreateDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func IsValidURL(url string) bool {
	// Regular expression to validate URL
	re := regexp.MustCompile(`^https?://[a-zA-Z0-9-\.]+\.[a-zA-Z]{2,}(\/.*)?$`)
	return re.MatchString(url)
}
