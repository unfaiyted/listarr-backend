// models/responses.go
package models

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message,omitempty" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}
