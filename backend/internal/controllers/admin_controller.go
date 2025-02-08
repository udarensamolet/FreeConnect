package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/repositories"
	"github.com/gin-gonic/gin"
)

type AdminController struct {
	userRepo repositories.UserRepository
	// add other repos if you need
}

func NewAdminController(userRepo repositories.UserRepository) *AdminController {
	return &AdminController{userRepo: userRepo}
}

// ListAllUsers => GET /api/admin/users
func (ac *AdminController) ListAllUsers(c *gin.Context) {
	db := ac.userRepo.GetDB()

	var users []struct {
		UserID   uint    `json:"user_id"`
		Email    string  `json:"email"`
		Name     string  `json:"name"`
		Role     string  `json:"role"`
		Bio      string  `json:"bio"`
		Earnings float64 `json:"earnings"`
	}
	if err := db.Raw(`SELECT user_id, email, name, role, bio, earnings FROM users`).Scan(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}

// ApproveUser => PUT /api/admin/users/:id/approve
// Example: set user role from "pending" to "freelancer" if your system had that flow.
func (ac *AdminController) ApproveUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	db := ac.userRepo.GetDB()
	// Example update: set role to 'freelancer'
	if err := db.Exec("UPDATE users SET role = 'freelancer' WHERE user_id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User approved successfully"})
}
