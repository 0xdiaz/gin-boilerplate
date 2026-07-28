package middlewares_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

const testAllowedOrigin = "http://localhost:3000"

func setupCORSTestRouter(allowedOrigins []string) *gin.Engine {
	router := setupTestRouter()
	router.Use(middlewares.CORSMiddlewareWithConfig(allowedOrigins))
	router.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "ok"})
	})
	return router
}

func TestCORSMiddleware_AllowedOrigin(t *testing.T) {
	router := setupCORSTestRouter([]string{testAllowedOrigin})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", testAllowedOrigin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, testAllowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "Origin", w.Header().Get("Vary"))
}

func TestCORSMiddleware_AllowedOrigin_Preflight(t *testing.T) {
	router := setupCORSTestRouter([]string{testAllowedOrigin})

	req, _ := http.NewRequest(http.MethodOptions, "/test", nil)
	req.Header.Set("Origin", testAllowedOrigin)
	req.Header.Set("Access-Control-Request-Method", http.MethodGet)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Equal(t, testAllowedOrigin, w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Headers"), "Authorization")
}

func TestCORSMiddleware_DisallowedOrigin(t *testing.T) {
	router := setupCORSTestRouter([]string{testAllowedOrigin})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_DisallowedOrigin_HandlerStillReached(t *testing.T) {
	// CORS is enforced by the browser, not the server; a non-preflight request
	// with a disallowed Origin should still reach the handler.
	router := setupCORSTestRouter([]string{testAllowedOrigin})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestCORSMiddleware_NoOriginHeader(t *testing.T) {
	router := setupCORSTestRouter([]string{testAllowedOrigin})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Empty(t, w.Header().Get("Access-Control-Allow-Origin"))
}

func TestCORSMiddleware_NeverSetsCredentialsWithWildcard(t *testing.T) {
	router := setupCORSTestRouter([]string{testAllowedOrigin})

	req, _ := http.NewRequest(http.MethodGet, "/test", nil)
	req.Header.Set("Origin", testAllowedOrigin)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Empty(t, w.Header().Get("Access-Control-Allow-Credentials"), "this API authenticates via Bearer JWT, not cookies")
	assert.NotEqual(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
}
