# Database Migrations

The **GORM structs are the single source of truth** for the schema. Versioned SQL
migrations are *generated* from those structs with [Atlas](https://atlasgo.io) and
checked into `./migrations`. The app no longer relies on `AutoMigrate` in
production.

## Why

`AutoMigrate` is additive-only: it adds columns/indexes and widens types, but it
**cannot** do type changes (e.g. `bigint → uuid`), column drops, renames, or data
backfills — it errors and crash-loops the service on boot. Versioned migrations
make every schema change explicit, reviewed, ordered, and reversible.

## Layout

| Path | What |
|---|---|
| `cmd/atlas-loader` | Prints the DDL for every model (`bootstrap.AllModels()`). Atlas's schema source. |
| `atlas.hcl` | Atlas config (`env "gorm"`): loader as source, `./migrations` as the dir. |
| `migrations/*.sql` | Versioned migrations + `atlas.sum` (integrity checksum). **Schema only.** |
| `seeds/` + `cmd/seed` | All seed rows, in prod/dev/sample buckets. See `seeds/README.md`. |
| `internal/adapters/database/migrations/` | Boot-time `AutoMigrate` runner (dev/test only, gated by `DB_AUTOMIGRATE`). |
| `.docker/entrypoint.sh` | Applies migrations on container start when `DB_MIGRATE=atlas`. |

## Prerequisites

- The **official** Atlas binary (the community build lacks the GORM provider):
  `curl -sSf https://atlasgo.sh | sh`
- Docker (Atlas spins an ephemeral Postgres to plan migrations).

## Creating a migration

1. Edit the GORM struct(s) as usual.
2. Generate the migration from the model diff:
   ```bash
   atlas migrate diff <short_name> --env gorm
   ```
   → writes `migrations/<timestamp>_<short_name>.sql` and updates `atlas.sum`.
3. **Review / edit the generated SQL.** Atlas nails mechanical changes (add
   column/index, widen type). For anything with intent it can't infer, hand-edit:
   - **Renames** (it emits drop+add → rewrite to `ALTER ... RENAME`).
   - **Data transforms / backfills** (e.g. `bigint → uuid`: add the new column,
     backfill, repoint FKs, drop the old — or, in dev, drop+recreate).
4. Verify it applies and the suite passes against a scratch DB:
   ```bash
   atlas migrate apply --dir file://migrations --url "$SCRATCH_DB_URL"
   go test ./...
   ```
5. Commit the struct change **and** the migration in the same PR.

CI (`.github/workflows/ci.yml`) runs:
- `atlas migrate validate` — directory integrity + `atlas.sum`.
- a **drift check** — regenerates the diff and fails if the models changed without
  a committed migration. (This is the guard that prevents the struct/DB drift that
  previously crash-looped dev.)

## Rollback

- **Incident response = roll back the *app***, not the schema.
- Schema down: `atlas migrate down --env gorm 1` (or `--to-version <v>`). Atlas
  plans the reverse. **A down migration cannot recover dropped data** — for
  destructive changes prefer **expand/contract** (add new → migrate app → drop old
  in a later migration) so app rollback alone is a safe undo.

## Applying in environments

- **Local/dev:** keep `DB_AUTOMIGRATE=true` for fast iteration; run
  `atlas migrate diff` when you're ready to commit a change.
- **Production:** set `DB_MIGRATE=atlas` and `DB_AUTOMIGRATE=false`. The entrypoint
  runs `atlas migrate apply` (advisory-locked) before the app starts; a bad
  migration aborts the deploy with a clear error instead of crash-looping.
- **Locking caveat (hot tables):** a plain `CREATE UNIQUE INDEX` takes a `SHARE`
  lock that blocks writes for the whole build, and `ALTER COLUMN … SET NOT NULL`
  takes `ACCESS EXCLUSIVE` + a full-table scan. On small tables this is sub-second;
  on large production tables it is a write outage. Before applying such a migration
  to a large prod table:
  - Use `CREATE UNIQUE INDEX CONCURRENTLY` (standalone, non-transactional file) instead.
  - Split `SET NOT NULL` into: add `CHECK ... NOT VALID` → `VALIDATE CONSTRAINT` → `SET NOT NULL`.
  - Apply in a maintenance window if neither option is feasible.

## Cutover for an existing database (baseline)

A database created by the old `AutoMigrate` path has the tables but no Atlas
revision history. Baseline it once so Atlas doesn't try to recreate existing
tables:

```bash
# Mark the first migration as already-applied on the existing DB:
atlas migrate apply --dir file://migrations --url "$DB_URL" --baseline <first_migration_timestamp>
```

Then set `DB_MIGRATE=atlas` + `DB_AUTOMIGRATE=false`. Subsequent migrations apply
normally.

## Seed data is not a migration

A migration may not insert rows. Reference/seed data lives in `seeds/` and is
applied by `./seed` (`START_COMMAND=./seed` in the container), which runs
*after* `atlas migrate apply`. The schema migration and the seed run are distinct
steps — never combine them.

`seeds/prod/` and `seeds/dev/` are independent, self-contained environments —
`--env dev` alone is enough to boot a development database without touching the
prod bucket. See `seeds/README.md` for the bucket model and per-file table.

If a migration file was originally created with seed rows that were later moved to
`seeds/`, leave the migration file in place (emptied) — Atlas keys applied
revisions by filename. Removing or renaming the file would break the revision
history on any database that has already applied it. Before emptying a migration
file and changing its hash, verify against your pinned Atlas version that
`migrate status` still reports OK on a database that had applied the original file.

## Note on Atlas licensing

The GORM provider (`external_schema`) requires the **official** Atlas binary
(free, runs offline; Atlas Cloud is **not** required). `migrate diff`, `apply`,
and `validate` work without a login. As of Atlas v0.38 `migrate lint` (deep
destructive-change analysis) requires a free `atlas login` or the Community
build — our CI uses `validate` + the drift check instead, which need no login. If
fully-OSS tooling is a hard requirement, `golang-migrate` (hand-written SQL) is
the alternative.
