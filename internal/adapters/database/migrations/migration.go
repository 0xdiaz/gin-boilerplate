package migrations

import (
	"errors"
	"fmt"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// migrationsPath points to the versioned SQL migration files. Each schema
// change must be added here as a new {version}_{name}.up.sql/.down.sql pair;
// see sql/README.md for the workflow.
const migrationsPath = "file://internal/adapters/database/migrations/sql"

// Migrate applies all pending versioned SQL migrations using golang-migrate.
//
// It reuses the already-open GORM connection instead of dialing a second one.
// Unlike GORM AutoMigrate, this is versioned (schema_migrations table),
// lock-protected (safe under concurrent instance startup), and supports
// rollback via the paired .down.sql files.
//
// Callers must treat a non-nil error as fatal: starting the server against a
// database that failed to migrate risks serving requests against a schema
// missing columns/tables the app code depends on.
func Migrate() error {
	sqlDB, err := database.DB.DB()
	if err != nil {
		return fmt.Errorf("get underlying sql.DB: %w", err)
	}

	driver, err := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
	if err != nil {
		return fmt.Errorf("initialize migrator: %w", err)
	}

	logger.Infof("Starting database migrations...")

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	logger.Infof("Database migrations completed successfully")
	return nil
}
