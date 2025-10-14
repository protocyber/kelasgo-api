package dto

import (
	"github.com/google/uuid"
)

// SubscriptionPlan DTOs
type CreateSubscriptionPlanRequest struct {
	Name           string  `json:"name" validate:"required,max=100"`
	PriceMonthly   float64 `json:"price_monthly" validate:"min=0"`
	PriceYearly    float64 `json:"price_yearly" validate:"min=0"`
	MaxStudents    *int    `json:"max_students,omitempty" validate:"omitempty,min=0"`
	MaxTeachers    *int    `json:"max_teachers,omitempty" validate:"omitempty,min=0"`
	StorageLimitMb *int    `json:"storage_limit_mb,omitempty" validate:"omitempty,min=0"`
	Features       string  `json:"features,omitempty"`
	IsActive       bool    `json:"is_active"`
}

type UpdateSubscriptionPlanRequest struct {
	Name           *string  `json:"name,omitempty" validate:"omitempty,max=100"`
	PriceMonthly   *float64 `json:"price_monthly,omitempty" validate:"omitempty,min=0"`
	PriceYearly    *float64 `json:"price_yearly,omitempty" validate:"omitempty,min=0"`
	MaxStudents    *int     `json:"max_students,omitempty" validate:"omitempty,min=0"`
	MaxTeachers    *int     `json:"max_teachers,omitempty" validate:"omitempty,min=0"`
	StorageLimitMb *int     `json:"storage_limit_mb,omitempty" validate:"omitempty,min=0"`
	Features       *string  `json:"features,omitempty"`
	IsActive       *bool    `json:"is_active,omitempty"`
}

type SubscriptionPlanQueryParams struct {
	QueryParams
	Name     string `query:"name" validate:"omitempty,max=100"`
	IsActive *bool  `query:"is_active" validate:"omitempty"`
}

type SubscriptionPlanResponse struct {
	ID             uuid.UUID `json:"id"`
	Name           string    `json:"name"`
	PriceMonthly   float64   `json:"price_monthly"`
	PriceYearly    float64   `json:"price_yearly"`
	MaxStudents    *int      `json:"max_students,omitempty"`
	MaxTeachers    *int      `json:"max_teachers,omitempty"`
	StorageLimitMb *int      `json:"storage_limit_mb,omitempty"`
	Features       string    `json:"features"`
	IsActive       bool      `json:"is_active"`
	CreatedAt      string    `json:"created_at"`
	UpdatedAt      string    `json:"updated_at"`
}
