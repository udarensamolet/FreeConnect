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

	// NEW: specialized logic to accept a proposal and update the project
	AcceptProposal(proposal *models.Proposal) error
}

type proposalService struct {
	repo repositories.ProposalRepository
}

func NewProposalService(repo repositories.ProposalRepository) ProposalService {
	return &proposalService{repo: repo}
}

// CreateProposal creates a new proposal
func (s *proposalService) CreateProposal(proposal *models.Proposal) error {
	// Additional validations can be added here if needed.
	return s.repo.Create(proposal)
}

// GetProposalByID returns a proposal by ID
func (s *proposalService) GetProposalByID(id uint) (*models.Proposal, error) {
	return s.repo.FindByID(id)
}

// GetProposalsByProject returns all proposals for a given project
func (s *proposalService) GetProposalsByProject(projectID uint) ([]models.Proposal, error) {
	return s.repo.FindByProject(projectID)
}

// UpdateProposal updates a given proposal
func (s *proposalService) UpdateProposal(proposal *models.Proposal) error {
	return s.repo.Update(proposal)
}

// DeleteProposal deletes a proposal by ID
func (s *proposalService) DeleteProposal(id uint) error {
	return s.repo.Delete(id)
}

// AcceptProposal marks a proposal as accepted and updates the linked project
func (s *proposalService) AcceptProposal(proposal *models.Proposal) error {
	// 1) Set proposal status
	proposal.Status = "accepted"
	if err := s.repo.Update(proposal); err != nil {
		return err
	}

	// 2) Retrieve the associated project
	db := s.repo.GetDB() // we can do this now that we added GetDB() in the repository
	projRepo := repositories.NewProjectRepository(db)
	project, err := projRepo.FindByID(proposal.ProjectID)
	if err != nil {
		return err
	}

	// 3) Update project: assign the freelancer and set status = "in_progress"
	project.FreelancerID = &proposal.FreelancerID
	project.Status = "in_progress"
	if err := projRepo.Update(project); err != nil {
		return err
	}

	// 4) (Optional) create a notification for the freelancer
	//    For example:
	/*
		notifRepo := repositories.NewNotificationRepository(db)
		notif := models.Notification{
			Message: "Your proposal has been accepted!",
			UserID:  proposal.FreelancerID,
			Type:    "proposal_update",
		}
		if err := notifRepo.Create(&notif); err != nil {
			return err
		}
	*/

	return nil
}
