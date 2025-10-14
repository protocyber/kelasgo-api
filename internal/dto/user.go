package dto

import (
	"time"

	"github.com/google/uuid"
)

// User DTOs
type CreateUserRequest struct {
	Username    string     `json:"username" validate:"required,min=3,max=50"`
	Password    string     `json:"password" validate:"required,min=6"`
	Email       string     `json:"email" validate:"omitempty,email,max=100"`
	FullName    string     `json:"full_name" validate:"required,max=100"`
	Gender      *string    `json:"gender" validate:"omitempty,oneof=Male Female"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Phone       *string    `json:"phone" validate:"omitempty,max=20"`
	Address     *string    `json:"address,omitempty"`
	RoleID      *uuid.UUID `json:"role_id,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

type UpdateUserRequest struct {
	Email       *string    `json:"email" validate:"omitempty,email,max=100"`
	FullName    *string    `json:"full_name" validate:"omitempty,max=100"`
	Gender      *string    `json:"gender" validate:"omitempty,oneof=Male Female"`
	DateOfBirth *time.Time `json:"date_of_birth,omitempty"`
	Phone       *string    `json:"phone" validate:"omitempty,max=20"`
	Address     *string    `json:"address,omitempty"`
	RoleID      *uuid.UUID `json:"role_id,omitempty"`
	IsActive    *bool      `json:"is_active,omitempty"`
}

type UserQueryParams struct {
	QueryParams
	RoleID   *uuid.UUID `query:"role_id"`
	IsActive *bool      `query:"is_active"`
}
