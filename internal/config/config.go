// Package config manages user accounts and application configuration,
// allowing creation, listing, modification, and deletion of accounts with API credentials,
// stored persistently in a local JSON configuration file.
package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/thxrsxm/harzmind-code/internal/acc"
)

// Config represents the configuration file structure.
type Config struct {
	path string
	data *configData
}

type configData struct {
	AccountManager acc.AccountManager `json:"accountManagement"`
}

// NewConfig creates a new configuration file.
func NewConfig(path string) (*Config, error) {
	// Create a new Config instance and populate it
	config := &Config{
		path: path,
	}
	config.data = newConfigData(config)
	// Marshal the Config struct to JSON
	jsonData, err := json.MarshalIndent(config.data, "", "  ")
	if err != nil {
		return nil, err
	}
	// Create or truncate the target file
	jsonFile, err := os.Create(config.path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	// Write the JSON data to the file
	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return nil, err
	}
	return config, nil
}

func newConfigData(config *Config) *configData {
	return &configData{
		AccountManager: *acc.NewAccountManager(func() error { return config.SaveConfig() }),
	}
}

// SaveConfig saves the configuration to a file.
func (c *Config) SaveConfig() error {
	// Marshal the Config struct to JSON
	jsonData, err := json.MarshalIndent(c.data, "", "  ")
	if err != nil {
		return err
	}
	// Create or truncate the target file
	jsonFile, err := os.Create(c.path)
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
	// Read the JSON file
	jsonFile, err := os.Open(path)
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
	var data configData
	err = json.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}
	config := Config{
		path: path,
		data: &data,
	}
	data.AccountManager.SetSave(func() error { return config.SaveConfig() })
	return &config, nil
}

func (c *Config) GetAccountManager() *acc.AccountManager {
	return &c.data.AccountManager
}
