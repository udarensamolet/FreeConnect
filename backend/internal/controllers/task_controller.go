package controllers

import (
	"net/http"
	"strconv"
	"time"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

// TaskController manages HTTP endpoints for tasks (creation, retrieval, update, deletion, etc.).
type TaskController struct {
	taskService services.TaskService
}

// NewTaskController creates a new instance of TaskController with the provided TaskService.
func NewTaskController(ts services.TaskService) *TaskController {
	return &TaskController{taskService: ts}
}

// CreateTask handles POST /api/tasks.
// It creates a new task for a project.
func (tc *TaskController) CreateTask(c *gin.Context) {
	// Define the JSON payload structure for creating a task.
	var payload struct {
		Title       string  `json:"title" binding:"required"`       // Task title is mandatory.
		Description string  `json:"description" binding:"required"` // Task description is mandatory.
		Deadline    string  `json:"deadline" binding:"required"`    // Deadline in RFC3339 format.
		Budget      float64 `json:"budget"`                         // Optional task budget.
		Status      string  `json:"status"`                         // Optional; defaults to "open".
		ProjectID   uint    `json:"project_id" binding:"required"`  // ID of the project this task belongs to.
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the deadline string into a time.Time value.
	deadline, err := time.Parse(time.RFC3339, payload.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339"})
		return
	}

	// Create a Task model instance.
	task := models.Task{
		Title:       payload.Title,
		Description: payload.Description,
		Deadline:    deadline,
		Budget:      payload.Budget,
		Status:      payload.Status,
		ProjectID:   payload.ProjectID,
	}

	// If no status provided, default to "open".
	if task.Status == "" {
		task.Status = "open"
	}

	// Call the service to persist the task in the database.
	if err := tc.taskService.CreateTask(&task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the created task with HTTP 201 status.
	c.JSON(http.StatusCreated, gin.H{"task": task})
}

// GetTask handles GET /api/tasks/:id.
// It retrieves a specific task by its ID.
func (tc *TaskController) GetTask(c *gin.Context) {
	// Get the task ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Retrieve the task using the TaskService.
	task, err := tc.taskService.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Return the task data as JSON.
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// GetTasksByProject handles GET /api/projects/:id/tasks.
// It returns all tasks that belong to a specific project.
func (tc *TaskController) GetTasksByProject(c *gin.Context) {
	// Get the project ID from the URL.
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Use the service to fetch tasks related to the project.
	tasks, err := tc.taskService.GetTasksByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of tasks.
	c.JSON(http.StatusOK, gin.H{"tasks": tasks})
}

// UpdateTask handles PUT /api/tasks/:id.
// This is a generic update endpoint for tasks.
func (tc *TaskController) UpdateTask(c *gin.Context) {
	// Parse the task ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Retrieve the task from the service.
	task, err := tc.taskService.GetTaskByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// Define a payload struct for the update.
	var payload struct {
		Title       string  `json:"title"`
		Description string  `json:"description"`
		Deadline    string  `json:"deadline"` // expects RFC3339 format
		Budget      float64 `json:"budget"`
		Status      string  `json:"status"`
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields if provided.
	if payload.Title != "" {
		task.Title = payload.Title
	}
	if payload.Description != "" {
		task.Description = payload.Description
	}
	if payload.Deadline != "" {
		parsedDeadline, err := time.Parse(time.RFC3339, payload.Deadline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid deadline format. Use RFC3339"})
			return
		}
		task.Deadline = parsedDeadline
	}
	if payload.Budget != 0 {
		task.Budget = payload.Budget
	}
	if payload.Status != "" {
		task.Status = payload.Status
	}

	// Call the service to update the task in the database.
	if err := tc.taskService.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated task.
	c.JSON(http.StatusOK, gin.H{"task": task})
}

// DeleteTask handles DELETE /api/tasks/:id.
// It deletes the specified task.
func (tc *TaskController) DeleteTask(c *gin.Context) {
	// Extract the task ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// Call the service to delete the task.
	if err := tc.taskService.DeleteTask(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})
}

// EditTask handles PUT /api/projects/:id/tasks/:taskId/edit.
// This endpoint specifically edits a task within a given project.
func (tc *TaskController) EditTask(c *gin.Context) {
	// 1) Parse the project ID from the URL parameter.
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// 2) Parse the task ID from the URL parameter.
	taskIDStr := c.Param("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 3) Retrieve the task from the service.
	task, err := tc.taskService.GetTaskByID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 4) Validate that the task belongs to the specified project.
	if task.ProjectID != uint(projectID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This task does not belong to the specified project"})
		return
	}

	// 5) (Optional) Here you can add ownership checks to ensure that
	//    only the project owner or an admin can edit the task.
	//    For example, fetch the project and compare project.ClientID with c.GetUint("userID").

	// 6) Define a payload for fields that can be updated.
	var payload struct {
		Title       string     `json:"title"`       // New title (if provided)
		Description string     `json:"description"` // New description (if provided)
		Deadline    *time.Time `json:"deadline"`    // New deadline (optional, pointer to detect absence)
		Budget      *float64   `json:"budget"`      // New budget (optional)
		Status      string     `json:"status"`      // New status (if provided)
	}

	// Bind the JSON payload to the structure.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 7) Update the task fields if new values are provided.
	if payload.Title != "" {
		task.Title = payload.Title
	}
	if payload.Description != "" {
		task.Description = payload.Description
	}
	if payload.Deadline != nil {
		task.Deadline = *payload.Deadline
	}
	if payload.Budget != nil {
		task.Budget = *payload.Budget
	}
	if payload.Status != "" {
		task.Status = payload.Status
	}

	// 8) Persist the updated task using the service.
	if err := tc.taskService.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 9) Respond with the updated task.
	c.JSON(http.StatusOK, gin.H{"task": task})
}
