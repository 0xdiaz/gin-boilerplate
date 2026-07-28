package middlewares

import (
	"net/http"

	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
)

const (
	corsAllowMethods  = "GET, POST, PUT, PATCH, DELETE, OPTIONS"
	corsAllowHeaders  = "Origin, Content-Type, Authorization, Accept, Accept-Encoding, X-CSRF-Token, Content-Length"
	corsExposeHeaders = "Content-Length, X-Request-ID"
	corsMaxAge        = "86400"
)

// CORSMiddleware reads the allowed origins from config.Get() (CORS_ALLOWED_ORIGINS).
func CORSMiddleware() gin.HandlerFunc {
	var allowedOrigins []string
	if cfg := config.Get(); cfg != nil {
		allowedOrigins = cfg.Server.CORSAllowedOrigins
	}
	return CORSMiddlewareWithConfig(allowedOrigins)
}

// CORSMiddlewareWithConfig creates a CORS middleware restricted to an explicit
// origin allowlist. It never combines Access-Control-Allow-Origin: * with
// Access-Control-Allow-Credentials, since the two are mutually exclusive per
// the Fetch/CORS spec and this API authenticates via Bearer JWT (not cookies).
//
// Example:
//
//	router.Use(middlewares.CORSMiddlewareWithConfig([]string{"https://app.example.com"}))
func CORSMiddlewareWithConfig(allowedOrigins []string) gin.HandlerFunc {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, o := range allowedOrigins {
		allowed[o] = struct{}{}
	}

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if _, ok := allowed[origin]; ok {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			ctx.Writer.Header().Set("Vary", "Origin")
		}

		ctx.Writer.Header().Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Writer.Header().Set("Access-Control-Expose-Headers", corsExposeHeaders)
		ctx.Writer.Header().Set("Access-Control-Max-Age", corsMaxAge)

		if ctx.Request.Method == http.MethodOptions {
			logger.Debugf("CORS preflight OPTIONS from origin: %s", origin)
			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}
