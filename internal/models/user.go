package models

import (
	"time"
)

// User represents a user in the FreeConnect system.
type User struct {
	ID           uint      `gorm:"column:user_id;primaryKey" json:"user_id"`
	Name         string    `gorm:"type:varchar(255);not null" json:"name"`
	Email        string    `gorm:"type:varchar(255);unique;not null" json:"email"`
	PasswordHash string    `gorm:"type:text;not null" json:"-"` // never return the hash in responses
	Role         string    `gorm:"type:varchar(50);not null;check:role IN ('admin','client','freelancer')" json:"role"`
	Bio          string    `gorm:"type:text" json:"bio,omitempty"`
	CompanyName  string    `gorm:"type:varchar(255)" json:"company_name,omitempty"`
	Rating       float64   `gorm:"type:decimal(3,2);check:rating BETWEEN 0 AND 5" json:"rating,omitempty"`
	HourlyRate   float64   `gorm:"type:decimal(10,2)" json:"hourly_rate,omitempty"`
	Availability bool      `gorm:"default:true" json:"availability"`
	TotalSpent   float64   `gorm:"type:decimal(10,2);default:0.0" json:"total_spent"`
	Earnings     float64   `gorm:"type:decimal(10,2);default:0.0" json:"earnings"`
	LastLogin    time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"last_login"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
