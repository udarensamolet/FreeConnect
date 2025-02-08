package models

import "time"

// Review represents a review given by one user to another for a project.
type Review struct {
	ID           uint      `gorm:"column:review_id;primaryKey" json:"review_id"`
	Rating       float64   `gorm:"type:decimal(3,2);not null;check:rating BETWEEN 0 AND 5" json:"rating"`
	Comment      string    `gorm:"type:text" json:"comment,omitempty"`
	ReviewedBy   uint      `json:"reviewed_by"`   // The user who leaves the review.
	ReviewedeeID uint      `json:"reviewedee_id"` // The user being reviewed.
	ProjectID    uint      `json:"project_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
