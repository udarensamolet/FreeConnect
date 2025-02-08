package services

import (
	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
)

type ProposalService interface {
	CreateProposal(proposal *models.Proposal) error
	GetProposalByID(id uint) (*models.Proposal, error)
	GetProposalsByProject(projectID uint) ([]models.Proposal, error)
	UpdateProposal(proposal *models.Proposal) error
	DeleteProposal(id uint) error
}

type proposalService struct {
	repo repositories.ProposalRepository
}

func NewProposalService(repo repositories.ProposalRepository) ProposalService {
	return &proposalService{repo: repo}
}

func (s *proposalService) CreateProposal(proposal *models.Proposal) error {
	// Additional validations can be added here if needed.
	return s.repo.Create(proposal)
}

func (s *proposalService) GetProposalByID(id uint) (*models.Proposal, error) {
	return s.repo.FindByID(id)
}

func (s *proposalService) GetProposalsByProject(projectID uint) ([]models.Proposal, error) {
	return s.repo.FindByProject(projectID)
}

func (s *proposalService) UpdateProposal(proposal *models.Proposal) error {
	return s.repo.Update(proposal)
}

func (s *proposalService) DeleteProposal(id uint) error {
	return s.repo.Delete(id)
}
