package models

import "time"

// Proposal references the project and the freelancer who submitted it.
type Proposal struct {
	ID                uint      `gorm:"column:proposal_id;primaryKey" json:"proposal_id"`
	ProposalText      string    `gorm:"type:text;not null" json:"proposal_text"`
	EstimatedDuration int       `gorm:"check:estimated_duration > 0" json:"estimated_duration"`
	BidAmount         float64   `gorm:"type:decimal(10,2);not null" json:"bid_amount"`
	SubmissionDate    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"submission_date"`
	Status            string    `gorm:"type:varchar(50);default:'pending';check:status IN ('pending','accepted','rejected')" json:"status"`

	ProjectID    uint    `json:"project_id"`
	Project      Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`
	FreelancerID uint    `json:"freelancer_id"`
	Freelancer   User    `gorm:"foreignKey:FreelancerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"freelancer,omitempty"`
}
