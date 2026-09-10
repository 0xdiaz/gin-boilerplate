# Database Migrations

Schema changes are **versioned SQL files** applied by [golang-migrate](https://github.com/golang-migrate/migrate).
GORM models are used for querying only — they never drive the schema.

For the full command reference (installing the CLI, forcing a dirty version, etc.), see
[`internal/adapters/database/migrations/sql/README.md`](../internal/adapters/database/migrations/sql/README.md).

## Why

`AutoMigrate` is additive-only: it adds columns, indexes and widens types, but it never drops a
column, never renames one, and never tells you what it is about to do. That is fine for a
scratch database and unacceptable for one holding real rows.

Versioned SQL gives us the three things AutoMigrate cannot:

- **A record of what ran** — the `schema_migrations` table names the current version.
- **A rollback path** — every change ships with its `.down.sql`.
- **Review** — the exact DDL that will hit production is in the diff.

It is also lock-protected, so several instances starting at once cannot race each other.

## Layout

| Path | Purpose |
|---|---|
| `internal/adapters/database/migrations/migration.go` | The runner. `Migrate()` reuses the open GORM connection and applies everything pending. |
| `internal/adapters/database/migrations/sql/` | The versioned files: `NNNNNN_description.up.sql` and `.down.sql`. |
| `internal/adapters/database/seeders/` | Development-only demo data. **Not** migrations — see below. |
| `main.go` | Calls `migrations.Migrate()` at startup and treats failure as fatal. |

Migrations run **automatically on every boot**, before the HTTP server starts. A failed
migration is fatal by design: the server must never serve requests against a schema that is
missing columns the code expects.

## Creating a migration

1. Pick the next version number — zero-padded, six digits, one higher than the last file in
   `sql/`.
2. Create **both** files:

```
internal/adapters/database/migrations/sql/000007_create_products_table.up.sql
internal/adapters/database/migrations/sql/000007_create_products_table.down.sql
```

3. Write the forward change:

```sql
-- Create products table
-- Migration: 000007_create_products_table
-- Created to match internal/domain/models/product_model.go

CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    -- Money is an integer in the smallest currency unit. Never FLOAT/DOUBLE.
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_products_name ON products (name);

COMMENT ON TABLE products IS 'Sellable products';
```

4. Write the rollback:

```sql
-- Rollback products table creation
-- Migration: 000007_create_products_table

DROP TABLE IF EXISTS products;
```

5. Keep the GORM model in `internal/domain/models/` in sync with the DDL — column names,
   nullability, and types.

6. Verify by starting the app against a scratch database and checking the applied version:

```bash
docker compose --env-file .env -f .docker/docker-compose-dev.yml up -d postgres_db
go run main.go
```

### Rules

- **Never edit a migration that has already been applied** anywhere. Add a new pair instead.
- **Every `.up.sql` has a `.down.sql`.** A rollback that only says `DROP TABLE` is fine; a
  missing one is not.
- **Constraints belong in the migration**, not only in Go. A `CHECK` on a status column and a
  `UNIQUE` on a natural key are enforced by the database, so a bug in the application cannot
  write a row that violates them.
- **Keep status constants in sync** with their `CHECK` constraint — e.g.
  `models.RefreshTokenStatusAvailable` and the `CHECK (status IN ('available','sold'))` in
  `000002_add_user_auth_tokens.up.sql`.

## Rollback

`Migrate()` only ever moves forward. To step back, use the golang-migrate CLI against the
same directory:

```bash
migrate -path internal/adapters/database/migrations/sql \
        -database "$DATABASE_URL" down 1
```

If a migration fails halfway, the database is left marked **dirty** and the app refuses to
start. Fix the SQL, then clear the flag with `migrate ... force <version>` before retrying.
See the [sql/README.md](../internal/adapters/database/migrations/sql/README.md) for details.

## Seed data is not a migration

Demo rows live in `internal/adapters/database/seeders/` and run from `main.go` **only when
`APP_ENV=development`**. Seeders must be idempotent — they check whether the data already
exists before inserting, so restarting the app never duplicates rows.

Reference data that must exist in every environment (lookup tables, roles) does not belong in
a seeder either; it belongs in a migration, so it is versioned and reviewable like any other
schema change.
