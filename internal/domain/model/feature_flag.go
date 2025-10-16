package model

import (
	"github.com/google/uuid"
)

// FeatureFlag represents the feature_flags table
type FeatureFlag struct {
	ID          uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4()" json:"id"`
	Code        string    `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Name        string    `gorm:"size:255;not null" json:"name"`
	Description *string   `gorm:"type:text" json:"description,omitempty"`

	// Relationships
	TenantFeatures []TenantFeature `gorm:"foreignKey:FeatureID;constraint:OnDelete:CASCADE" json:"tenant_features,omitempty"`
}

// TableName returns the table name for FeatureFlag
func (FeatureFlag) TableName() string {
	return "feature_flags"
}
