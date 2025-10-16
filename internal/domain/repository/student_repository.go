package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"github.com/rs/zerolog/log"
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
	BulkDelete(ids []uuid.UUID) error
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
		log.Error().
			Err(err).
			Str("student_number", student.StudentNumber).
			Str("tenant_id", student.TenantID.String()).
			Msg("Failed to set tenant context for student creation")
		return err
	}
	err := r.db.Write.Create(student).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("student_number", student.StudentNumber).
			Str("tenant_id", student.TenantID.String()).
			Str("tenant_user_id", student.TenantUserID.String()).
			Msg("Failed to create student in database")
	}
	return err
}

func (r *studentRepository) GetByID(id uuid.UUID) (*model.Student, error) {
	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").First(&student, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("student_id", id.String()).
				Msg("Student not found by ID")
			return nil, errors.New("student not found")
		}
		log.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Database error while getting student by ID")
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) GetByStudentNumber(studentNumber string, tenantID uuid.UUID) (*model.Student, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		log.Error().
			Err(err).
			Str("student_number", studentNumber).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to set tenant context for GetByStudentNumber")
		return nil, err
	}

	var student model.Student
	err := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("student_number = ? AND tenant_id = ?", studentNumber, tenantID).First(&student).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Debug().
				Str("student_number", studentNumber).
				Str("tenant_id", tenantID.String()).
				Msg("Student not found by student number")
			return nil, errors.New("student not found")
		}
		log.Error().
			Err(err).
			Str("student_number", studentNumber).
			Str("tenant_id", tenantID.String()).
			Msg("Database error in GetByStudentNumber")
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
			log.Debug().
				Str("tenant_user_id", tenantUserID.String()).
				Msg("Student not found by tenant user ID")
			return nil, errors.New("student not found")
		}
		log.Error().
			Err(err).
			Str("tenant_user_id", tenantUserID.String()).
			Msg("Database error in GetByTenantUserID")
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) Update(student *model.Student) error {
	if err := r.SetTenantContext(student.TenantID); err != nil {
		log.Error().
			Err(err).
			Str("student_id", student.ID.String()).
			Str("student_number", student.StudentNumber).
			Str("tenant_id", student.TenantID.String()).
			Msg("Failed to set tenant context for student update")
		return err
	}
	err := r.db.Write.Save(student).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("student_id", student.ID.String()).
			Str("student_number", student.StudentNumber).
			Str("tenant_id", student.TenantID.String()).
			Msg("Failed to update student in database")
	}
	return err
}

func (r *studentRepository) Delete(id uuid.UUID) error {
	err := r.db.Write.Delete(&model.Student{}, id).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to delete student from database")
	}
	return err
}

func (r *studentRepository) BulkDelete(ids []uuid.UUID) error {
	if len(ids) == 0 {
		return nil
	}

	err := r.db.Write.Where("id IN (?)", ids).Delete(&model.Student{}).Error
	if err != nil {
		log.Error().
			Err(err).
			Interface("ids", ids).
			Msg("Failed to bulk delete students from database")
	}
	return err
}

func (r *studentRepository) List(tenantID uuid.UUID, offset, limit int, search string) ([]model.Student, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Int("offset", offset).
			Int("limit", limit).
			Str("search", search).
			Msg("Failed to set tenant context for student list")
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
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("search", search).
			Msg("Failed to count students in List method")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Int("offset", offset).
			Int("limit", limit).
			Str("search", search).
			Msg("Failed to list students from database")
	}
	return students, total, err
}

func (r *studentRepository) GetByClass(tenantID, classID uuid.UUID, offset, limit int) ([]model.Student, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("class_id", classID.String()).
			Int("offset", offset).
			Int("limit", limit).
			Msg("Failed to set tenant context for GetByClass")
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("class_id = ? AND tenant_id = ?", classID, tenantID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("class_id", classID.String()).
			Msg("Failed to count students by class")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("class_id", classID.String()).
			Int("offset", offset).
			Int("limit", limit).
			Msg("Failed to get students by class from database")
	}
	return students, total, err
}

func (r *studentRepository) GetByParent(tenantID, parentID uuid.UUID, offset, limit int) ([]model.Student, int64, error) {
	if err := r.SetTenantContext(tenantID); err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("parent_id", parentID.String()).
			Int("offset", offset).
			Int("limit", limit).
			Msg("Failed to set tenant context for GetByParent")
		return nil, 0, err
	}

	var students []model.Student
	var total int64

	query := r.db.Read.Preload("TenantUser.User").Preload("Class").Preload("Parent").
		Where("parent_id = ? AND tenant_id = ?", parentID, tenantID)

	// Get total count
	if err := query.Model(&model.Student{}).Count(&total).Error; err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("parent_id", parentID.String()).
			Msg("Failed to count students by parent")
		return nil, 0, err
	}

	// Get paginated results
	err := query.Offset(offset).Limit(limit).Find(&students).Error
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("parent_id", parentID.String()).
			Int("offset", offset).
			Int("limit", limit).
			Msg("Failed to get students by parent from database")
	}
	return students, total, err
}
