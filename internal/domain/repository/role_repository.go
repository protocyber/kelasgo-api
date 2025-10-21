package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"gorm.io/gorm"
)

// RoleRepository interface defines role repository methods
type RoleRepository interface {
	Create(ctx context.Context, role *model.Role) error
	GetByID(ctx context.Context, id uuid.UUID) (*model.Role, error)
	GetByName(ctx context.Context, name string, tenantID uuid.UUID) (*model.Role, error)
	Update(ctx context.Context, role *model.Role) error
	Delete(ctx context.Context, id uuid.UUID) error
	List(ctx context.Context, tenantID uuid.UUID, offset, limit int, search string) ([]model.Role, int64, error)
}

// roleRepository implements RoleRepository
type roleRepository struct {
	*BaseRepository
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *database.DatabaseConnections) RoleRepository {
	return &roleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *roleRepository) Create(c context.Context, role *model.Role) error {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(role.TenantID); err != nil {
		return err
	}
	err := r.db.Write.Create(role).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "create_role").
			Msg("Database write operation failed")
	}
	return err
}

func (r *roleRepository) GetByID(c context.Context, id uuid.UUID) (*model.Role, error) {
	repoCtx := r.WithContext(c)
	var role model.Role
	err := r.db.Read.First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		repoCtx.logger.Error().
			Err(err).
			Str("role_id", id.String()).
			Msg("Database error while getting role by ID")
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(c context.Context, name string, tenantID uuid.UUID) (*model.Role, error) {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, err
	}

	var role model.Role
	err := r.db.Read.Where("name = ? AND tenant_id = ?", name, tenantID).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(c context.Context, role *model.Role) error {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(role.TenantID); err != nil {
		return err
	}
	return r.db.Write.Save(role).Error
}

func (r *roleRepository) Delete(c context.Context, id uuid.UUID) error {
	// repoCtx := r.WithContext(c)
	return r.db.Write.Delete(&model.Role{}, id).Error
}

func (r *roleRepository) List(c context.Context, tenantID uuid.UUID, offset, limit int, search string) ([]model.Role, int64, error) {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var roles []model.Role
	var total int64

	query := r.db.Read.Model(&model.Role{}).Where("tenant_id = ?", tenantID)

	if search != "" {
		query = query.Where("name ILIKE ? OR description ILIKE ?",
			"%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&roles).Error
	return roles, total, err
}
