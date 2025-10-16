package model

import (
	"time"

	"github.com/google/uuid"
)

// TenantUser represents the tenant_users table (junction table for users and tenants)
type TenantUser struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	Tenant          *Tenant          `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant,omitempty"`
	User            *User            `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Teacher         *Teacher         `gorm:"foreignKey:TenantUserID;constraint:OnDelete:CASCADE" json:"teacher,omitempty"`
	Student         *Student         `gorm:"foreignKey:TenantUserID;constraint:OnDelete:CASCADE" json:"student,omitempty"`
	TenantUserRoles []TenantUserRole `gorm:"foreignKey:TenantUserID;constraint:OnDelete:CASCADE" json:"tenant_user_roles,omitempty"`
}

// TableName returns the table name for TenantUser
func (TenantUser) TableName() string {
	return "tenant_users"
}
