package metrics

import (
	"sync/atomic"
	"time"
)

var (
	// startTime tracks when the application started
	startTime time.Time

	// totalRequests tracks total number of requests
	totalRequests int64

	// successRequests tracks successful requests (2xx, 3xx)
	successRequests int64

	// errorRequests tracks failed requests (4xx, 5xx)
	errorRequests int64
)

// Init initializes the metrics package.
//
// Should be called once at application startup.
func Init() {
	startTime = time.Now()
}

// RecordRequest records a completed HTTP request.
//
// Increments counters based on HTTP status code.
// Thread-safe using atomic operations.
func RecordRequest(statusCode int) {
	atomic.AddInt64(&totalRequests, 1)

	// Categorize by status code
	if statusCode >= 200 && statusCode < 400 {
		atomic.AddInt64(&successRequests, 1)
	} else if statusCode >= 400 {
		atomic.AddInt64(&errorRequests, 1)
	}
}

// GetTotalRequests returns the total number of requests handled.
func GetTotalRequests() int64 {
	return atomic.LoadInt64(&totalRequests)
}

// GetSuccessRequests returns the number of successful requests.
func GetSuccessRequests() int64 {
	return atomic.LoadInt64(&successRequests)
}

// GetErrorRequests returns the number of failed requests.
func GetErrorRequests() int64 {
	return atomic.LoadInt64(&errorRequests)
}

// GetUptime returns the application uptime in seconds.
func GetUptime() int64 {
	if startTime.IsZero() {
		return 0
	}
	return int64(time.Since(startTime).Seconds())
}

// Reset resets all metrics counters.
//
// Useful for testing purposes.
func Reset() {
	atomic.StoreInt64(&totalRequests, 0)
	atomic.StoreInt64(&successRequests, 0)
	atomic.StoreInt64(&errorRequests, 0)
	startTime = time.Now()
}
