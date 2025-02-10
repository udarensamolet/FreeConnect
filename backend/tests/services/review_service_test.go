package services_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"FreeConnect/internal/models"
	"FreeConnect/internal/repositories"
	"FreeConnect/internal/services"
	"FreeConnect/tests"
)

func TestReviewService(t *testing.T) {
	db := tests.SetupTestDB()
	revRepo := repositories.NewReviewRepository(db)
	revService := services.NewReviewService(revRepo)

	// 1) Create a Review
	review := models.Review{
		Rating:       4.5,
		Comment:      "Excellent work!",
		ReviewedBy:   1,
		ReviewedeeID: 2,
		ProjectID:    10, // optional
		CreatedAt:    time.Now(),
	}

	err := revService.CreateReview(&review)
	assert.NoError(t, err)
	assert.NotZero(t, review.ID)

	// 2) Retrieve
	got, err := revService.GetReviewByID(review.ID)
	assert.NoError(t, err)
	assert.Equal(t, 4.5, got.Rating)

	// 3) Update
	review.Comment = "Excellent and timely!"
	err = revService.UpdateReview(&review)
	assert.NoError(t, err)

	updated, err := revService.GetReviewByID(review.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Excellent and timely!", updated.Comment)

	// 4) Delete
	err = revService.DeleteReview(review.ID)
	assert.NoError(t, err)

	// 5) Confirm
	_, err = revService.GetReviewByID(review.ID)
	assert.Error(t, err, "Should return error after deletion")
}
