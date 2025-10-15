package util

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/dto"
)

// ParsePaginationParams parses pagination parameters from query string using app configuration
func ParsePaginationParams(c *gin.Context) dto.QueryParams {
	// Get app context from gin context values
	appCtxValue, exists := c.Get("app_context")

	// Default pagination values
	defaultLimit := 10
	maxLimit := 100

	// Use app config if available
	if exists {
		if appCtx, ok := appCtxValue.(*AppContext); ok {
			defaultLimit, maxLimit, _ = appCtx.GetPaginationDefaults()
		}
	}

	// Parse page
	page := 1
	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// Parse limit
	limit := defaultLimit
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			if l <= maxLimit {
				limit = l
			} else {
				limit = maxLimit
			}
		}
	}

	return dto.QueryParams{
		Page:    page,
		Limit:   limit,
		Search:  c.Query("search"),
		SortBy:  c.Query("sort_by"),
		SortDir: c.DefaultQuery("sort_dir", "asc"),
	}
}

// FormatTimeResponse formats a time response using the app timezone
func FormatTimeResponse(c *gin.Context, timeStr string) string {
	appCtxValue, exists := c.Get("app_context")
	if !exists {
		return timeStr
	}

	if appCtx, ok := appCtxValue.(*AppContext); ok {
		// Here you could parse and reformat the time in the app timezone
		// For now, just return as is - you can extend this based on your needs
		_ = appCtx // Use appCtx to avoid unused variable error
		return timeStr
	}

	return timeStr
}

// GetAppInfoForResponse returns app information for API responses
func GetAppInfoForResponse(c *gin.Context) map[string]interface{} {
	appCtxValue, exists := c.Get("app_context")
	if !exists {
		return nil
	}

	if appCtx, ok := appCtxValue.(*AppContext); ok {
		return map[string]interface{}{
			"app_name":    appCtx.GetAppName(),
			"app_version": appCtx.GetAppVersion(),
			"timezone":    appCtx.GetTimezone(),
			"locale":      appCtx.GetLocale(),
			"server_time": appCtx.Now(),
		}
	}

	return nil
}
