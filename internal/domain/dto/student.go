package dto

import (
	"time"

	"github.com/google/uuid"
)

// Student DTOs
type CreateStudentRequest struct {
	TenantUserID  uuid.UUID  `json:"tenant_user_id" validate:"required,uuid"`
	StudentNumber string     `json:"student_number" validate:"required,max=50"`
	AdmissionDate time.Time  `json:"admission_date" validate:"required"`
	ClassID       *uuid.UUID `json:"class_id" validate:"omitempty,uuid"`
	ParentID      *uuid.UUID `json:"parent_id" validate:"omitempty,uuid"`
}

type UpdateStudentRequest struct {
	StudentNumber *string    `json:"student_number" validate:"omitempty,max=50"`
	AdmissionDate *time.Time `json:"admission_date,omitempty"`
	ClassID       *uuid.UUID `json:"class_id" validate:"omitempty,uuid"`
	ParentID      *uuid.UUID `json:"parent_id" validate:"omitempty,uuid"`
}

type StudentQueryParams struct {
	QueryParams
	ClassID  *uuid.UUID `query:"class_id" validate:"omitempty,uuid"`
	ParentID *uuid.UUID `query:"parent_id" validate:"omitempty,uuid"`
}

type BulkDeleteStudentRequest struct {
	IDs []uuid.UUID `json:"ids" validate:"required,min=1,dive,required"`
}
