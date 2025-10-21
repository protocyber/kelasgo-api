package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/service"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// AuthHandler handles authentication related requests
type AuthHandler struct {
	BaseHandler
	authService service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, validator *validator.Validate, appCtx *util.AppContext) *AuthHandler {
	return &AuthHandler{
		BaseHandler: NewBaseHandler(appCtx),
		authService: authService,
		validator:   validator,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	logger := h.GetLogger(c)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to bind login request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Warn().
			Err(err).
			Str("email", req.Email).
			Msg("Login request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	serviceCtx := h.CreateServiceContext(c)
	response, err := h.authService.Login(serviceCtx, req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	logger := h.GetLogger(c)

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
			Msg("Failed to bind registration request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Warn().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Msg("Registration request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	serviceCtx := h.CreateServiceContext(c)
	user, err := h.authService.Register(serviceCtx, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	logger := h.GetLogger(c)

	userID, exists := h.ValidateUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to bind change password request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Warn().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Change password request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	serviceCtx := h.CreateServiceContext(c)
	err := h.authService.ChangePassword(serviceCtx, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to change password",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Password changed successfully",
	})
}

// SelectTenant handles tenant selection after authentication
func (h *AuthHandler) SelectTenant(c *gin.Context) {
	logger := h.GetLogger(c)

	userID, exists := h.ValidateUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	var req dto.TenantSelectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
			Str("user_id", userID.String()).
			Msg("Failed to bind tenant selection request JSON")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		logger.Warn().
			Err(err).
			Str("user_id", userID.String()).
			Str("tenant_id", req.TenantID).
			Msg("Tenant selection request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	serviceCtx := h.CreateServiceContext(c)
	response, err := h.authService.SelectTenant(serviceCtx, userID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant selection failed",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Tenant selected successfully",
		Data:    response,
	})
}

// GetUserTenants handles getting all tenants for the authenticated user
func (h *AuthHandler) GetUserTenants(c *gin.Context) {
	// logger := h.GetLogger(c)

	userID, exists := h.ValidateUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	serviceCtx := h.CreateServiceContext(c)
	tenants, err := h.authService.GetUserTenants(serviceCtx, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to get user tenants",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "User tenants retrieved successfully",
		Data:    tenants,
	})
}
