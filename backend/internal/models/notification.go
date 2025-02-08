package models

import "time"

// Notification represents a notification for a user.
type Notification struct {
	ID         uint      `gorm:"column:notification_id;primaryKey" json:"notification_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	Date       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	UserID     uint      `json:"user_id"`
	Type       string    `gorm:"type:varchar(50);check:type IN ('proposal_update','payment_received','project_status','admin_message')" json:"type"`
	ReadStatus bool      `gorm:"default:false" json:"read_status"`
}
