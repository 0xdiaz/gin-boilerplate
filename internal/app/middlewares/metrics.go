package middlewares

import (
	"github.com/0xdiaz/gin-boilerplate/pkg/metrics"
	"github.com/gin-gonic/gin"
)

// MetricsMiddleware records request metrics.
//
// Records:
// - Total request count
// - Success count (2xx, 3xx status codes)
// - Error count (4xx, 5xx status codes)
//
// Metrics are stored in memory and can be retrieved via /metrics endpoint.
// Uses atomic operations for thread-safety with minimal overhead.
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Process request
		c.Next()

		// Record metrics after request completes
		statusCode := c.Writer.Status()
		metrics.RecordRequest(statusCode)
	}
}
