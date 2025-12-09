// Package config provides functionality for handling accounts.
package config

import "fmt"

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
