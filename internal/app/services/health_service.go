package services

import (
	"context"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/0xdiaz/gin-boilerplate/pkg/metrics"
)

// HealthChecker is a dependency health check (e.g. database, cache, HTTP).
// Check returns status "ok" or "error" and an optional error.
type HealthChecker interface {
	Check(ctx context.Context) (string, error)
}

// HealthService handles health check operations.
type HealthService struct {
	checkers map[string]HealthChecker
}

// NewHealthService creates a new HealthService instance.
// Use AddChecker to register dependency checkers (e.g. database, Redis).
func NewHealthService() *HealthService {
	return &HealthService{
		checkers: make(map[string]HealthChecker),
	}
}

// AddChecker registers a health checker by name.
func (s *HealthService) AddChecker(name string, c HealthChecker) {
	if s.checkers == nil {
		s.checkers = make(map[string]HealthChecker)
	}
	s.checkers[name] = c
}

// DatabaseChecker implements HealthChecker for the application database.
var _ HealthChecker = (*DatabaseChecker)(nil)

type DatabaseChecker struct{}

func (DatabaseChecker) Check(ctx context.Context) (string, error) {
	if database.DB == nil {
		logger.Warnf("health check: database not initialized")
		return "error", nil
	}
	sqlDB, err := database.DB.DB()
	if err != nil {
		logger.Errorf("health check: failed to get database instance: %v", err)
		return "error", err
	}
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	if err := sqlDB.PingContext(pingCtx); err != nil {
		logger.Errorf("health check: database ping failed: %v", err)
		return "error", err
	}
	return "ok", nil
}

// CheckHealth performs health checks on registered dependencies.
//
// Returns "healthy" only if all checkers return "ok".
func (s *HealthService) CheckHealth(ctx context.Context) *dto.HealthResponse {
	ctx, start := logger.LogStart(ctx, "HealthService.CheckHealth")

	checks := make(map[string]string)
	for name, checker := range s.checkers {
		status, _ := checker.Check(ctx)
		checks[name] = status
	}

	overallStatus := "healthy"
	for _, status := range checks {
		if status != "ok" {
			overallStatus = "unhealthy"
			break
		}
	}

	logger.LogFinish(ctx, "HealthService.CheckHealth", nil, start)
	return &dto.HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now(),
		Checks:    checks,
		Uptime:    metrics.GetUptime(),
	}
}

// GetMetrics returns application metrics.
//
// Returns request counters and uptime statistics.
func (s *HealthService) GetMetrics(ctx context.Context) *dto.MetricsResponse {
	ctx, start := logger.LogStart(ctx, "HealthService.GetMetrics")

	logger.LogFinish(ctx, "HealthService.GetMetrics", nil, start)
	return &dto.MetricsResponse{
		TotalRequests:   metrics.GetTotalRequests(),
		SuccessRequests: metrics.GetSuccessRequests(),
		ErrorRequests:   metrics.GetErrorRequests(),
		UptimeSeconds:   metrics.GetUptime(),
		Timestamp:       time.Now(),
	}
}
