package dto

import (
	"github.com/google/uuid"
)

// Enrollment DTOs
type CreateEnrollmentRequest struct {
	StudentID      *uuid.UUID `json:"student_id" validate:"omitempty,uuid"`
	ClassSubjectID *uuid.UUID `json:"class_subject_id" validate:"omitempty,uuid"`
	AcademicYearID *uuid.UUID `json:"academic_year_id" validate:"omitempty,uuid"`
}

type UpdateEnrollmentRequest struct {
	StudentID      *uuid.UUID `json:"student_id" validate:"omitempty,uuid"`
	ClassSubjectID *uuid.UUID `json:"class_subject_id" validate:"omitempty,uuid"`
	AcademicYearID *uuid.UUID `json:"academic_year_id" validate:"omitempty,uuid"`
}

type EnrollmentQueryParams struct {
	QueryParams
	StudentID      *uuid.UUID `query:"student_id" validate:"omitempty,uuid"`
	ClassSubjectID *uuid.UUID `query:"class_subject_id" validate:"omitempty,uuid"`
	AcademicYearID *uuid.UUID `query:"academic_year_id" validate:"omitempty,uuid"`
}
