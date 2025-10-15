//go:build wireinject
// +build wireinject

//go:generate wire

package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/protocyber/kelasgo-api/internal/database"
	"github.com/protocyber/kelasgo-api/internal/handler"
	"github.com/protocyber/kelasgo-api/internal/repository"
	"github.com/protocyber/kelasgo-api/internal/service"
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

// ProviderSet contains all the wire providers
var ProviderSet = wire.NewSet(
	// Config
	config.Load,

	// Database
	database.NewConnections,

	// Validator
	ProvideValidator,

	// JWT Service
	ProvideJWTConfig,
	util.NewJWTService,

	// Repositories
	repository.NewUserRepository,
	repository.NewRoleRepository,
	repository.NewTenantUserRepository,
	repository.NewTenantUserRoleRepository,
	repository.NewStudentRepository,

	// Services
	service.NewAuthService,
	service.NewUserService,
	service.NewStudentService,

	// Handlers
	handler.NewAuthHandler,
	handler.NewUserHandler,
	handler.NewStudentHandler,

	// App
	NewApp,
)

// ProvideJWTConfig extracts JWT config from main config
func ProvideJWTConfig(cfg *config.Config) *config.JWTConfig {
	return (*config.JWTConfig)(&cfg.JWT)
}

// ProvideValidator creates a new validator instance
func ProvideValidator() *validator.Validate {
	return validator.New()
}

// NewApp creates a new App instance
func NewApp(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	studentHandler *handler.StudentHandler,
	dbConns *database.DatabaseConnections,
	jwtService *util.JWTService,
	cfg *config.Config,
) *App {
	return &App{
		AuthHandler:    authHandler,
		UserHandler:    userHandler,
		StudentHandler: studentHandler,
		DBConns:        dbConns,
		JWTService:     jwtService,
		Config:         cfg,
	}
}

// InitializeApp initializes the application with all dependencies
func InitializeApp() (*App, error) {
	wire.Build(ProviderSet)
	return &App{}, nil
}
