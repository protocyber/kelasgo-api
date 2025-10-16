package model

import (
	"github.com/google/uuid"
)

// TenantUserRole represents the tenant_user_roles table (many-to-many relationship between tenant_users and roles)
type TenantUserRole struct {
	TenantUserID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"tenant_user_id"`
	RoleID       uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"role_id"`

	// Relationships
	TenantUser *TenantUser `gorm:"foreignKey:TenantUserID;constraint:OnDelete:CASCADE" json:"tenant_user,omitempty"`
	Role       *Role       `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
}

// TableName returns the table name for TenantUserRole
func (TenantUserRole) TableName() string {
	return "tenant_user_roles"
}
