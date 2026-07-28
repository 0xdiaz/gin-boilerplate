package database_test

import (
	"os"
	"testing"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDatabaseConnection pings the database after connecting.
// Skip if TEST_DB_MASTER_DSN is not set (e.g. in CI without test DB).
// DSN format: postgres://user:password@host:port/dbname?sslmode=disable
// or use libpq: host=... user=... password=... dbname=... port=... sslmode=disable
func TestDatabaseConnection(t *testing.T) {
	dsn := os.Getenv("TEST_DB_MASTER_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_MASTER_DSN not set, skipping database integration test")
	}

	err := database.DbConnection(dsn, dsn)
	require.NoError(t, err)

	sqlDB, err := database.GetDB().DB()
	require.NoError(t, err)
	require.NotNil(t, sqlDB)
	err = sqlDB.Ping()
	assert.NoError(t, err)
}
