package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/app"
	"github.com/protocyber/kelasgo-api/internal/server/middleware"
	request_id "github.com/protocyber/kelasgo-api/pkg/gin-request-id"
)

// SetupRoutes configures all API routes
func SetupRoutes(r *gin.Engine, app *app.App) {
	var (
		cfg            = app.Config
		db             = app.DBConns
		jwtService     = app.JWTService
		authHandler    = app.AuthHandler
		userHandler    = app.UserHandler
		studentHandler = app.StudentHandler
	)

	// Middleware
	r.Use(middleware.AppContextMiddleware(cfg))
	r.Use(request_id.RequestID(nil))
	r.Use(middleware.CORSMiddleware(cfg.App.CORS))
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
	users.Use(middleware.TenantMiddleware(db))
	users.Use(middleware.RequireTenant())
	users.Use(middleware.RoleMiddleware("Admin", "Developer"))
	{
		users.POST("", userHandler.Create)
		users.GET("", userHandler.List)
		users.GET("/:id", userHandler.GetByID)
		users.PUT("/:id", userHandler.Update)
		users.DELETE("/:id", userHandler.Delete)
		users.DELETE("", userHandler.BulkDelete)
	}

	// Student routes (can be accessed by Teachers, Admin, Developer)
	students := protected.Group("/students")
	students.Use(middleware.TenantMiddleware(db))
	students.Use(middleware.RequireTenant())
	students.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		students.POST("", studentHandler.Create)
		students.GET("", studentHandler.List)
		students.GET("/:id", studentHandler.GetByID)
		students.PUT("/:id", studentHandler.Update)
		students.DELETE("/:id", studentHandler.Delete)
		students.DELETE("", studentHandler.BulkDelete)
		students.GET("/class/:class_id", studentHandler.GetByClass)
		students.GET("/parent/:parent_id", studentHandler.GetByParent)
	}

	// Teacher routes (can be accessed by Admin, Developer)
	teachers := protected.Group("/teachers")
	teachers.Use(middleware.TenantMiddleware(db))
	teachers.Use(middleware.RequireTenant())
	teachers.Use(middleware.RoleMiddleware("Admin", "Developer"))
	{
		// TODO: Add teacher handlers
	}

	// Class routes (can be accessed by Teachers, Admin, Developer)
	classes := protected.Group("/classes")
	classes.Use(middleware.TenantMiddleware(db))
	classes.Use(middleware.RequireTenant())
	classes.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add class handlers
	}

	// Subject routes (can be accessed by Teachers, Admin, Developer)
	subjects := protected.Group("/subjects")
	subjects.Use(middleware.TenantMiddleware(db))
	subjects.Use(middleware.RequireTenant())
	subjects.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add subject handlers
	}

	// Attendance routes (can be accessed by Teachers, Admin, Developer)
	attendance := protected.Group("/attendance")
	attendance.Use(middleware.TenantMiddleware(db))
	attendance.Use(middleware.RequireTenant())
	attendance.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add attendance handlers
	}

	// Grade routes (can be accessed by Teachers, Admin, Developer)
	grades := protected.Group("/grades")
	grades.Use(middleware.TenantMiddleware(db))
	grades.Use(middleware.RequireTenant())
	grades.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	{
		// TODO: Add grade handlers
	}

	// Fee routes (can be accessed by Staff, Admin, Developer)
	fees := protected.Group("/fees")
	fees.Use(middleware.TenantMiddleware(db))
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
