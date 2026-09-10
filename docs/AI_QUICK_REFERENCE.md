# AI Agent Quick Reference

> **Print this mentally before every code change!**

> 🧭 **Structure:** Code is organized by technical layer —
> `internal/app/controllers` → `internal/app/services` → `internal/domain/repositories` →
> `internal/domain/models`, wired in `internal/app/routers/index.go`. See
> **[MODULE_GUIDE.md](./MODULE_GUIDE.md)** (the source of truth for the layout). The templates
> below are the real shapes from the `auth` slice — copy that slice to start a new feature.

---

## ⚠️ FIRST TIME HERE?

**🚨 READ [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) FIRST!**

That file contains the absolute non-negotiable rules.
This file is for quick templates and checklists.

---

## ⚡ THE 5 COMMANDMENTS

```
1. 📏 File >300 lines?        → STOP. Split it.
2. 📐 Function >100 lines?    → STOP. Extract functions.
3. 🧪 No tests?               → STOP. Write tests first.
4. ❌ Error ignored (_, _)?   → STOP. Handle it.
5. 📝 Exported without docs?  → STOP. Document it.
```

**VIOLATE = CODE REJECTED**

---

## 🎯 Before Writing ANY Code

```bash
# Ask yourself:
□ Which layer am I in? (controller / service / repository / model)
□ Am I following dependency direction? (controller → service → repository → model)
□ Does my service depend on a repository INTERFACE, not a concrete type?
□ Will this file exceed 300 lines? → Plan to split (another file, SAME package)
□ Will this function exceed 100 lines? → Plan to extract
□ Do I need tests? → Yes, ALWAYS for services (tests/unit/services + a fake in tests/mocks)
□ Is this documented? → Required for exported items
```

---

## 📐 Size Limits (HARD LIMITS)

```
File:     MAX 300 lines  (warning at 250)
Function: MAX 100 lines  (warning at 80)
```

**Approaching limit?**
- Stop and refactor NOW
- Don't wait until you exceed
- Split proactively

---

## 🏗️ Architecture Cheat Sheet

```
Request flow:
routers/index.go → Register<Name>Routes → Controller → Service → Repository → Database

Layers (one file per feature, per directory):
┌─────────────────────────────────────────────────────────────────┐
│ Controller  │ internal/app/controllers/<name>_controller.go     │  HTTP only
│ Service     │ internal/app/services/<name>_service.go           │  business logic
│ Repository  │ internal/domain/repositories/<name>_repo.go       │  data access only
│ Model       │ internal/domain/models/<name>_model.go            │  GORM struct
│ DTO         │ internal/app/dto/<name>_dto.go                    │  request/response
│ Routes      │ internal/app/routers/<name>_routes.go             │  route registration
└─────────────────────────────────────────────────────────────────┘

Dependency direction (never upward):
Controller  →  Service  →  Repository  →  Model
```

**Never:**
- ❌ Controller with business logic
- ❌ Controller calling a repository directly
- ❌ Service touching `database.DB`
- ❌ Repository importing a service

**Reference slice:** `auth` — the only complete vertical slice in this codebase
(model → repository → dto → service → controller → routes → mock → test). Copy its shape to
start a new feature. Use the repo's actual import path (`github.com/0xdiaz/gin-boilerplate/...`).

---

## 🔥 Layer Templates

### `<name>_model.go` — the GORM struct

```go
// internal/domain/models/user_model.go
package models

import "time"

// User represents an application user.
type User struct {
    ID    uint   `json:"id" gorm:"primaryKey"`
    Name  string `json:"name" gorm:"type:varchar(255);not null"`
    Email string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`

    // Password is never serialised: the json tag is "-".
    Password string `json:"-" gorm:"not null"`

    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

// TableName specifies the database table name for User model.
func (u *User) TableName() string {
    return "users"
}
```

Every model needs a matching **versioned migration** —
`internal/adapters/database/migrations/sql/NNNNNN_create_<table>.up.sql` plus its `.down.sql`.
There is no AutoMigrate.

### `<name>_dto.go` — request/response types

```go
// internal/app/dto/auth_dto.go
package dto

// LoginRequest is the body of POST /api/v1/auth/login.
type LoginRequest struct {
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// UserResponse is the API representation of a user. It deliberately omits the
// password hash and any reset-token fields.
type UserResponse struct {
    ID    uint   `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

Request DTOs carry `binding:"..."` tags and are bound at the HTTP boundary with
`c.ShouldBindJSON(&req)`.

### `<name>_repo.go` — data access, the ONLY layer touching the database

```go
// internal/domain/repositories/user_repo.go
package repositories

// UserRepository defines data access for user entity (used by auth and others).
type UserRepository interface {
    GetUserByEmail(email string) (*models.User, error)
    GetUserByID(id uint) (*models.User, error)
    CreateUser(user *models.User) error
    UpdateUser(user *models.User) error
}

// userRepo is the default implementation of UserRepository.
type userRepo struct{}

// NewUserRepository returns a new UserRepository implementation.
func NewUserRepository() UserRepository {
    return &userRepo{}
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
    var user models.User
    err := database.DB.Where("email = ?", email).First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // a missing row is not an error at this layer
        }
        logger.Errorf("failed to get user by email: %v", err)
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }
    return &user, nil
}
```

### `<name>_service.go` — business logic + sentinel errors

A feature small enough for one file lives at `internal/app/services/<name>_service.go` in
`package services`. `auth` outgrew that, so it has its own package — the shape is identical.

```go
// internal/app/services/auth/auth_service.go
package auth

// ErrUserNotFound is returned when the requested user does not exist.
var ErrUserNotFound = errors.New("user not found")

// AuthService handles authentication business logic.
type AuthService struct {
    userRepo         repositories.UserRepository
    refreshTokenRepo repositories.RefreshTokenRepository
    mailer           EmailSender
}

// NewAuthService creates a new AuthService instance.
func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository, mailer EmailSender) *AuthService {
    return &AuthService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo, mailer: mailer}
}

// GetProfile returns the profile of the given user.
//
// Returns ErrUserNotFound when no user has that id.
func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error) {
    ctx, start := logger.LogStart(ctx, "AuthService.GetProfile")

    user, err := s.userRepo.GetUserByID(userID)
    if err != nil {
        logger.Errorf("failed to get user for profile: %v", err)
        logger.LogFinish(ctx, "AuthService.GetProfile", err, start)
        return nil, fmt.Errorf("failed to get profile: %w", err)
    }
    if user == nil {
        logger.LogFinish(ctx, "AuthService.GetProfile", ErrUserNotFound, start)
        return nil, ErrUserNotFound
    }

    logger.LogFinish(ctx, "AuthService.GetProfile", nil, start)
    return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}
```

Services take repository **interfaces** — that is what makes them testable without a database.

### `<name>_controller.go` — thin HTTP layer

```go
// internal/app/controllers/auth_controller.go
package controllers

// AuthController handles authentication-related HTTP requests.
type AuthController struct {
    service auth.AuthServicer
}

// NewAuthController creates a new AuthController instance.
func NewAuthController(service auth.AuthServicer) *AuthController {
    return &AuthController{service: service}
}

// Profile returns the authenticated user's profile.
//
// GET /api/v1/profile
func (ctrl *AuthController) Profile(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Profile")

    userID := c.GetUint("user_id") // set by AuthMiddleware
    profile, err := ctrl.service.GetProfile(ctx, userID)
    if err != nil {
        if apiErr := authErrToAPIError(err); apiErr != nil {
            logger.LogFinish(ctx, "AuthController.Profile", err, start)
            utils.RespondWithAPIError(c, apiErr)
            return
        }
        logger.Errorf("get profile failed: %v", err)
        logger.LogFinish(ctx, "AuthController.Profile", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve profile")
        return
    }

    logger.LogFinish(ctx, "AuthController.Profile", nil, start)
    utils.Ok(c, profile, "Profile retrieved successfully")
}
```

### `<name>_routes.go` — route registration for one feature

```go
// internal/app/routers/auth_routes.go
package routers

// RegisterAuthRoutes registers authentication routes under the given group.
func RegisterAuthRoutes(group *gin.RouterGroup, authService auth.AuthServicer) {
    authController := controllers.NewAuthController(authService)

    authRoutes := group.Group("/auth")
    {
        authRoutes.POST("/register", authController.Register)
        authRoutes.POST("/login", authController.Login)
        authRoutes.POST("/refresh", authController.RefreshToken)
    }
}
```

### Wire it — a few lines in `routers/index.go`

```go
// internal/app/routers/index.go — the only file that knows concrete types
func RegisterRoutes(route *gin.Engine) {
    RegisterHealthRoutes(route) // /health and /metrics at the ROOT

    apiV1 := route.Group("/api/v1")
    apiV1.Use(middlewares.RateLimitMiddleware())

    userRepo := repositories.NewUserRepository()
    refreshTokenRepo := repositories.NewRefreshTokenRepository()
    authService := auth.NewAuthService(userRepo, refreshTokenRepo, nil)

    RegisterAuthRoutes(apiV1, authService)

    // Protected routes sit behind the auth middleware.
    authController := controllers.NewAuthController(authService)
    protectedRoutes := apiV1.Group("")
    protectedRoutes.Use(middlewares.AuthMiddleware(authService))
    {
        protectedRoutes.GET("/profile", authController.Profile)
    }
}
```

**Rules:**
- ✅ Each layer's `New*` takes its dependency; `index.go` assembles the chain.
- ✅ Adding a feature touches `index.go` once and adds one `<name>_routes.go`.
- ❌ Never register handlers inline in `index.go`.

---

## ✅ Error Handling Pattern

```go
// ✅ ALWAYS do this:
result, err := someFunction()
if err != nil {
    logger.Errorf("context: %v", err)                    // Log
    return fmt.Errorf("operation failed: %w", err)       // Wrap with %w
}

// ❌ NEVER do this:
result, _ := someFunction()                              // Ignored!
result, err := someFunction()
if err != nil {
    panic(err)                                           // Panic!
}
result, err := someFunction()
return err                                               // Not wrapped!
```

**Sentinel errors** are declared in the service and translated in the controller:

```go
// service
var ErrUserNotFound = errors.New("user not found")

// controller
if errors.Is(err, services.ErrUserNotFound) {
    utils.NotFound(c, err, "User not found")
    return
}
```

The `auth` service package keeps its own sentinels (`auth.ErrUserNotFound`,
`auth.ErrInvalidCredentials`, …) and maps them centrally in `authErrToAPIError`.

---

## 📝 Documentation Pattern

```go
// ✅ CORRECT:
// GetProfile returns the profile of the given user.
//
// Returns ErrUserNotFound when no user has that id.
func (s *AuthService) Get(ctx context.Context, id uint) (*dto.UserResponse, error) {
    // implementation
}

// ❌ WRONG:
// Get profile
func (s *AuthService) Get(ctx context.Context, id uint) (*dto.UserResponse, error) {

// ❌ WRONG:
func (s *AuthService) Get(ctx context.Context, id uint) (*dto.UserResponse, error) {  // No comment
```

---

## 🧪 Testing Checklist

```
⚠️  Unit tests live in the tests/ tree, NOT next to the code
□ Create tests/unit/services/<name>_service_test.go
□ Use package services_test (black box) — drive the exported surface
□ Create a shared fake in tests/mocks/<name>_repo_mock.go
□ Assert the fake satisfies the interface:
     var _ repositories.UserRepository = (*MockUserRepository)(nil)
□ Test happy path
□ Test 2+ error cases (drive the fake's error fields)
□ Use table-driven subtests if >3 scenarios
□ Assert with testify (require for fatal, assert for the rest)
□ Run: make test   (or: go test ./tests/unit/...)
```

**Canonical fake + test** (from `tests/mocks/user_repo_mock.go` and
`tests/unit/services/auth_service_test.go`):

```go
// tests/mocks/user_repo_mock.go
package mocks

// MockUserRepository is an in-memory UserRepository for unit tests.
type MockUserRepository struct {
    mu     sync.RWMutex
    nextID uint
    byID   map[uint]*models.User

    // GetErr, when set, is returned instead of data.
    GetErr error
}

var _ repositories.UserRepository = (*MockUserRepository)(nil)
```

```go
// tests/unit/services/auth_service_test.go
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
```

Tests that need a real database go in `tests/integration/`.

---

## 🚨 Forbidden Patterns

```go
❌ panic() in business logic
❌ _, _ = someFunc()                    // Ignored error
❌ "SELECT * FROM " + table             // SQL injection
❌ if x { if y { if z { } } }           // Too nested (>3 levels)
❌ password := "hardcoded"              // Hardcoded secrets
❌ log.Printf()                         // Use logger.Infof()
❌ file size >300 lines
❌ function >100 lines
❌ c.JSON(...) for API responses        // Use pkg/utils (utils.Ok, utils.BadRequest, …)
❌ database.DB in a service/controller  // Only repositories touch the database
❌ returning gorm.ErrRecordNotFound     // Return (nil, nil); let the service decide
❌ editing an already-applied migration // Add a new .up.sql / .down.sql pair
❌ float types for money                // Use int64 in the smallest currency unit
❌ No tests for services
❌ Exported function without docs
```

---

## 🎨 Naming Conventions

```go
// Files (feature-prefixed, snake_case)
✅ auth_service.go, user_repo.go, auth_service_tokens.go
❌ AuthService.go, auth-service.go, authService.go

// Packages (named after the layer directory)
✅ package controllers, package services, package repositories, package models
✅ package auth  (a service that needs multiple files gets its own subpackage)
❌ package Services, package auth_svc

// Types
✅ AuthController, AuthService, UserRepository, MockUserRepository
❌ AuthCtrl, AuthSvc, UserRepoImpl

// Constructors
✅ NewAuthService, NewUserRepository, NewAuthController
❌ CreateAuthService, MakeUserRepo

// Variables
✅ user, userID, refreshTokenRepo
❌ u, usrID, refresh_token_repo

// Constants
✅ const TokenTypeBearer = "Bearer"
❌ const TOKEN_TYPE_BEARER = "Bearer"
```

---

## 🔍 Pre-Commit Checklist

```bash
□ All functions <100 lines?
□ All files <300 lines?
□ All errors handled and wrapped with %w?
□ All exported items documented?
□ Tests written and passing?
□ No hardcoded secrets?
□ No panic() in business logic?
□ No SQL string concatenation?
□ No ignored errors (_, _)?
□ New migration has a matching .down.sql?
□ gofmt applied?

# Run these:
gofmt -w .
go vet ./...
go build ./...
make test
```

---

## 💡 Common Patterns

### Transactions

Multi-step writes that must succeed or fail together go through a single GORM transaction,
inside the repository layer:

```go
func (r *userRepo) DoTwoThings(...) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(&a).Error; err != nil {
            return fmt.Errorf("create a: %w", err)
        }
        if err := tx.Model(&b).Update("x", y).Error; err != nil {
            return fmt.Errorf("update b: %w", err)
        }
        return nil
    })
}
```

Returning a non-nil error from the callback rolls the whole thing back.

### Validation (gin binding tags, bound at the HTTP boundary)

```go
// internal/app/dto/auth_dto.go
type RegisterRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
}

// internal/app/controllers/auth_controller.go
var req dto.RegisterRequest
if err := c.ShouldBindJSON(&req); err != nil {
    utils.BadRequest(c, err, "Invalid request data")
    return
}
```

### DataTables (server-side pagination/search/sort)

See `internal/domain/repositories/example_repo.go` for the wiring of the
`Datatables-Gin` helper.

---

## 📚 Quick Links

- [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) — the non-negotiables, read first
- [`MODULE_GUIDE.md`](./MODULE_GUIDE.md) — layout source of truth
- [`CODING_STANDARDS.md`](./CODING_STANDARDS.md) — full standards
- [`DESIGN_PATTERNS.md`](./DESIGN_PATTERNS.md) — patterns and rationale
- [`AUTHENTICATION.md`](./AUTHENTICATION.md) — JWT, refresh rotation, password reset
- [`MIGRATIONS.md`](./MIGRATIONS.md) — schema change workflow

---

## 🎯 Remember

**When in doubt, copy the `auth` slice.** It is the most recent feature and follows every
rule on this page:

```
internal/adapters/database/migrations/sql/000001_create_users_table.up.sql
internal/domain/models/user_model.go
internal/domain/repositories/user_repo.go
internal/app/dto/auth_dto.go
internal/app/services/auth/auth_service.go
internal/app/controllers/auth_controller.go
internal/app/routers/auth_routes.go
tests/mocks/user_repo_mock.go
tests/unit/services/auth_service_test.go
```
