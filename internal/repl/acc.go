package repl

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal"
	"github.com/thxrsxm/harzmind-code/internal/config"
)

// handleAccountCreation handles the creation of a new account.
// It prompts the user for account details and validates the input.
// If successful, it adds the new account to the configuration and saves it.
func (r *REPL) handleAccountCreation() error {
	r.out.Println("Create account")
	// Get account name
	r.out.Print("Name: ")
	name, err := r.readInput()
	if err != nil {
		return err
	}
	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return fmt.Errorf("name cannot be empty")
	}
	// Get API url
	r.out.Print("API Url: ")
	apiUrl, err := r.readInput()
	if err != nil {
		return err
	}
	// Check for empty API URL
	if strings.TrimSpace(apiUrl) == "" {
		return fmt.Errorf("api url cannot be empty")
	}
	// Validate API URL
	if !internal.IsValidURL(apiUrl) {
		return fmt.Errorf("invalid api url")
	}
	// Get API token (secure)
	r.out.Print("API Token: ")
	apiKey, err := r.readPassword()
	if err != nil {
		return err
	}
	// Check for empty API token
	if strings.TrimSpace(apiKey) == "" {
		return fmt.Errorf("api token cannot be empty")
	}
	r.out.Print("Model (optional): ")
	model, err := r.readInput()
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
