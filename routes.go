package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/handler"
	"github.com/protocyber/kelasgo-api/internal/middleware"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	r *gin.Engine,
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	jwtService *util.JWTService,
	dbConns *database.DatabaseConnections,
) {
	// Middleware
	r.Use(middleware.RequestLogger(cfg))
	r.Use(middleware.CORSMiddleware(cfg.App.CORS))
	// Note: TenantMiddleware is now optional and applied per route group as needed

	// API group
	api := r.Group("/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "healthy",
			"message": "KelasGo API is running",
		})
	})

	// Auth routes (public - no tenant context required)
	auth := api.Group("/auth")
	{
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(jwtService))

	// Auth protected routes (for authenticated users - no tenant context required)
	authProtected := protected.Group("/auth")
	{
		authProtected.POST("/change-password", authHandler.ChangePassword)
		authProtected.GET("/tenants", authHandler.GetUserTenants)      // Get user's available tenants
		authProtected.POST("/select-tenant", authHandler.SelectTenant) // Select a tenant and get new token
	}

	// User routes (Admin and Developer only - requires tenant context)
	users := protected.Group("/users")
	users.Use(middleware.TenantMiddleware(dbConns))
	users.Use(middleware.RequireTenant())
	users.Use(middleware.RoleMiddleware("Admin", "Developer"))
	{
		users.POST("", userHandler.Create)
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
	}

	// Student routes (can be accessed by Teachers, Admin, Developer)
	students := protected.Group("/students")
	students.Use(middleware.TenantMiddleware(dbConns))
	students.Use(middleware.RequireTenant())
	students.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add student handlers
	}

	// Teacher routes (can be accessed by Admin, Developer)
	teachers := protected.Group("/teachers")
	teachers.Use(middleware.TenantMiddleware(dbConns))
	teachers.Use(middleware.RequireTenant())
	teachers.Use(middleware.RoleMiddleware("Admin", "Developer"))
	{
		// TODO: Add teacher handlers
	}

	// Class routes (can be accessed by Teachers, Admin, Developer)
	classes := protected.Group("/classes")
	classes.Use(middleware.TenantMiddleware(dbConns))
	classes.Use(middleware.RequireTenant())
	classes.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add class handlers
	}

	// Subject routes (can be accessed by Teachers, Admin, Developer)
	subjects := protected.Group("/subjects")
	subjects.Use(middleware.TenantMiddleware(dbConns))
	subjects.Use(middleware.RequireTenant())
	subjects.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add subject handlers
	}

	// Attendance routes (can be accessed by Teachers, Admin, Developer)
	attendance := protected.Group("/attendance")
	attendance.Use(middleware.TenantMiddleware(dbConns))
	attendance.Use(middleware.RequireTenant())
	attendance.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add attendance handlers
	}

	// Grade routes (can be accessed by Teachers, Admin, Developer)
	grades := protected.Group("/grades")
	grades.Use(middleware.TenantMiddleware(dbConns))
	grades.Use(middleware.RequireTenant())
	grades.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add grade handlers
	}

	// Fee routes (can be accessed by Staff, Admin, Developer)
	fees := protected.Group("/fees")
	fees.Use(middleware.TenantMiddleware(dbConns))
	fees.Use(middleware.RequireTenant())
	fees.Use(middleware.RoleMiddleware("Staff", "Admin", "Developer"))
	{
		// TODO: Add fee handlers
	}

	// Notification routes (can be accessed by all authenticated users)
	// notifications := protected.Group("/notifications")
	// {
	// 	// TODO: Add notification handlers
	// }

	// Dashboard routes (role-based access)
	// dashboard := protected.Group("/dashboard")
	// {
	// 	// TODO: Add dashboard handlers with role-specific data
	// }
}
