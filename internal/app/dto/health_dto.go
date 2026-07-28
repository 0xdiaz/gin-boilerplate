package dto

import "time"

// HealthResponse represents the health check status response.
//
// Returns overall status and individual dependency checks.
type HealthResponse struct {
	// Status is either "healthy" or "unhealthy"
	Status string `json:"status"`

	// Timestamp of when the health check was performed
	Timestamp time.Time `json:"timestamp"`

	// Checks contains individual dependency health status
	Checks map[string]string `json:"checks"`

	// Uptime in seconds since application started
	Uptime int64 `json:"uptime_seconds,omitempty"`
}

// MetricsResponse represents basic application metrics.
//
// Returns request counters and basic statistics.
type MetricsResponse struct {
	// TotalRequests is the total number of requests handled
	TotalRequests int64 `json:"total_requests"`

	// SuccessRequests is the number of successful requests (2xx, 3xx)
	SuccessRequests int64 `json:"success_requests"`

	// ErrorRequests is the number of failed requests (4xx, 5xx)
	ErrorRequests int64 `json:"error_requests"`

	// UptimeSeconds is time since application started
	UptimeSeconds int64 `json:"uptime_seconds"`

	// Timestamp of when metrics were collected
	Timestamp time.Time `json:"timestamp"`
}
