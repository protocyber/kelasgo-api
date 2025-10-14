package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/rs/zerolog/log"
)

// Logger creates a structured logging middleware
func Logger() gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logEvent := log.Info()
		if param.StatusCode >= 400 {
			logEvent = log.Error()
		}

		logEvent.
			Str("method", param.Method).
			Str("uri", param.Path).
			Str("remote_ip", param.ClientIP).
			Str("user_agent", param.Request.UserAgent()).
			Int("status", param.StatusCode).
			Int64("bytes_in", param.Request.ContentLength).
			Dur("latency", param.Latency).
			Str("latency_human", param.Latency.String()).
			Msg("HTTP Request")

		return ""
	})
}

// RequestLogger returns a customized logger middleware with enhanced request logging
func RequestLogger(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Get status code
		statusCode := c.Writer.Status()

		logEvent := log.Info()
		if statusCode >= 400 {
			logEvent = log.Error()
		}

		// Enhanced logging based on environment
		if cfg.Server.Env == "development" {
			logEvent.
				Str("method", c.Request.Method).
				Str("uri", path).
				Str("remote_ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Int("status", statusCode).
				Int64("bytes_in", c.Request.ContentLength).
				Int("bytes_out", c.Writer.Size()).
				Dur("latency", latency).
				Str("latency_human", latency.String()).
				Msg("HTTP Request")
		} else {
			// More detailed logging for production
			uri := path
			if raw != "" {
				uri = path + "?" + raw
			}

			logEvent.
				Str("method", c.Request.Method).
				Str("uri", uri).
				Str("remote_ip", c.ClientIP()).
				Str("user_agent", c.Request.UserAgent()).
				Str("host", c.Request.Host).
				Str("referer", c.Request.Referer()).
				Int("status", statusCode).
				Int64("bytes_in", c.Request.ContentLength).
				Int("bytes_out", c.Writer.Size()).
				Dur("latency", latency).
				Str("latency_human", latency.String()).
				Str("protocol", c.Request.Proto).
				Msg("HTTP Request")
		}
	}
}
