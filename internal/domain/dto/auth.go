package dto

import (
	"time"

	"github.com/google/uuid"
)

// Auth DTOs
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type LoginResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,min=6"`
	FullName string `json:"full_name" validate:"required,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Phone    string `json:"phone" validate:"omitempty,max=20"`
}

type UserInfo struct {
	ID       uuid.UUID  `json:"id"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"` // Optional, null if no tenant selected
	Username string     `json:"username"`
	Email    string     `json:"email"`
	FullName string     `json:"full_name"`
	Role     string     `json:"role,omitempty"` // Optional, only present when tenant is selected
}

type TokenClaims struct {
	UserID   uuid.UUID  `json:"user_id"`
	TenantID *uuid.UUID `json:"tenant_id,omitempty"` // Optional, null if no tenant selected
	Username string     `json:"username"`
	Email    string     `json:"email"`
	Role     string     `json:"role,omitempty"` // Optional, only present when tenant is selected
}

// Tenant selection DTOs
type TenantSelectionRequest struct {
	TenantID string `json:"tenant_id" validate:"required,uuid"`
}

type TenantSelectionResponse struct {
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
	User         UserInfo  `json:"user"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
}
