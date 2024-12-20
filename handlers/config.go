// handlers/config.go
package handlers

import (
	"encoding/json"
	"listarr-backend/models"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gin-gonic/gin"
)

const (
	configDir  = "./config"
	configFile = "app.config.json"
)

var (
	config     *models.Configuration
	configLock sync.RWMutex
)

// getDefaultConfig returns the default configuration
func getDefaultConfig() models.Configuration {
	return models.Configuration{
		DatabaseURL:   "postgres://localhost:5432/listarr",
		ServerPort:    "8080",
		LogLevel:      "info",
		APIVersion:    "v1",
		MaxPageSize:   100,
		EnableCaching: true,
	}
}

// InitConfig initializes the configuration system
func InitConfig() error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	configPath := filepath.Join(configDir, configFile)
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default configuration if file doesn't exist
		if err := saveConfig(getDefaultConfig()); err != nil {
			return err
		}
	}

	// Load existing configuration
	return loadConfig()
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	return os.MkdirAll(configDir, 0755)
}

// loadConfig loads the configuration from file
func loadConfig() error {
	configLock.Lock()
	defer configLock.Unlock()

	data, err := os.ReadFile(filepath.Join(configDir, configFile))
	if err != nil {
		return err
	}

	var cfg models.Configuration
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	config = &cfg
	return nil
}

// saveConfig saves the configuration to file
func saveConfig(cfg models.Configuration) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(configDir, configFile), data, 0644)
}

// GetConfig godoc
// @Summary Get configuration
// @Description Retrieve current application configuration
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} models.Configuration
// @Failure 500 {object} models.ConfigResponse
// @Router /config [get]
func GetConfig(c *gin.Context) {
	configLock.RLock()
	defer configLock.RUnlock()

	if config == nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Configuration not initialized",
		})
		return
	}

	c.JSON(http.StatusOK, models.ConfigResponse{
		Data: config,
	})
}

// UpdateConfig godoc
// @Summary Update configuration
// @Description Update application configuration settings
// @Tags config
// @Accept json
// @Produce json
// @Param configuration body models.Configuration true "Configuration settings"
// @Success 200 {object} models.Configuration
// @Failure 400 {object} models.ConfigResponse
// @Failure 500 {object} models.ConfigResponse
// @Router /config [put]
func UpdateConfig(c *gin.Context) {
	var newConfig models.Configuration
	if err := c.ShouldBindJSON(&newConfig); err != nil {
		c.JSON(http.StatusBadRequest, models.ConfigResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	configLock.Lock()
	defer configLock.Unlock()

	if err := saveConfig(newConfig); err != nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Failed to save configuration: " + err.Error(),
		})
		return
	}

	config = &newConfig
	c.JSON(http.StatusOK, models.ConfigResponse{
		Data: config,
	})
}

// ResetConfig godoc
// @Summary Reset configuration
// @Description Reset configuration to default values
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} models.Configuration
// @Failure 500 {object} models.ConfigResponse
// @Router /config/reset [post]
func ResetConfig(c *gin.Context) {
	configLock.Lock()
	defer configLock.Unlock()

	if err := ensureConfigDir(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Failed to create config directory: " + err.Error(),
		})
		return
	}

	defaultConfig := getDefaultConfig()
	if err := saveConfig(defaultConfig); err != nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Failed to reset configuration: " + err.Error(),
		})
		return
	}

	config = &defaultConfig
	c.JSON(http.StatusOK, models.ConfigResponse{
		Data: config,
	})
}
