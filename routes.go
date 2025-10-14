package main

import (
	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/handler"
	"github.com/protocyber/kelasgo-api/internal/middleware"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	e *echo.Echo,
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	jwtService *util.JWTService,
	dbConns *database.DatabaseConnections,
) {
	// Middleware
	e.Use(middleware.RequestLogger(cfg))
	e.Use(middleware.CORSMiddleware(cfg.App.CORS))
	e.Use(middleware.TenantMiddleware(dbConns))

	// API group
	api := e.Group("/v1")

	// Health check
	api.GET("/health", func(c echo.Context) error {
		return c.JSON(200, map[string]string{
			"status":  "healthy",
			"message": "KelasGo API is running",
		})
	})

	// Auth routes (public)
	auth := api.Group("/auth")
	auth.POST("/login", authHandler.Login)
	auth.POST("/register", authHandler.Register)

	// Protected routes
	protected := api.Group("")
	protected.Use(middleware.JWTMiddleware(jwtService))

	// Auth protected routes
	authProtected := protected.Group("/auth")
	authProtected.POST("/change-password", authHandler.ChangePassword)

	// User routes (Admin and Developer only)
	users := protected.Group("/users")
	users.Use(middleware.RoleMiddleware("Admin", "Developer"))
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)

	// Student routes (can be accessed by Teachers, Admin, Developer)
	students := protected.Group("/students")
	students.Use(middleware.RequireTenant())
	students.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	// TODO: Add student handlers

	// Teacher routes (can be accessed by Admin, Developer)
	teachers := protected.Group("/teachers")
	teachers.Use(middleware.RequireTenant())
	teachers.Use(middleware.RoleMiddleware("Admin", "Developer"))
	// TODO: Add teacher handlers

	// Class routes (can be accessed by Teachers, Admin, Developer)
	classes := protected.Group("/classes")
	classes.Use(middleware.RequireTenant())
	classes.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	// TODO: Add class handlers

	// Subject routes (can be accessed by Teachers, Admin, Developer)
	subjects := protected.Group("/subjects")
	subjects.Use(middleware.RequireTenant())
	subjects.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	// TODO: Add subject handlers

	// Attendance routes (can be accessed by Teachers, Admin, Developer)
	attendance := protected.Group("/attendance")
	attendance.Use(middleware.RequireTenant())
	attendance.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	// TODO: Add attendance handlers

	// Grade routes (can be accessed by Teachers, Admin, Developer)
	grades := protected.Group("/grades")
	grades.Use(middleware.RequireTenant())
	grades.Use(middleware.RoleMiddleware("Teacher", "Admin", "Developer"))
	// TODO: Add grade handlers

	// Fee routes (can be accessed by Staff, Admin, Developer)
	fees := protected.Group("/fees")
	fees.Use(middleware.RequireTenant())
	fees.Use(middleware.RoleMiddleware("Staff", "Admin", "Developer"))
	// TODO: Add fee handlers

	// Notification routes (can be accessed by all authenticated users)
	// notifications := protected.Group("/notifications")
	// TODO: Add notification handlers

	// Dashboard routes (role-based access)
	// dashboard := protected.Group("/dashboard")
	// TODO: Add dashboard handlers with role-specific data
}
