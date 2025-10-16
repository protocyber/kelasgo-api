package examples

// EXAMPLE: Refactored auth_handler.go using ContextLogger
// This is a reference example showing how to migrate existing handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/protocyber/kelasgo-api/internal/domain/dto"
	"github.com/protocyber/kelasgo-api/internal/domain/service"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

// ExampleAuthHandler shows how to use ContextLogger in handlers
type ExampleAuthHandler struct {
	authService service.AuthService
	validator   *validator.Validate
}

// ExampleLogin demonstrates the refactored Login handler with ContextLogger
func (h *ExampleAuthHandler) ExampleLogin(c *gin.Context) {
	// Create context logger - automatically captures request_id
	logger := util.NewContextLogger(c)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Error log will include request_id automatically
		logger.Error().
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
		// Warning log will include request_id automatically
		logger.Warn().
			Err(err).
			Str("email", req.Email).
			Str("remote_ip", c.ClientIP()).
			Msg("Login request validation failed")

		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Pass request ID to service layer
	_ = logger.GetRequestID()                 // Get request ID for passing to service
	response, err := h.authService.Login(req) // TODO: Update service to accept requestID
	// requestID := logger.GetRequestID()
	// response, err := h.authService.Login(requestID, req) // After migration
	if err != nil {
		// Warning log will include request_id automatically
		logger.Warn().
			Err(err).
			Str("email", req.Email).
			Str("remote_ip", c.ClientIP()).
			Msg("Login attempt failed")

		c.JSON(http.StatusUnauthorized, dto.Response{
			Success: false,
			Message: "Login failed",
			Error:   err.Error(),
		})
		return
	}

	tenantIDStr := "none"
	if response.User.TenantID != nil {
		tenantIDStr = response.User.TenantID.String()
	}

	// Success log - use standard log (no request_id needed)
	log.Info().
		Str("user_id", response.User.ID.String()).
		Str("email", response.User.Email).
		Str("tenant_id", tenantIDStr).
		Str("remote_ip", c.ClientIP()).
		Msg("User logged in successfully")

	c.JSON(http.StatusOK, dto.Response{
		Success: true,
		Message: "Login successful",
		Data:    response,
	})
}

// ExampleRegister demonstrates the refactored Register handler with ContextLogger
func (h *ExampleAuthHandler) ExampleRegister(c *gin.Context) {
	logger := util.NewContextLogger(c)

	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error().
			Err(err).
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
		logger.Warn().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Str("remote_ip", c.ClientIP()).
			Msg("Registration request validation failed")

		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Validation failed",
			Error:   err.Error(),
		})
		return
	}

	// Pass request ID to service layer
	_ = logger.GetRequestID()                // Get request ID for passing to service
	user, err := h.authService.Register(req) // TODO: Update service to accept requestID
	// requestID := logger.GetRequestID()
	// user, err := h.authService.Register(requestID, req) // After migration
	if err != nil {
		logger.Error().
			Err(err).
			Str("username", req.Username).
			Str("email", req.Email).
			Str("remote_ip", c.ClientIP()).
			Msg("User registration failed")

		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Registration failed",
			Error:   err.Error(),
		})
		return
	}

	// Success log - use standard log
	log.Info().
		Str("user_id", user.ID.String()).
		Str("username", user.Username).
		Str("email", req.Email).
		Str("remote_ip", c.ClientIP()).
		Msg("User registered successfully")

	c.JSON(http.StatusCreated, dto.Response{
		Success: true,
		Message: "User registered successfully",
		Data:    user,
	})
}

// ExampleUsingConvenienceMethods shows alternative approach with convenience methods
func (h *ExampleAuthHandler) ExampleUsingConvenienceMethods(c *gin.Context) {
	logger := util.NewContextLogger(c)

	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		// Using convenience method
		logger.LogError(err, "Failed to bind login request JSON", map[string]interface{}{
			"remote_ip":  c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
		})

		c.JSON(http.StatusBadRequest, dto.Response{
			Success: false,
			Message: "Invalid request body",
			Error:   err.Error(),
		})
		return
	}

	// ... rest of the handler
}
