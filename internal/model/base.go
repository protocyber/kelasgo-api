package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/guregu/null/v5"
	"gorm.io/gorm"
)

// Enums to match database schema
type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type DayOfWeek string

const (
	DayMonday    DayOfWeek = "senin"
	DayTuesday   DayOfWeek = "selasa"
	DayWednesday DayOfWeek = "rabu"
	DayThursday  DayOfWeek = "kamis"
	DayFriday    DayOfWeek = "jumat"
	DaySaturday  DayOfWeek = "sabtu"
	DaySunday    DayOfWeek = "minggu"
)

// BaseModel contains common fields for all models with tenant support
type BaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TenantID  uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy null.Int  `json:"created_by,omitempty"`
	UpdatedBy null.Int  `json:"updated_by,omitempty"`
}

// GlobalBaseModel for tables without tenant isolation (like roles, subscription_plans, etc.)
type GlobalBaseModel struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedBy null.Int  `json:"created_by,omitempty"`
	UpdatedBy null.Int  `json:"updated_by,omitempty"`
}

// BeforeCreate hook to set audit fields
func (m *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate hook to update audit fields
func (m *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}

// BeforeCreate hook for GlobalBaseModel
func (m *GlobalBaseModel) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	now := time.Now()
	m.CreatedAt = now
	m.UpdatedAt = now
	return nil
}

// BeforeUpdate hook for GlobalBaseModel
func (m *GlobalBaseModel) BeforeUpdate(tx *gorm.DB) error {
	m.UpdatedAt = time.Now()
	return nil
}
