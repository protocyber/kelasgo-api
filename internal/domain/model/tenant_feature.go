package model

import (
	"github.com/google/uuid"
)

// TenantFeature represents the tenant_features table (many-to-many relationship between tenants and feature_flags)
type TenantFeature struct {
	TenantID  uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"tenant_id"`
	FeatureID uuid.UUID `gorm:"type:uuid;not null;primaryKey" json:"feature_id"`
	Enabled   bool      `gorm:"default:true" json:"enabled"`

	// Relationships
	Tenant  *Tenant      `gorm:"foreignKey:TenantID;constraint:OnDelete:CASCADE" json:"tenant,omitempty"`
	Feature *FeatureFlag `gorm:"foreignKey:FeatureID;constraint:OnDelete:CASCADE" json:"feature,omitempty"`
}

// TableName returns the table name for TenantFeature
func (TenantFeature) TableName() string {
	return "tenant_features"
}
