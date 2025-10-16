package model

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlanStatus represents the subscription plan status enum
type SubscriptionPlanStatus string

const (
	SubscriptionPlanStatusActive    SubscriptionPlanStatus = "active"
	SubscriptionPlanStatusInactive  SubscriptionPlanStatus = "inactive"
	SubscriptionPlanStatusCancelled SubscriptionPlanStatus = "cancelled"
	SubscriptionPlanStatusExpired   SubscriptionPlanStatus = "expired"
)

// Subscription represents the subscriptions table
type Subscription struct {
	ID            uuid.UUID              `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TenantID      *uuid.UUID             `gorm:"type:uuid" json:"tenant_id,omitempty"`
	PlanID        *uuid.UUID             `gorm:"type:uuid" json:"plan_id,omitempty"`
	StartDate     time.Time              `gorm:"type:date;not null" json:"start_date"`
	EndDate       time.Time              `gorm:"type:date;not null" json:"end_date"`
	IsTrial       bool                   `gorm:"default:false" json:"is_trial"`
	Status        SubscriptionPlanStatus `gorm:"type:subscription_plan_status_enum;default:'active'" json:"status"`
	AmountPaid    *float64               `gorm:"type:decimal(10,2)" json:"amount_paid,omitempty"`
	PaymentMethod *string                `gorm:"size:50" json:"payment_method,omitempty"`
	InvoiceID     *string                `gorm:"size:100" json:"invoice_id,omitempty"`
	CreatedAt     time.Time              `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Tenant   *Tenant           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant,omitempty"`
	Plan     *SubscriptionPlan `gorm:"foreignKey:PlanID;constraint:OnDelete:CASCADE" json:"plan,omitempty"`
	Invoices []Invoice         `gorm:"foreignKey:SubscriptionID;constraint:OnDelete:CASCADE" json:"invoices,omitempty"`
}

// TableName returns the table name for Subscription
func (Subscription) TableName() string {
	return "subscriptions"
}
