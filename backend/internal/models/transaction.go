package models

import "time"

type Transaction struct {
	ID            uint      `gorm:"column:transaction_id;primaryKey" json:"transaction_id"`
	Amount        float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Date          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	PaymentMethod string    `gorm:"type:varchar(50);check:payment_method IN ('credit_card','paypal','bank_transfer')" json:"payment_method"`
	Status        string    `gorm:"type:varchar(50);default:'pending';check:status IN ('pending','completed','failed')" json:"status"`

	ClientID     uint `json:"client_id"`
	Client       User `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"client,omitempty"`
	FreelancerID uint `json:"freelancer_id"`
	Freelancer   User `gorm:"foreignKey:FreelancerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"freelancer,omitempty"`

	ProjectID uint    `json:"project_id"`
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"project,omitempty"`
}
