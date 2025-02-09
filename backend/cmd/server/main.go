package main

import (
	"log"

	"FreeConnect/internal/config"
	"FreeConnect/internal/controllers"
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to PostgreSQL
	db, err := models.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// Auto-migrate all models
	if err := db.AutoMigrate(
		&models.User{},
		&models.Project{},
		&models.Skill{},
		&models.Proposal{},
		&models.Review{},
		&models.Transaction{},
		&models.Task{},
		&models.Notification{},
		&models.Invoice{},
	); err != nil {
		log.Fatalf("AutoMigrate failed: %v", err)
	}

	// Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	skillRepo := repositories.NewSkillRepository(db)
	proposalRepo := repositories.NewProposalRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)

	// Initialize services
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	skillService := services.NewSkillService(skillRepo)
	proposalService := services.NewProposalService(proposalRepo)
	reviewService := services.NewReviewService(reviewRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	taskService := services.NewTaskService(taskRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// Initialize controllers
	userController := controllers.NewUserController(userService)
	projectController := controllers.NewProjectController(projectService)
	skillController := controllers.NewSkillController(skillService)
	proposalController := controllers.NewProposalController(proposalService)
	reviewController := controllers.NewReviewController(reviewService)
	transactionController := controllers.NewTransactionController(transactionService)
	taskController := controllers.NewTaskController(taskService)
	notificationController := controllers.NewNotificationController(notificationService)
	invoiceController := controllers.NewInvoiceController(invoiceService)

	// NEW AUTH SETUP
	jwtService := services.NewJWTService()
	authController := controllers.NewAuthController(userService, jwtService)

	// Admin controller
	adminController := controllers.NewAdminController(userRepo)

	// Real-time SSE controller
	rtc := controllers.NewRealTimeController()

	// Setup Gin
	router := gin.Default()

	// 1) USE CORS MIDDLEWARE
	//    Adjust origins, methods, headers as needed.
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true, // If needed
	}))

	// 2) Define your routes
	api := router.Group("/api")

	// PUBLIC routes
	api.POST("/register", userController.RegisterUser)
	api.POST("/login", authController.Login)

	// (OPTIONAL) SSE
	rtc.RegisterRoutes(api)

	// Examples of other routes
	api.POST("/projects", projectController.CreateProject)
	api.POST("/proposals", proposalController.CreateProposal)
	api.GET("/users", adminController.ListAllUsers)
	api.PUT("/users/:id/approve", adminController.ApproveUser)

	api.GET("/users/:id", userController.GetUser)
	api.PUT("/users/:id", userController.UpdateUser)
	api.PUT("/users/:id/skills", userController.UpdateUserSkills)

	api.GET("/projects", projectController.GetAllProjects)
	api.GET("/projects/:id", projectController.GetProject)
	api.PUT("/projects/:id", projectController.UpdateProject)
	api.DELETE("/projects/:id", projectController.DeleteProject)

	api.POST("/skills", skillController.CreateSkill)
	api.GET("/skills", skillController.GetAllSkills)
	api.GET("/skills/:id", skillController.GetSkill)
	api.PUT("/skills/:id", skillController.UpdateSkill)
	api.DELETE("/skills/:id", skillController.DeleteSkill)

	api.GET("/proposals/:id", proposalController.GetProposal)
	api.GET("/projects/:id/proposals", proposalController.GetProposalsByProject)
	api.PUT("/proposals/:id", proposalController.UpdateProposal)
	api.DELETE("/proposals/:id", proposalController.DeleteProposal)
	api.POST("/proposals/:id/accept", proposalController.AcceptProposal)

	api.POST("/reviews", reviewController.CreateReview)
	api.GET("/reviews/:id", reviewController.GetReview)
	api.GET("/projects/:id/reviews", reviewController.GetReviewsByProject)
	api.PUT("/reviews/:id", reviewController.UpdateReview)
	api.DELETE("/reviews/:id", reviewController.DeleteReview)

	api.POST("/transactions", transactionController.CreateTransaction)
	api.GET("/transactions/:id", transactionController.GetTransaction)
	api.GET("/projects/:id/transactions", transactionController.GetTransactionsByProject)
	api.PUT("/transactions/:id", transactionController.UpdateTransaction)
	api.DELETE("/transactions/:id", transactionController.DeleteTransaction)

	api.POST("/tasks", taskController.CreateTask)
	api.GET("/tasks/:id", taskController.GetTask)
	api.GET("/projects/:id/tasks", taskController.GetTasksByProject)
	api.PUT("/tasks/:id", taskController.UpdateTask)
	api.DELETE("/tasks/:id", taskController.DeleteTask)

	api.POST("/notifications", notificationController.CreateNotification)
	api.GET("/notifications/:id", notificationController.GetNotification)
	api.GET("/notifications/user/:user_id", notificationController.GetNotificationsByUser)
	api.PUT("/notifications/:id", notificationController.UpdateNotification)
	api.DELETE("/notifications/:id", notificationController.DeleteNotification)

	api.POST("/invoices", invoiceController.CreateInvoice)
	api.GET("/invoices/:id", invoiceController.GetInvoice)
	api.GET("/projects/:id/invoices", invoiceController.GetInvoicesByProject)
	api.PUT("/invoices/:id", invoiceController.UpdateInvoice)
	api.DELETE("/invoices/:id", invoiceController.DeleteInvoice)

	// Run server
	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
