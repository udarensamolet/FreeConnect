package controllers

import (
	"net/http"
	"strconv"

	"FreeConnect/internal/models"
	"FreeConnect/internal/services"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService services.ReviewService
}

func NewReviewController(rs services.ReviewService) *ReviewController {
	return &ReviewController{reviewService: rs}
}

// CreateReview handles POST /api/reviews.
func (rc *ReviewController) CreateReview(c *gin.Context) {
	var payload struct {
		Rating       float64 `json:"rating" binding:"required"`
		Comment      string  `json:"comment"`
		ReviewedBy   uint    `json:"reviewed_by" binding:"required"`
		ReviewedeeID uint    `json:"reviewedee_id" binding:"required"`
		ProjectID    uint    `json:"project_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review := models.Review{
		Rating:       payload.Rating,
		Comment:      payload.Comment,
		ReviewedBy:   payload.ReviewedBy,
		ReviewedeeID: payload.ReviewedeeID,
		ProjectID:    payload.ProjectID,
	}

	if err := rc.reviewService.CreateReview(&review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"review": review})
}

// GetReview handles GET /api/reviews/:id.
func (rc *ReviewController) GetReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}
	review, err := rc.reviewService.GetReviewByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"review": review})
}

// GetReviewsByProject handles GET /api/projects/:id/reviews.
func (rc *ReviewController) GetReviewsByProject(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project ID"})
		return
	}
	reviews, err := rc.reviewService.GetReviewsByProject(uint(projectID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

// UpdateReview handles PUT /api/reviews/:id.
func (rc *ReviewController) UpdateReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}
	review, err := rc.reviewService.GetReviewByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review not found"})
		return
	}
	var payload struct {
		Rating  float64 `json:"rating"`
		Comment string  `json:"comment"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if payload.Rating != 0 {
		review.Rating = payload.Rating
	}
	if payload.Comment != "" {
		review.Comment = payload.Comment
	}
	if err := rc.reviewService.UpdateReview(review); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"review": review})
}

// DeleteReview handles DELETE /api/reviews/:id.
func (rc *ReviewController) DeleteReview(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid review ID"})
		return
	}
	if err := rc.reviewService.DeleteReview(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Review deleted successfully"})
}
