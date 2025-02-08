package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type TransactionService interface {
	CreateTransaction(transaction *models.Transaction) error
	GetTransactionByID(id uint) (*models.Transaction, error)
	GetTransactionsByProject(projectID uint) ([]models.Transaction, error)
	UpdateTransaction(transaction *models.Transaction) error
	DeleteTransaction(id uint) error
}

type transactionService struct {
	repo repositories.TransactionRepository
}

func NewTransactionService(repo repositories.TransactionRepository) TransactionService {
	return &transactionService{repo: repo}
}

// CreateTransaction creates a new transaction
// Optionally, we could finalize immediately if status == "completed"
func (s *transactionService) CreateTransaction(transaction *models.Transaction) error {
	return s.repo.Create(transaction)
}

// GetTransactionByID returns transaction by ID
func (s *transactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	return s.repo.FindByID(id)
}

// GetTransactionsByProject lists all transactions for a given project
func (s *transactionService) GetTransactionsByProject(projectID uint) ([]models.Transaction, error) {
	return s.repo.FindByProject(projectID)
}

// UpdateTransaction updates a transaction; if status changes to "completed",
// it updates the client/freelancer balances automatically.
func (s *transactionService) UpdateTransaction(transaction *models.Transaction) error {
	// Compare old vs new
	oldTx, err := s.repo.FindByID(transaction.ID)
	if err != nil {
		return err
	}

	// If we're transitioning to "completed", update user balances
	if oldTx.Status != "completed" && transaction.Status == "completed" {
		if err := s.updateUserBalances(transaction); err != nil {
			return err
		}
	}

	return s.repo.Update(transaction)
}

// DeleteTransaction deletes a transaction
func (s *transactionService) DeleteTransaction(id uint) error {
	return s.repo.Delete(id)
}

// updateUserBalances updates the client & freelancer after the transaction completes
func (s *transactionService) updateUserBalances(tx *models.Transaction) error {
	db := s.repo.GetDB() // from the transaction_repository

	// Fetch the client
	var client models.User
	if err := db.First(&client, tx.ClientID).Error; err != nil {
		return err
	}

	// Fetch the freelancer
	var freelancer models.User
	if err := db.First(&freelancer, tx.FreelancerID).Error; err != nil {
		return err
	}

	// Increase freelancer earnings
	freelancer.Earnings += tx.Amount
	// Increase client total_spent
	client.TotalSpent += tx.Amount

	// Save
	if err := db.Save(&freelancer).Error; err != nil {
		return err
	}
	if err := db.Save(&client).Error; err != nil {
		return err
	}

	return nil
}
