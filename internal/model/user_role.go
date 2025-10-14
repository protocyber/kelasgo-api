package model

import (
	"github.com/google/uuid"
)

// UserRole represents the user_roles table (many-to-many relationship between users and roles)
type UserRole struct {
	UserID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"user_id"`
	RoleID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"role_id"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
	Role *Role `gorm:"foreignKey:RoleID;constraint:OnDelete:CASCADE" json:"role,omitempty"`
}

// TableName returns the table name for UserRole
func (UserRole) TableName() string {
	return "user_roles"
}
