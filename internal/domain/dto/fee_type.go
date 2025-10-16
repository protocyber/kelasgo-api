package dto

// Fee Type DTOs
type CreateFeeTypeRequest struct {
	Name          string   `json:"name" validate:"required,max=100"`
	Description   *string  `json:"description,omitempty"`
	DefaultAmount *float64 `json:"default_amount,omitempty" validate:"omitempty,min=0"`
	IsMandatory   *bool    `json:"is_mandatory,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}

type UpdateFeeTypeRequest struct {
	Name          *string  `json:"name" validate:"omitempty,max=100"`
	Description   *string  `json:"description,omitempty"`
	DefaultAmount *float64 `json:"default_amount,omitempty" validate:"omitempty,min=0"`
	IsMandatory   *bool    `json:"is_mandatory,omitempty"`
	IsActive      *bool    `json:"is_active,omitempty"`
}
