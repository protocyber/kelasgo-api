package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/protocyber/kelasgo-api/internal/util"
)

// BaseHandler provides common functionality for all handlers
type BaseHandler struct {
	log *util.ContextLogger
}

// InitLogger manually sets the context (for backward compatibility or special cases)
func (b *BaseHandler) InitLogger(c *gin.Context) {
	b.log = util.NewContextLogger(c)
}
