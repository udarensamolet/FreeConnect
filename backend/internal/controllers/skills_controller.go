package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type SkillController struct {
	skillService services.SkillService
}

func NewSkillController(ss services.SkillService) *SkillController {
	return &SkillController{skillService: ss}
}

// CreateSkill handles POST /api/skills.
func (sc *SkillController) CreateSkill(c *gin.Context) {
	var payload struct {
		Name        string `json:"name" binding:"required"`
		Level       string `json:"level"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	skill := models.Skill{
		Name:        payload.Name,
		Level:       payload.Level,
		Description: payload.Description,
	}
	if err := sc.skillService.CreateSkill(&skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"skill": skill})
}

// GetSkill handles GET /api/skills/:id.
func (sc *SkillController) GetSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}
	skill, err := sc.skillService.GetSkillByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"skill": skill})
}

// GetAllSkills handles GET /api/skills.
func (sc *SkillController) GetAllSkills(c *gin.Context) {
	skills, err := sc.skillService.GetAllSkills()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"skills": skills})
}

// UpdateSkill handles PUT /api/skills/:id.
func (sc *SkillController) UpdateSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}
	skill, err := sc.skillService.GetSkillByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Skill not found"})
		return
	}
	var payload struct {
		Name        string `json:"name"`
		Level       string `json:"level"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Name != "" {
		skill.Name = payload.Name
	}
	if payload.Level != "" {
		skill.Level = payload.Level
	}
	if payload.Description != "" {
		skill.Description = payload.Description
	}
	if err := sc.skillService.UpdateSkill(skill); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"skill": skill})
}

// DeleteSkill handles DELETE /api/skills/:id.
func (sc *SkillController) DeleteSkill(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid skill ID"})
		return
	}
	if err := sc.skillService.DeleteSkill(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Skill deleted successfully"})
}
