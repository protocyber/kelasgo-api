package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

// AppWithDB extends App interface to include database connections
type AppWithDB interface {
	App
	GetDBConns() interface{ Close() error }
}

// Server represents the main application server
type Server struct {
	httpServer *HTTPServer
	app        App
}

// New creates a new server instance
func New(app App, setupRoutes SetupRoutesFunc) *Server {
	httpServer := NewHTTPServer(app, setupRoutes)

	return &Server{
		httpServer: httpServer,
		app:        app,
	}
}

// Start starts the server and handles graceful shutdown
func (s *Server) Start() error {
	// Start HTTP server in a goroutine
	serverErr := make(chan error, 1)
	go func() {
		if err := s.httpServer.Start(); err != nil {
			serverErr <- err
		}
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case <-quit:
		log.Info().Msg("Shutting down server...")
		return s.shutdown()
	}
}

// shutdown handles graceful shutdown of the server
func (s *Server) shutdown() error {
	// Create context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Gracefully shutdown the HTTP server
	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("Failed to gracefully shutdown HTTP server")
		return err
	}

	// Close database connections if the app has them
	if dbApp, ok := s.app.(AppWithDB); ok {
		if err := dbApp.GetDBConns().Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close database connections")
			// Don't return error here, just log it
		}
	}

	log.Info().Msg("Server shutdown complete")
	return nil
}
