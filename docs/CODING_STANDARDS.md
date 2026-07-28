# Coding Standards & Guidelines

**Version:** 2.0
**Last Updated:** 2026-06-10
**Language:** Go 1.24+

> **Architecture:** This service is organized **by business module** (package-by-feature),
> not by technical layer. The source of truth for the layout is
> [`MODULE_GUIDE.md`](./MODULE_GUIDE.md); this document gives the coding standards that
> apply *inside* that layout. Where older docs describe an `internal/app/...` +
> `internal/domain/...` layered tree, MODULE_GUIDE (and this doc) override them.

---

## ⚠️ FOR AI AGENTS - READ THIS FIRST

> **🚨 CRITICAL: Read [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) and
> [`MODULE_GUIDE.md`](./MODULE_GUIDE.md) BEFORE reading this document!**
>
> The critical rules file is short and contains NON-NEGOTIABLE patterns; MODULE_GUIDE
> describes the folder layout. This document is the detailed reference.
>
> **Don't make the same mistake:** Skipping critical sections = Code rejected.

### 🔥 Critical Sections (MUST READ)

After reading `00_AI_CRITICAL_RULES.md` and `MODULE_GUIDE.md`, focus on these sections:

- **[1. File Organization](#1-file-organization)** — the modular (package-by-feature) layout
- **[3. Code Structure](#3-code-structure)** — handler → service → repository, consumer-defined interfaces, cross-module calls
- **[7. Testing Requirements](#7-testing-requirements)** — co-located tests with a fake repo
- **[11.4 Response Utilities](#114-response-utilities)** — `pkg/utils` (MANDATORY)
- **[11.6 Routing](#116-routing-and-module-registration)** — modules register their own routes

### 📖 How to Use This Document

1. ✅ Read `00_AI_CRITICAL_RULES.md` and `MODULE_GUIDE.md` first
2. ✅ Read the critical sections listed above
3. ⚠️  Skim other sections for context
4. 📚 Use as reference when needed

---

## 📋 Table of Contents

### 🔥 Critical Sections (Read First)

- [1. File Organization](#1-file-organization) — `module`, `vertical slice`, `300 lines`
  - [1.5 DTO & Constant Placement](#15-dto-and-constant-placement-rules) — `dto.go`, module-local
- [3. Code Structure](#3-code-structure) — `handler`, `service`, `repository`, `interface`
- [5. Error Handling](#5-error-handling) — `fmt.Errorf`, `%w`, `domain errors`, `APIError`
- [7. Testing Requirements](#7-testing-requirements) — `co-located`, `fake repo`, `70%`, `coverage`
- [11. API Design](#11-api-design) — `REST`, `response`, `utils.Ok`, `routes`
  - [11.4 Response Utilities](#114-response-utilities) — `utils.Ok`, `utils.Created`, `utils.BadRequest`
  - [11.6 Routing & Module Registration](#116-routing-and-module-registration) — `RegisterRoutes`, `buildModules`

### 📚 All Sections

| # | Section | Keywords |
|---|---------|----------|
| 1 | [File Organization](#1-file-organization) | `file size`, `module`, `dto`, `vertical slice` |
| 1.1 | [File Size Limits](#11-file-size-limits) | `300 lines`, `100 lines function` |
| 1.2 | [Directory Structure](#12-directory-structure) | `internal/modules`, `pkg/`, `bootstrap` |
| 1.3 | [File Naming](#13-file-naming) | `snake_case`, `_test.go`, `naming` |
| 1.4 | [Package Organization](#14-package-organization) | `package`, `one module = one package` |
| 1.5 | [DTO & Constant Placement](#15-dto-and-constant-placement-rules) | `dto`, `constants`, `module-local` |
| 2 | [Naming Conventions](#2-naming-conventions) | `camelCase`, `PascalCase`, `naming` |
| 3 | [Code Structure](#3-code-structure) | `handler`, `service`, `repository`, `interface` |
| 4 | [Function Guidelines](#4-function-guidelines) | `function size`, `parameters`, `return` |
| 5 | [Error Handling](#5-error-handling) | `error`, `wrapping`, `domain errors`, `APIError` |
| 6 | [Documentation](#6-documentation) | `comment`, `godoc`, `docs` |
| 7 | [Testing Requirements](#7-testing-requirements) | `test`, `co-located`, `fake repo`, `70%` |
| 8 | [Logging Standards](#8-logging-standards) | `logger`, `LogStart`, `LogFinish`, `request_id` |
| 9 | [Security Guidelines](#9-security-guidelines) | `validation`, `SQL injection`, `bcrypt` |
| 10 | [Database Access](#10-database-access) | `GORM`, `model`, `migration`, `DI` |
| 11 | [API Design](#11-api-design) | `REST`, `API`, `response`, `routes` |
| 11.1 | [RESTful Conventions](#111-restful-conventions) | `REST`, `HTTP methods`, `endpoints` |
| 11.2 | [HTTP Status Codes](#112-http-status-codes) | `200`, `404`, `500`, `status` |
| 11.3 | [Response Format](#113-response-format) | `success`, `message`, `data`, `errors` |
| 11.4 | [Response Utilities](#114-response-utilities) | `utils.Ok`, `utils.Created`, `RespondWithAPIError` |
| 11.5 | [API Versioning](#115-api-versioning) | `v1`, `v2`, `versioning` |
| 11.6 | [Routing & Module Registration](#116-routing-and-module-registration) | `RegisterRoutes`, `buildModules` |
| 12 | [Configuration](#12-configuration) | `config`, `env`, `.env` |
| 13 | [Forbidden Practices](#13-forbidden-practices) | `panic`, `global`, `anti-pattern` |
| 14 | [Code Review Checklist](#14-code-review-checklist) | `checklist`, `review`, `validation` |
| 15 | [AI Agent Rules](#15-ai-agent-specific-rules) | `AI`, `agent`, `automation` |
| 16 | [Resources](#16-resources) | `links`, `references`, `docs` |

### 🎯 Quick Lookups by Task

**Creating a new module (handler/service/repository):**
- [3. Code Structure](#3-code-structure) — consumer-defined interfaces + constructor injection
- [11.4 Response Utilities](#114-response-utilities) — MUST use
- See: `MODULE_GUIDE.md` — "How to add a new module"

**Adding routes:**
- [11.6 Routing & Module Registration](#116-routing-and-module-registration)
- [11.1 RESTful Conventions](#111-restful-conventions)

**Writing tests:**
- [7. Testing Requirements](#7-testing-requirements)
- [1.1 File Size Limits](#11-file-size-limits)

**Error handling:**
- [5. Error Handling](#5-error-handling)
- [8. Logging Standards](#8-logging-standards)

**Database models:**
- [10. Database Access](#10-database-access)
- See: `MODULE_GUIDE.md` — "A module owns its tables"

---

## 1. FILE ORGANIZATION

### 1.1 File Size Limits

```
✅ MUST: Maximum 300 lines per file
⚠️  WARNING: 300-500 lines requires justification
❌ FORBIDDEN: >500 lines in a single file
```

**Action when exceeding:**
- Split into multiple focused files **within the same module package**.
- **No subfolders.** When one part of a module (service, handler, or repository) grows past the limit, add another file in the same module folder, keeping the same package. Example: the auth module splits its service into `internal/modules/auth/service.go` + `internal/modules/auth/service_tokens.go` — both `package auth`, no nested folder.
- The module *is* the feature folder; a split file does not get its own subfolder.

### 1.2 Directory Structure

Code is organized **by business module** (package-by-feature), not by technical layer. A
module owns its full vertical slice (`handler → service → repository → model`) in one
package. See [`MODULE_GUIDE.md`](./MODULE_GUIDE.md) for the canonical version.

**MUST follow this structure:**

```
.
├── main.go                      # 3 lines: bootstrap.Run()
├── internal/
│   ├── bootstrap/               # the ONLY per-service wiring
│   │   ├── bootstrap.go         #   Run(): config → db → modules → migrate → serve → shutdown
│   │   ├── modules.go           #   Module interface + buildModules() (the list of modules)
│   │   ├── server.go            #   gin engine + global middleware + route mounting
│   │   └── swagger.go           #   OpenAPI/Swagger UI (debug only)
│   ├── migrations/              # service-specific schema
│   │   ├── migration.go         #   Run(db, models) — AutoMigrate
│   │   └── sql/                 #   versioned SQL (production path)
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
├── pkg/                         # ← cross-cutting "kit"; MUST never import from internal/
│   ├── config/                  #   configuration (viper-backed)
│   ├── logger/                  #   structured logging + request-scoped tracing
│   ├── metrics/                 #   request counters + uptime
│   ├── types/                   #   shared response/error types
│   ├── utils/                   #   response helpers (utils.Ok, utils.BadRequest, …)
│   ├── middleware/              #   cors, request_id, request_log, metrics, rate_limit
│   └── database/                #   connect + read-replica resolver
└── internal/testsupport/        # shared test doubles/mocks imported by co-located tests
```

> **Testing layout:** unit tests are **co-located** with the code they test
> (`foo_test.go` beside `foo.go`, default `package <pkg>_test`); test fixtures live in a
> `testdata/` folder next to the test that uses them (Go convention). Shared mocks reusable
> across packages live in `internal/testsupport/` (`package testsupport`); a mock used by
> only one package is co-located there instead. There is **no root `tests/` tree** — DB /
> external-service integration tests belong in the cross-service harness `paprika-testing`,
> not in this repo.

> **Old → new (do not use the old paths):** `internal/app/{controllers,services,dto,middlewares,routers}`
> and `internal/domain/{models,repositories}` and `internal/adapters/database` no longer exist.
> Everything for a feature lives in `internal/modules/<name>/`; wiring lives in
> `internal/bootstrap/`; migrations in `internal/migrations/`; cross-cutting code in `pkg/`
> (note: `pkg/middleware` is **singular**).

#### 1.2.1 One Module = One Package = One Folder

```
✅ MUST: A feature is a single package under internal/modules/<name>/ holding
   model.go, dto.go (optional), repository.go, service.go, handler.go, module.go.
   This is a vertical slice — do NOT spread a feature across layer folders.
```

**Purpose:**
- One person (or one AI agent) can own a module without touching others.
- Navigation, code ownership, and the future "lift this module into its own service" cut are all trivial.
- A module owns its tables (declared in `Module.Models()`); no cross-module foreign keys.

**Example (the canonical `example` module):**
```
internal/modules/example/
├── model.go         # Example (GORM model) + TableName()
├── repository.go    # Repository{db *gorm.DB} + NewRepository(db)
├── service.go       # Service + the `repository` interface it consumes
├── handler.go       # Handler + the `service` interface it consumes
├── module.go        # Module: New(db), Name(), Models(), RegisterRoutes(api), API()
└── service_test.go  # co-located unit test using a fake repo (no DB)
```

**Naming Rules:**
- Folder name *is* the module/sub-domain: `auth/`, `members/`, `accounts/`, `balance/`.
- All files use `snake_case`; every file in the folder shares one package (`package <name>`).
- No nested feature subfolders inside a module — split into more files in the same folder instead.

### 1.3 File Naming

**Rules:**
```go
✅ CORRECT:
- service.go               (snake_case; the module's business logic)
- repository.go            (snake_case)
- service_tokens.go        (snake_case; a second service file in the same module)
- service_test.go          (co-located test)

❌ WRONG:
- Service.go               (PascalCase not allowed)
- service-tokens.go        (kebab-case not allowed)
- serviceTokens.go         (camelCase not allowed)
```

Inside a module the standard files are `model.go`, `dto.go`, `repository.go`,
`service.go`, `handler.go`, `module.go`. When a file outgrows the size limit, append a
descriptive suffix (`service_tokens.go`, `repository_query.go`) in the same folder.

### 1.4 Package Organization

**One module = one package = one directory:**
```go
// ✅ CORRECT — every file in the module shares the package
internal/modules/auth/
    ├── service.go         → package auth
    ├── service_tokens.go  → package auth
    ├── repository.go      → package auth
    └── handler.go         → package auth

// ❌ WRONG - multiple packages in one directory, or layer folders
internal/modules/auth/
    ├── service.go         → package service
    └── handler.go         → package handler
```

### 1.5 DTO and Constant Placement Rules

DTOs and constants are **module-local**. There is no shared `dto/` package and no
`pkg/enums/` package — a module owns its own request/response types and its own constants.

#### 1.5.1 DTOs (Data Transfer Objects)
```
✅ DTOs LIVE IN THE MODULE'S dto.go (package <module>)

Location rules:
- Request/Response structs → internal/modules/<name>/dto.go
- Any struct used for data transfer at the HTTP boundary → dto.go
- dto.go is optional: tiny modules may declare a request/response type in model.go

Examples:
✅ CORRECT:
internal/modules/auth/dto.go     → RegisterRequest, LoginRequest, AuthResponse, UserResponse
internal/modules/health/dto.go   → HealthResponse, MetricsResponse

❌ WRONG:
internal/app/dto/user_dto.go     → (the internal/app tree no longer exists)
pkg/types/request.go             → ApiRequest (pkg/types is for shared response/error types only)
```

DTOs carry their own validation via gin **binding tags** (see [§9.1](#91-input-validation)):

```go
// internal/modules/auth/dto.go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

#### 1.5.2 Constants and Enum-like Values
```
✅ CONSTANTS LIVE WITH THE MODULE THAT OWNS THEM

Location rules:
- A module's status/role/type constants → declared in that module's package
- A truly cross-cutting constant (used by pkg/ and multiple modules) → pkg/ (e.g. pkg/types)

There is NO pkg/enums package. Do not invent one. Keep a constant next to the code that
gives it meaning, in the module package, so it stays cohesive and testable.

Examples:
✅ CORRECT:
internal/modules/auth/service.go  → ErrInvalidCredentials, ErrEmailAlreadyExists (domain errors)
pkg/types/errors.go               → APIError, ErrNotFound, ErrUnauthorized (shared HTTP mapping)

❌ WRONG:
pkg/enums/role.go                 → (no such package)
```

#### 1.5.3 Import Direction (no cycles)

**CRITICAL: `pkg/` MUST never import from `internal/`. Modules MUST never import another
module's internals.**

```go
✅ ALLOWED:
internal/modules/<name> → pkg/utils, pkg/types, pkg/logger, pkg/config   ✓
internal/bootstrap      → internal/modules/<name>                         ✓
one module → another module's PUBLIC interface (injected via constructor) ✓

❌ FORBIDDEN (cause import cycles or break the layering):
pkg/anything → internal/...                                              ✗
internal/modules/a → internal/modules/b's unexported types               ✗
```

Cross-module calls go through the consumed module's exported interface (e.g.
`auth.Module.Auth()` returns `auth.Servicer`, `example.Module.API()` returns
`example.API`), injected in `buildModules()`. See [§3.5](#35-cross-module-communication).

---

## 2. NAMING CONVENTIONS

### 2.1 Package Names

```go
✅ MUST:
- Lowercase, single word
- No underscores, no dashes
- The package is named after the module/sub-domain it represents

✅ CORRECT:
package auth        // internal/modules/auth
package example     // internal/modules/example
package health      // internal/modules/health
package utils       // pkg/utils

❌ WRONG:
package user_auth        // No underscores
package Auth             // No capitals
package svc              // Too cryptic
package controllers      // Layer-named packages no longer exist
```

### 2.2 Variable Names

**Short names for short scopes:**
```go
✅ CORRECT:
// Short scope (1-5 lines)
for i, v := range items {
    fmt.Println(i, v)
}

// Medium scope (5-20 lines)
user := &models.User{}
client := repository.GetClient(id)

// Long scope or package-level
transactionRepository := NewTransactionRepository(db)
httpClientTimeout := 30 * time.Second

❌ WRONG:
u := &models.User{}              // Too short for long scope
transRepo := repo.GetTrans()     // Unclear abbreviation
HTTPClientTimeout := 30          // Unexported shouldn't be capitalized
```

### 2.3 Function/Method Names

```go
✅ CORRECT:
// Exported (public)
func GetUserByID(id uint) (*models.User, error)
func CreateTransaction(tx *models.Transaction) error
func ValidateHTTPRequest(req *http.Request) error

// Unexported (private)
func parseRequestBody(body []byte) (map[string]interface{}, error)
func sanitizeLogData(data string) string
func calculateTotalFee(amount float64) float64

❌ WRONG:
func get_user(id uint)                    // Snake_case not allowed
func GetUser(id uint)                     // Too generic (which user?)
func GU(id uint)                          // Too cryptic
func HTTPGETRequest()                     // Redundant "HTTP GET"
func GetUserByIdFromDatabase(id uint)     // Too verbose
```

### 2.4 Constants and Enums

```go
✅ CORRECT:
const (
    // Single constant - PascalCase
    MaxRetryAttempts = 3
    DefaultTimeout   = 30 * time.Second

    // Enum-like constants - Prefix with type
    StatusPending    Status = "pending"
    StatusProcessing Status = "processing"
    StatusCompleted  Status = "completed"

    // Private constants
    defaultPageSize = 10
    maxPageSize     = 100
)

❌ WRONG:
const MAX_RETRY = 3              // SCREAMING_SNAKE_CASE not idiomatic
const max_retry = 3              // snake_case not idiomatic
const Pending = "pending"        // Missing type prefix
```

### 2.5 Struct Names

```go
✅ CORRECT:
type User struct { }
type HTTPClient struct { }       // Acronyms in PascalCase: HTTP not Http
type APIResponse struct { }      // API not Api
type TransactionDTO struct { }   // Clear purpose with suffix

❌ WRONG:
type user struct { }             // Lowercase for exported
type HttpClient struct { }       // Should be HTTPClient
type TransactionDataTransferObject struct { }  // Too verbose
```

---

## 3. CODE STRUCTURE

Each module owns its full vertical slice. Inside a module the flow is
`handler → service → repository → model`. The HTTP layer is the **Handler** (type
`Handler`, constructor `NewHandler`), not a "controller".

### 3.1 The Vertical Slice (handler → service → repository → model)

**MUST follow this flow inside a module:**

```
HTTP Request
    ↓
[Module.RegisterRoutes] → mounts the route on the /api/v1 group
    ↓
[Handler] → thin HTTP layer (parse/bind request, call service, write response via pkg/utils)
    ↓
[Service] → business logic (validation, orchestration, domain errors)
    ↓
[Repository] → data access (CRUD / queries on the injected *gorm.DB)
    ↓
[Model] → GORM entity the module owns
```

**Rules:**
```go
✅ Handler SHOULD:
- Bind/parse the request (c.ShouldBindJSON)
- Call service methods (passing the request context)
- Write responses via pkg/utils (utils.Ok, utils.Created, utils.BadRequest, …)
- Map domain errors to HTTP (see §5.5)

❌ Handler MUST NOT:
- Contain business logic
- Access the database directly
- Call c.JSON directly for success/error (use pkg/utils)

✅ Service SHOULD:
- Contain all business logic and validation of business rules
- Define the `repository` interface it needs (consumer-defined; see §3.3)
- Return domain errors (e.g. auth.ErrInvalidCredentials)
- Take context.Context as its first parameter and use logger.LogStart/LogFinish

❌ Service MUST NOT:
- Handle HTTP concerns (gin.Context) — EXCEPTION: a third-party lib that requires it,
  e.g. DataTables binding (see §11.4)
- Import another module's internals

✅ Repository SHOULD:
- Perform CRUD / build queries on its injected *gorm.DB
- Be constructed with NewRepository(db) — no global DB access
- Return the module's models (or errors)

❌ Repository MUST NOT:
- Contain business logic
- Reach into another module's tables
```

### 3.2 Dependency Direction

**MUST follow (within a module):**
```
Handler → Service → Repository → Model
   ↓         ↓          ↓
 pkg/utils  pkg/logger  *gorm.DB
```

**FORBIDDEN:**
```
❌ Model importing Service
❌ Repository importing Service
❌ Service importing Handler
❌ A module importing another module's unexported types
❌ pkg/ importing internal/
❌ Circular dependencies
```

### 3.3 Interfaces are Consumer-Defined (DI + testability)

**MANDATORY: the consumer declares the interface it depends on, and receives the
implementation through its constructor.** This is what makes a module unit-testable with a
fake (no DB) and keeps coupling loose.

- The **service** defines the `repository` interface it needs (in `service.go`).
- The **handler** defines the `service` interface it needs (in `handler.go`).
- `module.go` wires the concrete types together via constructors.

```go
// internal/modules/example/service.go
// repository is the data-access contract THIS service needs (defined here, at the consumer).
type repository interface {
    List() ([]*Example, error)
    Datatables(c *gin.Context) (interface{}, error)
}

type Service struct {
    repo repository
}

func NewService(repo repository) *Service { return &Service{repo: repo} }
```

```go
// internal/modules/example/handler.go
// service is the business contract THIS handler needs (consumer-defined for testability).
type service interface {
    List(ctx context.Context) ([]*Example, error)
    Datatables(ctx context.Context, c *gin.Context) (interface{}, error)
}

type Handler struct {
    svc service
}

func NewHandler(svc service) *Handler { return &Handler{svc: svc} }
```

```go
// internal/modules/example/module.go — wires the concrete graph
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))
    return &Module{svc: svc, handler: NewHandler(svc)}
}
```

> **No standalone business functions.** Business logic lives on a `Service` method, behind
> a constructor-injected dependency — never as a free function operating on globals.

### 3.4 Single Responsibility Principle

**Each module owns one sub-domain; each file has ONE clear purpose:**

```go
✅ CORRECT:
// internal/modules/auth/service.go — auth business logic only
type Service struct {
    userRepo Repository    // interface
    mailer   EmailSender   // optional dependency, injected
}

func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
func (s *Service) ValidateToken(token string) (uint, error)

❌ WRONG:
// a "god" service mixing unrelated sub-domains
type GodService struct {
    userRepo    Repository
    paymentRepo PaymentRepository
    reportRepo  ReportRepository
}

func (s *GodService) Register() error
func (s *GodService) ProcessPayment() error
func (s *GodService) GenerateReport() error
// → these belong in separate modules (auth / payments / reporting)
```

### 3.5 Cross-module Communication

**Modules talk to each other only through a public interface, injected via the
constructor — never by importing another module's guts.** This keeps modules loosely
coupled and makes the seam easy to cut later (an in-process call becomes a REST + HMAC call
with nothing else changing).

```go
// auth exposes its public surface:
authMod := auth.New(db)
authMod.Middleware()   // gin.HandlerFunc — protect routes in other modules
authMod.Auth()         // auth.Servicer — ValidateToken, etc.

// A module that needs auth receives it injected in buildModules():
//   payments.New(db, authMod.Auth())   // payments depends on auth.Servicer, not *auth.Service
```

```go
// example exposes a minimal public contract for other modules:
// internal/modules/example/module.go
type API interface {
    List(ctx context.Context) ([]*Example, error)
}
func (m *Module) API() API { return m.svc }
```

**Rules:**
- Expose the smallest interface other modules need (`auth.Servicer`, `example.API`).
- Depend on the interface, not the concrete `*Service`.
- Wire the dependency in `buildModules()` (`internal/bootstrap/modules.go`) — the one place
  that knows the concrete module list.

---

## 4. FUNCTION GUIDELINES

### 4.1 Function Size Limits

```
✅ IDEAL: 20-50 lines per function
⚠️  ACCEPTABLE: 50-100 lines with justification
❌ FORBIDDEN: >100 lines per function
```

**If function exceeds 100 lines, MUST refactor:**

```go
❌ BEFORE (150 lines):
func CreateTransactionFromApiLog(log *models.ApiLog) error {
    // 150 lines of code doing everything
    // parsing, validation, transformation, saving, logging
}

✅ AFTER (split into focused functions):
func CreateTransactionFromApiLog(log *models.ApiLog) error {
    // Orchestration only - 15 lines
    reqData := parseRequestData(log.RequestBody)
    respData := parseResponseData(log.ResponseBody)

    tx := buildTransaction(reqData, respData)
    if err := validateTransaction(tx); err != nil {
        return err
    }

    return saveTransaction(tx)
}

func parseRequestData(body string) map[string]interface{} { }      // 20 lines
func parseResponseData(body string) map[string]interface{} { }     // 20 lines
func buildTransaction(req, resp map[string]interface{}) *Transaction { } // 30 lines
func validateTransaction(tx *Transaction) error { }                // 25 lines
func saveTransaction(tx *Transaction) error { }                    // 15 lines
```

### 4.2 Function Parameters

```go
✅ MAXIMUM: 4 parameters per function
⚠️  WARNING: 5-7 parameters (consider refactoring)
❌ FORBIDDEN: >7 parameters

// ❌ WRONG - Too many parameters
func CreateUser(
    name, email, phone, address, city, province, postalCode string,
    age int,
    isActive bool,
) error

// ✅ CORRECT - Use struct
type CreateUserParams struct {
    Name       string
    Email      string
    Phone      string
    Address    string
    City       string
    Province   string
    PostalCode string
    Age        int
    IsActive   bool
}

func CreateUser(params CreateUserParams) error
```

### 4.3 Return Values

```go
✅ CORRECT:
// Return value + error
func GetUser(id uint) (*models.User, error)

// Return multiple values (max 3)
func ParseTransaction(data string) (amount float64, currency string, err error)

// Return only error for void operations
func DeleteUser(id uint) error

❌ WRONG:
// Returning more than 3 values
func GetUserDetails(id uint) (string, string, string, int, bool, error)

// Not returning error when operation can fail
func GetUser(id uint) *models.User

// Returning error as first parameter
func GetUser(id uint) (error, *models.User)
```

### 4.4 Function Ordering in Files

**MUST follow this order:**

```go
// internal/modules/auth/service.go (package auth)

// 1. Consumer-defined interface(s) this file needs (see §3.3)
type repository interface {
    GetUserByEmail(email string) (*User, error)
    CreateUser(user *User) error
}

// 2. Type definition
type Service struct {
    userRepo Repository
    mailer   EmailSender
}

// 3. Constructor
func NewService(userRepo Repository, mailer EmailSender) *Service {
    return &Service{userRepo: userRepo, mailer: mailer}
}

// 4. Public methods (exported) — grouped by use case
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) { }
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error) { }
func (s *Service) ValidateToken(token string) (uint, error) { }

// 5. Private methods (unexported)
func (s *Service) hashPassword(password string) (string, error) { }
func (s *Service) verifyPassword(hashed, plain string) error { }
```

---

## 5. ERROR HANDLING

### 5.1 Error Handling Pattern

**MUST use this pattern:**

```go
✅ CORRECT:
func GetUser(id uint) (*models.User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        logger.Errorf("failed to get user %d: %v", id, err)
        return nil, fmt.Errorf("get user: %w", err)  // Wrap error
    }

    return user, nil
}

❌ WRONG:
func GetUser(id uint) (*models.User, error) {
    user, err := repo.FindByID(id)
    if err != nil {
        return nil, err  // Not wrapped, no context
    }
    return user, nil
}

❌ WRONG:
func GetUser(id uint) *models.User {
    user, _ := repo.FindByID(id)  // Error ignored!
    return user
}

❌ FORBIDDEN:
func GetUser(id uint) *models.User {
    user, err := repo.FindByID(id)
    if err != nil {
        panic(err)  // NEVER panic in business logic!
    }
    return user
}
```

### 5.2 Error Wrapping

**Always wrap errors with context:**

```go
✅ CORRECT:
import "fmt"

// Add context with %w (Go 1.13+)
if err != nil {
    return fmt.Errorf("failed to create transaction: %w", err)
}

// Multiple context layers
if err := service.CreateUser(user); err != nil {
    return fmt.Errorf("user registration failed for email %s: %w", user.Email, err)
}

❌ WRONG:
if err != nil {
    return err  // No context
}

if err != nil {
    return fmt.Errorf("failed to create transaction: %v", err)  // Use %w not %v
}
```

### 5.3 Domain Errors (module-local)

**Define a module's domain errors in that module's package**, as sentinel `error` values.
The handler maps them to HTTP (see §5.5). Shared, HTTP-shaped errors (`*types.APIError`)
live in `pkg/types/errors.go`.

```go
✅ CORRECT:
// internal/modules/auth/service.go (package auth)
import "errors"

var (
    ErrEmailAlreadyExists  = errors.New("email already exists")
    ErrInvalidCredentials  = errors.New("invalid email or password")
    ErrUserNotFound        = errors.New("user not found")
    ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
)

// Return a domain error from the service
if existingUser != nil {
    return nil, ErrEmailAlreadyExists
}

// Wrap with context for non-sentinel failures
if err := s.userRepo.CreateUser(user); err != nil {
    return nil, fmt.Errorf("failed to create user: %w", err)
}
```

```go
// pkg/types/errors.go — shared HTTP-shaped errors (used when mapping to a response)
var (
    ErrNotFound       = NewAPIError(http.StatusNotFound, "Resource not found", nil)
    ErrUnauthorized   = NewAPIError(http.StatusUnauthorized, "Unauthorized access", nil)
    ErrInvalidInput   = NewAPIError(http.StatusBadRequest, "Invalid input data", nil)
    ErrConflict       = NewAPIError(http.StatusConflict, "Resource already exists", nil)
    ErrInternalServer = NewAPIError(http.StatusInternalServerError, "Internal server error", nil)
)
```

### 5.4 Panic Usage

```
❌ FORBIDDEN in business logic (handlers, services, repositories)
⚠️  ACCEPTABLE only in:
    - Application startup (bootstrap.Run via logger.Fatalf)
    - Fatal unrecoverable startup errors (config load / DB connect failure)

✅ The gin engine installs gin.Recovery() in bootstrap.buildEngine(), so panics in a
   handler are caught and turned into 500 — but you MUST NOT rely on it: return errors.
```

```go
✅ CORRECT:
// internal/bootstrap/bootstrap.go — Fatalf is acceptable at startup
if err := config.SetupConfig(); err != nil {
    logger.Fatalf("config SetupConfig() error: %s", err)
}
if err := database.DbConnection(masterDSN, replicaDSN); err != nil {
    logger.Fatalf("database DbConnection error: %s", err)
}

❌ FORBIDDEN:
// service.go — NEVER panic in business logic
func (s *Service) Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error) {
    if req.Email == "" {
        panic("email is required")  // ❌ return an error instead
    }
    // ...
}
```

### 5.5 Error Strategy (HTTP Response and Domain Errors)

**HTTP responses:** All responses to the client (success and error) MUST go through
`pkg/utils`: `utils.Ok`, `utils.Created`, `utils.BadRequest`, `utils.Unauthorized`,
`utils.RespondWithAPIError`, etc. Do not call `c.JSON(code, ...)` directly.

**Domain errors:** Services return module-local domain errors (e.g.
`auth.ErrEmailAlreadyExists`, `auth.ErrInvalidCredentials`). The **handler** maps them to
HTTP using either:
- `utils.RespondWithAPIError(c, apiErr)` when the error is mapped to a `*types.APIError`
  (code, message, details), or
- `utils.InternalServerError(c, err, "…")` (and similar) for unmapped/generic errors.

The canonical pattern (from `internal/modules/auth/handler.go`) is a small mapping helper
plus `errors.Is`:

```go
// internal/modules/auth/handler.go
func errToAPIError(err error) *types.APIError {
    switch {
    case errors.Is(err, ErrEmailAlreadyExists):
        return &types.APIError{Code: http.StatusConflict, Message: "Email already exists"}
    case errors.Is(err, ErrInvalidCredentials):
        return &types.APIError{Code: http.StatusUnauthorized, Message: "Invalid email or password"}
    default:
        return nil
    }
}

func (h *Handler) Login(c *gin.Context) {
    // ... bind request ...
    response, err := h.service.Login(ctx, &req)
    if err != nil {
        if apiErr := errToAPIError(err); apiErr != nil {
            utils.RespondWithAPIError(c, apiErr)
            return
        }
        utils.InternalServerError(c, err, "Failed to authenticate user")
        return
    }
    utils.Ok(c, response, "Login successful")
}
```

**Validation errors:** Pass the bind error straight to `utils.BadRequest(c, err, "Invalid
request data")`. `utils.FormatValidationErrors` turns `validator.ValidationErrors` into a
field→message map automatically; no need to convert to `types.APIError`.

---

## 6. DOCUMENTATION

### 6.0 Comment Language Policy

```
✅ MUST: All code comments, package comments, and documentation MUST be written in English.
```

Rationale:
- Consistency across teams and open-source contribution
- Easier code review and adoption across regions
- Better processing by English-oriented tooling and AI

### 6.1 Package Documentation

**MUST include package comment in ONE file per package:**

```go
✅ CORRECT:
// Package auth implements the authentication module: registration, login, JWT
// issuance/validation, and password reset.
//
// It owns the users table and exposes a public surface to other modules via
// Module.Auth() (auth.Servicer) and Module.Middleware() (the JWT guard).
//
// Example usage (from bootstrap):
//
//	authMod := auth.New(db)
//	protected.Use(authMod.Middleware())
package auth
```

**File to add package comment:**
- Choose the module's main file (`module.go` or `service.go`) or create `doc.go`.

### 6.2 Function Documentation

**MUST document ALL exported functions:**

```go
✅ CORRECT:
// GetUserByID retrieves a user by their unique identifier.
//
// Returns ErrNotFound if the user does not exist.
// Returns ErrInvalidInput if id is 0.
func GetUserByID(id uint) (*models.User, error) {
    // Implementation
}

// CreateTransaction creates a new transaction record and processes payment.
//
// The function performs the following steps:
//  1. Validates transaction data
//  2. Checks client balance
//  3. Processes payment via external provider
//  4. Stores transaction record
//
// Returns the created transaction and nil error on success.
// Returns nil and error if any step fails (operation is rolled back).
func CreateTransaction(tx *models.Transaction) (*models.Transaction, error) {
    // Implementation
}

❌ WRONG:
// Get user
func GetUserByID(id uint) (*models.User, error) {  // Too brief

func CreateTransaction(tx *models.Transaction) error {  // No comment!
    // Implementation
}
```

### 6.3 Struct Documentation

**MUST document exported structs and important fields:**

```go
✅ CORRECT:
// Service handles authentication-related business logic.
type Service struct {
    userRepo Repository  // injected; see §3.3
    mailer   EmailSender // optional dependency
}

// RegisterRequest represents the payload for user registration.
type RegisterRequest struct {
    // Name is the full name of the user (required, 3-255 chars)
    Name string `json:"name" binding:"required,min=3,max=255"`

    // Email must be unique across all users (required, valid email format)
    Email string `json:"email" binding:"required,email"`

    // Password is the plaintext password to hash (required, min 8 chars)
    Password string `json:"password" binding:"required,min=8"`
}

❌ WRONG:
// Auth service
type Service struct {  // Too brief
    userRepo Repository
}

type RegisterRequest struct {  // No struct comment!
    Name     string `json:"name"`
    Email    string `json:"email"`
    Password string `json:"password"`
}
```

> **Validation uses gin binding tags** (`binding:"..."`), bound with `c.ShouldBindJSON`,
> not a separate `validate:"..."` + `validator.New()` call. See [§9.1](#91-input-validation).

### 6.4 Inline Comments

**Use sparingly for complex logic only (in English):**

```go
✅ CORRECT (complex logic needs explanation):
// Calculate fee with progressive rate:
// 0-1000: 1%, 1001-5000: 0.75%, >5000: 0.5%
var fee float64
if amount <= 1000 {
    fee = amount * 0.01
} else if amount <= 5000 {
    fee = 1000*0.01 + (amount-1000)*0.0075
} else {
    fee = 1000*0.01 + 4000*0.0075 + (amount-5000)*0.005
}

❌ WRONG (obvious code doesn't need comments):
// Loop through all users
for _, user := range users {
    // Print user name
    fmt.Println(user.Name)
}

// Set status to pending
status = StatusPending

// Add fee to total
total += fee
```

### 6.5 TODO Comments

```go
⚠️  ACCEPTABLE temporarily, but must have:
- Assignee
- Date
- Reason

✅ CORRECT:
// TODO(username, 2025-11-08): Implement caching for frequently accessed users
// to improve performance. Current response time is 200ms, target is <50ms.

❌ WRONG:
// TODO: fix this
// TODO: optimize
// FIXME
```

**MUST NOT commit TODOs older than 30 days** - convert to GitHub issues instead.

---

## 7. TESTING REQUIREMENTS

### 7.1 Test File Location — Prefer Co-located

**⚠️ CRITICAL: Prefer co-locating unit tests with the module they test, as
`internal/modules/<name>/*_test.go`.** Unit-test the service with a **fake repository**
(no DB). This is the canonical pattern; see `internal/modules/example/service_test.go`.

There is **no root `tests/` tree**. Prefer in-package fakes; a mock reused across packages
lives in `internal/testsupport/` (`package testsupport`) and is imported by the tests that
need it. Test fixtures live in a `testdata/` folder next to the test that uses them (Go
convention). DB / external-service **integration tests belong in the cross-service harness
`paprika-testing`**, not in this repo.

```
✅ CORRECT:
internal/modules/example/service_test.go   // co-located unit test, fake repo, no DB
internal/modules/auth/service_test.go      // co-located unit test
internal/testsupport/user_repo_mock.go     // shared mock, imported by co-located tests
internal/modules/example/testdata/…        // fixtures next to the test that reads them

❌ WRONG:
internal/app/services/auth_service_test.go  // the internal/app tree no longer exists
tests/integration/api/health_test.go        // no root tests/ tree; integration → paprika-testing
tests/mocks/…  tests/fixtures/…              // retired; use internal/testsupport + testdata/
```

**Package naming for co-located unit tests — white-box (`package <module>`):**

Co-located tests live in the **same package** so they can satisfy the module's
**unexported** consumer-defined interfaces with an in-package fake (e.g. the `repository`
interface the service depends on). This is the key enabler of DB-free unit tests.

```go
✅ CORRECT:
// internal/modules/example/service_test.go
package example  // same package — can implement the unexported `repository` interface

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

❌ WRONG:
// internal/modules/example/service_test.go
package example_test  // external package can't see the unexported repository interface
```

> For a public-API test that only needs the module's exported surface, the external
> `package <module>_test` (importing the real module path) is fine and is the default when
> white-box access to unexported interfaces is not required.

### 7.2 Test Coverage Requirements

```
✅ MINIMUM: 70% code coverage for the service (business logic) of each module
✅ MINIMUM: 50% code coverage for repositories
✅ MINIMUM: 60% code coverage for pkg/utils and other shared utilities
⚠️  Handlers: covered by handler tests (httptest + mocked service) and/or integration tests
```

### 7.3 Test File Naming

```go
✅ CORRECT:
service.go        → service_test.go
repository.go     → repository_test.go
service_tokens.go → service_tokens_test.go

❌ WRONG:
service.go → auth_test.go
service.go → test_service.go
```

### 7.4 Test Function Naming

```go
✅ CORRECT:
func TestCreateUser_Success(t *testing.T) { }
func TestCreateUser_DuplicateEmail(t *testing.T) { }
func TestCreateUser_InvalidInput(t *testing.T) { }

func TestGetUserByID_Found(t *testing.T) { }
func TestGetUserByID_NotFound(t *testing.T) { }

❌ WRONG:
func TestCreateUser(t *testing.T) { }  // Too generic
func Test_Create_User(t *testing.T) { }  // Wrong format
func createUserTest(t *testing.T) { }  // Must start with Test
```

### 7.5 Canonical Pattern — Co-located Service Test with a Fake Repo

**This is the pattern to copy** (from `internal/modules/example/service_test.go`). The fake
implements the module's **unexported `repository` interface**, so the service is tested with
no DB. Because the test is white-box (`package example`), it can see that interface.

```go
// internal/modules/example/service_test.go
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

**MUST use table-driven tests** when a method has multiple scenarios (success +
error cases), driving the fake's return values per case:

```go
func TestService_List_Cases(t *testing.T) {
    tests := []struct {
        name    string
        repo    *fakeRepo
        wantLen int
        wantErr bool
    }{
        {name: "returns items", repo: &fakeRepo{items: []*Example{{ID: 1, Data: "x"}}}, wantLen: 1},
        {name: "empty", repo: &fakeRepo{}, wantLen: 0},
        {name: "repo error", repo: &fakeRepo{err: assert.AnError}, wantErr: true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := NewService(tt.repo).List(context.Background())
            if (err != nil) != tt.wantErr {
                t.Fatalf("List() error = %v, wantErr %v", err, tt.wantErr)
            }
            assert.Len(t, got, tt.wantLen)
        })
    }
}
```

**Handler tests** use `httptest` + a mocked `service` interface (defined in `handler.go`).

### 7.6 Test Organization

- **Co-located unit tests** live next to the code: `internal/modules/<name>/*_test.go`.
- **`tests/`** holds integration tests, shared mocks, fixtures, and the testing README:

```
tests/
├── unit/                 # Legacy unit tests (being migrated to co-located *_test.go)
│   ├── services/        #   e.g. auth_service_test.go, health_service_test.go
│   ├── controllers/     #   e.g. auth_controller_test.go, health_controller_test.go
│   └── middlewares/     #   e.g. rate_limit_test.go, request_id_test.go
├── integration/          # Integration tests (real server / DB)
│   ├── api/             #   end-to-end API tests (e.g. health_test.go)
│   └── database/        #   database integration tests (e.g. connection_test.go)
├── mocks/                # Shared mocks for legacy tests (prefer in-package fakes for new code)
│   ├── user_repo_mock.go
│   └── auth_servicer_mock.go
├── fixtures/             # Test data (e.g. users.json)
└── README.md            # Testing documentation
```

> New code should be tested with co-located `internal/modules/<name>/*_test.go`. The
> `tests/unit/` tree holds older tests from the layered layout and is being migrated.

`make test` runs `./tests/unit/... ./internal/... ./pkg/...`, so co-located tests under
`internal/...` and `pkg/...` are included.

---

## 8. LOGGING STANDARDS

**MUST use the structured logger from `pkg/logger`.** For request-scoped tracing (START/FINISH with request_id and duration), use **LogStart** and **LogFinish** in every HTTP handler and service method; see [OBSERVABILITY.md](OBSERVABILITY.md) — section "Logger API (pkg/logger)" and "Using Request ID in Code" for full usage and examples.

### 8.1 Logger Usage

**MUST use structured logger from `pkg/logger`:**

```go
✅ CORRECT:
import "github.com/0xdiaz/gin-boilerplate/pkg/logger"

// Info level - normal operations
logger.Infof("User created successfully: ID=%d, Email=%s", user.ID, user.Email)

// Error level - operation failed but handled
logger.Errorf("Failed to send email to %s: %v", user.Email, err)

// Warning level - unexpected but not critical
logger.Warnf("Client %d approaching rate limit: %d/%d requests", clientID, current, limit)

// Debug level - detailed debugging info
logger.Debugf("Processing transaction: %+v", transaction)

❌ WRONG:
import "log"
log.Printf("User created: %v", user)  // Don't use standard log

import "fmt"
fmt.Println("Error:", err)  // Don't use fmt for logging

panic("Something went wrong")  // Don't panic for errors
```

**Request-scoped tracing (handlers and services):** At the start of each handler call
`ctx, start := logger.LogStart(c.Request.Context(), "<module>.Handler.Method")` and pass
`ctx` to the service. In each service method call
`ctx, start := logger.LogStart(ctx, "<module>.Service.Method")`. Before **every** return
(success or error), call `logger.LogFinish(ctx, "<module>.Handler.Method", err, start)` (or
the `Service` equivalent). The label convention is `<module>.<Type>.<Method>`, e.g.
`auth.Handler.Login`, `auth.Service.Login`, `example.Service.List`. `request_id` is injected
into the context by `RequestIDMiddleware`, so the START/FINISH lines carry it automatically.
For details and examples see **OBSERVABILITY.md**.

### 8.2 Log Levels

```
✅ DEBUG: Detailed debugging information (disabled in production)
✅ INFO:  Normal operations, state changes, successful operations
✅ WARN:  Unexpected situations that don't prevent operation
✅ ERROR: Operation failed, error occurred but system continues
❌ FATAL: ONLY at startup (bootstrap.Run) via logger.Fatalf, for unrecoverable errors
❌ PANIC: FORBIDDEN in application code
```

### 8.3 Log Message Format

```go
✅ CORRECT:
// Include relevant context
logger.Infof("[TrxID:%s] Transaction created: Amount=%.2f, Client=%d",
    trxID, amount, clientID)

logger.Errorf("Failed to update user %d: %v", userID, err)

// Use structured fields when available
logger.WithFields(logrus.Fields{
    "user_id":     userID,
    "transaction": trxID,
    "amount":      amount,
}).Info("Payment processed successfully")

❌ WRONG:
logger.Info("Transaction created")  // Missing context
logger.Errorf("Error: %v", err)     // Too generic
logger.Info("User:", user)          // Use structured fields instead
```

### 8.4 Sensitive Data Sanitization

**MUST sanitize before logging:**

```go
✅ CORRECT:
import "github.com/0xdiaz/gin-boilerplate/pkg/utils"

// Sanitize sensitive data
sanitizedBody := utils.SanitizeLogData(requestBody)
logger.Infof("Request body: %s", sanitizedBody)

// Mask specific fields
logger.Infof("User login: Email=%s, Password=***", user.Email)

❌ FORBIDDEN:
logger.Infof("Request: %s", requestBody)  // May contain passwords, tokens
logger.Debugf("User data: %+v", user)     // May expose sensitive fields
logger.Infof("API Key: %s", apiKey)       // NEVER log credentials
```

---

## 9. SECURITY GUIDELINES

### 9.1 Input Validation

**MUST validate ALL external input.** Declare validation rules as gin **binding tags** on
the DTO, and bind at the HTTP boundary with `c.ShouldBindJSON`. Gin uses
`go-playground/validator` under the hood; `utils.BadRequest`/`utils.FormatValidationErrors`
turn the resulting `validator.ValidationErrors` into a field→message map automatically.

```go
✅ CORRECT:
// internal/modules/<name>/dto.go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// internal/modules/<name>/handler.go
func (h *Handler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data") // auto-formats validation errors
        return
    }
    // req is validated; pass to the service
    resp, err := h.service.Register(c.Request.Context(), &req)
    // ...
}

❌ WRONG:
func (h *Handler) Register(c *gin.Context) {
    var req RegisterRequest
    _ = c.ShouldBindJSON(&req)           // bind error ignored — no validation enforced
    _, _ = h.service.Register(c.Request.Context(), &req)
}
```

### 9.2 SQL Injection Prevention

**MUST use GORM or parameterized queries:**

```go
✅ CORRECT:
// GORM automatically parameterizes
db.Where("email = ?", email).First(&user)
db.Where("age > ? AND status = ?", 18, "active").Find(&users)

❌ FORBIDDEN:
// String concatenation - SQL INJECTION RISK!
query := "SELECT * FROM users WHERE email = '" + email + "'"
db.Raw(query).Scan(&user)
```

### 9.3 Password Handling

**MUST use bcrypt for password hashing:**

```go
✅ CORRECT:
import "golang.org/x/crypto/bcrypt"

// Hash password before storing
func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hash), nil
}

// Verify password
func VerifyPassword(hash, password string) error {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

❌ FORBIDDEN:
// Storing plaintext passwords
user.Password = password

// Using weak hashing
hash := md5.Sum([]byte(password))  // MD5 is NOT secure!
```

### 9.4 Authentication & Authorization

**Protect routes with the auth module's middleware**, which validates the Bearer JWT and
sets `user_id` (a `uint`) in the gin context. Other modules obtain this guard via
`auth.Module.Middleware()` (injected through `buildModules()`); the auth module uses
`m.Middleware()` for its own protected routes.

```go
✅ CORRECT:
// module.go — mount a protected route behind the JWT guard
protected := api.Group("")
protected.Use(m.Middleware())      // auth.Module.Middleware() in other modules
protected.GET("/profile", m.handler.Profile)

// handler.go — read the authenticated user id set by the middleware
func (h *Handler) Profile(c *gin.Context) {
    userID := c.GetUint("user_id")  // set by authMiddleware; 0 if absent
    utils.Ok(c, gin.H{"user_id": userID}, "Profile retrieved successfully")
}

❌ WRONG:
func (h *Handler) Profile(c *gin.Context) {
    // No middleware on the route, and a panicking type assertion:
    user := c.MustGet("user").(*User)  // panics if missing/wrong type
    _ = user
}
```

> The middleware sets only `user_id`. There is no `"user"` context object and no built-in
> role/permission model — do not reference `c.MustGet("user")` or a `Role` field that does
> not exist. Add authorization checks in the service using the authenticated `user_id`.

### 9.5 Rate Limiting

**Rate limiting is applied to all `/api/v1` routes** in `bootstrap.buildEngine` — you do not
add it per route:

```go
✅ CORRECT (internal/bootstrap/server.go):
v1 := r.Group("/api/v1")
v1.Use(middleware.RateLimitMiddleware())   // per-client-IP limiter
for _, m := range mods {
    m.RegisterRoutes(v1)
}
```

`RateLimitMiddleware()` limits **per client IP** and reads `RATE_LIMIT_RPS` and
`RATE_LIMIT_BURST` from config (defaults 100 rps, 200 burst). Client IP is trustworthy only
when `TRUSTED_PROXIES` is configured correctly (empty = trust none, so `X-Forwarded-For`
cannot be spoofed). For a custom limit on a specific group, use
`middleware.RateLimitMiddlewareWithConfig(rps, burst)`. See
[CONFIGURATION.md](CONFIGURATION.md).

---

## 10. DATABASE ACCESS

### 10.1 Repository Pattern + Constructor Injection

**MUST use the repository pattern with an injected `*gorm.DB`.** The repository is
**module-local** (`internal/modules/<name>/repository.go`) and is built with
`NewRepository(db)`. It MUST NOT touch the global DB; the connection is passed in.

The **service** depends on the small `repository` interface it defines (see [§3.3](#33-interfaces-are-consumer-defined-di--testability)) — that
interface, not the concrete repository, is what gets faked in tests.

```go
✅ CORRECT:
// internal/modules/example/repository.go (package example)
type Repository struct {
    db *gorm.DB // injected, no global state
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

// internal/modules/example/module.go — DB injected at construction
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))
    return &Module{svc: svc, handler: NewHandler(svc)}
}

❌ WRONG:
// Reaching for the global connection inside a method
func (r *Repository) List() ([]*Example, error) {
    var list []*Example
    if err := database.GetDB().Find(&list).Error; err != nil {  // ❌ global DB, untestable
        return nil, err
    }
    return list, nil
}
```

> **DI rule:** repositories receive `*gorm.DB` via `NewRepository(db)`. The process-wide
> `database.DB` / `database.GetDB()` exists only for the connection lifecycle (opened in
> `bootstrap.Run` via `database.DbConnection(master, replica)`) and for integration tests —
> not for use inside repositories or services. The old global-DB helpers (a shared
> `repositories` package with `Save`/`Get`/`GetOne`) are GONE.

### 10.2 Transaction Management

**MUST use transactions for multi-step operations:**

```go
✅ CORRECT:
func (r *transactionRepo) CreateWithFee(tx *models.Transaction, fee *models.Fee) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        // Step 1: Create transaction
        if err := tx.Create(transaction).Error; err != nil {
            return fmt.Errorf("create transaction: %w", err)
        }

        // Step 2: Create fee record
        fee.TransactionID = transaction.ID
        if err := tx.Create(fee).Error; err != nil {
            return fmt.Errorf("create fee: %w", err)
        }

        // Step 3: Update client balance
        if err := tx.Model(&models.Client{}).
            Where("id = ?", transaction.ClientID).
            Update("balance", gorm.Expr("balance - ?", transaction.Amount)).
            Error; err != nil {
            return fmt.Errorf("update balance: %w", err)
        }

        // All succeed or all rollback
        return nil
    })
}

❌ WRONG:
func (r *transactionRepo) CreateWithFee(tx *models.Transaction, fee *models.Fee) error {
    // No transaction - partial failure possible!
    if err := r.db.Create(transaction).Error; err != nil {
        return err
    }

    if err := r.db.Create(fee).Error; err != nil {
        // Transaction already saved, fee failed - INCONSISTENT STATE!
        return err
    }

    return nil
}
```

### 10.3 Query Optimization

**MUST use eager loading when needed:**

```go
✅ CORRECT:
// Preload relationships to avoid N+1 queries
func (r *transactionRepo) ListWithDetails(page, pageSize int) ([]*models.Transaction, error) {
    var transactions []*models.Transaction

    err := r.db.
        Preload("Client").
        Preload("PaymentMethod").
        Preload("ApiLog").
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&transactions).Error

    return transactions, err
}

❌ WRONG:
// N+1 query problem
func (r *transactionRepo) ListWithDetails() ([]*models.Transaction, error) {
    var transactions []*models.Transaction
    r.db.Find(&transactions)  // 1 query

    for _, tx := range transactions {
        r.db.First(&tx.Client, tx.ClientID)  // N queries!
        r.db.First(&tx.PaymentMethod, tx.PaymentMethodID)  // N queries!
    }

    return transactions, nil
}
```

### 10.4 Soft Deletes

**MUST use soft deletes for important data:**

```go
✅ CORRECT:
// Model with soft delete
type User struct {
    ID        uint           `gorm:"primarykey"`
    Name      string         `gorm:"size:255;not null"`
    Email     string         `gorm:"size:255;unique;not null"`
    DeletedAt gorm.DeletedAt `gorm:"index"`  // Soft delete field
}

// Soft delete
db.Delete(&user)  // Sets deleted_at

// Permanent delete (use with caution)
db.Unscoped().Delete(&user)

// Include soft-deleted records
db.Unscoped().Where("id = ?", id).First(&user)

❌ WRONG:
// No soft delete - data lost permanently
type User struct {
    ID    uint   `gorm:"primarykey"`
    Name  string `gorm:"size:255"`
    Email string `gorm:"size:255"`
    // Missing DeletedAt field
}
```

---

## 11. API DESIGN

### 11.1 RESTful Conventions

**MUST follow REST principles:**

```go
✅ CORRECT:
GET    /api/v1/users           // List users
GET    /api/v1/users/:id       // Get single user
POST   /api/v1/users           // Create user
PUT    /api/v1/users/:id       // Update user (full replace)
PATCH  /api/v1/users/:id       // Update user (partial)
DELETE /api/v1/users/:id       // Delete user

// Nested resources
GET    /api/v1/users/:id/transactions      // User's transactions
POST   /api/v1/users/:id/transactions      // Create transaction for user

❌ WRONG:
POST   /api/v1/getUser         // Use GET
POST   /api/v1/createUser      // Use POST /users
GET    /api/v1/deleteUser/:id  // Use DELETE
PUT    /api/v1/user/update     // Use PUT /users/:id
```

### 11.2 HTTP Status Codes

**MUST use appropriate status codes:**

```go
✅ CORRECT:
200 OK                  // Successful GET, PUT, PATCH
201 Created             // Successful POST
204 No Content          // Successful DELETE
400 Bad Request         // Validation error, malformed request
401 Unauthorized        // Missing or invalid authentication
403 Forbidden           // Authenticated but not authorized
404 Not Found           // Resource doesn't exist
409 Conflict            // Duplicate resource, constraint violation
422 Unprocessable       // Validation failed
429 Too Many Requests   // Rate limit exceeded
500 Internal Error      // Server error

// Example usage — pick the right status via pkg/utils, never raw c.JSON
func (h *Handler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data")          // 400
        return
    }

    resp, err := h.service.Register(c.Request.Context(), &req)
    if err != nil {
        if apiErr := errToAPIError(err); apiErr != nil {          // e.g. 409 for ErrEmailAlreadyExists
            utils.RespondWithAPIError(c, apiErr)
            return
        }
        utils.InternalServerError(c, err, "Failed to register user") // 500
        return
    }

    utils.Created(c, resp, "User registered successfully")         // 201
}

❌ WRONG:
// Always returning 200 OK
c.JSON(200, gin.H{"success": false, "error": "User not found"})

// Wrong status for operation
c.JSON(404, gin.H{"message": "User created"})  // Should be 201

// Raw c.JSON instead of pkg/utils (inconsistent response shape)
c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
```

### 11.3 Response Format

**MUST use consistent response format:**

```go
✅ CORRECT:
// Success response
{
    "data": {
        "id": 123,
        "name": "John Doe",
        "email": "john@example.com"
    },
    "meta": {
        "request_id": "abc-123",
        "timestamp": "2025-11-08T10:30:00Z"
    }
}

// List response with pagination
{
    "data": [...],
    "meta": {
        "page": 1,
        "page_size": 20,
        "total": 150,
        "total_pages": 8
    }
}

// Error response
{
    "error": {
        "code": "VALIDATION_ERROR",
        "message": "Invalid input data",
        "details": {
            "email": "must be valid email format",
            "age": "must be greater than 0"
        }
    },
    "meta": {
        "request_id": "abc-123",
        "timestamp": "2025-11-08T10:30:00Z"
    }
}

❌ WRONG:
// Inconsistent format
{"success": true, "user": {...}}
{"data": {...}}
{"result": {...}, "status": "ok"}
```

### 11.4 Response Utilities

**MUST use standardized response utilities from `pkg/utils/response.go`:**

#### Success Responses

```go
✅ CORRECT - Use utility functions:
import "github.com/0xdiaz/gin-boilerplate/pkg/utils"

// 200 OK (from internal/modules/example/handler.go)
func (h *Handler) List(c *gin.Context) {
    data, err := h.svc.List(c.Request.Context())
    if err != nil {
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }
    utils.Ok(c, data, "Data retrieved successfully")
}

// 201 Created
func (h *Handler) Create(c *gin.Context) {
    item, err := h.svc.Create(c.Request.Context(), &req)
    if err != nil {
        utils.BadRequest(c, err, "Failed to create item")
        return
    }
    utils.Created(c, item, "Item created successfully")
}

// 204 No Content
func (h *Handler) Delete(c *gin.Context) {
    if err := h.svc.Delete(c.Request.Context(), id); err != nil {
        utils.InternalServerError(c, err, "Failed to delete item")
        return
    }
    utils.NoContent(c)
}

❌ WRONG - Direct c.JSON():
func (h *Handler) List(c *gin.Context) {
    data, _ := h.svc.List(c.Request.Context())
    c.JSON(200, gin.H{"data": data})  // Inconsistent format!
}
```

#### Error Responses

```go
✅ CORRECT - Use utility functions:

// 400 Bad Request
utils.BadRequest(c, err, "Invalid input data")

// 401 Unauthorized
utils.Unauthorized(c, err, "Invalid credentials")

// 403 Forbidden
utils.Forbidden(c, nil, "Access denied")

// 404 Not Found
utils.NotFound(c, err, "Resource not found")

// 409 Conflict
utils.Conflict(c, err, "Email already exists")

// 422 Unprocessable Entity
utils.UnprocessableEntity(c, err, "Validation failed")

// 429 Too Many Requests
utils.TooManyRequests(c, nil, "Rate limit exceeded")

// 500 Internal Server Error
utils.InternalServerError(c, err, "Something went wrong")

// 502 Bad Gateway
utils.BadGateway(c, err, "External service unavailable")

❌ WRONG - Direct c.JSON() with custom format:
c.JSON(400, gin.H{"error": "bad request"})
c.JSON(500, map[string]string{"message": "error"})
```

#### Available Utility Functions

**Success Functions:**
- `utils.Ok(c, data, message)` - 200 OK
- `utils.Created(c, data, message)` - 201 Created
- `utils.NoContent(c)` - 204 No Content

**Error Functions:**
- `utils.BadRequest(c, err, message)` - 400
- `utils.Unauthorized(c, err, message)` - 401
- `utils.Forbidden(c, err, message)` - 403
- `utils.NotFound(c, err, message)` - 404
- `utils.Conflict(c, err, message)` - 409
- `utils.UnprocessableEntity(c, err, message)` - 422
- `utils.TooManyRequests(c, err, message)` - 429
- `utils.InternalServerError(c, err, message)` - 500
- `utils.BadGateway(c, err, message)` - 502

**Generic Functions:**
- `utils.HandleSuccess(c, code, data, message)` - Custom success status
- `utils.HandleErrors(c, code, err, message)` - Custom error status
- `utils.HandleErrorsWithData(c, code, data, message)` - Custom error status with body (e.g. health details)
- `utils.ServiceUnavailableWithData(c, data, message)` - 503 with body; uses HandleErrorsWithData

#### Third-party response format exception

When an endpoint **must** follow an external format (e.g. DataTables server-side, third-party protocol), the response may be sent using that library’s helper as long as it is consistent and documented. Example: `example.Handler.Datatables` uses `datatables.JSON(c, dt)` for DataTables format, not `utils.Ok`. Error paths must still use response utils.

#### Service accepting gin.Context exception (DataTables)

General rule: services **MUST NOT** accept or use `gin.Context`. Exception: when a third-party library integration (e.g. Datatables-Gin) **requires** `*gin.Context` for binding query params (draw, start, length, search, order), the service method that calls that library may accept `(ctx context.Context, c *gin.Context)` to avoid duplicating the entire library API into DTOs. Example: `example.Service.Datatables(ctx, c)` and `example.Repository.Datatables(c)` — this exception is documented in [SERVICE_COMPLIANCE_AUDIT.md](SERVICE_COMPLIANCE_AUDIT.md). Other services must not use `gin.Context`.

#### Standard Response Format

All utility functions return standardized JSON format:

**Success Response:**
```json
{
    "success": true,
    "message": "User retrieved successfully",
    "data": {
        "id": 123,
        "name": "John Doe",
        "email": "john@example.com"
    },
    "errors": null
}
```

**Error Response:**
```json
{
    "success": false,
    "message": "Validation failed",
    "data": null,
    "errors": {
        "email": "email is required",
        "age": "age must be at least 18"
    }
}
```

**Exception — external protocols:** Endpoints that follow an external protocol (e.g. DataTables server-side API) may not use `utils.Ok`/SuccessResponse. The response format must follow the library or client contract. Example: `GET /datatables` returns DataTables-specific JSON, not the standard `{ success, message, data, errors }` shape.

#### Validation Errors

Validation errors are automatically formatted:

```go
✅ CORRECT:
func (h *Handler) Register(c *gin.Context) {
    var req RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        // Automatically formats validator.ValidationErrors
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    // ...
}

// Response automatically formatted as:
{
    "success": false,
    "message": "Invalid request data",
    "data": null,
    "errors": {
        "Email": "Email must be a valid email address",
        "Name": "Name is required",
        "Password": "Password must be at least 8"
    }
}
```

#### Error strategy (layers)

- **HTTP response:** All HTTP responses MUST go through `pkg/utils` (e.g. `utils.Ok`, `utils.BadRequest`, `utils.HandleErrors`). Handlers map service/domain errors to these utilities; no direct `c.JSON` for error/success.
- **Domain/business errors:** Use module-local sentinel errors in the service layer (e.g. `auth.ErrEmailAlreadyExists`, `auth.ErrInvalidCredentials`). Handlers use `errors.Is(err, auth.Err...)` (via a small `errToAPIError` mapper) to choose the right HTTP status and message. Keeps business rules in one place.
- **`pkg/types/errors`:** `APIError`, `ErrNotFound`, etc. are the shared, HTTP-shaped errors used when mapping domain errors to a response. Don't duplicate them with domain errors; prefer domain errors in services and map to `*types.APIError` in the handler. Document in code when using `types.APIError` for a given flow.

### 11.5 API Versioning

**MUST version API in URL path:**

```go
✅ CORRECT:
/api/v1/users
/api/v1/transactions
/api/v2/users  // New version with breaking changes

❌ WRONG:
/api/users (no version)
/users/v1 (version in wrong place)
```

### 11.6 Routing and Module Registration

**⚠️ CRITICAL: There is no central `routers/` package. Each module registers its own
routes in its `module.go` via `RegisterRoutes(api *gin.RouterGroup)`.** Bootstrap mounts
every module under `/api/v1`; the list of modules lives in `buildModules()`.

**How routing is wired (three places):**

```
internal/modules/<name>/module.go   →  RegisterRoutes(api)   // each module mounts its own routes
internal/bootstrap/modules.go       →  buildModules(db)      // the single list of modules
internal/bootstrap/server.go        →  buildEngine(db, mods) // mounts modules under /api/v1
```

**1. A module mounts its own routes:**

```go
// internal/modules/auth/module.go (package auth)
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    g := api.Group("/auth")
    g.POST("/register", m.handler.Register)
    g.POST("/login", m.handler.Login)
    g.POST("/refresh", m.handler.RefreshToken)
    g.POST("/forgot-password", m.handler.ForgotPassword)
    g.POST("/reset-password", m.handler.ResetPassword)

    // Protected routes use this module's own JWT guard
    protected := api.Group("")
    protected.Use(m.Middleware())
    protected.GET("/profile", m.handler.Profile)
}
```

```go
// internal/modules/example/module.go (package example)
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    api.GET("/examples", m.handler.List)
    api.GET("/datatables", m.handler.Datatables)
}
```

**2. The module list (the ONE place that knows the concrete modules):**

```go
// internal/bootstrap/modules.go
type Module interface {
    Name() string
    Models() []any
    RegisterRoutes(api *gin.RouterGroup)
}

// Adding a module is ONE line here.
func buildModules(db *gorm.DB) []Module {
    return []Module{
        auth.New(db),
        example.New(db),
    }
}
```

**3. Bootstrap mounts everything:**

```go
// internal/bootstrap/server.go
func buildEngine(db *gorm.DB, mods []Module) *gin.Engine {
    r := gin.New()
    r.Use(gin.Recovery())
    r.Use(middleware.CORSMiddleware())
    r.Use(middleware.RequestIDMiddleware())
    r.Use(middleware.RequestLogMiddleware())
    r.Use(middleware.MetricsMiddleware())

    r.NoRoute(func(c *gin.Context) { utils.NotFound(c, nil, "Route not found") })

    // System endpoints at ROOT (probes, monitoring) — NOT under /api/v1.
    health.New(db).RegisterSystem(r) // GET /health, GET /metrics

    // Business modules under /api/v1, rate-limited per IP.
    v1 := r.Group("/api/v1")
    v1.Use(middleware.RateLimitMiddleware())
    for _, m := range mods {
        m.RegisterRoutes(v1)
    }

    registerSwagger(r) // debug only
    return r
}
```

**Rules:**

1. **A module owns its routes.** Mount them in `Module.RegisterRoutes(api)`; never add a
   module's routes from `bootstrap` or another module.
2. **Business routes go under `/api/v1`** (the group passed to `RegisterRoutes`). Use
   sub-groups inside the module (e.g. `api.Group("/auth")`).
3. **System routes (`/health`, `/metrics`) mount at ROOT** via the health module's
   `RegisterSystem(r)`, not under `/api/v1`.
4. **Protected routes use the auth module's middleware**: `m.Middleware()` inside auth, or
   the injected `authMod.Middleware()` in another module.
5. **Add a module in exactly one line** in `buildModules()`. Models and routes are picked up
   automatically (migrations from `Module.Models()`, routes from `RegisterRoutes`).

> ❌ Do NOT create a `routers/` package, `index.go`, or `Register*Routes(router)` functions —
> that layout no longer exists. Routing is co-located with the module that owns it.

---

## 12. CONFIGURATION

### 12.1 Environment Variables

**MUST use environment variables for configuration:**

Configuration is loaded from `.env` by `pkg/config` (viper-backed) in `config.SetupConfig()`,
validated, and exposed via `config.Get()`. Access it through the typed `Configuration`
struct, not by reading env vars ad hoc.

```go
✅ CORRECT:
// .env file — real keys (see CONFIGURATION.md for the full list)
MASTER_DB_HOST=localhost
MASTER_DB_PORT=5432
MASTER_DB_USER=postgres
MASTER_DB_PASSWORD=secret
MASTER_DB_NAME=project_db

JWT_SECRET=change-me-to-a-32+char-random-secret   # required, min 32 chars
SERVER_HOST=0.0.0.0
SERVER_PORT=8000
TRUSTED_PROXIES=                                  # CSV of proxy CIDRs/IPs; empty = trust none
RATE_LIMIT_RPS=100                                # per-IP requests/sec (default 100)
RATE_LIMIT_BURST=200                              # per-IP burst (default 200)
REQUEST_TIMEOUT_SECONDS=30
SERVER_TIMEZONE=UTC

// pkg/config — typed configuration (read via config.Get())
type Configuration struct {
    Server   ServerConfiguration   // Host, Port, JWTSecret, TrustedProxies, RateLimitRPS, RateLimitBurst, …
    Database DatabaseConfiguration // master + replica DSN parts
}

cfg := config.Get()                 // nil-safe: returns nil if not loaded
secret := cfg.Server.JWTSecret

❌ WRONG:
// Hardcoded in code
const (
    DBHost     = "localhost"      // ❌ Hardcoded
    DBPassword = "secret123"      // ❌ NEVER commit passwords!
    JWTSecret  = "my-secret-key"  // ❌ Security risk
)

// Stale/removed env names — do NOT use:
// ALLOWED_HOSTS         → replaced by TRUSTED_PROXIES
// RATE_LIMIT            → replaced by RATE_LIMIT_RPS / RATE_LIMIT_BURST
// RATE_LIMIT_USE_USER   → removed (rate limit is per IP on all /api/v1)
```

### 12.2 Configuration File Structure

```
✅ REQUIRED files:
.env.example    // Template with all keys, no sensitive values
.env            // Actual config (in .gitignore)

❌ FORBIDDEN:
.env            // Committed to git with secrets
config.json     // With hardcoded passwords
```

### 12.3 Secrets Management

**MUST NOT commit secrets:**

```bash
✅ .gitignore MUST include:
.env
*.key
*.pem
secrets/
credentials.json
service-account.json

❌ FORBIDDEN in repository:
- API keys
- Database passwords
- JWT secrets
- Private keys
- Service account files
```

---

## 13. FORBIDDEN PRACTICES

### 13.1 Absolutely Forbidden

```go
❌ Global mutable state
var globalUser *User  // Race conditions!

❌ init() functions with side effects
func init() {
    db.Connect()  // Unpredictable initialization order
}

❌ Panic in business logic
func CreateUser(user *User) {
    if user.Email == "" {
        panic("email required")  // Use error instead!
    }
}

❌ Ignoring errors
user, _ := repo.GetUser(id)  // Error ignored!

❌ Naked returns in long functions
func ProcessTransaction(tx Transaction) (result Result, err error) {
    // ... 100 lines of code ...
    return  // What are we returning?
}

❌ Type assertions without checking
user := c.Get("user").(*User)  // Panic if type wrong!

❌ Goroutines without context/timeout
go processInBackground()  // No way to cancel!

❌ String concatenation for SQL
query := "SELECT * FROM users WHERE id = " + id  // SQL injection!

❌ Using == for float comparison
if amount == 100.50 { }  // Floating point precision issues!

❌ Modifying slice/map during iteration
for k := range m {
    delete(m, k)  // Undefined behavior!
}
```

### 13.2 Discouraged Practices

```go
⚠️ Deep nesting (>3 levels)
if x {
    if y {
        if z {
            if a {  // Too deep!
            }
        }
    }
}
// Refactor with early returns

⚠️ else after return
if err != nil {
    return err
} else {  // Unnecessary else
    return nil
}

⚠️ Single-letter variable names (except i, j, k in loops)
u := GetUser()  // Use 'user' instead

⚠️ Premature optimization
// Don't optimize until profiling shows it's needed

⚠️ Not using defer for cleanup
file, _ := os.Open("file.txt")
// ... many lines ...
file.Close()  // Might be skipped if error occurs
// Use: defer file.Close()
```

---

## 14. CODE REVIEW CHECKLIST

### Before Committing

- [ ] All functions under 100 lines
- [ ] All files under 300 lines
- [ ] No hardcoded secrets or passwords
- [ ] All exported functions documented
- [ ] All errors handled (no `_` for errors)
- [ ] No `panic()` in business logic
- [ ] No `TODO` comments older than 30 days
- [ ] All tests passing
- [ ] Code coverage >70% for new services
- [ ] No commented-out code
- [ ] Imports organized (stdlib, external, internal)
- [ ] `gofmt` applied
- [ ] `golint` passes
- [ ] `go vet` passes

### Code Quality Checks

```bash
# Format code
gofmt -w .

# Lint
golangci-lint run

# Vet
go vet ./...

# Test with coverage
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out

# Security check
gosec ./...
```

### Pull Request Checklist

- [ ] Descriptive PR title following convention: `feat:`, `fix:`, `refactor:`, `docs:`, `test:`
- [ ] Description explains WHAT and WHY
- [ ] Related issue linked
- [ ] Screenshots/logs for UI/API changes
- [ ] Database migrations included if schema changed
- [ ] Documentation updated
- [ ] Changelog updated
- [ ] No merge conflicts
- [ ] CI/CD pipeline passing

---

## 15. AI AGENT SPECIFIC RULES

### 15.1 When Writing New Code

1. **ALWAYS** start a new feature by copying the `example` module
   - `cp -r internal/modules/example internal/modules/<name>`, rename the package
   - Register it with one line in `buildModules()` (`internal/bootstrap/modules.go`)
   - See `MODULE_GUIDE.md` — "How to add a new module"

2. **ALWAYS** keep the feature in one module package (vertical slice)
   - `model.go`, `dto.go`, `repository.go`, `service.go`, `handler.go`, `module.go`
   - The service defines the `repository` interface it needs; the handler defines the
     `service` interface it needs (consumer-defined; see §3.3)

3. **ALWAYS** check file size (>250 lines → split into another file in the same package)
   and function size (>80 lines → split into smaller functions)

4. **ALWAYS** add documentation
   - Package comment if new module/package
   - Function comment for exported functions; struct comment for exported types

5. **ALWAYS** add error handling
   - Never ignore errors; always wrap with `%w` for context
   - Never use panic (except `bootstrap.Run` startup via `logger.Fatalf`)

6. **ALWAYS** write a co-located test
   - `internal/modules/<name>/service_test.go` using a fake repo (no DB), per
     `example/service_test.go`
   - At least happy path + 2 error cases

### 15.2 When Refactoring Code

1. **MUST** maintain backward compatibility
   - Don't change a module's public surface (`Module.API()`, `auth.Servicer`,
     `Module.Middleware()`) without discussion; other modules depend on it

2. **MUST** add tests before refactoring
   - Ensure existing functionality preserved

3. **MUST** refactor incrementally
   - Small, reviewable changes
   - One concept per commit

4. **MUST** update documentation
   - If function signature changes, update comment
   - If behavior changes, update docs

### 15.3 When Reviewing Code

1. **CHECK** all items in Code Review Checklist
2. **VERIFY** no forbidden practices used
3. **CONFIRM** test coverage adequate
4. **VALIDATE** error handling complete
5. **ENSURE** documentation present

---

## 16. RESOURCES

### Official Go Documentation
- Style Guide: https://google.github.io/styleguide/go/
- Effective Go: https://go.dev/doc/effective_go
- Code Review Comments: https://github.com/golang/go/wiki/CodeReviewComments

### Tools
- `gofmt` - Code formatter
- `golint` - Linter
- `go vet` - Static analyzer
- `golangci-lint` - Meta-linter
- `gosec` - Security checker

### Testing
- Standard library: `testing`
- Assertions: `github.com/stretchr/testify/assert` (used in `example/service_test.go`)
- Prefer in-package fakes over generated mocks for new modules; shared mocks live in `tests/mocks/`
- HTTP testing: `httptest`

### Project docs
- [`MODULE_GUIDE.md`](./MODULE_GUIDE.md) — the modular layout (source of truth)
- [`CONFIGURATION.md`](./CONFIGURATION.md) — environment variables
- [`OBSERVABILITY.md`](./OBSERVABILITY.md) — logger API, request IDs, tracing
- [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) — non-negotiable rules

---

**END OF CODING STANDARDS**

*Last updated: 2026-06-10*
*Version: 2.0 (modular architecture)*
*Review: Every 3 months or after major changes*
