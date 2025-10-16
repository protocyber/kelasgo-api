package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/service"
	"github.com/protocyber/kelasgo-api/internal/server/middleware"
)

// StudentHandler handles student related requests
type StudentHandler struct {
	BaseHandler
	studentService service.StudentService
	validator      *validator.Validate
}

// NewStudentHandler creates a new student handler
func NewStudentHandler(studentService service.StudentService, validator *validator.Validate) *StudentHandler {
	return &StudentHandler{
		studentService: studentService,
		validator:      validator,
	}
}

// Create handles student creation
func (h *StudentHandler) Create(c *gin.Context) {
	h.InitLogger(c)

	var req dto.CreateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().
			Err(err).
			Msg("Failed to bind create student request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.log.Warn().
			Err(err).
			Str("student_number", req.StudentNumber).
			Str("tenant_user_id", req.TenantUserID.String()).
			Msg("Create student request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == uuid.Nil {
		h.log.Error().
			Str("student_number", req.StudentNumber).
			Msg("Student creation attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "Student creation requires a valid tenant context",
		})
		return
	}

	student, err := h.studentService.Create(tenantID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to create student",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "Student created successfully",
		Data:    student,
	})
}

// GetByID handles getting student by ID
func (h *StudentHandler) GetByID(c *gin.Context) {
	h.InitLogger(c)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("id_param", idStr).
			Msg("Invalid student ID format in get request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid student ID format",
			Error:   err.Error(),
		})
		return
	}

	student, err := h.studentService.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "Student not found",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Student retrieved successfully",
		Data:    student,
	})
}

// Update handles student update
func (h *StudentHandler) Update(c *gin.Context) {
	h.InitLogger(c)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("id_param", idStr).
			Msg("Invalid student ID format in update request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid student ID format",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().
			Err(err).
			Str("student_id", id.String()).
			Msg("Failed to bind update student request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.log.Warn().
			Err(err).
			Str("student_id", id.String()).
			Msg("Update student request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	student, err := h.studentService.Update(id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to update student",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Student updated successfully",
		Data:    student,
	})
}

// Delete handles student deletion
func (h *StudentHandler) Delete(c *gin.Context) {
	h.InitLogger(c)

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("id_param", idStr).
			Msg("Invalid student ID format in delete request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid student ID format",
			Error:   err.Error(),
		})
		return
	}

	err = h.studentService.Delete(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to delete student",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Student deleted successfully",
	})
}

// BulkDelete handles bulk student deletion
func (h *StudentHandler) BulkDelete(c *gin.Context) {
	h.InitLogger(c)

	var req dto.BulkDeleteStudentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error().
			Err(err).
			Msg("Failed to bind bulk delete student request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.log.Warn().
			Err(err).
			Interface("student_ids", req.IDs).
			Msg("Bulk delete student request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == uuid.Nil {
		h.log.Error().
			Interface("student_ids", req.IDs).
			Msg("Bulk delete students attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "Student bulk deletion requires a valid tenant context",
		})
		return
	}

	err := h.studentService.BulkDelete(tenantID, req.IDs)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to bulk delete students",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Students bulk deleted successfully",
	})
}

// List handles student listing with pagination
func (h *StudentHandler) List(c *gin.Context) {
	h.InitLogger(c)

	var params dto.StudentQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.log.Error().
			Err(err).
			Msg("Failed to bind student list query parameters")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(params); err != nil {
		h.log.Warn().
			Err(err).
			Interface("params", params).
			Msg("Student list query parameters validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == uuid.Nil {
		h.log.Error().
			Msg("Student listing attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "Student listing requires a valid tenant context",
		})
		return
	}

	students, meta, err := h.studentService.List(tenantID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Failed to retrieve students",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Message: "Students retrieved successfully",
		Data:    students,
		Meta:    *meta,
	})
}

// GetByClass handles getting students by class ID
func (h *StudentHandler) GetByClass(c *gin.Context) {
	h.InitLogger(c)

	classIDStr := c.Param("class_id")
	classID, err := uuid.Parse(classIDStr)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("class_id_param", classIDStr).
			Msg("Invalid class ID format in get students by class request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid class ID format",
			Error:   err.Error(),
		})
		return
	}

	var params dto.QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.log.Error().
			Err(err).
			Msg("Failed to bind query parameters for students by class")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == uuid.Nil {
		h.log.Error().
			Str("class_id", classID.String()).
			Msg("Get students by class attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "Getting students by class requires a valid tenant context",
		})
		return
	}

	students, meta, err := h.studentService.GetByClass(tenantID, classID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Failed to retrieve students by class",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Message: "Students retrieved successfully",
		Data:    students,
		Meta:    *meta,
	})
}

// GetByParent handles getting students by parent ID
func (h *StudentHandler) GetByParent(c *gin.Context) {
	h.InitLogger(c)

	parentIDStr := c.Param("parent_id")
	parentID, err := uuid.Parse(parentIDStr)
	if err != nil {
		h.log.Error().
			Err(err).
			Str("parent_id_param", parentIDStr).
			Msg("Invalid parent ID format in get students by parent request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid parent ID format",
			Error:   err.Error(),
		})
		return
	}

	var params dto.QueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		h.log.Error().
			Err(err).
			Msg("Failed to bind query parameters for students by parent")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == uuid.Nil {
		h.log.Error().
			Str("parent_id", parentID.String()).
			Msg("Get students by parent attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "Getting students by parent requires a valid tenant context",
		})
		return
	}

	students, meta, err := h.studentService.GetByParent(tenantID, parentID, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Failed to retrieve students by parent",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Message: "Students retrieved successfully",
		Data:    students,
		Meta:    *meta,
	})
}
