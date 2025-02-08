package main

import (
	"log"

	"FreeConnect/internal/config"
	"FreeConnect/internal/controllers"
	"FreeConnect/internal/middleware"
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"

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

	// ---------- NEW AUTH SETUP ----------
	jwtService := services.NewJWTService()
	authController := controllers.NewAuthController(userService, jwtService)
	// Admin controller
	adminController := controllers.NewAdminController(userRepo)
	// Real-time SSE controller
	rtc := controllers.NewRealTimeController()

	// Setup Gin
	router := gin.Default()
	api := router.Group("/api")

	// PUBLIC routes ------------------
	api.POST("/register", userController.RegisterUser)
	api.POST("/login", authController.Login)

	// (OPTIONAL) SSE
	rtc.RegisterRoutes(api)

	// PROTECTED routes (must have valid JWT) ---
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtService)) // requires token
	{
		// Example: only "client" can create projects
		clientOnly := protected.Group("/client")
		clientOnly.Use(middleware.RoleMiddleware("client"))
		{
			// Project routes
			clientOnly.POST("/projects", projectController.CreateProject)
			// More client-only actions...
		}

		// Example: only "freelancer" can create proposals
		freelancerOnly := protected.Group("/freelancer")
		freelancerOnly.Use(middleware.RoleMiddleware("freelancer"))
		{
			freelancerOnly.POST("/proposals", proposalController.CreateProposal)
			// More freelancer-only endpoints...
		}

		// Example: Admin routes
		adminOnly := protected.Group("/admin")
		adminOnly.Use(middleware.RoleMiddleware("admin"))
		{
			adminOnly.GET("/users", adminController.ListAllUsers)
			adminOnly.PUT("/users/:id/approve", adminController.ApproveUser)
			// etc...
		}

		// ---------- Common routes (any logged-in user) ----------
		// e.g. user updates
		protected.GET("/users/:id", userController.GetUser)
		protected.PUT("/users/:id", userController.UpdateUser)
		protected.PUT("/users/:id/skills", userController.UpdateUserSkills)

		// Project (some might require role checks if you want)
		protected.GET("/projects", projectController.GetAllProjects)
		protected.GET("/projects/:id", projectController.GetProject)
		protected.PUT("/projects/:id", projectController.UpdateProject)
		protected.DELETE("/projects/:id", projectController.DeleteProject)

		// Skills
		protected.POST("/skills", skillController.CreateSkill)
		protected.GET("/skills", skillController.GetAllSkills)
		protected.GET("/skills/:id", skillController.GetSkill)
		protected.PUT("/skills/:id", skillController.UpdateSkill)
		protected.DELETE("/skills/:id", skillController.DeleteSkill)

		// Proposals
		protected.GET("/proposals/:id", proposalController.GetProposal)
		protected.GET("/projects/:id/proposals", proposalController.GetProposalsByProject)
		protected.PUT("/proposals/:id", proposalController.UpdateProposal)
		protected.DELETE("/proposals/:id", proposalController.DeleteProposal)
		// Accept proposal
		protected.POST("/proposals/:id/accept", proposalController.AcceptProposal)

		// Reviews
		protected.POST("/reviews", reviewController.CreateReview)
		protected.GET("/reviews/:id", reviewController.GetReview)
		protected.GET("/projects/:id/reviews", reviewController.GetReviewsByProject)
		protected.PUT("/reviews/:id", reviewController.UpdateReview)
		protected.DELETE("/reviews/:id", reviewController.DeleteReview)

		// Transactions
		protected.POST("/transactions", transactionController.CreateTransaction)
		protected.GET("/transactions/:id", transactionController.GetTransaction)
		protected.GET("/projects/:id/transactions", transactionController.GetTransactionsByProject)
		protected.PUT("/transactions/:id", transactionController.UpdateTransaction)
		protected.DELETE("/transactions/:id", transactionController.DeleteTransaction)

		// Tasks
		protected.POST("/tasks", taskController.CreateTask)
		protected.GET("/tasks/:id", taskController.GetTask)
		protected.GET("/projects/:id/tasks", taskController.GetTasksByProject)
		protected.PUT("/tasks/:id", taskController.UpdateTask)
		protected.DELETE("/tasks/:id", taskController.DeleteTask)

		// Notifications
		protected.POST("/notifications", notificationController.CreateNotification)
		protected.GET("/notifications/:id", notificationController.GetNotification)
		protected.GET("/notifications/user/:user_id", notificationController.GetNotificationsByUser)
		protected.PUT("/notifications/:id", notificationController.UpdateNotification)
		protected.DELETE("/notifications/:id", notificationController.DeleteNotification)

		// Invoices
		protected.POST("/invoices", invoiceController.CreateInvoice)
		protected.GET("/invoices/:id", invoiceController.GetInvoice)
		protected.GET("/projects/:id/invoices", invoiceController.GetInvoicesByProject)
		protected.PUT("/invoices/:id", invoiceController.UpdateInvoice)
		protected.DELETE("/invoices/:id", invoiceController.DeleteInvoice)
	}

	// Run server
	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
