package model

import (
	"time"
)

// User represents the users table
type User struct {
	GlobalBaseModel            // Users table doesn't have tenant_id since it's a global table
	Username        string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash    string     `gorm:"size:255;not null" json:"-"`
	Email           string     `gorm:"size:100;uniqueIndex" json:"email"`
	FullName        string     `gorm:"size:100;not null" json:"full_name"`
	Birthplace      *string    `gorm:"size:100" json:"birthplace,omitempty"`
	Birthday        *time.Time `gorm:"type:date" json:"birthday,omitempty"`
	Gender          *Gender    `gorm:"type:gender_enum" json:"gender,omitempty"`
	DateOfBirth     *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Phone           *string    `gorm:"size:20" json:"phone,omitempty"`
	Address         *string    `gorm:"type:text" json:"address,omitempty"`
	IsActive        bool       `gorm:"default:true" json:"is_active"`

	// Relationships
	TenantUsers   []TenantUser   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"tenant_users,omitempty"`
	UserRoles     []UserRole     `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"user_roles,omitempty"`
	Notifications []Notification `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE" json:"notifications,omitempty"`
	AuditLogs     []AuditLog     `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"audit_logs,omitempty"`
}

// TableName returns the table name for User
func (User) TableName() string {
	return "users"
}
