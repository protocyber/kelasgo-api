package util

import (
	"context"
	"time"

	"github.com/protocyber/kelasgo-api/internal/config"
)

// AppContextKey defines keys for context values
type AppContextKey string

const (
	AppConfigKey AppContextKey = "app_config"
	TimezoneKey  AppContextKey = "timezone"
	LocaleKey    AppContextKey = "locale"
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
	return context.WithValue(ctx, AppConfigKey, appCtx)
}

// GetAppContextFromContext extracts app context from context
func GetAppContextFromContext(ctx context.Context) (*AppContext, bool) {
	appCtx, ok := ctx.Value(AppConfigKey).(*AppContext)
	return appCtx, ok
}
