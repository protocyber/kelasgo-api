package dto

import (
	"github.com/google/uuid"
)

// UserRole DTOs
type CreateUserRoleRequest struct {
	UserID uuid.UUID `json:"user_id" validate:"required,uuid"`
	RoleID uuid.UUID `json:"role_id" validate:"required,uuid"`
}

type UserRoleQueryParams struct {
	QueryParams
	UserID uuid.UUID `query:"user_id" validate:"omitempty,uuid"`
	RoleID uuid.UUID `query:"role_id" validate:"omitempty,uuid"`
}

type UserRoleResponse struct {
	UserID uuid.UUID `json:"user_id"`
	RoleID uuid.UUID `json:"role_id"`
}
