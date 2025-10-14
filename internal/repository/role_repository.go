package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// RoleRepository interface defines role repository methods
type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uuid.UUID) (*model.Role, error)
	GetByName(name string, tenantID uuid.UUID) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uuid.UUID) error
	List(tenantID uuid.UUID, offset, limit int, search string) ([]model.Role, int64, error)
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

func (r *roleRepository) Create(role *model.Role) error {
	if err := r.SetTenantContext(role.TenantID); err != nil {
		log.Error().
			Err(err).
			Str("role_name", role.Name).
			Str("tenant_id", role.TenantID.String()).
			Msg("Failed to set tenant context for role creation")
		return err
	}
	err := r.db.Write.Create(role).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("role_name", role.Name).
			Str("tenant_id", role.TenantID.String()).
			Msg("Failed to create role in database")
	}
	return err
}

func (r *roleRepository) GetByID(id uuid.UUID) (*model.Role, error) {
	var role model.Role
	err := r.db.Read.First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("role_id", id.String()).
				Msg("Role not found by ID")
			return nil, errors.New("role not found")
		}
		log.Error().
			Err(err).
			Str("role_id", id.String()).
			Msg("Database error while getting role by ID")
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(name string, tenantID uuid.UUID) (*model.Role, error) {
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

func (r *roleRepository) Update(role *model.Role) error {
	if err := r.SetTenantContext(role.TenantID); err != nil {
		return err
	}
	return r.db.Write.Save(role).Error
}

func (r *roleRepository) Delete(id uuid.UUID) error {
	return r.db.Write.Delete(&model.Role{}, id).Error
}

func (r *roleRepository) List(tenantID uuid.UUID, offset, limit int, search string) ([]model.Role, int64, error) {
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
