package controllers

import (
	"net/http"
	"strconv"
	"time"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type InvoiceController struct {
	invoiceService services.InvoiceService
}

func NewInvoiceController(is services.InvoiceService) *InvoiceController {
	return &InvoiceController{invoiceService: is}
}

// CreateInvoice handles POST /api/invoices.
func (ic *InvoiceController) CreateInvoice(c *gin.Context) {
	var payload struct {
		InvoiceNumber string  `json:"invoice_number" binding:"required"`
		AmountDue     float64 `json:"amount_due" binding:"required"`
		PaymentStatus string  `json:"payment_status"`              // optional; default "pending"
		DueDate       string  `json:"due_date" binding:"required"` // expected RFC3339 format
		ProjectID     uint    `json:"project_id" binding:"required"`
		ClientID      uint    `json:"client_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dueDate, err := time.Parse(time.RFC3339, payload.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid due_date format. Use RFC3339"})
		return
	}

	invoice := models.Invoice{
		InvoiceNumber: payload.InvoiceNumber,
		AmountDue:     payload.AmountDue,
		PaymentStatus: payload.PaymentStatus,
		DueDate:       dueDate,
		ProjectID:     payload.ProjectID,
		ClientID:      payload.ClientID,
	}
	if invoice.PaymentStatus == "" {
		invoice.PaymentStatus = "pending"
	}
	if err := ic.invoiceService.CreateInvoice(&invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"invoice": invoice})
}

// GetInvoice handles GET /api/invoices/:id.
func (ic *InvoiceController) GetInvoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	invoice, err := ic.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// GetInvoicesByProject handles GET /api/projects/:id/invoices.
func (ic *InvoiceController) GetInvoicesByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	invoices, err := ic.invoiceService.GetInvoicesByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoices": invoices})
}

// UpdateInvoice handles PUT /api/invoices/:id.
func (ic *InvoiceController) UpdateInvoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	invoice, err := ic.invoiceService.GetInvoiceByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Invoice not found"})
		return
	}
	var payload struct {
		InvoiceNumber string  `json:"invoice_number"`
		AmountDue     float64 `json:"amount_due"`
		PaymentStatus string  `json:"payment_status"`
		DueDate       string  `json:"due_date"` // RFC3339
		ProjectID     uint    `json:"project_id"`
		ClientID      uint    `json:"client_id"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
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
	if err := ic.invoiceService.UpdateInvoice(invoice); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"invoice": invoice})
}

// DeleteInvoice handles DELETE /api/invoices/:id.
func (ic *InvoiceController) DeleteInvoice(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid invoice ID"})
		return
	}
	if err := ic.invoiceService.DeleteInvoice(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Invoice deleted successfully"})
}
