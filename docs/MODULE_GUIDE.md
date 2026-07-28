# Module Guide — Modular Service Layout

> **This is the source of truth for how code is organized in this service.**
> Where older docs (DESIGN_PATTERNS, CODING_STANDARDS, AI_AGENT_RULES) describe an
> `internal/app/...` + `internal/domain/...` layered layout, **this guide overrides them.**
> The boilerplate ships with the legacy layout as a starting point; migrate toward the
> modular layout below as new modules are added.

## TL;DR

- Each service is **one deployable** (one binary, one container, one Postgres).
- Code is organized **by business module** (package-by-feature), not by technical layer.
- A module owns its full vertical slice: `handler → service → repository → model`.
- Modules talk to each other through a **public interface + constructor injection**, never by importing each other's guts.
- Cross-cutting code lives in `pkg/` (the shared kit). Per-service wiring lives in `internal/bootstrap/`.

## Directory layout (target modular structure)

```
.
├── main.go                      # 3 lines: bootstrap.Run()
├── internal/
│   ├── bootstrap/               # the ONLY per-service wiring
│   │   ├── bootstrap.go         #   Run(): config → db → modules → migrate → serve → shutdown
│   │   ├── modules.go           #   Module interface + buildModules() (the list of modules)
│   │   ├── server.go            #   gin engine + global middleware + route mounting
│   │   └── swagger.go           #   OpenAPI/Swagger UI (debug only)
│   ├── adapters/
│   │   └── database/            # DB connect + AutoMigrate runner
│   │       ├── database.go
│   │       ├── migrations/      #   migration.go + sql/ versioned files
│   │       └── seeders/         #   seed data (rows, not schema)
│   └── modules/                 # ← business modules (work happens here)
│       ├── example/             #   THE reference module — copy this to make a new one
│       │   ├── model.go         #     GORM model(s) the module owns
│       │   ├── dto.go           #     request/response types (optional)
│       │   ├── repository.go    #     data access (holds *gorm.DB, no globals)
│       │   ├── service.go       #     business logic (defines the repo interface it needs)
│       │   ├── handler.go       #     HTTP layer (defines the service interface it needs)
│       │   ├── module.go        #     New(db), Name(), Models(), RegisterRoutes(), public API()
│       │   └── service_test.go  #     co-located test (no DB needed — uses a fake repo)
│       ├── auth/                #   real module: JWT auth, exposes Middleware() to others
│       └── health/              #   system module: mounts /health, /metrics at ROOT
└── pkg/                         # ← cross-cutting kit — must never import from internal/
    ├── config/                  #   env-based config struct + SetupConfig()
    ├── logger/                  #   structured logger
    ├── metrics/                 #   Prometheus metrics
    ├── types/                   #   shared error + response types
    └── utils/                   #   response helpers, validator, mask, search
```

## Current boilerplate layout (legacy layered)

The boilerplate ships with this layout for quick-start compatibility:

```
internal/
├── adapters/database/           # DB + migrations + seeders
├── app/
│   ├── controllers/             # HTTP handlers (auth, example, health)
│   ├── dto/                     # request/response structs
│   ├── middlewares/             # auth, cors, metrics, rate_limit, request_id, request_log
│   ├── routers/                 # route registration
│   └── services/                # business logic (auth, example, health)
└── domain/
    ├── models/                  # GORM models
    └── repositories/            # data access
```

This layout is fine for small services. **As the service grows, move each vertical
slice into `internal/modules/<name>/`** following the modular rules below.

## The rules

1. **One module = one Go package = one folder.** An AI agent or a person can own a module without touching others.
2. **A module owns its tables.** Declare them in `Module.Models()`. No cross-module foreign keys; no reaching into another module's tables.
3. **Call other modules only through their public surface** (the exported interface in `module.go`), injected via the constructor. Never import another module's unexported types.
4. **`pkg/` is the shared kit.** It must never import from `internal/`. Anything you'd want identical across repos goes here.
5. **`internal/bootstrap/` is the only place that knows the concrete module list.** Adding a module is one line in `buildModules()`.

## How to add a new module

1. `cp -r internal/modules/example internal/modules/<name>` and rename the package to `<name>`.
2. Replace the model(s) in `model.go`; list them in `Module.Models()`.
3. Implement `repository.go` (data access, holds the injected `*gorm.DB`).
4. Implement `service.go` (business logic). Define the `repository` interface it needs **in this file** (consumer-defined interface → testable without a DB).
5. Implement `handler.go` (HTTP). Define the `service` interface it needs **in this file**.
6. In `module.go`: wire `New(db)`, list `Models()`, mount routes in `RegisterRoutes(api)`, and expose a minimal public `API` interface for other modules.
7. Register it: add `<name>.New(db)` to `buildModules()` in `internal/bootstrap/modules.go`.
8. Add a co-located `*_test.go` using a fake repo (see `example/service_test.go`) — no DB required.

That's it — no other file changes. Routes mount under `/api/v1/...`; migrations pick up your models automatically.

## Cross-module communication (in-process)

A module exposes a small interface; the consumer receives it via its constructor.

```go
// auth exposes its JWT guard and token validation:
authMod := auth.New(db)
authMod.Middleware()   // gin.HandlerFunc — used to protect routes
authMod.Auth()         // auth.Servicer — ValidateToken, etc.

// A module that needs auth gets it injected in buildModules():
//   payments.New(db, authMod.Auth())
```

This keeps modules loosely coupled **and** makes the seam easy to cut later: if a module
must become its own service, the in-process interface call becomes a network call, and
nothing else changes.

## Testing

- **Co-locate tests** with the module (`internal/modules/<name>/*_test.go`).
- Unit-test the service with a **fake repository** (no DB) — see `example/service_test.go`.
- Handler tests use `httptest` + a mocked service interface.
- Shared mocks for legacy tests live in `tests/mocks/`; new modules should prefer in-package fakes.
- `make test` runs `go test ./...`.

## Per-service stamp checklist (when using this boilerplate for a new service)

1. Update `module` in `go.mod` to `github.com/<org>/<service-name>`.
2. Run find/replace `github.com/0xdiaz/gin-boilerplate` → new module path across `*.go`, `*.md`, `go.mod`.
3. Run `gofmt -w .` (import groups may resort after the rename).
4. Verify with `go build ./...`.
5. Rename the DB schema in `.env.example` (`MASTER_DB_SCHEMA`).
6. Update `OTEL_SERVICE_NAME` in `.env.example` and `bootstrap.go`.
7. Delete or replace the `example` module with your first real module.
