package routers

import (
	"github.com/0xdiaz/gin-boilerplate/internal/app/middlewares"
	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRoute() *gin.Engine {
	cfg := config.Get()
	debug := cfg != nil && cfg.Server.Debug
	var trustedProxies []string
	if cfg != nil {
		trustedProxies = cfg.Server.TrustedProxies
	}
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()
	router.SetTrustedProxies(trustedProxies)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.CORSMiddleware())
	router.Use(middlewares.RequestIDMiddleware())
	router.Use(middlewares.RequestLogMiddleware())
	router.Use(middlewares.MetricsMiddleware())

	RegisterRoutes(router) //routes register

	// Swagger UI and OpenAPI spec (only when Debug is true)
	RegisterSwaggerRoutes(router, true)

	return router
}
