package acc

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/common"
	"github.com/thxrsxm/harzmind-code/internal/input"
)

// handleAccountCreation prompts the user interactively to input account details.
// It reads: name, API URL, API token (securely), and optionally model.
// Returns a pointer to a new Account, or an error if validation fails.
//
// NOTE: The wizard expects valid input: non-empty name/token, valid URL.
func handleAccountCreation() (*Account, error) {
	fmt.Println("Create account")
	// Read and validate account name
	fmt.Print("Name: ")
	name, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	// Check for empty name
	if strings.TrimSpace(name) == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}
	// Read and validate API URL
	fmt.Print("API Url: ")
	apiURL, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	// Validate API URL
	if !common.IsValidURL(apiURL) {
		return nil, fmt.Errorf("invalid api url")
	}
	// Securely read API key/token (no echo)
	fmt.Print("API Token: ")
	apiKey, err := input.ReadPassword()
	if err != nil {
		return nil, err
	}
	// Check for empty API token
	if strings.TrimSpace(apiKey) == "" {
		return nil, fmt.Errorf("api token cannot be empty")
	}
	// Optionally read model (defaults to empty)
	fmt.Print("Model (optional): ")
	model, err := input.ReadInput(false)
	if err != nil {
		return nil, err
	}
	// Build and return the new account
	account := NewAccount(name, apiURL, apiKey, model)
	return account, nil
}
