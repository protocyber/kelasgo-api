package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// TenantUserRoleRepository interface defines tenant user role repository methods
type TenantUserRoleRepository interface {
	Create(tenantUserRole *model.TenantUserRole) error
	GetByTenantUserAndRole(tenantUserID, roleID uuid.UUID) (*model.TenantUserRole, error)
	GetRolesByTenantUser(tenantUserID uuid.UUID) ([]model.TenantUserRole, error)
	GetTenantUsersByRole(roleID uuid.UUID) ([]model.TenantUserRole, error)
	Delete(tenantUserID, roleID uuid.UUID) error
	DeleteAllTenantUserRoles(tenantUserID uuid.UUID) error
}

// tenantUserRoleRepository implements TenantUserRoleRepository
type tenantUserRoleRepository struct {
	*BaseRepository
}

// NewTenantUserRoleRepository creates a new tenant user role repository
func NewTenantUserRoleRepository(db *database.DatabaseConnections) TenantUserRoleRepository {
	return &tenantUserRoleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *tenantUserRoleRepository) Create(tenantUserRole *model.TenantUserRole) error {
	err := r.db.Write.Create(tenantUserRole).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_user_id", tenantUserRole.TenantUserID.String()).
			Str("role_id", tenantUserRole.RoleID.String()).
			Msg("Failed to create tenant user-role relationship in database")
	}
	return err
}

func (r *tenantUserRoleRepository) GetByTenantUserAndRole(tenantUserID, roleID uuid.UUID) (*model.TenantUserRole, error) {
	var tenantUserRole model.TenantUserRole
	err := r.db.Read.Preload("TenantUser").Preload("Role").
		Where("tenant_user_id = ? AND role_id = ?", tenantUserID, roleID).First(&tenantUserRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant user role not found")
		}
		return nil, err
	}
	return &tenantUserRole, nil
}

func (r *tenantUserRoleRepository) GetRolesByTenantUser(tenantUserID uuid.UUID) ([]model.TenantUserRole, error) {
	var tenantUserRoles []model.TenantUserRole
	err := r.db.Read.Preload("Role").Where("tenant_user_id = ?", tenantUserID).Find(&tenantUserRoles).Error
	if err != nil {
		return nil, err
	}
	return tenantUserRoles, nil
}

func (r *tenantUserRoleRepository) GetTenantUsersByRole(roleID uuid.UUID) ([]model.TenantUserRole, error) {
	var tenantUserRoles []model.TenantUserRole
	err := r.db.Read.Preload("TenantUser").Where("role_id = ?", roleID).Find(&tenantUserRoles).Error
	if err != nil {
		return nil, err
	}
	return tenantUserRoles, nil
}

func (r *tenantUserRoleRepository) Delete(tenantUserID, roleID uuid.UUID) error {
	return r.db.Write.Where("tenant_user_id = ? AND role_id = ?", tenantUserID, roleID).Delete(&model.TenantUserRole{}).Error
}

func (r *tenantUserRoleRepository) DeleteAllTenantUserRoles(tenantUserID uuid.UUID) error {
	err := r.db.Write.Where("tenant_user_id = ?", tenantUserID).Delete(&model.TenantUserRole{}).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_user_id", tenantUserID.String()).
			Msg("Failed to delete all tenant user roles from database")
	}
	return err
}
