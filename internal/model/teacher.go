package model

import (
	"time"

	"github.com/google/uuid"
)

// Teacher represents the teachers table
type Teacher struct {
	BaseModel
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	TenantUserID   uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_user_id"`
	EmployeeNumber *string    `gorm:"size:50;uniqueIndex" json:"employee_number,omitempty"`
	HireDate       *time.Time `gorm:"type:date" json:"hire_date,omitempty"`
	DepartmentID   *uuid.UUID `gorm:"type:uuid;index" json:"department_id,omitempty"`
	Qualification  *string    `gorm:"size:100" json:"qualification,omitempty"`
	Position       *string    `gorm:"size:100" json:"position,omitempty"`

	// Relationships
	TenantUser      *TenantUser    `gorm:"foreignKey:TenantUserID;constraint:OnDelete:CASCADE" json:"tenant_user,omitempty"`
	Department      *Department    `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL" json:"department,omitempty"`
	Classes         []Class        `gorm:"foreignKey:HomeroomTeacherID;constraint:OnDelete:SET NULL" json:"classes,omitempty"`
	ClassSubjects   []ClassSubject `gorm:"foreignKey:TeacherID;constraint:OnDelete:SET NULL" json:"class_subjects,omitempty"`
	HeadDepartments []Department   `gorm:"foreignKey:HeadTeacherID;constraint:OnDelete:SET NULL" json:"head_departments,omitempty"`
}

// TableName returns the table name for Teacher
func (Teacher) TableName() string {
	return "teachers"
}
