package repl

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/config"
)

// handleAccountCreation handles the creation of a new account.
// It prompts the user for account details and validates the input.
// If successful, it adds the new account to the configuration and saves it.
func (r *REPL) handleAccountCreation() error {
	fmt.Println("Create account")
	// Get account name
	fmt.Print("Name: ")
	name, err := r.readInput(false)
	if err != nil {
		return err
	}
	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// Get API url
	fmt.Print("API Url: ")
	apiUrl, err := r.readInput(false)
	if err != nil {
		return err
	}
	// Validate API URL
	if !common.IsValidURL(apiUrl) {
		return fmt.Errorf("invalid api url")
	}
	// Get API token (secure)
	fmt.Print("API Token: ")
	apiKey, err := r.readPassword()
	if err != nil {
		return err
	}
	// Check for empty API token
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("api token cannot be empty")
	}
	fmt.Print("Model (optional): ")
	model, err := r.readInput(false)
	if err != nil {
		return err
	}
	account := config.NewAccount(name, apiUrl, apiKey, model)
	err = r.config.AddAccount(*account)
	if err != nil {
		return err
	}
	return nil
}
