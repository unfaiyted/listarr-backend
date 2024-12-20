// models/config.go
package models

// Configuration represents the application configuration structure
type Configuration struct {
	DatabaseURL   string `json:"database_url" example:"postgres://localhost:5432/listarr" binding:"required"`
	ServerPort    string `json:"server_port" example:"8080" binding:"required"`
	LogLevel      string `json:"log_level" example:"info" binding:"required,oneof=debug info warn error"`
	APIVersion    string `json:"api_version" example:"v1" binding:"required"`
	MaxPageSize   int    `json:"max_page_size" example:"100" binding:"required,min=1,max=1000"`
	EnableCaching bool   `json:"enable_caching" example:"true"`
}

// ConfigResponse represents the response structure for configuration endpoints
type ConfigResponse struct {
	Data  *Configuration `json:"data,omitempty"`
	Error string         `json:"error,omitempty"`
}
