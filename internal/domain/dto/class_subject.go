package dto

import (
	"github.com/google/uuid"
)

// ClassSubject DTOs (linking class, subject, and teacher)
type CreateClassSubjectRequest struct {
	ClassID   *uuid.UUID `json:"class_id" validate:"omitempty,uuid"`
	SubjectID *uuid.UUID `json:"subject_id" validate:"omitempty,uuid"`
	TeacherID *uuid.UUID `json:"teacher_id" validate:"omitempty,uuid"`
}

type UpdateClassSubjectRequest struct {
	ClassID   *uuid.UUID `json:"class_id" validate:"omitempty,uuid"`
	SubjectID *uuid.UUID `json:"subject_id" validate:"omitempty,uuid"`
	TeacherID *uuid.UUID `json:"teacher_id" validate:"omitempty,uuid"`
}

type ClassSubjectQueryParams struct {
	QueryParams
	ClassID   *uuid.UUID `query:"class_id" validate:"omitempty,uuid"`
	SubjectID *uuid.UUID `query:"subject_id" validate:"omitempty,uuid"`
	TeacherID *uuid.UUID `query:"teacher_id" validate:"omitempty,uuid"`
}
