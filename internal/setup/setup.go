// Package setup handles initial configuration and project setup tasks,
// including creating necessary directories, initializing default files like HZMIND.md and .hzmignore,
// and ensuring the application's configuration structure is in place.
package setup

import (
	"os"

	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/config"
)

func SetupBinaryDataDir() error {
	return os.MkdirAll(common.PATH_DIR_BINARY_DATA, 0755)
}

func SetupConfigFile() (*config.Config, error) {
	// Check config file not exists
	if _, err := os.Stat(common.PATH_FILE_CONFIG); os.IsNotExist(err) {
		// Create new config file
		config, err := config.NewConfig(common.PATH_FILE_CONFIG)
		if err != nil {
			return nil, err
		}
		return config, nil
	}
	return config.LoadConfig(common.PATH_FILE_CONFIG)
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
