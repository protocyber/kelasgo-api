package dto

import (
	"github.com/google/uuid"
)

// TenantFeature DTOs
type CreateTenantFeatureRequest struct {
	TenantID  uuid.UUID `json:"tenant_id" validate:"required,uuid"`
	FeatureID uuid.UUID `json:"feature_id" validate:"required,uuid"`
	Enabled   bool      `json:"enabled"`
}

type UpdateTenantFeatureRequest struct {
	Enabled *bool `json:"enabled,omitempty"`
}

type TenantFeatureQueryParams struct {
	QueryParams
	TenantID  uuid.UUID `query:"tenant_id" validate:"omitempty,uuid"`
	FeatureID uuid.UUID `query:"feature_id" validate:"omitempty,uuid"`
	Enabled   *bool     `query:"enabled" validate:"omitempty"`
}

type TenantFeatureResponse struct {
	TenantID  uuid.UUID `json:"tenant_id"`
	FeatureID uuid.UUID `json:"feature_id"`
	Enabled   bool      `json:"enabled"`
}
