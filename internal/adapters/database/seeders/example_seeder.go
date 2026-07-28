package seeders

import (
	"fmt"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

// demoExamples seed rows for the example CRUD/datatables endpoints.
var demoExamples = []string{
	"Example row 1",
	"Example row 2",
}

func seedDemoExamples() error {
	var count int64
	if err := database.DB.Model(&models.Example{}).Count(&count).Error; err != nil {
		return fmt.Errorf("count examples: %w", err)
	}
	if count > 0 {
		return nil
	}

	rows := make([]models.Example, 0, len(demoExamples))
	for _, data := range demoExamples {
		rows = append(rows, models.Example{Data: data})
	}
	if err := database.DB.Create(&rows).Error; err != nil {
		return fmt.Errorf("create examples: %w", err)
	}

	logger.Infof("Seeded %d demo example rows", len(rows))
	return nil
}
