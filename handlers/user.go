// handlers/user.go
package handlers


import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

  "listarr-backend/models"
)

// @Summary     Create user
// @Description Create a new user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body User true "User Info"
// @Success     200 {object} User
// @Router      /users [post]
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if err := db.Create(&user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

// @Summary     Get users
// @Description Get all users
// @Tags        users
// @Produce     json
// @Success     200 {array} User
// @Router      /users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// Implement other handlers similarly...
func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}
