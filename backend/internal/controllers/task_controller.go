package controllers

import (
	"net/http"
	"strconv"
	"time"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type TaskController struct {
	taskService services.TaskService
}

func NewTaskController(ts services.TaskService) *TaskController {
	return &TaskController{taskService: ts}
}

// CreateTask handles POST /api/tasks.
func (tc *TaskController) CreateTask(c *gin.Context) {
	var payload struct {
		Title       string  `json:"title" binding:"required"`
		Description string  `json:"description" binding:"required"`
		Deadline    string  `json:"deadline" binding:"required"` // RFC3339 string
		Budget      float64 `json:"budget"`
		Status      string  `json:"status"` // optional; default "open"
		ProjectID   uint    `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	deadline, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339"})
		return
	}
	task := models.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Deadline:    deadline,
		Budget:      payload.Budget,
		Status:      payload.Status,
		ProjectID:   payload.ProjectID,
	}
	if task.Status == "" {
		task.Status = "open"
	}
	if err := tc.taskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// GetTask handles GET /api/tasks/:id.
func (tc *TaskController) GetTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	task, err := tc.taskService.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// GetTasksByProject handles GET /api/projects/:id/tasks.
func (tc *TaskController) GetTasksByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	tasks, err := tc.taskService.GetTasksByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// UpdateTask handles PUT /api/tasks/:id.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	task, err := tc.taskService.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}
	var payload struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Deadline    string  `json:"deadline"` // RFC3339
		Budget      float64 `json:"budget"`
		Status      string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Title != "" {
		task.Title = payload.Title
	}
	if payload.Description != "" {
		task.Description = payload.Description
	}
	if payload.Deadline != "" {
		deadline, err := time.Parse(time.RFC3339, payload.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339"})
			return
		}
		task.Deadline = deadline
	}
	if payload.Budget != 0 {
		task.Budget = payload.Budget
	}
	if payload.Status != "" {
		task.Status = payload.Status
	}
	if err := tc.taskService.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// DeleteTask handles DELETE /api/tasks/:id.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}
	if err := tc.taskService.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}
