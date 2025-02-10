package controllers

import (
	"net/http"
	"strconv"
	"time"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type ProjectController struct {
	projectService services.ProjectService
}

func NewProjectController(ps services.ProjectService) *ProjectController {
	return &ProjectController{projectService: ps}
}

// CreateProject handles POST /api/projects.
// FreelancerID is NO longer mandatory; we remove it from the payload entirely.
func (pc *ProjectController) CreateProject(c *gin.Context) {
	// Updated payload: no 'freelancer_id' field and no binding:"required" for it
	var payload struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Budget      float64 `json:"budget" binding:"required"`
		Duration    int     `json:"duration" binding:"required"`
		Status      string  `json:"status"` // optional; default 'open'
		ClientID    uint    `json:"client_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project := models.Project{
		Title:        payload.Title,
		Description:  payload.Description,
		Budget:       payload.Budget,
		Duration:     payload.Duration,
		Status:       payload.Status,
		ClientID:     payload.ClientID,
		CreationDate: time.Now(),
	}

	// If no status provided, default to "open"
	if project.Status == "" {
		project.Status = "open"
	}

	if err := pc.projectService.CreateProject(&project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// GetProject handles GET /api/projects/:id.
func (pc *ProjectController) GetProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := pc.projectService.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// GetAllProjects handles GET /api/projects?search=...&minBudget=...&maxBudget=...&status=...
func (pc *ProjectController) GetAllProjects(c *gin.Context) {
	search := c.Query("search")
	minBudgetStr := c.Query("minBudget")
	maxBudgetStr := c.Query("maxBudget")
	status := c.Query("status")

	projects, err := pc.projectService.SearchProjects(search, minBudgetStr, maxBudgetStr, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// UpdateProject handles PUT /api/projects/:id.
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	idStr := c.Param("projectId")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := pc.projectService.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	var payload struct {
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Budget       float64 `json:"budget"`
		Duration     int     `json:"duration"`
		Status       string  `json:"status"`
		ClientID     uint    `json:"client_id"`
		FreelancerID *uint   `json:"freelancer_id"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID := c.GetUint("userID")       // from JWT middleware
	userRole := c.GetString("userRole") // also from JWT
	if project.ClientID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this project"})
		return
	}

	if payload.Title != "" {
		project.Title = payload.Title
	}
	if payload.Description != "" {
		project.Description = payload.Description
	}
	if payload.Budget != 0 {
		project.Budget = payload.Budget
	}
	if payload.Duration != 0 {
		project.Duration = payload.Duration
	}
	if payload.Status != "" {
		project.Status = payload.Status
	}
	if payload.ClientID != 0 {
		project.ClientID = payload.ClientID
	}
	// If you still allow updating the freelancer via PUT:
	if payload.FreelancerID != nil {
		project.FreelancerID = payload.FreelancerID
	}

	if err := pc.projectService.UpdateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// SetProjectFreelancer handles POST /api/projects/:id/set-freelancer
// This method assigns a freelancer to a project once a proposal is accepted.
func (pc *ProjectController) SetProjectFreelancer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	project, err := pc.projectService.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Ensure only the project owner or an admin can assign a freelancer
	userID := c.GetUint("id")
	userRole := c.GetString("role")
	if project.ClientID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this project"})
		return
	}

	var payload struct {
		FreelancerID uint `json:"freelancer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.FreelancerID = &payload.FreelancerID
	// Optionally set the project status to 'in_progress'
	// once the freelancer is assigned:
	// project.Status = "in_progress"

	if err := pc.projectService.UpdateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"project": project})
}

// DeleteProject handles DELETE /api/projects/:id.
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	if err := pc.projectService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
