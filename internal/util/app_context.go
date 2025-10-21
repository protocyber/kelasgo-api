package util

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/config"
	request_id "github.com/protocyber/kelasgo-api/pkg/gin-request-id"
)

// AppContextKey defines keys for context values
type AppContextKeyType string

const (
	AppContextKey AppContextKeyType = "app_context"
)

// AppContext wraps application context with configuration
type AppContext struct {
	context.Context
	Config   *config.Config
	Location *time.Location
}

// NewAppContext creates a new application context with configuration
func NewAppContext(ctx context.Context, cfg *config.Config) (*AppContext, error) {
	// Parse timezone from config
	location, err := time.LoadLocation(cfg.App.Timezone)
	if err != nil {
		// Fallback to UTC if timezone parsing fails
		location = time.UTC
	}

	return &AppContext{
		Context:  ctx,
		Config:   cfg,
		Location: location,
	}, nil
}

// GetAppURL returns the configured application URL
func (ac *AppContext) GetAppURL() string {
	return ac.Config.App.URL
}

// GetTimezone returns the configured timezone
func (ac *AppContext) GetTimezone() string {
	return ac.Config.App.Timezone
}

// GetLocale returns the configured locale
func (ac *AppContext) GetLocale() string {
	return ac.Config.App.Locale
}

// GetLocation returns the timezone location
func (ac *AppContext) GetLocation() *time.Location {
	return ac.Location
}

// Now returns the current time in the configured timezone
func (ac *AppContext) Now() time.Time {
	return time.Now().In(ac.Location)
}

// ParseTimeInTimezone parses a time string in the configured timezone
func (ac *AppContext) ParseTimeInTimezone(layout, value string) (time.Time, error) {
	return time.ParseInLocation(layout, value, ac.Location)
}

// FormatTimeInTimezone formats a time in the configured timezone
func (ac *AppContext) FormatTimeInTimezone(t time.Time, layout string) string {
	return t.In(ac.Location).Format(layout)
}

// GetAppName returns the configured application name
func (ac *AppContext) GetAppName() string {
	return ac.Config.App.Name
}

// GetAppVersion returns the configured application version
func (ac *AppContext) GetAppVersion() string {
	return ac.Config.App.Version
}

// GetAppDescription returns the configured application description
func (ac *AppContext) GetAppDescription() string {
	return ac.Config.App.Description
}

// GetPaginationDefaults returns pagination configuration
func (ac *AppContext) GetPaginationDefaults() (defaultLimit, maxLimit int, enabled bool) {
	return ac.Config.App.Pagination.DefaultLimit, ac.Config.App.Pagination.MaxLimit, ac.Config.App.Pagination.Enabled
}

// WithAppContext adds app context to a regular context
func WithAppContext(ctx context.Context, appCtx *AppContext) context.Context {
	return context.WithValue(ctx, AppContextKey, appCtx)
}

// GetAppContextFromContext extracts app context from context
func GetAppContextFromContext(ctx context.Context) (*AppContext, bool) {
	appCtx, ok := ctx.Value(AppContextKey).(*AppContext)
	return appCtx, ok
}

// CreateServiceContext creates a context suitable for service layer with all necessary values
func CreateServiceContext(ctx context.Context, appCtx *AppContext) context.Context {
	c := WithAppContext(ctx, appCtx)
	c = context.WithValue(c, request_id.XRequestIDKey, ctx.Value(request_id.XRequestIDKey))
	c = context.WithValue(c, XTenantIDKey, ctx.Value(XTenantIDKey))
	c = context.WithValue(c, "user_id", ctx.Value("user_id"))
	return c
}

// CreateServiceContextFromGin creates a service context from gin context
func CreateServiceContextFromGin(ginCtx interface{}) context.Context {
	ctx := context.Background()

	// Extract app context if available
	if appCtx, exists := extractValue(ginCtx, "app_context"); exists {
		if ac, ok := appCtx.(*AppContext); ok {
			ctx = WithAppContext(ctx, ac)
		}
	}

	// Copy request ID
	if requestID, exists := extractValue(ginCtx, request_id.XRequestIDKey); exists {
		ctx = context.WithValue(ctx, request_id.XRequestIDKey, requestID)
	}

	// Copy tenant ID
	if tenantID, exists := extractValue(ginCtx, string(XTenantIDKey)); exists {
		ctx = context.WithValue(ctx, XTenantIDKey, tenantID)
	}

	// Copy user ID
	if userID, exists := extractValue(ginCtx, "user_id"); exists {
		ctx = context.WithValue(ctx, "user_id", userID)
	}

	return ctx
}

// GetTenantIDFromContext extracts tenant ID from context
func GetTenantIDFromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(XTenantIDKey).(string)
	return tenantID, ok && tenantID != ""
}

// GetUserIDFromContext extracts user ID from context
func GetUserIDFromContext(ctx context.Context) (interface{}, bool) {
	userID := ctx.Value("user_id")
	return userID, userID != nil
}

// GetUserIDAsUUID extracts user ID from context as UUID
func GetUserIDAsUUID(ctx context.Context) (uuid.UUID, bool) {
	userID := ctx.Value("user_id")
	if userID == nil {
		return uuid.Nil, false
	}

	switch v := userID.(type) {
	case uuid.UUID:
		return v, true
	case string:
		if parsed, err := uuid.Parse(v); err == nil {
			return parsed, true
		}
	}

	return uuid.Nil, false
}

// GetTenantIDAsUUID extracts tenant ID from context as UUID
func GetTenantIDAsUUID(ctx context.Context) (uuid.UUID, bool) {
	tenantID, ok := ctx.Value(XTenantIDKey).(string)
	if !ok || tenantID == "" {
		return uuid.Nil, false
	}

	if parsed, err := uuid.Parse(tenantID); err == nil {
		return parsed, true
	}

	return uuid.Nil, false
}

// Helper function to extract values from gin context-like interface
func extractValue(ginCtx interface{}, key string) (interface{}, bool) {
	if gc, ok := ginCtx.(interface {
		Get(string) (interface{}, bool)
	}); ok {
		return gc.Get(key)
	}
	return nil, false
}

// Backward compatibility - deprecated, use CreateServiceContextFromGin instead
func CopyContextFromContext(ctx context.Context) context.Context {
	appCtx, _ := GetAppContextFromContext(ctx)
	c := WithAppContext(ctx, appCtx)
	c = context.WithValue(c, request_id.XRequestIDKey, ctx.Value(request_id.XRequestIDKey))
	c = context.WithValue(c, XTenantIDKey, ctx.Value(XTenantIDKey))
	c = context.WithValue(c, "user_id", ctx.Value("user_id"))
	return c
}
