package model

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionPlan represents the subscription_plans table
type SubscriptionPlan struct {
	ID             uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name           string    `gorm:"size:100;not null" json:"name"`
	PriceMonthly   float64   `gorm:"type:decimal(10,2);default:0" json:"price_monthly"`
	PriceYearly    float64   `gorm:"type:decimal(10,2);default:0" json:"price_yearly"`
	MaxStudents    *int      `json:"max_students,omitempty"`
	MaxTeachers    *int      `json:"max_teachers,omitempty"`
	StorageLimitMb *int      `json:"storage_limit_mb,omitempty"`
	Features       string    `gorm:"type:jsonb;default:'{}'" json:"features"`
	IsActive       bool      `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt      time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`

	// Relationships
	Tenants       []Tenant       `gorm:"foreignKey:PlanID" json:"tenants,omitempty"`
	Subscriptions []Subscription `gorm:"foreignKey:PlanID" json:"subscriptions,omitempty"`
}

// TableName returns the table name for SubscriptionPlan
func (SubscriptionPlan) TableName() string {
	return "subscription_plans"
}
