package model

import (
	"github.com/google/uuid"
)

// Department represents the departments table
type Department struct {
	BaseModel
	TenantID      uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name          string     `gorm:"size:100;not null" json:"name"`
	Description   *string    `gorm:"type:text" json:"description,omitempty"`
	HeadTeacherID *uuid.UUID `gorm:"type:uuid;index" json:"head_teacher_id,omitempty"`

	// Relationships
	HeadTeacher *Teacher  `gorm:"foreignKey:HeadTeacherID;constraint:OnDelete:SET NULL" json:"head_teacher,omitempty"`
	Teachers    []Teacher `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL" json:"teachers,omitempty"`
	Subjects    []Subject `gorm:"foreignKey:DepartmentID;constraint:OnDelete:SET NULL" json:"subjects,omitempty"`
}

// TableName returns the table name for Department
func (Department) TableName() string {
	return "departments"
}
