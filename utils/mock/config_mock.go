package mock

import (
	"listarr-backend/models"
)

// mock_utils/mock_config.go
type MockConfigUtils interface {
	GetConfig() *models.Configuration
	GetFileConfig() *models.Configuration
	SaveFileConfig(config models.Configuration) error
	ResetFileConfig() error
}
