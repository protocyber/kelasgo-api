package examples

// EXAMPLE: Refactored student_service.go using ContextLogger
// This is a reference example showing how to migrate existing services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/model"
	"github.com/protocyber/kelasgo-api/internal/domain/repository"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

// ExampleStudentService shows how to use ContextLogger in services
type ExampleStudentService struct {
	studentRepo    repository.StudentRepository
	tenantUserRepo repository.TenantUserRepository
}

// ExampleCreate demonstrates the refactored Create method with ContextLogger
func (s *ExampleStudentService) ExampleCreate(requestID string, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error) {
	// Create logger with request ID and tenant ID
	logger := util.NewContextLoggerWithRequestID(requestID).
		WithTenantID(tenantID.String())

	// Check if tenant user exists
	tenantUser, err := s.tenantUserRepo.GetByID(req.TenantUserID)
	if err != nil {
		// Error log will include request_id and tenant_id
		logger.Error().
			Err(err).
			Str("tenant_user_id", req.TenantUserID.String()).
			Msg("Tenant user not found during student creation")
		return nil, errors.New("tenant user not found")
	}

	// Verify tenant user belongs to the correct tenant
	if tenantUser.TenantID != tenantID {
		// Warning log will include request_id and tenant_id
		logger.Warn().
			Str("tenant_user_id", req.TenantUserID.String()).
			Str("expected_tenant", tenantID.String()).
			Str("actual_tenant", tenantUser.TenantID.String()).
			Msg("Tenant user does not belong to the specified tenant")
		return nil, errors.New("tenant user does not belong to this tenant")
	}

	// Check if student number already exists within tenant
	existingStudent, _ := s.studentRepo.GetByStudentNumber(req.StudentNumber, tenantID)
	if existingStudent != nil {
		// Warning log will include request_id and tenant_id
		logger.Warn().
			Str("student_number", req.StudentNumber).
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

	err = s.studentRepo.Create(student)
	if err != nil {
		// Error log will include request_id and tenant_id
		logger.Error().
			Err(err).
			Str("student_number", req.StudentNumber).
			Msg("Failed to create student in database")
		return nil, errors.New("failed to create student")
	}

	// Success log - use standard log (no request_id)
	log.Info().
		Str("student_id", student.ID.String()).
		Str("student_number", student.StudentNumber).
		Str("tenant_id", tenantID.String()).
		Msg("Student created successfully")

	return student, nil
}

// ExampleUpdate demonstrates the refactored Update method with ContextLogger
func (s *ExampleStudentService) ExampleUpdate(requestID string, tenantID uuid.UUID, id uuid.UUID, req dto.UpdateStudentRequest) (*model.Student, error) {
	logger := util.NewContextLoggerWithRequestID(requestID).
		WithTenantID(tenantID.String())

	// Get existing student
	student, err := s.studentRepo.GetByID(id, tenantID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Student not found during update")
		return nil, errors.New("student not found")
	}

	// Update fields
	if req.StudentNumber != "" && req.StudentNumber != student.StudentNumber {
		existingStudent, _ := s.studentRepo.GetByStudentNumber(req.StudentNumber, tenantID)
		if existingStudent != nil && existingStudent.ID != id {
			logger.Warn().
				Str("student_number", req.StudentNumber).
				Str("existing_student_id", existingStudent.ID.String()).
				Msg("Student number already in use by another student")
			return nil, errors.New("student number already exists")
		}
		student.StudentNumber = req.StudentNumber
	}

	if req.ClassID != nil {
		student.ClassID = req.ClassID
	}

	if req.ParentID != nil {
		student.ParentID = req.ParentID
	}

	// Save changes
	err = s.studentRepo.Update(student)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to update student in database")
		return nil, errors.New("failed to update student")
	}

	// Success log - use standard log
	log.Info().
		Str("student_id", student.ID.String()).
		Str("student_number", student.StudentNumber).
		Str("tenant_id", tenantID.String()).
		Msg("Student updated successfully")

	return student, nil
}

// ExampleDelete demonstrates the refactored Delete method with ContextLogger
func (s *ExampleStudentService) ExampleDelete(requestID string, tenantID uuid.UUID, id uuid.UUID) error {
	logger := util.NewContextLoggerWithRequestID(requestID).
		WithTenantID(tenantID.String())

	// Verify student exists and belongs to tenant
	student, err := s.studentRepo.GetByID(id, tenantID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Student not found during deletion")
		return errors.New("student not found")
	}

	// Delete student
	err = s.studentRepo.Delete(id, tenantID)
	if err != nil {
		logger.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to delete student from database")
		return errors.New("failed to delete student")
	}

	// Success log - use standard log
	log.Info().
		Str("student_id", id.String()).
		Str("student_number", student.StudentNumber).
		Str("tenant_id", tenantID.String()).
		Msg("Student deleted successfully")

	return nil
}

// ExampleUsingConvenienceMethods shows alternative approach with convenience methods
func (s *ExampleStudentService) ExampleUsingConvenienceMethods(requestID string, tenantID uuid.UUID, req dto.CreateStudentRequest) (*model.Student, error) {
	logger := util.NewContextLoggerWithRequestID(requestID).
		WithTenantID(tenantID.String())

	// Check if tenant user exists
	tenantUser, err := s.tenantUserRepo.GetByID(req.TenantUserID)
	if err != nil {
		// Using convenience method
		logger.LogError(err, "Tenant user not found during student creation", map[string]interface{}{
			"tenant_user_id": req.TenantUserID.String(),
		})
		return nil, errors.New("tenant user not found")
	}

	// Verify tenant user belongs to the correct tenant
	if tenantUser.TenantID != tenantID {
		// Using convenience method for warning
		logger.LogWarn("Tenant user does not belong to the specified tenant", map[string]interface{}{
			"tenant_user_id":  req.TenantUserID.String(),
			"expected_tenant": tenantID.String(),
			"actual_tenant":   tenantUser.TenantID.String(),
		})
		return nil, errors.New("tenant user does not belong to this tenant")
	}

	// ... rest of the implementation
	return nil, nil
}
