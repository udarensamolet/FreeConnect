package controllers

import (
	"net/http" // Provides HTTP status codes.
	"strconv"  // Used for string-to-int conversion.
	"time"     // Used for parsing date/time strings.

	"FreeConnect/internal/models"   // Database models.
	"FreeConnect/internal/services" // Business logic layer for invoices.
	"github.com/gin-gonic/gin"      // Gin framework for routing.
)

// InvoiceController handles endpoints related to invoices.
type InvoiceController struct {
	invoiceService services.InvoiceService // Service to manage invoices.
}

// NewInvoiceController constructs a new InvoiceController by injecting the InvoiceService.
func NewInvoiceController(is services.InvoiceService) *InvoiceController {
	return &InvoiceController{invoiceService: is}
}

// CreateInvoice handles POST /api/invoices.
// It expects a JSON payload with the invoice number, amount due, payment status, due date,
// project ID, and client ID, then creates a new invoice in the database.
func (ic *InvoiceController) CreateInvoice(c *gin.Context) {
	// Define a struct to bind the incoming JSON payload.
	var payload struct {
		InvoiceNumber string  `json:"invoice_number" binding:"required"` // Unique invoice number.
		AmountDue     float64 `json:"amount_due" binding:"required"`     // Amount that is due.
		PaymentStatus string  `json:"payment_status"`                    // Optional; defaults to "pending".
		DueDate       string  `json:"due_date" binding:"required"`       // Due date in RFC3339 format.
		ProjectID     uint    `json:"project_id" binding:"required"`     // Associated project.
		ClientID      uint    `json:"client_id" binding:"required"`      // ID of the client to be charged.
	}

	// Bind the JSON payload to our struct.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// Return HTTP 400 if the payload is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse the due date string into a time.Time object.
	dueDate, err := time.Parse(time.RFC3339, payload.DueDate)
	if err != nil {
		// If parsing fails, return HTTP 400 with an appropriate message.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use RFC3339"})
		return
	}

	// Create a new Invoice model instance.
	invoice := models.Invoice{
		InvoiceNumber: payload.InvoiceNumber,
		AmountDue:     payload.AmountDue,
		PaymentStatus: payload.PaymentStatus,
		DueDate:       dueDate,
		ProjectID:     payload.ProjectID,
		ClientID:      payload.ClientID,
	}

	// If PaymentStatus is empty, default it to "pending".
	if invoice.PaymentStatus == "" {
		invoice.PaymentStatus = "pending"
	}

	// Call the service layer to save the invoice in the database.
	if err := ic.invoiceService.CreateInvoice(&invoice); err != nil {
		// Return HTTP 500 if saving fails.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// On success, respond with HTTP 201 and return the invoice.
	c.JSON(http.StatusCreated, gin.H{"invoice": invoice})
}

// GetInvoice handles GET /api/invoices/:id.
// It retrieves an invoice by its ID.
func (ic *InvoiceController) GetInvoice(c *gin.Context) {
	// Extract the invoice ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Return HTTP 400 if the invoice ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	// Retrieve the invoice using the service.
	invoice, err := ic.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		// Return HTTP 404 if the invoice is not found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	// Respond with the invoice data.
	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// GetInvoicesByProject handles GET /api/projects/:id/invoices.
// It returns all invoices related to a specific project.
func (ic *InvoiceController) GetInvoicesByProject(c *gin.Context) {
	// Extract the project ID from the URL.
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		// Return HTTP 400 if the project ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve invoices for the given project.
	invoices, err := ic.invoiceService.GetInvoicesByProject(uint(projectID))
	if err != nil {
		// Return HTTP 500 if an error occurs during retrieval.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of invoices.
	c.JSON(http.StatusOK, gin.H{"invoices": invoices})
}

// UpdateInvoice handles PUT /api/invoices/:id.
// It updates an existing invoice with new data.
func (ic *InvoiceController) UpdateInvoice(c *gin.Context) {
	// Extract the invoice ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	// Retrieve the current invoice from the database.
	invoice, err := ic.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}

	// Define a payload for fields that can be updated.
	var payload struct {
		InvoiceNumber string  `json:"invoice_number"`
		AmountDue     float64 `json:"amount_due"`
		PaymentStatus string  `json:"payment_status"`
		DueDate       string  `json:"due_date"` // Expected in RFC3339 format.
		ProjectID     uint    `json:"project_id"`
		ClientID      uint    `json:"client_id"`
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update invoice fields if new values are provided.
	if payload.InvoiceNumber != "" {
		invoice.InvoiceNumber = payload.InvoiceNumber
	}
	if payload.AmountDue != 0 {
		invoice.AmountDue = payload.AmountDue
	}
	if payload.PaymentStatus != "" {
		invoice.PaymentStatus = payload.PaymentStatus
	}
	if payload.DueDate != "" {
		dueDate, err := time.Parse(time.RFC3339, payload.DueDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use RFC3339"})
			return
		}
		invoice.DueDate = dueDate
	}
	if payload.ProjectID != 0 {
		invoice.ProjectID = payload.ProjectID
	}
	if payload.ClientID != 0 {
		invoice.ClientID = payload.ClientID
	}

	// Update the invoice using the service.
	if err := ic.invoiceService.UpdateInvoice(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated invoice.
	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// DeleteInvoice handles DELETE /api/invoices/:id.
// It deletes the specified invoice from the database.
func (ic *InvoiceController) DeleteInvoice(c *gin.Context) {
	// Extract the invoice ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}

	// Call the service to delete the invoice.
	if err := ic.invoiceService.DeleteInvoice(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
