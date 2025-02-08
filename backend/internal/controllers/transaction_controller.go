package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type TransactionController struct {
	transactionService services.TransactionService
}

func NewTransactionController(ts services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: ts}
}

// CreateTransaction handles POST /api/transactions.
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	var payload struct {
		Amount        float64 `json:"amount" binding:"required"`
		PaymentMethod string  `json:"payment_method" binding:"required"`
		ClientID      uint    `json:"client_id" binding:"required"`
		FreelancerID  uint    `json:"freelancer_id" binding:"required"`
		ProjectID     uint    `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	transaction := models.Transaction{
		Amount:        payload.Amount,
		PaymentMethod: payload.PaymentMethod,
		Status:        "pending",
		ClientID:      payload.ClientID,
		FreelancerID:  payload.FreelancerID,
		ProjectID:     payload.ProjectID,
	}
	if err := tc.transactionService.CreateTransaction(&transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"transaction": transaction})
}

// GetTransaction handles GET /api/transactions/:id.
func (tc *TransactionController) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	transaction, err := tc.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// GetTransactionsByProject handles GET /api/projects/:id/transactions.
func (tc *TransactionController) GetTransactionsByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	transactions, err := tc.transactionService.GetTransactionsByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// UpdateTransaction handles PUT /api/transactions/:id.
func (tc *TransactionController) UpdateTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	transaction, err := tc.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}
	var payload struct {
		Amount        float64 `json:"amount"`
		PaymentMethod string  `json:"payment_method"`
		Status        string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Amount != 0 {
		transaction.Amount = payload.Amount
	}
	if payload.PaymentMethod != "" {
		transaction.PaymentMethod = payload.PaymentMethod
	}
	if payload.Status != "" {
		transaction.Status = payload.Status
	}
	if err := tc.transactionService.UpdateTransaction(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// DeleteTransaction handles DELETE /api/transactions/:id.
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}
	if err := tc.transactionService.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
