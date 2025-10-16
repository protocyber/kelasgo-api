package model

import (
	"time"

	"github.com/google/uuid"
)

// FeeStatus represents the fee status enum
type FeeStatus string

const (
	FeeStatusPaid    FeeStatus = "paid"
	FeeStatusUnpaid  FeeStatus = "unpaid"
	FeeStatusPartial FeeStatus = "partial"
	FeeStatusOverdue FeeStatus = "overdue"
)

// StudentFee represents the student_fees table
type StudentFee struct {
	BaseModel
	TenantID       uuid.UUID  `gorm:"type:uuid;not null;index" json:"tenant_id"`
	StudentID      *uuid.UUID `gorm:"type:uuid;index" json:"student_id,omitempty"`
	FeeTypeID      *uuid.UUID `gorm:"type:uuid;index" json:"fee_type_id,omitempty"`
	AcademicYearID *uuid.UUID `gorm:"type:uuid;index" json:"academic_year_id,omitempty"`
	Amount         float64    `gorm:"type:decimal(10,2);not null;check:amount >= 0" json:"amount"`
	DueDate        time.Time  `gorm:"type:date;not null" json:"due_date"`
	Status         FeeStatus  `gorm:"type:fee_status_enum;default:'unpaid'" json:"status"`
	PaymentDate    *time.Time `gorm:"type:date" json:"payment_date,omitempty"`
	PaymentMethod  *string    `gorm:"size:50" json:"payment_method,omitempty"`
	Notes          *string    `gorm:"type:text" json:"notes,omitempty"`

	// Relationships
	Student      *Student      `gorm:"foreignKey:StudentID;constraint:OnDelete:CASCADE" json:"student,omitempty"`
	FeeType      *FeeType      `gorm:"foreignKey:FeeTypeID;constraint:OnDelete:CASCADE" json:"fee_type,omitempty"`
	AcademicYear *AcademicYear `gorm:"foreignKey:AcademicYearID;constraint:OnDelete:CASCADE" json:"academic_year,omitempty"`
}

// TableName returns the table name for StudentFee
func (StudentFee) TableName() string {
	return "student_fees"
}
