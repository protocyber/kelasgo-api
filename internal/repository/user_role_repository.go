package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"gorm.io/gorm"
)

// UserRoleRepository interface defines user role repository methods
type UserRoleRepository interface {
	Create(userRole *model.UserRole) error
	GetByUserAndRole(userID, roleID uuid.UUID) (*model.UserRole, error)
	GetRolesByUser(userID uuid.UUID) ([]model.UserRole, error)
	GetUsersByRole(roleID uuid.UUID) ([]model.UserRole, error)
	Delete(userID, roleID uuid.UUID) error
	DeleteAllUserRoles(userID uuid.UUID) error
}

// userRoleRepository implements UserRoleRepository
type userRoleRepository struct {
	*BaseRepository
}

// NewUserRoleRepository creates a new user role repository
func NewUserRoleRepository(db *database.DatabaseConnections) UserRoleRepository {
	return &userRoleRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *userRoleRepository) Create(userRole *model.UserRole) error {
	return r.db.Write.Create(userRole).Error
}

func (r *userRoleRepository) GetByUserAndRole(userID, roleID uuid.UUID) (*model.UserRole, error) {
	var userRole model.UserRole
	err := r.db.Read.Preload("User").Preload("Role").
		Where("user_id = ? AND role_id = ?", userID, roleID).First(&userRole).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user role not found")
		}
		return nil, err
	}
	return &userRole, nil
}

func (r *userRoleRepository) GetRolesByUser(userID uuid.UUID) ([]model.UserRole, error) {
	var userRoles []model.UserRole
	err := r.db.Read.Preload("Role").Where("user_id = ?", userID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (r *userRoleRepository) GetUsersByRole(roleID uuid.UUID) ([]model.UserRole, error) {
	var userRoles []model.UserRole
	err := r.db.Read.Preload("User").Where("role_id = ?", roleID).Find(&userRoles).Error
	if err != nil {
		return nil, err
	}
	return userRoles, nil
}

func (r *userRoleRepository) Delete(userID, roleID uuid.UUID) error {
	return r.db.Write.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&model.UserRole{}).Error
}

func (r *userRoleRepository) DeleteAllUserRoles(userID uuid.UUID) error {
	return r.db.Write.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error
}
