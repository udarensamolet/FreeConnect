package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type InvoiceRepository interface {
	Create(invoice *models.Invoice) error
	FindByID(id uint) (*models.Invoice, error)
	FindByProject(projectID uint) ([]models.Invoice, error)
	Update(invoice *models.Invoice) error
	Delete(id uint) error
}

type invoiceRepository struct {
	db *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &invoiceRepository{db: db}
}

func (r *invoiceRepository) Create(invoice *models.Invoice) error {
	return r.db.Create(invoice).Error
}

func (r *invoiceRepository) FindByID(id uint) (*models.Invoice, error) {
	var invoice models.Invoice
	if err := r.db.First(&invoice, id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) FindByProject(projectID uint) ([]models.Invoice, error) {
	var invoices []models.Invoice
	if err := r.db.Where("project_id = ?", projectID).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}

func (r *invoiceRepository) Update(invoice *models.Invoice) error {
	return r.db.Save(invoice).Error
}

func (r *invoiceRepository) Delete(id uint) error {
	return r.db.Delete(&models.Invoice{}, id).Error
}
