package util

import (
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
)

// TenantContextManager helps manage tenant context for database operations
type TenantContextManager struct {
	db *database.DatabaseConnections
}

// NewTenantContextManager creates a new tenant context manager
func NewTenantContextManager(db *database.DatabaseConnections) *TenantContextManager {
	return &TenantContextManager{db: db}
}

// SetTenantContext sets the PostgreSQL session variable for the given tenant
func (tcm *TenantContextManager) SetTenantContext(tenantID uuid.UUID) error {
	sql := "SELECT set_config('app.current_tenant', ?, false)"

	// Set on write connection
	if err := tcm.db.Write.Exec(sql, tenantID.String()).Error; err != nil {
		return err
	}

	// Set on read connection (if different from write)
	if tcm.db.Read != tcm.db.Write {
		if err := tcm.db.Read.Exec(sql, tenantID.String()).Error; err != nil {
			return err
		}
	}

	return nil
}

// ClearTenantContext clears the tenant context
func (tcm *TenantContextManager) ClearTenantContext() error {
	sql := "SELECT set_config('app.current_tenant', '', false)"

	// Clear on write connection
	if err := tcm.db.Write.Exec(sql).Error; err != nil {
		return err
	}

	// Clear on read connection (if different from write)
	if tcm.db.Read != tcm.db.Write {
		if err := tcm.db.Read.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}

// ExecuteWithTenant executes a function with tenant context set
func (tcm *TenantContextManager) ExecuteWithTenant(tenantID uuid.UUID, fn func() error) error {
	// Set tenant context
	if err := tcm.SetTenantContext(tenantID); err != nil {
		return err
	}

	// Execute the function
	return fn()
}
