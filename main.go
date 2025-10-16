package main

import (
	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/server"
	"github.com/protocyber/kelasgo-api/internal/util"
	"github.com/rs/zerolog/log"
)

func main() {
	// Initialize the application with all dependencies
	app, err := InitializeApp()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize application")
	}

	// Setup logger
	util.SetupLogger(app.Config)

	// Perform database health check
	if err := app.DBConns.HealthCheck(); err != nil {
		log.Fatal().Err(err).Msg("Database health check failed")
	}
	log.Info().Msg("Database connections healthy")

	// Create and start server
	srv := server.New(app, func(r *gin.Engine, a server.App) {
		SetupRoutes(r, a.(*App))
	})

	// Start server with graceful shutdown handling
	if err := srv.Start(); err != nil {
		log.Fatal().Err(err).Msg("Failed to start server")
	}
}
