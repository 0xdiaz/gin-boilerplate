# Coding Standards & Guidelines

**Version:** 1.0
**Last Updated:** 2025-11-08
**Language:** Go 1.24+

---

## ⚠️ FOR AI AGENTS - READ THIS FIRST

> **🚨 CRITICAL: Read [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) BEFORE reading this document!**
>
> **This document is 1955 lines.** The critical rules file is 100 lines and contains NON-NEGOTIABLE patterns.
>
> **Don't make the same mistake:** Skipping critical sections = Code rejected.

### 🔥 Critical Sections in This Document (MUST READ)

After reading `00_AI_CRITICAL_RULES.md`, focus on these sections:

- **Lines 900-1100:** Struct-based Controller/Service Pattern (MANDATORY)
- **Lines 1429-1475:** Standard Response Format (MANDATORY)
- **Lines 1479-1584:** Response Utilities - pkg/utils/response.go (MANDATORY)
- **Lines 840-1009:** Testing Requirements (tests/ directory structure)

### 📖 How to Use This Document

1. ✅ Read `00_AI_CRITICAL_RULES.md` first (100 lines)
2. ✅ Read critical sections listed above
3. ⚠️  Skim other sections for context
4. 📚 Use as reference when needed

---

## 📋 Table of Contents

> **Quick Navigation:** Use Ctrl+F with line numbers (e.g., "L1429") or search keywords

### 🔥 Critical Sections (Read First)

| Section | Lines | Time | Keywords |
|---------|-------|------|----------|
| [1. File Organization](#1-file-organization) | L54-218 | 3 min | `300 lines`, `dto`, `enums`, `imports` |
| → [1.5 DTO & Enum Placement](#15-dto-and-enum-placement-rules) | L132-217 | 2 min | `internal/app/dto`, `pkg/enums`, `import cycle` |
| [5. Error Handling](#5-error-handling) | L562-711 | 3 min | `fmt.Errorf`, `%w`, `logger.Errorf` |
| [7. Testing Requirements](#7-testing-requirements) | L863-1036 | 4 min | `tests/`, `_test`, `70%`, `coverage` |
| [11. API Design](#11-api-design) | L1421-1953 | 8 min | `REST`, `response`, `utils.Ok`, `routes` |
| → [11.4 Response Utilities](#114-response-utilities) | L1547-1711 | 3 min | `utils.Ok`, `utils.Created`, `utils.BadRequest` |
| → [11.6 Router Organization](#116-router-organization) | L1727-1952 | 4 min | `Register*Routes`, `{feature}_routes.go` |

### 📚 All Sections

| # | Section | Lines | Keywords |
|---|---------|-------|----------|
| 1 | [File Organization](#1-file-organization) | L54-218 | `file size`, `directory`, `dto`, `enums` |
| 1.1 | [File Size Limits](#11-file-size-limits) | L56-67 | `300 lines`, `100 lines function` |
| 1.2 | [Directory Structure](#12-directory-structure) | L68-100 | `internal/`, `pkg/`, `cmd/` |
| 1.3 | [File Naming](#13-file-naming) | L101-115 | `snake_case`, `_test.go`, `naming` |
| 1.4 | [Package Organization](#14-package-organization) | L116-131 | `package`, `organization` |
| 1.5 | [DTO & Enum Placement](#15-dto-and-enum-placement-rules) | L132-217 | `dto`, `enums`, `import cycle` |
| 2 | [Naming Conventions](#2-naming-conventions) | L219-328 | `camelCase`, `PascalCase`, `naming` |
| 3 | [Code Structure](#3-code-structure) | L329-435 | `package`, `import`, `struct` |
| 4 | [Function Guidelines](#4-function-guidelines) | L436-561 | `function size`, `parameters`, `return` |
| 5 | [Error Handling](#5-error-handling) | L562-711 | `error`, `wrapping`, `logging` |
| 6 | [Documentation](#6-documentation) | L712-862 | `comment`, `godoc`, `docs` |
| 7 | [Testing Requirements](#7-testing-requirements) | L863-1036 | `test`, `tests/`, `coverage`, `70%` |
| 8 | [Logging Standards](#8-logging-standards) | L1037-1125 | `logger`, `Infof`, `Errorf`, `Warnf` |
| 9 | [Security Guidelines](#9-security-guidelines) | L1126-1263 | `validation`, `SQL injection`, `XSS` |
| 10 | [Database Access](#10-database-access) | L1264-1420 | `GORM`, `model`, `migration` |
| 11 | [API Design](#11-api-design) | L1421-1953 | `REST`, `API`, `response`, `routes` |
| 11.1 | [RESTful Conventions](#111-restful-conventions) | L1423-1446 | `REST`, `HTTP methods`, `endpoints` |
| 11.2 | [HTTP Status Codes](#112-http-status-codes) | L1447-1493 | `200`, `404`, `500`, `status` |
| 11.3 | [Response Format](#113-response-format) | L1494-1546 | `success`, `message`, `data`, `errors` |
| 11.4 | [Response Utilities](#114-response-utilities) | L1547-1711 | `utils.Ok`, `utils.Created`, `HandleErrors` |
| 11.5 | [API Versioning](#115-api-versioning) | L1712-1726 | `v1`, `v2`, `versioning` |
| 11.6 | [Router Organization](#116-router-organization) | L1727-1952 | `routes`, `Register*Routes`, `index.go` |
| 12 | [Configuration](#12-configuration) | L1955-2035 | `config`, `env`, `.env` |
| 13 | [Forbidden Practices](#13-forbidden-practices) | L2036-2118 | `panic`, `global`, `anti-pattern` |
| 14 | [Code Review Checklist](#14-code-review-checklist) | L2119-2171 | `checklist`, `review`, `validation` |
| 15 | [AI Agent Rules](#15-ai-agent-specific-rules) | L2172-2221 | `AI`, `agent`, `automation` |
| 16 | [Resources](#16-resources) | L2222-2226 | `links`, `references`, `docs` |

### 🎯 Quick Lookups by Task

**Creating Controllers/Services:**
- L132-157: DTO Placement
- L1547-1711: Response Utilities (MUST use)
- See: `DESIGN_PATTERNS.md` L900-1016

**Adding Routes:**
- L1727-1952: Router Organization
- L1423-1446: RESTful Conventions

**Writing Tests:**
- L863-1036: Testing Requirements
- L56-67: File Size Limits

**Error Handling:**
- L562-711: Error Handling Patterns
- L1037-1125: Logging Standards

**Database Models:**
- L1264-1420: Database Access
- See: `DESIGN_PATTERNS.md` L1017-1085

---

## 1. FILE ORGANIZATION

### 1.1 File Size Limits

```
✅ MUST: Maximum 300 lines per file
⚠️  WARNING: 300-500 lines requires justification
❌ FORBIDDEN: >500 lines in a single file
```

**Action when exceeding:**
- Split into multiple focused files
- **Subfolder required when splitting:** When one feature (service, controller, or repository) is split into more than one file, move **all** files for that feature into a new subfolder named after the feature. The parent folder should only contain features that are still single-file. Example: `auth_service.go` + `auth_service_tokens.go` → move to `services/auth/auth_service.go` + `services/auth/auth_service_tokens.go` (package `auth`). This keeps split features grouped and the structure clean.
- Another example: `controllers/user_controller.go` + `controllers/user_validation.go` → `controllers/user/user_controller.go` + `controllers/user/user_validation.go` (package `user`).

### 1.2 Directory Structure

**MUST follow this structure:**

```
project/
├── cmd/                        # Commands (main applications)
│   ├── migrate/               # Database migrations
│   └── seeder/                # Data seeding
├── internal/                   # Private application code
│   ├── adapters/              # External adapters (DB, cache, etc)
│   ├── app/
│   │   ├── controllers/       # HTTP handlers (thin layer)
│   │   ├── dto/               # Data Transfer Objects
│   │   ├── middlewares/       # Gin middlewares
│   │   ├── routers/           # Route definitions
│   │   └── services/          # Business logic; split features use subfolder (e.g. services/auth/)
│   └── domain/
│       ├── models/            # Database entities
│       └── repositories/      # Data access layer
├── pkg/                       # Public reusable packages
│   ├── config/               # Configuration management
│   ├── enums/                # Enumerations
│   ├── logger/               # Logging infrastructure
│   ├── types/                # Shared types
│   └── utils/                # Utility functions
├── tests/                    # Test files (NEW - REQUIRED)
│   ├── unit/
│   ├── integration/
│   └── fixtures/
└── main.go                   # Application entry point
```

#### 1.2.1 Feature Grouping for Multi-file Controllers/Services

```
✅ MUST: If a controller or service is split into multiple files (due to the 300-line/file or 100-line/function limits), those files MUST be grouped under a single feature folder to avoid scattering.
```

**Purpose:**
- Keep files organized by feature
- Improve navigation and code ownership
- Avoid mixing many feature files in a flat directory

**Examples:**
```
internal/app/controllers/
├── auth/
│   ├── auth_controller.go
│   ├── auth_controller_handlers.go
│   └── auth_controller_validators.go
├── user/
│   ├── user_controller.go
│   └── user_controller_helpers.go
└── health_controller.go

internal/app/services/
├── auth/
│   ├── auth_service.go
│   ├── auth_service_reset.go
│   └── auth_service_validation.go
└── user/
    ├── user_service.go
    └── user_service_rules.go
```

**Naming Rules:**
- Folder name must be the feature: `auth/`, `user/`, `order/`
- All files in the folder use `snake_case` and keep a consistent `package` (`package controllers` or `package services`)
- Avoid deep nested folders unless necessary (max 1 level for feature grouping)

### 1.3 File Naming

**Rules:**
```go
✅ CORRECT:
- user_service.go          (snake_case)
- transaction_repo.go      (snake_case, abbreviated)
- api_client.go            (lowercase acronyms in filename)

❌ WRONG:
- UserService.go           (PascalCase not allowed)
- user-service.go          (kebab-case not allowed)
- userService.go           (camelCase not allowed)
```

### 1.4 Package Organization

**One package per directory:**
```go
// ✅ CORRECT
internal/app/services/
    ├── user_service.go        → package services
    ├── transaction_service.go → package services
    └── hub_service.go         → package services

// ❌ WRONG - multiple packages in one directory
internal/app/
    ├── user.go     → package users
    └── client.go   → package clients
```

### 1.5 DTO and Enum Placement Rules

**CRITICAL RULES - MUST ALWAYS FOLLOW:**

#### 1.5.1 DTOs (Data Transfer Objects)
```
✅ ALL STRUCTS MUST GO TO: internal/app/dto/

Location rules:
- Request/Response structs → dto/
- API payload structs → dto/
- Service-specific structs → dto/
- Any struct used for data transfer → dto/

Examples:
✅ CORRECT:
internal/app/dto/
    ├── notification_dto.go → NotificationRequest, NotificationResponse
    ├── user_dto.go         → CreateUserRequest, UserResponse
    ├── order_dto.go        → OrderRequest, OrderResponse

❌ WRONG:
internal/app/services/notification_types.go → NotificationRoute (should be in dto/)
internal/app/services/user_types.go         → UserRequest (should be in dto/)
pkg/types/request.go                        → ApiRequest (should be in dto/)
```

#### 1.5.2 Enums and Constants
```
✅ ALL CONSTANTS MUST GO TO: pkg/enums/

Location rules:
- String constants → pkg/enums/
- Numeric constants → pkg/enums/
- Array/slice constants → pkg/enums/
- Status codes → pkg/enums/
- Configuration values (if constant) → pkg/enums/

⚠️  EXCEPTION - Service-specific mappings:
- Map[string]Struct using DTOs → stays in services/ (to avoid import cycles)
- Example: ServiceRouteMapping → internal/app/services/

Examples:
✅ CORRECT:
pkg/enums/
    ├── payment_constants.go → PaymentMethodCard, PaymentStatusSuccess
    ├── status.go            → StatusActive, StatusInactive
    ├── role.go              → RoleAdmin, RoleUser, RoleModerator

internal/app/services/
    └── notification_router.go → NotificationRouteMapping (map using dto.NotificationRoute)

❌ WRONG:
internal/app/services/hub_types.go         → EndpointBalance (should be in enums/)
internal/app/controllers/constants.go      → StatusCodes (should be in enums/)
pkg/enums/hub_constants.go                 → ServiceRouteMapping (causes import cycle)
```

#### 1.5.3 Import Cycle Prevention

**CRITICAL: pkg/ packages CANNOT import from internal/**

```go
✅ ALLOWED:
internal/app/dto → pkg/enums        ✓
internal/app/services → pkg/enums   ✓
internal/app/services → dto          ✓

❌ FORBIDDEN (causes import cycle):
pkg/enums → internal/app/dto         ✗
pkg/enums → internal/app/services    ✗
pkg/utils → internal/app/models      ✗

Solution for maps using DTOs:
// ❌ WRONG - causes import cycle
// pkg/enums/notification_constants.go
import "internal/app/dto"
var NotificationRouteMapping = map[string]dto.NotificationRoute{...}

// ✅ CORRECT - keep in services
// internal/app/services/notification_router.go
import "internal/app/dto"
var NotificationRouteMapping = map[string]dto.NotificationRoute{...}
```

---

## 2. NAMING CONVENTIONS

### 2.1 Package Names

```go
✅ MUST:
- Lowercase, single word
- No underscores, no dashes
- Descriptive but concise

✅ CORRECT:
package controllers
package services
package repositories
package utils

❌ WRONG:
package user_services    // No underscores
package UserServices     // No capitals
package svc             // Too cryptic
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

### 3.1 Layered Architecture

**MUST follow this flow:**

```
HTTP Request
    ↓
[Router] → routes request to controller
    ↓
[Controller] → thin layer, only handles HTTP concerns
    ↓         (parse params, call service, return response)
[Service] → business logic layer (FAT)
    ↓         (validation, transformation, orchestration)
[Repository] → data access layer
    ↓         (CRUD operations only)
[Model] → database entity
```

**Rules:**
```go
✅ Controller SHOULD:
- Parse request parameters
- Call service methods
- Format HTTP responses
- Handle HTTP-specific errors

❌ Controller MUST NOT:
- Contain business logic
- Access database directly
- Perform complex transformations
- Have more than 50 lines per function

✅ Service SHOULD:
- Contain all business logic
- Orchestrate multiple repositories
- Validate business rules
- Transform data between layers

❌ Service MUST NOT:
- Handle HTTP concerns (gin.Context)
- Import controller packages
- Exceed 400 lines per file

✅ Repository SHOULD:
- Perform CRUD operations
- Build database queries
- Handle transaction management
- Return domain models

❌ Repository MUST NOT:
- Contain business logic
- Call other repositories directly
- Import service packages
```

### 3.2 Dependency Direction

**MUST follow:**
```
Controllers → Services → Repositories → Models
     ↓            ↓            ↓
    DTOs       Utils        Database
```

**FORBIDDEN:**
```
❌ Models importing Services
❌ Repositories importing Services
❌ Services importing Controllers
❌ Circular dependencies
```

### 3.3 Single Responsibility Principle

**Each file should have ONE clear purpose:**

```go
✅ CORRECT:
// user_service.go - User business logic
type UserService struct {
    userRepo repositories.UserRepository
}

func (s *UserService) CreateUser(dto dto.CreateUserRequest) error
func (s *UserService) GetUserByID(id uint) (*models.User, error)
func (s *UserService) UpdateUser(id uint, dto dto.UpdateUserRequest) error
func (s *UserService) DeleteUser(id uint) error

❌ WRONG:
// god_service.go - Too many responsibilities
type GodService struct {
    userRepo    repositories.UserRepository
    clientRepo  repositories.ClientRepository
    transRepo   repositories.TransactionRepository
}

func (s *GodService) CreateUser() error
func (s *GodService) ProcessTransaction() error
func (s *GodService) SendEmail() error
func (s *GodService) GenerateReport() error
func (s *GodService) ValidatePayment() error
```

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
// 1. Type definitions
type UserService struct {
    repo repositories.UserRepository
}

// 2. Constructor
func NewUserService(repo repositories.UserRepository) *UserService {
    return &UserService{repo: repo}
}

// 3. Public methods (exported) - alphabetical order
func (s *UserService) CreateUser(dto dto.CreateUserRequest) error { }
func (s *UserService) DeleteUser(id uint) error { }
func (s *UserService) GetUserByID(id uint) (*models.User, error) { }
func (s *UserService) UpdateUser(id uint, dto dto.UpdateUserRequest) error { }

// 4. Private methods (unexported) - alphabetical order
func (s *UserService) buildUser(dto dto.CreateUserRequest) *models.User { }
func (s *UserService) validateUserData(user *models.User) error { }

// 5. Standalone helper functions
func generateUserID() string { }
func hashPassword(password string) (string, error) { }
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

### 5.3 Custom Error Types

**Define custom errors in `/pkg/types/errors.go`:**

```go
✅ CORRECT:
// pkg/types/errors.go
package types

import "errors"

var (
    ErrNotFound          = errors.New("resource not found")
    ErrUnauthorized      = errors.New("unauthorized access")
    ErrInvalidInput      = errors.New("invalid input data")
    ErrDuplicateEntry    = errors.New("duplicate entry")
    ErrExternalAPIFailed = errors.New("external API call failed")
)

// Usage in code
if user == nil {
    return types.ErrNotFound
}

// With context
if user == nil {
    return fmt.Errorf("user %d: %w", id, types.ErrNotFound)
}
```

### 5.4 Panic Usage

```
❌ FORBIDDEN in business logic (controllers, services, repositories)
⚠️  ACCEPTABLE only in:
    - Application initialization (main.go)
    - Configuration validation (startup only)
    - Fatal unrecoverable errors (database connection failure at startup)

✅ ALWAYS use defer/recover in HTTP handlers
```

```go
✅ CORRECT:
// main.go - acceptable for fatal startup errors
func main() {
    db, err := database.Connect()
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)  // OK at startup
    }

    logger.Init()
    // ...
}

// controller.go - never panic, always return error
func GetUser(c *gin.Context) {
    defer func() {
        if r := recover(); r != nil {
            logger.Errorf("Panic recovered: %v", r)
            c.JSON(500, gin.H{"error": "Internal server error"})
        }
    }()

    // Handler code
}

❌ FORBIDDEN:
// service.go - NEVER panic in business logic
func CreateUser(user *models.User) error {
    if user.Email == "" {
        panic("email is required")  // ❌ NEVER DO THIS!
    }
    // ...
}
```

### 5.5 Error Strategy (HTTP Response and Domain Errors)

**HTTP responses:** All error responses to the client MUST go through `pkg/utils`: `utils.BadRequest`, `utils.Unauthorized`, `utils.RespondWithAPIError`, etc. Do not call `c.JSON(code, ...)` directly for errors.

**Domain errors:** Services return domain-specific errors (e.g. `auth.ErrEmailAlreadyExists`, `auth.ErrInvalidCredentials`). Controllers map these to HTTP using either:
- `utils.RespondWithAPIError(c, apiErr)` when the error is mapped to a `*types.APIError` (code, message, details), or
- `utils.InternalServerError(c, err, "…")` (and similar) for unmapped or generic errors.

**pkg/types/errors:** `types.APIError` and the predefined `Err*` values are used when mapping domain errors to HTTP. Map domain errors (e.g. with `errors.Is`) to `*types.APIError` in the controller, then call `utils.RespondWithAPIError`. Do not duplicate status-code logic in multiple places.

**Validation errors:** Keep using `utils.BadRequest(c, err, "Invalid request data")` for `validator.ValidationErrors`; no need to convert to `types.APIError`.

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
// Package services implements the business logic layer of the application.
//
// This package contains all business rules, validation logic, and orchestration
// between repositories. Each service should focus on a single domain entity.
//
// Example usage:
//
//	userService := services.NewUserService(userRepo)
//	user, err := userService.GetUserByID(123)
package services
```

**File to add package comment:**
- Choose the main file or create `doc.go`
- Usually: `user_service.go`, `transaction_service.go`, or `doc.go`

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
// UserService handles all user-related business logic.
type UserService struct {
    userRepo repositories.UserRepository
    logger   *logrus.Logger
}

// CreateUserRequest represents the data required to create a new user.
type CreateUserRequest struct {
    // Name is the full name of the user (required, max 255 chars)
    Name string `json:"name" validate:"required,max=255"`

    // Email must be unique across all users (required, valid email format)
    Email string `json:"email" validate:"required,email"`

    // Phone in international format with country code (optional)
    Phone string `json:"phone" validate:"omitempty,e164"`
}

❌ WRONG:
// User service
type UserService struct {  // Too brief
    userRepo repositories.UserRepository
}

type CreateUserRequest struct {  // No struct comment!
    Name  string `json:"name"`
    Email string `json:"email"`
    Phone string `json:"phone"`
}
```

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

### 7.1 Test File Location

**⚠️ CRITICAL: ALL test files MUST be in the `tests/` directory, NOT co-located with source code.**

```
✅ CORRECT:
tests/unit/services/auth_service_test.go
tests/unit/controllers/auth_controller_test.go
tests/integration/api/user_api_test.go
tests/fixtures/users.json

❌ WRONG:
internal/app/services/auth_service_test.go     // Co-located with source
internal/app/controllers/auth_controller_test.go
internal/domain/repositories/user_repo_test.go
```

**Package naming for test files:**

```go
✅ CORRECT:
// tests/unit/services/auth_service_test.go
package services_test  // External test package with _test suffix

import (
    "testing"
    "your-project/internal/app/services"  // Import the package being tested
)

❌ WRONG:
// tests/unit/services/auth_service_test.go
package services  // Missing _test suffix
```

**Why separate tests directory:**
- ✅ Cleaner project structure
- ✅ Easier to run all tests (`go test ./tests/...`)
- ✅ Clear separation of concerns
- ✅ Better organization for large projects
- ✅ Explicit external testing (tests public API only)

### 7.2 Test Coverage Requirements

```
✅ MINIMUM: 70% code coverage for services layer
✅ MINIMUM: 50% code coverage for repositories layer
✅ MINIMUM: 60% code coverage for utilities
⚠️  OPTIONAL: Controllers (covered by integration tests)
```

### 7.3 Test File Naming

```go
✅ CORRECT:
user_service.go       → user_service_test.go
transaction_repo.go   → transaction_repo_test.go
helper.go             → helper_test.go

❌ WRONG:
user_service.go → user_test.go
user_service.go → test_user_service.go
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

### 7.5 Test Structure (Table-Driven Tests)

**MUST use table-driven tests for multiple scenarios:**

```go
✅ CORRECT:
func TestCreateUser(t *testing.T) {
    tests := []struct {
        name    string
        input   dto.CreateUserRequest
        want    *models.User
        wantErr bool
    }{
        {
            name: "valid user",
            input: dto.CreateUserRequest{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            want: &models.User{
                Name:  "John Doe",
                Email: "john@example.com",
            },
            wantErr: false,
        },
        {
            name: "duplicate email",
            input: dto.CreateUserRequest{
                Name:  "Jane Doe",
                Email: "existing@example.com",
            },
            want:    nil,
            wantErr: true,
        },
        {
            name: "invalid email",
            input: dto.CreateUserRequest{
                Name:  "Invalid User",
                Email: "not-an-email",
            },
            want:    nil,
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := service.CreateUser(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("CreateUser() error = %v, wantErr %v", err, tt.wantErr)
                return
            }

            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("CreateUser() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### 7.6 Test Organization

**MUST follow this structure in `tests/` directory:**

```
tests/
├── unit/                  # Unit tests (isolated, mocked dependencies)
│   ├── controllers/       # Controller tests
│   │   └── auth_controller_test.go
│   ├── services/          # Service layer tests
│   │   ├── user_service_test.go
│   │   └── auth_service_test.go
│   ├── repositories/      # Repository tests
│   │   └── user_repo_test.go
│   └── utils/             # Utility tests
│       └── helper_test.go
├── integration/           # Integration tests (real DB, external services)
│   ├── api/              # End-to-end API tests
│   │   ├── user_api_test.go
│   │   └── auth_api_test.go
│   └── database/         # Database integration tests
│       └── migration_test.go
├── fixtures/             # Test data
│   ├── users.json        # Sample user data
│   └── transactions.json # Sample transaction data
└── README.md            # Testing documentation
```

---

## 8. LOGGING STANDARDS

**MUST use the structured logger from `pkg/logger`.** For request-scoped tracing (START/FINISH with request_id and duration), use **LogStart** and **LogFinish** in every HTTP handler and service method; see [OBSERVABILITY.md](OBSERVABILITY.md) — section "Logger API (pkg/logger)" and "Using Request ID in Code" for full usage and examples.

### 8.1 Logger Usage

**MUST use structured logger from `pkg/logger`:**

```go
✅ CORRECT:
import "your-project/pkg/logger"

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

**Request-scoped tracing (handlers and services):** At the start of each handler call `ctx, start := logger.LogStart(c.Request.Context(), "ControllerName.Method")` and pass `ctx` to services. In each service method call `ctx, start := logger.LogStart(ctx, "ServiceName.Method")`. Before **every** return (success or error), call `logger.LogFinish(ctx, "ControllerName.Method", err, start)` or `logger.LogFinish(ctx, "ServiceName.Method", err, start)`. This produces START/FINISH lines with request_id and duration. For details and examples see **OBSERVABILITY.md**.

### 8.2 Log Levels

```
✅ DEBUG: Detailed debugging information (disabled in production)
✅ INFO:  Normal operations, state changes, successful operations
✅ WARN:  Unexpected situations that don't prevent operation
✅ ERROR: Operation failed, error occurred but system continues
❌ FATAL: ONLY in main.go for unrecoverable startup errors
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
import "your-project/pkg/utils"

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

**MUST validate ALL external input:**

```go
✅ CORRECT:
import "github.com/go-playground/validator/v10"

type CreateUserRequest struct {
    Name  string `json:"name" validate:"required,min=3,max=255"`
    Email string `json:"email" validate:"required,email"`
    Age   int    `json:"age" validate:"gte=0,lte=150"`
}

func CreateUser(dto CreateUserRequest) error {
    validate := validator.New()
    if err := validate.Struct(dto); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    // Process...
}

❌ WRONG:
func CreateUser(dto CreateUserRequest) error {
    // No validation - directly using input!
    user := &models.User{
        Name:  dto.Name,
        Email: dto.Email,
    }
    return repo.Create(user)
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

**MUST check permissions:**

```go
✅ CORRECT:
func DeleteUser(c *gin.Context) {
    // Get authenticated user
    currentUser := c.MustGet("user").(*models.User)

    // Check permission
    targetUserID := c.Param("id")
    if currentUser.Role != models.RoleAdmin && currentUser.ID != targetUserID {
        utils.HandleErrors(c, http.StatusForbidden, nil, "Unauthorized")
        return
    }

    // Proceed with deletion
}

❌ WRONG:
func DeleteUser(c *gin.Context) {
    // No authentication check!
    // No authorization check!
    id := c.Param("id")
    service.DeleteUser(id)
}
```

### 9.5 Rate Limiting

**MUST implement rate limiting for public endpoints (e.g. auth):**

```go
✅ CORRECT:
// In index.go — applied once on the apiV1 group; all sub-groups (auth, example, ...) inherit it.
// middleware reads RATE_LIMIT_RPS and RATE_LIMIT_BURST from config
apiV1.Use(middlewares.RateLimitMiddleware())
```

```go
❌ WRONG:
// Do NOT call RateLimitMiddleware() again on a sub-group (e.g. authRoutes) — it shares the
// same singleton limiter/key as apiV1, so this silently halves the effective burst capacity
// for that sub-group instead of adding real protection.
authRoutes.Use(middlewares.RateLimitMiddleware())
```

Limits are read inside `RateLimitMiddleware()` from env vars `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST` (defaults 100 rps, 200 burst if unset). See [CONFIGURATION.md](CONFIGURATION.md) — Environment Variables (Optional) and Example 3. If a route group genuinely needs a stricter limit (e.g. brute-force protection on login), use `RateLimitMiddlewareWithConfig(rps, burst)` with its own limiter instance instead of calling `RateLimitMiddleware()` twice.

---

## 10. DATABASE ACCESS

### 10.1 Repository Pattern

**MUST use repository pattern:**

```go
✅ CORRECT:
// internal/domain/repositories/user_repository.go
type UserRepository interface {
    FindByID(id uint) (*models.User, error)
    FindByEmail(email string) (*models.User, error)
    Create(user *models.User) error
    Update(user *models.User) error
    Delete(id uint) error
    List(page, pageSize int) ([]*models.User, int64, error)
}

// internal/adapters/database/user_repo_impl.go
type userRepositoryImpl struct {
    db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repositories.UserRepository {
    return &userRepositoryImpl{db: db}
}

❌ WRONG:
// service.go - Direct database access!
func (s *UserService) GetUser(id uint) (*models.User, error) {
    var user models.User
    if err := database.DB.First(&user, id).Error; err != nil {  // ❌ Direct DB access
        return nil, err
    }
    return &user, nil
}
```

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

// Example usage
func CreateUser(c *gin.Context) {
    var dto dto.CreateUserRequest
    if err := c.ShouldBindJSON(&dto); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    user, err := service.CreateUser(dto)
    if err != nil {
        if errors.Is(err, types.ErrDuplicateEntry) {
            c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal error"})
        return
    }

    c.JSON(http.StatusCreated, user)
}

❌ WRONG:
// Always returning 200 OK
c.JSON(200, gin.H{"success": false, "error": "User not found"})

// Wrong status for operation
c.JSON(404, gin.H{"message": "User created"})  // Should be 201
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

// 200 OK
func GetUser(c *gin.Context) {
    user, err := service.GetUserByID(id)
    if err != nil {
        utils.NotFound(c, err, "User not found")
        return
    }
    utils.Ok(c, user, "User retrieved successfully")
}

// 201 Created
func CreateUser(c *gin.Context) {
    user, err := service.CreateUser(dto)
    if err != nil {
        utils.BadRequest(c, err, "Failed to create user")
        return
    }
    utils.Created(c, user, "User created successfully")
}

// 204 No Content
func DeleteUser(c *gin.Context) {
    if err := service.DeleteUser(id); err != nil {
        utils.InternalServerError(c, err, "Failed to delete user")
        return
    }
    utils.NoContent(c)
}

❌ WRONG - Direct c.JSON():
func GetUser(c *gin.Context) {
    user, _ := service.GetUserByID(id)
    c.JSON(200, gin.H{"data": user})  // Inconsistent format!
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

When an endpoint **must** follow an external format (e.g. DataTables server-side, third-party protocol), the response may be sent using that library’s helper as long as it is consistent and documented. Example: `ExampleController.GetDataDatatables` uses `datatables.JSON(c, data)` for DataTables format, not `utils.Ok`. Error paths must still use response utils.

#### Service accepting gin.Context exception (DataTables)

General rule: services **MUST NOT** accept or use `gin.Context`. Exception: when a third-party library integration (e.g. Datatables-Gin) **requires** `*gin.Context` for binding query params (draw, start, length, search, order), the service method that calls that library may accept `(ctx context.Context, c *gin.Context)` to avoid duplicating the entire library API into DTOs. Example: `ExampleService.GetDataDatatables(ctx, c)` and `repositories.GetDataDatatables(c)` — this exception is documented in [SERVICE_COMPLIANCE_AUDIT.md](SERVICE_COMPLIANCE_AUDIT.md). Other services must not use `gin.Context`.

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
func CreateUser(c *gin.Context) {
    var dto dto.CreateUserRequest
    if err := c.ShouldBindJSON(&dto); err != nil {
        // Automatically formats validator.ValidationErrors
        utils.BadRequest(c, err, "Validation failed")
        return
    }
    // ...
}

// Response automatically formatted as:
{
    "success": false,
    "message": "Validation failed",
    "data": null,
    "errors": {
        "email": "email must be a valid email address",
        "name": "name is required",
        "age": "age must be at least 18"
    }
}
```

#### Error strategy (layers)

- **HTTP response:** All HTTP responses MUST go through `pkg/utils` (e.g. `utils.Ok`, `utils.BadRequest`, `utils.HandleErrors`). Controllers map service/domain errors to these utilities; no direct `c.JSON` for error/success.
- **Domain/business errors:** Use typed errors in the service layer (e.g. `auth.ErrEmailAlreadyExists`, `auth.ErrInvalidCredentials`). Controllers use `errors.Is(err, auth.Err...)` to choose the right HTTP status and message. Keeps business rules in one place.
- **`pkg/types/errors`:** `APIError`, `ErrNotFound`, etc. are for specific use cases (e.g. client SDK, internal mapping). Do not duplicate with domain errors; prefer domain errors in services and map to utils in controllers. Document in code when using `types.APIError` for a given flow.

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

### 11.6 Router Organization

**⚠️ CRITICAL: Router files MUST be organized by controller/feature, NOT all in one file.**

**File Structure:**

```
internal/app/routers/
├── index.go           # Main router - registers all route groups
├── auth_routes.go     # Authentication routes
├── user_routes.go     # User management routes
├── product_routes.go  # Product routes
└── order_routes.go    # Order routes
```

**Rules:**

```go
✅ CORRECT:
1. One route file per controller/feature
2. File name pattern: {feature}_routes.go
3. Each route file has a Register function
4. Main router calls all Register functions

// internal/app/routers/auth_routes.go
package routers

import (
    "github.com/your-org/project/internal/app/controllers"
    "github.com/your-org/project/internal/app/middlewares"
    "github.com/your-org/project/internal/app/services"
    "github.com/gin-gonic/gin"
)

// RegisterAuthRoutes registers authentication routes
func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
    authController := controllers.NewAuthController(authService)

    authRoutes := router.Group("/auth")
    {
        authRoutes.POST("/register", authController.Register)
        authRoutes.POST("/login", authController.Login)
        authRoutes.POST("/refresh", authController.RefreshToken)
    }
}

// internal/app/routers/user_routes.go
package routers

import (
    "github.com/your-org/project/internal/app/controllers"
    "github.com/your-org/project/internal/app/middlewares"
    "github.com/your-org/project/internal/app/services"
    "github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers user management routes
func RegisterUserRoutes(router *gin.Engine, userService *services.UserService, authService *services.AuthService) {
    userController := controllers.NewUserController(userService)

    // Protected routes - require authentication
    userRoutes := router.Group("/api/v1/users")
    userRoutes.Use(middlewares.AuthMiddleware(authService))
    {
        userRoutes.GET("", userController.GetAll)
        userRoutes.GET("/:id", userController.GetByID)
        userRoutes.PUT("/:id", userController.Update)
        userRoutes.DELETE("/:id", userController.Delete)
    }
}

// internal/app/routers/index.go
package routers

import (
    "net/http"

    "github.com/your-org/project/internal/app/services"
    "github.com/gin-gonic/gin"
)

// RegisterRoutes registers all application routes
func RegisterRoutes(router *gin.Engine) {
    // 404 handler
    router.NoRoute(func(ctx *gin.Context) {
        ctx.JSON(http.StatusNotFound, gin.H{
            "status":  http.StatusNotFound,
            "message": "Route Not Found",
        })
    })

    // Health check (public)
    router.GET("/health", func(ctx *gin.Context) {
        ctx.JSON(http.StatusOK, gin.H{"live": "ok"})
    })

    // Initialize services
    authService := services.NewAuthService()
    userService := services.NewUserService()
    productService := services.NewProductService()

    // Register route groups
    RegisterAuthRoutes(router, authService)
    RegisterUserRoutes(router, userService, authService)
    RegisterProductRoutes(router, productService, authService)
    RegisterOrderRoutes(router, orderService, authService)
}

❌ WRONG:
// All routes in index.go (BAD - file akan jadi terlalu panjang)
package routers

func RegisterRoutes(router *gin.Engine) {
    // Auth routes
    authRoutes := router.Group("/auth")
    {
        authRoutes.POST("/register", ...)
        authRoutes.POST("/login", ...)
    }

    // User routes
    userRoutes := router.Group("/users")
    {
        userRoutes.GET("", ...)
        userRoutes.GET("/:id", ...)
        userRoutes.PUT("/:id", ...)
        // ... 50+ more routes ...
    }

    // Product routes
    productRoutes := router.Group("/products")
    {
        // ... 30+ more routes ...
    }

    // This file will become 500+ lines! ❌
}
```

**Route Organization Principles:**

1. **Separation by Feature**
   - One file per feature/controller
   - Max 100 lines per route file
   - Clear naming: `{feature}_routes.go`

2. **Consistent Pattern**
   ```go
   // Function naming
   RegisterAuthRoutes()      ✅
   RegisterUserRoutes()      ✅
   AuthRoutes()              ❌ (missing Register prefix)
   SetupAuthRoutes()         ❌ (use Register)
   InitAuthRoutes()          ❌ (use Register)
   ```

3. **Dependency Injection**
   ```go
   ✅ CORRECT:
   func RegisterUserRoutes(router *gin.Engine, userService *services.UserService) {
       userController := controllers.NewUserController(userService)
       // ...
   }

   ❌ WRONG:
   func RegisterUserRoutes(router *gin.Engine) {
       // Creating services inside route function
       userService := services.NewUserService()  // Should be passed as param
   }
   ```

4. **Route Grouping**
   ```go
   ✅ CORRECT:
   // Group related routes
   authRoutes := router.Group("/auth")
   {
       authRoutes.POST("/register", authController.Register)
       authRoutes.POST("/login", authController.Login)
   }

   // Protected group
   protectedRoutes := router.Group("/api/v1")
   protectedRoutes.Use(middlewares.AuthMiddleware(authService))
   {
       protectedRoutes.GET("/profile", userController.GetProfile)
   }

   ❌ WRONG:
   // No grouping - repetitive; also prefer versioned path /api/v1/auth
   router.POST("/auth/register", authController.Register)
   router.POST("/auth/login", authController.Login)
   router.POST("/auth/refresh", authController.RefreshToken)
   ```

5. **File Size Limit**
   - Single route file MAX 100 lines
   - If exceeding, split by sub-feature
   - Example: `user_routes.go` → `user_admin_routes.go` + `user_public_routes.go`

**Why This Matters:**

- ✅ **Maintainability**: Easy to find routes for specific feature
- ✅ **Scalability**: Add new features without modifying existing files
- ✅ **Readability**: Each file focuses on one responsibility
- ✅ **Team Work**: Multiple developers can work on different route files
- ✅ **File Size**: Prevents index.go from becoming 1000+ lines
- ❌ **Without organization**: index.go becomes unmaintainable mess

**Example Project Structure:**

```
internal/app/routers/
├── index.go              # 50 lines - main router
├── auth_routes.go        # 40 lines - auth endpoints
├── user_routes.go        # 60 lines - user CRUD
├── product_routes.go     # 70 lines - product CRUD
├── order_routes.go       # 80 lines - order management
├── payment_routes.go     # 55 lines - payment processing
└── admin_routes.go       # 45 lines - admin panel routes

Total: ~400 lines across 7 files (average 57 lines per file)
vs
Single file: 400 lines (unmaintainable)
```

---

## 12. CONFIGURATION

### 12.1 Environment Variables

**MUST use environment variables for configuration:**

```go
✅ CORRECT:
// .env file
MASTER_DB_HOST=localhost
MASTER_DB_PORT=5432
MASTER_DB_USER=postgres
MASTER_DB_PASSWORD=secret
MASTER_DB_NAME=project_db

JWT_SECRET=your-secret-key
API_TIMEOUT=30s
RATE_LIMIT=100

// config.go
type Config struct {
    Database struct {
        Host     string
        Port     int
        User     string
        Password string
        Name     string
    }
    JWT struct {
        Secret string
    }
    API struct {
        Timeout    time.Duration
        RateLimit  int
    }
}

❌ WRONG:
// Hardcoded in code
const (
    DBHost     = "localhost"      // ❌ Hardcoded
    DBPassword = "secret123"      // ❌ NEVER commit passwords!
    JWTSecret  = "my-secret-key"  // ❌ Security risk
)
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

1. **ALWAYS** check file size before adding code
   - If file >250 lines, create new file instead

2. **ALWAYS** check function size
   - If function >80 lines, split into smaller functions

3. **ALWAYS** add documentation
   - Package comment if new package
   - Function comment for exported functions
   - Struct comment for exported types

4. **ALWAYS** add error handling
   - Never ignore errors
   - Always wrap errors with context
   - Never use panic (except main.go startup)

5. **ALWAYS** write tests
   - Test file for every new service/repository
   - At least happy path + 2 error cases

### 15.2 When Refactoring Code

1. **MUST** maintain backward compatibility
   - Don't change public API without discussion

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
- Assertions: `github.com/stretchr/testify`
- Mocking: `github.com/stretchr/testify/mock`
- HTTP testing: `httptest`

---

**END OF CODING STANDARDS**

*Last updated: 2025-11-08*
*Version: 1.0*
*Review: Every 3 months or after major changes*
