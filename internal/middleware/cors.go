package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/config"
)

// CORSMiddleware creates a CORS middleware based on configuration
func CORSMiddleware(corsConfig config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CORS if disabled
		if !corsConfig.Enable {
			c.Next()
			return
		}

		// Set CORS headers
		c.Header("Access-Control-Allow-Origin", corsConfig.AllowedOrigins)
		c.Header("Access-Control-Allow-Methods", corsConfig.AllowedMethods)
		c.Header("Access-Control-Allow-Headers", corsConfig.AllowedHeaders)
		c.Header("Access-Control-Allow-Credentials", strconv.FormatBool(corsConfig.AllowCredentials))

		if corsConfig.MaxAgeSeconds > 0 {
			c.Header("Access-Control-Max-Age", strconv.Itoa(corsConfig.MaxAgeSeconds))
		}

		// Handle preflight requests
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}

// CORSWithDefaults creates a CORS middleware with default settings (backward compatibility)
func CORSWithDefaults() gin.HandlerFunc {
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
