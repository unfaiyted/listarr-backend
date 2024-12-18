// models/user.go
package models

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user model
type User struct {
	ID       uint   `json:"id" gorm:"primarykey" example:"1"`
	Name     string `json:"name" gorm:"not null" example:"John Doe" binding:"required" minLength:"2" maxLength:"100"`
	Email    string `json:"email" gorm:"uniqueIndex;not null" example:"john@example.com" binding:"required,email"`
	Password string `json:"password,omitempty" gorm:"not null" example:"strongpassword123" binding:"required,min=8" swaggertype:"string" format:"password"` // omitempty will exclude it from JSON responses
}

// BeforeSave hook to hash password before saving to database
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// For API responses, we want to exclude the password
type UserResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// ToResponse converts User to UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
