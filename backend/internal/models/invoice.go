package models

import "time"

// Invoice represents the invoices table.
type Invoice struct {
	ID            uint      `gorm:"column:invoice_id;primaryKey" json:"invoice_id"`
	InvoiceNumber string    `gorm:"type:varchar(50);unique;not null" json:"invoice_number"`
	Date          time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"date"`
	AmountDue     float64   `gorm:"type:decimal(10,2);not null" json:"amount_due"`
	PaymentStatus string    `gorm:"type:varchar(50);default:'pending';check:payment_status IN ('pending','paid','overdue')" json:"payment_status"`
	DueDate       time.Time `gorm:"not null" json:"due_date"`
	ProjectID     uint      `json:"project_id"`
	ClientID      uint      `json:"client_id"`
}
