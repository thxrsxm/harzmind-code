package config

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	API     string   `yaml:"api"`
	Model   string   `yaml:"model"`
	Outfile bool     `yaml:"outfile"`
	Ignore  []string `yaml:"ignore"`
}

// String implements the fmt.Stringer interface for the Config struct.
func (c Config) String() string {
	return fmt.Sprintf("API: %q\nModel: %q\nOutfile: %t\nIgnore: %v",
		c.API,
		c.Model,
		c.Outfile,
		c.Ignore)
}

func LoadConfig(path string) (*Config, error) {
	// Read the YAML file
	yamlFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer yamlFile.Close()
	// Read all content from the file
	byteValue, err := io.ReadAll(yamlFile)
	if err != nil {
		return nil, err
	}
	// Unmarshal the YAML content into a Config struct
	var config Config
	err = yaml.Unmarshal(byteValue, &config)
	if err != nil {
		return nil, err
	}
	if len(config.API) == 0 {
		return nil, fmt.Errorf("API field is empty")
	}
	if len(config.Model) == 0 {
		return nil, fmt.Errorf("model field is empty")
	}
	return &config, nil
}

func (c *Config) Save(path string) error {
	yamlData, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(path, yamlData, 0644)
}

func CreateConfig(path string) error {
	// Create a new Config instance and populate it
	config := Config{
		API:   "https://api.openai.com/v1",
		Model: "gpt-4o",
		Ignore: []string{
			"hzmind",
		},
	}
	// Marshal the Config struct to YAML
	yamlData, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	// Create or truncate the target file
	yamlFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer yamlFile.Close()
	// Write the YAML data to the file
	_, err = yamlFile.Write(yamlData)
	if err != nil {
		return err
	}
	return nil
}
