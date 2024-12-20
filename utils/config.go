// utils/config.go
package utils

import (
	"encoding/json"
	"fmt"
	"listarr-backend/models"
	"os"
	"path/filepath"
)

// ConfigLoader handles loading and validation of configuration files
type ConfigLoader struct {
	configDir  string
	configFile string
}

// NewConfigLoader creates a new configuration loader
func NewConfigLoader(configDir, configFile string) *ConfigLoader {
	return &ConfigLoader{
		configDir:  configDir,
		configFile: configFile,
	}
}

// ValidateConfig checks if the configuration file exists and is valid
func (cl *ConfigLoader) ValidateConfig() error {
	configPath := filepath.Join(cl.configDir, cl.configFile)

	// Check if file exists
	if _, err := os.Stat(configPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("configuration file does not exist: %s", configPath)
		}
		return fmt.Errorf("error checking configuration file: %v", err)
	}

	// Try to read and parse the file to validate JSON
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading configuration file: %v", err)
	}

	var config models.Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("invalid configuration format: %v", err)
	}

	// Additional validation could be added here
	if config.MaxPageSize < 1 || config.MaxPageSize > 1000 {
		return fmt.Errorf("invalid max_page_size: must be between 1 and 1000")
	}

	return nil
}

// RepairConfig attempts to repair or recreate the configuration file
func (cl *ConfigLoader) RepairConfig(defaultConfig models.Configuration) error {
	// Ensure config directory exists
	if err := os.MkdirAll(cl.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Marshal default config to JSON
	data, err := json.MarshalIndent(defaultConfig, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %v", err)
	}

	// Write default config to file
	configPath := filepath.Join(cl.configDir, cl.configFile)
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write default config: %v", err)
	}

	return nil
}

// LoadConfig loads and returns the current configuration
func (cl *ConfigLoader) LoadConfig() (*models.Configuration, error) {
	configPath := filepath.Join(cl.configDir, cl.configFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading configuration: %v", err)
	}

	var config models.Configuration
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing configuration: %v", err)
	}

	return &config, nil
}
