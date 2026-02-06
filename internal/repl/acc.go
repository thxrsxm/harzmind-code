package repl

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/config"
	"github.com/thxrsxm/harzmind-code/internal/input"
)

// handleAccountCreation handles the creation of a new account.
// It prompts the user for account details and validates the input.
// If successful, it adds the new account to the configuration and saves it.
func (r *REPL) handleAccountCreation() (*config.Account, error) {
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
	apiUrl, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	// Validate API URL
	if !common.IsValidURL(apiUrl) {
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
	account := config.NewAccount(name, apiUrl, apiKey, model)
	err = r.config.AddAccount(*account)
	if err != nil {
		return nil, err
	}
	return account, nil
}
