package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type ProposalController struct {
	proposalService services.ProposalService
}

func NewProposalController(ps services.ProposalService) *ProposalController {
	return &ProposalController{proposalService: ps}
}

// CreateProposal handles POST /api/proposals.
func (pc *ProposalController) CreateProposal(c *gin.Context) {
	var payload struct {
		ProposalText      string  `json:"proposal_text" binding:"required"`
		EstimatedDuration int     `json:"estimated_duration" binding:"required"`
		BidAmount         float64 `json:"bid_amount" binding:"required"`
		ProjectID         uint    `json:"project_id" binding:"required"`
		FreelancerID      uint    `json:"freelancer_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	proposal := models.Proposal{
		ProposalText:      payload.ProposalText,
		EstimatedDuration: payload.EstimatedDuration,
		BidAmount:         payload.BidAmount,
		Status:            "pending",
		ProjectID:         payload.ProjectID,
		FreelancerID:      payload.FreelancerID,
	}

	if err := pc.proposalService.CreateProposal(&proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"proposal": proposal})
}

// GetProposal handles GET /api/proposals/:id.
func (pc *ProposalController) GetProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}
	proposal, err := pc.proposalService.GetProposalByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"proposal": proposal})
}

// GetProposalsByProject handles GET /api/projects/:id/proposals.
func (pc *ProposalController) GetProposalsByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	proposals, err := pc.proposalService.GetProposalsByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"proposals": proposals})
}

// UpdateProposal handles PUT /api/proposals/:id.
func (pc *ProposalController) UpdateProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}
	proposal, err := pc.proposalService.GetProposalByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}
	var payload struct {
		ProposalText      string  `json:"proposal_text"`
		EstimatedDuration int     `json:"estimated_duration"`
		BidAmount         float64 `json:"bid_amount"`
		Status            string  `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.ProposalText != "" {
		proposal.ProposalText = payload.ProposalText
	}
	if payload.EstimatedDuration != 0 {
		proposal.EstimatedDuration = payload.EstimatedDuration
	}
	if payload.BidAmount != 0 {
		proposal.BidAmount = payload.BidAmount
	}
	if payload.Status != "" {
		proposal.Status = payload.Status
	}
	if err := pc.proposalService.UpdateProposal(proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"proposal": proposal})
}

// DeleteProposal handles DELETE /api/proposals/:id.
func (pc *ProposalController) DeleteProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}
	if err := pc.proposalService.DeleteProposal(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Proposal deleted successfully"})
}

func (pc *ProposalController) AcceptProposal(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}

	proposal, err := pc.proposalService.GetProposalByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	if err := pc.proposalService.AcceptProposal(proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Proposal accepted successfully"})
}
