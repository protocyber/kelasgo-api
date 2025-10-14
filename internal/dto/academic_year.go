package dto

import (
	"time"
)

// Academic Year DTOs
type CreateAcademicYearRequest struct {
	Name      string    `json:"name" validate:"required,max=50"`
	StartDate time.Time `json:"start_date" validate:"required"`
	EndDate   time.Time `json:"end_date" validate:"required"`
	IsActive  *bool     `json:"is_active,omitempty"`
}

type UpdateAcademicYearRequest struct {
	Name      *string    `json:"name" validate:"omitempty,max=50"`
	StartDate *time.Time `json:"start_date,omitempty"`
	EndDate   *time.Time `json:"end_date,omitempty"`
	IsActive  *bool      `json:"is_active,omitempty"`
}
