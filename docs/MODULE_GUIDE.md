# Project Layout Guide

> **This is the source of truth for how code is organized in this service.**
> It describes the layout that is actually in the repository. If any other document
> disagrees with what is written here, this document wins and the other document is
> the one that needs fixing.

## TL;DR

- One deployable: one binary, one container, one Postgres.
- Code is organized **by technical layer**: HTTP → service → repository → model.
- Wiring happens in **one place**: `internal/app/routers/index.go`.
- Data access goes through the package-level `database.DB` handle in
  `internal/adapters/database`.
- Cross-cutting, domain-agnostic code lives in `pkg/`, which must never import `internal/`.
- Unit tests live in `tests/unit/<layer>/`, with shared fakes in `tests/mocks/`.

## Directory layout

```
.
├── main.go                            # config → db → migrate → seed (dev) → serve → graceful shutdown
├── internal/
│   ├── adapters/
│   │   └── database/
│   │       ├── database.go            #   DbConnection(master, replica), GetDB(), package-level DB
│   │       ├── migrations/
│   │       │   ├── migration.go       #     golang-migrate runner, fail-fast on startup
│   │       │   └── sql/               #     versioned NNNNNN_name.up.sql / .down.sql pairs
│   │       └── seeders/               #   idempotent demo data, development only
│   ├── app/
│   │   ├── controllers/               #   HTTP layer: <name>_controller.go
│   │   ├── dto/                       #   request/response types: <name>_dto.go
│   │   ├── middlewares/               #   auth, cors, metrics, rate_limit, request_id, request_log
│   │   ├── routers/
│   │   │   ├── router.go              #     gin engine + global middleware
│   │   │   ├── index.go               #     RegisterRoutes(): builds repos + services, mounts groups
│   │   │   ├── <name>_routes.go       #     Register<Name>Routes(group, service)
│   │   │   └── swagger.go             #     OpenAPI/Swagger UI (debug only)
│   │   └── services/                  #   business logic: <name>_service.go
│   │       └── auth/                  #     multi-file service package (its own errors + interface)
│   └── domain/
│       ├── models/                    #   GORM models: <name>_model.go, each with TableName()
│       └── repositories/              #   data access: <name>_repo.go
├── pkg/                               # cross-cutting kit — must never import internal/
│   ├── config/  logger/  metrics/  types/  utils/
└── tests/
    ├── unit/{controllers,services,middlewares}/   # package <x>_test, no database
    ├── integration/{api,database}/                # needs a real database
    ├── mocks/                                     # shared in-memory fakes
    └── fixtures/
```

## The rules

1. **One feature = one file per layer**, named after the feature: `auth_controller.go`,
   `auth_service.go`, `user_repo.go`, `user_model.go`, `auth_dto.go`, `auth_routes.go`.
2. **Layers only call downward:** controller → service → repository → model. A controller
   never touches a repository; a repository never imports a service.
3. **Repositories are the only place that touches `database.DB`.** Everything above them
   receives a repository interface.
4. **`internal/app/routers/index.go` is the only place that knows concrete types.** It
   constructs repositories, injects them into services, and hands services to the
   per-feature `Register*Routes` functions.
5. **`pkg/` is the shared kit.** It must never import from `internal/`. Anything you would
   want identical across services goes here.
6. **`/health` and `/metrics` are mounted at the root**, not under `/api/v1`.

## How to add a new feature

Follow the `auth` slice — it is the most recent and closest to these rules.

1. **Migration** — add a versioned pair in `internal/adapters/database/migrations/sql/`:
   `NNNNNN_create_<table>.up.sql` and the matching `.down.sql`. Never edit an applied migration.
2. **Model** — `internal/domain/models/<name>_model.go`, with a `TableName()` method.
   Keep the struct tags in sync with the migration.
3. **Repository** — `internal/domain/repositories/<name>_repo.go`. Declare an exported
   interface, an unexported struct implementing it, and a `New<Name>Repository()`
   constructor that returns the interface. Wrap every error with `fmt.Errorf("...: %w", err)`
   and log it. A missing row is `(nil, nil)`, not an error.
4. **DTO** — `internal/app/dto/<name>_dto.go` for request and response shapes.
5. **Service** — `internal/app/services/<name>_service.go`. A struct holding the repository
   interfaces, a `New<Name>Service(...)` constructor, and `logger.LogStart`/`LogFinish`
   around each exported method. Declare sentinel errors (`ErrUserNotFound`) here.
6. **Controller** — `internal/app/controllers/<name>_controller.go`. A struct holding the
   service, a `New<Name>Controller(service)` constructor, and methods taking `*gin.Context`.
   Map sentinel errors to responses; always answer through `pkg/utils`.
7. **Routes** — `internal/app/routers/<name>_routes.go` with
   `Register<Name>Routes(group *gin.RouterGroup, service *services.<Name>Service)`.
8. **Wire it** — in `internal/app/routers/index.go`: construct the repositories and service,
   then call `Register<Name>Routes(apiV1, <name>Service)`.
9. **Tests** — a fake in `tests/mocks/<name>_repo_mock.go` (with a
   `var _ repositories.<Name>Repository = (*Mock<Name>Repository)(nil)` assertion) and unit
   tests in `tests/unit/services/<name>_service_test.go`.
10. **Seeder** (optional) — `internal/adapters/database/seeders/<name>_seeder.go`, made
    idempotent and called from `seeders.Run()`.

## Wiring example

```go
// internal/app/routers/index.go — the single place that knows concrete types
func RegisterRoutes(route *gin.Engine) {
    RegisterHealthRoutes(route) // root-level probes

    apiV1 := route.Group("/api/v1")
    apiV1.Use(middlewares.RateLimitMiddleware())

    userRepo := repositories.NewUserRepository()
    refreshTokenRepo := repositories.NewRefreshTokenRepository()
    authService := services.NewAuthService(userRepo, refreshTokenRepo)

    RegisterAuthRoutes(apiV1, authService)

    // Protected routes go behind the auth middleware.
    protected := apiV1.Group("")
    protected.Use(middlewares.AuthMiddleware(authService))
}
```

```go
// internal/app/routers/auth_routes.go — one file per feature
func RegisterAuthRoutes(group *gin.RouterGroup, authService *services.AuthService) {
    authController := controllers.NewAuthController(authService)
    authRoutes := group.Group("/auth")
    authRoutes.POST("/login", authController.Login)
}
```

## Testing

- Unit tests live in `tests/unit/<layer>/`, in package `<layer>_test` (black box).
- Services are tested against the in-memory fakes in `tests/mocks/` — no database.
- Integration tests that need a real database live in `tests/integration/`.
- `make test` runs the unit suite; `make test-coverage` produces a coverage report.

## Per-service stamp checklist (when reusing this codebase for a new service)

1. Update `module` in `go.mod` to `github.com/<org>/<service-name>`.
2. Find/replace `github.com/0xdiaz/gin-boilerplate` → the new module path across `*.go`, `*.md`, `go.mod`.
3. Run `gofmt -w .` (import groups resort after the rename).
4. Update `container_name`, volume names, and host ports in `.docker/docker-compose-*.yml`
   so the new service does not collide with other services on the same machine.
5. Verify with `go build ./...`.
6. Update the database name in `.env.example`.
7. Delete or replace the `example` slice with your first real feature.

## Planned direction (not implemented)

A package-by-feature layout — `internal/modules/<name>/{model,repository,service,handler,module}.go`
with per-module route registration and co-located tests — was drafted for this codebase but
**was never implemented**. No `internal/modules/` or `internal/bootstrap/` directory exists.

It is recorded here only so the intent is not lost. Until that migration actually happens,
the layout documented above is the one to follow, and any example showing `internal/modules/...`
is describing code that does not exist.
