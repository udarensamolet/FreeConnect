package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByID(id uint) (*models.Transaction, error)
	FindByProject(projectID uint) ([]models.Transaction, error)
	Update(transaction *models.Transaction) error
	Delete(id uint) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *models.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var t models.Transaction
	if err := r.db.First(&t, id).Error; err != nil {
		return nil, err
	}
	return &t, nil
}

func (r *transactionRepository) FindByProject(projectID uint) ([]models.Transaction, error) {
	var transactions []models.Transaction
	if err := r.db.Where("project_id = ?", projectID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *transactionRepository) Update(transaction *models.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Transaction{}, id).Error
}
