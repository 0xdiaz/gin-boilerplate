package routers

import (
	"github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RegisterHealthRoutes registers health check and metrics endpoints.
//
// These routes are public (no authentication required) and used for:
// - Kubernetes/Docker health checks
// - Load balancer health probes
// - Monitoring systems
// - Observability metrics
func RegisterHealthRoutes(router *gin.Engine) {
	healthService := services.NewHealthService()
	healthService.AddChecker("database", &services.DatabaseChecker{})
	healthController := controllers.NewHealthController(healthService)

	// Health check routes (no middleware needed)
	router.GET("/health", healthController.Health)
	router.GET("/metrics", healthController.Metrics)
}
