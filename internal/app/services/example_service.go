package services

import (
	"context"

	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
)

// ExampleService handles example/datatables-related business logic.
type ExampleService struct{}

// NewExampleService creates a new ExampleService instance.
func NewExampleService() *ExampleService {
	return &ExampleService{}
}

// GetDataDatatables returns data in DataTables format for server-side processing.
//
// Expects DataTables request query params from gin.Context; returns datatables result or error.
func (s *ExampleService) GetDataDatatables(ctx context.Context, c *gin.Context) (data interface{}, err error) {
	ctx, start := logger.LogStart(ctx, "ExampleService.GetDataDatatables")

	data, err = repositories.GetDataDatatables(c)
	logger.LogFinish(ctx, "ExampleService.GetDataDatatables", err, start)
	return data, err
}
