package handlers

import (
	"fmt"
	"listarr-backend/models"
	"listarr-backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetConfig godoc
// @Summary Get configuration
// @Description Retrieve current application configuration
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} models.ConfigResponse
// @Failure 500 {object} models.ConfigResponse
// @Router /config [get]
func GetConfig(c *gin.Context) {
	currentConfig := utils.GetConfig()
	if currentConfig == nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Configuration not initialized",
		})
		return
	}

	// Only return the file-based configuration, not environment overrides
	// fileConfig := utils.GetFileConfig()
	c.JSON(http.StatusOK, models.ConfigResponse{
		Data: currentConfig,
	})
}

// UpdateConfig godoc
// @Summary Update configuration
// @Description Update application configuration settings in app.config.json
// @Tags config
// @Accept json
// @Produce json
// @Param configuration body models.Configuration true "Configuration settings"
// @Success 200 {object} models.ConfigResponse
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

	// Validate the new configuration
	if err := validateConfig(newConfig); err != nil {
		c.JSON(http.StatusBadRequest, models.ConfigResponse{
			Error: "Invalid configuration: " + err.Error(),
		})
		return
	}

	// Save only to app.config.json
	if err := utils.SaveFileConfig(newConfig); err != nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Failed to save configuration: " + err.Error(),
		})
		return
	}

	// Return the file-based configuration
	c.JSON(http.StatusOK, models.ConfigResponse{
		Data: utils.GetFileConfig(),
	})
}

// ResetConfig godoc
// @Summary Reset configuration
// @Description Reset app.config.json to default values
// @Tags config
// @Accept json
// @Produce json
// @Success 200 {object} models.ConfigResponse
// @Failure 500 {object} models.ConfigResponse
// @Router /config/reset [post]
func ResetConfig(c *gin.Context) {
	if err := utils.ResetFileConfig(); err != nil {
		c.JSON(http.StatusInternalServerError, models.ConfigResponse{
			Error: "Failed to reset configuration: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, models.ConfigResponse{
		Data: utils.GetFileConfig(),
	})
}

// validateConfig performs basic validation of configuration values
func validateConfig(cfg models.Configuration) error {
	// Add your validation logic here
	fmt.Println(cfg)
	// For example:
	// - Check required fields
	// - Validate value ranges
	// - Check format of URLs, etc.
	return nil
}
