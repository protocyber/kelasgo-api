package model

import (
	"time"

	"github.com/google/uuid"
)

// AcademicYear represents the academic_years table
type AcademicYear struct {
	BaseModel
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`

	// Relationships
	Classes     []Class      `gorm:"foreignKey:AcademicYearID;constraint:OnDelete:SET NULL" json:"classes,omitempty"`
	Enrollments []Enrollment `gorm:"foreignKey:AcademicYearID;constraint:OnDelete:CASCADE" json:"enrollments,omitempty"`
	StudentFees []StudentFee `gorm:"foreignKey:AcademicYearID;constraint:OnDelete:CASCADE" json:"student_fees,omitempty"`
}

// TableName returns the table name for AcademicYear
func (AcademicYear) TableName() string {
	return "academic_years"
}
