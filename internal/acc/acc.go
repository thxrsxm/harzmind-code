// Package acc manages user accounts for API access, including authentication, configuration,
// session state (login/logout), and interactive account management commands.
package acc

import (
	"fmt"
	"strings"

	"github.com/thxrsxm/harzmind-code/internal/logger"
	"github.com/thxrsxm/harzmind-code/internal/output"
	"github.com/thxrsxm/rnbw"
)

// Account represents a user account with API credentials and model information.
type Account struct {
	Name   string `json:"name"`
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
	Model  string `json:"model"`
}

// NewAccount creates a new Account instance with the given parameters.
func NewAccount(name, apiUrl, apiKey, model string) *Account {
	return &Account{
		Name:   name,
		ApiUrl: apiUrl,
		ApiKey: apiKey,
		Model:  model,
	}
}

// String returns a string representation of the Account instance.
func (a Account) String() string {
	return fmt.Sprintf("Name: %s\nAPI Url: %s\nModel: %s",
		a.Name,
		a.ApiUrl,
		a.Model,
	)
}

// AccountManager manages a collection of accounts and tracks the currently logged-in account.
type AccountManager struct {
	CurrentAccountName string    `json:"currentAccount"`
	Accounts           []Account `json:"accounts"`
	save               func() error
}

// NewAccountManager initializes a new AccountManager with an initial empty account list and no current account.
// The provided `save` function is used internally to persist changes after mutations.
func NewAccountManager(save func() error) *AccountManager {
	return &AccountManager{
		CurrentAccountName: "",
		Accounts:           []Account{},
		save:               save,
	}
}

// SetSave updates the persistent save callback used to commit changes.
func (m *AccountManager) SetSave(save func() error) {
	m.save = save
}

// GetAccount retrieves an account by name.
func (m *AccountManager) GetAccount(name string) (*Account, error) {
	if m.Accounts == nil {
		return nil, fmt.Errorf("no accounts")
	}
	for i := range m.Accounts {
		if m.Accounts[i].Name == name {
			return &m.Accounts[i], nil
		}
	}
	return nil, fmt.Errorf("account %s not found", name)
}

// GetCurrentAccount retrieves the currently active account.
// Returns an error if no account is currently logged in.
func (m *AccountManager) GetCurrentAccount() (*Account, error) {
	if len(m.CurrentAccountName) == 0 {
		return nil, fmt.Errorf("no current account")
	}
	return m.GetAccount(m.CurrentAccountName)
}

// AddAccount adds a new account to the manager's list, avoiding duplicates.
// It returns an error if an account with the same name already exists.
// Persists the updated list via the save callback.
func (m *AccountManager) AddAccount(account Account) error {
	// Check for existing account to prevent duplicates
	if _, err := m.GetAccount(account.Name); err == nil {
		return fmt.Errorf("account %s already exists", account.Name)
	}
	m.Accounts = append(m.Accounts, account)
	return m.save()
}

// RemoveAccount removes an account by name.
// If the removed account is the current one, it automatically logs out (clears CurrentAccountName).
// Persists the updated list via the save callback.
func (m *AccountManager) RemoveAccount(name string) error {
	for i := range m.Accounts {
		if m.Accounts[i].Name == name {
			// Remove account
			m.Accounts = append(m.Accounts[:i], m.Accounts[i+1:]...)
			if m.CurrentAccountName == name {
				// Logout
				m.CurrentAccountName = ""
			}
			break
		}
	}
	return m.save()
}

// Login sets the given account name as the current active account.
// Returns an error if the account does not exist.
// Persists the updated session via the save callback.
func (m *AccountManager) Login(accountName string) error {
	if _, err := m.GetAccount(accountName); err != nil {
		return err
	}
	m.CurrentAccountName = accountName
	return m.save()
}

// Logout clears the current account and returns the name of the account just logged out from.
// Returns an error if no account was currently logged in.
func (m *AccountManager) Logout() (string, error) {
	name := m.CurrentAccountName
	if m.CurrentAccountName == "" {
		return "", fmt.Errorf("not logged in")
	}
	m.CurrentAccountName = ""
	return name, m.save()
}

// PrintAccount prints the details of a named account using the output module.
// Returns an error if the account does not exist.
func (m *AccountManager) PrintAccount(accountName string) error {
	account, err := m.GetAccount(accountName)
	if err != nil {
		return err
	}
	output.Println(account)
	return nil
}

// PrintAllAccounts prints all registered accounts, separated by blank lines.
// If no accounts exist, prints "no accounts".
func (m *AccountManager) PrintAllAccounts() {
	if len(m.Accounts) == 0 {
		output.Println("no accounts")
		return
	}
	for i := range m.Accounts {
		output.Println(m.Accounts[i])
		// Blank line between accounts
		if i < len(m.Accounts)-1 {
			output.Println()
		}
	}
}

// HandleCommands parses and executes account-related commands from a string input.
// Supported single-word commands:
//   - `new`: invokes account creation wizard
//   - `logout`: logs out from current account
//
// Supported two-word commands:
//   - `login <name>`
//   - `remove <name>`
//   - `info <name>`
func (m *AccountManager) HandleCommands(input string) error {
	if len(input) == 0 {
		m.PrintAllAccounts()
		return nil
	}
	args := strings.Split(input, " ")
	if len(args) == 1 {
		switch args[0] {
		case "new":
			// Create a new account via interactive wizard
			account, err := handleAccountCreation()
			if err != nil {
				return err
			}
			if err := m.AddAccount(*account); err != nil {
				return err
			}
			rnbw.ForgroundColor(rnbw.Green)
			output.Printf("\nSuccessfully created the account '%s'\n", account.Name)
			rnbw.ResetColor()
			logger.Log(logger.INFO, "created account '%s'", account.Name)
			return nil
		case "logout":
			// Logout from current account
			name, err := m.Logout()
			if err != nil {
				return err
			}
			rnbw.ForgroundColor(rnbw.Green)
			output.Printf("Successfully logged out from '%s'\n", name)
			rnbw.ResetColor()
			logger.Log(logger.INFO, "logged out from '%s'", name)
			return nil
		default:
			return fmt.Errorf("command not found")
		}
	} else if len(args) == 2 {
		if len(args[1]) == 0 {
			return fmt.Errorf("argument is missing")
		}
		switch args[0] {
		case "login":
			// Login to specified account
			if err := m.Login(args[1]); err != nil {
				return err
			}
			rnbw.ForgroundColor(rnbw.Green)
			output.Printf("Successfully logged in to '%s'\n", args[1])
			rnbw.ResetColor()
			logger.Log(logger.INFO, "logged in to '%s'", args[1])
			return nil
		case "remove":
			// Remove specified account
			if err := m.RemoveAccount(args[1]); err != nil {
				return err
			}
			rnbw.ForgroundColor(rnbw.Green)
			output.Printf("Successfully removed account '%s'\n", args[1])
			rnbw.ResetColor()
			logger.Log(logger.WARNING, "removed account '%s'", args[1])
			return nil
		case "info":
			// Show details of specified account
			if err := m.PrintAccount(args[1]); err != nil {
				return err
			}
			return nil
		default:
			return fmt.Errorf("command not found")
		}
	}
	return fmt.Errorf("command not found")
}
