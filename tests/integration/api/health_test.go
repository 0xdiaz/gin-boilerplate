package api_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/routers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// TestHealthEndpoint tests GET /health without initializing database.
// When DB is not initialized, health returns 503 and status "unhealthy".
// For 200 "healthy", run with database configured (e.g. in CI with test DB).
func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	routers.RegisterHealthRoutes(router)

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	router.ServeHTTP(w, req)

	// Without DB init, we expect 503 unhealthy; with DB, would be 200 healthy
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusServiceUnavailable,
		"expected 200 or 503, got %d", w.Code)
	assert.Contains(t, w.Body.String(), "status",
		"response body should contain status field")
}
