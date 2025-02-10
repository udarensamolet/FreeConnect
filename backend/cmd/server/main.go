package main

import (
	"log"

	"FreeConnect/internal/config"
	"FreeConnect/internal/controllers"
	"FreeConnect/internal/middleware"
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 1) Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 2) Connect to PostgreSQL
	db, err := models.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	// 3) Auto-migrate models
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

	// 4) Initialize repositories
	userRepo := repositories.NewUserRepository(db)
	projectRepo := repositories.NewProjectRepository(db)
	skillRepo := repositories.NewSkillRepository(db)
	proposalRepo := repositories.NewProposalRepository(db)
	reviewRepo := repositories.NewReviewRepository(db)
	transactionRepo := repositories.NewTransactionRepository(db)
	taskRepo := repositories.NewTaskRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	invoiceRepo := repositories.NewInvoiceRepository(db)

	// 5) Initialize services
	userService := services.NewUserService(userRepo)
	projectService := services.NewProjectService(projectRepo)
	skillService := services.NewSkillService(skillRepo)
	proposalService := services.NewProposalService(proposalRepo)
	reviewService := services.NewReviewService(reviewRepo)
	transactionService := services.NewTransactionService(transactionRepo)
	taskService := services.NewTaskService(taskRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// 6) Initialize controllers
	userController := controllers.NewUserController(userService)
	projectController := controllers.NewProjectController(projectService)
	skillController := controllers.NewSkillController(skillService)
	proposalController := controllers.NewProposalController(proposalService)
	reviewController := controllers.NewReviewController(reviewService)
	transactionController := controllers.NewTransactionController(transactionService)
	taskController := controllers.NewTaskController(taskService)
	notificationController := controllers.NewNotificationController(notificationService)
	invoiceController := controllers.NewInvoiceController(invoiceService)

	// Auth & Admin controllers
	jwtService := services.NewJWTService()
	authController := controllers.NewAuthController(userService, jwtService)
	adminController := controllers.NewAdminController(userRepo)

	// Real-time SSE controller
	rtc := controllers.NewRealTimeController()

	// 7) Setup Gin + CORS
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	//--------------------------------------------------------------------
	// PUBLIC ROUTES (No auth required)
	//--------------------------------------------------------------------
	public := router.Group("/api")
	{
		// 1) Registration & Login
		public.POST("/register", userController.RegisterUser)
		public.POST("/login", authController.Login)

		// 2) Browse / View Projects (Anonymous can see them)
		public.GET("/projects", projectController.GetAllProjects)
		public.GET("/projects/:id", projectController.GetProject)

		// 3) Skills: (If desired, anonymous can see skill listings)
		public.GET("/skills", skillController.GetAllSkills)
		public.GET("/skills/:id", skillController.GetSkill)

		// 4) Real-time SSE (optional if you want it public)
		rtc.RegisterRoutes(public)
	}

	//--------------------------------------------------------------------
	// PROTECTED ROUTES (Require a valid JWT)
	//--------------------------------------------------------------------
	// The AuthMiddleware ensures any request here has a valid token,
	// so c.Get("userRole") and c.Get("userID") will be set if needed.
	secure := router.Group("/api")
	secure.Use(middleware.AuthMiddleware(jwtService))
	{
		// ---------------- ADMIN EXAMPLES ----------------
		// (If you want a stricter admin-only approach, use RoleMiddleware("admin"))
		secure.GET("/users", adminController.ListAllUsers)
		secure.PUT("/users/:id/approve", adminController.ApproveUser)

		secure.GET("/users/:id", userController.GetUser)
		secure.PUT("/users/:id", userController.UpdateUser)
		secure.PUT("/users/:id/skills", userController.UpdateUserSkills)

		// ---------------- PROJECTS ----------------
		// NOTE: Now creation does NOT require a freelancer_id.
		secure.POST("/projects", projectController.CreateProject)
		secure.PUT("/projects/:projectId", projectController.UpdateProject)
		secure.DELETE("/projects/:id", projectController.DeleteProject)

		// ADDITIONAL: route for setting the freelancer
		// e.g. POST /api/projects/:id/set-freelancer
		secure.POST("/projects/:id/set-freelancer", projectController.SetProjectFreelancer)

		// ---------------- PROPOSALS ----------------
		secure.POST("/proposals", proposalController.CreateProposal)
		secure.GET("/proposals/:id", proposalController.GetProposal)
		secure.GET("/projects/:id/proposals", proposalController.GetProposalsByProject)
		secure.PUT("/proposals/:id", proposalController.UpdateProposal)
		secure.DELETE("/proposals/:id", proposalController.DeleteProposal)
		secure.POST("/proposals/:id/accept", proposalController.AcceptProposal)

		// ---------------- REVIEWS ----------------
		secure.POST("/reviews", reviewController.CreateReview)
		secure.GET("/reviews/:id", reviewController.GetReview)
		secure.GET("/projects/:id/reviews", reviewController.GetReviewsByProject)
		secure.PUT("/reviews/:id", reviewController.UpdateReview)
		secure.DELETE("/reviews/:id", reviewController.DeleteReview)

		// ---------------- TRANSACTIONS ----------------
		secure.POST("/transactions", transactionController.CreateTransaction)
		secure.GET("/transactions/:id", transactionController.GetTransaction)
		secure.GET("/projects/:id/transactions", transactionController.GetTransactionsByProject)
		secure.PUT("/transactions/:id", transactionController.UpdateTransaction)
		secure.DELETE("/transactions/:id", transactionController.DeleteTransaction)

		// ---------------- TASKS ----------------
		secure.POST("/tasks", taskController.CreateTask)
		secure.GET("/tasks/:id", taskController.GetTask)
		secure.GET("/projects/:id/tasks", taskController.GetTasksByProject)
		secure.PUT("/tasks/:id", taskController.UpdateTask)
		secure.DELETE("/tasks/:id", taskController.DeleteTask)
		secure.PUT("/projects/:projectId/tasks/:taskId/edit", taskController.EditTask)
		// ---------------- NOTIFICATIONS ----------------
		secure.POST("/notifications", notificationController.CreateNotification)
		secure.GET("/notifications/:id", notificationController.GetNotification)
		secure.GET("/notifications/user/:user_id", notificationController.GetNotificationsByUser)
		secure.PUT("/notifications/:id", notificationController.UpdateNotification)
		secure.DELETE("/notifications/:id", notificationController.DeleteNotification)

		// ---------------- INVOICES ----------------
		secure.POST("/invoices", invoiceController.CreateInvoice)
		secure.GET("/invoices/:id", invoiceController.GetInvoice)
		secure.GET("/projects/:id/invoices", invoiceController.GetInvoicesByProject)
		secure.PUT("/invoices/:id", invoiceController.UpdateInvoice)
		secure.DELETE("/invoices/:id", invoiceController.DeleteInvoice)
	}

	//--------------------------------------------------------------------
	// START THE SERVER
	//--------------------------------------------------------------------
	log.Printf("Server running on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to run: %v", err)
	}
}
