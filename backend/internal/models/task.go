package models

import (
	"time"
)

type Task struct {
	ID          uint      `gorm:"column:task_id;primaryKey" json:"task_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Deadline    time.Time `gorm:"not null" json:"deadline"`
	Budget      float64   `gorm:"type:decimal(10,2)" json:"budget"`
	Status      string    `gorm:"type:varchar(50);default:'open';check:status IN ('open','in_progress','completed','cancelled')" json:"status"`

	ProjectID uint    `json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`
}
