package routers

import (
	"github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
	"github.com/0xdiaz/gin-boilerplate/internal/app/middlewares"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services/auth"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
	"github.com/0xdiaz/gin-boilerplate/pkg/utils"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes adds all routing list here and delegates to Register* per feature.
// API versioning: /health and /metrics at root; /api/v1/* for versioned API.
func RegisterRoutes(route *gin.Engine) {
	route.NoRoute(func(ctx *gin.Context) {
		utils.NotFound(ctx, nil, "Route not found")
	})

	RegisterHealthRoutes(route)

	apiV1 := route.Group("/api/v1")
	apiV1.Use(middlewares.RateLimitMiddleware())
	userRepo := repositories.NewUserRepository()
	refreshTokenRepo := repositories.NewRefreshTokenRepository()
	authService := auth.NewAuthService(userRepo, refreshTokenRepo, nil)
	exampleService := services.NewExampleService()

	RegisterAuthRoutes(apiV1, authService)
	RegisterExampleRoutes(apiV1, exampleService)

	// Protected routes (require authentication)
	authController := controllers.NewAuthController(authService)
	protectedRoutes := apiV1.Group("")
	protectedRoutes.Use(middlewares.AuthMiddleware(authService))
	{
		protectedRoutes.GET("/profile", authController.Profile)
		protectedRoutes.POST("/logout-all", authController.LogoutAll)
	}
}
