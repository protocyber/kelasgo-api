package server

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/app"
	"github.com/rs/zerolog/log"
)

// HTTPServer wraps the HTTP server functionality
type HTTPServer struct {
	server *http.Server
	router *gin.Engine
}

// SetupRoutesFunc is a function type for setting up routes
type SetupRoutesFunc func(*gin.Engine, *app.App)

// NewHTTPServer creates a new HTTP server instance
func NewHTTPServer(app *app.App, setupRoutes SetupRoutesFunc) *HTTPServer {
	// Set Gin mode based on environment
	cfg := app.Config
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	// Create Gin router
	g := gin.New()
	g.Use(gin.Logger())

	// Setup routes
	setupRoutes(g, app)

	// Get server address
	serverAddr := cfg.GetServerAddress()

	// Create HTTP server
	srv := &http.Server{
		Addr:    serverAddr,
		Handler: g,
	}

	return &HTTPServer{
		server: srv,
		router: g,
	}
}

// Start starts the HTTP server
func (s *HTTPServer) Start() error {
	log.Info().Msgf("Starting server on %s", s.server.Addr)
	if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully shuts down the HTTP server
func (s *HTTPServer) Shutdown(ctx context.Context) error {
	log.Info().Msg("Shutting down server...")
	return s.server.Shutdown(ctx)
}

// GetRouter returns the Gin router instance
func (s *HTTPServer) GetRouter() *gin.Engine {
	return s.router
}
