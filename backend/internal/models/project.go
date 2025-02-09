package models

import (
	"time"
)

// Project references the client who creates it, and an optional freelancer if awarded
type Project struct {
	ID           uint      `gorm:"column:project_id;primaryKey" json:"project_id"`
	Title        string    `gorm:"type:varchar(255);not null" json:"title"`
	Description  string    `gorm:"type:text;not null" json:"description"`
	Budget       float64   `gorm:"type:decimal(10,2);not null" json:"budget"`
	Duration     int       `gorm:"check:duration > 0" json:"duration"` // in days
	Status       string    `gorm:"type:varchar(50);default:'open';check:status IN ('open','in_progress','completed','cancelled')" json:"status"`
	CreationDate time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"creation_date"`

	ClientID uint `json:"client_id"`
	Client   User `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"client,omitempty"`

	// Once a proposal is accepted, set FreelancerID
	FreelancerID *uint `json:"freelancer_id,omitempty"`
	Freelancer   *User `gorm:"foreignKey:FreelancerID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"freelancer,omitempty"`

	Tasks []Task `json:"tasks,omitempty" gorm:"foreignKey:ProjectID"`
}
