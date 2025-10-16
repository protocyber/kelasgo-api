package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// TenantUserRepository interface defines tenant user repository methods
type TenantUserRepository interface {
	Create(tenantUser *model.TenantUser) error
	GetByID(id uuid.UUID) (*model.TenantUser, error)
	GetByTenantAndUser(tenantID, userID uuid.UUID) (*model.TenantUser, error)
	GetByTenant(tenantID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error)
	GetByUser(userID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error)
	Update(tenantUser *model.TenantUser) error
	Delete(id uuid.UUID) error
	BulkDelete(ids []uuid.UUID) error
	ActivateUser(tenantID, userID uuid.UUID) error
	DeactivateUser(tenantID, userID uuid.UUID) error
}

// tenantUserRepository implements TenantUserRepository
type tenantUserRepository struct {
	*BaseRepository
}

// NewTenantUserRepository creates a new tenant user repository
func NewTenantUserRepository(db *database.DatabaseConnections) TenantUserRepository {
	return &tenantUserRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *tenantUserRepository) Create(tenantUser *model.TenantUser) error {
	if err := r.SetTenantContext(tenantUser.TenantID); err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantUser.TenantID.String()).
			Str("user_id", tenantUser.UserID.String()).
			Msg("Failed to set tenant context for tenant-user creation")
		return err
	}
	err := r.db.Write.Create(tenantUser).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantUser.TenantID.String()).
			Str("user_id", tenantUser.UserID.String()).
			Msg("Failed to create tenant-user relationship in database")
	}
	return err
}

func (r *tenantUserRepository) GetByID(id uuid.UUID) (*model.TenantUser, error) {
	var tenantUser model.TenantUser
	err := r.db.Read.Preload("User").Preload("Teacher").Preload("Student").First(&tenantUser, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant user not found")
		}
		return nil, err
	}
	return &tenantUser, nil
}

func (r *tenantUserRepository) GetByTenantAndUser(tenantID, userID uuid.UUID) (*model.TenantUser, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, err
	}

	var tenantUser model.TenantUser
	err := r.db.Read.Preload("User").Preload("Teacher").Preload("Student").
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).First(&tenantUser).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("tenant user not found")
		}
		return nil, err
	}
	return &tenantUser, nil
}

func (r *tenantUserRepository) GetByTenant(tenantID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var tenantUsers []model.TenantUser
	var count int64

	query := r.db.Read.Model(&model.TenantUser{}).Where("tenant_id = ?", tenantID)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Teacher").Preload("Student").
		Offset(offset).Limit(limit).Find(&tenantUsers).Error; err != nil {
		return nil, 0, err
	}

	return tenantUsers, count, nil
}

func (r *tenantUserRepository) GetByUser(userID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error) {
	var tenantUsers []model.TenantUser
	var count int64

	query := r.db.Read.Model(&model.TenantUser{}).Where("user_id = ?", userID)

	if err := query.Count(&count).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Preload("Teacher").Preload("Student").
		Offset(offset).Limit(limit).Find(&tenantUsers).Error; err != nil {
		return nil, 0, err
	}

	return tenantUsers, count, nil
}

func (r *tenantUserRepository) Update(tenantUser *model.TenantUser) error {
	if err := r.SetTenantContext(tenantUser.TenantID); err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantUser.TenantID.String()).
			Str("user_id", tenantUser.UserID.String()).
			Msg("Failed to set tenant context for tenant-user update")
		return err
	}
	err := r.db.Write.Save(tenantUser).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantUser.TenantID.String()).
			Str("user_id", tenantUser.UserID.String()).
			Msg("Failed to update tenant-user relationship in database")
	}
	return err
}

func (r *tenantUserRepository) Delete(id uuid.UUID) error {
	return r.db.Write.Delete(&model.TenantUser{}, id).Error
}

func (r *tenantUserRepository) BulkDelete(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Write.Where("id IN (?)", ids).Delete(&model.TenantUser{}).Error
	if err != nil {
		log.Error().
			Err(err).
			Interface("ids", ids).
			Msg("Failed to bulk delete tenant users from database")
	}
	return err
}

func (r *tenantUserRepository) ActivateUser(tenantID, userID uuid.UUID) error {
	if err := r.SetTenantContext(tenantID); err != nil {
		return err
	}

	return r.db.Write.Model(&model.TenantUser{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Update("is_active", true).Error
}

func (r *tenantUserRepository) DeactivateUser(tenantID, userID uuid.UUID) error {
	if err := r.SetTenantContext(tenantID); err != nil {
		return err
	}

	return r.db.Write.Model(&model.TenantUser{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Update("is_active", false).Error
}
