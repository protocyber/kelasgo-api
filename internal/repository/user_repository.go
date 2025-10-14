package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"gorm.io/gorm"
)

// UserRepository interface defines user repository methods
type UserRepository interface {
	Create(user *model.User) error
	GetByID(id uuid.UUID) (*model.User, error)
	GetByUsername(username string) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
	Update(user *model.User) error
	Delete(id uuid.UUID) error
	List(offset, limit int, search string) ([]model.User, int64, error)
	GetUsersByTenant(tenantID uuid.UUID, offset, limit int, search string) ([]model.User, int64, error)
	GetUsersByRole(roleID uuid.UUID, offset, limit int) ([]model.User, int64, error)
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
	return r.db.Write.Create(user).Error
}

func (r *userRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	err := r.db.Read.Preload("TenantUsers").Preload("UserRoles.Role").First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByUsername(username string) (*model.User, error) {
	var user model.User
	err := r.db.Read.Preload("TenantUsers").Preload("UserRoles.Role").Where("username = ?", username).First(&user).Error
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
	err := r.db.Read.Preload("TenantUsers").Preload("UserRoles.Role").Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Update(user *model.User) error {
	return r.db.Write.Save(user).Error
}

func (r *userRepository) Delete(id uuid.UUID) error {
	return r.db.Write.Delete(&model.User{}, id).Error
}

func (r *userRepository) List(offset, limit int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers").Preload("UserRoles.Role")

	if search != "" {
		query = query.Where("full_name ILIKE ? OR username ILIKE ? OR email ILIKE ?",
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

func (r *userRepository) GetUsersByTenant(tenantID uuid.UUID, offset, limit int, search string) ([]model.User, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var users []model.User
	var total int64

	query := r.db.Read.Preload("TenantUsers").Preload("UserRoles.Role").
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

	query := r.db.Read.Preload("TenantUsers").Preload("UserRoles.Role").
		Joins("JOIN user_roles ON users.id = user_roles.user_id").
		Where("user_roles.role_id = ?", roleID)

	// Get total count
	if err := query.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}
