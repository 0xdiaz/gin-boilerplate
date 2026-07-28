package middlewares

import (
	"strings"

	"github.com/0xdiaz/gin-boilerplate/internal/app/services/auth"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/0xdiaz/gin-boilerplate/pkg/utils"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates JWT token and protects routes.
//
// Expects Authorization header with format: "Bearer <token>"
// On success, sets "user_id" in gin.Context
// On failure, returns 401 Unauthorized
func AuthMiddleware(authService auth.AuthServicer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Get Authorization header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warnf("missing authorization header")
			utils.Unauthorized(ctx, nil, "Authorization header required")
			ctx.Abort()
			return
		}

		// Check Bearer token format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			logger.Warnf("invalid authorization header format")
			utils.Unauthorized(ctx, nil, "Invalid authorization header format")
			ctx.Abort()
			return
		}

		token := parts[1]

		// Validate token
		userID, err := authService.ValidateToken(token)
		if err != nil {
			logger.Warnf("invalid token: %v", err)
			utils.Unauthorized(ctx, err, "Invalid or expired token")
			ctx.Abort()
			return
		}

		// Set user ID in context for downstream handlers
		ctx.Set("user_id", userID)

		// Continue to next handler
		ctx.Next()
	}
}
