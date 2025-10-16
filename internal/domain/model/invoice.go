package model

import (
	"time"

	"github.com/google/uuid"
)

// InvoiceStatus represents the invoice status enum
type InvoiceStatus string

const (
	InvoiceStatusDraft     InvoiceStatus = "draft"
	InvoiceStatusSent      InvoiceStatus = "sent"
	InvoiceStatusPaid      InvoiceStatus = "paid"
	InvoiceStatusUnpaid    InvoiceStatus = "unpaid"
	InvoiceStatusOverdue   InvoiceStatus = "overdue"
	InvoiceStatusCancelled InvoiceStatus = "cancelled"
)

// Invoice represents the invoices table
type Invoice struct {
	ID               uuid.UUID     `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TenantID         *uuid.UUID    `gorm:"type:uuid" json:"tenant_id,omitempty"`
	SubscriptionID   *uuid.UUID    `gorm:"type:uuid" json:"subscription_id,omitempty"`
	InvoiceNumber    string        `gorm:"size:50;uniqueIndex;not null" json:"invoice_number"`
	Amount           float64       `gorm:"type:decimal(10,2);not null" json:"amount"`
	Currency         string        `gorm:"size:10;default:'Rp'" json:"currency"`
	IssueDate        time.Time     `gorm:"type:date;default:CURRENT_DATE" json:"issue_date"`
	DueDate          *time.Time    `gorm:"type:date" json:"due_date,omitempty"`
	Status           InvoiceStatus `gorm:"type:invoice_status_enum;default:'unpaid'" json:"status"`
	PaymentDate      *time.Time    `gorm:"type:date" json:"payment_date,omitempty"`
	PaymentReference *string       `gorm:"size:100" json:"payment_reference,omitempty"`

	// Relationships
	Tenant       *Tenant       `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant,omitempty"`
	Subscription *Subscription `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"subscription,omitempty"`
}

// TableName returns the table name for Invoice
func (Invoice) TableName() string {
	return "invoices"
}
