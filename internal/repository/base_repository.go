package repository

import (
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/util"
	"gorm.io/gorm"
)

// BaseRepository provides common database operations with tenant context
type BaseRepository struct {
	db     *database.DatabaseConnections
	helper *util.TenantHelper
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *database.DatabaseConnections) *BaseRepository {
	return &BaseRepository{
		db:     db,
		helper: util.NewTenantHelper(db),
	}
}

// SetTenantContext sets the tenant context for database operations
func (r *BaseRepository) SetTenantContext(tenantID uuid.UUID) error {
	return r.helper.SetTenantContext(tenantID)
}

// ClearTenantContext clears the tenant context
func (r *BaseRepository) ClearTenantContext() error {
	return r.helper.GetTenantContextManager().ClearTenantContext()
}

// ExecuteWithTenant executes a function with tenant context set
func (r *BaseRepository) ExecuteWithTenant(tenantID uuid.UUID, fn func() error) error {
	return r.helper.ExecuteWithTenant(tenantID, fn)
}

// GetReadDB returns the read database connection
func (r *BaseRepository) GetReadDB() *gorm.DB {
	return r.db.Read
}

// GetWriteDB returns the write database connection
func (r *BaseRepository) GetWriteDB() *gorm.DB {
	return r.db.Write
}
