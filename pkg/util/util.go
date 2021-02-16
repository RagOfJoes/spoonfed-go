package util

import (
	"context"

	"github.com/gin-gonic/gin"
)

// AddToContext adds a key/value pair on to context
func AddToContext(c *gin.Context, key ContextKey, value interface{}) *gin.Context {
	// return c.Request.WithContext(context.WithValue(c.Request.Context(), key, value))
	ctx := context.WithValue(c.Request.Context(), key, value)
	c.Request = c.Request.WithContext(ctx)
	return c
}
