package dto

import (
	"github.com/google/uuid"
)

// Department DTOs
type CreateDepartmentRequest struct {
	Name          string     `json:"name" validate:"required,max=100"`
	Description   *string    `json:"description,omitempty"`
	HeadTeacherID *uuid.UUID `json:"head_teacher_id" validate:"omitempty,uuid"`
}

type UpdateDepartmentRequest struct {
	Name          *string    `json:"name" validate:"omitempty,max=100"`
	Description   *string    `json:"description,omitempty"`
	HeadTeacherID *uuid.UUID `json:"head_teacher_id" validate:"omitempty,uuid"`
}
