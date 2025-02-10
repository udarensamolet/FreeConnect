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

// task_controller.go

func (tc *TaskController) EditTask(c *gin.Context) {
	// 1) Parse route params: project :id and task :taskId
	projectIDStr := c.Param("projectId")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	taskIDStr := c.Param("taskId")
	taskID, err := strconv.Atoi(taskIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid task ID"})
		return
	}

	// 2) Fetch the task from DB
	task, err := tc.taskService.GetTaskByID(uint(taskID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	// 3) Optional: Confirm the task truly belongs to that project
	if task.ProjectID != uint(projectID) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "This task does not belong to the specified project"})
		return
	}

	// 4) (Recommended) Check user ownership or admin
	//   - Retrieve the project to see project.ClientID
	//   - Compare with c.GetUint("userID") or userRole == "admin"
	//   - For brevity, we skip that or only do it if your logic demands it.

	// 5) Bind JSON payload for updated fields
	var payload struct {
		Title       string     `json:"title"`
		Description string     `json:"description"`
		Deadline    *time.Time `json:"deadline"`
		Budget      *float64   `json:"budget"`
		Status      string     `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 6) Update fields if provided
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

	// 7) Call service to persist the changes
	if err := tc.taskService.UpdateTask(task); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 8) Return updated task
	c.JSON(http.StatusOK, gin.H{"task": task})
}
