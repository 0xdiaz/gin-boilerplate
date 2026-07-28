# AI Agent Quick Reference

> **Print this mentally before every code change!**

> 🧭 **Structure:** Code lives in `internal/modules/<name>/` — one package per feature, holding
> `model.go`, `dto.go` (optional), `repository.go`, `service.go`, `handler.go`, `module.go`. See
> **[MODULE_GUIDE.md](./MODULE_GUIDE.md)** (the source of truth for the layout). The templates below
> are the real shapes from the canonical `example` module — copy that module to start a new one.

---

## ⚠️ FIRST TIME HERE?

**🚨 READ [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md) FIRST!**

That file contains the absolute non-negotiable rules (100 lines).
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
□ Which module am I in? Which file in the slice? (handler/service/repository/model)
□ Am I following dependency direction? (handler → service → repository → model)
□ Am I depending on a consumer-defined interface, not a concrete type?
□ Will this file exceed 300 lines? → Plan to split (another file, SAME module package)
□ Will this function exceed 100 lines? → Plan to extract
□ Do I need tests? → Yes, ALWAYS for services (co-located, fake repo)
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
Request Flow (one module = one vertical slice):
Module.RegisterRoutes → Handler → Service → Repository → Database

Layers (co-located in internal/modules/<name>/):
┌──────────────┐
│   Handler    │  ← HTTP only, <50 lines/function (defines the `service` interface)
├──────────────┤
│   Service    │  ← Business logic, <100 lines/function (defines the `repository` interface)
├──────────────┤
│  Repository  │  ← CRUD only on the injected *gorm.DB, return models
├──────────────┤
│    Model     │  ← Data structures only (the tables this module owns)
└──────────────┘

Dependencies:
Handler  →  Service  →  Repository  →  Model
    ↓          ↓            ↓
   DTO      Utils      injected *gorm.DB

Cross-module: only via a module's PUBLIC interface, injected in buildModules().
```

**Forbidden:**
- ❌ Handler with business logic
- ❌ Service accessing the database directly (go through the repository)
- ❌ Repository reaching for the global DB or another module's tables
- ❌ Importing another module's unexported types
- ❌ `pkg/` importing `internal/`
- ❌ Circular dependencies

---

## 🔥 Module Templates

A module is one folder = one package = one vertical slice. These are the real shapes from
`internal/modules/example/` — copy that module to start a new one. Use the repo's actual import
path `github.com/0xdiaz/gin-boilerplate`.

### model.go — the tables this module owns
```go
package example

import "time"

// Example is the module's GORM model. Each module owns its own tables.
type Example struct {
    ID        int        `json:"id" gorm:"primaryKey"`
    Data      string     `json:"data" binding:"required"`
    CreatedAt *time.Time `json:"created_at"`
    UpdatedAt *time.Time `json:"updated_at"`
}

func (e *Example) TableName() string { return "examples" }
```

### dto.go (optional) — request/response types with binding tags
```go
package example

// CreateRequest is the payload for creating an Example (validated by gin binding tags).
type CreateRequest struct {
    Data string `json:"data" binding:"required,min=1,max=255"`
}
```

### repository.go — data access on the INJECTED *gorm.DB (no globals)
```go
package example

import "gorm.io/gorm"

// Repository holds an injected *gorm.DB (no global state) so it is testable and reusable.
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

### service.go — business logic + the `repository` interface IT needs
```go
package example

import (
    "context"

    "github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

// repository is the data-access contract this service needs (consumer-defined → testable).
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

### handler.go — thin HTTP layer + the `service` interface IT needs
```go
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

type Handler struct {
    svc service
}

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

### module.go — wires the slice + the Module contract (Name/Models/RegisterRoutes/API)
```go
package example

import (
    "context"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

// API is the minimal public surface other modules depend on (never import internals).
type API interface {
    List(ctx context.Context) ([]*Example, error)
}

type Module struct {
    svc     *Service
    handler *Handler
}

// New assembles repository → service → handler from an injected DB connection.
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db))
    return &Module{svc: svc, handler: NewHandler(svc)}
}

func (m *Module) Name() string  { return "example" }
func (m *Module) Models() []any { return []any{&Example{}} }

// RegisterRoutes mounts the module's routes under the given API group (e.g. /api/v1).
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    api.GET("/examples", m.handler.List)
}

// API exposes this module's public contract to other modules.
func (m *Module) API() API { return m.svc }
```

### Register the module — ONE line in buildModules()
```go
// internal/bootstrap/modules.go
func buildModules(db *gorm.DB) []Module {
    return []Module{
        auth.New(db),
        example.New(db),
        // newmodule.New(db),  ← add here
    }
}
```

**Rules:**
- ✅ One module = one package = one folder; split big files into more files in the SAME folder.
- ✅ Each layer's `New*` takes a consumer-defined interface; `module.go` assembles the chain.
- ✅ Adding a module is one line in `buildModules()`.
- ❌ DON'T create a central router, layer folders, or subfolders inside a module.

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

---

## 📝 Documentation Pattern

```go
// ✅ CORRECT:
// CreateUser creates a new user account with validation.
//
// Returns ErrDuplicateEntry if email exists.
// Returns ErrValidation if input is invalid.
func CreateUser(dto CreateUserRequest) (*User, error) {
    // implementation
}

// ❌ WRONG:
// Create user
func CreateUser(dto CreateUserRequest) (*User, error) {

// ❌ WRONG:
func CreateUser(dto CreateUserRequest) (*User, error) {  // No comment
```

---

## 🧪 Testing Checklist

```go
⚠️  CRITICAL: unit tests are CO-LOCATED with the module
□ Create internal/modules/<name>/{filename}_test.go
□ Use package <module> (white-box) — so it can implement the unexported repository interface
□ Implement an in-package fakeRepo (no DB needed)
□ Test happy path
□ Test 2+ error cases (drive the fake's return values)
□ Use table-driven tests if >3 scenarios
□ Assert with testify (assert.NoError, assert.Equal, …)
□ Run: go test ./internal/... ./pkg/...  (or: make test)
□ Coverage >70% for services
```

**Example test location:**
```
internal/modules/example/service.go
→ internal/modules/example/service_test.go  (package example — white-box, fake repo)

internal/modules/auth/service.go
→ internal/modules/auth/service_test.go     (package auth)
```

**Canonical fake-repo test** (from `internal/modules/example/service_test.go`):
```go
package example

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
)

// fakeRepo satisfies the unexported repository interface — no DB needed.
type fakeRepo struct {
    items []*Example
    err   error
}

func (f *fakeRepo) List() ([]*Example, error) { return f.items, f.err }

func TestService_List(t *testing.T) {
    svc := NewService(&fakeRepo{items: []*Example{{ID: 1, Data: "x"}}})

    got, err := svc.List(context.Background())

    assert.NoError(t, err)
    assert.Len(t, got, 1)
    assert.Equal(t, "x", got[0].Data)
}
```

> Integration tests (real server/DB) and shared legacy mocks live under `tests/`.

---

## 🚨 Forbidden Patterns

```go
❌ panic() in business logic
❌ _, _ = someFunc()              // Ignored error
❌ "SELECT * FROM " + table       // SQL injection
❌ if x { if y { if z { } } }     // Too nested (>3 levels)
❌ password := "hardcoded"        // Hardcoded secrets
❌ log.Printf()                   // Use logger.Infof()
❌ file size >300 lines
❌ function >100 lines
❌ subfolders inside a module      // Split into more files in the SAME module package
❌ c.JSON(...) for API responses   // Use pkg/utils (utils.Ok, utils.BadRequest, …)
❌ database.GetDB() in a repo       // Use the injected *gorm.DB (NewRepository(db))
❌ importing another module's internals  // Depend on its public interface, injected
❌ No tests for services
❌ Exported function without docs
```

---

## 🎨 Naming Conventions

```go
// Files
✅ service.go, repository.go, service_tokens.go
❌ Service.go, service-tokens.go, serviceTokens.go

// Packages (named after the module/sub-domain)
✅ package example, package auth, package health
❌ package services, package user_auth, package Auth

// Variables
✅ user, userID, httpClient
❌ u, usrID, http_client

// Functions
✅ GetUserByID, CreateTransaction
❌ get_user, GetUser (too generic)

// Constants
✅ const MaxRetryAttempts = 3
✅ const StatusPending Status = "pending"
❌ const MAX_RETRY = 3
```

---

## 🔍 Pre-Commit Checklist

```bash
□ All functions <100 lines?
□ All files <300 lines?
□ All errors handled?
□ All exported items documented?
□ Tests written and passing?
□ No hardcoded secrets?
□ No panic() in business logic?
□ No SQL string concatenation?
□ No ignored errors (_, _)?
□ gofmt applied?

# Run these:
gofmt -w .
go vet ./...
go test ./...
```

---

## 🚀 When Refactoring Large Files

**If file >300 lines:**

1. **Identify boundaries**
   - Group related functions
   - Find logical separations

2. **Create new files**
   ```
   service.go           → service.go (main)
                        → service_helpers.go
                        → service_validators.go
                        → service_transformers.go
   ```

3. **Move code**
   - Keep related functions together
   - Maintain package cohesion

4. **Update imports**

5. **Run tests**
   ```bash
   go test ./...
   ```

**If function >100 lines:**

1. **Extract logical blocks**
   ```go
   // Before: 200 lines
   func Process() { ... }

   // After: Multiple focused functions
   func Process() {           // 20 lines - orchestration
       data := parse()
       validated := validate(data)
       transformed := transform(validated)
       save(transformed)
   }

   func parse() { }           // 30 lines
   func validate() { }        // 25 lines
   func transform() { }       // 40 lines
   func save() { }            // 20 lines
   ```

---

## 💡 Common Patterns

### Pagination
```go
func List(page, pageSize int) ([]*Model, int64, error) {
    if page < 1 { page = 1 }
    if pageSize < 1 || pageSize > 100 { pageSize = 20 }

    var items []*Model
    var total int64

    db := r.db.Model(&Model{})
    db.Count(&total)

    err := db.
        Offset((page - 1) * pageSize).
        Limit(pageSize).
        Find(&items).
        Error

    return items, total, err
}
```

### Transactions
```go
func Process(data Data) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := step1(tx, data); err != nil {
            return err  // Auto rollback
        }
        if err := step2(tx, data); err != nil {
            return err  // Auto rollback
        }
        return nil  // Auto commit
    })
}
```

### Validation (gin binding tags, bound at the HTTP boundary)
```go
// dto.go — declare rules as gin binding tags
type Request struct {
    Name  string `json:"name" binding:"required,min=3,max=255"`
    Email string `json:"email" binding:"required,email"`
    Age   int    `json:"age" binding:"gte=0,lte=150"`
}

// handler.go — bind + validate in one step; utils.BadRequest formats field errors
func (h *Handler) Create(c *gin.Context) {
    var req Request
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    // req is validated; pass to the service
}
```

---

## 📚 Quick Links

- **Layout (source of truth):** [`MODULE_GUIDE.md`](./MODULE_GUIDE.md)
- **Full Standards:** [`CODING_STANDARDS.md`](./CODING_STANDARDS.md)
- **Design Patterns:** [`DESIGN_PATTERNS.md`](./DESIGN_PATTERNS.md)
- **AI Rules:** [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md), [`AI_AGENT_RULES.md`](./AI_AGENT_RULES.md)
- **Quality Checks:** `make test` (runs `./tests/unit/... ./internal/... ./pkg/...`)

---

## 🎯 Remember

```
Small files    = Easy to understand
Small functions = Easy to test
Good tests     = Confident refactoring
Good docs      = Happy developers

✅ Quality > Quantity
✅ Simple > Complex
✅ Clear > Clever
```

---

**Print this before every commit!**
**Follow the rules strictly!**
**Your future self will thank you!**
