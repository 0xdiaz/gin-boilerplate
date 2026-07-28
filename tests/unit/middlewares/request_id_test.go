package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestRequestIDMiddleware tests the request ID middleware
func TestRequestIDMiddleware(t *testing.T) {
	t.Run("Generates request ID if not present", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RequestIDMiddleware())

		var capturedRequestID string
		router.GET("/test", func(c *gin.Context) {
			capturedRequestID = c.GetString("request_id")
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check request ID was generated and stored in context
		assert.NotEmpty(t, capturedRequestID, "Request ID should be generated")

		// Check UUID format (36 characters with dashes)
		assert.Len(t, capturedRequestID, 36, "Request ID should be UUID format")
		assert.Contains(t, capturedRequestID, "-", "Request ID should contain dashes")
	})

	t.Run("Adds request ID to response header", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RequestIDMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		responseID := w.Header().Get("X-Request-ID")
		assert.NotEmpty(t, responseID, "X-Request-ID header should be present in response")
		assert.Len(t, responseID, 36, "X-Request-ID should be UUID format")
	})

	t.Run("Uses existing request ID from header", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RequestIDMiddleware())

		var capturedRequestID string
		router.GET("/test", func(c *gin.Context) {
			capturedRequestID = c.GetString("request_id")
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		existingID := "test-request-id-123"
		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		req.Header.Set("X-Request-ID", existingID)

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Should use existing ID, not generate new one
		assert.Equal(t, existingID, capturedRequestID, "Should use existing request ID from header")
		assert.Equal(t, existingID, w.Header().Get("X-Request-ID"), "Should echo back existing request ID")
	})

	t.Run("Request ID is accessible in handlers", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RequestIDMiddleware())

		var handlerAccessedID string
		router.GET("/test", func(c *gin.Context) {
			handlerAccessedID = c.GetString("request_id")
			c.JSON(http.StatusOK, gin.H{"request_id": handlerAccessedID})
		})

		req, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.NotEmpty(t, handlerAccessedID, "Handler should be able to access request ID")
		assert.Equal(t, w.Header().Get("X-Request-ID"), handlerAccessedID, "Handler and header should have same ID")
	})

	t.Run("Different requests get different IDs", func(t *testing.T) {
		router := setupTestRouter()
		router.Use(middlewares.RequestIDMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "ok"})
		})

		// Make two requests
		req1, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w1 := httptest.NewRecorder()
		router.ServeHTTP(w1, req1)

		req2, _ := http.NewRequest(http.MethodGet, "/test", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		id1 := w1.Header().Get("X-Request-ID")
		id2 := w2.Header().Get("X-Request-ID")

		assert.NotEqual(t, id1, id2, "Different requests should get different IDs")
		assert.NotEmpty(t, id1)
		assert.NotEmpty(t, id2)
	})
}

// BenchmarkRequestIDMiddleware benchmarks the request ID middleware
func BenchmarkRequestIDMiddleware(b *testing.B) {
	router := setupTestRouter()
	router.Use(middlewares.RequestIDMiddleware())
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
