package handlers

import (
	"bytes"
	"encoding/json"
	"listarr-backend/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Setup test router and any mock dependencies
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	return r
}

func TestGetConfig(t *testing.T) {
	// Setup
	r := setupTestRouter()
	r.GET("/config", GetConfig)

	// Create a mock request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/config", nil)

	// Perform the request
	r.ServeHTTP(w, req)

	// Assert status code
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse response body
	var response models.ConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestUpdateConfig(t *testing.T) {
	// Setup
	r := setupTestRouter()
	r.PUT("/config", UpdateConfig)

	// Create test configuration
	testConfig := models.Configuration{
		// Fill with test data
	}

	// Convert to JSON
	jsonData, _ := json.Marshal(testConfig)

	// Create request with JSON body
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	// Perform request
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)

	var response models.ConfigResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Data)
}

func TestResetConfig(t *testing.T) {
	// Setup
	r := setupTestRouter()
	r.POST("/config/reset", ResetConfig)

	// Create request
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/config/reset", nil)

	// Perform request
	r.ServeHTTP(w, req)

	// Assertions
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetConfig_Error(t *testing.T) {
	// Setup
	r := setupTestRouter()
	r.GET("/config", GetConfig)

	// Mock utils.GetConfig to return nil
	// You'll need to implement mocking here

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/config", nil)

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestUpdateConfig_InvalidJSON(t *testing.T) {
	r := setupTestRouter()
	r.PUT("/config", UpdateConfig)

	// Send invalid JSON
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/config", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateConfig_ValidationCases(t *testing.T) {
	tests := []struct {
		name          string
		config        models.Configuration
		expectedCode  int
		expectedError string
	}{
		{
			name:   "valid config",
			config: models.Configuration{
				// valid configuration
			},
			expectedCode: http.StatusOK,
		},
		{
			name:   "invalid config",
			config: models.Configuration{
				// invalid configuration
			},
			expectedCode:  http.StatusBadRequest,
			expectedError: "Invalid configuration",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupTestRouter()
			r.PUT("/config", UpdateConfig)

			jsonData, _ := json.Marshal(tt.config)
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("PUT", "/config", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")

			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedError != "" {
				var response models.ConfigResponse
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response.Error, tt.expectedError)
			}
		})
	}
}
