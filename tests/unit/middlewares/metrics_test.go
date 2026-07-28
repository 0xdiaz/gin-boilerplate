package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/middlewares"
	"github.com/0xdiaz/gin-boilerplate/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestMetricsMiddleware tests the metrics collection middleware
func TestMetricsMiddleware(t *testing.T) {
	t.Run("Records successful request (200)", func(t *testing.T) {
		metrics.Reset()

		router := setupTestRouter()
		router.Use(middlewares.MetricsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, int64(1), metrics.GetTotalRequests())
		assert.Equal(t, int64(1), metrics.GetSuccessRequests())
		assert.Equal(t, int64(0), metrics.GetErrorRequests())
	})

	t.Run("Records client error (400)", func(t *testing.T) {
		metrics.Reset()

		router := setupTestRouter()
		router.Use(middlewares.MetricsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, int64(1), metrics.GetTotalRequests())
		assert.Equal(t, int64(0), metrics.GetSuccessRequests())
		assert.Equal(t, int64(1), metrics.GetErrorRequests())
	})

	t.Run("Records server error (500)", func(t *testing.T) {
		metrics.Reset()

		router := setupTestRouter()
		router.Use(middlewares.MetricsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "server error"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, int64(1), metrics.GetTotalRequests())
		assert.Equal(t, int64(0), metrics.GetSuccessRequests())
		assert.Equal(t, int64(1), metrics.GetErrorRequests())
	})

	t.Run("Records multiple requests correctly", func(t *testing.T) {
		metrics.Reset()

		router := setupTestRouter()
		router.Use(middlewares.MetricsMiddleware())
		router.GET("/success", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})
		router.GET("/error", func(c *gin.Context) {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		})

		// Make 3 successful requests
		for i := 0; i < 3; i++ {
			req, _ := http.NewRequest(http.MethodGet, "/success", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		// Make 2 error requests
		for i := 0; i < 2; i++ {
			req, _ := http.NewRequest(http.MethodGet, "/error", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		assert.Equal(t, int64(5), metrics.GetTotalRequests())
		assert.Equal(t, int64(3), metrics.GetSuccessRequests())
		assert.Equal(t, int64(2), metrics.GetErrorRequests())
	})

	t.Run("Records 2xx status codes as success", func(t *testing.T) {
		metrics.Reset()

		statusCodes := []int{200, 201, 202, 204}
		for _, code := range statusCodes {
			router := setupTestRouter()
			router.Use(middlewares.MetricsMiddleware())

			currentCode := code // Capture current value
			router.GET("/test", func(c *gin.Context) {
				c.Status(currentCode)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		assert.Equal(t, int64(4), metrics.GetTotalRequests())
		assert.Equal(t, int64(4), metrics.GetSuccessRequests())
		assert.Equal(t, int64(0), metrics.GetErrorRequests())
	})

	t.Run("Records 3xx status codes as success", func(t *testing.T) {
		metrics.Reset()

		router := setupTestRouter()
		router.Use(middlewares.MetricsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "/new-location")
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, int64(1), metrics.GetTotalRequests())
		assert.Equal(t, int64(1), metrics.GetSuccessRequests())
		assert.Equal(t, int64(0), metrics.GetErrorRequests())
	})

	t.Run("Records 4xx and 5xx as errors", func(t *testing.T) {
		metrics.Reset()

		errorCodes := []int{400, 401, 403, 404, 409, 429, 500, 502, 503}
		for _, code := range errorCodes {
			router := setupTestRouter()
			router.Use(middlewares.MetricsMiddleware())

			currentCode := code // Capture current value
			router.GET("/test", func(c *gin.Context) {
				c.Status(currentCode)
			})

			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
		}

		assert.Equal(t, int64(9), metrics.GetTotalRequests())
		assert.Equal(t, int64(0), metrics.GetSuccessRequests())
		assert.Equal(t, int64(9), metrics.GetErrorRequests())
	})

	t.Run("Does not affect request handling", func(t *testing.T) {
		metrics.Reset()

		router := setupTestRouter()
		router.Use(middlewares.MetricsMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Verify response is not affected
		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "success")

		// Verify metrics recorded
		assert.Equal(t, int64(1), metrics.GetTotalRequests())
	})
}

// TestMetricsThreadSafety tests concurrent access to metrics
func TestMetricsThreadSafety(t *testing.T) {
	metrics.Reset()

	router := setupTestRouter()
	router.Use(middlewares.MetricsMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Make 100 concurrent requests
	done := make(chan bool)
	for i := 0; i < 100; i++ {
		go func() {
			req, _ := http.NewRequest(http.MethodGet, "/test", nil)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			done <- true
		}()
	}

	// Wait for all requests to complete
	for i := 0; i < 100; i++ {
		<-done
	}

	// All 100 requests should be counted
	assert.Equal(t, int64(100), metrics.GetTotalRequests(), "All concurrent requests should be counted")
	assert.Equal(t, int64(100), metrics.GetSuccessRequests())
}

// BenchmarkMetricsMiddleware benchmarks the metrics middleware
func BenchmarkMetricsMiddleware(b *testing.B) {
	metrics.Reset()

	router := setupTestRouter()
	router.Use(middlewares.MetricsMiddleware())
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
