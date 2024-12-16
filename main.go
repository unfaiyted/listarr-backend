 // main.go
package main

import (
	"log"
  "fmt"
  "os"
	// "net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"gorm.io/gorm"
	"gorm.io/driver/postgres"

   ginSwagger "github.com/swaggo/gin-swagger"
   swaggerFiles "github.com/swaggo/files"
   _ "listarr-backend/docs"

   "listarr-backend/models"
   "listarr-backend/handlers"

)

// @title           Listarr 
// @version         1.0
// @description     API Server for Listarr application. 
// @host            localhost:8080
// @BasePath        /api/v1
func main() {
	// Initialize DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	// Auto Migrate the schema
	db.AutoMigrate(&models.User{})

	// Initialize Gin
	r := gin.Default()

	// CORS Configuration
	config := cors.DefaultConfig()
  config.AllowOrigins = []string{"http://localhost:3000", "http://192.168.0.126:3000"} // Your frontend URL
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Authorization", "Content-Type"}
	r.Use(cors.New(config))

	// API v1 routes
	v1 := r.Group("/api/v1")
	{
		// Users routes
		users := v1.Group("/users")
		{
			users.POST("", handlers.CreateUser(db))
			users.GET("", handlers.GetUsers(db))
			users.GET("/:id", handlers.GetUser(db))
			users.PUT("/:id", handlers.UpdateUser(db))
			users.DELETE("/:id", handlers.DeleteUser(db))
		}
	}

// Then in your main() function, add:
r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

// Start server
r.Run(":8080")
}
