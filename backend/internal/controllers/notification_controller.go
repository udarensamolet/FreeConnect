package controllers

import (
	"net/http" // For HTTP status codes.
	"strconv"  // To convert string parameters to integer.

	"FreeConnect/internal/models"   // For the Notification model.
	"FreeConnect/internal/services" // For the NotificationService.
	"github.com/gin-gonic/gin"      // Gin framework for routing.
)

// NotificationController handles endpoints for creating, retrieving, updating, and deleting notifications.
type NotificationController struct {
	notificationService services.NotificationService // Service layer for notification operations.
}

// NewNotificationController constructs a new NotificationController with the provided service.
func NewNotificationController(ns services.NotificationService) *NotificationController {
	return &NotificationController{notificationService: ns}
}

// CreateNotification handles POST /api/notifications.
// It expects a JSON payload with a message, user_id, and type, then creates a new notification.
func (nc *NotificationController) CreateNotification(c *gin.Context) {
	// Define a payload structure to bind the JSON input.
	var payload struct {
		Message string `json:"message" binding:"required"` // The text of the notification.
		UserID  uint   `json:"user_id" binding:"required"` // The recipient's user ID.
		Type    string `json:"type" binding:"required"`    // The type of notification.
	}

	// Bind the JSON payload to the struct.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// Return HTTP 400 if the payload is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Notification model instance with the provided data.
	notification := models.Notification{
		Message: payload.Message,
		UserID:  payload.UserID,
		Type:    payload.Type,
	}

	// Call the service to create the notification in the database.
	if err := nc.notificationService.CreateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with HTTP 201 (Created) and return the created notification.
	c.JSON(http.StatusCreated, gin.H{"notification": notification})
}

// GetNotification handles GET /api/notifications/:id.
// It retrieves a notification by its ID.
func (nc *NotificationController) GetNotification(c *gin.Context) {
	// Extract the notification ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Return HTTP 400 if the ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Retrieve the notification from the service layer.
	notification, err := nc.notificationService.GetNotificationByID(uint(id))
	if err != nil {
		// Return HTTP 404 if the notification is not found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	// Return the notification as JSON.
	c.JSON(http.StatusOK, gin.H{"notification": notification})
}

// GetNotificationsByUser handles GET /api/notifications/user/:user_id.
// It retrieves all notifications for a given user.
func (nc *NotificationController) GetNotificationsByUser(c *gin.Context) {
	// Extract the user_id from the URL.
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		// Return HTTP 400 if the user_id is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	// Retrieve notifications for the given user from the service.
	notifications, err := nc.notificationService.GetNotificationsByUser(uint(userID))
	if err != nil {
		// Return HTTP 500 if an error occurs.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the notifications as JSON.
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// UpdateNotification handles PUT /api/notifications/:id.
// It updates a notification's message, type, or read status.
func (nc *NotificationController) UpdateNotification(c *gin.Context) {
	// Extract the notification ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Retrieve the notification from the service.
	notification, err := nc.notificationService.GetNotificationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}

	// Define a payload for updatable fields.
	var payload struct {
		Message    string `json:"message"`     // New message text.
		ReadStatus *bool  `json:"read_status"` // New read status; pointer to detect if provided.
		Type       string `json:"type"`        // New notification type.
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the notification fields if provided.
	if payload.Message != "" {
		notification.Message = payload.Message
	}
	if payload.ReadStatus != nil {
		notification.ReadStatus = *payload.ReadStatus
	}
	if payload.Type != "" {
		notification.Type = payload.Type
	}

	// Update the notification in the database via the service.
	if err := nc.notificationService.UpdateNotification(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated notification.
	c.JSON(http.StatusOK, gin.H{"notification": notification})
}

// DeleteNotification handles DELETE /api/notifications/:id.
// It deletes a notification from the database.
func (nc *NotificationController) DeleteNotification(c *gin.Context) {
	// Extract the notification ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Call the service to delete the notification.
	if err := nc.notificationService.DeleteNotification(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
