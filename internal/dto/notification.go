package dto

import (
	"github.com/google/uuid"
)

// Notification DTOs
type CreateNotificationRequest struct {
	UserID  *uuid.UUID `json:"user_id" validate:"omitempty,uuid"`
	Title   string     `json:"title" validate:"required,max=100"`
	Message string     `json:"message" validate:"required"`
}

type UpdateNotificationRequest struct {
	Title   *string `json:"title" validate:"omitempty,max=100"`
	Message *string `json:"message,omitempty"`
	IsRead  *bool   `json:"is_read,omitempty"`
}
