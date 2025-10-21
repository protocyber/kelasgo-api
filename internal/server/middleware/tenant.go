package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

// TenantMiddleware extracts tenant ID from various sources and adds it to context
// It also sets the PostgreSQL session variable for Row Level Security
func TenantMiddleware(db *database.DatabaseConnections) gin.HandlerFunc {
	return func(c *gin.Context) {
		var tenantID uuid.UUID
		var err error

		// Try to get tenant ID from different sources in order of priority:
		// 1. Header (X-Tenant-ID)
		// 2. Query parameter (tenant_id)
		// 3. Subdomain (for subdomain-based tenancy)

		// 1. Check header
		tenantIDStr := c.GetHeader(string(util.XTenantIDKey))

		// 2. Check query parameter if header is empty
		if tenantIDStr == "" {
			tenantIDStr = c.Query(string(util.TenantIDRequestKey))
		}

		// 3. Extract from subdomain if still empty
		// if tenantIDStr == "" {
		// 	host := c.Request.Host
		// 	// Example: tenant.example.com -> tenant
		// 	// This is a basic implementation, adjust based on your domain structure
		// 	if subdomain := extractSubdomain(host); subdomain != "" && subdomain != "www" && subdomain != "api" {
		// 		tenantIDStr = subdomain
		// 	}
		// }

		// Parse tenant ID if found
		if tenantIDStr != "" {
			tenantID, err = uuid.Parse(tenantIDStr)
			if err != nil {
				log.Error().
					Err(err).
					Str("tenant_id_str", tenantIDStr).
					Str("remote_ip", c.ClientIP()).
					Str("uri", c.Request.URL.Path).
					Str("host", c.Request.Host).
					Msg("Invalid tenant ID format provided")
				// If tenant ID is not a valid UUID, you might want to look it up by name/slug
				// For now, we'll return an error
				c.JSON(http.StatusBadRequest, gin.H{
					"error":   "Invalid tenant ID format",
					"message": "Tenant ID must be a valid UUID",
				})
				c.Abort()
				return
			}

			// Set PostgreSQL session variable for Row Level Security
			if err := setTenantContext(db, tenantID); err != nil {
				log.Error().
					Err(err).
					Str("tenant_id", tenantID.String()).
					Str("remote_ip", c.ClientIP()).
					Str("uri", c.Request.URL.Path).
					Msg("Failed to set tenant context in database")
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "Failed to set tenant context",
					"message": "Unable to establish tenant isolation",
				})
				c.Abort()
				return
			}

			log.Debug().
				Str("tenant_id", tenantID.String()).
				Str("uri", c.Request.URL.Path).
				Msg("Tenant context established successfully")
		}

		// Add tenant ID to context (even if empty - some operations might not require tenant)
		ctx := context.WithValue(c.Request.Context(), util.XTenantIDKey, tenantID)
		c.Request = c.Request.WithContext(ctx)

		// Also set in Gin context for easier access
		c.Set(string(util.XTenantIDKey), tenantID)

		c.Next()
	}
}

// RequireTenant is a middleware that ensures a tenant ID is present
func RequireTenant() gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantID, exists := c.Get(string(util.XTenantIDKey))
		if !exists || tenantID == nil || tenantID == uuid.Nil {
			log.Warn().
				Str("remote_ip", c.ClientIP()).
				Str("uri", c.Request.URL.Path).
				Str("method", c.Request.Method).
				Interface("tenant_id", tenantID).
				Bool("exists", exists).
				Msg("Request blocked due to missing tenant ID")
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Tenant ID required",
				"message": "This operation requires a valid tenant ID",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetTenantID extracts tenant ID from Gin context
func GetTenantID(c *gin.Context) uuid.UUID {
	if tenantID, exists := c.Get(string(util.XTenantIDKey)); exists && tenantID != nil {
		if tid, ok := tenantID.(uuid.UUID); ok {
			return tid
		}
	}
	return uuid.Nil
}

// GetTenantIDFromContext extracts tenant ID from standard context
func GetTenantIDFromContext(ctx context.Context) uuid.UUID {
	if tenantID, ok := ctx.Value(util.XTenantIDKey).(uuid.UUID); ok {
		return tenantID
	}
	return uuid.Nil
}

// extractSubdomain extracts subdomain from host
// This is a basic implementation - adjust based on your domain structure
// func extractSubdomain(host string) string {
// 	// Remove port if present
// 	if colonIndex := len(host) - 1; colonIndex > 0 {
// 		for i := len(host) - 1; i >= 0; i-- {
// 			if host[i] == ':' {
// 				host = host[:i]
// 				break
// 			}
// 		}
// 	}

// 	// Split by dots and get the first part
// 	for i := 0; i < len(host); i++ {
// 		if host[i] == '.' {
// 			return host[:i]
// 		}
// 	}

// 	// No subdomain found
// 	return ""
// }

// setTenantContext sets the PostgreSQL session variable for Row Level Security
func setTenantContext(db *database.DatabaseConnections, tenantID uuid.UUID) error {
	// Set the current tenant for Row Level Security on both read and write connections
	sql := "SELECT set_config('app.current_tenant', ?, false)"

	// Set on write connection
	if err := db.Write.Exec(sql, tenantID.String()).Error; err != nil {
		log.Error().
			Err(err).
			Str("tenant_id", tenantID.String()).
			Str("connection", "write").
			Msg("Failed to set tenant context on write connection")
		return err
	}

	// Set on read connection (if different from write)
	if db.Read != db.Write {
		if err := db.Read.Exec(sql, tenantID.String()).Error; err != nil {
			log.Error().
				Err(err).
				Str("tenant_id", tenantID.String()).
				Str("connection", "read").
				Msg("Failed to set tenant context on read connection")
			return err
		}
	}

	return nil
}

// ClearTenantContext clears the PostgreSQL session variable
func ClearTenantContext(db *database.DatabaseConnections) error {
	// Clear the current tenant context
	sql := "SELECT set_config('app.current_tenant', '', false)"

	// Clear on write connection
	if err := db.Write.Exec(sql).Error; err != nil {
		return err
	}

	// Clear on read connection (if different from write)
	if db.Read != db.Write {
		if err := db.Read.Exec(sql).Error; err != nil {
			return err
		}
	}

	return nil
}
