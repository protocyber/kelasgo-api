package dto

import (
	"time"

	"github.com/google/uuid"
)

// Student Fee DTOs
type CreateStudentFeeRequest struct {
	StudentID      *uuid.UUID `json:"student_id" validate:"omitempty,uuid"`
	FeeTypeID      *uuid.UUID `json:"fee_type_id" validate:"omitempty,uuid"`
	AcademicYearID *uuid.UUID `json:"academic_year_id" validate:"omitempty,uuid"`
	Amount         float64    `json:"amount" validate:"required,min=0"`
	DueDate        time.Time  `json:"due_date" validate:"required"`
	Status         *string    `json:"status" validate:"omitempty,oneof=paid unpaid partial overdue"`
	PaymentDate    *time.Time `json:"payment_date,omitempty"`
	PaymentMethod  *string    `json:"payment_method" validate:"omitempty,max=50"`
	Notes          *string    `json:"notes,omitempty"`
}

type UpdateStudentFeeRequest struct {
	Amount        *float64   `json:"amount" validate:"omitempty,min=0"`
	DueDate       *time.Time `json:"due_date,omitempty"`
	Status        *string    `json:"status" validate:"omitempty,oneof=paid unpaid partial overdue"`
	PaymentDate   *time.Time `json:"payment_date,omitempty"`
	PaymentMethod *string    `json:"payment_method" validate:"omitempty,max=50"`
	Notes         *string    `json:"notes,omitempty"`
}

type FeeQueryParams struct {
	QueryParams
	StudentID      *uuid.UUID `query:"student_id" validate:"omitempty,uuid"`
	FeeTypeID      *uuid.UUID `query:"fee_type_id" validate:"omitempty,uuid"`
	AcademicYearID *uuid.UUID `query:"academic_year_id" validate:"omitempty,uuid"`
	Status         *string    `query:"status" validate:"omitempty,oneof=paid unpaid partial overdue"`
}
