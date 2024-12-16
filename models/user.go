// models/user.go
package models
// User model for example
import (
	"gorm.io/gorm"
)

//
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // '-' means this won't appear in JSON responses
}


