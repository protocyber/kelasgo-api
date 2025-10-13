package middleware

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/dto"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// JWTMiddleware creates a JWT authentication middleware
func JWTMiddleware(jwtService *util.JWTService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")

			tokenString, err := util.ExtractTokenFromAuthHeader(authHeader)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.Response{
					Success: false,
					Message: "Unauthorized",
					Error:   err.Error(),
				})
			}

			claims, err := jwtService.ValidateToken(tokenString)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, dto.Response{
					Success: false,
					Message: "Unauthorized",
					Error:   "Invalid or expired token",
				})
			}

			// Set user claims in context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role", claims.Role)
			c.Set("claims", claims)

			return next(c)
		}
	}
}

// RoleMiddleware creates a role-based authorization middleware
func RoleMiddleware(allowedRoles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userRole, ok := c.Get("role").(string)
			if !ok {
				return c.JSON(http.StatusUnauthorized, dto.Response{
					Success: false,
					Message: "Unauthorized",
					Error:   "User role not found in context",
				})
			}

			// Check if user role is in allowed roles
			allowed := false
			for _, role := range allowedRoles {
				if strings.EqualFold(userRole, role) {
					allowed = true
					break
				}
			}

			if !allowed {
				return c.JSON(http.StatusForbidden, dto.Response{
					Success: false,
					Message: "Forbidden",
					Error:   "Insufficient permissions",
				})
			}

			return next(c)
		}
	}
}
