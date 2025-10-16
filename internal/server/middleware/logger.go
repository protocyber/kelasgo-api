package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/config"
	"github.com/rs/zerolog/log"
)

// RequestLogger returns a customized logger middleware with enhanced request logging and context
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

			// Add request ID to log only if error occurs
			if requestID := GetRequestID(c); requestID != "" {
				logEvent = logEvent.Str("request_id", requestID)
			}

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
