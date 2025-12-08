package setup

import (
	"os"
	"path/filepath"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/config"
)

func SetupConfigFile() error {
	// Get binary path
	binDir, err := internal.GetBinaryPath()
	if err != nil {
		return err
	}
	configPath := filepath.Join(binDir, internal.PATH_FILE_CONFIG)
	// Check config file not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create new config file
		err := config.CreateConfig(internal.PATH_FILE_CONFIG)
		if err != nil {
			return nil
		}
	}
	return nil
}

func SetupProjectDir() error {
	// Create hzmind directory
	err := internal.CreateDirIfNotExists(internal.DIR_MAIN)
	if err != nil {
		return err
	}
	// Check HZMIND.md file
	err = internal.CreateFileIfNotExists(internal.PATH_FILE_README)
	if err != nil {
		return err
	}
	// Check .hzmignore file
	err = internal.CreateFileIfNotExists(internal.PATH_FILE_IGNORE)
	if err != nil {
		return err
	}
	return nil
}
