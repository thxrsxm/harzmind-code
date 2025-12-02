package setup

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/config"
)

func Setup() (string, error) {
	// Load .env file variables
	_ = godotenv.Load()
	// Load API token
	apiToken := os.Getenv(internal.API_TOKEN_NAME)
	if len(apiToken) == 0 {
		return "", fmt.Errorf("API token is missing (%s)", internal.API_TOKEN_NAME)
	}
	return apiToken, nil
}

func SetupWorkingDir() error {
	// Check hzmind directory not exists
	if _, err := os.Stat(internal.DIR_MAIN); os.IsNotExist(err) {
		// Create hzmind directory
		err := os.Mkdir(internal.DIR_MAIN, 0755)
		if err != nil {
			return err
		}
	}
	// Check config file not exists
	if _, err := os.Stat(internal.PATH_FILE_CONFIG); os.IsNotExist(err) {
		// Create config file
		err := config.CreateConfig(internal.PATH_FILE_CONFIG)
		if err != nil {
			return nil
		}
	}
	// Check HZMIND.md not exists
	if _, err := os.Stat(internal.PATH_FILE_README); os.IsNotExist(err) {
		// Create HZMIND.md
		readme, err := os.Create(internal.PATH_FILE_README)
		if err != nil {
			return err
		}
		defer readme.Close()
	}
	return nil
}
