package controllers

import (
	"net/http" // For HTTP status codes and responses.

	"FreeConnect/internal/services" // Business logic for authentication and user management.
	"github.com/gin-gonic/gin"      // Gin framework for HTTP routing and handling.
)

// AuthController handles authentication-related endpoints (login, token generation, etc.).
type AuthController struct {
	userService services.UserService // Service used for user-related operations (e.g., verifying credentials).
	jwtService  services.JWTService  // Service used for generating and validating JWT tokens.
}

// NewAuthController is a constructor that returns a new AuthController instance.
// It injects both the UserService and JWTService into the controller.
func NewAuthController(us services.UserService, js services.JWTService) *AuthController {
	return &AuthController{
		userService: us,
		jwtService:  js,
	}
}

// Login handles the POST /api/login endpoint.
// It expects a JSON payload containing "email" and "password", validates them,
// and returns a JWT token along with the user data if authentication succeeds.
func (ac *AuthController) Login(c *gin.Context) {
	// Define an inline struct to bind the incoming JSON payload.
	var creds struct {
		Email    string `json:"email" binding:"required,email"` // Email must be provided and valid.
		Password string `json:"password" binding:"required"`    // Password must be provided.
	}

	// Bind the JSON input to the creds struct.
	if err := c.ShouldBindJSON(&creds); err != nil {
		// If the binding fails, respond with HTTP 400 (Bad Request) and the error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the user's credentials using the UserService.
	user, err := ac.userService.VerifyCredentials(creds.Email, creds.Password)
	if err != nil {
		// If verification fails, respond with HTTP 401 (Unauthorized).
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate a JWT token for the authenticated user using the JWTService.
	token, err := ac.jwtService.GenerateToken(user)
	if err != nil {
		// If token generation fails, respond with HTTP 500 (Internal Server Error).
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	// On success, respond with HTTP 200 (OK) and return the token and user information in JSON format.
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
