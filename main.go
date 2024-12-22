// main.go
package main

import (
	"fmt"
	"listarr-backend/handlers"
	"listarr-backend/models"
	"listarr-backend/utils"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	_ "listarr-backend/docs"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

//	@title			Listarr API
//	@version		1.0
//	@description	API Server for Listarr application
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		localhost:8080
// @BasePath	/api/v1
// @schemes	http
// @openapi	3.0.0
func main() {
	if err := utils.InitConfig(); err != nil {
		log.Fatalf("Failed to initilize conifg: %v", err)
	}

	appConfig := utils.GetConfig()

	// Initialize DB
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		appConfig.Db.Host,
		appConfig.Db.User,
		appConfig.Db.Password,
		appConfig.Db.Name,
		appConfig.Db.Port)
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
	config.AllowOrigins = []string{"http://localhost:5173", "http://localhost:5174", "http://localhost:3000", "http://192.168.0.126:3000"} // Your frontend URL
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
		v1.GET("/config", handlers.GetConfig)
		v1.PUT("/config", handlers.UpdateConfig)
		v1.POST("/config/reset", handlers.ResetConfig)

	}

	// Then in your main() function, add:
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Start server
	r.Run(":8080")
}
