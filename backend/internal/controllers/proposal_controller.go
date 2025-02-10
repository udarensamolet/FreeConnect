package controllers

import (
	"net/http" // For HTTP status codes.
	"strconv"  // For converting URL parameters to integers.

	"FreeConnect/internal/models"   // For the Proposal model.
	"FreeConnect/internal/services" // For the ProposalService.
	"github.com/gin-gonic/gin"      // Gin framework for routing and HTTP responses.
)

// ProposalController handles endpoints for proposals.
type ProposalController struct {
	proposalService services.ProposalService // Service layer for proposals.
}

// NewProposalController creates a new instance of ProposalController with the provided ProposalService.
func NewProposalController(ps services.ProposalService) *ProposalController {
	return &ProposalController{proposalService: ps}
}

// CreateProposal handles POST /api/proposals.
// It allows freelancers to submit a proposal for a project.
// The expected JSON payload must include proposal text, estimated duration, bid amount,
// project ID, and freelancer ID.
func (pc *ProposalController) CreateProposal(c *gin.Context) {
	// Define a payload structure for the incoming JSON.
	var payload struct {
		ProposalText      string  `json:"proposal_text" binding:"required"`      // Proposal description.
		EstimatedDuration int     `json:"estimated_duration" binding:"required"` // Estimated duration (days).
		BidAmount         float64 `json:"bid_amount" binding:"required"`         // Proposed bid amount.
		ProjectID         uint    `json:"project_id" binding:"required"`         // ID of the project.
		FreelancerID      uint    `json:"freelancer_id" binding:"required"`      // ID of the freelancer submitting the proposal.
	}

	// Bind the JSON payload to the struct.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// Respond with 400 Bad Request if the payload is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Proposal model with the provided values.
	proposal := models.Proposal{
		ProposalText:      payload.ProposalText,
		EstimatedDuration: payload.EstimatedDuration,
		BidAmount:         payload.BidAmount,
		Status:            "pending", // Default status for a new proposal.
		ProjectID:         payload.ProjectID,
		FreelancerID:      payload.FreelancerID,
	}

	// Check if the user role from context is "freelancer".
	userRole := c.GetString("role")
	if userRole != "freelancer" {
		// Only freelancers can create proposals.
		c.JSON(http.StatusForbidden, gin.H{"error": "Only freelancers can create proposals"})
		return
	}

	// Use the proposal service to create the proposal in the database.
	if err := pc.proposalService.CreateProposal(&proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// On success, respond with HTTP 201 (Created) and return the new proposal.
	c.JSON(http.StatusCreated, gin.H{"proposal": proposal})
}

// GetProposal handles GET /api/proposals/:id.
// It retrieves a single proposal by its ID.
func (pc *ProposalController) GetProposal(c *gin.Context) {
	// Extract the proposal ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// Respond with HTTP 400 if the ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}

	// Retrieve the proposal using the service.
	proposal, err := pc.proposalService.GetProposalByID(uint(id))
	if err != nil {
		// Respond with HTTP 404 if the proposal is not found.
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	// Return the proposal data.
	c.JSON(http.StatusOK, gin.H{"proposal": proposal})
}

// GetProposalsByProject handles GET /api/projects/:id/proposals.
// It retrieves all proposals associated with a specific project.
func (pc *ProposalController) GetProposalsByProject(c *gin.Context) {
	// Extract the project ID from the URL.
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve proposals for this project via the service.
	proposals, err := pc.proposalService.GetProposalsByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of proposals.
	c.JSON(http.StatusOK, gin.H{"proposals": proposals})
}

// UpdateProposal handles PUT /api/proposals/:id.
// It updates fields of an existing proposal.
func (pc *ProposalController) UpdateProposal(c *gin.Context) {
	// Extract the proposal ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}

	// Retrieve the current proposal.
	proposal, err := pc.proposalService.GetProposalByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	// Define a payload struct for fields that can be updated.
	var payload struct {
		ProposalText      string  `json:"proposal_text"`      // New proposal text.
		EstimatedDuration int     `json:"estimated_duration"` // New estimated duration.
		BidAmount         float64 `json:"bid_amount"`         // New bid amount.
		Status            string  `json:"status"`             // New status.
	}

	// Bind the incoming JSON.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the proposal's fields if new values are provided.
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

	// Use the service to update the proposal.
	if err := pc.proposalService.UpdateProposal(proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the updated proposal.
	c.JSON(http.StatusOK, gin.H{"proposal": proposal})
}

// DeleteProposal handles DELETE /api/proposals/:id.
// It deletes a proposal from the database.
func (pc *ProposalController) DeleteProposal(c *gin.Context) {
	// Extract the proposal ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}

	// Call the service layer to delete the proposal.
	if err := pc.proposalService.DeleteProposal(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Proposal deleted successfully"})
}

// AcceptProposal handles POST /api/proposals/:id/accept.
// It marks the proposal as accepted and triggers related business logic (e.g., assigning a freelancer to the project).
func (pc *ProposalController) AcceptProposal(c *gin.Context) {
	// Extract the proposal ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid proposal ID"})
		return
	}

	// Retrieve the proposal using the service.
	proposal, err := pc.proposalService.GetProposalByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Proposal not found"})
		return
	}

	// Call the service layer to accept the proposal.
	// This should update the proposal's status and also update the project by assigning the freelancer.
	if err := pc.proposalService.AcceptProposal(proposal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Proposal accepted successfully"})
}
