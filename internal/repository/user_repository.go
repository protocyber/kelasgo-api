package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

// UserRepository interface defines user repository methods
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	GetByEmailGlobal(email string) (*model.User, error) // Global email lookup without tenant context
	GetByUsernameAndTenant(username string, tenantID uuid.UUID) (*model.User, error)
	GetByEmailAndTenant(email string, tenantID uuid.UUID) (*model.User, error)
	GetUserTenants(userID uuid.UUID) ([]model.TenantUser, error) // Get all tenants for a user
	Update(user *model.User) error
	Delete(id uuid.UUID) error
	BulkDelete(ids []uuid.UUID) error
	List(offset, limit int, search string) ([]model.User, int64, error)
	GetUsersByTenant(tenantID uuid.UUID, offset, limit int, search string) ([]model.User, int64, error)
	GetUsersByRole(roleID uuid.UUID, offset, limit int) ([]model.User, int64, error)
	GetByRole(tenantID uuid.UUID, roleID uuid.UUID, offset, limit int) ([]model.User, int64, error)
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

func (r *userRepository) Create(user *model.User) error {
	err := r.db.Write.Create(user).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("username", user.Username).
			Str("email", user.Email).
			Msg("Failed to create user in database")
	}
	return err
}

func (r *userRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Read.Preload("TenantUsers").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("user_id", id.String()).
				Msg("User not found in database")
			return nil, errors.New("user not found")
		}
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Database error while getting user by ID")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
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

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
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

func (r *userRepository) GetByEmailGlobal(email string) (*model.User, error) {
	var user model.User
	err := r.db.Read.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("email", email).
				Msg("User not found by email globally")
			return nil, errors.New("user not found")
		}
		log.Error().
			Err(err).
			Str("email", email).
			Msg("Database error while getting user by email globally")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetUserTenants(userID uuid.UUID) ([]model.TenantUser, error) {
	var tenantUsers []model.TenantUser
	err := r.db.Read.Preload("Tenant").Where("user_id = ? AND is_active = true", userID).Find(&tenantUsers).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to get user tenants from database")
		return nil, err
	}
	return tenantUsers, nil
}

func (r *userRepository) Update(user *model.User) error {
	err := r.db.Write.Save(user).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", user.ID.String()).
			Str("username", user.Username).
			Msg("Failed to update user in database")
	}
	return err
}

func (r *userRepository) Delete(id uuid.UUID) error {
	err := r.db.Write.Delete(&model.User{}, id).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Msg("Failed to delete user from database")
	}
	return err
}

func (r *userRepository) BulkDelete(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Write.Where("id IN (?)", ids).Delete(&model.User{}).Error
	if err != nil {
		log.Error().
			Err(err).
			Interface("ids", ids).
			Msg("Failed to bulk delete users from database")
	}
	return err
}

func (r *userRepository) List(offset, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers")

	if search != "" {
		query = query.Where("full_name ILIKE ? OR username ILIKE ? OR email ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		log.Error().
			Err(err).
			Str("search", search).
			Msg("Failed to count users in List method")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		log.Error().
			Err(err).
			Int("offset", offset).
			Int("limit", limit).
			Str("search", search).
			Msg("Failed to list users from database in List method")
	}
	return users, total, err
}

func (r *userRepository) GetUsersByTenant(tenantID uuid.UUID, offset, limit int, search string) ([]model.User, int64, error) {
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

func (r *userRepository) GetUsersByRole(roleID uuid.UUID, offset, limit int) ([]model.User, int64, error) {
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

func (r *userRepository) GetByUsernameAndTenant(username string, tenantID uuid.UUID) (*model.User, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		log.Error().
			Err(err).
			Str("username", username).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to set tenant context for GetByUsernameAndTenant")
		return nil, err
	}

	var user model.User
	err := r.db.Read.Preload("TenantUsers").
		Joins("JOIN tenant_users ON users.id = tenant_users.user_id").
		Where("users.username = ? AND tenant_users.tenant_id = ? AND tenant_users.is_active = true", username, tenantID).
		First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("username", username).
				Str("tenant_id", tenantID.String()).
				Msg("User not found by username and tenant")
			return nil, errors.New("user not found")
		}
		log.Error().
			Err(err).
			Str("username", username).
			Str("tenant_id", tenantID.String()).
			Msg("Database error in GetByUsernameAndTenant")
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByEmailAndTenant(email string, tenantID uuid.UUID) (*model.User, error) {
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

func (r *userRepository) GetByRole(tenantID uuid.UUID, roleID uuid.UUID, offset, limit int) ([]model.User, int64, error) {
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
