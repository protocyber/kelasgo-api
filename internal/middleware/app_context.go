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
		c.Set("app_context", appCtx)
		c.Set("app_config", cfg)
		c.Set("app_url", cfg.App.URL)
		c.Set("timezone", cfg.App.Timezone)
		c.Set("locale", cfg.App.Locale)

		c.Next()
	}
}

// GetAppContext is a helper to get app context from gin context
func GetAppContext(c *gin.Context) (*util.AppContext, bool) {
	appCtx, exists := c.Get("app_context")
	if !exists {
		return nil, false
	}
	return appCtx.(*util.AppContext), true
}

// GetAppConfig is a helper to get app config from gin context
func GetAppConfig(c *gin.Context) (*config.Config, bool) {
	cfg, exists := c.Get("app_config")
	if !exists {
		return nil, false
	}
	return cfg.(*config.Config), true
}
