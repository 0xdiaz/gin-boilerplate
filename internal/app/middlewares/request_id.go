package middlewares

import (
	"context"

	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RequestIDMiddleware adds a unique request ID to each request.
//
// The request ID is:
// - Read from header X-Request-ID or generated using UUID v4
// - Set in gin context as "request_id"
// - Injected into c.Request.Context() for use in services
// - Added to response header as "X-Request-ID"
func RequestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("request_id", requestID)
		c.Request = c.Request.WithContext(context.WithValue(c.Request.Context(), logger.RequestIDContextKey, requestID))
		c.Header("X-Request-ID", requestID)

		c.Next()
	}
}
