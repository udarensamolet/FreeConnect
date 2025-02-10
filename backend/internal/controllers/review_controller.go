package controllers

import (
	"net/http" // Provides HTTP status codes.
	"strconv"  // For converting URL parameters from string to int.

	"FreeConnect/internal/models"   // Contains the Review model definition.
	"FreeConnect/internal/services" // Provides business logic for review operations.
	"github.com/gin-gonic/gin"      // Gin framework for HTTP routing.
)

// ReviewController handles endpoints related to reviews.
type ReviewController struct {
	reviewService services.ReviewService // Service layer for review operations.
}

// NewReviewController constructs a new ReviewController by injecting the given ReviewService.
func NewReviewController(rs services.ReviewService) *ReviewController {
	return &ReviewController{reviewService: rs}
}

// CreateReview handles POST /api/reviews.
// It expects a JSON payload containing rating, comment, reviewed_by, reviewedee_id, and project_id.
func (rc *ReviewController) CreateReview(c *gin.Context) {
	// Define a payload struct to bind the incoming JSON.
	var payload struct {
		Rating       float64 `json:"rating" binding:"required"`        // Rating value (e.g., from 0 to 5).
		Comment      string  `json:"comment"`                          // Optional comment text.
		ReviewedBy   uint    `json:"reviewed_by" binding:"required"`   // ID of the reviewer.
		ReviewedeeID uint    `json:"reviewedee_id" binding:"required"` // ID of the reviewed user.
		ProjectID    uint    `json:"project_id" binding:"required"`    // ID of the project associated with the review.
	}

	// Bind the JSON payload to the payload struct.
	if err := c.ShouldBindJSON(&payload); err != nil {
		// If binding fails, respond with HTTP 400 and the error message.
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new Review model instance with the provided data.
	review := models.Review{
		Rating:       payload.Rating,
		Comment:      payload.Comment,
		ReviewedBy:   payload.ReviewedBy,
		ReviewedeeID: payload.ReviewedeeID,
		ProjectID:    payload.ProjectID,
	}

	// Call the ReviewService to create the review in the database.
	if err := rc.reviewService.CreateReview(&review); err != nil {
		// If an error occurs, respond with HTTP 500.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with HTTP 201 (Created) and return the newly created review.
	c.JSON(http.StatusCreated, gin.H{"review": review})
}

// GetReview handles GET /api/reviews/:id.
// It retrieves a review based on its ID.
func (rc *ReviewController) GetReview(c *gin.Context) {
	// Extract the review ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If conversion fails, return HTTP 400.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Retrieve the review using the ReviewService.
	review, err := rc.reviewService.GetReviewByID(uint(id))
	if err != nil {
		// If the review is not found, return HTTP 404.
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Return the review as JSON.
	c.JSON(http.StatusOK, gin.H{"review": review})
}

// GetReviewsByProject handles GET /api/projects/:id/reviews.
// It retrieves all reviews associated with a specific project.
func (rc *ReviewController) GetReviewsByProject(c *gin.Context) {
	// Extract the project ID from the URL.
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		// Return HTTP 400 if the project ID is invalid.
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}

	// Retrieve the list of reviews for the project using the ReviewService.
	reviews, err := rc.reviewService.GetReviewsByProject(uint(projectID))
	if err != nil {
		// If an error occurs, return HTTP 500.
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the list of reviews.
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// UpdateReview handles PUT /api/reviews/:id.
// It updates the rating and comment fields of an existing review.
func (rc *ReviewController) UpdateReview(c *gin.Context) {
	// Extract the review ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Retrieve the existing review using the ReviewService.
	review, err := rc.reviewService.GetReviewByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}

	// Define a payload struct for updating the review.
	var payload struct {
		Rating  float64 `json:"rating"`  // New rating (optional).
		Comment string  `json:"comment"` // New comment (optional).
	}

	// Bind the JSON payload.
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update the review fields if new values are provided.
	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}
	if payload.Comment != "" {
		review.Comment = payload.Comment
	}

	// Call the service to update the review in the database.
	if err := rc.reviewService.UpdateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated review.
	c.JSON(http.StatusOK, gin.H{"review": review})
}

// DeleteReview handles DELETE /api/reviews/:id.
// It deletes the specified review from the database.
func (rc *ReviewController) DeleteReview(c *gin.Context) {
	// Extract the review ID from the URL.
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}

	// Call the ReviewService to delete the review.
	if err := rc.reviewService.DeleteReview(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
