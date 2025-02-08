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

func (s *transactionService) CreateTransaction(transaction *models.Transaction) error {
	return s.repo.Create(transaction)
}

func (s *transactionService) GetTransactionByID(id uint) (*models.Transaction, error) {
	return s.repo.FindByID(id)
}

func (s *transactionService) GetTransactionsByProject(projectID uint) ([]models.Transaction, error) {
	return s.repo.FindByProject(projectID)
}

func (s *transactionService) UpdateTransaction(transaction *models.Transaction) error {
	return s.repo.Update(transaction)
}

func (s *transactionService) DeleteTransaction(id uint) error {
	return s.repo.Delete(id)
}
