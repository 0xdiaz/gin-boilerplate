package controllers

import (
	"github.com/0xdiaz/gin-boilerplate/internal/app/services"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/0xdiaz/gin-boilerplate/pkg/utils"
	"github.com/gin-gonic/gin"
)

// HealthController handles health check and metrics endpoints.
type HealthController struct {
	service *services.HealthService
}

// NewHealthController creates a new HealthController instance.
func NewHealthController(service *services.HealthService) *HealthController {
	return &HealthController{
		service: service,
	}
}

// Health performs application health check.
//
// GET /health
// Returns health status of the application and its dependencies.
func (ctrl *HealthController) Health(c *gin.Context) {
	ctx, start := logger.LogStart(c.Request.Context(), "HealthController.Health")

	response := ctrl.service.CheckHealth(ctx)

	if response.Status == "unhealthy" {
		logger.LogFinish(ctx, "HealthController.Health", nil, start)
		utils.ServiceUnavailableWithData(c, response, "Service is unhealthy")
		return
	}

	logger.LogFinish(ctx, "HealthController.Health", nil, start)
	utils.Ok(c, response, "Service is healthy")
}

// Metrics returns application metrics.
//
// GET /metrics
// Returns basic request counters and uptime statistics.
func (ctrl *HealthController) Metrics(c *gin.Context) {
	ctx, start := logger.LogStart(c.Request.Context(), "HealthController.Metrics")

	response := ctrl.service.GetMetrics(ctx)
	logger.LogFinish(ctx, "HealthController.Metrics", nil, start)
	utils.Ok(c, response, "Metrics retrieved successfully")
}
