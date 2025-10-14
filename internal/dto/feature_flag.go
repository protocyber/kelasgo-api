package dto

import (
	"github.com/google/uuid"
)

// FeatureFlag DTOs
type CreateFeatureFlagRequest struct {
	Code        string  `json:"code" validate:"required,max=100"`
	Name        string  `json:"name" validate:"required,max=255"`
	Description *string `json:"description,omitempty"`
}

type UpdateFeatureFlagRequest struct {
	Code        *string `json:"code,omitempty" validate:"omitempty,max=100"`
	Name        *string `json:"name,omitempty" validate:"omitempty,max=255"`
	Description *string `json:"description,omitempty"`
}

type FeatureFlagQueryParams struct {
	QueryParams
	Code string `query:"code" validate:"omitempty,max=100"`
	Name string `query:"name" validate:"omitempty,max=255"`
}

type FeatureFlagResponse struct {
	ID          uuid.UUID `json:"id"`
	Code        string    `json:"code"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
}
