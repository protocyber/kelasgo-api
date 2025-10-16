package dto

import (
	"github.com/google/uuid"
)

// TenantUser DTOs
type CreateTenantUserRequest struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required,uuid"`
	UserID   uuid.UUID `json:"user_id" validate:"required,uuid"`
	IsActive *bool     `json:"is_active,omitempty"`
}

type UpdateTenantUserRequest struct {
	IsActive *bool `json:"is_active,omitempty"`
}

type TenantUserQueryParams struct {
	QueryParams
	TenantID uuid.UUID `query:"tenant_id" validate:"omitempty,uuid"`
	UserID   uuid.UUID `query:"user_id" validate:"omitempty,uuid"`
	IsActive *bool     `query:"is_active"`
}

type TenantUserResponse struct {
	ID        uuid.UUID `json:"id"`
	TenantID  uuid.UUID `json:"tenant_id"`
	UserID    uuid.UUID `json:"user_id"`
	IsActive  bool      `json:"is_active"`
	CreatedAt string    `json:"created_at"`
}
