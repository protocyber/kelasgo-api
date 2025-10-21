package handler

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	appCtx *util.AppContext
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(appCtx *util.AppContext) BaseHandler {
	return BaseHandler{
		appCtx: appCtx,
	}
}

// GetLogger extracts context logger from gin context
func (b *BaseHandler) GetLogger(c *gin.Context) *util.ContextLogger {
	if logger, exists := c.Get("context_logger"); exists {
		if contextLogger, ok := logger.(*util.ContextLogger); ok {
			return contextLogger
		}
	}
	// Fallback to creating new logger
	return util.NewContextLogger(c)
}

// GetAppContext returns the app context
func (b *BaseHandler) GetAppContext() *util.AppContext {
	return b.appCtx
}

// CreateServiceContext creates a context suitable for service layer
func (b *BaseHandler) CreateServiceContext(c *gin.Context) context.Context {
	return util.CreateServiceContextFromGin(c)
}

// GetUserID extracts user ID from gin context as UUID
func (b *BaseHandler) GetUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDInterface, exists := c.Get("user_id")
	if !exists || userIDInterface == nil {
		return uuid.Nil, false
	}

	if userID, ok := userIDInterface.(uuid.UUID); ok {
		return userID, true
	}

	return uuid.Nil, false
}

// GetTenantID extracts tenant ID from gin context as string
func (b *BaseHandler) GetTenantID(c *gin.Context) (string, bool) {
	tenantIDInterface, exists := c.Get(string(util.XTenantIDKey))
	if !exists || tenantIDInterface == nil {
		return "", false
	}

	if tenantID, ok := tenantIDInterface.(string); ok && tenantID != "" {
		return tenantID, true
	}

	return "", false
}

// GetTenantIDAsUUID extracts tenant ID from gin context as UUID
func (b *BaseHandler) GetTenantIDAsUUID(c *gin.Context) (uuid.UUID, bool) {
	tenantID, exists := b.GetTenantID(c)
	if !exists {
		return uuid.Nil, false
	}

	if parsed, err := uuid.Parse(tenantID); err == nil {
		return parsed, true
	}

	return uuid.Nil, false
}

// ValidateUserID checks if user ID exists in context and logs error if not
func (b *BaseHandler) ValidateUserID(c *gin.Context) (uuid.UUID, bool) {
	logger := b.GetLogger(c)
	userID, exists := b.GetUserID(c)

	if !exists {
		logger.Error().
			Msg("User ID not found in context")
		return uuid.Nil, false
	}

	return userID, true
}

// Deprecated: Use GetLogger and CreateServiceContext instead
func (b *BaseHandler) ExtractContext(c *gin.Context) {
	// This method is kept for backward compatibility
	// New code should use GetLogger(c) and CreateServiceContext(c) directly
}
