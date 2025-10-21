package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"gorm.io/gorm"
)

// StudentRepository interface defines student repository methods
type StudentRepository interface {
	Create(c context.Context, student *model.Student) error
	GetByID(c context.Context, id uuid.UUID) (*model.Student, error)
	GetByStudentNumber(c context.Context, studentNumber string, tenantID uuid.UUID) (*model.Student, error)
	GetByTenantUserID(c context.Context, tenantUserID uuid.UUID) (*model.Student, error)
	Update(c context.Context, student *model.Student) error
	Delete(c context.Context, id uuid.UUID) error
	BulkDelete(c context.Context, ids []uuid.UUID) error
	List(c context.Context, tenantID uuid.UUID, offset, limit int, search string) ([]model.Student, int64, error)
	GetByClass(c context.Context, tenantID, classID uuid.UUID, offset, limit int) ([]model.Student, int64, error)
	GetByParent(c context.Context, tenantID, parentID uuid.UUID, offset, limit int) ([]model.Student, int64, error)
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

func (r *studentRepository) Create(c context.Context, student *model.Student) error {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(student.TenantID); err != nil {
		return err
	}
	err := r.db.Write.Create(student).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "create_student").
			Msg("Database write operation failed")
	}
	return err
}

func (r *studentRepository) GetByID(c context.Context, id uuid.UUID) (*model.Student, error) {
	repoCtx := r.WithContext(c)
	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").First(&student, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		repoCtx.logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Database error while getting student by ID")
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByStudentNumber(c context.Context, studentNumber string, tenantID uuid.UUID) (*model.Student, error) {
	repoCtx := r.WithContext(c)
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
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "get_student_by_number").
			Msg("Database query failed")
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByTenantUserID(c context.Context, tenantUserID uuid.UUID) (*model.Student, error) {
	repoCtx := r.WithContext(c)
	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("tenant_user_id = ?", tenantUserID).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("student not found")
		}
		repoCtx.logger.Error().
			Err(err).
			Str("tenant_user_id", tenantUserID.String()).
			Msg("Database error in GetByTenantUserID")
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(c context.Context, student *model.Student) error {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(student.TenantID); err != nil {
		return err
	}
	err := r.db.Write.Save(student).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "update_student").
			Msg("Database write operation failed")
	}
	return err
}

func (r *studentRepository) Delete(c context.Context, id uuid.UUID) error {
	repoCtx := r.WithContext(c)
	err := r.db.Write.Delete(&model.Student{}, id).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "delete_student").
			Msg("Database write operation failed")
	}
	return err
}

func (r *studentRepository) BulkDelete(c context.Context, ids []uuid.UUID) error {
	repoCtx := r.WithContext(c)
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Write.Where("id IN (?)", ids).Delete(&model.Student{}).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "bulk_delete_students").
			Int("count", len(ids)).
			Msg("Database write operation failed")
	}
	return err
}

func (r *studentRepository) List(c context.Context, tenantID uuid.UUID, offset, limit int, search string) ([]model.Student, int64, error) {
	repoCtx := r.WithContext(c)
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
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "count_students").
			Msg("Database query failed")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "list_students").
			Msg("Database query failed")
	}
	return students, total, err
}

func (r *studentRepository) GetByClass(c context.Context, tenantID, classID uuid.UUID, offset, limit int) ([]model.Student, int64, error) {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("class_id = ? AND tenant_id = ?", classID, tenantID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "count_students_by_class").
			Msg("Database query failed")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "get_students_by_class").
			Msg("Database query failed")
	}
	return students, total, err
}

func (r *studentRepository) GetByParent(c context.Context, tenantID, parentID uuid.UUID, offset, limit int) ([]model.Student, int64, error) {
	repoCtx := r.WithContext(c)
	if err := r.SetTenantContext(tenantID); err != nil {
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("parent_id = ? AND tenant_id = ?", parentID, tenantID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "count_students_by_parent").
			Msg("Database query failed")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		repoCtx.logger.Error().
			Err(err).
			Str("operation", "get_students_by_parent").
			Msg("Database query failed")
	}
	return students, total, err
}
