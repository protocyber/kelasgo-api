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

// AuthHandler handles authentication related requests
type AuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(authService service.AuthService, validator *validator.Validate) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		validator:   validator,
	}
}

// Login handles user login
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("remote_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Failed to bind login request JSON")
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
			Str("remote_ip", c.ClientIP()).
			Msg("Login request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		log.Warn().
			Err(err).
			Str("username", req.Username).
			Str("tenant_id", req.TenantID).
			Str("remote_ip", c.ClientIP()).
			Msg("Login attempt failed")
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", response.User.ID.String()).
		Str("username", response.User.Username).
		Str("tenant_id", response.User.TenantID.String()).
		Str("remote_ip", c.ClientIP()).
		Msg("User logged in successfully")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// Register handles user registration
func (h *AuthHandler) Register(c *gin.Context) {
	// Get tenant ID from middleware context
	tenantID := middleware.GetTenantID(c)
	if tenantID == uuid.Nil {
		log.Error().
			Str("remote_ip", c.ClientIP()).
			Str("user_agent", c.Request.UserAgent()).
			Msg("Registration attempt without valid tenant ID")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Tenant ID required",
			Error:   "Registration requires a valid tenant context",
		})
		return
	}

	var req dto.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Failed to bind registration request JSON")
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
			Str("tenant_id", tenantID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Registration request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	user, err := h.authService.Register(tenantID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Str("tenant_id", tenantID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("User registration failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", user.ID.String()).
		Str("username", user.Username).
		Str("tenant_id", tenantID.String()).
		Str("remote_ip", c.ClientIP()).
		Msg("User registered successfully")

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// ChangePassword handles password change
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDInterface, exists := c.Get("user_id")
	if !exists || userIDInterface == nil {
		log.Error().
			Str("remote_ip", c.ClientIP()).
			Bool("exists", exists).
			Interface("user_id", userIDInterface).
			Msg("User ID not found in context during password change")
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
			Error:   "User ID not found in context",
		})
		return
	}

	userID, ok := userIDInterface.(uuid.UUID)
	if !ok {
		log.Error().
			Str("remote_ip", c.ClientIP()).
			Interface("user_id", userIDInterface).
			Msg("Invalid user ID format in context during password change")
		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Unauthorized",
			Error:   "Invalid user ID format in context",
		})
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Failed to bind change password request JSON")
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
			Str("user_id", userID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Change password request validation failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	err := h.authService.ChangePassword(userID, req)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID.String()).
			Str("remote_ip", c.ClientIP()).
			Msg("Password change failed")
		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Failed to change password",
			Error:   err.Error(),
		})
		return
	}

	log.Info().
		Str("user_id", userID.String()).
		Str("remote_ip", c.ClientIP()).
		Msg("Password changed successfully via handler")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Password changed successfully",
	})
}
