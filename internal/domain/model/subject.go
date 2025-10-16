package model

import (
	"github.com/google/uuid"
)

// Subject represents the subjects table
type Subject struct {
	BaseModel
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Code         string     `gorm:"size:50;not null" json:"code"`
	Description  *string    `gorm:"type:text" json:"description,omitempty"`
	DepartmentID *uuid.UUID `gorm:"type:uuid;index" json:"department_id,omitempty"`
	Credit       int        `gorm:"default:0" json:"credit"`

	// Relationships
	Department    *Department    `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL" json:"department,omitempty"`
	ClassSubjects []ClassSubject `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE" json:"class_subjects,omitempty"`
}

// TableName returns the table name for Subject
func (Subject) TableName() string {
	return "subjects"
}
