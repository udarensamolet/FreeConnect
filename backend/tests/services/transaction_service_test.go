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

func TestTransactionService(t *testing.T) {
	db := tests.SetupTestDB()
	txRepo := repositories.NewTransactionRepository(db)
	txService := services.NewTransactionService(txRepo)

	// 1) Create a Transaction
	tx := models.Transaction{
		Amount:        250.0,
		Date:          time.Now(),
		PaymentMethod: "bank_transfer",
		Status:        "pending",
		ClientID:      1,
		FreelancerID:  2,
		ProjectID:     3,
	}

	err := txService.CreateTransaction(&tx)
	assert.NoError(t, err)
	assert.NotZero(t, tx.ID)

	// 2) Retrieve
	got, err := txService.GetTransactionByID(tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, 250.0, got.Amount)

	// 3) Update to completed
	tx.Status = "completed"
	err = txService.UpdateTransaction(&tx)
	assert.NoError(t, err)

	updated, err := txService.GetTransactionByID(tx.ID)
	assert.NoError(t, err)
	assert.Equal(t, "completed", updated.Status)

	// 4) Delete
	err = txService.DeleteTransaction(tx.ID)
	assert.NoError(t, err)

	// 5) Confirm
	_, err = txService.GetTransactionByID(tx.ID)
	assert.Error(t, err)
}
