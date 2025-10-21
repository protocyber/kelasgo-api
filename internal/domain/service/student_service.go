package service

import (
	"context"
	"errors"
	"math"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/domain/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// StudentService interface defines student service methods
type StudentService interface {
	Create(c context.Context, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error)
	GetByID(c context.Context, id uuid.UUID) (*model.Student, error)
	Update(c context.Context, id uuid.UUID, req dto.UpdateStudentRequest) (*model.Student, error)
	Delete(c context.Context, id uuid.UUID) error
	BulkDelete(c context.Context, tenantID uuid.UUID, ids []uuid.UUID) error
	List(c context.Context, tenantID uuid.UUID, params dto.StudentQueryParams) ([]model.Student, *dto.PaginationMeta, error)
	GetByClass(c context.Context, tenantID, classID uuid.UUID, params dto.QueryParams) ([]model.Student, *dto.PaginationMeta, error)
	GetByParent(c context.Context, tenantID, parentID uuid.UUID, params dto.QueryParams) ([]model.Student, *dto.PaginationMeta, error)
}

// studentService implements StudentService
type studentService struct {
	studentRepo    repository.StudentRepository
	tenantUserRepo repository.TenantUserRepository
}

// NewStudentService creates a new student service
func NewStudentService(
	studentRepo repository.StudentRepository,
	tenantUserRepo repository.TenantUserRepository,
) StudentService {
	return &studentService{
		studentRepo:    studentRepo,
		tenantUserRepo: tenantUserRepo,
	}
}

func (s *studentService) Create(c context.Context, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Check if tenant user exists
	tenantUser, err := s.tenantUserRepo.GetByID(c, req.TenantUserID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_user_id", req.TenantUserID.String()).
			Str("tenant_id", tenantID.String()).
			Msg("Tenant user not found during student creation")
		return nil, errors.New("tenant user not found")
	}

	// Verify tenant user belongs to the correct tenant
	if tenantUser.TenantID != tenantID {
		logger.Warn().
			Str("tenant_user_id", req.TenantUserID.String()).
			Str("expected_tenant", tenantID.String()).
			Str("actual_tenant", tenantUser.TenantID.String()).
			Msg("Tenant user does not belong to the specified tenant")
		return nil, errors.New("tenant user does not belong to this tenant")
	}

	// Check if student number already exists within tenant
	existingStudent, _ := s.studentRepo.GetByStudentNumber(c, req.StudentNumber, tenantID)
	if existingStudent != nil {
		logger.Warn().
			Str("student_number", req.StudentNumber).
			Str("tenant_id", tenantID.String()).
			Msg("Student creation attempt with existing student number")
		return nil, errors.New("student number already exists")
	}

	// Create student
	student := &model.Student{
		TenantID:      tenantID,
		TenantUserID:  req.TenantUserID,
		StudentNumber: req.StudentNumber,
		AdmissionDate: req.AdmissionDate,
		ClassID:       req.ClassID,
		ParentID:      req.ParentID,
	}

	err = s.studentRepo.Create(c, student)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_number", req.StudentNumber).
			Str("tenant_id", tenantID.String()).
			Msg("Failed to create student in database")
		return nil, errors.New("failed to create student")
	}

	return student, nil
}

func (s *studentService) GetByID(c context.Context, id uuid.UUID) (*model.Student, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	student, err := s.studentRepo.GetByID(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to get student by ID")
		return nil, errors.New("student not found")
	}
	return student, nil
}

func (s *studentService) Update(c context.Context, id uuid.UUID, req dto.UpdateStudentRequest) (*model.Student, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Get existing student
	student, err := s.studentRepo.GetByID(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Student not found during update")
		return nil, err
	}

	// Check if student number already exists (if changed and provided)
	if req.StudentNumber != nil && *req.StudentNumber != "" && *req.StudentNumber != student.StudentNumber {
		existingStudent, _ := s.studentRepo.GetByStudentNumber(c, *req.StudentNumber, student.TenantID)
		if existingStudent != nil && existingStudent.ID != id {
			logger.Warn().
				Str("student_number", *req.StudentNumber).
				Str("student_id", id.String()).
				Str("tenant_id", student.TenantID.String()).
				Msg("Student update attempt with existing student number")
			return nil, errors.New("student number already exists")
		}
	}

	// Update fields
	if req.StudentNumber != nil && *req.StudentNumber != "" {
		student.StudentNumber = *req.StudentNumber
	}
	if req.AdmissionDate != nil {
		student.AdmissionDate = *req.AdmissionDate
	}
	if req.ClassID != nil {
		student.ClassID = req.ClassID
	}
	if req.ParentID != nil {
		student.ParentID = req.ParentID
	}

	err = s.studentRepo.Update(c, student)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to update student in database")
		return nil, errors.New("failed to update student")
	}

	return student, nil
}

func (s *studentService) Delete(c context.Context, id uuid.UUID) error {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Check if student exists
	_, err := s.studentRepo.GetByID(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Student not found during delete")
		return err
	}

	err = s.studentRepo.Delete(c, id)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to delete student from database")
		return err
	}

	return nil
}

func (s *studentService) BulkDelete(c context.Context, tenantID uuid.UUID, ids []uuid.UUID) error {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	if len(ids) == 0 {
		return errors.New("no student IDs provided for bulk delete")
	}

	// Get students that belong to the tenant to validate they exist and log properly
	students, _, err := s.studentRepo.List(c, tenantID, 0, len(ids)*2, "")
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Interface("student_ids", ids).
			Msg("Failed to validate students for bulk delete")
		return errors.New("failed to validate students for bulk delete")
	}

	// Create a set of valid student IDs that belong to the tenant
	validStudentMap := make(map[uuid.UUID]bool)
	for _, student := range students {
		validStudentMap[student.ID] = true
	}

	// Filter IDs to only include students that belong to the tenant
	var validIDs []uuid.UUID
	var invalidIDs []uuid.UUID
	for _, id := range ids {
		if validStudentMap[id] {
			validIDs = append(validIDs, id)
		} else {
			invalidIDs = append(invalidIDs, id)
		}
	}

	if len(invalidIDs) > 0 {
		logger.Warn().
			Str("tenant_id", tenantID.String()).
			Interface("invalid_ids", invalidIDs).
			Msg("Some student IDs do not belong to the tenant or do not exist")
	}

	if len(validIDs) == 0 {
		return errors.New("no valid student IDs found for bulk delete in this tenant")
	}

	// Perform bulk delete
	err = s.studentRepo.BulkDelete(c, validIDs)
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Interface("student_ids", validIDs).
			Msg("Failed to bulk delete students from database")
		return errors.New("failed to bulk delete students")
	}

	return nil
}

func (s *studentService) List(c context.Context, tenantID uuid.UUID, params dto.StudentQueryParams) ([]model.Student, *dto.PaginationMeta, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	var students []model.Student
	var total int64
	var err error

	if params.ClassID != nil {
		students, total, err = s.studentRepo.GetByClass(c, tenantID, *params.ClassID, offset, params.Limit)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Str("class_id", params.ClassID.String()).
				Interface("params", params).
				Msg("Failed to get students by class")
		}
	} else if params.ParentID != nil {
		students, total, err = s.studentRepo.GetByParent(c, tenantID, *params.ParentID, offset, params.Limit)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Str("parent_id", params.ParentID.String()).
				Interface("params", params).
				Msg("Failed to get students by parent")
		}
	} else {
		students, total, err = s.studentRepo.List(c, tenantID, offset, params.Limit, params.Search)
		if err != nil {
			logger.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Interface("params", params).
				Msg("Failed to get students by tenant")
		}
	}

	if err != nil {
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	meta := &dto.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalRows:  total,
		TotalPages: totalPages,
	}

	return students, meta, nil
}

func (s *studentService) GetByClass(c context.Context, tenantID, classID uuid.UUID, params dto.QueryParams) ([]model.Student, *dto.PaginationMeta, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	students, total, err := s.studentRepo.GetByClass(c, tenantID, classID, offset, params.Limit)
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("class_id", classID.String()).
			Interface("params", params).
			Msg("Failed to get students by class")
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	meta := &dto.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalRows:  total,
		TotalPages: totalPages,
	}

	return students, meta, nil
}

func (s *studentService) GetByParent(c context.Context, tenantID, parentID uuid.UUID, params dto.QueryParams) ([]model.Student, *dto.PaginationMeta, error) {
	// Create context logger for service
	logger := util.NewServiceLogger(c)

	// Set defaults
	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = 10
	}

	offset := (params.Page - 1) * params.Limit

	students, total, err := s.studentRepo.GetByParent(c, tenantID, parentID, offset, params.Limit)
	if err != nil {
		logger.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("parent_id", parentID.String()).
			Interface("params", params).
			Msg("Failed to get students by parent")
		return nil, nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(params.Limit)))

	meta := &dto.PaginationMeta{
		Page:       params.Page,
		Limit:      params.Limit,
		TotalRows:  total,
		TotalPages: totalPages,
	}

	return students, meta, nil
}
