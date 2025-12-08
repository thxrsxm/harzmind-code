package config

import "fmt"

type Account struct {
	Name   string `json:"name"`
	ApiUrl string `json:"apiUrl"`
	ApiKey string `json:"apiKey"`
	Model  string `json:"model"`
}

func NewAccount(name, apiUrl, apiKey, model string) *Account {
	return &Account{
		Name:   name,
		ApiUrl: apiUrl,
		ApiKey: apiKey,
		Model:  model,
	}
}

func (a Account) String() string {
	return fmt.Sprintf("Name: %s\nAPI Url: %s\nModel: %s",
		a.Name,
		a.ApiUrl,
		a.Model,
	)
}
