package controllers

import (
	"net/http"

	"FreeConnect/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	userService services.UserService
	jwtService  services.JWTService
}

// NewAuthController creates a new AuthController with the required services.
func NewAuthController(us services.UserService, js services.JWTService) *AuthController {
	return &AuthController{
		userService: us,
		jwtService:  js,
	}
}

// Login handles POST /api/login
// Expects JSON payload: { "email": "...", "password": "..." }
func (ac *AuthController) Login(c *gin.Context) {
	var creds struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify user exists + password match
	user, err := ac.userService.VerifyCredentials(creds.Email, creds.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Generate JWT
	token, err := ac.jwtService.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
