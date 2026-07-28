package middlewares_test

import (
	"github.com/gin-gonic/gin"
)

// setupTestRouter creates a test Gin router
// Shared helper for all middleware tests
func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	return router
}
