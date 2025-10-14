package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/middleware"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the application with all dependencies
	app, err := InitializeApp()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	// Setup logger
	middleware.SetupLogger(app.Config)

	// Perform database health check
	if err := app.DBConns.HealthCheck(); err != nil {
		log.Fatal().Err(err).Msg("Database health check failed")
	}
	log.Info().Msg("Database connections healthy")

	// Create Echo instance
	e := echo.New()

	// Setup routes
	SetupRoutes(e, app.Config, app.AuthHandler, app.UserHandler, app.JWTService, app.DBConns)

	// Get server address
	serverAddr := app.Config.GetServerAddress()

	// Start server in a goroutine
	go func() {
		log.Info().Msgf("Starting server on %s", serverAddr)
		if err := e.Start(serverAddr); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Shutting down server...")

	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the server
	if err := e.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Failed to gracefully shutdown server")
	}

	// Close database connections
	if err := app.DBConns.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close database connections")
	}

	log.Info().Msg("Server shutdown complete")
}
