package services_test

import (
	"context"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/app/services"
	"github.com/stretchr/testify/assert"
)

// mockHealthChecker implements services.HealthChecker for tests.
type mockHealthChecker struct {
	status string
}

func (m *mockHealthChecker) Check(ctx context.Context) (string, error) {
	return m.status, nil
}

func TestNewHealthService(t *testing.T) {
	svc := services.NewHealthService()
	assert.NotNil(t, svc)
}

func TestHealthService_AddChecker(t *testing.T) {
	svc := services.NewHealthService()
	svc.AddChecker("db", &mockHealthChecker{status: "ok"})
	svc.AddChecker("cache", &mockHealthChecker{status: "error"})
	ctx := context.Background()
	resp := svc.CheckHealth(ctx)
	assert.Equal(t, "ok", resp.Checks["db"])
	assert.Equal(t, "error", resp.Checks["cache"])
}

func TestHealthService_AddChecker_nilMapInit(t *testing.T) {
	svc := &services.HealthService{} // nil checkers map; AddChecker initializes it
	svc.AddChecker("a", &mockHealthChecker{status: "ok"})
	ctx := context.Background()
	resp := svc.CheckHealth(ctx)
	assert.Equal(t, "healthy", resp.Status)
	assert.Equal(t, "ok", resp.Checks["a"])
}

func TestHealthService_CheckHealth_allOk(t *testing.T) {
	svc := services.NewHealthService()
	svc.AddChecker("db", &mockHealthChecker{status: "ok"})
	svc.AddChecker("cache", &mockHealthChecker{status: "ok"})
	ctx := context.Background()
	resp := svc.CheckHealth(ctx)
	assert.NotNil(t, resp)
	assert.Equal(t, "healthy", resp.Status)
	assert.Equal(t, "ok", resp.Checks["db"])
	assert.Equal(t, "ok", resp.Checks["cache"])
}

func TestHealthService_CheckHealth_oneError(t *testing.T) {
	svc := services.NewHealthService()
	svc.AddChecker("db", &mockHealthChecker{status: "ok"})
	svc.AddChecker("cache", &mockHealthChecker{status: "error"})
	ctx := context.Background()
	resp := svc.CheckHealth(ctx)
	assert.NotNil(t, resp)
	assert.Equal(t, "unhealthy", resp.Status)
	assert.Equal(t, "ok", resp.Checks["db"])
	assert.Equal(t, "error", resp.Checks["cache"])
}

func TestHealthService_CheckHealth_noCheckers(t *testing.T) {
	svc := services.NewHealthService()
	ctx := context.Background()
	resp := svc.CheckHealth(ctx)
	assert.NotNil(t, resp)
	assert.Equal(t, "healthy", resp.Status)
	assert.Empty(t, resp.Checks)
}

func TestHealthService_GetMetrics(t *testing.T) {
	svc := services.NewHealthService()
	ctx := context.Background()
	resp := svc.GetMetrics(ctx)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Timestamp)
	// Global metrics; we only assert struct shape
	assert.GreaterOrEqual(t, resp.TotalRequests, int64(0))
	assert.GreaterOrEqual(t, resp.SuccessRequests, int64(0))
	assert.GreaterOrEqual(t, resp.ErrorRequests, int64(0))
	assert.GreaterOrEqual(t, resp.UptimeSeconds, int64(0))
}
