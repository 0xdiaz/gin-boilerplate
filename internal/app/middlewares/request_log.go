package middlewares

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/0xdiaz/gin-boilerplate/pkg/utils"
	"github.com/gin-gonic/gin"
)

// responseRecorder wraps gin.ResponseWriter to capture status, size, and body for logging.
type responseRecorder struct {
	gin.ResponseWriter
	body   *bytes.Buffer
	status int
	size   int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	r.body.Write(b)
	n, err := r.ResponseWriter.Write(b)
	r.size += n
	return n, err
}

func (r *responseRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// RequestLogMiddleware logs ARRIVED REQUEST (with masked body/query) and RESPONSE SENT (with masked body).
// Must run after RequestIDMiddleware so request_id is available.
func RequestLogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetString("request_id")
		if requestID == "" {
			c.Next()
			return
		}
		reqLog := logger.WithRequestID(requestID)

		// Read and restore request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			var readErr error
			bodyBytes, readErr = io.ReadAll(c.Request.Body)
			if readErr != nil {
				logger.Warnf("failed to read request body for logging: %v", readErr)
			}
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}
		bodyStr := string(bodyBytes)
		if bodyStr == "" {
			bodyStr = "{}"
		}
		queryStr := c.Request.URL.RawQuery
		maskedBody := utils.MaskSensitiveJSON(bodyStr)
		maskedQuery := utils.MaskQueryString(queryStr)

		reqLog.Infof("ARRIVED REQUEST from IP=%s method=%s path=%s query=%s body=%s",
			c.ClientIP(), c.Request.Method, c.Request.URL.Path, maskedQuery, maskedBody)

		// Wrap writer to capture response
		rec := &responseRecorder{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
			status:         http.StatusOK,
		}
		c.Writer = rec

		c.Next()

		durMs := time.Since(start).Seconds() * 1000
		maskedRespBody := utils.MaskSensitiveJSON(rec.body.String())
		if maskedRespBody == "" {
			maskedRespBody = "{}"
		}
		reqLog.Infof("RESPONSE SENT status=%d duration=%.2fms size=%d bytes body=%s",
			rec.status, durMs, rec.size, maskedRespBody)
	}
}
