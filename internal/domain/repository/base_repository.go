package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"github.com/protocyber/kelasgo-api/internal/util"
	"gorm.io/gorm"
)

// BaseRepository provides common database operations with tenant context
type BaseRepository struct {
	db     *database.DatabaseConnections
	helper *util.TenantHelper
	logger *util.ContextLogger
}

// NewBaseRepository creates a new base repository
func NewBaseRepository(db *database.DatabaseConnections) *BaseRepository {
	return &BaseRepository{
		db:     db,
		helper: util.NewTenantHelper(db),
	}
}

// WithContext sets the context for the repository operations and creates a context logger
func (r *BaseRepository) WithContext(ctx context.Context) *BaseRepository {
	// Create a copy of the repository with context logger
	return &BaseRepository{
		db:     r.db,
		helper: r.helper,
		logger: util.NewServiceLogger(ctx),
	}
}

// GetLogger returns the context logger, creating a fallback if none exists
func (r *BaseRepository) GetLogger() *util.ContextLogger {
	if r.logger != nil {
		return r.logger
	}
	// Fallback logger without context - should not happen in normal flow
	return &util.ContextLogger{}
}

// SetTenantContext sets the tenant context for database operations
func (r *BaseRepository) SetTenantContext(tenantID uuid.UUID) error {
	if r.logger != nil {
		r.logger.Debug().
			Str("tenant_id", tenantID.String()).
			Msg("Setting tenant context for repository operation")
	}
	return r.helper.SetTenantContext(tenantID)
}

// ClearTenantContext clears the tenant context
func (r *BaseRepository) ClearTenantContext() error {
	if r.logger != nil {
		r.logger.Debug().Msg("Clearing tenant context for repository operation")
	}
	return r.helper.GetTenantContextManager().ClearTenantContext()
}

// ExecuteWithTenant executes a function with tenant context set
func (r *BaseRepository) ExecuteWithTenant(tenantID uuid.UUID, fn func() error) error {
	if r.logger != nil {
		r.logger.Debug().
			Str("tenant_id", tenantID.String()).
			Msg("Executing repository operation with tenant context")
	}
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
