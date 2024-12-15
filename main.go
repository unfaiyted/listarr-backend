 // main.go
package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"

   ginSwagger "github.com/swaggo/gin-swagger"
   swaggerFiles "github.com/swaggo/files"
   _ "listarr-backend/docs"
)

// User model for example
type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"-"` // '-' means this won't appear in JSON responses
}

// @title           Listarr 
// @version         1.0
// @description     API Server for Listarr application. 
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// Initialize DB
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto Migrate the schema
	db.AutoMigrate(&User{})

	// Initialize Gin
	r := gin.Default()

	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000"} // Your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	r.Use(cors.New(config))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Users routes
		users := v1.Group("/users")
		{
			users.POST("", createUser(db))
			users.GET("", getUsers(db))
			users.GET("/:id", getUser(db))
			users.PUT("/:id", updateUser(db))
			users.DELETE("/:id", deleteUser(db))
		}
	}

// Then in your main() function, add:
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// Start server
r.Run(":8080")
}

// @Summary     Create user
// @Description Create a new user
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       user body User true "User Info"
// @Success     200 {object} User
// @Router      /users [post]
func createUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user User
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
func getUsers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var users []User
		if err := db.Find(&users).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, users)
	}
}

// Implement other handlers similarly...
func getUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func updateUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}

func deleteUser(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Implementation
	}
}
