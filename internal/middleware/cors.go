package middleware

import (
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/config"
)

// CORSMiddleware creates a CORS middleware based on configuration
func CORSMiddleware(corsConfig config.CORSConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip CORS if disabled
			if !corsConfig.Enable {
				return next(c)
			}

			// Set CORS headers
			c.Response().Header().Set("Access-Control-Allow-Origin", corsConfig.AllowedOrigins)
			c.Response().Header().Set("Access-Control-Allow-Methods", corsConfig.AllowedMethods)
			c.Response().Header().Set("Access-Control-Allow-Headers", corsConfig.AllowedHeaders)
			c.Response().Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(corsConfig.AllowCredentials))

			if corsConfig.MaxAgeSeconds > 0 {
				c.Response().Header().Set("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAgeSeconds))
			}

			// Handle preflight requests
			if c.Request().Method == "OPTIONS" {
				return c.NoContent(204)
			}

			return next(c)
		}
	}
}

// CORSWithDefaults creates a CORS middleware with default settings (backward compatibility)
func CORSWithDefaults() echo.MiddlewareFunc {
	defaultConfig := config.CORSConfig{
		Enable:           true,
		AllowCredentials: true,
		AllowedHeaders:   "Origin, Content-Type, Accept, Authorization, X-Request-ID",
		AllowedMethods:   "GET, POST, PUT, DELETE, OPTIONS",
		AllowedOrigins:   "*",
		MaxAgeSeconds:    300,
	}
	return CORSMiddleware(defaultConfig)
}

// ParseOrigins parses comma-separated origins and returns them as a slice
func ParseOrigins(origins string) []string {
	if origins == "" {
		return []string{}
	}

	parts := strings.Split(origins, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}

// ParseMethods parses comma-separated HTTP methods and returns them as a slice
func ParseMethods(methods string) []string {
	if methods == "" {
		return []string{}
	}

	parts := strings.Split(methods, ",")
	result := make([]string, 0, len(parts))

	for _, part := range parts {
		trimmed := strings.ToUpper(strings.TrimSpace(part))
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return result
}
