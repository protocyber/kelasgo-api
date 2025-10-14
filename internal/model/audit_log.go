package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// AuditLog represents the audit_logs table
type AuditLog struct {
	ID        int64            `gorm:"primaryKey;autoIncrement" json:"id"`
	TenantID  uuid.UUID        `gorm:"type:uuid;not null;index" json:"tenant_id"`
	UserID    *uuid.UUID       `gorm:"type:uuid;index" json:"user_id,omitempty"`
	Table     string           `gorm:"column:table_name;size:255;not null" json:"table_name"`
	RecordID  *uuid.UUID       `gorm:"type:uuid" json:"record_id,omitempty"`
	Action    string           `gorm:"size:50;not null;check:action IN ('INSERT','UPDATE','DELETE')" json:"action"`
	OldData   *json.RawMessage `gorm:"type:jsonb" json:"old_data,omitempty"`
	NewData   *json.RawMessage `gorm:"type:jsonb" json:"new_data,omitempty"`
	CreatedAt time.Time        `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`

	// Relationships
	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

// TableName returns the table name for AuditLog
func (AuditLog) TableName() string {
	return "audit_logs"
}
