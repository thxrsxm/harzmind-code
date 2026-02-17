// Package config manages user accounts and application configurations.
package config

import (
	"encoding/json"
	"io"
	"os"

	"github.com/thxrsxm/harzmind-code/internal/acc"
)

// Config represents the application's persistent configuration.
// It wraps a configData struct and provides methods to load, save, and mutate configuration.
// Fields are unexported to enforce controlled access via methods.
type Config struct {
	path string
	data *configData
}

// configData holds the serialized configuration structure.
// It is designed for JSON marshaling/unmarshaling and decoupled from runtime state.
type configData struct {
	AccountManager acc.AccountManager `json:"accountManagement"`
}

// NewConfig creates a new configuration file at the given path with default/empty state.
// It initializes an empty account manager and writes the initial JSON structure to disk.
// If the file already exists, it will be truncated and overwritten.
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

// newConfigData constructs a configData instance and initializes its AccountManager.
// The AccountManager is configured with a callback to SaveConfig, ensuring that
// any account modification (add, update, delete) persists changes to disk automatically.
func newConfigData(config *Config) *configData {
	return &configData{
		AccountManager: *acc.NewAccountManager(func() error { return config.SaveConfig() }),
	}
}

// SaveConfig writes the current configuration state to the file path specified during creation.
// It overwrites the file with a new JSON representation, preserving indentation for readability.
// Returns an error if the file cannot be created or written.
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

// LoadConfig reads and deserializes a configuration from the specified file path.
// It validates the JSON structure and initializes the AccountManager with a save hook.
// If the file is missing, unreadable, or contains invalid JSON, an error is returned.
func LoadConfig(path string) (*Config, error) {
	// Open the configuration file
	jsonFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()
	// Read entire file content
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
	// Re-hook the AccountManager's save callback to this config's SaveConfig method.
	data.AccountManager.SetSave(func() error { return config.SaveConfig() })
	return &config, nil
}

// GetAccountManager returns the underlying AccountManager instance.
func (c *Config) GetAccountManager() *acc.AccountManager {
	return &c.data.AccountManager
}
