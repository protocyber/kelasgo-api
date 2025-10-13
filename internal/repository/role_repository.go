package repository

import (
	"errors"

	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"gorm.io/gorm"
)

// RoleRepository interface defines role repository methods
type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List(offset, limit int, search string) ([]model.Role, int64, error)
}

// roleRepository implements RoleRepository
type roleRepository struct {
	db *database.DatabaseConnections
}

// NewRoleRepository creates a new role repository
func NewRoleRepository(db *database.DatabaseConnections) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Write.Create(role).Error
}

func (r *roleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Read.First(&role, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Read.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("role not found")
		}
		return nil, err
	}
	return &role, nil
}

func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Write.Save(role).Error
}

func (r *roleRepository) Delete(id uint) error {
	return r.db.Write.Delete(&model.Role{}, id).Error
}

func (r *roleRepository) List(offset, limit int, search string) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.Read.Model(&model.Role{})

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
