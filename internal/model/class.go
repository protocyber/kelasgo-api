package model

import (
	"github.com/google/uuid"
)

// Class represents the classes table
type Class struct {
	BaseModel
	TenantID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name              string     `gorm:"size:50;not null" json:"name"`
	GradeLevel        *int       `json:"grade_level,omitempty"`
	HomeroomTeacherID *uuid.UUID `gorm:"type:uuid;index" json:"homeroom_teacher_id,omitempty"`
	AcademicYearID    *uuid.UUID `gorm:"type:uuid;index" json:"academic_year_id,omitempty"`

	// Relationships
	HomeroomTeacher *Teacher       `gorm:"foreignKey:HomeroomTeacherID;constraint:OnDelete:SET NULL" json:"homeroom_teacher,omitempty"`
	AcademicYear    *AcademicYear  `gorm:"foreignKey:AcademicYearID;constraint:OnDelete:SET NULL" json:"academic_year,omitempty"`
	Students        []Student      `gorm:"foreignKey:ClassID;constraint:OnDelete:SET NULL" json:"students,omitempty"`
	ClassSubjects   []ClassSubject `gorm:"foreignKey:ClassID;constraint:OnDelete:CASCADE" json:"class_subjects,omitempty"`
}

// TableName returns the table name for Class
func (Class) TableName() string {
	return "classes"
}
