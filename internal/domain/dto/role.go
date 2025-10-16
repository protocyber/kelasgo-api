package dto

// Role DTOs
type CreateRoleRequest struct {
	Name        string  `json:"name" validate:"required,max=50"`
	Description *string `json:"description,omitempty"`
}

type UpdateRoleRequest struct {
	Name        *string `json:"name" validate:"omitempty,max=50"`
	Description *string `json:"description,omitempty"`
}
