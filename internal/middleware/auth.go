package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

// JWTMiddleware creates a JWT authentication middleware
func JWTMiddleware(jwtService *util.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		tokenString, err := util.ExtractTokenFromAuthHeader(authHeader)
		if err != nil {
			log.Warn().
				Err(err).
				Str("remote_ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Str("uri", c.Request.URL.Path).
				Msg("Failed to extract token from authorization header")
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "Unauthorized",
				Error:   err.Error(),
			})
			c.Abort()
			return
		}

		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			log.Warn().
				Err(err).
				Str("remote_ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Str("uri", c.Request.URL.Path).
				Msg("JWT token validation failed")
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Set user claims in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("claims", claims)

		c.Next()
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			log.Error().
				Str("remote_ip", c.ClientIP()).
				Str("uri", c.Request.URL.Path).
				Strs("allowed_roles", allowedRoles).
				Msg("User role not found in context during role check")
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "Unauthorized",
				Error:   "User role not found in context",
			})
			c.Abort()
			return
		}

		role, ok := userRole.(string)
		if !ok {
			log.Error().
				Interface("user_role", userRole).
				Str("remote_ip", c.ClientIP()).
				Str("uri", c.Request.URL.Path).
				Msg("Invalid user role format in context")
			c.JSON(http.StatusUnauthorized, dto.Response{
				Success: false,
				Message: "Unauthorized",
				Error:   "Invalid user role format in context",
			})
			c.Abort()
			return
		}

		// Check if user role is in allowed roles
		allowed := false
		for _, allowedRole := range allowedRoles {
			if strings.EqualFold(role, allowedRole) {
				allowed = true
				break
			}
		}

		if !allowed {
			log.Warn().
				Str("user_role", role).
				Strs("allowed_roles", allowedRoles).
				Str("remote_ip", c.ClientIP()).
				Str("uri", c.Request.URL.Path).
				Msg("Insufficient permissions for role-based access")
			c.JSON(http.StatusForbidden, dto.Response{
				Success: false,
				Message: "Forbidden",
				Error:   "Insufficient permissions",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
