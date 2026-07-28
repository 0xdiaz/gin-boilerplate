package controllers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services"
	"github.com/0xdiaz/gin-boilerplate/pkg/metrics"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// setupHealthTestRouter creates a test Gin router for health checks
func setupHealthTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}

// TestHealthController_Health tests the health check endpoint
func TestHealthController_Health(t *testing.T) {
	metrics.Init()

	service := services.NewHealthService()
	service.AddChecker("database", &services.DatabaseChecker{})
	controller := controllers.NewHealthController(service)

	t.Run("Returns unhealthy status when database not initialized (test env)", func(t *testing.T) {
		router := setupHealthTestRouter()
		router.GET("/health", controller.Health)

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// In test environment without database, should return unhealthy (503)
		assert.Equal(t, http.StatusServiceUnavailable, w.Code)

		var response struct {
			Success bool               `json:"success"`
			Message string             `json:"message"`
			Data    dto.HealthResponse `json:"data"`
			Errors  interface{}        `json:"errors"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.False(t, response.Success)
		assert.Equal(t, "unhealthy", response.Data.Status)
		assert.NotNil(t, response.Data.Checks)
		assert.Equal(t, "error", response.Data.Checks["database"])
		assert.NotZero(t, response.Data.Timestamp)
	})

	t.Run("Returns standard response format", func(t *testing.T) {
		svc := services.NewHealthService()
		svc.AddChecker("database", &services.DatabaseChecker{})
		ctrl := controllers.NewHealthController(svc)
		router := setupHealthTestRouter()
		router.GET("/health", ctrl.Health)

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		// Check standard response fields
		assert.Contains(t, w.Body.String(), "success")
		assert.Contains(t, w.Body.String(), "message")
		assert.Contains(t, w.Body.String(), "data")
		assert.Contains(t, w.Body.String(), "checks")
	})

	t.Run("Includes uptime in response", func(t *testing.T) {
		router := setupHealthTestRouter()
		router.GET("/health", controller.Health)

		req, _ := http.NewRequest(http.MethodGet, "/health", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response struct {
			Data dto.HealthResponse `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, response.Data.Uptime, int64(0))
	})
}

// TestHealthController_Metrics tests the metrics endpoint
func TestHealthController_Metrics(t *testing.T) {
	// Reset and initialize metrics
	metrics.Reset()
	metrics.Init()

	service := services.NewHealthService()
	controller := controllers.NewHealthController(service)

	t.Run("Returns metrics successfully", func(t *testing.T) {
		router := setupHealthTestRouter()
		router.GET("/metrics", controller.Metrics)

		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response struct {
			Success bool                `json:"success"`
			Data    dto.MetricsResponse `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.True(t, response.Success)
		assert.NotNil(t, response.Data)
	})

	t.Run("Metrics include all required fields", func(t *testing.T) {
		router := setupHealthTestRouter()
		router.GET("/metrics", controller.Metrics)

		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response struct {
			Data dto.MetricsResponse `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, response.Data.TotalRequests, int64(0))
		assert.GreaterOrEqual(t, response.Data.SuccessRequests, int64(0))
		assert.GreaterOrEqual(t, response.Data.ErrorRequests, int64(0))
		assert.GreaterOrEqual(t, response.Data.UptimeSeconds, int64(0))
		assert.NotZero(t, response.Data.Timestamp)
	})

	t.Run("Metrics counter increments correctly", func(t *testing.T) {
		metrics.Reset()

		// Simulate some requests
		metrics.RecordRequest(200)
		metrics.RecordRequest(201)
		metrics.RecordRequest(404)
		metrics.RecordRequest(500)

		router := setupHealthTestRouter()
		router.GET("/metrics", controller.Metrics)

		req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		var response struct {
			Data dto.MetricsResponse `json:"data"`
		}

		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, int64(4), response.Data.TotalRequests)
		assert.Equal(t, int64(2), response.Data.SuccessRequests) // 200, 201
		assert.Equal(t, int64(2), response.Data.ErrorRequests)   // 404, 500
	})
}

// BenchmarkHealthController_Health benchmarks the health check endpoint
func BenchmarkHealthController_Health(b *testing.B) {
	metrics.Init()
	service := services.NewHealthService()
	service.AddChecker("database", &services.DatabaseChecker{})
	controller := controllers.NewHealthController(service)

	router := setupHealthTestRouter()
	router.GET("/health", controller.Health)

	req, _ := http.NewRequest(http.MethodGet, "/health", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}

// BenchmarkHealthController_Metrics benchmarks the metrics endpoint
func BenchmarkHealthController_Metrics(b *testing.B) {
	metrics.Init()
	service := services.NewHealthService()
	controller := controllers.NewHealthController(service)

	router := setupHealthTestRouter()
	router.GET("/metrics", controller.Metrics)

	req, _ := http.NewRequest(http.MethodGet, "/metrics", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)
	}
}
