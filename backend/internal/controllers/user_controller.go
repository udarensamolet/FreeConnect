package controllers

import (
	"net/http" // For HTTP status codes.
	"strconv"  // For converting string parameters to integers.

	"FreeConnect/internal/models"   // Contains the User model.
	"FreeConnect/internal/services" // Provides the UserService for user operations.
	"github.com/gin-gonic/gin"      // Gin framework for HTTP routing.
)

// UserController handles endpoints related to user operations such as registration,
// profile retrieval, updating profile information, and updating user skills.
type UserController struct {
	userService services.UserService // Service layer for user operations.
}

// NewUserController creates a new UserController by injecting the provided UserService.
func NewUserController(us services.UserService) *UserController {
	return &UserController{userService: us}
}

// RegisterUser handles POST /api/register.
// It expects a JSON payload with name, email, password, confirmPassword, and role.
// It verifies that the password and confirmPassword match and then creates a new user.
func (uc *UserController) RegisterUser(c *gin.Context) {
	// Define a payload structure to bind the incoming JSON.
	var payload struct {
		Name            string `json:"name" binding:"required"`                  // User's full name.
		Email           string `json:"email" binding:"required,email"`           // Valid email address.
		Password        string `json:"password" binding:"required,min=6"`        // Password (minimum 6 characters).
		ConfirmPassword string `json:"confirmPassword" binding:"required,min=6"` // Confirmation of the password.
		Role            string `json:"role" binding:"required"`                  // Role: admin, client, or freelancer.
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the password and confirmPassword fields match.
	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	// Create a new User model instance.
	user := models.User{
		Name:  payload.Name,
		Email: payload.Email,
		Role:  payload.Role,
	}

	// Call the UserService to register the user.
	// Note: The password is hashed within the Register function.
	if err := uc.userService.Register(&user, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Respond with HTTP 201 (Created) and return the created user.
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// GetUser handles GET /api/users/:id.
// It retrieves a user's information by their ID.
func (uc *UserController) GetUser(c *gin.Context) {
	// Extract the user ID from the URL parameter.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the user using the UserService.
	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		// If the user is not found, respond with HTTP 404.
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Return the user's information.
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUser handles PUT /api/users/:id.
// It updates a user's profile information based on the provided JSON payload.
func (uc *UserController) UpdateUser(c *gin.Context) {
	// Extract the user ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve the current user from the UserService.
	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// Define a payload structure for updating user fields.
	var payload struct {
		Name         string  `json:"name"`         // Updated name.
		Bio          string  `json:"bio"`          // Updated biography.
		CompanyName  string  `json:"company_name"` // Updated company name.
		Rating       float64 `json:"rating"`       // Updated rating.
		HourlyRate   float64 `json:"hourly_rate"`  // Updated hourly rate.
		Availability *bool   `json:"availability"` // Updated availability; pointer to detect absence.
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update user fields if new values are provided.
	if payload.Name != "" {
		user.Name = payload.Name
	}
	if payload.Bio != "" {
		user.Bio = payload.Bio
	}
	if payload.CompanyName != "" {
		user.CompanyName = payload.CompanyName
	}
	if payload.Rating != 0 {
		user.Rating = payload.Rating
	}
	if payload.HourlyRate != 0 {
		user.HourlyRate = payload.HourlyRate
	}
	if payload.Availability != nil {
		user.Availability = *payload.Availability
	}

	// Use the UserService to update the user in the database.
	if err := uc.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user information.
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserSkills handles PUT /api/users/:id/skills.
// It expects a JSON payload containing an array of skill IDs to update the user's skills.
func (uc *UserController) UpdateUserSkills(c *gin.Context) {
	// Extract the user ID from the URL parameter.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Define a payload struct to capture the list of skill IDs.
	var payload struct {
		SkillIDs []uint `json:"skill_ids" binding:"required"` // An array of skill IDs is required.
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Call the UserService to update the user's skills.
	if err := uc.userService.UpdateUserSkills(uint(id), payload.SkillIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Retrieve the updated user profile.
	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated user.
	c.JSON(http.StatusOK, gin.H{"user": user})
}
