package dto

import (
	"time"

	"github.com/google/uuid"
)

// Teacher DTOs
type CreateTeacherRequest struct {
	TenantUserID   uuid.UUID  `json:"tenant_user_id" validate:"required,uuid"`
	EmployeeNumber *string    `json:"employee_number" validate:"omitempty,max=50"`
	HireDate       *time.Time `json:"hire_date,omitempty"`
	DepartmentID   *uuid.UUID `json:"department_id" validate:"omitempty,uuid"`
	Qualification  *string    `json:"qualification" validate:"omitempty,max=100"`
	Position       *string    `json:"position" validate:"omitempty,max=100"`
}

type UpdateTeacherRequest struct {
	EmployeeNumber *string    `json:"employee_number" validate:"omitempty,max=50"`
	HireDate       *time.Time `json:"hire_date,omitempty"`
	DepartmentID   *uuid.UUID `json:"department_id" validate:"omitempty,uuid"`
	Qualification  *string    `json:"qualification" validate:"omitempty,max=100"`
	Position       *string    `json:"position" validate:"omitempty,max=100"`
}

type TeacherQueryParams struct {
	QueryParams
	DepartmentID *uuid.UUID `query:"department_id" validate:"omitempty,uuid"`
}
