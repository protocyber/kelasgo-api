package model

import (
	"time"
)

// Parent represents the parents table
type Parent struct {
	BaseModel
	FullName     string     `gorm:"size:100;not null" json:"full_name"`
	Phone        *string    `gorm:"size:20" json:"phone,omitempty"`
	Email        *string    `gorm:"size:100" json:"email,omitempty"`
	Address      *string    `gorm:"type:text" json:"address,omitempty"`
	Relationship *string    `gorm:"size:50" json:"relationship,omitempty"`
	Birthplace   *string    `gorm:"size:100" json:"birthplace,omitempty"`
	Birthday     *time.Time `gorm:"type:date" json:"birthday,omitempty"`
	Gender       *Gender    `gorm:"type:gender_enum" json:"gender,omitempty"`

	// Relationships
	Students []Student `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"students,omitempty"`
}

// TableName returns the table name for Parent
func (Parent) TableName() string {
	return "parents"
}
