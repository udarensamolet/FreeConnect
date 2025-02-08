package main

import (
	"log"

	"FreeConnect/internal/config"
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

	// Connect to PostgreSQL.
	db, err := models.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Auto-migrate the User and Project models.
	if err := db.AutoMigrate(&models.User{}); err != nil {
		log.Fatalf("User migration failed: %v", err)
	}
	if err := db.AutoMigrate(&models.Project{}); err != nil {
		log.Fatalf("Project migration failed: %v", err)
	}

	// Initialize repositories, services, and controllers.
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectController := controllers.NewProjectController(projectService)

	// Setup Gin router and register routes.
	router := gin.Default()
	api := router.Group("/api")
	{
		// User routes.
		api.POST("/register", userController.RegisterUser)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)

		// Project routes.
		api.POST("/projects", projectController.CreateProject)
		api.GET("/projects", projectController.GetAllProjects)
		api.GET("/projects/:id", projectController.GetProject)
		api.PUT("/projects/:id", projectController.UpdateProject)
		api.DELETE("/projects/:id", projectController.DeleteProject)
	}

	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
