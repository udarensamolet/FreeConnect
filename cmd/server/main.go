package main

import (
	"FreeConnect/internal/config"
	"log"

	"FreeConnect/internal/controllers"
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration.
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to PostgreSQL using DSN from config.
	db, err := models.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Auto-migrate the User model.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Initialize repository, service, and controller.
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Setup Gin router and register user routes.
	router := gin.Default()
	api := router.Group("/api")
	{
		api.POST("/register", userController.RegisterUser)
		api.GET("/users/:id", userController.GetUser)
	}

	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
