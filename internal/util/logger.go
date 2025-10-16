package util

import (
	"os"
	"strings"
	"time"

	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

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
	if cfg.Logger.Format == "console" || cfg.IsDevelopment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		})
	} else {
		// JSON format for production
		log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	}
}
