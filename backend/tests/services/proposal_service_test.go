package services_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"
	"FreeConnect/tests"
)

func TestProposalService(t *testing.T) {
	db := tests.SetupTestDB()
	propRepo := repositories.NewProposalRepository(db)
	propService := services.NewProposalService(propRepo)

	// 1) Create a Proposal
	proposal := models.Proposal{
		ProposalText:      "I can finish this in 10 days",
		EstimatedDuration: 10,
		BidAmount:         800.0,
		ProjectID:         1, // If a project with ID=1 exists
		FreelancerID:      2, // If user with ID=2 is a freelancer
	}

	err := propService.CreateProposal(&proposal)
	assert.NoError(t, err)
	assert.NotZero(t, proposal.ID)

	// 2) Retrieve
	got, err := propService.GetProposalByID(proposal.ID)
	assert.NoError(t, err)
	assert.Equal(t, 800.0, got.BidAmount)

	// 3) Update
	proposal.BidAmount = 900.0
	err = propService.UpdateProposal(&proposal)
	assert.NoError(t, err)

	updated, err := propService.GetProposalByID(proposal.ID)
	assert.NoError(t, err)
	assert.Equal(t, 900.0, updated.BidAmount)

	// 4) Accept
	err = propService.AcceptProposal(&proposal)
	assert.NoError(t, err, "AcceptProposal might fail if Project or Freelancer are missing in DB")

	// 5) Delete
	err = propService.DeleteProposal(proposal.ID)
	assert.NoError(t, err)
}
