package models

import "time"

type Review struct {
	ID        uint      `gorm:"column:review_id;primaryKey" json:"review_id"`
	Rating    float64   `gorm:"type:decimal(3,2);not null;check:rating BETWEEN 0 AND 5" json:"rating"`
	Comment   string    `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	ReviewedBy uint `json:"reviewed_by"` // the user leaving the review
	Reviewer   User `gorm:"foreignKey:ReviewedBy;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"reviewer,omitempty"`

	ReviewedeeID uint `json:"reviewedee_id"` // the user being reviewed
	Reviewedee   User `gorm:"foreignKey:ReviewedeeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"reviewedee,omitempty"`

	ProjectID uint    `json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`
}
