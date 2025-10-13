package repository

import (
	"errors"

	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/model"
	"gorm.io/gorm"
)

// StudentRepository interface defines student repository methods
type StudentRepository interface {
	Create(student *model.Student) error
	GetByID(id uint) (*model.Student, error)
	GetByStudentNumber(studentNumber string) (*model.Student, error)
	GetByUserID(userID uint) (*model.Student, error)
	Update(student *model.Student) error
	Delete(id uint) error
	List(offset, limit int, search string) ([]model.Student, int64, error)
	GetByClass(classID uint, offset, limit int) ([]model.Student, int64, error)
	GetByParent(parentID uint, offset, limit int) ([]model.Student, int64, error)
}

// studentRepository implements StudentRepository
type studentRepository struct {
	db *database.DatabaseConnections
}

// NewStudentRepository creates a new student repository
func NewStudentRepository(db *database.DatabaseConnections) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(student *model.Student) error {
	return r.db.Write.Create(student).Error
}

func (r *studentRepository) GetByID(id uint) (*model.Student, error) {
	var student model.Student
	err := r.db.Read.Preload("User").Preload("Class").Preload("Parent").First(&student, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByStudentNumber(studentNumber string) (*model.Student, error) {
	var student model.Student
	err := r.db.Read.Preload("User").Preload("Class").Preload("Parent").
		Where("student_number = ?", studentNumber).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByUserID(userID uint) (*model.Student, error) {
	var student model.Student
	err := r.db.Read.Preload("User").Preload("Class").Preload("Parent").
		Where("user_id = ?", userID).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(student *model.Student) error {
	return r.db.Write.Save(student).Error
}

func (r *studentRepository) Delete(id uint) error {
	return r.db.Write.Delete(&model.Student{}, id).Error
}

func (r *studentRepository) List(offset, limit int, search string) ([]model.Student, int64, error) {
	var students []model.Student
	var total int64

	query := r.db.Read.Preload("User").Preload("Class").Preload("Parent")

	if search != "" {
		query = query.Joins("JOIN users ON users.id = students.user_id").
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

func (r *studentRepository) GetByClass(classID uint, offset, limit int) ([]model.Student, int64, error) {
	var students []model.Student
	var total int64

	query := r.db.Read.Preload("User").Preload("Class").Preload("Parent").
		Where("class_id = ?", classID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	return students, total, err
}

func (r *studentRepository) GetByParent(parentID uint, offset, limit int) ([]model.Student, int64, error) {
	var students []model.Student
	var total int64

	query := r.db.Read.Preload("User").Preload("Class").Preload("Parent").
		Where("parent_id = ?", parentID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	return students, total, err
}
