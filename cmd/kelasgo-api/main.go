package main

import (
	"github.com/protocyber/kelasgo-api/internal/app"
	"github.com/protocyber/kelasgo-api/internal/server"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the application with all dependencies
	application, err := app.NewApp()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	// Setup logger
	util.SetupLogger(application.Config)

	// Perform database health check
	if err := application.DBConns.HealthCheck(); err != nil {
		log.Fatal().Err(err).Msg("Database health check failed")
	}
	log.Info().Msg("Database connections healthy")

	// Create and start server
	srv := server.New(application, server.SetupRoutes)

	// Start server with graceful shutdown handling
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
