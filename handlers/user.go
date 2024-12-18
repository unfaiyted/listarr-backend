// handlers/user.go
package handlers

import (
	"listarr-backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateUser godoc
//	@Summary		Create a new user
//	@Description	Create a new user in the system
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			user	body		models.User	true	"User data"
//	@Success		201		{object}	models.UserResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users [post]
func CreateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		result := db.Create(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
			return
		}

		c.JSON(http.StatusCreated, user.ToResponse())
	}
}

// GetUsers godoc
//	@Summary		List users
//	@Description	Get all users in the system
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		models.UserResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users [get]
func GetUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []models.User
		result := db.Find(&users)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, models.ErrorResponse{Error: result.Error.Error()})
			return
		}

		userResponses := make([]models.UserResponse, len(users))
		for i, user := range users {
			userResponses[i] = user.ToResponse()
		}

		c.JSON(http.StatusOK, userResponses)
	}
}

// GetUser godoc
//	@Summary		Get a user
//	@Description	Get a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		200	{object}	models.UserResponse
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/{id} [get]
func GetUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
			return
		}
		c.JSON(http.StatusOK, user.ToResponse())
	}
}

// UpdateUser godoc
//	@Summary		Update a user
//	@Description	Update a user's information
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int			true	"User ID"
//	@Param			user	body		models.User	true	"User data"
//	@Success		200		{object}	models.UserResponse
//	@Failure		400		{object}	models.ErrorResponse
//	@Failure		404		{object}	models.ErrorResponse
//	@Failure		500		{object}	models.ErrorResponse
//	@Router			/users/{id} [put]
func UpdateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
			return
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{Error: err.Error()})
			return
		}

		db.Save(&user)
		c.JSON(http.StatusOK, user.ToResponse())
	}
}

// DeleteUser godoc
//	@Summary		Delete a user
//	@Description	Delete a user by ID
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"User ID"
//	@Success		204	{object}	nil
//	@Failure		404	{object}	models.ErrorResponse
//	@Failure		500	{object}	models.ErrorResponse
//	@Router			/users/{id} [delete]
func DeleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := db.First(&user, c.Param("id")).Error; err != nil {
			c.JSON(http.StatusNotFound, models.ErrorResponse{Error: "User not found"})
			return
		}

		db.Delete(&user)
		c.Status(http.StatusNoContent)
	}
}
