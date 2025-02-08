package models

import "time"

// Proposal represents a freelancer's proposal for a project.
type Proposal struct {
	ID                uint      `gorm:"column:proposal_id;primaryKey" json:"proposal_id"`
	ProposalText      string    `gorm:"type:text;not null" json:"proposal_text"`
	EstimatedDuration int       `gorm:"check:estimated_duration > 0" json:"estimated_duration"` // in days
	BidAmount         float64   `gorm:"type:decimal(10,2);not null" json:"bid_amount"`
	SubmissionDate    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"submission_date"`
	Status            string    `gorm:"type:varchar(50);default:'pending';check:status IN ('pending','accepted','rejected')" json:"status"`
	ProjectID         uint      `json:"project_id"`
	FreelancerID      uint      `json:"freelancer_id"`
}
