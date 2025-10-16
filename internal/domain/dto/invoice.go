package dto

import (
	"time"

	"github.com/google/uuid"
)

// Invoice DTOs
type CreateInvoiceRequest struct {
	TenantID         *uuid.UUID `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	SubscriptionID   *uuid.UUID `json:"subscription_id,omitempty" validate:"omitempty,uuid"`
	InvoiceNumber    string     `json:"invoice_number" validate:"required,max=50"`
	Amount           float64    `json:"amount" validate:"required,min=0"`
	Currency         string     `json:"currency,omitempty" validate:"omitempty,max=10"`
	IssueDate        time.Time  `json:"issue_date,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	Status           string     `json:"status,omitempty" validate:"omitempty,oneof=draft sent paid unpaid overdue cancelled"`
	PaymentDate      *time.Time `json:"payment_date,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty" validate:"omitempty,max=100"`
}

type UpdateInvoiceRequest struct {
	TenantID         *uuid.UUID `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	SubscriptionID   *uuid.UUID `json:"subscription_id,omitempty" validate:"omitempty,uuid"`
	InvoiceNumber    *string    `json:"invoice_number,omitempty" validate:"omitempty,max=50"`
	Amount           *float64   `json:"amount,omitempty" validate:"omitempty,min=0"`
	Currency         *string    `json:"currency,omitempty" validate:"omitempty,max=10"`
	IssueDate        *time.Time `json:"issue_date,omitempty"`
	DueDate          *time.Time `json:"due_date,omitempty"`
	Status           *string    `json:"status,omitempty" validate:"omitempty,oneof=draft sent paid unpaid overdue cancelled"`
	PaymentDate      *time.Time `json:"payment_date,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty" validate:"omitempty,max=100"`
}

type InvoiceQueryParams struct {
	QueryParams
	TenantID       string `query:"tenant_id" validate:"omitempty,uuid"`
	SubscriptionID string `query:"subscription_id" validate:"omitempty,uuid"`
	Status         string `query:"status" validate:"omitempty,oneof=draft sent paid unpaid overdue cancelled"`
	IssueDate      string `query:"issue_date" validate:"omitempty"`
	DueDate        string `query:"due_date" validate:"omitempty"`
	InvoiceNumber  string `query:"invoice_number" validate:"omitempty,max=50"`
}

type InvoiceResponse struct {
	ID               uuid.UUID  `json:"id"`
	TenantID         *uuid.UUID `json:"tenant_id,omitempty"`
	SubscriptionID   *uuid.UUID `json:"subscription_id,omitempty"`
	InvoiceNumber    string     `json:"invoice_number"`
	Amount           float64    `json:"amount"`
	Currency         string     `json:"currency"`
	IssueDate        string     `json:"issue_date"`
	DueDate          *string    `json:"due_date,omitempty"`
	Status           string     `json:"status"`
	PaymentDate      *string    `json:"payment_date,omitempty"`
	PaymentReference *string    `json:"payment_reference,omitempty"`
}
