package controllers

import (
	"net/http"
	"strconv"
	"time"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

// ProjectController handles HTTP requests related to projects.
// It communicates with the ProjectService to perform business logic.
type ProjectController struct {
	projectService services.ProjectService
}

// NewProjectController creates a new instance of ProjectController with the provided service.
// This constructor function is used to inject dependencies.
func NewProjectController(ps services.ProjectService) *ProjectController {
	return &ProjectController{projectService: ps}
}

// CreateProject handles POST /api/projects.
// This endpoint creates a new project. Note that freelancer_id is no longer required at creation.
func (pc *ProjectController) CreateProject(c *gin.Context) {
	// Define a payload structure to bind the JSON input.
	// Notice that FreelancerID is removed from the payload.
	var payload struct {
		Title       string  `json:"title" binding:"required"`       // Project title is mandatory.
		Description string  `json:"description" binding:"required"` // Project description is mandatory.
		Budget      float64 `json:"budget" binding:"required"`      // Project budget is mandatory.
		Duration    int     `json:"duration" binding:"required"`    // Duration (in days) is mandatory.
		Status      string  `json:"status"`                         // Optional; defaults to "open" if not provided.
		ClientID    uint    `json:"client_id" binding:"required"`   // The client ID that creates the project.
	}

	// Bind the JSON from the request into the payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// If binding fails, respond with HTTP 400 and the error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Project model instance using the input data.
	project := models.Project{
		Title:        payload.Title,
		Description:  payload.Description,
		Budget:       payload.Budget,
		Duration:     payload.Duration,
		Status:       payload.Status,
		ClientID:     payload.ClientID,
		CreationDate: time.Now(), // Set the current time as the creation date.
	}

	// If the status is empty, default it to "open".
	if project.Status == "" {
		project.Status = "open"
	}

	// Call the service layer to create the project in the database.
	if err := pc.projectService.CreateProject(&project); err != nil {
		// If an error occurs during creation, respond with HTTP 500.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with HTTP 201 and return the created project as JSON.
	c.JSON(http.StatusCreated, gin.H{"project": project})
}

// GetProject handles GET /api/projects/:id.
// It retrieves a single project by its ID.
func (pc *ProjectController) GetProject(c *gin.Context) {
	// Extract the project ID from the URL parameter.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Respond with HTTP 400 if the project ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve the project from the service layer.
	project, err := pc.projectService.GetProjectByID(uint(id))
	if err != nil {
		// Respond with HTTP 404 if the project is not found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Return the project as JSON with HTTP 200.
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// GetAllProjects handles GET /api/projects with optional filtering parameters.
// It accepts query parameters such as search, minBudget, maxBudget, and status.
func (pc *ProjectController) GetAllProjects(c *gin.Context) {
	// Retrieve query parameters from the URL.
	search := c.Query("search")
	minBudgetStr := c.Query("minBudget")
	maxBudgetStr := c.Query("maxBudget")
	status := c.Query("status")

	// Call the service method to perform the search/filter.
	projects, err := pc.projectService.SearchProjects(search, minBudgetStr, maxBudgetStr, status)
	if err != nil {
		// Respond with HTTP 500 if an error occurs.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the filtered projects as JSON.
	c.JSON(http.StatusOK, gin.H{"projects": projects})
}

// UpdateProject handles PUT /api/projects/:id.
// It allows updating project fields; only the project owner or admin can update.
func (pc *ProjectController) UpdateProject(c *gin.Context) {
	// Get the project ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve the project from the database.
	project, err := pc.projectService.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Define a payload structure to receive updated values.
	var payload struct {
		Title        string  `json:"title"`
		Description  string  `json:"description"`
		Budget       float64 `json:"budget"`
		Duration     int     `json:"duration"`
		Status       string  `json:"status"`
		ClientID     uint    `json:"client_id"`
		FreelancerID *uint   `json:"freelancer_id"`
	}

	// Bind the incoming JSON to the payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Retrieve user details from context (set by JWT middleware)
	userID := c.GetUint("id")
	userRole := c.GetString("role")
	// Only allow the client who created the project or an admin to update.
	if project.ClientID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this project"})
		return
	}

	// Update fields if provided in the payload.
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
	// If a new freelancer ID is provided, update it.
	if payload.FreelancerID != nil {
		project.FreelancerID = payload.FreelancerID
	}

	// Persist the updated project using the service.
	if err := pc.projectService.UpdateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated project as JSON.
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// SetProjectFreelancer handles POST /api/projects/:id/set-freelancer.
// This endpoint is used to assign a freelancer to an existing project after creation.
func (pc *ProjectController) SetProjectFreelancer(c *gin.Context) {
	// Parse the project ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve the project from the database.
	project, err := pc.projectService.GetProjectByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		return
	}

	// Ensure that only the project owner (or an admin) can assign a freelancer.
	userID := c.GetUint("id")
	userRole := c.GetString("role")
	if project.ClientID != userID && userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not the owner of this project"})
		return
	}

	// Define a payload that expects a freelancer_id.
	var payload struct {
		FreelancerID uint `json:"freelancer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set the freelancer for the project.
	project.FreelancerID = &payload.FreelancerID

	// Optionally, you can set project status to "in_progress" here.
	// project.Status = "in_progress"

	// Update the project in the database.
	if err := pc.projectService.UpdateProject(project); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated project.
	c.JSON(http.StatusOK, gin.H{"project": project})
}

// DeleteProject handles DELETE /api/projects/:id.
// It removes the project from the database.
func (pc *ProjectController) DeleteProject(c *gin.Context) {
	// Extract the project ID from the URL parameter.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Call the service layer to delete the project.
	if err := pc.projectService.DeleteProject(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Project deleted successfully"})
}
