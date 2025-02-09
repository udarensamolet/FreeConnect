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
// user_controller.go
func (uc *UserController) RegisterUser(c *gin.Context) {
	var payload struct {
		Name            string `json:"name" binding:"required"`
		Email           string `json:"email" binding:"required,email"`
		Password        string `json:"password" binding:"required,min=6"`
		ConfirmPassword string `json:"confirmPassword" binding:"required,min=6"`
		Role            string `json:"role" binding:"required"` // 'admin', 'client', or 'freelancer'
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure password == confirmPassword
	if payload.Password != payload.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Passwords do not match"})
		return
	}

	user := models.User{
		Name:  payload.Name,
		Email: payload.Email,
		Role:  payload.Role,
	}
	// We only store the main password in the DB (hashed), confirm_password is only for verification
	if err := uc.userService.Register(&user, payload.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"user": user})
}

// GetUser handles GET /api/users/:id.
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

// UpdateUser handles PUT /api/users/:id.
func (uc *UserController) UpdateUser(c *gin.Context) {
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

	var payload struct {
		Name         string  `json:"name"`
		Bio          string  `json:"bio"`
		CompanyName  string  `json:"company_name"`
		Rating       float64 `json:"rating"`
		HourlyRate   float64 `json:"hourly_rate"`
		Availability *bool   `json:"availability"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

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

	if err := uc.userService.UpdateUser(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}

// UpdateUserSkills handles PUT /api/users/:id/skills.
// Expects a JSON payload with an array of skill IDs.
func (uc *UserController) UpdateUserSkills(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var payload struct {
		SkillIDs []uint `json:"skill_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := uc.userService.UpdateUserSkills(uint(id), payload.SkillIDs); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user, err := uc.userService.GetUserByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user})
}
