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
- **[3. Code Structure](#3-code-structure)** — controller → service → repository, repository interfaces, injection
- **[7. Testing Requirements](#7-testing-requirements)** — tests/ tree, shared fake repos
- **[11.4 Response Utilities](#114-response-utilities)** — `pkg/utils` (MANDATORY)
- **[11.6 Routing](#116-routing-and-route-registration)** — one <name>_routes.go per feature

### 📖 How to Use This Document

1. ✅ Read `00_AI_CRITICAL_RULES.md` and `MODULE_GUIDE.md` first
2. ✅ Read the critical sections listed above
3. ⚠️  Skim other sections for context
4. 📚 Use as reference when needed

---

## 📋 Table of Contents

### 🔥 Critical Sections (Read First)

- [1. File Organization](#1-file-organization) — `layers`, `feature prefix`, `300 lines`
  - [1.5 DTO & Constant Placement](#15-dto-and-constant-placement-rules) — `internal/app/dto`
- [3. Code Structure](#3-code-structure) — `handler`, `service`, `repository`, `interface`
- [5. Error Handling](#5-error-handling) — `fmt.Errorf`, `%w`, `domain errors`, `APIError`
- [7. Testing Requirements](#7-testing-requirements) — `tests/unit`, `fake repo`, `70%`, `coverage`
- [11. API Design](#11-api-design) — `REST`, `response`, `utils.Ok`, `routes`
  - [11.4 Response Utilities](#114-response-utilities) — `utils.Ok`, `utils.Created`, `utils.BadRequest`
  - [11.6 Routing & Route Registration](#116-routing-and-route-registration) — `Register*Routes`, `index.go`

### 📚 All Sections

| # | Section | Keywords |
|---|---------|----------|
| 1 | [File Organization](#1-file-organization) | `file size`, `layers`, `dto`, `feature prefix` |
| 1.1 | [File Size Limits](#11-file-size-limits) | `300 lines`, `100 lines function` |
| 1.2 | [Directory Structure](#12-directory-structure) | `internal/app`, `internal/domain`, `pkg/` |
| 1.3 | [File Naming](#13-file-naming) | `snake_case`, `_test.go`, `naming` |
| 1.4 | [Package Organization](#14-package-organization) | `package`, `one layer = one package` |
| 1.5 | [DTO & Constant Placement](#15-dto-and-constant-placement-rules) | `dto`, `constants`, `placement` |
| 2 | [Naming Conventions](#2-naming-conventions) | `camelCase`, `PascalCase`, `naming` |
| 3 | [Code Structure](#3-code-structure) | `handler`, `service`, `repository`, `interface` |
| 4 | [Function Guidelines](#4-function-guidelines) | `function size`, `parameters`, `return` |
| 5 | [Error Handling](#5-error-handling) | `error`, `wrapping`, `domain errors`, `APIError` |
| 6 | [Documentation](#6-documentation) | `comment`, `godoc`, `docs` |
| 7 | [Testing Requirements](#7-testing-requirements) | `test`, `tests/unit`, `fake repo`, `70%` |
| 8 | [Logging Standards](#8-logging-standards) | `logger`, `LogStart`, `LogFinish`, `request_id` |
| 9 | [Security Guidelines](#9-security-guidelines) | `validation`, `SQL injection`, `bcrypt` |
| 10 | [Database Access](#10-database-access) | `GORM`, `model`, `migration`, `DI` |
| 11 | [API Design](#11-api-design) | `REST`, `API`, `response`, `routes` |
| 11.1 | [RESTful Conventions](#111-restful-conventions) | `REST`, `HTTP methods`, `endpoints` |
| 11.2 | [HTTP Status Codes](#112-http-status-codes) | `200`, `404`, `500`, `status` |
| 11.3 | [Response Format](#113-response-format) | `success`, `message`, `data`, `errors` |
| 11.4 | [Response Utilities](#114-response-utilities) | `utils.Ok`, `utils.Created`, `RespondWithAPIError` |
| 11.5 | [API Versioning](#115-api-versioning) | `v1`, `v2`, `versioning` |
| 11.6 | [Routing & Route Registration](#116-routing-and-route-registration) | `Register*Routes`, `index.go` |
| 12 | [Configuration](#12-configuration) | `config`, `env`, `.env` |
| 13 | [Forbidden Practices](#13-forbidden-practices) | `panic`, `global`, `anti-pattern` |
| 14 | [Code Review Checklist](#14-code-review-checklist) | `checklist`, `review`, `validation` |
| 15 | [AI Agent Rules](#15-ai-agent-specific-rules) | `AI`, `agent`, `automation` |
| 16 | [Resources](#16-resources) | `links`, `references`, `docs` |

### 🎯 Quick Lookups by Task

**Creating a new feature (controller/service/repository):**
- [3. Code Structure](#3-code-structure) — consumer-defined interfaces + constructor injection
- [11.4 Response Utilities](#114-response-utilities) — MUST use
- See: `MODULE_GUIDE.md` — "How to add a new module"

**Adding routes:**
- [11.6 Routing & Route Registration](#116-routing-and-route-registration)
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
- Split into multiple focused files **within the same package**.
- **No subfolders inside a layer.** When one file grows past the limit, add another file in the same folder with the same package. Example: the auth service splits into `internal/app/services/auth/auth_service.go` + `internal/app/services/auth/auth_service_tokens.go` — both `package auth`.
- The filename prefix identifies the feature; a split file does not get its own subfolder.

### 1.2 Directory Structure

Code is organized **by technical layer**, not by feature folder. A feature is a horizontal
slice: one file in each layer directory, tied together by a shared filename prefix. See
[`MODULE_GUIDE.md`](./MODULE_GUIDE.md) for the canonical version.

**MUST follow this structure:**

```
.
├── main.go                          # config → db → migrate → seed (dev) → serve → shutdown
├── internal/
│   ├── adapters/
│   │   └── database/
│   │       ├── database.go          #   DbConnection(master, replica), GetDB(), package-level DB
│   │       ├── migrations/
│   │       │   ├── migration.go     #     golang-migrate runner, fatal on failure
│   │       │   └── sql/             #     versioned NNNNNN_name.up.sql / .down.sql pairs
│   │       └── seeders/             #   idempotent demo data, development only
│   ├── app/
│   │   ├── controllers/             #   HTTP layer — <name>_controller.go
│   │   ├── dto/                     #   request/response types — <name>_dto.go
│   │   ├── middlewares/             #   auth, cors, metrics, rate_limit, request_id, request_log
│   │   ├── routers/
│   │   │   ├── router.go            #     gin engine + global middleware
│   │   │   ├── index.go             #     RegisterRoutes(): builds repos + services, mounts groups
│   │   │   ├── <name>_routes.go     #     Register<Name>Routes(group, service), one per feature
│   │   │   └── swagger.go           #     OpenAPI/Swagger UI (debug only)
│   │   └── services/                #   business logic — <name>_service.go
│   │       └── auth/                #     a service that outgrew one file gets its own package
│   └── domain/
│       ├── models/                  #   GORM structs with TableName() — <name>_model.go
│       └── repositories/            #   interface + unexported impl + New*Repository()
├── pkg/                             # ← cross-cutting "kit"; MUST never import from internal/
│   ├── config/                      #   configuration (viper-backed)
│   ├── logger/                      #   structured logging + request-scoped tracing
│   ├── metrics/                     #   request counters + uptime
│   ├── types/                       #   shared response/error types
│   └── utils/                       #   response helpers (utils.Ok, utils.BadRequest, …)
└── tests/
    ├── unit/{controllers,services,middlewares}/   # package <layer>_test, no database
    ├── integration/{api,database}/                # needs a real database
    ├── mocks/                                     # shared in-memory fakes
    └── fixtures/                                  # test data
```

> **Testing layout:** unit tests live in the root `tests/` tree, **not** beside the code they
> test. A service test goes in `tests/unit/services/<name>_service_test.go` (package
> `services_test`), and the fake it drives goes in `tests/mocks/<name>_repo_mock.go` so other
> tests can reuse it. Tests that need a real database go in `tests/integration/`. Fixture data
> lives in `tests/fixtures/`.

#### 1.2.1 One Layer = One Package = One Folder

```
✅ MUST: A feature contributes one file per layer, all sharing the feature name as a prefix:
   internal/domain/models/<name>_model.go
   internal/domain/repositories/<name>_repo.go
   internal/app/dto/<name>_dto.go
   internal/app/services/<name>_service.go
   internal/app/controllers/<name>_controller.go
   internal/app/routers/<name>_routes.go
```

**Purpose:**
- Each layer has one obvious home, so a reviewer knows where to look for a rule.
- The filename prefix keeps the slice greppable: `ls internal/app/**/auth_*` shows the feature.
- Cross-layer boundaries stay enforceable: only repositories import `database`.

**Example (the canonical `auth` slice):**
```
internal/adapters/database/migrations/sql/000001_create_users_table.up.sql
internal/domain/models/user_model.go            # User (GORM model) + TableName()
internal/domain/repositories/user_repo.go       # UserRepository interface + userRepo impl
internal/app/dto/auth_dto.go                    # UserResponse
internal/app/services/auth/auth_service.go           # AuthService + ErrUserNotFound
internal/app/controllers/auth_controller.go     # AuthController
internal/app/routers/auth_routes.go             # RegisterAuthRoutes(group, service)
tests/mocks/user_repo_mock.go                   # MockUserRepository
tests/unit/services/auth_service_test.go        # unit test, no DB
```

**Naming Rules:**
- The filename prefix *is* the feature: `auth_`, `user_`, `health_`.
- All files use `snake_case`; every file in a layer folder shares that layer's package.
- No feature subfolders inside a layer — the one exception is a service that needs several
  files, which gets its own package (`internal/app/services/auth/`).

### 1.3 File Naming

**Rules:**
```go
✅ CORRECT:
- auth_service.go          (snake_case; the feature's business logic)
- user_repo.go             (snake_case)
- auth_service_tokens.go   (snake_case; a second service file in the same package)
- auth_service_test.go     (test, under tests/unit/services/)

❌ WRONG:
- AuthService.go          (PascalCase not allowed)
- auth-service.go          (kebab-case not allowed)
- authService.go          (camelCase not allowed)
```

The standard files for a feature are `<name>_model.go`, `<name>_repo.go`, `<name>_dto.go`,
`<name>_service.go`, `<name>_controller.go`, `<name>_routes.go`. When a file outgrows the size
limit, append a descriptive suffix (`auth_service_tokens.go`, `user_repo_query.go`) in the
same folder.

### 1.4 Package Organization

**One layer = one package = one directory:**
```go
// ✅ CORRECT — every file in the folder shares the layer package
internal/app/services/
    ├── health_service.go  → package services
    ├── health_service.go  → package services
    └── auth/              → package auth  (a service with several files)
        ├── auth_service.go
        └── auth_service_tokens.go

// ❌ WRONG - a subfolder per feature inside a layer
internal/app/services/
    ├── health/service.go  → package health
    └── health/service.go  → package health
```

### 1.5 DTO and Constant Placement Rules

DTOs live in one place (`internal/app/dto`). Constants live next to the thing they describe.
There is no `pkg/enums/` package.

#### 1.5.1 DTOs (Data Transfer Objects)
```
✅ DTOs LIVE IN internal/app/dto/<name>_dto.go (package dto)

Location rules:
- Request/Response structs → internal/app/dto/<name>_dto.go
- Any struct used for data transfer at the HTTP boundary → the dto package

Examples:
✅ CORRECT:
internal/app/dto/auth_dto.go     → RegisterRequest, LoginRequest, AuthResponse, UserResponse
internal/app/dto/auth_dto.go    → UserResponse
internal/app/dto/health_dto.go   → HealthResponse, MetricsResponse

❌ WRONG:
internal/domain/models/user_model.go → UserResponse (models are GORM structs, not API shapes)
pkg/types/request.go                  → ApiRequest (pkg/types is for shared response/error types only)
```

DTOs carry their own validation via gin **binding tags** (see [§9.1](#91-input-validation)):

```go
// internal/app/dto/auth_dto.go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}
```

#### 1.5.2 Constants and Enum-like Values
```
✅ CONSTANTS LIVE NEXT TO THE CODE THAT GIVES THEM MEANING

Location rules:
- A model's status/role/type constants → that model's file, in package models
- A service's sentinel errors           → that service's file
- A truly cross-cutting constant        → pkg/ (e.g. pkg/types)

There is NO pkg/enums package. Do not invent one.

Examples:
✅ CORRECT:
internal/domain/models/<name>_model.go      → status/role/type constants for that model
                                              (must match any CHECK constraint in the migration)
internal/app/services/auth/auth_service.go  → ErrInvalidCredentials, ErrUserNotFound,
                                              ErrEmailAlreadyExists, ErrMailerNotConfigured
pkg/types/errors.go                         → APIError and the shared HTTP mapping

❌ WRONG:
pkg/enums/role.go                 → (no such package)
```

#### 1.5.3 Import Direction (no cycles)

**CRITICAL: `pkg/` MUST never import from `internal/`. A layer MUST never import a layer above
it.**

```go
✅ ALLOWED:
internal/...                 → pkg/utils, pkg/types, pkg/logger, pkg/config   ✓
internal/app/controllers     → internal/app/services, internal/app/dto        ✓
internal/app/services        → internal/domain/repositories, internal/domain/models ✓
internal/domain/repositories → internal/domain/models, internal/adapters/database   ✓
internal/app/routers         → controllers, services, repositories, middlewares     ✓

❌ FORBIDDEN (cause import cycles or break the layering):
pkg/anything                 → internal/...                                   ✗
internal/domain/repositories → internal/app/services                          ✗
internal/app/controllers     → internal/domain/repositories (skipping a layer) ✗
```

Dependencies are injected, never constructed in place: a service receives repository
interfaces through `New*Service(...)`, assembled in `internal/app/routers/index.go`.
See [§3.5](#35-crossing-features).

---

## 2. NAMING CONVENTIONS

### 2.1 Package Names

```go
✅ MUST:
- Lowercase, single word
- No underscores, no dashes
- Layer packages are plural; a service with its own package uses the singular feature name

✅ CORRECT:
package controllers   // internal/app/controllers
package services      // internal/app/services
package repositories  // internal/domain/repositories
package models        // internal/domain/models
package auth          // internal/app/services/auth
package utils         // pkg/utils

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

A feature spans the layers. The flow is
`handler → service → repository → model`. The HTTP layer is the **Handler** (type
`Handler`, constructor `NewAuthController`), not a "controller".

### 3.1 The Vertical Slice (handler → service → repository → model)

**MUST follow this flow:**

```
HTTP Request
    ↓
[Register<Name>Routes] → mounts the route on the /api/v1 group
    ↓
[Controller] → thin HTTP layer (parse/bind request, call service, write response via pkg/utils)
    ↓
[Service] → business logic (validation, orchestration, domain errors)
    ↓
[Repository] → data access (CRUD / queries on the injected *gorm.DB)
    ↓
[Model] → GORM entity in internal/domain/models
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
- Import the service or controller layer

✅ Repository SHOULD:
- Perform CRUD / build queries on its injected *gorm.DB
- Be constructed with NewRepository(db) — no global DB access
- Return models from internal/domain/models (or errors)

❌ Repository MUST NOT:
- Contain business logic
- Contain business rules
```

### 3.2 Dependency Direction

**MUST follow:**
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
❌ A repository importing a service
❌ pkg/ importing internal/
❌ Circular dependencies
```

### 3.3 Interfaces are Consumer-Defined (DI + testability)

**MANDATORY: dependencies are declared as interfaces and handed in through the
constructor.** This is what makes a service unit-testable with a fake (no DB) and keeps
coupling loose.

- The **repository** publishes the interface it satisfies (in `<name>_repo.go`).
- The **service** depends on that interface, received via `New<Name>Service(...)`.
- The **controller** holds the service, received via `New<Name>Controller(...)`.
- `internal/app/routers/index.go` wires the concrete types together.

```go
// internal/domain/repositories/user_repo.go
// UserRepository is the data-access contract the service depends on.
type UserRepository interface {
    List() ([]*models.User, error)
    GetByID(id uint) (*models.User, error)
}

type userRepo struct{}

func NewUserRepository() UserRepository { return &userRepo{} }
```

```go
// internal/app/services/auth/auth_service.go
type AuthService struct {
    userRepo  repositories.UserRepository
    refreshTokenRepo repositories.RefreshTokenRepository
}

func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository) *AuthService {
    return &AuthService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo}
}
```

```go
// internal/app/controllers/auth_controller.go
type AuthController struct {
    service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
    return &AuthController{service: service}
}
```

```go
// internal/app/routers/index.go — the single wiring point
userRepo := repositories.NewUserRepository()
refreshTokenRepo := repositories.NewRefreshTokenRepository()
authService := services.NewAuthService(userRepo, refreshTokenRepo)
RegisterAuthRoutes(apiV1, authService)
```

### 3.4 Single Responsibility Principle

**Each file has ONE clear purpose:**

```go
✅ CORRECT:
// internal/app/services/auth/auth_service.go — auth business logic only
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

### 3.5 Crossing Features

**A service that needs another feature depends on an interface, injected through its
constructor — never by constructing the dependency itself or reaching for a global.**

```go
// The auth service publishes an interface so the middleware can depend on the contract
// rather than the concrete *AuthService:
// internal/app/services/auth/interface.go
type AuthServicer interface {
    ValidateToken(tokenString string) (*Claims, error)
    // ...
}

// internal/app/middlewares/auth.go
func AuthMiddleware(authService auth.AuthServicer) gin.HandlerFunc { /* ... */ }
```

```go
// A service that needs another feature's data takes its repository interface:
func NewOrderService(orderRepo repositories.OrderRepository, refreshTokenRepo repositories.RefreshTokenRepository) *OrderService {
    return &OrderService{orderRepo: orderRepo, refreshTokenRepo: refreshTokenRepo}
}
```

**Rules:**
- Depend on the smallest interface that does the job (`auth.AuthServicer`, `UserRepository`).
- Depend on the interface, not the concrete `*Service` or `*repo`.
- Wire the dependency in `internal/app/routers/index.go` — the one place that knows concrete
  types.
- Never call `database.DB` from a service to shortcut another feature's repository.

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
// internal/app/services/auth/auth_service.go (package auth)

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

### 5.3 Domain Errors (declared by the service)

**Define a feature's domain errors in the service that returns them**, as sentinel `error` values.
The handler maps them to HTTP (see §5.5). Shared, HTTP-shaped errors (`*types.APIError`)
live in `pkg/types/errors.go`.

```go
✅ CORRECT:
// internal/app/services/auth/auth_service.go (package auth)
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
    - Application startup (main.go via logger.Fatalf)
    - Fatal unrecoverable startup errors (config load / DB connect failure)

✅ The gin engine installs gin.Recovery() in routers.SetupRoute(), so panics in a
   handler are caught and turned into 500 — but you MUST NOT rely on it: return errors.
```

```go
✅ CORRECT:
// main.go — Fatalf is acceptable at startup
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

**Domain errors:** Services return their own sentinel errors (e.g.
`auth.ErrEmailAlreadyExists`, `auth.ErrInvalidCredentials`). The **handler** maps them to
HTTP using either:
- `utils.RespondWithAPIError(c, apiErr)` when the error is mapped to a `*types.APIError`
  (code, message, details), or
- `utils.InternalServerError(c, err, "…")` (and similar) for unmapped/generic errors.

The canonical pattern (from `internal/app/controllers/auth_controller.go`) is a small mapping helper
plus `errors.Is`:

```go
// internal/app/controllers/auth_controller.go
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

func (ctrl *AuthController) Login(c *gin.Context) {
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
// Package auth implements authentication: registration, login, JWT
// issuance/validation, refresh token rotation, and password reset.
//
// It exposes the AuthServicer contract so the HTTP layer and the auth
// middleware depend on the interface rather than the concrete *AuthService.
//
// Example usage (from internal/app/routers/index.go):
//
//	authService := auth.NewAuthService(userRepo, refreshTokenRepo, nil)
//	protectedRoutes.Use(middlewares.AuthMiddleware(authService))
package auth
```

**File to add package comment:**
- Choose the feature's main file (e.g. `<name>_service.go`) or create `doc.go`.

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

### 7.1 Test File Location — the `tests/` tree

**⚠️ CRITICAL: unit tests live in `tests/unit/<layer>/`, not beside the code they test.**
Unit-test the service against a **fake repository** from `tests/mocks/` (no DB). This is the
canonical pattern; see `tests/mocks/user_repo_mock.go` and
`tests/unit/services/auth_service_test.go`.

Fakes are shared rather than re-declared per test file, so one interface change surfaces in
one place. Fixture data lives in `tests/fixtures/`. Tests that need a real database or a
running server live in `tests/integration/`.

```
✅ CORRECT:
tests/unit/services/auth_service_test.go      // unit test, fake repo, no DB
tests/unit/controllers/auth_controller_test.go // httptest against the real controller
tests/mocks/user_repo_mock.go                 // shared fake, imported by the tests
tests/fixtures/users.json                      // fixture data
tests/integration/database/connection_test.go  // needs a real database

❌ WRONG:
internal/app/services/auth_service_test.go     // tests do not live beside the code here
internal/domain/repositories/user_repo_test.go
```

**Package naming for unit tests — black box (`package <layer>_test`):**

Tests live in the external test package so they exercise the code exactly as production
callers do, through its exported surface. This works because repositories publish **exported**
interfaces, so a fake in `tests/mocks/` can satisfy them from outside.

```go
✅ CORRECT:
// tests/unit/services/auth_service_test.go
package services_test  // black box — drives the exported API

import (
    "context"
    "testing"

    "github.com/0xdiaz/gin-boilerplate/internal/app/services"
    "github.com/0xdiaz/gin-boilerplate/tests/mocks"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

❌ WRONG:
// tests/unit/services/auth_service_test.go
package services  // white-box test placed inside the layer package
```

Every fake asserts at compile time that it still matches the real interface:

```go
// tests/mocks/user_repo_mock.go
var _ repositories.UserRepository = (*MockUserRepository)(nil)
```

That one line is what turns an interface change into a compile error instead of a silently
stale test.

### 7.2 Test Coverage Requirements

```
✅ MINIMUM: 70% code coverage for the service (business logic) of each feature
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

### 7.5 Canonical Pattern — Service Test with a Shared Fake Repo

**This is the pattern to copy** (from `tests/mocks/user_repo_mock.go` and
`tests/unit/services/auth_service_test.go`). The fake implements the **exported
`repositories.UserRepository` interface**, so the service is tested with no DB — and because
the interface is exported, the test can stay black-box (`package services_test`).

```go
// tests/unit/services/auth_service_test.go
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

**Controller tests** use `httptest` against the real controller, with a service built on the
fakes in `tests/mocks/`. See `tests/unit/controllers/auth_controller_test.go`.

### 7.6 Test Organization

All tests live under `tests/`:

```
tests/
├── unit/                 # Unit tests, package <layer>_test, no database
│   ├── services/        #   e.g. auth_service_test.go, health_service_test.go
│   ├── controllers/     #   e.g. auth_controller_test.go, health_controller_test.go
│   └── middlewares/     #   e.g. rate_limit_test.go, request_id_test.go
├── integration/          # Integration tests (real server / DB)
│   ├── api/             #   end-to-end API tests (e.g. health_test.go)
│   └── database/        #   database integration tests (e.g. connection_test.go)
├── mocks/                # Shared in-memory fakes
│   ├── user_repo_mock.go
│   ├── refresh_token_repo_mock.go
│   ├── user_repo_mock.go
│   └── auth_servicer_mock.go
├── fixtures/             # Test data (e.g. users.json)
└── README.md            # Testing documentation
```

Adding a feature means adding one fake to `tests/mocks/` and one test file to
`tests/unit/services/`.

`make test` runs the unit suite under `./tests/unit/...`; `make test-coverage` produces a
coverage report.

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
`ctx, start := logger.LogStart(c.Request.Context(), "<Type>.<Method>")` and pass
`ctx` to the service. In each service method call
`ctx, start := logger.LogStart(ctx, "<Type>.<Method>")`. Before **every** return
(success or error), call `logger.LogFinish(ctx, "<Type>.<Method>", err, start)` (or
the `Service` equivalent). The label convention is `<Type>.<Method>`, e.g.
`auth.Handler.Login`, `auth.Service.Login`, `example.Service.List`. `request_id` is injected
into the context by `RequestIDMiddleware`, so the START/FINISH lines carry it automatically.
For details and examples see **OBSERVABILITY.md**.

### 8.2 Log Levels

```
✅ DEBUG: Detailed debugging information (disabled in production)
✅ INFO:  Normal operations, state changes, successful operations
✅ WARN:  Unexpected situations that don't prevent operation
✅ ERROR: Operation failed, error occurred but system continues
❌ FATAL: ONLY at startup (main.go) via logger.Fatalf, for unrecoverable errors
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
// internal/app/dto/<name>_dto.go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// internal/app/controllers/<name>_controller.go
func (ctrl *AuthController) Register(c *gin.Context) {
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
func (ctrl *AuthController) Register(c *gin.Context) {
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

**Protect routes with `middlewares.AuthMiddleware(authService)`**, which validates the
Bearer JWT and sets `user_id` (a `uint`) in the gin context. It takes the
`auth.AuthServicer` interface, so it depends on the contract rather than the concrete
service.

```go
✅ CORRECT:
// internal/app/routers/index.go — mount protected routes behind the JWT guard
authController := controllers.NewAuthController(authService)
protectedRoutes := apiV1.Group("")
protectedRoutes.Use(middlewares.AuthMiddleware(authService))
{
    protectedRoutes.GET("/profile", authController.Profile)
}

// internal/app/controllers/auth_controller.go — read the user id set by the middleware
func (ctrl *AuthController) Profile(c *gin.Context) {
    userID := c.GetUint("user_id")  // set by authMiddleware; 0 if absent
    utils.Ok(c, gin.H{"user_id": userID}, "Profile retrieved successfully")
}

❌ WRONG:
func (ctrl *AuthController) Profile(c *gin.Context) {
    // No middleware on the route, and a panicking type assertion:
    user := c.MustGet("user").(*User)  // panics if missing/wrong type
    _ = user
}
```

> The middleware sets only `user_id`. There is no `"user"` context object and no built-in
> role/permission model — do not reference `c.MustGet("user")` or a `Role` field that does
> not exist. Add authorization checks in the service using the authenticated `user_id`.

### 9.5 Rate Limiting

**Rate limiting is applied to all `/api/v1` routes** in `RegisterRoutes` (`internal/app/routers/index.go`) — you do not
add it per route:

```go
✅ CORRECT (internal/app/routers/router.go):
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

**MUST use the repository pattern.** A repository lives at
`internal/domain/repositories/<name>_repo.go` and publishes an **exported interface**, an
unexported struct implementing it, and a `New<Name>Repository()` constructor returning the
interface.

Repositories are the **only** layer allowed to touch the package-level `database.DB` handle.
The **service** depends on the repository interface — that interface, not the concrete struct,
is what gets faked in tests.

```go
✅ CORRECT:
// internal/domain/repositories/user_repo.go (package repositories)

// UserRepository defines data access for the user entity.
type UserRepository interface {
    List() ([]*models.User, error)
    // GetUserByID returns the user, or (nil, nil) when it does not exist.
    GetByID(id uint) (*models.User, error)
}

type userRepo struct{}

// NewUserRepository returns a new UserRepository implementation.
func NewUserRepository() UserRepository {
    return &userRepo{}
}

func (r *userRepo) List() ([]*models.User, error) {
    var user models.User
    err := database.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        logger.Errorf("failed to get user by email: %v", err)
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    return &user, nil
}

// internal/app/services/auth/auth_service.go — the service depends on the INTERFACE
type AuthService struct {
    userRepo repositories.UserRepository
}

func NewAuthService(userRepo repositories.UserRepository) *AuthService {
    return &AuthService{userRepo: userRepo}
}

// internal/app/routers/index.go — wired once
authService := services.NewAuthService(repositories.NewUserRepository())

❌ WRONG:
// A service reaching for the database itself
func (s *AuthService) List(ctx context.Context) ([]*models.User, error) {
    var users []*models.User
    database.DB.Find(&users)   // ❌ only repositories may touch database.DB
    return users, nil
}

// A repository returning the raw GORM sentinel
func (r *userRepo) GetByID(id uint) (*models.User, error) {
    var user models.User
    return &user, database.DB.First(&user, id).Error  // ❌ leaks gorm.ErrRecordNotFound
}
```

> **DI rule:** the testability seam is the **interface**, not the connection. `database.DB`
> is opened once in `main.go` via `database.DbConnection(master, replica)` and used only
> inside `internal/domain/repositories` and `internal/adapters/database`. Services and
> controllers never reference it, which is exactly why swapping in `tests/mocks/` works
> without a database.

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
func (ctrl *AuthController) Register(c *gin.Context) {
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

// 200 OK (from internal/app/controllers/auth_controller.go)
func (ctrl *AuthController) List(c *gin.Context) {
    data, err := ctrl.service.List(c.Request.Context())
    if err != nil {
        utils.InternalServerError(c, err, "Failed to retrieve data")
        return
    }
    utils.Ok(c, data, "Data retrieved successfully")
}

// 201 Created
func (ctrl *AuthController) Create(c *gin.Context) {
    item, err := ctrl.service.Create(c.Request.Context(), &req)
    if err != nil {
        utils.BadRequest(c, err, "Failed to create item")
        return
    }
    utils.Created(c, item, "Item created successfully")
}

// 204 No Content
func (ctrl *AuthController) Delete(c *gin.Context) {
    if err := ctrl.service.Delete(c.Request.Context(), id); err != nil {
        utils.InternalServerError(c, err, "Failed to delete item")
        return
    }
    utils.NoContent(c)
}

❌ WRONG - Direct c.JSON():
func (ctrl *AuthController) List(c *gin.Context) {
    data, _ := ctrl.service.List(c.Request.Context())
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
func (ctrl *AuthController) Register(c *gin.Context) {
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
- **Domain/business errors:** Use sentinel errors in the service layer (e.g. `auth.ErrEmailAlreadyExists`, `auth.ErrInvalidCredentials`). Handlers use `errors.Is(err, auth.Err...)` (via a small `errToAPIError` mapper) to choose the right HTTP status and message. Keeps business rules in one place.
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

### 11.6 Routing and Route Registration

**⚠️ CRITICAL: routing lives in `internal/app/routers/`. Each feature gets its own
`<name>_routes.go` with a `Register<Name>Routes(group, service)` function, and
`index.go` is the single place that builds dependencies and calls them.**

**How routing is wired (three places):**

```
internal/app/routers/router.go        →  SetupRoute()          // gin engine + global middleware
internal/app/routers/index.go         →  RegisterRoutes(route) // builds repos+services, mounts groups
internal/app/routers/<name>_routes.go →  Register<Name>Routes  // one file per feature
```

**1. A feature declares its own routes:**

```go
// internal/app/routers/auth_routes.go (package routers)
func RegisterAuthRoutes(group *gin.RouterGroup, authService *services.AuthService) {
    authController := controllers.NewAuthController(authService)
    authRoutes := group.Group("/auth")
    authRoutes.POST("/login", authController.Login)
}
```

**2. `index.go` builds the dependencies and mounts everything:**

```go
// internal/app/routers/index.go
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

**3. `router.go` creates the engine and applies global middleware:**

```go
// internal/app/routers/router.go
func SetupRoute() *gin.Engine {
    route := gin.New()
    route.Use(middlewares.RequestID(), middlewares.RequestLog(), middlewares.CORS(), middlewares.Metrics())

    RegisterRoutes(route)
    registerSwagger(route) // debug only
    return route
}
```

**Rules:**

1. **One file per feature.** Put its routes in `internal/app/routers/<name>_routes.go`;
   never declare handlers inline in `index.go`.
2. **Business routes go under `/api/v1`** (the group passed to `Register<Name>Routes`),
   which carries the rate limiter.
3. **System routes (`/health`, `/metrics`) mount at ROOT** via `RegisterHealthRoutes(route)`,
   not under `/api/v1`.
4. **Protected routes use `middlewares.AuthMiddleware(authService)`**, applied to a group in
   `index.go`.
5. **Adding a feature touches `index.go` once**: build its repository and service, then call
   its `Register<Name>Routes`. Migrations are separate — versioned SQL, never derived from
   the models.

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

1. **ALWAYS** start a new feature by copying the `auth` slice
   - Copy each file and rename the prefix: `user_model.go` → `<name>_model.go`, and so on
   - Wire it in `internal/app/routers/index.go` and add `<name>_routes.go`
   - See `MODULE_GUIDE.md` — "How to add a new feature"

2. **ALWAYS** put each layer's file in its layer folder
   - `internal/domain/models/<name>_model.go`, `internal/domain/repositories/<name>_repo.go`,
     `internal/app/dto/<name>_dto.go`, `internal/app/services/<name>_service.go`,
     `internal/app/controllers/<name>_controller.go`, `internal/app/routers/<name>_routes.go`
   - The repository publishes its interface; the service depends on that interface,
     injected through its constructor (see §3.3)
   - Every schema change is a versioned `.up.sql` / `.down.sql` pair

3. **ALWAYS** check file size (>250 lines → split into another file in the same package)
   and function size (>80 lines → split into smaller functions)

4. **ALWAYS** add documentation
   - Package comment if a new package
   - Function comment for exported functions; struct comment for exported types

5. **ALWAYS** add error handling
   - Never ignore errors; always wrap with `%w` for context
   - Never use panic (except startup in `main.go` via `logger.Fatalf`)

6. **ALWAYS** write a test in `tests/unit/`
   - `tests/unit/services/<name>_service_test.go` using a fake repo (no DB), per
     `example/service_test.go`
   - At least happy path + 2 error cases

### 15.2 When Refactoring Code

1. **MUST** maintain backward compatibility
   - Don't change a published interface (`auth.AuthServicer`, `repositories.*Repository`)
     without discussion; other packages depend on it

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
