package dto

import "time"

// Parent DTOs
type CreateParentRequest struct {
	FullName     string     `json:"full_name" validate:"required,max=100"`
	Phone        *string    `json:"phone" validate:"omitempty,max=20"`
	Email        *string    `json:"email" validate:"omitempty,email,max=100"`
	Address      *string    `json:"address,omitempty"`
	Relationship *string    `json:"relationship" validate:"omitempty,max=50"`
	Birthplace   *string    `json:"birthplace" validate:"omitempty,max=100"`
	Birthday     *time.Time `json:"birthday,omitempty"`
	Gender       *Gender    `json:"gender" validate:"omitempty,oneof=male female"`
}

type UpdateParentRequest struct {
	FullName     *string    `json:"full_name" validate:"omitempty,max=100"`
	Phone        *string    `json:"phone" validate:"omitempty,max=20"`
	Email        *string    `json:"email" validate:"omitempty,email,max=100"`
	Address      *string    `json:"address,omitempty"`
	Relationship *string    `json:"relationship" validate:"omitempty,max=50"`
	Birthplace   *string    `json:"birthplace" validate:"omitempty,max=100"`
	Birthday     *time.Time `json:"birthday,omitempty"`
	Gender       *Gender    `json:"gender" validate:"omitempty,oneof=male female"`
}
