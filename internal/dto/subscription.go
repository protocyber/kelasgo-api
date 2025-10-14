package dto

import (
	"time"

	"github.com/google/uuid"
)

// Subscription DTOs
type CreateSubscriptionRequest struct {
	TenantID      *uuid.UUID `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	PlanID        *uuid.UUID `json:"plan_id,omitempty" validate:"omitempty,uuid"`
	StartDate     time.Time  `json:"start_date" validate:"required"`
	EndDate       time.Time  `json:"end_date" validate:"required"`
	IsTrial       bool       `json:"is_trial"`
	Status        string     `json:"status,omitempty" validate:"omitempty,oneof=active inactive cancelled expired"`
	AmountPaid    *float64   `json:"amount_paid,omitempty" validate:"omitempty,min=0"`
	PaymentMethod *string    `json:"payment_method,omitempty" validate:"omitempty,max=50"`
	InvoiceID     *string    `json:"invoice_id,omitempty" validate:"omitempty,max=100"`
}

type UpdateSubscriptionRequest struct {
	TenantID      *uuid.UUID `json:"tenant_id,omitempty" validate:"omitempty,uuid"`
	PlanID        *uuid.UUID `json:"plan_id,omitempty" validate:"omitempty,uuid"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	IsTrial       *bool      `json:"is_trial,omitempty"`
	Status        *string    `json:"status,omitempty" validate:"omitempty,oneof=active inactive cancelled expired"`
	AmountPaid    *float64   `json:"amount_paid,omitempty" validate:"omitempty,min=0"`
	PaymentMethod *string    `json:"payment_method,omitempty" validate:"omitempty,max=50"`
	InvoiceID     *string    `json:"invoice_id,omitempty" validate:"omitempty,max=100"`
}

type SubscriptionQueryParams struct {
	QueryParams
	TenantID  string `query:"tenant_id" validate:"omitempty,uuid"`
	PlanID    string `query:"plan_id" validate:"omitempty,uuid"`
	Status    string `query:"status" validate:"omitempty,oneof=active inactive cancelled expired"`
	IsTrial   *bool  `query:"is_trial" validate:"omitempty"`
	StartDate string `query:"start_date" validate:"omitempty"`
	EndDate   string `query:"end_date" validate:"omitempty"`
}

type SubscriptionResponse struct {
	ID            uuid.UUID  `json:"id"`
	TenantID      *uuid.UUID `json:"tenant_id,omitempty"`
	PlanID        *uuid.UUID `json:"plan_id,omitempty"`
	StartDate     string     `json:"start_date"`
	EndDate       string     `json:"end_date"`
	IsTrial       bool       `json:"is_trial"`
	Status        string     `json:"status"`
	AmountPaid    *float64   `json:"amount_paid,omitempty"`
	PaymentMethod *string    `json:"payment_method,omitempty"`
	InvoiceID     *string    `json:"invoice_id,omitempty"`
	CreatedAt     string     `json:"created_at"`
}
