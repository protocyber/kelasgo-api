package model

import (
	"github.com/google/uuid"
)

// ClassSubject represents the class_subjects table (linking class, subject, teacher)
type ClassSubject struct {
	BaseModel
	TenantID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	ClassID   *uuid.UUID `gorm:"type:uuid;index" json:"class_id,omitempty"`
	SubjectID *uuid.UUID `gorm:"type:uuid;index" json:"subject_id,omitempty"`
	TeacherID *uuid.UUID `gorm:"type:uuid;index" json:"teacher_id,omitempty"`

	// Relationships
	Class       *Class       `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class,omitempty"`
	Subject     *Subject     `gorm:"foreignKey:SubjectID;constraint:OnDelete:CASCADE" json:"subject,omitempty"`
	Teacher     *Teacher     `gorm:"foreignKey:TeacherID;constraint:OnDelete:SET NULL" json:"teacher,omitempty"`
	Schedules   []Schedule   `gorm:"foreignKey:ClassSubjectID;constraint:OnDelete:CASCADE" json:"schedules,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:ClassSubjectID;constraint:OnDelete:CASCADE" json:"enrollments,omitempty"`
}

// TableName returns the table name for ClassSubject
func (ClassSubject) TableName() string {
	return "class_subjects"
}
