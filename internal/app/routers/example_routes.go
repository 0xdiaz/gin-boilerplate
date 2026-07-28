package routers

import (
	"github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services"
	"github.com/gin-gonic/gin"
)

// RegisterExampleRoutes registers example/datatables routes under the given group (e.g. /api/v1).
// Full path will be groupPrefix/datatables.
func RegisterExampleRoutes(group *gin.RouterGroup, exampleService *services.ExampleService) {
	exampleController := controllers.NewExampleController(exampleService)
	group.GET("/datatables", exampleController.GetDataDatatables)
}
