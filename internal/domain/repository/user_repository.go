package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"gorm.io/gorm"
)

// UserRepository interface defines user repository methods
type UserRepository interface {
	Create(c context.Context, user *model.User) error
	GetByID(c context.Context, id uuid.UUID) (*model.User, error)
	GetByUsername(c context.Context, username string) (*model.User, error)
	GetByEmail(c context.Context, email string) (*model.User, error)
	GetByEmailGlobal(c context.Context, email string) (*model.User, error) // Global email lookup without tenant context
	GetByUsernameAndTenant(c context.Context, username string, tenantID uuid.UUID) (*model.User, error)
	GetByEmailAndTenant(c context.Context, email string, tenantID uuid.UUID) (*model.User, error)
	GetUserTenants(c context.Context, userID uuid.UUID) ([]model.TenantUser, error) // Get all tenants for a user
	Update(c context.Context, user *model.User) error
	Delete(c context.Context, id uuid.UUID) error
	BulkDelete(c context.Context, ids []uuid.UUID) error
	List(c context.Context, offset, limit int, search string) ([]model.User, int64, error)
	GetUsersByTenant(c context.Context, tenantID uuid.UUID, offset, limit int, search string) ([]model.User, int64, error)
	GetUsersByRole(c context.Context, roleID uuid.UUID, offset, limit int) ([]model.User, int64, error)
	GetByRole(c context.Context, tenantID uuid.UUID, roleID uuid.UUID, offset, limit int) ([]model.User, int64, error)
}

// userRepository implements UserRepository
type userRepository struct {
	*BaseRepository
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *database.DatabaseConnections) UserRepository {
	return &userRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *userRepository) Create(c context.Context, user *model.User) error {
	repoCtx := r.WithContext(c)
	err := r.db.Write.Create(user).Error
	if err != nil {
		repoCtx.GetLogger().Error().
			Err(err).
			Str("operation", "create_user").
			Msg("Database write operation failed")
	}
	return err
}

func (r *userRepository) GetByID(c context.Context, id uuid.UUID) (*model.User, error) {
	repoCtx := r.WithContext(c)
	var user model.User
	err := r.db.Read.Preload("TenantUsers").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		repoCtx.GetLogger().Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Database error while getting user by ID")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(c context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.Read.Preload("TenantUsers").Where("username = ?", username).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmail(c context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.Read.Preload("TenantUsers").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmailGlobal(c context.Context, email string) (*model.User, error) {
	repoCtx := r.WithContext(c)
	var user model.User
	err := r.db.Read.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		repoCtx.logger.Error().
			Err(err).
			Str("email", email).
			Msg("Database error while getting user by email (global)")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserTenants(c context.Context, userID uuid.UUID) ([]model.TenantUser, error) {
	repoCtx := r.WithContext(c)
	var tenantUsers []model.TenantUser
	err := r.db.Read.Preload("Tenant").Where("user_id = ? AND is_active = true", userID).Find(&tenantUsers).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "get_user_tenants").
			Msg("Database query failed")
		return nil, err
	}
	return tenantUsers, nil
}

func (r *userRepository) Update(c context.Context, user *model.User) error {
	repoCtx := r.WithContext(c)
	err := r.db.Write.Save(user).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "update_user").
			Msg("Database write operation failed")
	}
	return err
}

func (r *userRepository) Delete(c context.Context, id uuid.UUID) error {
	repoCtx := r.WithContext(c)
	err := r.db.Write.Delete(&model.User{}, id).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "delete_user").
			Msg("Database write operation failed")
	}
	return err
}

func (r *userRepository) BulkDelete(c context.Context, ids []uuid.UUID) error {
	repoCtx := r.WithContext(c)
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Write.Where("id IN (?)", ids).Delete(&model.User{}).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "bulk_delete_users").
			Int("count", len(ids)).
			Msg("Database write operation failed")
	}
	return err
}

func (r *userRepository) List(c context.Context, offset, limit int, search string) ([]model.User, int64, error) {
	repoCtx := r.WithContext(c)
	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers")

	if search != "" {
		query = query.Where("full_name ILIKE ? OR username ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "count_users").
			Msg("Database query failed")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "list_users").
			Msg("Database query failed")
	}
	return users, total, err
}

func (r *userRepository) GetUsersByTenant(c context.Context, tenantID uuid.UUID, offset, limit int, search string) ([]model.User, int64, error) {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers").
		Joins("JOIN tenant_users ON users.id = tenant_users.user_id").
		Where("tenant_users.tenant_id = ?", tenantID)

	if search != "" {
		query = query.Where("users.full_name ILIKE ? OR users.username ILIKE ? OR users.email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) GetUsersByRole(c context.Context, roleID uuid.UUID, offset, limit int) ([]model.User, int64, error) {
	// repoCtx := r.WithContext(c)
	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers").
		Joins("JOIN tenant_users ON users.id = tenant_users.user_id").
		Joins("JOIN tenant_user_roles ON tenant_users.id = tenant_user_roles.tenant_user_id").
		Where("tenant_user_roles.role_id = ? AND tenant_users.is_active = true", roleID)

	// Get total count
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) GetByUsernameAndTenant(c context.Context, username string, tenantID uuid.UUID) (*model.User, error) {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, err
	}

	var user model.User
	err := r.db.Read.Preload("TenantUsers").
		Joins("JOIN tenant_users ON users.id = tenant_users.user_id").
		Where("users.username = ? AND tenant_users.tenant_id = ? AND tenant_users.is_active = true", username, tenantID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "get_user_by_username_tenant").
			Msg("Database query failed")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmailAndTenant(c context.Context, email string, tenantID uuid.UUID) (*model.User, error) {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, err
	}

	var user model.User
	err := r.db.Read.Preload("TenantUsers").
		Joins("JOIN tenant_users ON users.id = tenant_users.user_id").
		Where("users.email = ? AND tenant_users.tenant_id = ? AND tenant_users.is_active = true", email, tenantID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByRole(c context.Context, tenantID uuid.UUID, roleID uuid.UUID, offset, limit int) ([]model.User, int64, error) {
	// repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers").
		Joins("JOIN tenant_users ON users.id = tenant_users.user_id").
		Joins("JOIN tenant_user_roles ON tenant_users.id = tenant_user_roles.tenant_user_id").
		Where("tenant_users.tenant_id = ? AND tenant_user_roles.role_id = ? AND tenant_users.is_active = true", tenantID, roleID)

	// Get total count
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}
