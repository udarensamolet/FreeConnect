package controllers

import (
	"net/http" // Provides HTTP status codes.
	"strconv"  // Used for converting URL parameters to integers.

	"FreeConnect/internal/models"   // Contains the Skill model.
	"FreeConnect/internal/services" // Contains the SkillService for business logic.
	"github.com/gin-gonic/gin"      // Gin framework for HTTP routing.
)

// SkillController handles HTTP endpoints for skills.
type SkillController struct {
	skillService services.SkillService // Service to manage skill operations.
}

// NewSkillController creates a new SkillController with the provided SkillService.
func NewSkillController(ss services.SkillService) *SkillController {
	return &SkillController{skillService: ss}
}

// CreateSkill handles POST /api/skills.
// It creates a new skill record using the provided JSON payload.
func (sc *SkillController) CreateSkill(c *gin.Context) {
	// Define a payload structure for the JSON input.
	var payload struct {
		Name        string `json:"name" binding:"required"` // Skill name is required.
		Level       string `json:"level"`                   // Optional skill level.
		Description string `json:"description"`             // Optional skill description.
	}

	// Bind the JSON request to the payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// If binding fails, respond with HTTP 400.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Skill model instance.
	skill := models.Skill{
		Name:        payload.Name,
		Level:       payload.Level,
		Description: payload.Description,
	}

	// Call the SkillService to persist the new skill.
	if err := sc.skillService.CreateSkill(&skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with HTTP 201 (Created) and return the newly created skill.
	c.JSON(http.StatusCreated, gin.H{"skill": skill})
}

// GetSkill handles GET /api/skills/:id.
// It retrieves a single skill by its ID.
func (sc *SkillController) GetSkill(c *gin.Context) {
	// Extract the skill ID from the URL parameter.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If conversion fails, respond with HTTP 400.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Retrieve the skill using the SkillService.
	skill, err := sc.skillService.GetSkillByID(uint(id))
	if err != nil {
		// If not found, respond with HTTP 404.
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	// Return the skill as JSON.
	c.JSON(http.StatusOK, gin.H{"skill": skill})
}

// GetAllSkills handles GET /api/skills.
// It retrieves all skills from the database.
func (sc *SkillController) GetAllSkills(c *gin.Context) {
	// Call the SkillService to get all skills.
	skills, err := sc.skillService.GetAllSkills()
	if err != nil {
		// If an error occurs, respond with HTTP 500.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of skills.
	c.JSON(http.StatusOK, gin.H{"skills": skills})
}

// UpdateSkill handles PUT /api/skills/:id.
// It updates a specific skill with new data provided in the JSON payload.
func (sc *SkillController) UpdateSkill(c *gin.Context) {
	// Extract the skill ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Retrieve the current skill using the SkillService.
	skill, err := sc.skillService.GetSkillByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}

	// Define a payload structure for the fields that can be updated.
	var payload struct {
		Name        string `json:"name"`        // New name (optional).
		Level       string `json:"level"`       // New level (optional).
		Description string `json:"description"` // New description (optional).
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the skill fields if new values are provided.
	if payload.Name != "" {
		skill.Name = payload.Name
	}
	if payload.Level != "" {
		skill.Level = payload.Level
	}
	if payload.Description != "" {
		skill.Description = payload.Description
	}

	// Call the SkillService to persist the updated skill.
	if err := sc.skillService.UpdateSkill(skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated skill.
	c.JSON(http.StatusOK, gin.H{"skill": skill})
}

// DeleteSkill handles DELETE /api/skills/:id.
// It deletes the specified skill from the database.
func (sc *SkillController) DeleteSkill(c *gin.Context) {
	// Extract the skill ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}

	// Call the SkillService to delete the skill.
	if err := sc.skillService.DeleteSkill(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}
