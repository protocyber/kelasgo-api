package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/middleware"
	"github.com/protocyber/kelasgo-api/internal/service"
	"github.com/rs/zerolog/log"
)

// UserHandler handles user related requests
type UserHandler struct {
	userService service.UserService
	validator   *validator.Validate
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, validator *validator.Validate) *UserHandler {
	return &UserHandler{
		userService: userService,
		validator:   validator,
	}
}

// Create handles user creation
func (h *UserHandler) Create(c *gin.Context) {
	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("remote_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Failed to bind create user request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Warn().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Str("remote_ip", c.ClientIP()).
			Msg("Create user request validation failed")
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
		log.Error().
			Str("username", req.Username).
			Str("remote_ip", c.ClientIP()).
			Msg("User creation attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "User creation requires a valid tenant context",
		})
		return
	}

	user, err := h.userService.Create(tenantID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Str("tenant_id", tenantID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("User creation failed in handler")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to create user",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("username", user.Username).
		Str("tenant_id", tenantID.String()).
		Str("remote_ip", c.ClientIP()).
		Msg("User created successfully via handler")

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User created successfully",
		Data:    user,
	})
}

// GetByID handles getting user by ID
func (h *UserHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("id_param", idStr).
			Str("remote_ip", c.ClientIP()).
			Msg("Invalid user ID format in get request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.GetByID(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Failed to get user by ID in handler")
		c.JSON(http.StatusNotFound, dto.Response{
			Success: false,
			Message: "User not found",
			Error:   err.Error(),
		})
		return
	}

	log.Debug().
		Str("user_id", id.String()).
		Str("username", user.Username).
		Str("remote_ip", c.ClientIP()).
		Msg("User retrieved successfully")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// Update handles user update
func (h *UserHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("id_param", idStr).
			Str("remote_ip", c.ClientIP()).
			Msg("Invalid user ID format in update request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
		return
	}

	var req dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Failed to bind update user request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Warn().
			Err(err).
			Str("user_id", id.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Update user request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.userService.Update(id, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("User update failed in handler")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to update user",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", id.String()).
		Str("username", user.Username).
		Str("remote_ip", c.ClientIP()).
		Msg("User updated successfully via handler")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User updated successfully",
		Data:    user,
	})
}

// Delete handles user deletion
func (h *UserHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		log.Error().
			Err(err).
			Str("id_param", idStr).
			Str("remote_ip", c.ClientIP()).
			Msg("Invalid user ID format in delete request")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid user ID format",
			Error:   err.Error(),
		})
		return
	}

	err = h.userService.Delete(id)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", id.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("User deletion failed in handler")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to delete user",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", id.String()).
		Str("remote_ip", c.ClientIP()).
		Msg("User deleted successfully via handler")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User deleted successfully",
	})
}

// BulkDelete handles bulk user deletion
func (h *UserHandler) BulkDelete(c *gin.Context) {
	var req dto.BulkDeleteUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("remote_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Failed to bind bulk delete user request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		log.Warn().
			Err(err).
			Interface("user_ids", req.IDs).
			Str("remote_ip", c.ClientIP()).
			Msg("Bulk delete user request validation failed")
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
		log.Error().
			Interface("user_ids", req.IDs).
			Str("remote_ip", c.ClientIP()).
			Msg("Bulk delete users attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "User bulk deletion requires a valid tenant context",
		})
		return
	}

	err := h.userService.BulkDelete(tenantID, req.IDs)
	if err != nil {
		log.Error().
			Err(err).
			Interface("user_ids", req.IDs).
			Str("tenant_id", tenantID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Bulk user deletion failed in handler")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to bulk delete users",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Interface("user_ids", req.IDs).
		Str("tenant_id", tenantID.String()).
		Str("remote_ip", c.ClientIP()).
		Msg("Users bulk deleted successfully via handler")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Users bulk deleted successfully",
	})
}

// List handles user listing with pagination
func (h *UserHandler) List(c *gin.Context) {
	var params dto.UserQueryParams
	if err := c.ShouldBindQuery(&params); err != nil {
		log.Error().
			Err(err).
			Str("remote_ip", c.ClientIP()).
			Msg("Failed to bind user list query parameters")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid query parameters",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(params); err != nil {
		log.Warn().
			Err(err).
			Int("page", params.Page).
			Int("limit", params.Limit).
			Str("search", params.Search).
			Str("remote_ip", c.ClientIP()).
			Msg("User list query parameters validation failed")
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
		log.Error().
			Str("remote_ip", c.ClientIP()).
			Msg("User listing attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "User listing requires a valid tenant context",
		})
		return
	}

	users, meta, err := h.userService.List(tenantID, params)
	if err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Int("page", params.Page).
			Int("limit", params.Limit).
			Str("search", params.Search).
			Str("remote_ip", c.ClientIP()).
			Msg("User listing failed in handler")
		c.JSON(http.StatusInternalServerError, dto.Response{
			Success: false,
			Message: "Failed to retrieve users",
			Error:   err.Error(),
		})
		return
	}

	log.Debug().
		Str("tenant_id", tenantID.String()).
		Int("page", params.Page).
		Int("limit", params.Limit).
		Int64("total_users", meta.TotalRows).
		Str("remote_ip", c.ClientIP()).
		Msg("Users listed successfully")

	c.JSON(http.StatusOK, dto.PaginatedResponse{
		Success: true,
		Message: "Users retrieved successfully",
		Data:    users,
		Meta:    *meta,
	})
}
