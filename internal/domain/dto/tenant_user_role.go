package dto

import (
	"github.com/google/uuid"
)

// TenantUserRole DTOs
type CreateTenantUserRoleRequest struct {
	TenantUserID uuid.UUID `json:"tenant_user_id" validate:"required,uuid"`
	RoleID       uuid.UUID `json:"role_id" validate:"required,uuid"`
}

type TenantUserRoleQueryParams struct {
	QueryParams
	TenantUserID uuid.UUID `query:"tenant_user_id" validate:"omitempty,uuid"`
	RoleID       uuid.UUID `query:"role_id" validate:"omitempty,uuid"`
}

type TenantUserRoleResponse struct {
	TenantUserID uuid.UUID `json:"tenant_user_id"`
	RoleID       uuid.UUID `json:"role_id"`
}
