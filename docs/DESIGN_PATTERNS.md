# Go Project Design Patterns & Architecture Blueprint

**Version:** 2.0
**Last Updated:** 2026-06-10
**Project:** Go Gin Modular Service — Universal Starter Kit

---

## ⚠️ FOR AI AGENTS - READ THIS FIRST

> **🚨 SOURCE OF TRUTH: [`MODULE_GUIDE.md`](./MODULE_GUIDE.md) describes how code is organized in this service.**
> Where this document and `MODULE_GUIDE.md` ever disagree on **folder structure**, `MODULE_GUIDE.md` wins.
>
> **The canonical reference is the `auth` slice.** Copy it to create a new feature.
> Every code example below is modelled on it.

### 🔥 Non-negotiable patterns (MUST follow)

- **Organize by technical layer.** A feature is one file per layer, feature-prefixed:
  `internal/app/controllers/<name>_controller.go`, `internal/app/services/<name>_service.go`,
  `internal/domain/repositories/<name>_repo.go`, `internal/domain/models/<name>_model.go`.
- **Interfaces at the data boundary + constructor injection.** Repositories declare an
  exported interface; services depend on that interface and receive it through `New*Service(...)`,
  never reached for via globals.
- **Only repositories touch the database.** Services and controllers never reference
  `database.DB`.
- **All HTTP responses go through `pkg/utils`** (`utils.Ok`, `utils.BadRequest`,
  `utils.RespondWithAPIError`, …). Never write `c.JSON(...)` by hand for API responses.
- **`pkg/` is the shared kit and must never import from `internal/`.**

### 📖 How to Use This Document

1. ✅ Read `MODULE_GUIDE.md` first (the layout source of truth, ~100 lines).
2. ✅ Skim the `auth` slice — it is the living version of this doc.
3. ⚠️  Read the Implementation Patterns and Anti-Patterns sections in full.
4. 📚 Use the rest as reference for detailed patterns.

---

## 📋 Table of Contents

> **Navigation:** anchors are stable; use Ctrl+F on a keyword or click the section links.

### 🔥 Critical Sections (MUST READ)

| Section | Keywords |
|---------|----------|
| [6. Implementation Patterns](#6-implementation-patterns) | `controller`, `service`, `repository`, `routes` |
| → [6.1 Controller Pattern](#61-controller-pattern) | `Controller`, `New*Controller`, `injected service` |
| → [6.2 Service Pattern](#62-service-pattern) | `Service`, `NewService`, `repository interface` |
| → [6.3 Repository Pattern](#63-repository-pattern) | `Repository`, `NewRepository(db)`, injected `*gorm.DB` |
| → [6.4 Route Wiring](#64-route-wiring-pattern) | `Register*Routes`, `index.go`, `New*` |
| [13. Anti-Patterns](#13-anti-patterns-to-avoid) | `wrong`, `bad`, `avoid`, `anti-pattern` |

### 📚 All Sections

| # | Section | Keywords |
|---|---------|----------|
| 1 | [Overview](#1-overview) | `philosophy`, `goals`, `principles` |
| 1.1 | [Architecture Philosophy](#11-architecture-philosophy) | `package-by-feature`, `why patterns` |
| 2 | [Project Architecture](#2-project-architecture) | `modules`, `vertical slice`, `separation` |
| 2.1 | [The Layers](#21-the-layers) | `controller`, `service`, `repository`, `model` |
| 2.2 | [Dependency Flow Rules](#22-dependency-flow-rules) | `dependency`, `direction`, `flow`, `layering` |
| 3 | [Core Design Patterns](#3-core-design-patterns) | `patterns`, `repository`, `service`, `DI` |
| 3.1 | [Repository Pattern](#31-repository-pattern) | `repository`, `data access`, `injected db` |
| 3.2 | [Service Layer Pattern](#32-service-layer-pattern) | `service`, `business logic`, `orchestration` |
| 3.3 | [DTO Pattern](#33-dto-data-transfer-object-pattern) | `DTO`, `request`, `response` |
| 3.4 | [Constructor Pattern](#34-constructor-pattern) | `New`, `constructor`, `injection` |
| 3.5 | [Middleware Pattern](#35-middleware-pattern) | `middleware`, `gin.HandlerFunc`, `auth guard` |
| 3.6 | [Dependency Injection](#36-dependency-injection-pattern) | `DI`, `consumer-defined interface`, `injection` |
| 4 | [Directory Structure](#4-directory-structure-standard) | `directory`, `folder`, `structure`, `tree` |
| 4.1 | [Complete Project Structure](#41-complete-project-structure) | `project tree`, `modular layout` |
| 4.2 | [Package Organization](#42-package-organization-rules) | `package`, `internal`, `pkg`, `layer` |
| 5 | [Layer Responsibilities](#5-layer-responsibilities) | `responsibilities`, `what`, `where` |
| 5.1 | [Controller Layer](#51-controller-layer-thin-layer) | `controller`, `thin`, `HTTP`, `validation` |
| 5.2 | [Service Layer](#52-service-layer-fat-layer) | `service`, `fat`, `business logic` |
| 5.3 | [Repository Layer](#53-repository-layer-data-layer) | `repository`, `database`, `CRUD`, `GORM` |
| 6 | [Implementation Patterns](#6-implementation-patterns) | `how to`, `implementation`, `code` |
| 6.1 | [Controller Pattern](#61-controller-pattern) | `Controller`, `New*Controller`, `methods` |
| 6.2 | [Service Pattern](#62-service-pattern) | `Service`, `NewService`, `methods` |
| 6.3 | [Repository Pattern](#63-repository-pattern) | `Repository`, `injected db`, `CRUD` |
| 6.4 | [Route Wiring Pattern](#64-route-wiring-pattern) | `Register*Routes`, `index.go`, `New*` |
| 6.5 | [Response Utility Pattern](#65-response-utility-pattern) | `utils.Ok`, `response`, `standard format` |
| 7 | [Request Flow Patterns](#7-request-flow-patterns) | `flow`, `request`, `lifecycle`, `pipeline` |
| 7.1 | [Standard CRUD Flow](#71-standard-crud-flow) | `CRUD`, `create`, `read`, `update`, `delete` |
| 7.2 | [Authentication Flow](#72-authentication-flow-pattern) | `auth`, `login`, `JWT`, `token` |
| 7.3 | [Transaction Flow](#73-transaction-flow-pattern) | `transaction`, `Begin()`, `Commit()`, `Rollback()` |
| 8 | [Data Flow Patterns](#8-data-flow-patterns) | `data`, `transformation`, `mapping` |
| 8.1 | [Request → Response Transform](#81-request--response-data-transformation) | `transform`, `DTO`, `model`, `mapping` |
| 8.2 | [DTO vs Model Usage](#82-dto-vs-model-usage) | `when`, `DTO`, `model`, `difference` |
| 9 | [Error Handling Patterns](#9-error-handling-patterns) | `error`, `handling`, `recovery`, `logging` |
| 9.1 | [Error Flow](#91-error-flow-pattern) | `error flow`, `propagation` |
| 9.2 | [Sentinel Errors & APIError Mapping](#92-sentinel-errors--apierror-mapping) | `sentinel error`, `APIError`, `mapping` |
| 9.3 | [Error Wrapping](#93-error-wrapping-pattern) | `%w`, `fmt.Errorf`, `wrapping` |
| 10 | [Testing Patterns](#10-testing-patterns) | `test`, `testing`, `fake`, `coverage` |
| 10.1 | [Service Layer Testing](#101-service-layer-testing-pattern) | `service test`, `fake repo`, `no DB` |
| 10.2 | [Table-Driven Testing](#102-table-driven-testing-pattern) | `table test`, `subtests`, `t.Run` |
| 11 | [Feature Implementation Guide](#11-complete-feature-implementation-guide) | `step by step`, `guide`, `new feature` |
| 11.x | [Step-by-Step: New Feature](#step-by-step-adding-a-new-feature) | `complete example`, `full feature` |
| 12 | [Pattern Examples](#12-pattern-examples-from-codebase) | `examples`, `real code`, `reference` |
| 12.1 | [Auth Pattern Example](#121-auth-pattern-from-the-auth-service) | `auth example`, `authentication` |
| 12.2 | [DataTable Example](#122-datatable-pattern-from-the-example-slice) | `datatable`, `server-side` |
| 13 | [Anti-Patterns to Avoid](#13-anti-patterns-to-avoid) | `wrong`, `bad`, `avoid`, `don't` |
| 13.1 | [Business Logic in Controller](#131--business-logic-in-controller) | `controller anti-pattern`, `fat controller` |
| 13.2 | [Direct DB in Service](#132--direct-database-access-in-service) | `service anti-pattern`, `tight coupling` |
| 13.3 | [Standalone Functions](#133--standalone-handler-functions) | `standalone`, `function anti-pattern` |
| 13.4 | [God Service](#134--god-service-too-many-responsibilities) | `god service`, `SRP violation` |
| 13.5 | [Skipping a Layer](#135--skipping-a-layer) | `layering`, `dependency direction` |

### 🎯 Quick Lookups by Task

**Implementing Handlers (HTTP layer):**
- [6.1 Controller Pattern](#61-controller-pattern)
- [5.1 Controller Responsibilities](#51-controller-layer-thin-layer)
- [13.1 What NOT to do](#131--business-logic-in-controller)

**Implementing Services:**
- [6.2 Service Pattern](#62-service-pattern)
- [5.2 Service Responsibilities](#52-service-layer-fat-layer)
- [13.2 What NOT to do](#132--direct-database-access-in-service)

**Implementing Repositories:**
- [6.3 Repository Pattern](#63-repository-pattern)
- [5.3 Repository Responsibilities](#53-repository-layer-data-layer)
- [3.1 Repository Pattern Theory](#31-repository-pattern)

**Wiring a Feature:**
- [6.4 Route Wiring Pattern](#64-route-wiring-pattern)
- [11 Step-by-Step: New Feature](#step-by-step-adding-a-new-feature)

**Understanding Flow:**
- [7.1 Standard CRUD Flow](#71-standard-crud-flow)
- [7.2 Authentication Flow](#72-authentication-flow-pattern)
- [7.3 Transaction Flow](#73-transaction-flow-pattern)

---

## 1. OVERVIEW

### 1.1 Architecture Philosophy

This project follows a **layered (package-by-layer) architecture**. The service is one
deployable (one binary, one container, one Postgres), and the code is organized **by technical
layer**: HTTP in `internal/app/controllers`, business logic in `internal/app/services`, data
access in `internal/domain/repositories`, and GORM structs in `internal/domain/models`.

The Clean Architecture **principles** — separation of concerns, dependency inversion,
dependencies pointing inward, testability — are the law. Each feature is a horizontal slice
across those directories, with the file name (`auth_controller.go`, `auth_service.go`,
`user_repo.go`) tying the slice together.

```
┌─────────────────────────────────────────────┐
│          External Systems / HTTP            │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│   one feature, spread across the layers      │
│                                             │
│   <name>_controller.go  (Thin Layer)        │  ← HTTP: parse, call service, respond
│        │                                    │
│        ▼                                    │
│   <name>_service.go     (Fat Layer)         │  ← Business rules, orchestration, logging
│        │                                    │
│        ▼                                    │
│   <name>_repo.go        (Data Layer)        │  ← CRUD against database.DB
│        │                                    │
│        ▼                                    │
│   <name>_model.go       (the GORM struct)   │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│              Database (Postgres)             │
└─────────────────────────────────────────────┘
```

**Key Principles:**
1. **One feature = one file per layer**, named after the feature.
2. **Dependency Inversion** — services depend on a repository *interface*, not a concrete type.
3. **Separation of Concerns** — each file in the slice has ONE responsibility.
4. **Dependency Direction** — `controller → service → repository`; never the reverse.
5. **Single wiring point** — `internal/app/routers/index.go` is the only file that knows
   concrete types.
6. **Testability** — services are tested with the fakes in `tests/mocks/`, no DB required.

> The canonical implementation of every principle above is the `auth` slice.
> `MODULE_GUIDE.md` is the source of truth for the layout.

---

## 2. PROJECT ARCHITECTURE

### 2.1 The Layers

Each layer is a directory; a feature contributes one file to each:

```
┌──────────────────────────────────────────────────────┐
│                  PRESENTATION                        │
│  internal/app/controllers/  — HTTP (gin), DTO binding│
│  internal/app/dto/          — request/response types │
└──────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────┐
│                  APPLICATION                          │
│  internal/app/services/     — business logic,        │
│                   orchestration, sentinel errors     │
└──────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────┐
│                  DATA ACCESS                          │
│  internal/domain/repositories/ — CRUD via database.DB│
│  internal/domain/models/       — the GORM structs    │
└──────────────────────────────────────────────────────┘

   internal/app/routers/index.go — wires it all:
   repositories → services → Register<Name>Routes(group)
```

Cross-cutting infrastructure lives outside those layers:

```
┌──────────────────────────────────────────────────────┐
│  pkg/                  (the shared "kit-in-waiting")  │
│  • config/ logger/ metrics/ types/ utils/            │
│    MUST never import internal/                        │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│  internal/app/middlewares/                            │
│  • cors, request_id, request_log, metrics,           │
│    rate_limit, auth                                   │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│  internal/adapters/database/                          │
│  • database.go  — connect + replica resolver         │
│  • migrations/  — golang-migrate + versioned SQL     │
│  • seeders/     — idempotent demo data (dev only)    │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│  main.go — config → db → migrate → seed → serve      │
└──────────────────────────────────────────────────────┘
```

### 2.2 Dependency Flow Rules

**✅ ALLOWED:**
```
Controller → Service → Repository → Database
     ↓           ↓           ↓
    DTO     Utils/Types    Model

internal/app/routers → controllers, services, repositories  // wiring knows concrete types
internal/*           → pkg/*                                // everything uses the shared kit
```

**❌ FORBIDDEN:**
```
Service → Controller            // Services cannot depend on the HTTP layer
Repository → Service         // Repositories cannot depend on business logic
Model → Repository           // Models are pure data structures
pkg/* → internal/*           // The shared kit must NEVER import internal
Controller → Repository      // Never skip the service layer
Service → database.DB        // Only repositories touch the database
```

**🔑 Cardinal wiring rule:** nothing constructs its own dependencies. A service receives the
repository **interfaces** it needs through its constructor, and everything is assembled in one
place:

```go
// in internal/app/routers/index.go
userRepo := repositories.NewUserRepository()
refreshTokenRepo := repositories.NewRefreshTokenRepository()
authService := services.NewAuthService(userRepo, refreshTokenRepo)

RegisterAuthRoutes(apiV1, authService)
```

This keeps every layer testable in isolation: the same service runs against Postgres in
production and against `tests/mocks/` in unit tests, with no code change.

---

## 3. CORE DESIGN PATTERNS

### 3.1 Repository Pattern

**Purpose:** Abstract data access from business logic.

**Structure:** A repository is a **struct that holds an injected `*gorm.DB`** — no package globals.
The *interface* the repository satisfies is **defined by the service that consumes it** (see §3.6),
not by the repository itself.

```go
// File: internal/domain/repositories/user_repo.go
package example

import "gorm.io/gorm"

// userRepo is the data-access layer for users.
// It holds an injected *gorm.DB (no global state) so it is easy to test and reuse.
type Repository struct {
    db *gorm.DB
}

// NewRepository creates a Repository bound to the given connection.
func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

// List returns all example records.
func (r *Repository) List() ([]*Example, error) {
    var list []*Example
    if err := r.db.Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}
```

**✅ DO:**
- Hold the injected `*gorm.DB` on the struct (`NewRepository(db)`); never reach for a global.
- Keep repositories focused on data access ONLY.
- Use GORM or parameterized queries (never string concatenation).
- Return models from internal/domain/models, not DTOs.
- Translate database-specific errors here (e.g. `gorm.ErrRecordNotFound` → `nil, nil`).
- Use transactions for multi-step operations.

**❌ DON'T:**
- Use the package-level `database.DB` global (it exists only for the connection lifecycle/tests).
- Put business logic in repositories.
- Reach past the repository to build queries elsewhere.
- Import the service layer.
- Transform data for API responses (that is the DTO's job).

---

### 3.2 Service Layer Pattern

**Purpose:** Encapsulate business logic and orchestrate operations.

**Structure:** A service holds a **consumer-defined `repository` interface** (declared in
`service.go` itself, see §3.6). It logs each operation with `logger.LogStart`/`LogFinish` and
propagates the request-scoped `context.Context`.

```go
// File: internal/app/services/auth/auth_service.go
package example

import (
    "context"

    "github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

// repository is the data-access contract this service needs.
// Defining it HERE (at the consumer) keeps the service testable with fakes.
type repository interface {
    List() ([]*Example, error)
}

// AuthService holds the authentication business logic.
type Service struct {
    repo repository
}

// NewService creates a Service backed by the given repository.
func NewService(repo repository) *Service {
    return &Service{repo: repo}
}

// List returns all example records.
func (s *Service) List(ctx context.Context) ([]*Example, error) {
    ctx, start := logger.LogStart(ctx, "AuthService.List")
    list, err := s.repo.List()
    logger.LogFinish(ctx, "AuthService.List", err, start)
    return list, err
}
```

**✅ DO:**
- Implement ALL business logic and rules here.
- Define the `repository` interface this service needs *in this file* (dependency inversion).
- Accept and propagate `context.Context` as the first argument.
- Wrap each operation in `logger.LogStart` / `logger.LogFinish` for tracing.
- Validate business constraints and orchestrate multiple repository calls.
- Return service-defined sentinel errors (e.g. `ErrInvalidCredentials`) for known conditions.

**❌ DON'T:**
- Handle HTTP concerns (`gin.Context`, status codes) — except where a repository helper such as
  DataTables genuinely needs `*gin.Context` for server-side paging.
- Access the database directly (go through the repository).
- Import the handler layer.
- Return HTTP responses.

---

### 3.3 DTO (Data Transfer Object) Pattern

**Purpose:** Define contracts for API requests/responses and prevent tight coupling.

**Structure:** DTOs live in `internal/app/dto/<name>_dto.go` (the file is optional; small modules
may keep request/response types alongside the model). They are pure data structures with binding tags.

```go
// File: internal/app/dto/auth_dto.go
package auth

// RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// AuthResponse represents the response after successful authentication.
type AuthResponse struct {
    User         UserResponse `json:"user"`
    AccessToken  string       `json:"access_token"`
    RefreshToken string       `json:"refresh_token"`
    TokenType    string       `json:"token_type"`
}

// UserResponse represents user information in API responses (never includes password).
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

**✅ DO:**
- Define separate DTOs for Request and Response.
- Use struct tags for validation (`binding:"required,email"`, …).
- Keep DTOs in `internal/app/dto/<name>_dto.go`, never inside a model file.
- Exclude sensitive fields from responses (e.g. the auth `User` model tags `Password` as `json:"-"`).

**❌ DON'T:**
- Expose database models directly via API when they carry secrets.
- Add business logic to DTOs.
- Reuse a request DTO as a response DTO.
- Put DTOs in `internal/domain/models`.

---

### 3.4 Constructor Pattern

**Purpose:** Centralize object creation and wire a feature's layers together.

**Structure:** Every layer has a `New*` constructor that takes its dependencies.
`internal/app/routers/index.go` is the single place that assembles
repository → service → controller.

```go
// File: internal/app/routers/index.go
userRepo := repositories.NewUserRepository()          // data layer
refreshTokenRepo := repositories.NewRefreshTokenRepository()
authService := services.NewAuthService(userRepo, refreshTokenRepo) // repo → service
RegisterAuthRoutes(apiV1, authService)                // service → controller

// File: internal/app/routers/auth_routes.go
func RegisterAuthRoutes(group *gin.RouterGroup, authService auth.AuthServicer) {
    authController := controllers.NewAuthController(authService)

    authRoutes := group.Group("/auth")
    authRoutes.POST("/login", authController.Login)
}
```

**✅ DO:**
- Use `New*` functions for all constructors; pass dependencies in.
- Assemble the chain exactly once, in `internal/app/routers/index.go`.
- Have `New*Repository()` return the **interface**; services depend on it, not the struct.

**❌ DON'T:**
- Build objects with `&Struct{}` directly in business logic.
- Use `init()` functions for dependency initialization.
- Create global singletons for feature state.

---

### 3.5 Middleware Pattern

**Purpose:** Handle cross-cutting concerns (CORS, request IDs, logging, metrics, rate limiting, auth).

**Two kinds of middleware:**

1. **Generic middleware lives in `internal/app/middlewares/`** and is applied globally in
   `internal/app/routers/router.go`:
   - `middlewares.CORSMiddleware()`
   - `middlewares.RequestIDMiddleware()` — flows a request ID into the context
   - `middlewares.RequestLogMiddleware()`
   - `middlewares.MetricsMiddleware()`
   - `middlewares.RateLimitMiddleware()` / `middlewares.RateLimitMiddlewareWithConfig(rps, burst)`
     — per **client IP**, applied to all of `/api/v1`.

2. **The JWT auth guard lives in `internal/app/middlewares/auth.go`** like the rest, but it
   takes a dependency: it is constructed with the `auth.AuthServicer` interface, so it depends
   on the contract rather than the concrete service.

```go
// internal/app/middlewares/auth.go — validates the Bearer JWT, sets "user_id" in the context
func AuthMiddleware(authService auth.AuthServicer) gin.HandlerFunc { /* ... */ }

// Protected routes, mounted in internal/app/routers/index.go:
authController := controllers.NewAuthController(authService)
protectedRoutes := apiV1.Group("")
protectedRoutes.Use(middlewares.AuthMiddleware(authService))
protectedRoutes.GET("/profile", authController.Profile)
```

```go
// Generic middleware is mounted globally in internal/app/routers/router.go:
r.Use(gin.Recovery())
r.Use(middleware.CORSMiddleware())
r.Use(middleware.RequestIDMiddleware())
r.Use(middleware.RequestLogMiddleware())
r.Use(middleware.MetricsMiddleware())

v1 := r.Group("/api/v1")
v1.Use(middleware.RateLimitMiddleware()) // per-IP rate limit on all business routes
```

> **Note:** trusted-proxy handling is configuration, not middleware. `server.go` calls
> `r.SetTrustedProxies(cfg.Server.TrustedProxies)` (the `TRUSTED_PROXIES` env). Empty = trust none,
> so `c.ClientIP()` uses the real peer address and `X-Forwarded-For` cannot be spoofed.

---

### 3.6 Dependency Injection Pattern

**Purpose:** Decouple components and make every layer testable in isolation.

The cornerstone is the **repository interface**: the data layer publishes an exported interface,
and the service depends on that rather than on a concrete type. Wiring happens in
`internal/app/routers/index.go`.

```go
// 1. The repository publishes the interface it satisfies (user_repo.go):
type UserRepository interface {
    List() ([]*models.User, error)
    GetByID(id uint) (*models.User, error)
}

type userRepo struct{}

func NewUserRepository() UserRepository { return &userRepo{} }

// 2. The service depends on that interface (auth_service.go):
type AuthService struct {
    userRepo  repositories.UserRepository
    refreshTokenRepo repositories.RefreshTokenRepository
}

func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository) *AuthService {
    return &AuthService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

// 3. The controller holds the service (auth_controller.go):
type AuthController struct{ service *services.AuthService }

func NewAuthController(service *services.AuthService) *AuthController {
    return &AuthController{service: service}
}

// 4. index.go wires the concrete chain (repo → service → controller):
authService := services.NewAuthService(repositories.NewUserRepository(), repositories.NewRefreshTokenRepository())
RegisterAuthRoutes(apiV1, authService)
```

> The auth service shows the variant where the *service itself* is consumed through an interface:
> it publishes `AuthServicer` (`internal/app/services/auth/interface.go`) so the auth middleware
> can depend on the contract rather than on the concrete `*AuthService`.

**✅ Benefits:**
- Easy to test (inject a fake) — see `tests/unit/services/auth_service_test.go`.
- The interface lists exactly what callers use, nothing more.
- No hidden state; dependencies are explicit in the constructor.
- Swapping the data layer never touches business logic.

---

## 4. DIRECTORY STRUCTURE STANDARD

### 4.1 Complete Project Structure

> Mirrors `MODULE_GUIDE.md` — that file is the source of truth.

```
gin-boilerplate/
│
├── main.go                          # config → db → migrate → seed (dev) → serve → shutdown
│
├── internal/                        # Private application code
│   │
│   ├── adapters/
│   │   └── database/
│   │       ├── database.go          #   DbConnection(master, replica), GetDB(), package-level DB
│   │       ├── migrations/
│   │       │   ├── migration.go     #     golang-migrate runner, fatal on failure
│   │       │   └── sql/             #     versioned NNNNNN_name.up.sql / .down.sql pairs
│   │       └── seeders/             #   idempotent demo data, development only
│   │
│   ├── app/
│   │   ├── controllers/             #   HTTP layer — <name>_controller.go
│   │   │   ├── auth_controller.go
│   │   │   ├── auth_controller.go   #     ← reference slice
│   │   │   └── health_controller.go
│   │   ├── dto/                     #   request/response types — <name>_dto.go
│   │   ├── middlewares/             #   auth, cors, metrics, rate_limit, request_id, request_log
│   │   ├── routers/
│   │   │   ├── router.go            #     gin engine + global middleware
│   │   │   ├── index.go             #     RegisterRoutes(): builds repos + services, mounts groups
│   │   │   ├── auth_routes.go       #     Register<Name>Routes(group, service), one per feature
│   │   │   ├── example_routes.go
│   │   │   ├── health_routes.go
│   │   │   └── swagger.go           #     OpenAPI/Swagger UI (debug only)
│   │   └── services/                #   business logic — <name>_service.go
│   │       ├── example_service.go
│   │       ├── health_service.go
│   │       └── auth/                #     a service that outgrew one file gets its own package
│   │           ├── auth_service.go  #       sentinel errors + core logic
│   │           ├── auth_service_tokens.go
│   │           ├── interface.go     #       AuthServicer, consumed by the middleware
│   │           └── mailer.go        #       EmailSender interface
│   │
│   └── domain/
│       ├── models/                  #   GORM structs with TableName() — <name>_model.go
│       └── repositories/            #   interface + unexported impl + New*Repository()
│
├── pkg/                             # Public reusable kit — must NEVER import internal/
│   ├── config/                      # Configuration loader + environment helpers
│   ├── logger/                      # Logging + LogStart/LogFinish tracing helpers
│   ├── metrics/                     # Request counters, uptime
│   ├── types/                       # SuccessResponse, ErrorResponse, APIError
│   └── utils/                       # Response helpers (Ok, BadRequest, RespondWithAPIError, …)
│
├── api/                             # Embedded OpenAPI spec (served by swagger.go)
├── docs/                            # Documentation
│   ├── MODULE_GUIDE.md              # ← source of truth for the layout
│   ├── DESIGN_PATTERNS.md           # This file
│   ├── CODING_STANDARDS.md
│   ├── AI_AGENT_RULES.md
│   └── ...
│
├── tests/                           # All tests live here
│   ├── unit/                        # package <layer>_test, no database
│   │   ├── controllers/  services/  middlewares/
│   ├── integration/                 # needs a real database
│   │   ├── api/  database/
│   ├── mocks/                       # shared in-memory fakes
│   └── fixtures/                    # test data
│
├── scripts/                         # Dev/ops scripts
├── .docker/                         # Dockerfiles + compose (dev and prod)
├── go.mod                           # Module path: github.com/0xdiaz/gin-boilerplate
├── go.sum
├── Makefile                         # `make test` → unit suite
└── README.md
```

> **Module path note:** this service uses
> `github.com/0xdiaz/gin-boilerplate`. Each stamped service should set its own path in `go.mod`
> (e.g. `github.com/your-org/your-service`) and update imports — a one-shot find/replace verified
> with `go build ./...`.

### 4.2 Package Organization Rules

**✅ RULES:**

1. **One layer = one Go package = one folder.**
   ```
   ✅ internal/app/controllers/   → package controllers  (auth_controller.go, health_controller.go)
   ✅ internal/app/services/auth/ → package auth         (a service that needs several files)
   ❌ internal/app/controllers/auth/ → subfolder per feature inside a layer
   ```

2. **A feature is a file name, not a folder.** `auth_controller.go`, `auth_service.go`,
   `user_repo.go`, `user_model.go` — the prefix is what ties the slice together.

3. **`pkg/` is the shared kit.** It must never import from `internal/`. Anything you want identical
   across sibling services goes here.

4. **`internal/app/routers/index.go` is the only place that knows concrete types.** Adding a
   feature means building its repository and service there, then calling its
   `Register<Name>Routes(...)`.

5. **Package naming:**
   - Lowercase, single word, no underscores or dashes.
   - Layer packages are plural: `controllers`, `services`, `repositories`, `models`, `middlewares`.
   - A service with its own package uses the singular feature name: `auth`.
   - Utility packages are singular: `logger`, `config`, `utils`, `metrics`.

6. **Grouping imports:**
   ```go
   import (
       // Standard library
       "context"
       "time"

       // External dependencies
       "github.com/gin-gonic/gin"
       "gorm.io/gorm"

       // Internal packages
       "github.com/0xdiaz/gin-boilerplate/pkg/logger"
       "github.com/0xdiaz/gin-boilerplate/pkg/utils"
   )
   ```

---

## 5. LAYER RESPONSIBILITIES

### 5.1 Controller Layer (Thin Layer)

**Responsibility:** Handle HTTP concerns ONLY.

**What handlers SHOULD do:**
```go
// File: internal/app/controllers/auth_controller.go
func (ctrl *AuthController) List(c *gin.Context) {
    // 1. Open a trace span from the request context
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.List")

    // 2. Call the service
    data, err := ctrl.service.List(ctx)
    if err != nil {
        logger.LogFinish(ctx, "AuthController.List", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }

    // 3. Respond with a response utility
    logger.LogFinish(ctx, "AuthController.List", nil, start)
    utils.Ok(c, data, "Data retrieved successfully")
}
```

**✅ Handlers SHOULD:**
- Bind request data (`c.ShouldBindJSON`, params, query).
- Open/close a trace span (`logger.LogStart` / `logger.LogFinish`).
- Call the service through the consumer-defined `service` interface.
- Map known sentinel errors to HTTP via `utils.RespondWithAPIError` (see §9.2).
- Use response utilities (`utils.Ok`, `utils.BadRequest`, `utils.Created`, …).

**❌ Handlers MUST NOT:**
- Contain business logic.
- Access the database or call repositories directly.
- Write `c.JSON(...)` by hand for API responses.
- Import the service or controller layer.

**Size Limits:**
- **Maximum 50 lines** per handler method.
- **Maximum 300 lines** per `handler.go` file (split helpers out if larger).

---

### 5.2 Service Layer (Fat Layer)

**Responsibility:** Implement ALL business logic.

**What services SHOULD do:**
```go
// File: internal/app/services/auth/auth_service.go (abridged)
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    ctx, start := logger.LogStart(ctx, "AuthService.Register")

    // 1. Business constraint: email must be unique
    existing, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        logger.LogFinish(ctx, "AuthService.Register", err, start)
        return nil, fmt.Errorf("failed to check email: %w", err)
    }
    if existing != nil {
        logger.LogFinish(ctx, "AuthService.Register", ErrEmailAlreadyExists, start)
        return nil, ErrEmailAlreadyExists // sentinel error
    }

    // 2. Apply business logic (hash password, build model)
    hashed, err := s.hashPassword(req.Password)
    if err != nil {
        logger.LogFinish(ctx, "AuthService.Register", err, start)
        return nil, fmt.Errorf("failed to process password: %w", err)
    }
    user := &User{Name: req.Name, Email: req.Email, Password: hashed}

    // 3. Orchestrate repository calls
    if err = s.userRepo.CreateUser(user); err != nil {
        logger.LogFinish(ctx, "AuthService.Register", err, start)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 4. Build response, log the business event
    logger.Infof("user registered successfully: %s", user.Email)
    logger.LogFinish(ctx, "AuthService.Register", nil, start)
    return &AuthResponse{ /* ... */ }, nil
}
```

**✅ Services SHOULD:**
- Implement ALL business logic and rules.
- Accept and propagate `context.Context`.
- Wrap operations in `logger.LogStart` / `logger.LogFinish`.
- Orchestrate multiple repository calls; call external services through injected interfaces.
- Return service sentinel errors for known conditions, wrap unexpected errors with `%w`.

**❌ Services MUST NOT:**
- Handle HTTP status codes or write responses.
- Access the database directly (use the repository).
- Import the handler layer.

**Size Limits:**
- **Maximum 100 lines** per service method.
- **Maximum 400 lines** per service file; split into `service.go`, `service_<topic>.go`
  (the auth service splits token logic into `auth_service_tokens.go`).

---

### 5.3 Repository Layer (Data Layer)

**Responsibility:** Handle data persistence ONLY, against the **injected** `*gorm.DB`.

**What repositories SHOULD do:**
```go
// File: internal/domain/repositories/user_repo.go (abridged)
type repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
    return &repository{db: db}
}

func (r *repository) CreateUser(user *User) error {
    if err := r.db.Create(user).Error; err != nil {
        logger.Errorf("failed to create user: %v", err)
        return fmt.Errorf("failed to create user: %w", err)
    }
    return nil
}

func (r *repository) GetUserByEmail(email string) (*User, error) {
    var user User
    err := r.db.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // not found is not an error here
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    return &user, nil
}
```

> Every repository follows the same shape: an **exported** `<Name>Repository` interface, an
> unexported struct implementing it, and a `New<Name>Repository()` constructor returning the
> interface. That is what lets `tests/mocks/` substitute a fake without touching the service.

**✅ Repositories SHOULD:**
- Query through the package-level `database.DB` handle — repositories are the only layer that may.
- Perform CRUD and build queries with GORM/parameterized SQL.
- Translate database-specific errors (`gorm.ErrRecordNotFound` → `nil, nil`).
- Use transactions for multi-step operations.
- Optimize queries (`Preload`, `Select`, paging).
- Return models from `internal/domain/models`.
- Log and wrap every error with `%w`.

**❌ Repositories MUST NOT:**
- Contain business logic or validate business rules.
- Return `gorm.ErrRecordNotFound` to the caller — return `(nil, nil)`.
- Import the service or controller layer.

**Size Limits:**
- **Maximum 30 lines** per repository method.
- **Maximum 300 lines** per repository file.

---

## 6. IMPLEMENTATION PATTERNS

### 6.1 Controller Pattern

**✅ REQUIRED Pattern (struct + consumer-defined `service` interface):**
```go
// File: internal/app/controllers/auth_controller.go
package example

import (
    "context"

    "github.com/0xdiaz/gin-boilerplate/pkg/logger"
    "github.com/0xdiaz/gin-boilerplate/pkg/utils"
    "github.com/gin-gonic/gin"
)

// service is the business contract this handler needs (consumer-defined for testability).
type service interface {
    List(ctx context.Context) ([]*Example, error)
}

// AuthController is the HTTP layer for authentication.
type AuthController struct {
    svc service
}

// NewAuthController creates an AuthController backed by the given service.
func NewAuthController(service *services.AuthService) *AuthController {
    return &AuthController{service: service}
}

// List handles GET /examples.
func (ctrl *AuthController) List(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.List")
    data, err := ctrl.service.List(ctx)
    if err != nil {
        logger.LogFinish(ctx, "AuthController.List", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }
    logger.LogFinish(ctx, "AuthController.List", nil, start)
    utils.Ok(c, data, "Data retrieved successfully")
}
```

**❌ WRONG Pattern (standalone function on a global):**
```go
// DON'T DO THIS
func ListExamples(c *gin.Context) {
    var list []Example
    database.DB.Find(&list) // ❌ global DB + business logic in the HTTP layer
    c.JSON(200, list)       // ❌ hand-written response
}
```

---

### 6.2 Service Pattern

**✅ REQUIRED Pattern (struct + consumer-defined `repository` interface):**
```go
// File: internal/app/services/auth/auth_service.go
package example

import (
    "context"

    "github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

// repository is the data-access contract this service needs.
type repository interface {
    List() ([]*Example, error)
}

type Service struct {
    repo repository
}

func NewService(repo repository) *Service {
    return &Service{repo: repo}
}

func (s *Service) List(ctx context.Context) ([]*Example, error) {
    ctx, start := logger.LogStart(ctx, "AuthService.List")
    list, err := s.repo.List()
    logger.LogFinish(ctx, "AuthService.List", err, start)
    return list, err
}
```

---

### 6.3 Repository Pattern

**✅ REQUIRED Pattern (struct holding the injected `*gorm.DB`):**
```go
// File: internal/domain/repositories/user_repo.go
package example

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
    return &Repository{db: db}
}

func (r *Repository) List() ([]*Example, error) {
    var list []*Example
    if err := r.db.Find(&list).Error; err != nil {
        return nil, err
    }
    return list, nil
}
```

**❌ WRONG Pattern (package functions on the global `database.DB`):**
```go
// DON'T DO THIS — globals make the repository untestable and hide dependencies.
func ListExamples() ([]*Example, error) {
    var list []*Example
    return list, database.DB.Find(&list).Error // ❌ global state
}
```

---

### 6.4 Route Wiring Pattern

Every feature gets one `internal/app/routers/<name>_routes.go` that turns a service into
mounted routes:

```go
// File: internal/app/routers/auth_routes.go
package routers

import (
    "github.com/0xdiaz/gin-boilerplate/internal/app/controllers"
    "github.com/0xdiaz/gin-boilerplate/internal/app/services"
    "github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes under the given group
// (e.g. /api/v1). Full paths are groupPrefix/auth/login, /auth/register, and so on.
func RegisterAuthRoutes(group *gin.RouterGroup, authService auth.AuthServicer) {
    authController := controllers.NewAuthController(authService)

    authRoutes := group.Group("/auth")
    {
        authRoutes.POST("/register", authController.Register)
        authRoutes.POST("/login", authController.Login)
    }
}
```

**Assembling everything** — `internal/app/routers/index.go` is the single place that constructs
concrete types:

```go
// File: internal/app/routers/index.go
func RegisterRoutes(route *gin.Engine) {
    route.NoRoute(func(ctx *gin.Context) {
        utils.NotFound(ctx, nil, "Route not found")
    })

    RegisterHealthRoutes(route) // /health and /metrics at the ROOT

    apiV1 := route.Group("/api/v1")
    apiV1.Use(middlewares.RateLimitMiddleware())

    userRepo := repositories.NewUserRepository()
    refreshTokenRepo := repositories.NewRefreshTokenRepository()
    authService := auth.NewAuthService(userRepo, refreshTokenRepo, nil)

    userRepo := repositories.NewUserRepository()
    refreshTokenRepo := repositories.NewRefreshTokenRepository()
    authService := services.NewAuthService(userRepo, refreshTokenRepo)

    RegisterAuthRoutes(apiV1, authService)
    RegisterAuthRoutes(apiV1, authService)

    // Protected routes sit behind the JWT guard.
    authController := controllers.NewAuthController(authService)
    protectedRoutes := apiV1.Group("")
    protectedRoutes.Use(middlewares.AuthMiddleware(authService))
    {
        protectedRoutes.GET("/profile", authController.Profile)
        protectedRoutes.POST("/logout-all", authController.LogoutAll)
    }
}
```

**Engine setup** happens in `internal/app/routers/router.go`: it creates the gin engine, applies
the global middleware (request ID, request log, CORS, metrics), and calls `RegisterRoutes`.

Business routes live under `/api/v1` with rate limiting; `/health` and `/metrics` mount at the
**root** via `RegisterHealthRoutes(route)`.

Migrations are **not** derived from the code: they are versioned SQL files applied by
`migrations.Migrate()` from `main.go`, and a failure is fatal.

---

### 6.5 Response Utility Pattern

**✅ MANDATORY: Always use the response utilities in `pkg/utils`.**

```go
import "github.com/0xdiaz/gin-boilerplate/pkg/utils"

func (ctrl *AuthController) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data") // 400 + field errors
        return
    }

    data, err := ctrl.service.Create(c.Request.Context(), &req)
    if err != nil {
        if apiErr := errToAPIError(err); apiErr != nil {
            utils.RespondWithAPIError(c, apiErr) // mapped sentinel error
            return
        }
        utils.InternalServerError(c, err, "Failed to create")
        return
    }

    utils.Created(c, data, "Created successfully") // 201
}
```

**Available response functions (`pkg/utils/response.go`):**
```go
// Success
utils.Ok(c, data, message)                          // 200 OK
utils.Created(c, data, message)                     // 201 Created
utils.NoContent(c)                                  // 204 No Content

// Errors
utils.BadRequest(c, err, message)                   // 400  (formats validation field errors)
utils.Unauthorized(c, err, message)                 // 401
utils.Forbidden(c, err, message)                    // 403
utils.NotFound(c, err, message)                     // 404
utils.Conflict(c, err, message)                     // 409
utils.UnprocessableEntity(c, err, message)          // 422
utils.TooManyRequests(c, err, message)              // 429
utils.InternalServerError(c, err, message)          // 500
utils.BadGateway(c, err, message)                   // 502
utils.ServiceUnavailable(c, err, message)           // 503
utils.ServiceUnavailableWithData(c, data, message)  // 503 with a body (health uses this)

// From a mapped domain error
utils.RespondWithAPIError(c, apiErr)                // uses apiErr.Code / .Message / .Details

// Generic
utils.HandleSuccess(c, statusCode, data, message)
utils.HandleErrors(c, statusCode, err, message)
```

**Response Format (`pkg/types` — `SuccessResponse` / `ErrorResponse`):**
```json
{
    "success": true,
    "message": "Data retrieved successfully",
    "data": { "...": "..." },
    "errors": null
}
```

---

## 7. REQUEST FLOW PATTERNS

### 7.1 Standard CRUD Flow

**Complete request flow for a LIST/CREATE operation (the `auth` slice):**

```
1. HTTP GET /api/v1/examples
       ↓
2. [Global Middleware] (internal/app/routers/router.go)
   • gin.Recovery
   • CORSMiddleware
   • RequestIDMiddleware   (request id → context)
   • RequestLogMiddleware
   • MetricsMiddleware
       ↓
3. [/api/v1 group]
   • RateLimitMiddleware   (per client IP)
       ↓
4. [Feature route]  RegisterAuthRoutes → authController.List
       ↓
5. [Controller] internal/app/controllers/auth_controller.go
   func (ctrl *AuthController) List(c *gin.Context) {
       ctx, start := logger.LogStart(c.Request.Context(), "AuthController.List")
       data, err := ctrl.service.List(ctx)            // call the injected service
       if err != nil { utils.InternalServerError(c, err, "Failed to retrieve data"); return }
       logger.LogFinish(ctx, "AuthController.List", nil, start)
       utils.Ok(c, data, "Data retrieved successfully")
   }
       ↓
6. [Service] internal/app/services/auth/auth_service.go
   func (s *Service) List(ctx context.Context) ([]*Example, error) {
       ctx, start := logger.LogStart(ctx, "AuthService.List")
       list, err := s.repo.List()              // call repo via `repository` interface
       logger.LogFinish(ctx, "AuthService.List", err, start)
       return list, err
   }
       ↓
7. [Repository] internal/domain/repositories/user_repo.go
   func (r *Repository) List() ([]*Example, error) {
       var list []*Example
       return list, r.db.Find(&list).Error     // injected *gorm.DB
   }
       ↓
8. [Database] SELECT * FROM examples
       ↓
9. Response back up the stack → HTTP 200
   { "success": true, "message": "Data retrieved successfully", "data": [...], "errors": null }
```

For CREATE: the handler binds the DTO (`c.ShouldBindJSON`), the service validates business rules
and builds the model, the repository persists it, and the handler responds with `utils.Created`.

---

### 7.2 Authentication Flow Pattern

Auth uses **JWT directly** (no OTP step). Routes are mounted by
`RegisterAuthRoutes(apiV1, authService)` under `/api/v1`.

```
PUBLIC routes (no auth):
  POST /api/v1/auth/register         → AuthController.Register
  POST /api/v1/auth/login            → AuthController.Login
  POST /api/v1/auth/refresh          → AuthController.RefreshToken
  POST /api/v1/auth/forgot-password  → AuthController.ForgotPassword
  POST /api/v1/auth/reset-password   → AuthController.ResetPassword

PROTECTED route (middlewares.AuthMiddleware(authService) guard):
  GET  /api/v1/profile               → AuthController.Profile
```

**Login flow:**
```
1. POST /api/v1/auth/login   {"email": "...", "password": "..."}
       ↓
2. [Controller.Login] bind LoginRequest → service.Login(ctx, &req)
       ↓
3. [Service.Login]
   • userRepo.GetUserByEmail(email)
   • bcrypt.CompareHashAndPassword(...)        (invalid → ErrInvalidCredentials)
   • generateToken(user)                        (signed JWT, user_id + email, 15m default)
   • generateRefreshToken()                     (32 random bytes → hex)
   • userRepo.UpdateUser(user)                  (persist refresh token)
   • return *AuthResponse{User, AccessToken, RefreshToken, TokenType:"Bearer"}
       ↓
4. [Controller.Login] on error: errToAPIError(err) → utils.RespondWithAPIError
                   on success: utils.Ok(c, response, "Login successful")
       ↓
5. Response: { "success": true, "data": { "user": {...},
               "access_token": "eyJ...", "refresh_token": "...", "token_type": "Bearer" } }

Subsequent protected requests:
6. GET /api/v1/profile   Header: Authorization: Bearer eyJ...
       ↓
7. [middlewares.AuthMiddleware]  internal/app/middlewares/auth.go
   • split "Bearer <token>"
   • service.ValidateToken(token) → userID   (HS256, checks signing method & exp)
   • c.Set("user_id", userID); c.Next()       (401 via utils.Unauthorized on failure)
       ↓
8. [Controller.Profile]  userID := c.GetUint("user_id")
```

> The JWT secret comes from config (`config.Get().Server.JWTSecret`). The guard lives in
> `internal/app/middlewares/auth.go` and takes the `auth.AuthServicer` interface, so it depends
> on the contract rather than the concrete service.

---

### 7.3 Transaction Flow Pattern

**Multi-step operations use a GORM transaction inside the repository**, via `database.DB`:

```go
// File: internal/domain/repositories/<name>_repo.go
func (r *Repository) ProcessPayment(clientID uint, amount float64) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Step 1: create the transaction record
        txn := &Transaction{ClientID: clientID, Amount: amount, Status: "pending"}
        if err := tx.Create(txn).Error; err != nil {
            return fmt.Errorf("create transaction: %w", err)
        }

        // Step 2: deduct balance
        if err := tx.Model(&Account{}).
            Where("client_id = ?", clientID).
            Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
            return fmt.Errorf("update balance: %w", err)
        }

        // All steps succeeded — returning nil commits; any error rolls back automatically.
        return nil
    })
}
```

The service orchestrates and decides *when* to call this; the repository owns the transaction
mechanics. Keep the closure short and let the returned error drive commit/rollback.

---

## 8. DATA FLOW PATTERNS

### 8.1 Request → Response Data Transformation

```
HTTP Request (JSON)
       ↓
[Controller] c.ShouldBindJSON → DTO
       ↓
RegisterRequest {
    Name:     "Jane Doe"
    Email:    "jane@example.com"
    Password: "plaintext"
}
       ↓
[Service] transform DTO → domain model (apply business logic)
       ↓
User {
    Name:     "Jane Doe"
    Email:    "jane@example.com"
    Password: "$2a$10$hashed..."  (bcrypt)
}
       ↓
[Repository] r.db.Create(user)  → DB row
       ↓
[Service] build a response DTO (never leak secrets)
       ↓
AuthResponse {
    User:        UserResponse{ID:1, Name:"Jane Doe", Email:"jane@example.com"}
    AccessToken: "eyJ..."
    // Password NEVER included
}
       ↓
[Controller] utils.Ok / utils.Created → standard JSON envelope
{
    "success": true,
    "message": "...",
    "data": { "user": {...}, "access_token": "..." },
    "errors": null
}
```

**Key Transformations:**
1. **Controller:** JSON → DTO (binding/validation at the boundary).
2. **Service:** DTO → domain model (business logic); domain model → response DTO.
3. **Repository:** domain model ↔ database.
4. **Controller:** response DTO → standard JSON envelope via `pkg/utils`.

---

### 8.2 DTO vs Model Usage

**✅ When to use DTOs:**
- API request payloads (binding/validation).
- API response payloads (shape the client sees; strip secrets).
- External/inter-service communication.

**✅ When to use Models:**
- Internal business logic and domain rules.
- Database operations in the repository.

**Example (the auth service):**
```go
// DTO for the API boundary (dto.go)
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// Model for the domain & DB (model.go). Password is never serialized.
type User struct {
    ID        uint       `json:"id" gorm:"primaryKey"`
    Name      string     `json:"name" gorm:"type:varchar(255);not null"`
    Email     string     `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    Password  string     `json:"-" gorm:"type:varchar(255);not null"` // bcrypt, never exposed
    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// DTO for the API response — only safe fields.
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

---

## 9. ERROR HANDLING PATTERNS

### 9.1 Error Flow Pattern

```
[Repository] DB error
       ↓  translate (gorm.ErrRecordNotFound → nil,nil) or wrap with %w
    return nil, fmt.Errorf("failed to get user by email: %w", err)
       ↓
[Service] add business context; return a SENTINEL error for known cases
       ↓
    if existing != nil { return nil, ErrEmailAlreadyExists }
    return nil, fmt.Errorf("failed to check email: %w", err)
       ↓
[Controller] map sentinel → APIError, else 500
       ↓
    if apiErr := errToAPIError(err); apiErr != nil {
        utils.RespondWithAPIError(c, apiErr)   // e.g. 409 Conflict
        return
    }
    utils.InternalServerError(c, err, "Failed to register user")
       ↓
HTTP Response: 409 Conflict
{ "success": false, "message": "Email already exists", "data": null, "errors": null }
```

### 9.2 Sentinel Errors & APIError Mapping

The pattern is: **each service defines sentinel errors; the controller maps them to `types.APIError`
and responds via `utils.RespondWithAPIError`.** `pkg/types` provides the `APIError` type
(`{Code, Message, Details}`) and a few shared predefined errors.

```go
// 1. The service defines sentinel errors (internal/app/services/auth/auth_service.go):
var (
    ErrEmailAlreadyExists  = errors.New("email already exists")
    ErrInvalidCredentials  = errors.New("invalid email or password")
    ErrUserNotFound        = errors.New("user not found")
    ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
    // ...
)

// 2. The handler maps them to types.APIError (internal/app/controllers/auth_controller.go):
func errToAPIError(err error) *types.APIError {
    switch {
    case errors.Is(err, ErrEmailAlreadyExists):
        return &types.APIError{Code: http.StatusConflict, Message: "Email already exists"}
    case errors.Is(err, ErrInvalidCredentials):
        return &types.APIError{Code: http.StatusUnauthorized, Message: "Invalid email or password"}
    case errors.Is(err, ErrUserNotFound):
        return &types.APIError{Code: http.StatusNotFound, Message: "User not found"}
    default:
        return nil // unknown → caller falls back to 500
    }
}

// 3. The handler uses it:
response, err := h.service.Login(ctx, &req)
if err != nil {
    if apiErr := errToAPIError(err); apiErr != nil {
        utils.RespondWithAPIError(c, apiErr)
        return
    }
    utils.InternalServerError(c, err, "Failed to authenticate user")
    return
}
```

`pkg/types` also ships ready-made `APIError` values (`types.ErrNotFound`, `types.ErrUnauthorized`,
`types.ErrConflict`, `types.ErrRateLimitExceeded`, …) for cases where a service does not need its
own message. All HTTP error responses go through `pkg/utils` — do not duplicate status-code logic
elsewhere.

### 9.3 Error Wrapping Pattern

**✅ ALWAYS wrap unexpected errors with context (`%w`) so the chain stays inspectable:**

```go
// Repository: what failed + which resource
func (r *repository) GetUserByEmail(email string) (*User, error) {
    var user User
    if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    return &user, nil
}

// Service: add business context, or return a sentinel for known cases
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    existing, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        return nil, fmt.Errorf("failed to check email: %w", err)
    }
    if existing != nil {
        return nil, ErrEmailAlreadyExists // sentinel — handler maps to 409
    }
    // ...
}
```

**Error chain example:**
```
Original:    record not found
Repository:  failed to get user by email: record not found
Service:     failed to check email: failed to get user by email: record not found
Controller:  HTTP 500 (unknown) — or a mapped APIError for a known sentinel
```

---

## 10. TESTING PATTERNS

> **Unit tests live in `tests/unit/<layer>/`**, in package `<layer>_test`, driving the code
> through its exported surface. Services are tested against the shared in-memory fakes in
> `tests/mocks/` — no database. Controller tests use `httptest` against the real controller with
> a service built on those fakes. Tests needing a real database go in `tests/integration/`.
> `make test` runs the unit suite.

### 10.1 Service Layer Testing Pattern

The canonical example is `tests/mocks/user_repo_mock.go` plus
`tests/unit/services/auth_service_test.go`: the fake satisfies the exported repository
interface, so no database is needed.

```go
// File: tests/mocks/user_repo_mock.go
package mocks

// MockUserRepository is an in-memory UserRepository for unit tests.
type MockUserRepository struct {
    mu     sync.RWMutex
    nextID uint
    byID   map[uint]*models.User

    // ListErr and GetErr, when set, are returned instead of data.
    ListErr error
    GetErr  error
}

// Compile-time proof the fake still matches the real interface.
var _ repositories.UserRepository = (*MockUserRepository)(nil)

// AddUserByEmail inserts a user directly, for setting up "already existing" test state.
func (m *MockUserRepository) AddUserByEmail(email string, user *models.User) { /* ... */ }

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
    if m.GetErr != nil {
        return nil, m.GetErr
    }
    user, ok := m.byID[id]
    if !ok {
        return nil, nil // matches the real repo: a missing row is (nil, nil)
    }
    return user, nil
}
```

```go
// File: tests/unit/services/auth_service_test.go
package services_test

func TestAuthServiceGet(t *testing.T) {
    userRepo := mocks.NewMockUserRepository()
    refreshTokenRepo := mocks.NewMockRefreshTokenRepository()
    userRepo.AddUserByEmail("user@example.com", &models.User{ID: 1, Email: "user@example.com", Name: "User"})

    service := services.NewAuthService(userRepo, refreshTokenRepo)

    got, err := service.GetProfile(context.Background(), 1)

    require.NoError(t, err)
    assert.Equal(t, "user@example.com", got.Email)
}

// A repository failure must surface as a failure, not as "not found".
func TestAuthServiceGet_RepoError(t *testing.T) {
    userRepo := mocks.NewMockUserRepository()
    userRepo.GetErr = errors.New("database is down")

    service := services.NewAuthService(userRepo, mocks.NewMockRefreshTokenRepository())

    got, err := service.Get(context.Background(), 1)

    assert.Nil(t, got)
    require.Error(t, err)
    assert.False(t, errors.Is(err, services.ErrUserNotFound))
}
```

Because the fake is shared in `tests/mocks/`, every service test reuses it. To test an error
path, set the fake's error field and assert it propagates.

### 10.2 Table-Driven Testing Pattern

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {name: "valid email", email: "user@example.com", wantErr: false},
        {name: "missing @", email: "userexample.com", wantErr: true},
        {name: "empty email", email: "", wantErr: true},
        {name: "missing domain", email: "user@", wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if (err != nil) != tt.wantErr {
                t.Errorf("ValidateEmail() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

---

## 11. COMPLETE FEATURE IMPLEMENTATION GUIDE

### Step-by-Step: Adding a New Feature

> This follows `MODULE_GUIDE.md`'s recipe: **copy the `auth` slice**. Every step below mirrors
> a real file in it.

#### Step 1: Write the migration

```sql
-- File: internal/adapters/database/migrations/sql/000007_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,
    -- Money is an integer in the smallest currency unit. Never FLOAT/DOUBLE.
    price_cents BIGINT NOT NULL CHECK (price_cents >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```sql
-- File: internal/adapters/database/migrations/sql/000007_create_products_table.down.sql
DROP TABLE IF EXISTS products;
```

Every `.up.sql` needs its `.down.sql`. Never edit a migration that has already been applied.

#### Step 2: Define the model

```go
// File: internal/domain/models/product_model.go
package models

import "time"

// Product is a sellable product.
type Product struct {
    ID   uint   `json:"id" gorm:"primaryKey"`
    Name string `json:"name" gorm:"type:varchar(200);not null"`

    // PriceCents is the price in the smallest currency unit.
    PriceCents int64 `json:"price_cents" gorm:"not null"`

    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the database table name for Product model.
func (p *Product) TableName() string { return "products" }
```

Keep the struct tags in sync with the migration.

#### Step 3: Define DTOs

```go
// File: internal/app/dto/product_dto.go
package dto

// CreateProductRequest is the body of POST /api/v1/products.
type CreateProductRequest struct {
    Name       string `json:"name" binding:"required,max=200"`
    PriceCents int64  `json:"price_cents" binding:"required,min=0"`
}

// ProductResponse is the API representation of a product.
type ProductResponse struct {
    ID         uint   `json:"id"`
    Name       string `json:"name"`
    PriceCents int64  `json:"price_cents"`
}
```

#### Step 4: Implement the repository

```go
// File: internal/domain/repositories/product_repo.go
package repositories

// ProductRepository defines data access for products.
type ProductRepository interface {
    Create(product *models.Product) error
    // GetByID returns the product, or (nil, nil) when it does not exist.
    GetByID(id uint) (*models.Product, error)
}

type productRepo struct{}

// NewProductRepository returns a new ProductRepository implementation.
func NewProductRepository() ProductRepository {
    return &productRepo{}
}

func (r *productRepo) Create(product *models.Product) error {
    if err := database.DB.Create(product).Error; err != nil {
        logger.Errorf("failed to create product: %v", err)
        return fmt.Errorf("failed to create product: %w", err)
    }
    return nil
}

func (r *productRepo) GetByID(id uint) (*models.Product, error) {
    var product models.Product
    err := database.DB.Where("id = ?", id).First(&product).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        logger.Errorf("failed to get product by id: %v", err)
        return nil, fmt.Errorf("failed to get product by id: %w", err)
    }
    return &product, nil
}
```

Repositories are the only layer allowed to touch `database.DB`.

#### Step 5: Implement the service

```go
// File: internal/app/services/product_service.go
package services

// ErrProductNotFound is returned when the requested product does not exist.
var ErrProductNotFound = errors.New("product not found")

// ProductService handles product business logic.
type ProductService struct {
    productRepo repositories.ProductRepository
}

// NewProductService creates a new ProductService instance.
func NewProductService(productRepo repositories.ProductRepository) *ProductService {
    return &ProductService{productRepo: productRepo}
}

// Get returns a single product.
//
// Returns ErrProductNotFound when no product has the given id.
func (s *ProductService) Get(ctx context.Context, id uint) (response *dto.ProductResponse, err error) {
    ctx, start := logger.LogStart(ctx, "ProductService.Get")

    product, err := s.productRepo.GetByID(id)
    if err != nil {
        logger.LogFinish(ctx, "ProductService.Get", err, start)
        return nil, fmt.Errorf("failed to get product: %w", err)
    }
    if product == nil {
        logger.LogFinish(ctx, "ProductService.Get", ErrProductNotFound, start)
        return nil, ErrProductNotFound
    }

    logger.LogFinish(ctx, "ProductService.Get", nil, start)
    return &dto.ProductResponse{ID: product.ID, Name: product.Name, PriceCents: product.PriceCents}, nil
}
```

The service depends on the repository **interface**, so it can be tested without a database.

#### Step 6: Implement the controller

```go
// File: internal/app/controllers/product_controller.go
package controllers

// ProductController handles product endpoints.
type ProductController struct {
    service *services.ProductService
}

// NewProductController creates a new ProductController instance.
func NewProductController(service *services.ProductService) *ProductController {
    return &ProductController{service: service}
}

// Get returns a single product.
//
// GET /api/v1/products/:id
func (ctrl *ProductController) Get(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "ProductController.Get")

    id, err := strconv.ParseUint(c.Param("id"), 10, 64)
    if err != nil {
        logger.LogFinish(ctx, "ProductController.Get", err, start)
        utils.BadRequest(c, err, "Invalid product id")
        return
    }

    product, err := ctrl.service.Get(ctx, uint(id))
    if err != nil {
        if errors.Is(err, services.ErrProductNotFound) {
            logger.LogFinish(ctx, "ProductController.Get", err, start)
            utils.NotFound(c, err, "Product not found")
            return
        }
        logger.Errorf("failed to get product: %v", err)
        logger.LogFinish(ctx, "ProductController.Get", err, start)
        utils.InternalServerError(c, err, "Failed to get product")
        return
    }

    logger.LogFinish(ctx, "ProductController.Get", nil, start)
    utils.Ok(c, product, "Product retrieved successfully")
}
```

#### Step 7: Declare the routes

```go
// File: internal/app/routers/product_routes.go
package routers

// RegisterProductRoutes registers product routes under the given group.
func RegisterProductRoutes(group *gin.RouterGroup, productService *services.ProductService) {
    productController := controllers.NewProductController(productService)
    group.GET("/products/:id", productController.Get)
}
```

#### Step 8: Wire it in the router

```go
// File: internal/app/routers/index.go
func RegisterRoutes(route *gin.Engine) {
    // ... existing wiring ...

    productRepo := repositories.NewProductRepository()
    productService := services.NewProductService(productRepo)

    RegisterProductRoutes(apiV1, productService)
}
```

That is the only place concrete types are constructed.

#### Step 9: Add the fake and the test

```go
// File: tests/mocks/product_repo_mock.go
package mocks

// MockProductRepository is an in-memory ProductRepository for unit tests.
type MockProductRepository struct {
    mu     sync.RWMutex
    nextID uint
    byID   map[uint]*models.Product

    // GetErr, when set, is returned instead of data.
    GetErr error
}

// Compile-time proof the fake still matches the real interface.
var _ repositories.ProductRepository = (*MockProductRepository)(nil)
```

```go
// File: tests/unit/services/product_service_test.go
package services_test

func TestProductServiceGet(t *testing.T) {
    productRepo := mocks.NewMockProductRepository()
    product := productRepo.Seed(&models.Product{Name: "Starter plan", PriceCents: 15000000})

    service := services.NewProductService(productRepo)

    got, err := service.Get(context.Background(), product.ID)

    require.NoError(t, err)
    assert.Equal(t, "Starter plan", got.Name)
    assert.Equal(t, int64(15000000), got.PriceCents)
}

func TestProductServiceGet_NotFound(t *testing.T) {
    service := services.NewProductService(mocks.NewMockProductRepository())

    got, err := service.Get(context.Background(), 999)

    assert.Nil(t, got)
    assert.True(t, errors.Is(err, services.ErrProductNotFound))
}
```

#### Step 10: Verify

```bash
gofmt -w .
go build ./...
go vet ./...
make test
```

---

## 12. PATTERN EXAMPLES FROM CODEBASE

### 12.1 Auth Pattern (from the `auth` service)

**✅ Controller maps domain errors and uses response utilities:**
```go
// internal/app/controllers/auth_controller.go
func (ctrl *AuthController) Login(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Login")

    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.LogFinish(ctx, "AuthController.Login", err, start)
        utils.BadRequest(c, err, "Invalid request data")
        return
    }

    response, err := h.service.Login(ctx, &req)
    if err != nil {
        if apiErr := errToAPIError(err); apiErr != nil {
            logger.LogFinish(ctx, "AuthController.Login", err, start)
            utils.RespondWithAPIError(c, apiErr) // e.g. 401 invalid credentials
            return
        }
        logger.LogFinish(ctx, "AuthController.Login", err, start)
        utils.InternalServerError(c, err, "Failed to authenticate user")
        return
    }

    logger.LogFinish(ctx, "AuthController.Login", nil, start)
    utils.Ok(c, response, "Login successful")
}
```

**✅ Service depends on the `Repository` interface, returns sentinel errors:**
```go
// internal/app/services/auth/auth_service.go
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
    ctx, start := logger.LogStart(ctx, "AuthService.Login")

    user, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        logger.LogFinish(ctx, "AuthService.Login", err, start)
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    if user == nil {
        logger.LogFinish(ctx, "AuthService.Login", ErrInvalidCredentials, start)
        return nil, ErrInvalidCredentials
    }
    if err = s.verifyPassword(user.Password, req.Password); err != nil {
        logger.LogFinish(ctx, "AuthService.Login", ErrInvalidCredentials, start)
        return nil, ErrInvalidCredentials
    }

    accessToken, _ := s.generateToken(user)
    refreshToken, _ := s.generateRefreshToken()
    user.RefreshToken = refreshToken
    _ = s.userRepo.UpdateUser(user)

    logger.LogFinish(ctx, "AuthService.Login", nil, start)
    return &AuthResponse{
        User:         UserResponse{ID: user.ID, Name: user.Name, Email: user.Email},
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
    }, nil
}
```

**✅ The service publishes its contract so the middleware can depend on it:**
```go
// internal/app/routers/auth_routes.go
// internal/app/services/auth/interface.go
type AuthServicer interface {
    ValidateToken(tokenString string) (*Claims, error)
    // ...
}

// internal/app/middlewares/auth.go — the guard takes the interface
func AuthMiddleware(authService auth.AuthServicer) gin.HandlerFunc { /* JWT guard */ }
```

### 12.2 DataTable Pattern (from the `example` slice)

The example slice shows server-side DataTables: the repository runs the query, the service wraps
it in a trace span, and the handler renders it.

```go
// internal/domain/repositories/user_repo.go
func (r *Repository) Datatables(c *gin.Context) (interface{}, error) {
    var rows []*Example
    return datatables.OfReturn(
        c,
        r.db.Model(&Example{}),
        &rows,
        []string{"id", "data"},
        map[string]string{"id": "id", "data": "data"},
        datatables.NewOptions().WithIndex("DT_RowIndex", false),
    )
}

// internal/app/controllers/auth_controller.go
func (ctrl *AuthController) Datatables(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Datatables")
    data, err := ctrl.service.Datatables(ctx, c)
    if err != nil {
        logger.LogFinish(ctx, "AuthController.Datatables", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }
    dt, ok := data.(dto.Datatables)
    if !ok {
        utils.InternalServerError(c, nil, "Failed to render data")
        return
    }
    logger.LogFinish(ctx, "AuthController.Datatables", nil, start)
    datatables.JSON(c, dt)
}
```

Routes are mounted in `internal/app/routers/example_routes.go`: `group.GET("/datatables", ...)`.

---

## 13. ANTI-PATTERNS TO AVOID

### 13.1 ❌ Business Logic in Controller

**WRONG:**
```go
func (ctrl *AuthController) Register(c *gin.Context) {
    var req RegisterRequest
    c.ShouldBindJSON(&req)

    // ❌ Password hashing in the handler!
    hashed, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

    // ❌ Building & persisting the model in the HTTP layer!
    user := &User{Name: req.Name, Email: req.Email, Password: string(hashed)}
    h.db.Create(user)

    c.JSON(200, user) // ❌ hand-written response
}
```

**CORRECT:**
```go
func (ctrl *AuthController) Register(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Register")
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.LogFinish(ctx, "AuthController.Register", err, start)
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    // ✅ Delegate everything to the service
    response, err := h.service.Register(ctx, &req)
    if err != nil {
        if apiErr := errToAPIError(err); apiErr != nil {
            utils.RespondWithAPIError(c, apiErr)
            return
        }
        utils.InternalServerError(c, err, "Failed to register user")
        return
    }
    logger.LogFinish(ctx, "AuthController.Register", nil, start)
    utils.Created(c, response, "User registered successfully")
}
```

---

### 13.2 ❌ Direct Database Access in Service

**WRONG:**
```go
func (s *Service) Register(ctx context.Context, req *RegisterRequest) error {
    // ❌ Direct DB access in the service — bypasses the repository, untestable, hidden dependency!
    var existing User
    database.DB.Where("email = ?", req.Email).First(&existing)

    user := &User{ /* ... */ }
    database.DB.Create(user)
    return nil
}
```

**CORRECT:**
```go
func (s *Service) Register(ctx context.Context, req *RegisterRequest) error {
    // ✅ Go through the injected repository interface
    existing, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        return fmt.Errorf("failed to check email: %w", err)
    }
    if existing != nil {
        return ErrEmailAlreadyExists
    }
    return s.userRepo.CreateUser(&User{ /* ... */ })
}
```

---

### 13.3 ❌ Standalone Handler Functions (no struct, no injection)

**WRONG:**
```go
// ❌ Standalone function, no struct, no injected dependency
func Login(c *gin.Context) {
    // direct implementation, reaches for globals
}
api.POST("/auth/login", Login)
```

**CORRECT:**
```go
// ✅ Struct controller with an injected service
type AuthController struct{ service auth.AuthServicer }
func NewAuthController(service auth.AuthServicer) *AuthController { return &AuthController{service: service} }
func (ctrl *AuthController) Login(c *gin.Context) { /* ... */ }

// Wired in internal/app/routers/auth_routes.go:
func RegisterAuthRoutes(group *gin.RouterGroup, authService auth.AuthServicer) {
    authController := controllers.NewAuthController(authService)
    group.POST("/auth/login", authController.Login)
}
```

---

### 13.4 ❌ God Service (Too Many Responsibilities)

**WRONG:**
```go
// ❌ One service doing everything
type AppService struct{}
func (s *AppService) CreateUser() error      { /* ... */ }
func (s *AppService) ProcessPayment() error  { /* ... */ }
func (s *AppService) SendEmail() error       { /* ... */ }
func (s *AppService) GenerateReport() error  { /* ... */ }
// ... 50 more methods
```

**CORRECT:**
```go
// ✅ Separate services, each with a single responsibility
internal/app/services/auth/            // authentication
internal/app/services/payment_service.go  // payment processing
internal/app/services/report_service.go   // reporting
```
A service that needs another one receives it through its `New*Service(...)` constructor (see §13.5).

---

### 13.5 ❌ Skipping a Layer

**WRONG:**
```go
// ❌ A service building its own repository, or reaching for the database directly
func (s *AuthService) doThing() {
    var users []*models.User
    database.DB.Find(&users)           // ❌ services never touch the database
    repo := repositories.NewUserRepository() // ❌ dependencies are injected, not constructed here
    _ = repo
}
```

**CORRECT:**
```go
// ✅ Depend on interfaces, injected through the constructor in routers/index.go
type AuthService struct {
    userRepo  repositories.UserRepository
    refreshTokenRepo repositories.RefreshTokenRepository
}

func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository) *AuthService {
    return &AuthService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}

// internal/app/routers/index.go — the single wiring point
userRepo := repositories.NewUserRepository()
refreshTokenRepo := repositories.NewRefreshTokenRepository()
authService := services.NewAuthService(userRepo, refreshTokenRepo)
```
Because the service only knows interfaces, the same code runs against Postgres in production
and against `tests/mocks/` in unit tests.

---

## 🎯 SUMMARY: Key Takeaways

### ✅ DO:
1. **Organize by layer** — `controllers` → `services` → `repositories` → `models`.
2. **Name the slice consistently** — `auth_controller.go`, `auth_service.go`, `user_repo.go`.
3. **Declare the interface at the repository** — services depend on it, never on the concrete type.
4. **Inject dependencies via `New*(...)`** — assemble once in `internal/app/routers/index.go`.
5. **Only repositories touch `database.DB`** — services and controllers never do.
6. **Versioned SQL migrations** — every `.up.sql` has a `.down.sql`; never edit an applied one.
7. **Use response utilities** (`pkg/utils`) — never hand-write `c.JSON`.
8. **Trace every operation** — `logger.LogStart` / `logger.LogFinish`, span name `<Type>.<Method>`.
9. **Map sentinel errors → HTTP status** in the controller; wrap unexpected errors with `%w`.
10. **Test in `tests/unit/`** — shared fakes from `tests/mocks/`, no DB.

### ❌ DON'T:
1. **Business logic in handlers** — move it to the service.
2. **Database access in services** — go through the repository.
3. **Standalone handler functions** — use struct handlers with injected interfaces.
4. **God services** — split by sub-domain.
5. **Skip a layer** — a controller must never call a repository.
6. **Use `database.DB` outside a repository** — that is what the repository is for.
7. **Money as `float64`** — use `int64` in the smallest currency unit.
8. **`pkg/` importing `internal/`** — the shared kit must stay clean.
9. **Skip tests** — every service gets a fake-backed test in `tests/unit/services/`.
10. **Large files/functions** — split per the size limits in §5.

---

**END OF DESIGN PATTERNS DOCUMENT**

*Last Updated: 2026-06-10*
*Version: 2.1 (layered architecture — aligned with the code)*
*Source of truth for layout: [`MODULE_GUIDE.md`](./MODULE_GUIDE.md). Reference slice: `auth`.*

---

## 📚 Related Documentation

- **[MODULE_GUIDE.md](./MODULE_GUIDE.md)** — Source of truth for the modular layout
- **CODING_STANDARDS.md** — Detailed coding standards and guidelines
- **AI_AGENT_RULES.md** — Quick reference rules for AI agents
- **CONTRACTS.md** — Stable API & configuration contracts
- **README.md** — Project overview and setup instructions
