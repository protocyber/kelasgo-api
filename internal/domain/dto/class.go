package dto

import (
	"github.com/google/uuid"
)

// Class DTOs
type CreateClassRequest struct {
	Name              string     `json:"name" validate:"required,max=50"`
	GradeLevel        *int       `json:"grade_level" validate:"omitempty,min=1,max=12"`
	HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id" validate:"omitempty,uuid"`
	AcademicYearID    *uuid.UUID `json:"academic_year_id" validate:"omitempty,uuid"`
}

type UpdateClassRequest struct {
	Name              *string    `json:"name" validate:"omitempty,max=50"`
	GradeLevel        *int       `json:"grade_level" validate:"omitempty,min=1,max=12"`
	HomeroomTeacherID *uuid.UUID `json:"homeroom_teacher_id" validate:"omitempty,uuid"`
	AcademicYearID    *uuid.UUID `json:"academic_year_id" validate:"omitempty,uuid"`
}
