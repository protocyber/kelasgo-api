package model

import (
	"github.com/google/uuid"
)

// FeeType represents the fee_types table
type FeeType struct {
	BaseModel
	TenantID      uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	Name          string    `gorm:"size:100;not null" json:"name"`
	Description   *string   `gorm:"type:text" json:"description,omitempty"`
	DefaultAmount *float64  `gorm:"type:decimal(10,2);default:0;check:default_amount >= 0" json:"default_amount,omitempty"`
	IsMandatory   bool      `gorm:"default:true" json:"is_mandatory"`
	IsActive      bool      `gorm:"default:true" json:"is_active"`

	// Relationships
	StudentFees []StudentFee `gorm:"foreignKey:FeeTypeID;constraint:OnDelete:CASCADE" json:"student_fees,omitempty"`
}

// TableName returns the table name for FeeType
func (FeeType) TableName() string {
	return "fee_types"
}
