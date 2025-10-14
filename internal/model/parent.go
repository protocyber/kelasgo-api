package model

import (
	"github.com/google/uuid"
)

// Parent represents the parents table
type Parent struct {
	BaseModel
	TenantID     uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	FullName     string    `gorm:"size:100;not null" json:"full_name"`
	Phone        *string   `gorm:"size:20" json:"phone,omitempty"`
	Email        *string   `gorm:"size:100" json:"email,omitempty"`
	Address      *string   `gorm:"type:text" json:"address,omitempty"`
	Relationship *string   `gorm:"size:50" json:"relationship,omitempty"`

	// Relationships
	Students []Student `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"students,omitempty"`
}

// TableName returns the table name for Parent
func (Parent) TableName() string {
	return "parents"
}
