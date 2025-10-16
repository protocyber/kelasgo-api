package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/middleware"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, app *App) {
	// Middleware
	// r.Use(middleware.RequestLogger(app.Config))
	r.Use(middleware.AppContextMiddleware(app.Config)) // Add app context middleware
	r.Use(middleware.CORSMiddleware(app.Config.App.CORS))
	// Note: TenantMiddleware is now optional and applied per route group as needed

	// API group
	api := r.Group("/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		// Example of using app context
		appCtx, _ := middleware.GetAppContext(c)

		response := gin.H{
			"status":  "healthy",
			"message": "KelasGo API is running",
		}

		// Add app info if context is available
		if appCtx != nil {
			response["app"] = gin.H{
				"name":        appCtx.GetAppName(),
				"version":     appCtx.GetAppVersion(),
				"description": appCtx.GetAppDescription(),
				"url":         appCtx.GetAppURL(),
				"timezone":    appCtx.GetTimezone(),
				"locale":      appCtx.GetLocale(),
				"server_time": appCtx.Now(),
			}
		}

		c.JSON(http.StatusOK, response)
	})

	// Auth routes (public - no tenant context required)
	auth := api.Group("/auth")
	{
		auth.POST("/login", app.AuthHandler.Login)
		auth.POST("/register", app.AuthHandler.Register)
	}

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(app.JWTService))

	// Auth protected routes (for authenticated users - no tenant context required)
	authProtected := protected.Group("/auth")
	{
		authProtected.POST("/change-password", app.AuthHandler.ChangePassword)
		authProtected.GET("/tenants", app.AuthHandler.GetUserTenants)      // Get user's available tenants
		authProtected.POST("/select-tenant", app.AuthHandler.SelectTenant) // Select a tenant and get new token
	}

	// User routes (Admin and Developer only - requires tenant context)
	users := protected.Group("/users")
	users.Use(middleware.TenantMiddleware(app.DBConns))
	users.Use(middleware.RequireTenant())
	users.Use(middleware.RoleMiddleware("Admin", "Developer"))
	{
		users.POST("", app.UserHandler.Create)
		users.GET("", app.UserHandler.List)
		users.GET("/:id", app.UserHandler.GetByID)
		users.PUT("/:id", app.UserHandler.Update)
		users.DELETE("/:id", app.UserHandler.Delete)
		users.DELETE("", app.UserHandler.BulkDelete)
	}

	// Student routes (can be accessed by Teachers, Admin, Developer)
	students := protected.Group("/students")
	students.Use(middleware.TenantMiddleware(app.DBConns))
	students.Use(middleware.RequireTenant())
	students.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		students.POST("", app.StudentHandler.Create)
		students.GET("", app.StudentHandler.List)
		students.GET("/:id", app.StudentHandler.GetByID)
		students.PUT("/:id", app.StudentHandler.Update)
		students.DELETE("/:id", app.StudentHandler.Delete)
		students.DELETE("", app.StudentHandler.BulkDelete)
		students.GET("/class/:class_id", app.StudentHandler.GetByClass)
		students.GET("/parent/:parent_id", app.StudentHandler.GetByParent)
	}

	// Teacher routes (can be accessed by Admin, Developer)
	teachers := protected.Group("/teachers")
	teachers.Use(middleware.TenantMiddleware(app.DBConns))
	teachers.Use(middleware.RequireTenant())
	teachers.Use(middleware.RoleMiddleware("Admin", "Developer"))
	{
		// TODO: Add teacher handlers
	}

	// Class routes (can be accessed by Teachers, Admin, Developer)
	classes := protected.Group("/classes")
	classes.Use(middleware.TenantMiddleware(app.DBConns))
	classes.Use(middleware.RequireTenant())
	classes.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add class handlers
	}

	// Subject routes (can be accessed by Teachers, Admin, Developer)
	subjects := protected.Group("/subjects")
	subjects.Use(middleware.TenantMiddleware(app.DBConns))
	subjects.Use(middleware.RequireTenant())
	subjects.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add subject handlers
	}

	// Attendance routes (can be accessed by Teachers, Admin, Developer)
	attendance := protected.Group("/attendance")
	attendance.Use(middleware.TenantMiddleware(app.DBConns))
	attendance.Use(middleware.RequireTenant())
	attendance.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add attendance handlers
	}

	// Grade routes (can be accessed by Teachers, Admin, Developer)
	grades := protected.Group("/grades")
	grades.Use(middleware.TenantMiddleware(app.DBConns))
	grades.Use(middleware.RequireTenant())
	grades.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add grade handlers
	}

	// Fee routes (can be accessed by Staff, Admin, Developer)
	fees := protected.Group("/fees")
	fees.Use(middleware.TenantMiddleware(app.DBConns))
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
