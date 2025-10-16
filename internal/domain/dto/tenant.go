package dto

import (
	"github.com/google/uuid"
)

// Tenant DTOs
type CreateTenantRequest struct {
	Name               string     `json:"name" validate:"required,max=255"`
	Domain             *string    `json:"domain,omitempty" validate:"omitempty,max=255"`
	ContactEmail       *string    `json:"contact_email,omitempty" validate:"omitempty,email,max=255"`
	Phone              *string    `json:"phone,omitempty" validate:"omitempty,max=50"`
	Address            *string    `json:"address,omitempty"`
	LogoURL            *string    `json:"logo_url,omitempty" validate:"omitempty,url,max=255"`
	PlanID             *uuid.UUID `json:"plan_id,omitempty" validate:"omitempty,uuid"`
	SubscriptionStatus string     `json:"subscription_status,omitempty" validate:"omitempty,oneof=active inactive cancelled expired trial"`
}

type UpdateTenantRequest struct {
	Name               *string    `json:"name,omitempty" validate:"omitempty,max=255"`
	Domain             *string    `json:"domain,omitempty" validate:"omitempty,max=255"`
	ContactEmail       *string    `json:"contact_email,omitempty" validate:"omitempty,email,max=255"`
	Phone              *string    `json:"phone,omitempty" validate:"omitempty,max=50"`
	Address            *string    `json:"address,omitempty"`
	LogoURL            *string    `json:"logo_url,omitempty" validate:"omitempty,url,max=255"`
	PlanID             *uuid.UUID `json:"plan_id,omitempty" validate:"omitempty,uuid"`
	SubscriptionStatus *string    `json:"subscription_status,omitempty" validate:"omitempty,oneof=active inactive cancelled expired trial"`
}

type TenantQueryParams struct {
	QueryParams
	Name               string `query:"name" validate:"omitempty,max=255"`
	Domain             string `query:"domain" validate:"omitempty,max=255"`
	SubscriptionStatus string `query:"subscription_status" validate:"omitempty,oneof=active inactive cancelled expired trial"`
	PlanID             string `query:"plan_id" validate:"omitempty,uuid"`
}

type TenantResponse struct {
	ID                 uuid.UUID  `json:"id"`
	Name               string     `json:"name"`
	Domain             *string    `json:"domain,omitempty"`
	ContactEmail       *string    `json:"contact_email,omitempty"`
	Phone              *string    `json:"phone,omitempty"`
	Address            *string    `json:"address,omitempty"`
	LogoURL            *string    `json:"logo_url,omitempty"`
	PlanID             *uuid.UUID `json:"plan_id,omitempty"`
	SubscriptionStatus string     `json:"subscription_status"`
	CreatedAt          string     `json:"created_at"`
	CreatedBy          *uuid.UUID `json:"created_by,omitempty"`
}
