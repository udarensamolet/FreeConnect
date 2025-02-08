package models

import "time"

// Skill represents the skills table.
type Skill struct {
	ID          uint      `gorm:"column:skill_id;primaryKey" json:"skill_id"`
	Name        string    `gorm:"type:varchar(255);not null" json:"name"`
	Level       string    `gorm:"type:varchar(50)" json:"level,omitempty"`
	Description string    `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
