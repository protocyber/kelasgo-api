package model

import (
	"github.com/google/uuid"
)

// Enrollment represents the enrollments table
type Enrollment struct {
	BaseModel
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StudentID      *uuid.UUID `gorm:"type:uuid;index" json:"student_id,omitempty"`
	ClassSubjectID *uuid.UUID `gorm:"type:uuid;index" json:"class_subject_id,omitempty"`
	AcademicYearID *uuid.UUID `gorm:"type:uuid;index" json:"academic_year_id,omitempty"`

	// Relationships
	Student      *Student      `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student,omitempty"`
	ClassSubject *ClassSubject `gorm:"foreignKey:ClassSubjectID;constraint:OnDelete:CASCADE" json:"class_subject,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID;constraint:OnDelete:CASCADE" json:"academic_year,omitempty"`
	Grades       []Grade       `gorm:"foreignKey:EnrollmentID;constraint:OnDelete:CASCADE" json:"grades,omitempty"`
}

// TableName returns the table name for Enrollment
func (Enrollment) TableName() string {
	return "enrollments"
}
