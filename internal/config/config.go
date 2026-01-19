// Package config manages user accounts and application configuration,
// allowing creation, listing, modification, and deletion of accounts with API credentials,
// stored persistently in a local JSON configuration file.
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/thxrsxm/harzmind-code/internal/common"
)

// Config represents the configuration file structure.
type Config struct {
	CurrentAccountName string    `json:"currentAccount"`
	Accounts           []Account `json:"accounts"`
}

// GetAccount retrieves an account by name.
func (c *Config) GetAccount(name string) (*Account, error) {
	if c.Accounts == nil {
		return nil, fmt.Errorf("no accounts")
	}
	for i := range c.Accounts {
		if c.Accounts[i].Name == name {
			return &c.Accounts[i], nil
		}
	}
	return nil, fmt.Errorf("account %s not found", name)
}

// GetCurrentAccount retrieves the currently active account.
func (c *Config) GetCurrentAccount() (*Account, error) {
	if len(c.CurrentAccountName) == 0 {
		return nil, fmt.Errorf("no current account")
	}
	return c.GetAccount(c.CurrentAccountName)
}

// AddAccount adds a new account to the configuration.
func (c *Config) AddAccount(account Account) error {
	// Check for existing account to prevent duplicates
	if _, err := c.GetAccount(account.Name); err == nil {
		return fmt.Errorf("account %s already exists", account.Name)
	}
	c.Accounts = append(c.Accounts, account)
	return nil
}

// RemoveAccount removes an account by name.
func (c *Config) RemoveAccount(name string) {
	index := -1
	for i, v := range c.Accounts {
		if v.Name == name {
			index = i
			break
		}
	}
	// Account not found
	if index == -1 {
		return
	}
	// Remove account
	c.Accounts = append(c.Accounts[:index], c.Accounts[index+1:]...)
	if c.CurrentAccountName == name {
		// Logout
		c.CurrentAccountName = ""
	}
}

// SaveConfig saves the configuration to a file.
func (c *Config) SaveConfig(path string) error {
	// Marshal the Config struct to JSON
	jsonData, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	// Get binary path
	binDir, err := common.GetBinaryPath()
	if err != nil {
		return err
	}
	configPath := filepath.Join(binDir, path)
	// Create or truncate the target file
	jsonFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	// Write the JSON data to the file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

// LoadConfig loads the configuration from a file.
func LoadConfig(path string) (*Config, error) {
	// Get binary path
	binDir, err := common.GetBinaryPath()
	if err != nil {
		return nil, err
	}
	configPath := filepath.Join(binDir, path)
	// Read the JSON file
	jsonFile, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	// Read all content from the file
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}
	// Unmarshal the JSON content into a Config struct
	var config Config
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

// CreateConfig creates a new configuration file.
func CreateConfig(path string) error {
	// Create a new Config instance and populate it
	config := Config{
		CurrentAccountName: "",
		Accounts:           []Account{},
	}
	// Marshal the Config struct to JSON
	jsonData, err := json.MarshalIndent(&config, "", "  ")
	if err != nil {
		return err
	}
	// Get binary path
	binDir, err := common.GetBinaryPath()
	if err != nil {
		return err
	}
	configPath := filepath.Join(binDir, path)
	// Create or truncate the target file
	jsonFile, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	// Write the JSON data to the file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}
