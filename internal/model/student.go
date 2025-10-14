package model

import (
	"time"

	"github.com/google/uuid"
)

// Student represents the students table
type Student struct {
	BaseModel
	TenantID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	TenantUserID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_user_id"`
	StudentNumber string     `gorm:"size:50;not null" json:"student_number"`
	AdmissionDate time.Time  `gorm:"type:date;not null" json:"admission_date"`
	ClassID       *uuid.UUID `gorm:"type:uuid;index" json:"class_id,omitempty"`
	ParentID      *uuid.UUID `gorm:"type:uuid;index" json:"parent_id,omitempty"`

	// Relationships
	TenantUser  *TenantUser  `gorm:"foreignKey:TenantUserID;constraint:OnDelete:CASCADE" json:"tenant_user,omitempty"`
	Class       *Class       `gorm:"foreignKey:ClassID;constraint:OnDelete:SET NULL" json:"class,omitempty"`
	Parent      *Parent      `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL" json:"parent,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"enrollments,omitempty"`
	Attendance  []Attendance `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"attendance,omitempty"`
	StudentFees []StudentFee `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student_fees,omitempty"`
}

// TableName returns the table name for Student
func (Student) TableName() string {
	return "students"
}
