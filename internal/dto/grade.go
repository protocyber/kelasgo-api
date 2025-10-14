package dto

import (
	"github.com/google/uuid"
)

// Grade DTOs
type CreateGradeRequest struct {
	EnrollmentID *uuid.UUID `json:"enrollment_id" validate:"omitempty,uuid"`
	GradeType    string     `json:"grade_type" validate:"required,oneof=Assignment Midterm Final Other"`
	Score        *float64   `json:"score,omitempty" validate:"omitempty,min=0,max=100"`
	Remarks      *string    `json:"remarks,omitempty"`
}

type UpdateGradeRequest struct {
	GradeType *string  `json:"grade_type" validate:"omitempty,oneof=Assignment Midterm Final Other"`
	Score     *float64 `json:"score,omitempty" validate:"omitempty,min=0,max=100"`
	Remarks   *string  `json:"remarks,omitempty"`
}
