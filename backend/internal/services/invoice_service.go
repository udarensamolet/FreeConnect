package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type InvoiceService interface {
	CreateInvoice(invoice *models.Invoice) error
	GetInvoiceByID(id uint) (*models.Invoice, error)
	GetInvoicesByProject(projectID uint) ([]models.Invoice, error)
	UpdateInvoice(invoice *models.Invoice) error
	DeleteInvoice(id uint) error
}

type invoiceService struct {
	repo repositories.InvoiceRepository
}

func NewInvoiceService(repo repositories.InvoiceRepository) InvoiceService {
	return &invoiceService{repo: repo}
}

func (s *invoiceService) CreateInvoice(invoice *models.Invoice) error {
	return s.repo.Create(invoice)
}

func (s *invoiceService) GetInvoiceByID(id uint) (*models.Invoice, error) {
	return s.repo.FindByID(id)
}

func (s *invoiceService) GetInvoicesByProject(projectID uint) ([]models.Invoice, error) {
	return s.repo.FindByProject(projectID)
}

func (s *invoiceService) UpdateInvoice(invoice *models.Invoice) error {
	return s.repo.Update(invoice)
}

func (s *invoiceService) DeleteInvoice(id uint) error {
	return s.repo.Delete(id)
}
