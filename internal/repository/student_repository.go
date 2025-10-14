package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"gorm.io/gorm"
)

// StudentRepository interface defines student repository methods
type StudentRepository interface {
	Create(student *model.Student) error
	GetByID(id uuid.UUID) (*model.Student, error)
	GetByStudentNumber(studentNumber string, tenantID uuid.UUID) (*model.Student, error)
	GetByTenantUserID(tenantUserID uuid.UUID) (*model.Student, error)
	Update(student *model.Student) error
	Delete(id uuid.UUID) error
	List(tenantID uuid.UUID, offset, limit int, search string) ([]model.Student, int64, error)
	GetByClass(tenantID, classID uuid.UUID, offset, limit int) ([]model.Student, int64, error)
	GetByParent(tenantID, parentID uuid.UUID, offset, limit int) ([]model.Student, int64, error)
}

// studentRepository implements StudentRepository
type studentRepository struct {
	*BaseRepository
}

// NewStudentRepository creates a new student repository
func NewStudentRepository(db *database.DatabaseConnections) StudentRepository {
	return &studentRepository{
		BaseRepository: NewBaseRepository(db),
	}
}

func (r *studentRepository) Create(student *model.Student) error {
	if err := r.SetTenantContext(student.TenantID); err != nil {
		return err
	}
	return r.db.Write.Create(student).Error
}

func (r *studentRepository) GetByID(id uuid.UUID) (*model.Student, error) {
	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").First(&student, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByStudentNumber(studentNumber string, tenantID uuid.UUID) (*model.Student, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, err
	}

	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("student_number = ? AND tenant_id = ?", studentNumber, tenantID).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByTenantUserID(tenantUserID uuid.UUID) (*model.Student, error) {
	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("tenant_user_id = ?", tenantUserID).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(student *model.Student) error {
	if err := r.SetTenantContext(student.TenantID); err != nil {
		return err
	}
	return r.db.Write.Save(student).Error
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	return r.db.Write.Delete(&model.Student{}, id).Error
}

func (r *studentRepository) List(tenantID uuid.UUID, offset, limit int, search string) ([]model.Student, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("students.tenant_id = ?", tenantID)

	if search != "" {
		query = query.Joins("JOIN tenant_users ON tenant_users.id = students.tenant_user_id").
			Joins("JOIN users ON users.id = tenant_users.user_id").
			Where("users.full_name ILIKE ? OR students.student_number ILIKE ?",
				"%"+search+"%", "%"+search+"%")
	}

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	return students, total, err
}

func (r *studentRepository) GetByClass(tenantID, classID uuid.UUID, offset, limit int) ([]model.Student, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("class_id = ? AND tenant_id = ?", classID, tenantID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	return students, total, err
}

func (r *studentRepository) GetByParent(tenantID, parentID uuid.UUID, offset, limit int) ([]model.Student, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("parent_id = ? AND tenant_id = ?", parentID, tenantID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	return students, total, err
}
