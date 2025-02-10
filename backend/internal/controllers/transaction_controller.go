package controllers

import (
	"net/http" // Provides HTTP status codes.
	"strconv"  // For converting URL parameters from strings to integers.

	"FreeConnect/internal/models"   // Contains the Transaction model.
	"FreeConnect/internal/services" // Contains the TransactionService.
	"github.com/gin-gonic/gin"      // Gin framework for HTTP routing.
)

// TransactionController handles endpoints related to payment transactions.
type TransactionController struct {
	transactionService services.TransactionService // Service layer for transaction operations.
}

// NewTransactionController creates a new TransactionController with the given TransactionService.
func NewTransactionController(ts services.TransactionService) *TransactionController {
	return &TransactionController{transactionService: ts}
}

// CreateTransaction handles POST /api/transactions.
// It creates a new transaction using the provided JSON payload.
func (tc *TransactionController) CreateTransaction(c *gin.Context) {
	// Define a payload struct for binding the JSON request.
	var payload struct {
		Amount        float64 `json:"amount" binding:"required"`         // Amount to be transacted.
		PaymentMethod string  `json:"payment_method" binding:"required"` // Payment method (e.g., credit_card, paypal).
		ClientID      uint    `json:"client_id" binding:"required"`      // ID of the client making the payment.
		FreelancerID  uint    `json:"freelancer_id" binding:"required"`  // ID of the freelancer receiving the payment.
		ProjectID     uint    `json:"project_id" binding:"required"`     // ID of the project associated with the transaction.
	}

	// Bind the JSON payload to the struct.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// Respond with 400 Bad Request if payload is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Transaction model instance.
	transaction := models.Transaction{
		Amount:        payload.Amount,
		PaymentMethod: payload.PaymentMethod,
		Status:        "pending", // Default status is pending.
		ClientID:      payload.ClientID,
		FreelancerID:  payload.FreelancerID,
		ProjectID:     payload.ProjectID,
	}

	// Call the TransactionService to create the transaction.
	if err := tc.transactionService.CreateTransaction(&transaction); err != nil {
		// Respond with 500 Internal Server Error if creation fails.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with 201 Created and return the created transaction.
	c.JSON(http.StatusCreated, gin.H{"transaction": transaction})
}

// GetTransaction handles GET /api/transactions/:id.
// It retrieves a specific transaction by its ID.
func (tc *TransactionController) GetTransaction(c *gin.Context) {
	// Extract the transaction ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Respond with 400 if the ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	// Retrieve the transaction using the TransactionService.
	transaction, err := tc.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		// Respond with 404 if not found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Return the transaction data.
	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// GetTransactionsByProject handles GET /api/projects/:id/transactions.
// It retrieves all transactions for a specific project.
func (tc *TransactionController) GetTransactionsByProject(c *gin.Context) {
	// Extract the project ID from the URL.
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve the transactions using the TransactionService.
	transactions, err := tc.transactionService.GetTransactionsByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of transactions.
	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}

// UpdateTransaction handles PUT /api/transactions/:id.
// It updates fields of an existing transaction.
func (tc *TransactionController) UpdateTransaction(c *gin.Context) {
	// Extract the transaction ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	// Retrieve the current transaction.
	transaction, err := tc.transactionService.GetTransactionByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Transaction not found"})
		return
	}

	// Define a payload for the fields that can be updated.
	var payload struct {
		Amount        float64 `json:"amount"`         // Updated amount.
		PaymentMethod string  `json:"payment_method"` // Updated payment method.
		Status        string  `json:"status"`         // Updated status.
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the transaction fields if new values are provided.
	if payload.Amount != 0 {
		transaction.Amount = payload.Amount
	}
	if payload.PaymentMethod != "" {
		transaction.PaymentMethod = payload.PaymentMethod
	}
	if payload.Status != "" {
		transaction.Status = payload.Status
	}

	// Persist the updated transaction using the TransactionService.
	if err := tc.transactionService.UpdateTransaction(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated transaction.
	c.JSON(http.StatusOK, gin.H{"transaction": transaction})
}

// DeleteTransaction handles DELETE /api/transactions/:id.
// It deletes the specified transaction from the database.
func (tc *TransactionController) DeleteTransaction(c *gin.Context) {
	// Extract the transaction ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid transaction ID"})
		return
	}

	// Call the TransactionService to delete the transaction.
	if err := tc.transactionService.DeleteTransaction(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Transaction deleted successfully"})
}
