package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/database"
)

// TenantContextKey is the key used to store tenant ID in context
type TenantContextKey string

const TenantIDKey TenantContextKey = "tenant_id"

// TenantMiddleware extracts tenant ID from various sources and adds it to context
// It also sets the PostgreSQL session variable for Row Level Security
func TenantMiddleware(db *database.DatabaseConnections) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			var tenantID uuid.UUID
			var err error

			// Try to get tenant ID from different sources in order of priority:
			// 1. Header (X-Tenant-ID)
			// 2. Query parameter (tenant_id)
			// 3. Subdomain (for subdomain-based tenancy)

			// 1. Check header
			tenantIDStr := c.Request().Header.Get("X-Tenant-ID")

			// 2. Check query parameter if header is empty
			if tenantIDStr == "" {
				tenantIDStr = c.QueryParam("tenant_id")
			}

			// 3. Extract from subdomain if still empty
			if tenantIDStr == "" {
				host := c.Request().Host
				// Example: tenant.example.com -> tenant
				// This is a basic implementation, adjust based on your domain structure
				if subdomain := extractSubdomain(host); subdomain != "" && subdomain != "www" && subdomain != "api" {
					tenantIDStr = subdomain
				}
			}

			// Parse tenant ID if found
			if tenantIDStr != "" {
				tenantID, err = uuid.Parse(tenantIDStr)
				if err != nil {
					// If tenant ID is not a valid UUID, you might want to look it up by name/slug
					// For now, we'll return an error
					return c.JSON(http.StatusBadRequest, map[string]interface{}{
						"error":   "Invalid tenant ID format",
						"message": "Tenant ID must be a valid UUID",
					})
				}

				// Set PostgreSQL session variable for Row Level Security
				if err := setTenantContext(db, tenantID); err != nil {
					return c.JSON(http.StatusInternalServerError, map[string]interface{}{
						"error":   "Failed to set tenant context",
						"message": "Unable to establish tenant isolation",
					})
				}
			}

			// Add tenant ID to context (even if empty - some operations might not require tenant)
			ctx := context.WithValue(c.Request().Context(), TenantIDKey, tenantID)
			c.SetRequest(c.Request().WithContext(ctx))

			// Also set in Echo context for easier access
			c.Set(string(TenantIDKey), tenantID)

			return next(c)
		}
	}
}

// RequireTenant is a middleware that ensures a tenant ID is present
func RequireTenant() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tenantID := c.Get(string(TenantIDKey))
			if tenantID == nil || tenantID == uuid.Nil {
				return c.JSON(http.StatusBadRequest, map[string]interface{}{
					"error":   "Tenant ID required",
					"message": "This operation requires a valid tenant ID",
				})
			}

			return next(c)
		}
	}
}

// GetTenantID extracts tenant ID from Echo context
func GetTenantID(c echo.Context) uuid.UUID {
	if tenantID := c.Get(string(TenantIDKey)); tenantID != nil {
		if tid, ok := tenantID.(uuid.UUID); ok {
			return tid
		}
	}
	return uuid.Nil
}

// GetTenantIDFromContext extracts tenant ID from standard context
func GetTenantIDFromContext(ctx context.Context) uuid.UUID {
	if tenantID, ok := ctx.Value(TenantIDKey).(uuid.UUID); ok {
		return tenantID
	}
	return uuid.Nil
}

// extractSubdomain extracts subdomain from host
// This is a basic implementation - adjust based on your domain structure
func extractSubdomain(host string) string {
	// Remove port if present
	if colonIndex := len(host) - 1; colonIndex > 0 {
		for i := len(host) - 1; i >= 0; i-- {
			if host[i] == ':' {
				host = host[:i]
				break
			}
		}
	}

	// Split by dots and get the first part
	for i := 0; i < len(host); i++ {
		if host[i] == '.' {
			return host[:i]
		}
	}

	// No subdomain found
	return ""
}

// setTenantContext sets the PostgreSQL session variable for Row Level Security
func setTenantContext(db *database.DatabaseConnections, tenantID uuid.UUID) error {
	// Set the current tenant for Row Level Security on both read and write connections
	sql := "SELECT set_config('app.current_tenant', ?, false)"

	// Set on write connection
	if err := db.Write.Exec(sql, tenantID.String()).Error; err != nil {
		return err
	}

	// Set on read connection (if different from write)
	if db.Read != db.Write {
		if err := db.Read.Exec(sql, tenantID.String()).Error; err != nil {
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
