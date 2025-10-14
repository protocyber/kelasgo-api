package util

import (
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
)

// TenantHelper provides utility functions for tenant operations
type TenantHelper struct {
	db  *database.DatabaseConnections
	tcm *TenantContextManager
}

// NewTenantHelper creates a new tenant helper
func NewTenantHelper(db *database.DatabaseConnections) *TenantHelper {
	return &TenantHelper{
		db:  db,
		tcm: NewTenantContextManager(db),
	}
}

// SetTenantContext sets tenant context using tenant ID
func (th *TenantHelper) SetTenantContext(tenantID uuid.UUID) error {
	return th.tcm.SetTenantContext(tenantID)
}

// GetTenantContextManager returns the tenant context manager for advanced operations
func (th *TenantHelper) GetTenantContextManager() *TenantContextManager {
	return th.tcm
}

// ExecuteWithTenant executes a function with tenant context
func (th *TenantHelper) ExecuteWithTenant(tenantID uuid.UUID, fn func() error) error {
	return th.tcm.ExecuteWithTenant(tenantID, fn)
}
