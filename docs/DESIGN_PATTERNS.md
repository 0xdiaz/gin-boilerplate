# Go Project Design Patterns & Architecture Blueprint

**Version:** 2.0
**Last Updated:** 2026-06-10
**Project:** Go Gin Modular Service — Universal Starter Kit

---

## ⚠️ FOR AI AGENTS - READ THIS FIRST

> **🚨 SOURCE OF TRUTH: [`MODULE_GUIDE.md`](./MODULE_GUIDE.md) describes how code is organized in this service.**
> Where this document and `MODULE_GUIDE.md` ever disagree on **folder structure**, `MODULE_GUIDE.md` wins.
>
> **The canonical reference module is [`internal/modules/example`](../internal/modules/example).**
> Copy that module to create a new one. Every code example below is modelled on it.

### 🔥 Non-negotiable patterns (MUST follow)

- **Package-by-feature, not by layer.** A module owns its full vertical slice
  (`handler → service → repository → model`) in one folder: `internal/modules/<name>/`.
- **Consumer-defined interfaces + constructor injection.** The handler defines the
  `service` interface it needs; the service defines the `repository` interface it needs.
  Dependencies are passed in via `New(...)`, never reached for via globals.
- **Cross-module access only through a module's public interface** (`Module.API()`,
  `auth.Servicer`), injected through the constructor. Never import another module's
  unexported types.
- **All HTTP responses go through `pkg/utils`** (`utils.Ok`, `utils.BadRequest`,
  `utils.RespondWithAPIError`, …). Never write `c.JSON(...)` by hand for API responses.
- **`pkg/` is the shared kit and must never import from `internal/`.**

### 📖 How to Use This Document

1. ✅ Read `MODULE_GUIDE.md` first (the layout source of truth, ~100 lines).
2. ✅ Skim the reference module `internal/modules/example/` — it is the living version of this doc.
3. ⚠️  Read the Implementation Patterns and Anti-Patterns sections in full.
4. 📚 Use the rest as reference for detailed patterns.

---

## 📋 Table of Contents

> **Navigation:** anchors are stable; use Ctrl+F on a keyword or click the section links.

### 🔥 Critical Sections (MUST READ)

| Section | Keywords |
|---------|----------|
| [6. Implementation Patterns](#6-implementation-patterns) | `module`, `handler`, `service`, `repository` |
| → [6.1 Handler Pattern](#61-handler-pattern) | `Handler`, `NewHandler`, `service interface` |
| → [6.2 Service Pattern](#62-service-pattern) | `Service`, `NewService`, `repository interface` |
| → [6.3 Repository Pattern](#63-repository-pattern) | `Repository`, `NewRepository(db)`, injected `*gorm.DB` |
| → [6.4 Module Wiring](#64-module-wiring-pattern) | `Module`, `New(db)`, `Models`, `RegisterRoutes`, `API` |
| [13. Anti-Patterns](#13-anti-patterns-to-avoid) | `wrong`, `bad`, `avoid`, `anti-pattern` |

### 📚 All Sections

| # | Section | Keywords |
|---|---------|----------|
| 1 | [Overview](#1-overview) | `philosophy`, `goals`, `principles` |
| 1.1 | [Architecture Philosophy](#11-architecture-philosophy) | `package-by-feature`, `why patterns` |
| 2 | [Project Architecture](#2-project-architecture) | `modules`, `vertical slice`, `separation` |
| 2.1 | [Layers Inside a Module](#21-layers-inside-a-module) | `handler`, `service`, `repository`, `model` |
| 2.2 | [Dependency Flow Rules](#22-dependency-flow-rules) | `dependency`, `direction`, `flow`, `cross-module` |
| 3 | [Core Design Patterns](#3-core-design-patterns) | `patterns`, `repository`, `module`, `DI` |
| 3.1 | [Repository Pattern](#31-repository-pattern) | `repository`, `data access`, `injected db` |
| 3.2 | [Service Layer Pattern](#32-service-layer-pattern) | `service`, `business logic`, `orchestration` |
| 3.3 | [DTO Pattern](#33-dto-data-transfer-object-pattern) | `DTO`, `request`, `response` |
| 3.4 | [Constructor / Module Factory Pattern](#34-constructor--module-factory-pattern) | `New`, `constructor`, `module` |
| 3.5 | [Middleware Pattern](#35-middleware-pattern) | `middleware`, `gin.HandlerFunc`, `auth guard` |
| 3.6 | [Dependency Injection](#36-dependency-injection-pattern) | `DI`, `consumer-defined interface`, `injection` |
| 4 | [Directory Structure](#4-directory-structure-standard) | `directory`, `folder`, `structure`, `tree` |
| 4.1 | [Complete Project Structure](#41-complete-project-structure) | `project tree`, `modular layout` |
| 4.2 | [Package Organization](#42-package-organization-rules) | `package`, `internal`, `pkg`, `module` |
| 5 | [Layer Responsibilities](#5-layer-responsibilities) | `responsibilities`, `what`, `where` |
| 5.1 | [Handler Layer](#51-handler-layer-thin-layer) | `handler`, `thin`, `HTTP`, `validation` |
| 5.2 | [Service Layer](#52-service-layer-fat-layer) | `service`, `fat`, `business logic` |
| 5.3 | [Repository Layer](#53-repository-layer-data-layer) | `repository`, `database`, `CRUD`, `GORM` |
| 6 | [Implementation Patterns](#6-implementation-patterns) | `how to`, `implementation`, `code` |
| 6.1 | [Handler Pattern](#61-handler-pattern) | `Handler`, `NewHandler`, `methods` |
| 6.2 | [Service Pattern](#62-service-pattern) | `Service`, `NewService`, `methods` |
| 6.3 | [Repository Pattern](#63-repository-pattern) | `Repository`, `injected db`, `CRUD` |
| 6.4 | [Module Wiring Pattern](#64-module-wiring-pattern) | `Module`, `New(db)`, `RegisterRoutes`, `API` |
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
| 11 | [Feature Implementation Guide](#11-complete-feature-implementation-guide) | `step by step`, `guide`, `new module` |
| 11.x | [Step-by-Step: New Module](#step-by-step-adding-a-new-module) | `complete example`, `full feature` |
| 12 | [Pattern Examples](#12-pattern-examples-from-codebase) | `examples`, `real code`, `reference` |
| 12.1 | [Auth Pattern Example](#121-auth-pattern-from-the-auth-module) | `auth example`, `authentication` |
| 12.2 | [DataTable Example](#122-datatable-pattern-from-the-example-module) | `datatable`, `server-side` |
| 13 | [Anti-Patterns to Avoid](#13-anti-patterns-to-avoid) | `wrong`, `bad`, `avoid`, `don't` |
| 13.1 | [Business Logic in Handler](#131--business-logic-in-handler) | `handler anti-pattern`, `fat handler` |
| 13.2 | [Direct DB in Service](#132--direct-database-access-in-service) | `service anti-pattern`, `tight coupling` |
| 13.3 | [Standalone Functions](#133--standalone-handler-functions) | `standalone`, `function anti-pattern` |
| 13.4 | [God Service](#134--god-service-too-many-responsibilities) | `god service`, `SRP violation` |
| 13.5 | [Reaching Into Another Module](#135--reaching-into-another-modules-internals) | `import cycle`, `cross-module` |

### 🎯 Quick Lookups by Task

**Implementing Handlers (HTTP layer):**
- [6.1 Handler Pattern](#61-handler-pattern)
- [5.1 Handler Responsibilities](#51-handler-layer-thin-layer)
- [13.1 What NOT to do](#131--business-logic-in-handler)

**Implementing Services:**
- [6.2 Service Pattern](#62-service-pattern)
- [5.2 Service Responsibilities](#52-service-layer-fat-layer)
- [13.2 What NOT to do](#132--direct-database-access-in-service)

**Implementing Repositories:**
- [6.3 Repository Pattern](#63-repository-pattern)
- [5.3 Repository Responsibilities](#53-repository-layer-data-layer)
- [3.1 Repository Pattern Theory](#31-repository-pattern)

**Wiring a Module:**
- [6.4 Module Wiring Pattern](#64-module-wiring-pattern)
- [11 Step-by-Step: New Module](#step-by-step-adding-a-new-module)

**Understanding Flow:**
- [7.1 Standard CRUD Flow](#71-standard-crud-flow)
- [7.2 Authentication Flow](#72-authentication-flow-pattern)
- [7.3 Transaction Flow](#73-transaction-flow-pattern)

---

## 1. OVERVIEW

### 1.1 Architecture Philosophy

This project follows a **modular (package-by-feature) architecture**. Each service is one
deployable (one binary, one container, one Postgres), and the code is organized **by business
module**, not by technical layer. A module owns its full vertical slice — `handler → service →
repository → model` — inside a single folder.

The Clean Architecture **principles** (separation of concerns, dependency inversion, dependencies
pointing inward, testability) are still the law. What changed is *where* those layers live: they
now sit **together inside one module** rather than being scattered across `internal/app/...` and
`internal/domain/...` folders.

```
┌─────────────────────────────────────────────┐
│          External Systems / HTTP            │
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│   internal/modules/<name>/  (one vertical slice)
│                                             │
│   handler.go     (Thin Layer)               │  ← HTTP: parse, call service, respond
│        │                                    │
│        ▼                                    │
│   service.go     (Fat Layer)                │  ← Business rules, orchestration, logging
│        │                                    │
│        ▼                                    │
│   repository.go  (Data Layer)               │  ← CRUD on an injected *gorm.DB
│        │                                    │
│        ▼                                    │
│   model.go       (the tables this module owns)
└──────────────────┬──────────────────────────┘
                   │
┌──────────────────▼──────────────────────────┐
│              Database (Postgres)             │
└─────────────────────────────────────────────┘
```

**Key Principles:**
1. **One module = one Go package = one folder** — owned end-to-end without touching others.
2. **Dependency Inversion** — each layer depends on a *consumer-defined interface*, not a concrete type.
3. **Separation of Concerns** — each file in the slice has ONE responsibility.
4. **Dependency Direction** — `handler → service → repository`; never the reverse.
5. **Cross-module communication** — only through a module's public interface, injected via the constructor.
6. **Testability** — services and handlers are tested with fakes, no DB required.

> The canonical implementation of every principle above is `internal/modules/example`.
> `MODULE_GUIDE.md` is the source of truth for the layout.

---

## 2. PROJECT ARCHITECTURE

### 2.1 Layers Inside a Module

Within `internal/modules/<name>/`, the layers are co-located files:

```
┌──────────────────────────────────────────────────────┐
│                  PRESENTATION                        │
│  • handler.go   — HTTP handlers (gin), DTO binding   │
│  • dto.go       — request/response structures        │
│  • model.go     — the GORM model(s) the module owns  │
└──────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────┐
│                  APPLICATION                          │
│  • service.go   — business logic & orchestration     │
│                   (defines the `repository` it needs)│
└──────────────────────────────────────────────────────┘
                         ↓
┌──────────────────────────────────────────────────────┐
│                  DATA ACCESS                          │
│  • repository.go — CRUD on the injected *gorm.DB     │
└──────────────────────────────────────────────────────┘

         module.go — wires it all: New(db), Name(),
         Models(), RegisterRoutes(api), public API()
```

Cross-cutting infrastructure lives outside the module:

```
┌──────────────────────────────────────────────────────┐
│  pkg/                  (the shared "kit-in-waiting")  │
│  • middleware/  — cors, request_id, request_log,     │
│                   metrics, rate_limit                 │
│  • database/    — connect + replica resolver         │
│  • config/ logger/ metrics/ types/ utils/            │
└──────────────────────────────────────────────────────┘
┌──────────────────────────────────────────────────────┐
│  internal/bootstrap/   (per-service wiring ONLY)     │
│  • bootstrap.go — Run(): config→db→modules→migrate   │
│  • modules.go   — Module interface + buildModules()  │
│  • server.go    — gin engine + global middleware     │
│  • swagger.go   — OpenAPI/Swagger UI (debug only)    │
└──────────────────────────────────────────────────────┘
```

### 2.2 Dependency Flow Rules

**✅ ALLOWED:**
```
Handler → Service → Repository → Database
   ↓         ↓           ↓
  DTO    Utils/Types   Model

internal/bootstrap → internal/modules/*     // wiring knows the concrete modules
internal/modules/* → pkg/*                  // modules use the shared kit
```

**❌ FORBIDDEN:**
```
Service → Handler            // Services cannot depend on the HTTP layer
Repository → Service         // Repositories cannot depend on business logic
Model → Repository           // Models are pure data structures
pkg/* → internal/*           // The shared kit must NEVER import internal
moduleA → moduleB's guts     // Never import another module's unexported types
```

**🔑 Cardinal cross-module rule:** modules never import each other's internals. A module that
needs another module receives that module's **public interface** (e.g. `auth.Servicer` via
`authMod.Auth()`, or the consumer module's own `API`) through its **constructor**, wired in
`buildModules()`:

```go
// in internal/bootstrap/modules.go
authMod := auth.New(db)
payments := payments.New(db, authMod.Auth())   // auth.Servicer injected
```

This keeps modules loosely coupled and makes the seam easy to cut later: if a module must become
its own service, the in-process interface call becomes a network call (REST + HMAC) and nothing
else changes.

---

## 3. CORE DESIGN PATTERNS

### 3.1 Repository Pattern

**Purpose:** Abstract data access from business logic.

**Structure:** A repository is a **struct that holds an injected `*gorm.DB`** — no package globals.
The *interface* the repository satisfies is **defined by the service that consumes it** (see §3.6),
not by the repository itself.

```go
// File: internal/modules/example/repository.go
package example

import "gorm.io/gorm"

// Repository is the data-access layer for the example module.
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
- Return the module's own models, not DTOs.
- Translate database-specific errors here (e.g. `gorm.ErrRecordNotFound` → `nil, nil`).
- Use transactions for multi-step operations.

**❌ DON'T:**
- Use the package-level `database.DB` global (it exists only for the connection lifecycle/tests).
- Put business logic in repositories.
- Call another module's repository directly.
- Import the service layer.
- Transform data for API responses (that is the DTO's job).

---

### 3.2 Service Layer Pattern

**Purpose:** Encapsulate business logic and orchestrate operations.

**Structure:** A service holds a **consumer-defined `repository` interface** (declared in
`service.go` itself, see §3.6). It logs each operation with `logger.LogStart`/`LogFinish` and
propagates the request-scoped `context.Context`.

```go
// File: internal/modules/example/service.go
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

// Service holds the example module's business logic.
type Service struct {
    repo repository
}

// NewService creates a Service backed by the given repository.
func NewService(repo repository) *Service {
    return &Service{repo: repo}
}

// List returns all example records.
func (s *Service) List(ctx context.Context) ([]*Example, error) {
    ctx, start := logger.LogStart(ctx, "example.Service.List")
    list, err := s.repo.List()
    logger.LogFinish(ctx, "example.Service.List", err, start)
    return list, err
}
```

**✅ DO:**
- Implement ALL business logic and rules here.
- Define the `repository` interface this service needs *in this file* (dependency inversion).
- Accept and propagate `context.Context` as the first argument.
- Wrap each operation in `logger.LogStart` / `logger.LogFinish` for tracing.
- Validate business constraints and orchestrate multiple repository calls.
- Return module-defined sentinel errors (e.g. `ErrInvalidCredentials`) for known conditions.

**❌ DON'T:**
- Handle HTTP concerns (`gin.Context`, status codes) — except where a repository helper such as
  DataTables genuinely needs `*gin.Context` for server-side paging.
- Access the database directly (go through the repository).
- Import the handler layer.
- Return HTTP responses.

---

### 3.3 DTO (Data Transfer Object) Pattern

**Purpose:** Define contracts for API requests/responses and prevent tight coupling.

**Structure:** DTOs live in `internal/modules/<name>/dto.go` (the file is optional; small modules
may keep request/response types alongside the model). They are pure data structures with binding tags.

```go
// File: internal/modules/auth/dto.go
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
- Keep DTOs in the owning module (`dto.go`), never in a shared `dto` package.
- Exclude sensitive fields from responses (e.g. the auth `User` model tags `Password` as `json:"-"`).

**❌ DON'T:**
- Expose database models directly via API when they carry secrets.
- Add business logic to DTOs.
- Reuse a request DTO as a response DTO.
- Put DTOs in another module's package.

---

### 3.4 Constructor / Module Factory Pattern

**Purpose:** Centralize object creation and wire a module's slice together.

**Structure:** Every layer has a `New*` constructor that takes its dependencies. The module's
`New(db)` is the single factory that assembles repository → service → handler.

```go
// File: internal/modules/example/module.go
package example

import "gorm.io/gorm"

// New builds the module from an injected DB connection.
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))     // repo → service
    return &Module{svc: svc, handler: NewHandler(svc)} // service → handler
}
```

**✅ DO:**
- Use `New*` functions for all constructors; pass dependencies in.
- Assemble the slice exactly once in `Module.New(db)`.
- Return concrete types from `New(db)` (the `*Module`); expose the public surface via `Module.API()`.

**❌ DON'T:**
- Build objects with `&Struct{}` directly in business logic.
- Use `init()` functions for dependency initialization.
- Create global singletons for module state.

---

### 3.5 Middleware Pattern

**Purpose:** Handle cross-cutting concerns (CORS, request IDs, logging, metrics, rate limiting, auth).

**Two kinds of middleware:**

1. **Generic middleware lives in `pkg/middleware/`** and is applied globally in
   `internal/bootstrap/server.go`:
   - `middleware.CORSMiddleware()`
   - `middleware.RequestIDMiddleware()` — flows a request ID into the context
   - `middleware.RequestLogMiddleware()`
   - `middleware.MetricsMiddleware()`
   - `middleware.RateLimitMiddleware()` / `middleware.RateLimitMiddlewareWithConfig(rps, burst)`
     — per **client IP**, applied to all of `/api/v1`.

2. **The JWT auth guard is NOT in `pkg/middleware`.** It lives **inside the auth module** and is
   exposed via `auth.Module.Middleware()`. A module that needs authentication takes the guard from
   the auth module rather than implementing its own.

```go
// The auth module owns its guard (internal/modules/auth/middleware.go) and
// exposes it (internal/modules/auth/module.go):
func (m *Module) Middleware() gin.HandlerFunc {
    return authMiddleware(m.svc) // validates Bearer JWT, sets "user_id" in the context
}

// Protected routes inside the auth module:
protected := api.Group("")
protected.Use(m.Middleware())
protected.GET("/profile", m.handler.Profile)
```

```go
// Generic middleware is mounted globally in internal/bootstrap/server.go:
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

The cornerstone is the **consumer-defined interface**: the layer that *uses* a dependency declares
the (usually unexported) interface it needs, in its own file. The concrete implementation satisfies
that interface implicitly. Wiring happens in `module.go`.

```go
// 1. The handler declares the `service` interface it needs (handler.go):
type service interface {
    List(ctx context.Context) ([]*Example, error)
}
type Handler struct{ svc service }
func NewHandler(svc service) *Handler { return &Handler{svc: svc} }

// 2. The service declares the `repository` interface it needs (service.go):
type repository interface {
    List() ([]*Example, error)
}
type Service struct{ repo repository }
func NewService(repo repository) *Service { return &Service{repo: repo} }

// 3. module.go wires the concrete chain (repo → service → handler):
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))
    return &Module{svc: svc, handler: NewHandler(svc)}
}
```

> The auth module shows the *exported* variant: it publishes a `Servicer` interface
> (`internal/modules/auth/servicer.go`) so its handler, middleware, **and other modules** can all
> depend on the contract rather than the concrete `*Service`.

**✅ Benefits:**
- Easy to test (inject a fake) — see `example/service_test.go`.
- The interface lists exactly what the consumer uses, nothing more.
- No global state; dependencies are explicit in the constructor.
- Cross-module seams are clean and can later become network calls.

---

## 4. DIRECTORY STRUCTURE STANDARD

### 4.1 Complete Project Structure

> Mirrors `MODULE_GUIDE.md` — that file is the source of truth.

```
gin-boilerplate/
│
├── main.go                          # 3 lines: bootstrap.Run()
│
├── internal/                        # Private application code
│   │
│   ├── bootstrap/                   # The ONLY per-service wiring
│   │   ├── bootstrap.go             #   Run(): config → db → modules → migrate → serve → shutdown
│   │   ├── modules.go               #   Module interface + buildModules() (the module list)
│   │   ├── server.go                #   gin engine + global middleware + route mounting
│   │   └── swagger.go               #   OpenAPI/Swagger UI (debug only)
│   │
│   ├── migrations/                  # Service-specific schema
│   │   ├── migration.go             #   Run(db, models) — AutoMigrate
│   │   └── sql/                     #   versioned SQL (production path)
│   │
│   └── modules/                     # ← business modules (work happens here)
│       │
│       ├── example/                 #   THE reference module — copy this to make a new one
│       │   ├── model.go             #     GORM model(s) the module owns
│       │   ├── dto.go               #     request/response types (optional)
│       │   ├── repository.go        #     data access (holds *gorm.DB, no globals)
│       │   ├── service.go           #     business logic (defines the repo interface it needs)
│       │   ├── handler.go           #     HTTP layer (defines the service interface it needs)
│       │   ├── module.go            #     New(db), Name(), Models(), RegisterRoutes(), public API()
│       │   └── service_test.go      #     co-located test (no DB — uses a fake repo)
│       │
│       ├── auth/                    #   real module: JWT auth, exposes Middleware() and Auth()
│       │   ├── model.go  dto.go  repository.go
│       │   ├── service.go  service_tokens.go  servicer.go
│       │   ├── handler.go  middleware.go  mailer.go
│       │   └── module.go
│       │
│       └── health/                  #   system module: mounts /health, /metrics at ROOT
│           ├── service.go  handler.go  dto.go
│           └── module.go
│
├── pkg/                             # Public reusable kit — must NEVER import internal/
│   ├── config/                      # Configuration loader
│   ├── logger/                      # Logging + LogStart/LogFinish tracing helpers
│   ├── metrics/                     # Request counters, uptime
│   ├── types/                       # SuccessResponse, ErrorResponse, APIError
│   ├── utils/                       # Response helpers (Ok, BadRequest, RespondWithAPIError, …)
│   ├── middleware/                  # cors, request_id, request_log, metrics, rate_limit
│   └── database/                    # DbConnection(master, replica), GetDB(), package global DB
│
├── api/                             # Embedded OpenAPI spec (served by swagger.go)
├── docs/                            # Documentation
│   ├── MODULE_GUIDE.md              # ← source of truth for the layout
│   ├── DESIGN_PATTERNS.md           # This file
│   ├── CODING_STANDARDS.md
│   ├── AI_AGENT_RULES.md
│   └── ...
│
├── tests/                           # Cross-cutting tests & helpers
│   ├── unit/                        # Unit tests
│   ├── integration/                 # Integration tests
│   ├── mocks/                       # Shared mocks (legacy; new modules prefer in-package fakes)
│   └── fixtures/                    # Test data
│
├── scripts/                         # Dev/ops scripts
├── go.mod                           # Module path: github.com/0xdiaz/gin-boilerplate
├── go.sum
├── Makefile                         # `make test` → ./tests/unit/... ./internal/... ./pkg/...
└── README.md
```

> **Module path note:** this boilerplate uses
> `github.com/0xdiaz/gin-boilerplate`. Each stamped service should set its own path in `go.mod`
> (e.g. `github.com/your-org/your-service`) and update imports — a one-shot find/replace verified
> with `go build ./...`.

### 4.2 Package Organization Rules

**✅ RULES:**

1. **One module = one Go package = one folder.**
   ```
   ✅ internal/modules/example/   → package example  (handler, service, repository, model, module)
   ❌ internal/modules/           → multiple feature packages mixed at one level
   ```

2. **A module owns its tables.** Declare them in `Module.Models()`. No cross-module foreign keys;
   no reaching into another module's tables.

3. **`pkg/` is the shared kit.** It must never import from `internal/`. Anything you want identical
   across sibling services goes here.

4. **`internal/bootstrap/` is the only place that knows the concrete module list.** Adding a module
   is one line in `buildModules()`.

5. **Package naming:**
   - Lowercase, single word, no underscores or dashes.
   - Module packages are singular feature names: `example`, `auth`, `health`.
   - Utility packages are singular: `logger`, `config`, `utils`, `middleware`.

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

### 5.1 Handler Layer (Thin Layer)

**Responsibility:** Handle HTTP concerns ONLY.

**What handlers SHOULD do:**
```go
// File: internal/modules/example/handler.go
func (h *Handler) List(c *gin.Context) {
    // 1. Open a trace span from the request context
    ctx, start := logger.LogStart(c.Request.Context(), "example.Handler.List")

    // 2. Call the service
    data, err := h.svc.List(ctx)
    if err != nil {
        logger.LogFinish(ctx, "example.Handler.List", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }

    // 3. Respond with a response utility
    logger.LogFinish(ctx, "example.Handler.List", nil, start)
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
- Import another module's internals.

**Size Limits:**
- **Maximum 50 lines** per handler method.
- **Maximum 300 lines** per `handler.go` file (split helpers out if larger).

---

### 5.2 Service Layer (Fat Layer)

**Responsibility:** Implement ALL business logic.

**What services SHOULD do:**
```go
// File: internal/modules/auth/service.go (abridged)
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    ctx, start := logger.LogStart(ctx, "auth.Service.Register")

    // 1. Business constraint: email must be unique
    existing, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        logger.LogFinish(ctx, "auth.Service.Register", err, start)
        return nil, fmt.Errorf("failed to check email: %w", err)
    }
    if existing != nil {
        logger.LogFinish(ctx, "auth.Service.Register", ErrEmailAlreadyExists, start)
        return nil, ErrEmailAlreadyExists // sentinel error
    }

    // 2. Apply business logic (hash password, build model)
    hashed, err := s.hashPassword(req.Password)
    if err != nil {
        logger.LogFinish(ctx, "auth.Service.Register", err, start)
        return nil, fmt.Errorf("failed to process password: %w", err)
    }
    user := &User{Name: req.Name, Email: req.Email, Password: hashed}

    // 3. Orchestrate repository calls
    if err = s.userRepo.CreateUser(user); err != nil {
        logger.LogFinish(ctx, "auth.Service.Register", err, start)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    // 4. Build response, log the business event
    logger.Infof("user registered successfully: %s", user.Email)
    logger.LogFinish(ctx, "auth.Service.Register", nil, start)
    return &AuthResponse{ /* ... */ }, nil
}
```

**✅ Services SHOULD:**
- Implement ALL business logic and rules.
- Accept and propagate `context.Context`.
- Wrap operations in `logger.LogStart` / `logger.LogFinish`.
- Orchestrate multiple repository calls; call external services through injected interfaces.
- Return module sentinel errors for known conditions, wrap unexpected errors with `%w`.

**❌ Services MUST NOT:**
- Handle HTTP status codes or write responses.
- Access the database directly (use the repository).
- Import the handler layer.

**Size Limits:**
- **Maximum 100 lines** per service method.
- **Maximum 400 lines** per service file; split into `service.go`, `service_<topic>.go`
  (the auth module splits token logic into `service_tokens.go`).

---

### 5.3 Repository Layer (Data Layer)

**Responsibility:** Handle data persistence ONLY, against the **injected** `*gorm.DB`.

**What repositories SHOULD do:**
```go
// File: internal/modules/auth/repository.go (abridged)
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

> The auth module declares an **exported** `Repository` interface plus an unexported `repository`
> struct implementation, because the interface is its module-internal data-access contract. The
> example module keeps it simpler — a plain exported `Repository` struct, with the *interface* it
> satisfies defined by the consuming service. Both are valid; prefer the example's simpler shape
> for new modules unless you need to swap implementations.

**✅ Repositories SHOULD:**
- Hold the injected `*gorm.DB`; never use the `database.DB` global.
- Perform CRUD and build queries with GORM/parameterized SQL.
- Translate database-specific errors (`gorm.ErrRecordNotFound` → `nil, nil`).
- Use transactions for multi-step operations.
- Optimize queries (`Preload`, `Select`, paging).
- Return the module's own models.

**❌ Repositories MUST NOT:**
- Contain business logic or validate business rules.
- Call another module's repository.
- Import the service or handler layer.

**Size Limits:**
- **Maximum 30 lines** per repository method.
- **Maximum 300 lines** per repository file.

---

## 6. IMPLEMENTATION PATTERNS

### 6.1 Handler Pattern

**✅ REQUIRED Pattern (struct + consumer-defined `service` interface):**
```go
// File: internal/modules/example/handler.go
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

// Handler is the HTTP layer for the example module.
type Handler struct {
    svc service
}

// NewHandler creates a Handler backed by the given service.
func NewHandler(svc service) *Handler {
    return &Handler{svc: svc}
}

// List handles GET /examples.
func (h *Handler) List(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "example.Handler.List")
    data, err := h.svc.List(ctx)
    if err != nil {
        logger.LogFinish(ctx, "example.Handler.List", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }
    logger.LogFinish(ctx, "example.Handler.List", nil, start)
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
// File: internal/modules/example/service.go
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
    ctx, start := logger.LogStart(ctx, "example.Service.List")
    list, err := s.repo.List()
    logger.LogFinish(ctx, "example.Service.List", err, start)
    return list, err
}
```

---

### 6.3 Repository Pattern

**✅ REQUIRED Pattern (struct holding the injected `*gorm.DB`):**
```go
// File: internal/modules/example/repository.go
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

### 6.4 Module Wiring Pattern

A module is the single object `bootstrap` touches. It satisfies the `Module` interface
(`Name() string`, `Models() []any`, `RegisterRoutes(api *gin.RouterGroup)`) and exposes a minimal
public `API` for other modules.

```go
// File: internal/modules/example/module.go
package example

import (
    "context"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// API is the public surface other modules depend on. Keep it minimal — other
// modules call this, they never import the example module's internals.
type API interface {
    List(ctx context.Context) ([]*Example, error)
}

type Module struct {
    svc     *Service
    handler *Handler
}

// New builds the module from an injected DB connection.
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))
    return &Module{svc: svc, handler: NewHandler(svc)}
}

func (m *Module) Name() string   { return "example" }
func (m *Module) Models() []any  { return []any{&Example{}} }

// RegisterRoutes mounts the module's routes under the given API group (e.g. /api/v1).
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    api.GET("/examples", m.handler.List)
    api.GET("/datatables", m.handler.Datatables)
}

// API exposes this module's public contract to other modules.
func (m *Module) API() API { return m.svc }
```

**Registering the module** — one line in `buildModules()`:
```go
// File: internal/bootstrap/modules.go
func buildModules(db *gorm.DB) []Module {
    return []Module{
        auth.New(db),
        example.New(db),
        // newmodule.New(db),  ← add here
    }
}
```

**Mounting** happens in `internal/bootstrap/server.go`: business modules go under `/api/v1` with
rate limiting; the `health` **system** module is special — it implements `New(db)` +
`RegisterSystem(r)` and mounts `/health` and `/metrics` at the **root**, not under `/api/v1`.

```go
// File: internal/bootstrap/server.go (abridged)
health.New(db).RegisterSystem(r)            // /health, /metrics at ROOT

v1 := r.Group("/api/v1")
v1.Use(middleware.RateLimitMiddleware())
for _, m := range mods {
    m.RegisterRoutes(v1)                    // every business module under /api/v1
}
```

Migrations are automatic: `bootstrap` collects every module's `Models()` and passes them to
`migrations.Run(db, models)`.

---

### 6.5 Response Utility Pattern

**✅ MANDATORY: Always use the response utilities in `pkg/utils`.**

```go
import "github.com/0xdiaz/gin-boilerplate/pkg/utils"

func (h *Handler) Create(c *gin.Context) {
    var req CreateRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data") // 400 + field errors
        return
    }

    data, err := h.svc.Create(c.Request.Context(), &req)
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

**Complete request flow for a LIST/CREATE operation (example module):**

```
1. HTTP GET /api/v1/examples
       ↓
2. [Global Middleware] (internal/bootstrap/server.go)
   • gin.Recovery
   • CORSMiddleware
   • RequestIDMiddleware   (request id → context)
   • RequestLogMiddleware
   • MetricsMiddleware
       ↓
3. [/api/v1 group]
   • RateLimitMiddleware   (per client IP)
       ↓
4. [Module route]   example.Module.RegisterRoutes → m.handler.List
       ↓
5. [Handler] internal/modules/example/handler.go
   func (h *Handler) List(c *gin.Context) {
       ctx, start := logger.LogStart(c.Request.Context(), "example.Handler.List")
       data, err := h.svc.List(ctx)            // call service via `service` interface
       if err != nil { utils.InternalServerError(c, err, "Failed to retrieve data"); return }
       logger.LogFinish(ctx, "example.Handler.List", nil, start)
       utils.Ok(c, data, "Data retrieved successfully")
   }
       ↓
6. [Service] internal/modules/example/service.go
   func (s *Service) List(ctx context.Context) ([]*Example, error) {
       ctx, start := logger.LogStart(ctx, "example.Service.List")
       list, err := s.repo.List()              // call repo via `repository` interface
       logger.LogFinish(ctx, "example.Service.List", err, start)
       return list, err
   }
       ↓
7. [Repository] internal/modules/example/repository.go
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

The auth module uses **JWT directly** (no OTP step). Routes are mounted by
`auth.Module.RegisterRoutes` under `/api/v1`.

```
PUBLIC routes (no auth):
  POST /api/v1/auth/register         → Handler.Register
  POST /api/v1/auth/login            → Handler.Login
  POST /api/v1/auth/refresh          → Handler.RefreshToken
  POST /api/v1/auth/forgot-password  → Handler.ForgotPassword
  POST /api/v1/auth/reset-password   → Handler.ResetPassword

PROTECTED route (auth.Module.Middleware() guard):
  GET  /api/v1/profile               → Handler.Profile
```

**Login flow:**
```
1. POST /api/v1/auth/login   {"email": "...", "password": "..."}
       ↓
2. [Handler.Login] bind LoginRequest → service.Login(ctx, &req)
       ↓
3. [Service.Login]
   • userRepo.GetUserByEmail(email)
   • bcrypt.CompareHashAndPassword(...)        (invalid → ErrInvalidCredentials)
   • generateToken(user)                        (signed JWT, user_id + email, 24h)
   • generateRefreshToken()                     (32 random bytes → hex)
   • userRepo.UpdateUser(user)                  (persist refresh token)
   • return *AuthResponse{User, AccessToken, RefreshToken, TokenType:"Bearer"}
       ↓
4. [Handler.Login] on error: errToAPIError(err) → utils.RespondWithAPIError
                   on success: utils.Ok(c, response, "Login successful")
       ↓
5. Response: { "success": true, "data": { "user": {...},
               "access_token": "eyJ...", "refresh_token": "...", "token_type": "Bearer" } }

Subsequent protected requests:
6. GET /api/v1/profile   Header: Authorization: Bearer eyJ...
       ↓
7. [auth.Module.Middleware()]  internal/modules/auth/middleware.go
   • split "Bearer <token>"
   • service.ValidateToken(token) → userID   (HS256, checks signing method & exp)
   • c.Set("user_id", userID); c.Next()       (401 via utils.Unauthorized on failure)
       ↓
8. [Handler.Profile]  userID := c.GetUint("user_id")
```

> The JWT secret comes from config (`config.Get().Server.JWTSecret`). The guard lives in the auth
> module and is exposed via `Module.Middleware()`; other modules reuse it rather than reimplementing.

---

### 7.3 Transaction Flow Pattern

**Multi-step operations use a GORM transaction inside the repository**, on the injected `*gorm.DB`:

```go
// File: internal/modules/<name>/repository.go
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
[Handler] c.ShouldBindJSON → DTO
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
[Handler] utils.Ok / utils.Created → standard JSON envelope
{
    "success": true,
    "message": "...",
    "data": { "user": {...}, "access_token": "..." },
    "errors": null
}
```

**Key Transformations:**
1. **Handler:** JSON → DTO (binding/validation at the boundary).
2. **Service:** DTO → domain model (business logic); domain model → response DTO.
3. **Repository:** domain model ↔ database.
4. **Handler:** response DTO → standard JSON envelope via `pkg/utils`.

---

### 8.2 DTO vs Model Usage

**✅ When to use DTOs:**
- API request payloads (binding/validation).
- API response payloads (shape the client sees; strip secrets).
- External/inter-service communication.

**✅ When to use Models:**
- Internal business logic and domain rules.
- Database operations in the repository.

**Example (auth module):**
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
[Handler] map sentinel → APIError, else 500
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

The pattern is: **each module defines sentinel errors; the handler maps them to `types.APIError`
and responds via `utils.RespondWithAPIError`.** `pkg/types` provides the `APIError` type
(`{Code, Message, Details}`) and a few shared predefined errors.

```go
// 1. The module defines sentinel errors (internal/modules/auth/service.go):
var (
    ErrEmailAlreadyExists  = errors.New("email already exists")
    ErrInvalidCredentials  = errors.New("invalid email or password")
    ErrUserNotFound        = errors.New("user not found")
    ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
    // ...
)

// 2. The handler maps them to types.APIError (internal/modules/auth/handler.go):
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
`types.ErrConflict`, `types.ErrRateLimitExceeded`, …) for cases where a module does not need its
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
Handler:     HTTP 500 (unknown) — or a mapped APIError for a known sentinel
```

---

## 10. TESTING PATTERNS

> **Prefer co-locating tests with the module** (`internal/modules/<name>/*_test.go`). Unit-test the
> service with a hand-written **fake repository** (no DB). Handler tests use `httptest` + a mocked
> service interface. Shared mocks for legacy tests live in `tests/mocks/`; new modules should prefer
> in-package fakes. `make test` runs `./tests/unit/... ./internal/... ./pkg/...`.

### 10.1 Service Layer Testing Pattern

The canonical example is `internal/modules/example/service_test.go`: a `fakeRepo` satisfies the
unexported `repository` interface, so no database is needed.

```go
// File: internal/modules/example/service_test.go
package example

import (
    "context"
    "testing"

    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

// fakeRepo satisfies the unexported repository interface — no DB needed.
type fakeRepo struct {
    items []*Example
    err   error
}

func (f *fakeRepo) List() ([]*Example, error)                      { return f.items, f.err }
func (f *fakeRepo) Datatables(c *gin.Context) (interface{}, error) { return nil, f.err }

func TestService_List(t *testing.T) {
    svc := NewService(&fakeRepo{items: []*Example{{ID: 1, Data: "x"}}})

    got, err := svc.List(context.Background())

    assert.NoError(t, err)
    assert.Len(t, got, 1)
    assert.Equal(t, "x", got[0].Data)
}

func TestModule_Meta(t *testing.T) {
    m := New(nil)
    assert.Equal(t, "example", m.Name())
    assert.Len(t, m.Models(), 1)
}
```

Because the test is **in-package** (`package example`), the fake can satisfy the unexported
`repository` interface directly. To test an error path, set `fakeRepo.err` and assert it propagates.

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

### Step-by-Step: Adding a New Module

> This follows `MODULE_GUIDE.md`'s recipe: **copy `internal/modules/example`**. The example module
> is the canonical reference; every step below mirrors a real file in it.

#### Step 0: Copy the reference module
```bash
cp -r internal/modules/example internal/modules/product
# rename the package from `example` to `product` in every file
```

#### Step 1: Define the model (`model.go`)
```go
// File: internal/modules/product/model.go
package product

import "time"

// Product is a model this module owns. Each module owns its own tables.
type Product struct {
    ID          int        `json:"id" gorm:"primaryKey"`
    Name        string     `json:"name" binding:"required"`
    Price       float64    `json:"price"`
    CreatedAt   *time.Time `json:"created_at"`
    UpdatedAt   *time.Time `json:"updated_at"`
}

func (p *Product) TableName() string { return "products" }
```

#### Step 2: Define DTOs (`dto.go`, optional)
```go
// File: internal/modules/product/dto.go
package product

type CreateProductRequest struct {
    Name  string  `json:"name" binding:"required,min=3,max=255"`
    Price float64 `json:"price" binding:"required,gt=0"`
}
```

#### Step 3: Implement the repository (`repository.go`)
```go
// File: internal/modules/product/repository.go
package product

import "gorm.io/gorm"

type Repository struct {
    db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(p *Product) error { return r.db.Create(p).Error }

func (r *Repository) List() ([]*Product, error) {
    var list []*Product
    return list, r.db.Find(&list).Error
}
```

#### Step 4: Implement the service (`service.go`) — define the `repository` interface here
```go
// File: internal/modules/product/service.go
package product

import (
    "context"
    "fmt"

    "github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

// repository is the data-access contract this service needs (consumer-defined).
type repository interface {
    Create(p *Product) error
    List() ([]*Product, error)
}

type Service struct {
    repo repository
}

func NewService(repo repository) *Service { return &Service{repo: repo} }

func (s *Service) Create(ctx context.Context, req *CreateProductRequest) (*Product, error) {
    ctx, start := logger.LogStart(ctx, "product.Service.Create")
    if req.Price <= 0 {
        err := fmt.Errorf("price must be greater than 0")
        logger.LogFinish(ctx, "product.Service.Create", err, start)
        return nil, err
    }
    p := &Product{Name: req.Name, Price: req.Price}
    err := s.repo.Create(p)
    logger.LogFinish(ctx, "product.Service.Create", err, start)
    return p, err
}

func (s *Service) List(ctx context.Context) ([]*Product, error) {
    ctx, start := logger.LogStart(ctx, "product.Service.List")
    list, err := s.repo.List()
    logger.LogFinish(ctx, "product.Service.List", err, start)
    return list, err
}
```

#### Step 5: Implement the handler (`handler.go`) — define the `service` interface here
```go
// File: internal/modules/product/handler.go
package product

import (
    "context"

    "github.com/0xdiaz/gin-boilerplate/pkg/logger"
    "github.com/0xdiaz/gin-boilerplate/pkg/utils"
    "github.com/gin-gonic/gin"
)

// service is the business contract this handler needs (consumer-defined).
type service interface {
    Create(ctx context.Context, req *CreateProductRequest) (*Product, error)
    List(ctx context.Context) ([]*Product, error)
}

type Handler struct {
    svc service
}

func NewHandler(svc service) *Handler { return &Handler{svc: svc} }

func (h *Handler) Create(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "product.Handler.Create")
    var req CreateProductRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.LogFinish(ctx, "product.Handler.Create", err, start)
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    p, err := h.svc.Create(ctx, &req)
    if err != nil {
        logger.LogFinish(ctx, "product.Handler.Create", err, start)
        utils.InternalServerError(c, err, "Failed to create product")
        return
    }
    logger.LogFinish(ctx, "product.Handler.Create", nil, start)
    utils.Created(c, p, "Product created successfully")
}

func (h *Handler) List(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "product.Handler.List")
    list, err := h.svc.List(ctx)
    if err != nil {
        logger.LogFinish(ctx, "product.Handler.List", err, start)
        utils.InternalServerError(c, err, "Failed to list products")
        return
    }
    logger.LogFinish(ctx, "product.Handler.List", nil, start)
    utils.Ok(c, list, "Products retrieved successfully")
}
```

#### Step 6: Wire it (`module.go`)
```go
// File: internal/modules/product/module.go
package product

import (
    "context"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type API interface {
    List(ctx context.Context) ([]*Product, error)
}

type Module struct {
    svc     *Service
    handler *Handler
}

func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))
    return &Module{svc: svc, handler: NewHandler(svc)}
}

func (m *Module) Name() string  { return "product" }
func (m *Module) Models() []any { return []any{&Product{}} }

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    g := api.Group("/products")
    g.GET("", m.handler.List)
    g.POST("", m.handler.Create)
}

func (m *Module) API() API { return m.svc }
```

#### Step 7: Register the module — one line
```go
// File: internal/bootstrap/modules.go
func buildModules(db *gorm.DB) []Module {
    return []Module{
        auth.New(db),
        example.New(db),
        product.New(db),   // ← added
    }
}
```

#### Step 8: Add a co-located test (`service_test.go`) using a fake repo
```go
// File: internal/modules/product/service_test.go
package product

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

type fakeRepo struct {
    items []*Product
    err   error
}

func (f *fakeRepo) Create(p *Product) error      { return f.err }
func (f *fakeRepo) List() ([]*Product, error)    { return f.items, f.err }

func TestService_List(t *testing.T) {
    svc := NewService(&fakeRepo{items: []*Product{{ID: 1, Name: "Widget"}}})
    got, err := svc.List(context.Background())
    assert.NoError(t, err)
    assert.Len(t, got, 1)
}
```

That's it — no other file changes. Routes mount under `/api/v1/products`; migrations pick up the
`Product` model automatically via `Module.Models()`.

---

## 12. PATTERN EXAMPLES FROM CODEBASE

### 12.1 Auth Pattern (from the `auth` module)

**✅ Handler maps domain errors and uses response utilities:**
```go
// internal/modules/auth/handler.go
func (h *Handler) Login(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "auth.Handler.Login")

    var req LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.LogFinish(ctx, "auth.Handler.Login", err, start)
        utils.BadRequest(c, err, "Invalid request data")
        return
    }

    response, err := h.service.Login(ctx, &req)
    if err != nil {
        if apiErr := errToAPIError(err); apiErr != nil {
            logger.LogFinish(ctx, "auth.Handler.Login", err, start)
            utils.RespondWithAPIError(c, apiErr) // e.g. 401 invalid credentials
            return
        }
        logger.LogFinish(ctx, "auth.Handler.Login", err, start)
        utils.InternalServerError(c, err, "Failed to authenticate user")
        return
    }

    logger.LogFinish(ctx, "auth.Handler.Login", nil, start)
    utils.Ok(c, response, "Login successful")
}
```

**✅ Service depends on the `Repository` interface, returns sentinel errors:**
```go
// internal/modules/auth/service.go
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) {
    ctx, start := logger.LogStart(ctx, "auth.Service.Login")

    user, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        logger.LogFinish(ctx, "auth.Service.Login", err, start)
        return nil, fmt.Errorf("authentication failed: %w", err)
    }
    if user == nil {
        logger.LogFinish(ctx, "auth.Service.Login", ErrInvalidCredentials, start)
        return nil, ErrInvalidCredentials
    }
    if err = s.verifyPassword(user.Password, req.Password); err != nil {
        logger.LogFinish(ctx, "auth.Service.Login", ErrInvalidCredentials, start)
        return nil, ErrInvalidCredentials
    }

    accessToken, _ := s.generateToken(user)
    refreshToken, _ := s.generateRefreshToken()
    user.RefreshToken = refreshToken
    _ = s.userRepo.UpdateUser(user)

    logger.LogFinish(ctx, "auth.Service.Login", nil, start)
    return &AuthResponse{
        User:         UserResponse{ID: user.ID, Name: user.Name, Email: user.Email},
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
    }, nil
}
```

**✅ The module exposes its guard and contract to other modules:**
```go
// internal/modules/auth/module.go
func (m *Module) Middleware() gin.HandlerFunc { return authMiddleware(m.svc) } // JWT guard
func (m *Module) Auth() Servicer              { return m.svc }                 // token validation, etc.
```

### 12.2 DataTable Pattern (from the `example` module)

The example module shows server-side DataTables: the repository runs the query, the service wraps
it in a trace span, and the handler renders it.

```go
// internal/modules/example/repository.go
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

// internal/modules/example/handler.go
func (h *Handler) Datatables(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "example.Handler.Datatables")
    data, err := h.svc.Datatables(ctx, c)
    if err != nil {
        logger.LogFinish(ctx, "example.Handler.Datatables", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }
    dt, ok := data.(dto.Datatables)
    if !ok {
        utils.InternalServerError(c, nil, "Failed to render data")
        return
    }
    logger.LogFinish(ctx, "example.Handler.Datatables", nil, start)
    datatables.JSON(c, dt)
}
```

Routes are mounted by the module: `api.GET("/examples", ...)` and `api.GET("/datatables", ...)`.

---

## 13. ANTI-PATTERNS TO AVOID

### 13.1 ❌ Business Logic in Handler

**WRONG:**
```go
func (h *Handler) Register(c *gin.Context) {
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
func (h *Handler) Register(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "auth.Handler.Register")
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        logger.LogFinish(ctx, "auth.Handler.Register", err, start)
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
    logger.LogFinish(ctx, "auth.Handler.Register", nil, start)
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

### 13.3 ❌ Standalone Handler Functions

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
// ✅ Struct handler with an injected, consumer-defined service interface
type Handler struct{ service Servicer }
func NewHandler(service Servicer) *Handler { return &Handler{service: service} }
func (h *Handler) Login(c *gin.Context) { /* ... */ }

// Wired in module.go, mounted in RegisterRoutes:
g := api.Group("/auth")
g.POST("/login", m.handler.Login)
```

---

### 13.4 ❌ God Service (Too Many Responsibilities)

**WRONG:**
```go
// ❌ One service (or one module) doing everything
type AppService struct{}
func (s *AppService) CreateUser() error      { /* ... */ }
func (s *AppService) ProcessPayment() error  { /* ... */ }
func (s *AppService) SendEmail() error       { /* ... */ }
func (s *AppService) GenerateReport() error  { /* ... */ }
// ... 50 more methods
```

**CORRECT:**
```go
// ✅ Separate modules, each a vertical slice with a single responsibility
internal/modules/auth/        // Service: authentication
internal/modules/payments/    // Service: payment processing
internal/modules/reports/     // Service: reporting
```
Each module talks to others only through their public `API` / `Servicer` interface (see §13.5).

---

### 13.5 ❌ Reaching Into Another Module's Internals

**WRONG:**
```go
// ❌ Importing another module's unexported types / repository directly
import "github.com/0xdiaz/gin-boilerplate/internal/modules/auth"

func (s *Service) doThing() {
    repo := auth.NewRepository(s.db)   // ❌ using auth's data layer from outside
    user, _ := repo.GetUserByEmail("x@y.com")
    _ = user
}
```

**CORRECT:**
```go
// ✅ Depend on auth's PUBLIC interface, injected through the constructor in buildModules()
type Service struct {
    auth auth.Servicer // public contract only
}
func NewService(repo repository, auth auth.Servicer) *Service {
    return &Service{ /* repo, */ auth: auth}
}

// internal/bootstrap/modules.go
authMod := auth.New(db)
mine := mymodule.New(db, authMod.Auth())
```
If a module must later become its own service, this in-process interface call becomes a network
(REST + HMAC) call and nothing else changes.

---

## 🎯 SUMMARY: Key Takeaways

### ✅ DO:
1. **Organize by module** — one feature = one folder = one Go package (`internal/modules/<name>/`).
2. **Keep the vertical slice together** — `handler → service → repository → model` in one place.
3. **Define interfaces at the consumer** — handler defines `service`, service defines `repository`.
4. **Inject dependencies via `New(...)`** — assemble the slice once in `Module.New(db)`.
5. **Cross modules only via public interfaces** (`Module.API()`, `auth.Servicer`), injected.
6. **Hold the injected `*gorm.DB`** in the repository — never the `database.DB` global.
7. **Use response utilities** (`pkg/utils`) — never hand-write `c.JSON`.
8. **Trace every operation** — `logger.LogStart` / `logger.LogFinish` around service & handler work.
9. **Map sentinel errors → APIError** in the handler; wrap unexpected errors with `%w`.
10. **Co-locate tests** — fake repo, no DB (`example/service_test.go`).

### ❌ DON'T:
1. **Business logic in handlers** — move it to the service.
2. **Database access in services** — go through the repository.
3. **Standalone handler functions** — use struct handlers with injected interfaces.
4. **God services/modules** — split by sub-domain.
5. **Reach into another module's internals** — use its public interface.
6. **Use the `database.DB` global in modules** — inject `*gorm.DB`.
7. **Put the auth guard in `pkg/middleware`** — it lives in the auth module (`Module.Middleware()`).
8. **`pkg/` importing `internal/`** — the shared kit must stay clean.
9. **Skip tests** — co-locate a fake-backed test with every module.
10. **Large files/functions** — split per the size limits in §5.

---

**END OF DESIGN PATTERNS DOCUMENT**

*Last Updated: 2026-06-10*
*Version: 2.0 (modular architecture)*
*Source of truth for layout: [`MODULE_GUIDE.md`](./MODULE_GUIDE.md). Reference module: `internal/modules/example`.*

---

## 📚 Related Documentation

- **[MODULE_GUIDE.md](./MODULE_GUIDE.md)** — Source of truth for the modular layout
- **CODING_STANDARDS.md** — Detailed coding standards and guidelines
- **AI_AGENT_RULES.md** — Quick reference rules for AI agents
- **CONTRACTS.md** — Stable API & configuration contracts
- **README.md** — Project overview and setup instructions
