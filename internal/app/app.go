package app

import (
	"github.com/go-playground/validator/v10"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/protocyber/kelasgo-api/internal/domain/handler"
	"github.com/protocyber/kelasgo-api/internal/domain/repository"
	"github.com/protocyber/kelasgo-api/internal/domain/service"
	"github.com/protocyber/kelasgo-api/internal/infrastructure/database"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// App represents the main application structure
type App struct {
	AuthHandler    *handler.AuthHandler
	UserHandler    *handler.UserHandler
	StudentHandler *handler.StudentHandler
	DBConns        *database.DatabaseConnections
	JWTService     *util.JWTService
	Config         *config.Config
}

// NewApp creates and initializes a new App instance with all dependencies
func NewApp() (*App, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Initialize database connections
	dbConns, err := database.NewConnections(cfg)
	if err != nil {
		return nil, err
	}

	// Initialize validator
	validator := validator.New()

	// Initialize JWT service
	jwtConfig := &config.JWTConfig{
		Secret:     cfg.JWT.Secret,
		ExpireTime: cfg.JWT.ExpireTime,
	}
	jwtService := util.NewJWTService(jwtConfig)

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbConns)
	roleRepo := repository.NewRoleRepository(dbConns)
	tenantUserRepo := repository.NewTenantUserRepository(dbConns)
	tenantUserRoleRepo := repository.NewTenantUserRoleRepository(dbConns)
	studentRepo := repository.NewStudentRepository(dbConns)

	// Initialize services
	authService := service.NewAuthService(userRepo, roleRepo, tenantUserRepo, tenantUserRoleRepo, jwtService)
	userService := service.NewUserService(userRepo, roleRepo, tenantUserRepo, tenantUserRoleRepo)
	studentService := service.NewStudentService(studentRepo, tenantUserRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService, validator)
	userHandler := handler.NewUserHandler(userService, validator)
	studentHandler := handler.NewStudentHandler(studentService, validator)

	// Create and return the app
	return &App{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		StudentHandler: studentHandler,
		DBConns:        dbConns,
		JWTService:     jwtService,
		Config:         cfg,
	}, nil
}
