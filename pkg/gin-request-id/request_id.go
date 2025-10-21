package request_id

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/oklog/ulid/v2"
)

const (
	XRequestIDKey = "X-Request-ID"
)

// generator a function type that returns string.
type generator func() string

var (
	random = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
)

// RequestID is a middleware that injects a 'RequestID' into the context and header of each request.
func RequestID(gen generator) gin.HandlerFunc {
	return func(c *gin.Context) {
		var requestID string
		if gen != nil {
			requestID = gen()
		} else {
			requestID = ulid.Make().String()
		}
		c.Header(XRequestIDKey, requestID)
		c.Set(XRequestIDKey, requestID)
		// fmt.Printf("[GIN-debug] %s [%s] - \"%s %s\"\n", time.Now().Format(time.RFC3339), xRequestID, c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

// GetLoggerConfig return gin.LoggerConfig which will write the logs to specified io.Writer with given gin.LogFormatter.
// By default gin.DefaultWriter = os.Stdout
// reference: https://github.com/gin-gonic/gin#custom-log-format
func GetLoggerConfig(formatter gin.LogFormatter, output io.Writer, skipPaths []string) gin.LoggerConfig {
	if formatter == nil {
		formatter = GetDefaultLogFormatterWithRequestID()
	}
	return gin.LoggerConfig{
		Formatter: formatter,
		Output:    output,
		SkipPaths: skipPaths,
	}
}

// GetDefaultLogFormatterWithRequestID returns gin.LogFormatter with 'RequestID'
func GetDefaultLogFormatterWithRequestID() gin.LogFormatter {
	return func(param gin.LogFormatterParams) string {

		var statusColor, methodColor, resetColor string
		if param.IsOutputColor() {
			statusColor = param.StatusCodeColor()
			methodColor = param.MethodColor()
			resetColor = param.ResetColor()
		}

		if param.Latency > time.Minute {
			param.Latency = param.Latency.Truncate(time.Second)
		}
		// Get request ID from context (param.Keys contains values set via c.Set())
		requestID := ""
		if param.Keys != nil {
			if v, ok := param.Keys[XRequestIDKey]; ok {
				if id, ok := v.(string); ok {
					requestID = id
				}
			}
		}

		return fmt.Sprintf("[GIN] %v |%s|%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
			param.TimeStamp.Format("2006/01/02 - 15:04:05"),
			requestID,
			statusColor, param.StatusCode, resetColor,
			param.Latency,
			param.ClientIP,
			methodColor, param.Method, resetColor,
			param.Path,
			param.ErrorMessage,
		)
	}
}

// GetRequestIDFromContext returns 'RequestID' from the given context if present.
func GetRequestIDFromContext(c *gin.Context) string {
	if v, ok := c.Get(XRequestIDKey); ok {
		if requestID, ok := v.(string); ok {
			return requestID
		}
	}
	return ""
}

// GetRequestIDFromHeaders returns 'RequestID' from the headers if present.
func GetRequestIDFromHeaders(c *gin.Context) string {
	return c.Request.Header.Get(XRequestIDKey)
}
