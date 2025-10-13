package middleware

import (
	"os"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Logger creates a structured logging middleware
func Logger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			latency := time.Since(start)

			logEvent := log.Info()
			if err != nil {
				logEvent = log.Error().Err(err)
			}

			logEvent.
				Str("method", req.Method).
				Str("uri", req.RequestURI).
				Str("remote_ip", c.RealIP()).
				Str("user_agent", req.UserAgent()).
				Int("status", res.Status).
				Int64("bytes_in", req.ContentLength).
				Int64("bytes_out", res.Size).
				Dur("latency", latency).
				Str("latency_human", latency.String()).
				Msg("HTTP Request")

			return err
		}
	}
}

// RequestLogger returns a customized logger middleware with enhanced request logging
func RequestLogger(cfg *config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()

			err := next(c)

			req := c.Request()
			res := c.Response()

			latency := time.Since(start)

			logEvent := log.Info()
			if err != nil {
				logEvent = log.Error().Err(err)
			}

			// Enhanced logging based on environment
			if cfg.Server.Env == "development" {
				logEvent.
					Str("method", req.Method).
					Str("uri", req.RequestURI).
					Str("remote_ip", c.RealIP()).
					Str("user_agent", req.UserAgent()).
					Int("status", res.Status).
					Int64("bytes_in", req.ContentLength).
					Int64("bytes_out", res.Size).
					Dur("latency", latency).
					Str("latency_human", latency.String()).
					Msg("HTTP Request")
			} else {
				// More detailed logging for production
				logEvent.
					Str("method", req.Method).
					Str("uri", req.RequestURI).
					Str("remote_ip", c.RealIP()).
					Str("user_agent", req.UserAgent()).
					Str("host", req.Host).
					Str("referer", req.Referer()).
					Int("status", res.Status).
					Int64("bytes_in", req.ContentLength).
					Int64("bytes_out", res.Size).
					Dur("latency", latency).
					Str("latency_human", latency.String()).
					Str("protocol", req.Proto).
					Msg("HTTP Request")
			}

			return err
		}
	}
}

// SetupLogger configures zerolog based on configuration
func SetupLogger(cfg *config.Config) {
	zerolog.TimeFieldFormat = time.RFC3339

	// Set log level
	level := strings.ToLower(cfg.Logger.Level)
	switch level {
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel) // default
	}

	// Set log format
	if cfg.Logger.Format == "console" || cfg.Server.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// JSON format for production
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
