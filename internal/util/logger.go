package util

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	request_id "github.com/protocyber/kelasgo-api/pkg/gin-request-id"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// ContextLogger provides logging methods with automatic request ID inclusion
type ContextLogger struct {
	requestID string
	tenantID  string
	userID    string
}

// NewContextLogger creates a new context logger from gin context
func NewContextLogger(c *gin.Context) *ContextLogger {
	logger := &ContextLogger{}

	// Extract request ID from context
	if requestID, exists := c.Get(request_id.XRequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			logger.requestID = id
		}
	}

	// Extract tenant ID if available
	if tenantID, exists := c.Get(string(XTenantIDKey)); exists {
		if id, ok := tenantID.(string); ok {
			logger.tenantID = id
		}
	}

	// Extract user ID if available
	if userID, exists := c.Get("user_id"); exists {
		switch v := userID.(type) {
		case string:
			logger.userID = v
		case uuid.UUID:
			logger.userID = v.String()
		}
	}

	return logger
}

// NewServiceLogger creates a new context logger from service context
func NewServiceLogger(ctx context.Context) *ContextLogger {
	logger := &ContextLogger{}

	// Extract request ID from context
	if requestID := ctx.Value(request_id.XRequestIDKey); requestID != nil {
		if id, ok := requestID.(string); ok {
			logger.requestID = id
		}
	}

	// Extract tenant ID if available
	if tenantID := ctx.Value(XTenantIDKey); tenantID != nil {
		if id, ok := tenantID.(string); ok {
			logger.tenantID = id
		}
	}

	// Extract user ID if available
	if userID := ctx.Value("user_id"); userID != nil {
		switch v := userID.(type) {
		case string:
			logger.userID = v
		case uuid.UUID:
			logger.userID = v.String()
		}
	}

	return logger
}

// WithTenantID adds tenant ID to the logger context
func (cl *ContextLogger) WithTenantID(tenantID string) *ContextLogger {
	cl.tenantID = tenantID
	return cl
}

// WithUserID adds user ID to the logger context
func (cl *ContextLogger) WithUserID(userID string) *ContextLogger {
	cl.userID = userID
	return cl
}

// enrichEvent adds context fields to a log event
func (cl *ContextLogger) enrichEvent(event *zerolog.Event) *zerolog.Event {
	if cl.requestID != "" {
		event = event.Str("request_id", cl.requestID)
	}
	if cl.tenantID != "" {
		event = event.Str("tenant_id", cl.tenantID)
	}
	if cl.userID != "" {
		event = event.Str("user_id", cl.userID)
	}
	return event
}

// Error creates an error level log event with context
func (cl *ContextLogger) Error() *zerolog.Event {
	return cl.enrichEvent(log.Error())
}

// Warn creates a warn level log event with context
func (cl *ContextLogger) Warn() *zerolog.Event {
	return cl.enrichEvent(log.Warn())
}

// Info creates an info level log event with context
func (cl *ContextLogger) Info() *zerolog.Event {
	return cl.enrichEvent(log.Info())
}

// Debug creates a debug level log event with context
func (cl *ContextLogger) Debug() *zerolog.Event {
	return cl.enrichEvent(log.Debug())
}

// GetRequestID returns the request ID
func (cl *ContextLogger) GetRequestID() string {
	return cl.requestID
}

// LogError is a convenience method for logging errors with standard fields
func (cl *ContextLogger) LogError(err error, msg string, fields map[string]interface{}) {
	logEvent(cl.Error().Err(err), fields).Msg(msg)
}

// LogWarn is a convenience method for logging warnings with standard fields
func (cl *ContextLogger) LogWarn(msg string, fields map[string]interface{}) {
	logEvent(cl.Warn(), fields).Msg(msg)
}

// LogInfo is a convenience method for logging info with standard fields
func (cl *ContextLogger) LogInfo(msg string, fields map[string]interface{}) {
	logEvent(cl.Info(), fields).Msg(msg)
}

func logEvent(event *zerolog.Event, fields map[string]interface{}) *zerolog.Event {
	for key, value := range fields {
		switch v := value.(type) {
		case string:
			event = event.Str(key, v)
		case int:
			event = event.Int(key, v)
		case int64:
			event = event.Int64(key, v)
		case bool:
			event = event.Bool(key, v)
		case float64:
			event = event.Float64(key, v)
		default:
			event = event.Interface(key, v)
		}
	}
	return event
}
