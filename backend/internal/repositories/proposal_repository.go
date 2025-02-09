package repositories

import (
	"FreeConnect/internal/models"

	"gorm.io/gorm"
)

type ProposalRepository interface {
	Create(proposal *models.Proposal) error
	FindByID(id uint) (*models.Proposal, error)
	FindByProject(projectID uint) ([]models.Proposal, error)
	Update(proposal *models.Proposal) error
	Delete(id uint) error

	// NEW: Return the underlying *gorm.DB
	GetDB() *gorm.DB
}

type proposalRepository struct {
	db *gorm.DB
}

func NewProposalRepository(db *gorm.DB) ProposalRepository {
	return &proposalRepository{db: db}
}

func (r *proposalRepository) Create(proposal *models.Proposal) error {
	return r.db.Create(proposal).Error
}

func (r *proposalRepository) FindByID(id uint) (*models.Proposal, error) {
	var proposal models.Proposal
	if err := r.db.First(&proposal, id).Error; err != nil {
		return nil, err
	}
	return &proposal, nil
}

func (r *proposalRepository) FindByProject(projectID uint) ([]models.Proposal, error) {
	var proposals []models.Proposal
	if err := r.db.Where("project_id = ?", projectID).Find(&proposals).Error; err != nil {
		return nil, err
	}
	return proposals, nil
}

func (r *proposalRepository) Update(proposal *models.Proposal) error {
	return r.db.Save(proposal).Error
}

func (r *proposalRepository) Delete(id uint) error {
	return r.db.Delete(&models.Proposal{}, id).Error
}

// GetDB returns the underlying *gorm.DB instance
func (r *proposalRepository) GetDB() *gorm.DB {
	return r.db
}
