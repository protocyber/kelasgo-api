package model

import (
	"github.com/google/uuid"
)

// Notification represents the notifications table
type Notification struct {
	BaseModel
	TenantID uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID   *uuid.UUID `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Title    string     `gorm:"size:100" json:"title"`
	Message  string     `gorm:"type:text" json:"message"`
	IsRead   bool       `gorm:"default:false" json:"is_read"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user,omitempty"`
}

// TableName returns the table name for Notification
func (Notification) TableName() string {
	return "notifications"
}
