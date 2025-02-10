package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"
	"FreeConnect/tests"
)

func TestInvoiceService(t *testing.T) {
	db := tests.SetupTestDB()
	invoiceRepo := repositories.NewInvoiceRepository(db)
	invoiceService := services.NewInvoiceService(invoiceRepo)

	// 1) Create a sample Invoice
	invoice := models.Invoice{
		InvoiceNumber: "INV-001",
		AmountDue:     500.0,
		PaymentStatus: "pending",
		DueDate:       time.Now().AddDate(0, 0, 7),
		ProjectID:     1, // optional: if a Project with ID=1 exists/required
		ClientID:      2, // optional: if a User with ID=2 exists/required
	}

	err := invoiceService.CreateInvoice(&invoice)
	assert.NoError(t, err)
	assert.NotZero(t, invoice.ID)

	// 2) Retrieve Invoice
	retrieved, err := invoiceService.GetInvoiceByID(invoice.ID)
	assert.NoError(t, err)
	assert.Equal(t, "INV-001", retrieved.InvoiceNumber)

	// 3) Update Invoice
	invoice.PaymentStatus = "paid"
	err = invoiceService.UpdateInvoice(&invoice)
	assert.NoError(t, err)

	updated, err := invoiceService.GetInvoiceByID(invoice.ID)
	assert.NoError(t, err)
	assert.Equal(t, "paid", updated.PaymentStatus)

	// 4) Delete Invoice
	err = invoiceService.DeleteInvoice(invoice.ID)
	assert.NoError(t, err)

	// 5) Confirm deletion
	_, err = invoiceService.GetInvoiceByID(invoice.ID)
	assert.Error(t, err, "Should fail after deletion")
}
