# Database Migrations

This directory contains SQL migration files for version-controlled database schema management.

## Current Migration Strategy

The project uses **golang-migrate** to apply these versioned SQL files. `migrations.Migrate()` (in `internal/adapters/database/migrations/migration.go`) runs `migrate.Up()` automatically on app startup and **fails fast** (`main.go` calls `logger.Fatalf`) if any migration fails, so the server never starts against a broken/incomplete schema.

GORM model structs (`internal/domain/models`) are used for querying only — schema changes must be made here as new SQL files, not by relying on GORM AutoMigrate.

## Migration Files Format

Migrations follow the naming convention:
```
{version}_{description}.up.sql   - Forward migration
{version}_{description}.down.sql - Rollback migration
```

Example:
- `000001_create_users_table.up.sql`
- `000001_create_users_table.down.sql`

## Using golang-migrate (Recommended for Production)

### Installation

```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
mv migrate /usr/local/bin/

# Or install as Go package
go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

### Running Migrations

```bash
# Set database URL
export DATABASE_URL="postgres://username:password@localhost:5432/dbname?sslmode=disable"

# Run all pending migrations
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql up

# Rollback last migration
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql down 1

# Check migration version
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql version

# Force set version (use with caution)
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql force 1
```

### Creating New Migrations

```bash
# Create new migration files
migrate create -ext sql -dir internal/adapters/database/migrations/sql -seq create_products_table
```

This will create:
- `000002_create_products_table.up.sql`
- `000002_create_products_table.down.sql`

## Integration with Application

Migrations run automatically on every app startup via `migrations.Migrate()`, which reuses the app's existing GORM/Postgres connection (`database.DB.DB()`) instead of opening a second one:

```go
func Migrate() error {
    sqlDB, _ := database.DB.DB()
    driver, _ := migratepostgres.WithInstance(sqlDB, &migratepostgres.Config{})
    m, _ := migrate.NewWithDatabaseInstance(migrationsPath, "postgres", driver)
    return m.Up() // ErrNoChange treated as success; any other error is fatal to boot
}
```

For manual operations (rollback, scaffolding new migrations), use the Makefile targets, which wrap the `migrate` CLI:

```bash
make migrate-up            # apply pending migrations
make migrate-down          # rollback last migration
make migrate-version       # show current version
make migrate-create NAME=add_products_table
```

## Best Practices

1. **Always write both up and down migrations** - Enable rollback capability
2. **Keep migrations small and focused** - One change per migration
3. **Test migrations on staging first** - Never run untested migrations in production
4. **Never modify existing migrations** - Create new migrations for changes
5. **Use transactions when possible** - Ensure atomic operations
6. **Include comments in SQL** - Document the purpose of changes
7. **Version control everything** - Commit migrations with code changes

## Migration Checklist

Before deploying:
- [ ] Migration files are properly named
- [ ] Both up and down migrations exist
- [ ] Migrations tested on local database
- [ ] Migrations tested on staging environment
- [ ] Backward compatibility verified
- [ ] Rollback procedure documented
- [ ] Team notified of schema changes

## Common Issues

### Migration out of sync
```bash
# Check current version
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql version

# Force version (last resort)
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql force {version}
```

### Dirty migration state
```bash
# Fix dirty migration
migrate -database ${DATABASE_URL} -path internal/adapters/database/migrations/sql force {last_good_version}
```

## Additional Resources

- [golang-migrate documentation](https://github.com/golang-migrate/migrate)
- [GORM AutoMigrate docs](https://gorm.io/docs/migration.html)
- [PostgreSQL ALTER TABLE](https://www.postgresql.org/docs/current/sql-altertable.html)
