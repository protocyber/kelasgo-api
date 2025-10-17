package middleware

import (
	"github.com/gin-gonic/gin"
)

// RequestID middleware generates a unique request ID for each request
// and adds it to the request context and response header
func RequestID() gin.HandlerFunc {
	return nil
}

// GetRequestID retrieves the request ID from the context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
