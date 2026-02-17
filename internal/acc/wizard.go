package acc

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/input"
)

func handleAccountCreation() (*Account, error) {
	fmt.Println("Create account")
	// Get account name
	fmt.Print("Name: ")
	name, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	// Get API url
	fmt.Print("API Url: ")
	apiURL, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	// Validate API URL
	if !common.IsValidURL(apiURL) {
		return nil, fmt.Errorf("invalid api url")
	}
	// Get API token (secure)
	fmt.Print("API Token: ")
	apiKey, err := input.ReadPassword()
	if err != nil {
		return nil, err
	}
	// Check for empty API token
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("api token cannot be empty")
	}
	fmt.Print("Model (optional): ")
	model, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	account := NewAccount(name, apiURL, apiKey, model)
	return account, nil
}
