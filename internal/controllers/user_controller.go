package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userService services.UserService
}

func NewUserController(us services.UserService) *UserController {
	return &UserController{userService: us}
}

// RegisterUser handles POST /api/register.
// The expected JSON payload should include "name", "email" and "password".
func (uc *UserController) RegisterUser(c *gin.Context) {
	// Define a temporary payload structure.
	var payload struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required"` // Should be one of 'admin', 'client', 'freelancer'
		// You can add optional fields like bio, company_name, etc.
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a User instance (only setting the fields we require).
	user := models.User{
		Name:  payload.Name,
		Email: payload.Email,
		Role:  payload.Role,
		// Other fields (bio, company_name, etc.) can be added if provided.
	}

	// Call service.Register with the plain password.
	if err := uc.userService.Register(&user, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// GetUser handles GET /api/users/:id
func (uc *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
