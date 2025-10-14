package model

import (
	"time"
)

// User represents the users table
type User struct {
	BaseModel
	Username     string     `gorm:"size:50;uniqueIndex;not null" json:"username"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Email        string     `gorm:"size:100;uniqueIndex" json:"email"`
	FullName     string     `gorm:"size:100;not null" json:"full_name"`
	Gender       *string    `gorm:"size:10;check:gender IN ('Male', 'Female')" json:"gender,omitempty"`
	DateOfBirth  *time.Time `gorm:"type:date" json:"date_of_birth,omitempty"`
	Phone        *string    `gorm:"size:20" json:"phone,omitempty"`
	Address      *string    `gorm:"type:text" json:"address,omitempty"`
	IsActive     bool       `gorm:"default:true" json:"is_active"`

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
