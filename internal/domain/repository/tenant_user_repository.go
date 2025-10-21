package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"gorm.io/gorm"
)

// TenantUserRepository interface defines tenant user repository methods
type TenantUserRepository interface {
	Create(c context.Context, tenantUser *model.TenantUser) error
	GetByID(c context.Context, id uuid.UUID) (*model.TenantUser, error)
	GetByTenantAndUser(c context.Context, tenantID, userID uuid.UUID) (*model.TenantUser, error)
	GetByTenant(c context.Context, tenantID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error)
	GetByUser(c context.Context, userID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error)
	Update(c context.Context, tenantUser *model.TenantUser) error
	Delete(c context.Context, id uuid.UUID) error
	BulkDelete(c context.Context, ids []uuid.UUID) error
	ActivateUser(c context.Context, tenantID, userID uuid.UUID) error
	DeactivateUser(c context.Context, tenantID, userID uuid.UUID) error
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

func (r *tenantUserRepository) Create(c context.Context, tenantUser *model.TenantUser) error {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantUser.TenantID); err != nil {
		return err
	}
	err := r.db.Write.Create(tenantUser).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "create_tenant_user").
			Msg("Database write operation failed")
	}
	return err
}

func (r *tenantUserRepository) GetByID(c context.Context, id uuid.UUID) (*model.TenantUser, error) {
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

func (r *tenantUserRepository) GetByTenantAndUser(c context.Context, tenantID, userID uuid.UUID) (*model.TenantUser, error) {
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

func (r *tenantUserRepository) GetByTenant(c context.Context, tenantID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error) {
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

func (r *tenantUserRepository) GetByUser(c context.Context, userID uuid.UUID, offset, limit int) ([]model.TenantUser, int64, error) {
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

func (r *tenantUserRepository) Update(c context.Context, tenantUser *model.TenantUser) error {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantUser.TenantID); err != nil {
		return err
	}
	err := r.db.Write.Save(tenantUser).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "update_tenant_user").
			Msg("Database write operation failed")
	}
	return err
}

func (r *tenantUserRepository) Delete(c context.Context, id uuid.UUID) error {
	// repoCtx := r.WithContext(c)
	return r.db.Write.Delete(&model.TenantUser{}, id).Error
}

func (r *tenantUserRepository) BulkDelete(c context.Context, ids []uuid.UUID) error {
	repoCtx := r.WithContext(c)
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Write.Where("id IN (?)", ids).Delete(&model.TenantUser{}).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "bulk_delete_tenant_users").
			Int("count", len(ids)).
			Msg("Database write operation failed")
	}
	return err
}

func (r *tenantUserRepository) ActivateUser(c context.Context, tenantID, userID uuid.UUID) error {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return err
	}

	return r.db.Write.Model(&model.TenantUser{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Update("is_active", true).Error
}

func (r *tenantUserRepository) DeactivateUser(c context.Context, tenantID, userID uuid.UUID) error {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return err
	}

	return r.db.Write.Model(&model.TenantUser{}).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		Update("is_active", false).Error
}
