package controllers

import (
	"net/http" // For HTTP status codes and response writing.
	"strconv"  // For converting string parameters to integers.

	"FreeConnect/internal/repositories" // Importing the repository layer to access the database.
	"github.com/gin-gonic/gin"          // Gin framework for routing and HTTP handling.
)

// AdminController handles endpoints that are reserved for administrator actions,
// such as viewing and approving users.
type AdminController struct {
	// userRepo is used to interact with user-related data in the database.
	userRepo repositories.UserRepository
	// Additional repositories (e.g., for projects or transactions) can be added here if needed.
}

// NewAdminController is a constructor function for creating a new AdminController.
// It accepts a UserRepository as a dependency, allowing for dependency injection.
func NewAdminController(userRepo repositories.UserRepository) *AdminController {
	return &AdminController{userRepo: userRepo}
}

// ListAllUsers handles the GET /api/admin/users endpoint.
// It retrieves a list of all users from the database and returns selected fields as JSON.
func (ac *AdminController) ListAllUsers(c *gin.Context) {
	// Retrieve the underlying *gorm.DB instance from the user repository.
	db := ac.userRepo.GetDB()

	// Define an anonymous struct to hold the selected user fields.
	var users []struct {
		UserID   uint    `json:"user_id"`  // Unique user identifier.
		Email    string  `json:"email"`    // User email address.
		Name     string  `json:"name"`     // Full name.
		Role     string  `json:"role"`     // Role: admin, client, or freelancer.
		Bio      string  `json:"bio"`      // Short biography.
		Earnings float64 `json:"earnings"` // Total earnings (if applicable).
	}

	// Execute a raw SQL query to fetch the desired fields from the 'users' table.
	// Note: You may choose to use GORM methods instead of raw SQL if you prefer.
	if err := db.Raw(`SELECT user_id, email, name, role, bio, earnings FROM users`).Scan(&users).Error; err != nil {
		// If an error occurs, respond with a 500 Internal Server Error and the error message.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of users as a JSON response with a 200 OK status.
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ApproveUser handles the PUT /api/admin/users/:id/approve endpoint.
// In this example, it updates a user's role from "pending" to "freelancer".
// This is a sample flow to illustrate administrative approval of new user accounts.
func (ac *AdminController) ApproveUser(c *gin.Context) {
	// Extract the "id" parameter from the URL (the user's ID).
	idStr := c.Param("id")
	// Convert the string ID to an integer.
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If conversion fails, respond with a 400 Bad Request.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the underlying *gorm.DB instance from the user repository.
	db := ac.userRepo.GetDB()

	// Execute a raw SQL update to change the user's role.
	// In a real application, you might include additional checks or use a GORM method.
	if err := db.Exec("UPDATE users SET role = 'freelancer' WHERE user_id = ?", id).Error; err != nil {
		// If an error occurs during the update, respond with a 500 Internal Server Error.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message indicating the user was approved.
	c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}
