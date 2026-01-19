// Package setup handles initial configuration and project setup tasks,
// including creating necessary directories, initializing default files like HZMIND.md and .hzmignore,
// and ensuring the application's configuration structure is in place.
package setup

import (
	"os"
	"path/filepath"

	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/config"
)

// SetupConfigFile sets up the configuration file.
// It checks if the configuration file exists, and if not, creates a new one.
func SetupConfigFile() error {
	// Get binary path
	binDir, err := common.GetBinaryPath()
	if err != nil {
		return err
	}
	configPath := filepath.Join(binDir, common.PATH_FILE_CONFIG)
	// Check config file not exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create new config file
		err := config.CreateConfig(common.PATH_FILE_CONFIG)
		if err != nil {
			return nil
		}
	}
	return nil
}

// SetupProjectDir sets up the project directory.
// It creates the main directory, README file, and ignore file if they do not exist.
func SetupProjectDir() error {
	// Create hzmind directory
	err := common.CreateDirIfNotExists(common.DIR_MAIN)
	if err != nil {
		return err
	}
	// Check HZMIND.md file
	err = common.CreateFileIfNotExists(common.PATH_FILE_README)
	if err != nil {
		return err
	}
	// Check .hzmignore file
	err = common.CreateFileIfNotExists(common.PATH_FILE_IGNORE)
	if err != nil {
		return err
	}
	return nil
}
