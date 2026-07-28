package routers

import (
	"github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services/auth"
	"github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes under the given group (e.g. /api/v1).
// Full path will be groupPrefix/auth/... Rate limiting is inherited from the parent
// group (see RateLimitMiddleware applied on apiV1 in index.go); do not re-apply it here.
func RegisterAuthRoutes(group *gin.RouterGroup, authService *auth.AuthService) {
	authController := controllers.NewAuthController(authService)
	authRoutes := group.Group("/auth")
	{
		authRoutes.POST("/register", authController.Register)
		authRoutes.POST("/login", authController.Login)
		authRoutes.POST("/refresh", authController.RefreshToken)
		authRoutes.POST("/logout", authController.Logout)
		authRoutes.POST("/forgot-password", authController.ForgotPassword)
		authRoutes.POST("/reset-password", authController.ResetPassword)
	}
}
