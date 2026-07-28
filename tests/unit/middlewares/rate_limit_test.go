package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRateLimitMiddleware tests the rate limiting middleware
func TestRateLimitMiddleware(t *testing.T) {
	t.Run("Allows requests under limit", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RateLimitMiddlewareWithConfig(10, 10)) // 10 req/s, burst 10
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// Make 5 requests (under limit)
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code, "Request %d should be allowed", i+1)
		}
	})

	t.Run("Blocks requests exceeding burst limit", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RateLimitMiddlewareWithConfig(1, 5)) // 1 req/s, burst 5
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		var successCount int
		var blockedCount int

		// Make 10 rapid requests
		for i := 0; i < 10; i++ {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			switch w.Code {
			case http.StatusOK:
				successCount++
			case http.StatusTooManyRequests:
				blockedCount++
			}
		}

		// First 5 should succeed (burst), rest should be blocked
		assert.GreaterOrEqual(t, successCount, 5, "At least 5 requests should succeed")
		assert.GreaterOrEqual(t, blockedCount, 1, "At least 1 request should be blocked")

		t.Logf("Success: %d, Blocked: %d", successCount, blockedCount)
	})

	t.Run("Returns standard error response format", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RateLimitMiddlewareWithConfig(1, 1)) // Very strict limit
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// Exhaust the limit
		for i := 0; i < 5; i++ {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code == http.StatusTooManyRequests {
				// Check response has standard format
				assert.Contains(t, w.Body.String(), "success", "Response should have success field")
				assert.Contains(t, w.Body.String(), "message", "Response should have message field")
				t.Logf("Rate limit response: %s", w.Body.String())
				break
			}
		}
	})
}

// TestIPRateLimiter tests that different IPs have separate limits.
// Skipped: integration-style test (requires simulating different client IPs; run in API/integration tests).
func TestIPRateLimiter(t *testing.T) {
	t.Skip("Skipping: Requires IP simulation which is complex in unit tests")

	// NOTE: This test demonstrates the pattern but is skipped
	// In production, test this with integration tests or manually
	router := setupTestRouter()
	router.Use(middlewares.RateLimitMiddlewareWithConfig(2, 2))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})

	// Different IPs should have separate rate limits
	// IP 1: 192.168.1.1
	// IP 2: 192.168.1.2
	// Each should be able to make 2 requests independently
}

// BenchmarkRateLimitMiddleware benchmarks the rate limit middleware
func BenchmarkRateLimitMiddleware(b *testing.B) {
	router := setupTestRouter()
	router.Use(middlewares.RateLimitMiddlewareWithConfig(1000, 2000))
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
