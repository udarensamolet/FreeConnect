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

	// Auto-migrate all models.
	if err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Skill{},
		&models.Proposal{},
		&models.Review{},
		&models.Transaction{},
		&models.Task{},
		&models.Notification{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Initialize repositories, services, and controllers for each module.

	// Users.
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Projects.
	projectRepo := repositories.NewProjectRepository(db)
	projectService := services.NewProjectService(projectRepo)
	projectController := controllers.NewProjectController(projectService)

	// Skills.
	skillRepo := repositories.NewSkillRepository(db)
	skillService := services.NewSkillService(skillRepo)
	skillController := controllers.NewSkillController(skillService)

	// Proposals.
	proposalRepo := repositories.NewProposalRepository(db)
	proposalService := services.NewProposalService(proposalRepo)
	proposalController := controllers.NewProposalController(proposalService)

	// Reviews.
	reviewRepo := repositories.NewReviewRepository(db)
	reviewService := services.NewReviewService(reviewRepo)
	reviewController := controllers.NewReviewController(reviewService)

	// Transactions.
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionController := controllers.NewTransactionController(transactionService)

	// Tasks.
	taskRepo := repositories.NewTaskRepository(db)
	taskService := services.NewTaskService(taskRepo)
	taskController := controllers.NewTaskController(taskService)

	// Notifications.
	notificationRepo := repositories.NewNotificationRepository(db)
	notificationService := services.NewNotificationService(notificationRepo)
	notificationController := controllers.NewNotificationController(notificationService)

	// Setup Gin router and register routes.
	router := gin.Default()
	api := router.Group("/api")
	{
		// User routes.
		api.POST("/register", userController.RegisterUser)
		api.GET("/users/:id", userController.GetUser)
		api.PUT("/users/:id", userController.UpdateUser)
		api.PUT("/users/:id/skills", userController.UpdateUserSkills)

		// Project routes.
		api.POST("/projects", projectController.CreateProject)
		api.GET("/projects", projectController.GetAllProjects)
		api.GET("/projects/:id", projectController.GetProject)
		api.PUT("/projects/:id", projectController.UpdateProject)
		api.DELETE("/projects/:id", projectController.DeleteProject)

		// Skill routes.
		api.POST("/skills", skillController.CreateSkill)
		api.GET("/skills", skillController.GetAllSkills)
		api.GET("/skills/:id", skillController.GetSkill)
		api.PUT("/skills/:id", skillController.UpdateSkill)
		api.DELETE("/skills/:id", skillController.DeleteSkill)

		// Proposal routes.
		api.POST("/proposals", proposalController.CreateProposal)
		api.GET("/proposals/:id", proposalController.GetProposal)
		api.GET("/projects/:id/proposals", proposalController.GetProposalsByProject)
		api.PUT("/proposals/:id", proposalController.UpdateProposal)
		api.DELETE("/proposals/:id", proposalController.DeleteProposal)

		// Review routes.
		api.POST("/reviews", reviewController.CreateReview)
		api.GET("/reviews/:id", reviewController.GetReview)
		api.GET("/projects/:id/reviews", reviewController.GetReviewsByProject)
		api.PUT("/reviews/:id", reviewController.UpdateReview)
		api.DELETE("/reviews/:id", reviewController.DeleteReview)

		// Transaction routes.
		api.POST("/transactions", transactionController.CreateTransaction)
		api.GET("/transactions/:id", transactionController.GetTransaction)
		api.GET("/projects/:id/transactions", transactionController.GetTransactionsByProject)
		api.PUT("/transactions/:id", transactionController.UpdateTransaction)
		api.DELETE("/transactions/:id", transactionController.DeleteTransaction)

		// Task routes.
		api.POST("/tasks", taskController.CreateTask)
		api.GET("/tasks/:id", taskController.GetTask)
		api.GET("/projects/:id/tasks", taskController.GetTasksByProject)
		api.PUT("/tasks/:id", taskController.UpdateTask)
		api.DELETE("/tasks/:id", taskController.DeleteTask)

		// Notification routes.
		api.POST("/notifications", notificationController.CreateNotification)
		api.GET("/notifications/:id", notificationController.GetNotification)
		api.GET("/notifications/user/:user_id", notificationController.GetNotificationsByUser)
		api.PUT("/notifications/:id", notificationController.UpdateNotification)
		api.DELETE("/notifications/:id", notificationController.DeleteNotification)
	}

	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
