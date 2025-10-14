package dto

import (
	"github.com/google/uuid"
)

// Subject DTOs
type CreateSubjectRequest struct {
	Name         string     `json:"name" validate:"required,max=100"`
	Code         string     `json:"code" validate:"required,max=50"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id" validate:"omitempty,uuid"`
	Credit       *int       `json:"credit" validate:"omitempty,min=0"`
}

type UpdateSubjectRequest struct {
	Name         *string    `json:"name" validate:"omitempty,max=100"`
	Code         *string    `json:"code" validate:"omitempty,max=50"`
	Description  *string    `json:"description,omitempty"`
	DepartmentID *uuid.UUID `json:"department_id" validate:"omitempty,uuid"`
	Credit       *int       `json:"credit" validate:"omitempty,min=0"`
}
