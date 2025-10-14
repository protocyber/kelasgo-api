package model

import (
	"time"

	"github.com/google/uuid"
)

// SubscriptionStatus represents the subscription status enum
type SubscriptionStatus string

const (
	SubscriptionStatusActive    SubscriptionStatus = "active"
	SubscriptionStatusInactive  SubscriptionStatus = "inactive"
	SubscriptionStatusCancelled SubscriptionStatus = "cancelled"
	SubscriptionStatusExpired   SubscriptionStatus = "expired"
	SubscriptionStatusTrial     SubscriptionStatus = "trial"
)

// Tenant represents the tenants table
type Tenant struct {
	ID                 uuid.UUID          `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Name               string             `gorm:"size:255;not null" json:"name"`
	Domain             *string            `gorm:"size:255;uniqueIndex" json:"domain,omitempty"`
	ContactEmail       *string            `gorm:"size:255" json:"contact_email,omitempty"`
	Phone              *string            `gorm:"size:50" json:"phone,omitempty"`
	Address            *string            `gorm:"type:text" json:"address,omitempty"`
	LogoURL            *string            `gorm:"size:255" json:"logo_url,omitempty"`
	PlanID             *uuid.UUID         `gorm:"type:uuid" json:"plan_id,omitempty"`
	SubscriptionStatus SubscriptionStatus `gorm:"type:subscription_status_enum;default:'active'" json:"subscription_status"`
	CreatedAt          time.Time          `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	CreatedBy          *uuid.UUID         `gorm:"type:uuid" json:"created_by,omitempty"`

	// Relationships
	Plan           *SubscriptionPlan `gorm:"foreignKey:PlanID;constraint:OnDelete:SET NULL" json:"plan,omitempty"`
	Creator        *User             `gorm:"foreignKey:CreatedBy;constraint:OnDelete:SET NULL" json:"creator,omitempty"`
	TenantUsers    []TenantUser      `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant_users,omitempty"`
	Subscriptions  []Subscription    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"subscriptions,omitempty"`
	Invoices       []Invoice         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"invoices,omitempty"`
	Departments    []Department      `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"departments,omitempty"`
	Teachers       []Teacher         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"teachers,omitempty"`
	Students       []Student         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"students,omitempty"`
	Parents        []Parent          `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"parents,omitempty"`
	Classes        []Class           `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"classes,omitempty"`
	Subjects       []Subject         `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"subjects,omitempty"`
	AcademicYears  []AcademicYear    `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"academic_years,omitempty"`
	TenantFeatures []TenantFeature   `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant_features,omitempty"`
}

// TableName returns the table name for Tenant
func (Tenant) TableName() string {
	return "tenants"
}
