package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// AppContextMiddleware injects application context into request context
func AppContextMiddleware(cfg *config.Config) gin.HandlerFunc {
	// Create app context once during middleware initialization
	appCtx, err := util.NewAppContext(context.Background(), cfg)
	if err != nil {
		// Log error but continue with basic context
		appCtx = &util.AppContext{
			Config: cfg,
		}
	}

	return func(c *gin.Context) {
		// Add app context to the request context
		ctx := util.WithAppContext(c.Request.Context(), appCtx)
		c.Request = c.Request.WithContext(ctx)

		// Also make it available directly in gin context for convenience
		c.Set(string(util.AppContextKey), appCtx)

		// Create and set context logger
		logger := util.NewContextLogger(c)
		c.Set("context_logger", logger)

		c.Next()
	}
}

// GetAppContext is a helper to get app context from gin context
func GetAppContext(c *gin.Context) (*util.AppContext, bool) {
	appCtx, exists := c.Get(string(util.AppContextKey))
	if !exists {
		return nil, false
	}
	return appCtx.(*util.AppContext), true
}

// GetContextLogger helper function to extract logger from gin context
func GetContextLogger(c *gin.Context) *util.ContextLogger {
	if logger, exists := c.Get("context_logger"); exists {
		if contextLogger, ok := logger.(*util.ContextLogger); ok {
			return contextLogger
		}
	}
	// Fallback to creating new logger
	return util.NewContextLogger(c)
}
