package model

import (
	"github.com/google/uuid"
)

// Role represents the roles table
type Role struct {
	BaseModel
	TenantID    uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	// Relationships
	Users []User `gorm:"foreignKey:RoleID;constraint:OnDelete:SET NULL" json:"users,omitempty"`
}

// TableName returns the table name for Role
func (Role) TableName() string {
	return "roles"
}
