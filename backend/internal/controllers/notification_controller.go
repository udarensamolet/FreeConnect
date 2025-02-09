package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notificationService services.NotificationService
}

func NewNotificationController(ns services.NotificationService) *NotificationController {
	return &NotificationController{notificationService: ns}
}

// CreateNotification handles POST /api/notifications.
func (nc *NotificationController) CreateNotification(c *gin.Context) {
	var payload struct {
		Message string `json:"message" binding:"required"`
		UserID  uint   `json:"user_id" binding:"required"`
		Type    string `json:"type" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	notification := models.Notification{
		Message: payload.Message,
		UserID:  payload.UserID,
		Type:    payload.Type,
	}
	if err := nc.notificationService.CreateNotification(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"notification": notification})
}

// GetNotification handles GET /api/notifications/:id.
func (nc *NotificationController) GetNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	notification, err := nc.notificationService.GetNotificationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notification": notification})
}

// GetNotificationsByUser handles GET /api/notifications/user/:user_id.
func (nc *NotificationController) GetNotificationsByUser(c *gin.Context) {
	userIDStr := c.Param("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	notifications, err := nc.notificationService.GetNotificationsByUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notifications": notifications})
}

// UpdateNotification handles PUT /api/notifications/:id.
func (nc *NotificationController) UpdateNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	notification, err := nc.notificationService.GetNotificationByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found"})
		return
	}
	var payload struct {
		Message    string `json:"message"`
		ReadStatus *bool  `json:"read_status"`
		Type       string `json:"type"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Message != "" {
		notification.Message = payload.Message
	}
	if payload.ReadStatus != nil {
		notification.ReadStatus = *payload.ReadStatus
	}
	if payload.Type != "" {
		notification.Type = payload.Type
	}
	if err := nc.notificationService.UpdateNotification(notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"notification": notification})
}

// DeleteNotification handles DELETE /api/notifications/:id.
func (nc *NotificationController) DeleteNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}
	if err := nc.notificationService.DeleteNotification(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted successfully"})
}
