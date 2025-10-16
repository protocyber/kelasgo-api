package model

import (
	"github.com/google/uuid"
)

// Grade represents the grades table
type Grade struct {
	BaseModel
	TenantID     uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	EnrollmentID *uuid.UUID `gorm:"type:uuid;index" json:"enrollment_id,omitempty"`
	GradeType    string     `gorm:"size:50;check:grade_type IN ('Assignment','Midterm','Final','Other')" json:"grade_type"`
	Score        *float64   `gorm:"type:decimal(5,2)" json:"score,omitempty"`
	Remarks      *string    `gorm:"type:text" json:"remarks,omitempty"`

	// Relationships
	Enrollment *Enrollment `gorm:"foreignKey:EnrollmentID;constraint:OnDelete:CASCADE" json:"enrollment,omitempty"`
}

// TableName returns the table name for Grade
func (Grade) TableName() string {
	return "grades"
}
