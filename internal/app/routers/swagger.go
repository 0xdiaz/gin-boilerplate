package routers

import (
	"net/http"

	"github.com/0xdiaz/gin-boilerplate/api"
	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/gin-gonic/gin"
)

// swaggerUIHTML is a minimal HTML page that loads Swagger UI from CDN and points to our OpenAPI spec.
const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>API Documentation</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function() {
      window.ui = SwaggerUIBundle({
        url: "/swagger/openapi.yaml",
        dom_id: "#swagger-ui",
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIBundle.SwaggerUIStandalonePreset
        ]
      });
    };
  </script>
</body>
</html>
`

// RegisterSwaggerRoutes registers routes for OpenAPI spec and Swagger UI.
// When onlyInDebug is true, routes are only registered if config.Server.Debug is true.
func RegisterSwaggerRoutes(router *gin.Engine, onlyInDebug bool) {
	if onlyInDebug {
		cfg := config.Get()
		if cfg == nil || !cfg.Server.Debug {
			return
		}
	}

	// Serve OpenAPI spec (embedded from api package)
	router.GET("/swagger/openapi.yaml", func(c *gin.Context) {
		c.Header("Content-Type", "application/x-yaml")
		c.Data(http.StatusOK, "application/x-yaml", api.OpenAPIYAML)
	})

	// Serve Swagger UI (HTML that loads spec from /swagger/openapi.yaml)
	router.GET("/swagger/index.html", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerUIHTML)
	})
	router.GET("/swagger/", func(c *gin.Context) {
		c.Header("Content-Type", "text/html; charset=utf-8")
		c.String(http.StatusOK, swaggerUIHTML)
	})
}
