package models

import "time"

type Notification struct {
	ID         uint      `gorm:"column:notification_id;primaryKey" json:"notification_id"`
	Message    string    `gorm:"type:text;not null" json:"message"`
	Date       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	ReadStatus bool      `gorm:"default:false" json:"read_status"`

	Type   string `gorm:"type:varchar(50);check:type IN ('proposal_update','payment_received','project_status','admin_message')" json:"type"`
	UserID uint   `json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
}
