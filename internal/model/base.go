package model

import (
	"github.com/google/uuid"
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
	ID       uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	TenantID uuid.UUID `gorm:"type:uuid;not null;index" json:"tenant_id"`
}

// GlobalBaseModel for tables without tenant isolation (like roles, subscription_plans, etc.)
type GlobalBaseModel struct {
	ID uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
}
