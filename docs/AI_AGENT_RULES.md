# AI Agent Rules

> **CRITICAL**: These rules are MANDATORY for all AI agents working on this codebase.
> Violating these rules will result in rejected code.

> 🧭 **Structure:** This is a **modular service** — code lives in `internal/modules/<name>/`
> (one vertical slice per module), cross-cutting code in `pkg/`, wiring in `internal/bootstrap/`.
> See **[MODULE_GUIDE.md](./MODULE_GUIDE.md)** (the source of truth for the layout). The mandatory
> rules below apply *inside* that layout.

---

## 🚨 CRITICAL RULES (Never Break These)

### 1. File Size - HARD LIMIT
```
✅ MAX 300 lines per file
❌ >300 lines → MUST split into multiple files
```

**If file exceeds 300 lines:**
1. STOP immediately
2. Split into logical components
3. Create separate files (e.g., `service.go` → `service.go` + `service_helpers.go`)

### 2. Function Size - HARD LIMIT
```
✅ MAX 100 lines per function
❌ >100 lines → MUST split into smaller functions
```

**If function exceeds 100 lines:**
1. Extract logical blocks into separate functions
2. Each function should do ONE thing
3. Use descriptive names for extracted functions

### 3. NO Testing = NO Merge
```
❌ FORBIDDEN to create/modify services without tests
✅ REQUIRED: Test file for every service file
✅ MINIMUM: 70% coverage for services
```

**When creating/modifying a service:**
1. Create `_test.go` file immediately
2. Write at least 3 test cases (happy path + 2 errors)
3. Run tests before committing

### 4. Documentation - MANDATORY
```
✅ MUST document ALL exported functions/structs
❌ NO exported code without documentation
```

**Required documentation:**
- Package comment (in main file or doc.go)
- Function comment (what it does, params, returns, errors)
- Struct comment (purpose, usage)

### 5. Error Handling - ZERO TOLERANCE
```
❌ FORBIDDEN: Ignoring errors with _, _
❌ FORBIDDEN: panic() in business logic
✅ REQUIRED: Wrap errors with context
```

**Always:**
```go
if err != nil {
    logger.Errorf("context: %v", err)
    return fmt.Errorf("operation failed: %w", err)
}
```

**Never:**
```go
result, _ := someFunc()  // ❌ FORBIDDEN
panic("error")           // ❌ FORBIDDEN in services
```

### 6. DTO and Constant Placement - STRICT RULES
```
✅ DTOs are MODULE-LOCAL → internal/modules/<name>/dto.go (package <module>)
✅ CONSTANTS live with the module that owns them (package <module>)
✅ Truly cross-cutting shared types/constants → pkg/types
❌ There is NO shared internal/app/dto and NO pkg/enums package — do not invent them
```

**When creating any DTO (request/response struct):**
```
✅ CORRECT: internal/modules/auth/dto.go     → RegisterRequest, LoginRequest, AuthResponse
            (dto.go is optional; tiny modules may keep the type in model.go)
❌ WRONG:   internal/app/dto/user_dto.go      → (the internal/app tree no longer exists)
❌ WRONG:   pkg/types/request.go              → (pkg/types is for shared response/error types only)
```

**When creating any constant or enum-like value:**
```
✅ CORRECT: internal/modules/auth/service.go  → ErrInvalidCredentials, ErrEmailAlreadyExists
            (a module's status/role/type constants live in that module's package)
❌ WRONG:   pkg/enums/role.go                 → (no such package)
```

**IMPORT DIRECTION (no cycles):**
```
✅ ALLOWED:
   internal/modules/<name> → pkg/utils, pkg/types, pkg/logger, pkg/config
   internal/bootstrap      → internal/modules/<name>
   one module → another module's PUBLIC interface (injected via the constructor)

❌ FORBIDDEN (import cycle / breaks layering):
   pkg/anything            → internal/...
   internal/modules/a      → internal/modules/b's unexported types
```

**Cross-module data:** never reach into another module's package for its DTOs/constants. A
module that needs another module receives its **public interface** (e.g. `auth.Servicer` via
`authMod.Auth()`) injected in `buildModules()`. See [DESIGN_PATTERNS.md](./DESIGN_PATTERNS.md)
§3.6 (Dependency Injection) and §13.5 (Reaching Into Another Module).

### 7. API Responses - MANDATORY UTILS
```
✅ MUST use pkg/utils response functions
❌ NEVER use c.JSON() directly in handlers (for API responses)
✅ ALWAYS import "pkg/utils" in handlers
```

**When sending success responses:**
```go
// ✅ CORRECT — handler methods on a struct, using pkg/utils
import "github.com/0xdiaz/gin-boilerplate/pkg/utils"

func (h *Handler) GetUser(c *gin.Context) {
    user, err := h.svc.GetUserByID(c.Request.Context(), id)
    if err != nil {
        utils.NotFound(c, err, "User not found")
        return
    }
    utils.Ok(c, user, "User retrieved successfully")
}

func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    user, err := h.svc.CreateUser(c.Request.Context(), &req)
    if err != nil {
        utils.InternalServerError(c, err, "Failed to create user")
        return
    }
    utils.Created(c, user, "User created successfully")
}

// ❌ WRONG - direct c.JSON()
func (h *Handler) GetUser(c *gin.Context) {
    user, _ := h.svc.GetUserByID(c.Request.Context(), id)
    c.JSON(200, gin.H{"data": user})  // Inconsistent format!
}
```

**Available functions:**
```
Success:
- utils.Ok(c, data, message)           // 200
- utils.Created(c, data, message)      // 201
- utils.NoContent(c)                   // 204

Errors:
- utils.BadRequest(c, err, message)    // 400
- utils.Unauthorized(c, err, message)  // 401
- utils.Forbidden(c, err, message)     // 403
- utils.NotFound(c, err, message)      // 404
- utils.Conflict(c, err, message)      // 409
- utils.InternalServerError(c, err, message) // 500
```

**Why this rule:**
- ✅ Consistent response format across all endpoints
- ✅ Automatic validation error formatting
- ✅ Standard JSON structure for clients
- ✅ Easier to maintain and modify response format globally

---

## 🎯 STEP-BY-STEP WORKFLOW

### Before Writing Any Code

1. **Check existing code**
   - Read related files
   - Understand current patterns
   - Follow existing style

2. **Plan the structure**
   - Estimate lines of code
   - If >300 lines → plan multiple files
   - If >100 lines per function → plan extraction

3. **Check dependencies**
   - Which module am I in? Which file in the slice? (handler/service/repository/model)
   - Am I following dependency direction?
   - Am I depending on a consumer-defined interface, not a concrete type?

### While Writing Code

1. **Follow the architecture**
   ```
   Handler (HTTP) → Service (Logic) → Repository (Data) → Model
   ```

2. **Keep count of lines**
   - Function approaching 80 lines? Plan to split
   - File approaching 250 lines? Plan new file

3. **Document as you write**
   - Write function comment BEFORE implementation
   - This helps clarify what function should do

4. **Handle errors immediately**
   - Never write `_, err :=` without handling err
   - Always log errors
   - Always wrap errors with context

### After Writing Code

1. **Self-review checklist**
   - [ ] All functions <100 lines?
   - [ ] All files <300 lines?
   - [ ] All errors handled?
   - [ ] All exported items documented?
   - [ ] No hardcoded secrets?
   - [ ] No panic() in business logic?

2. **Write tests**
   - [ ] Test file created?
   - [ ] Happy path tested?
   - [ ] Error cases tested?
   - [ ] Tests passing?

3. **Run checks**
   ```bash
   gofmt -w .
   go vet ./...
   go test ./...
   ```

---

## 📐 ARCHITECTURE RULES

### Handler Layer (the HTTP layer — not a "controller")
```go
✅ DO:
- Bind/parse the HTTP request (params, body, headers)
- Open/close a trace span (logger.LogStart / logger.LogFinish)
- Call the service through the consumer-defined `service` interface
- Respond via pkg/utils (utils.Ok, utils.Created, utils.BadRequest, …)

❌ DON'T:
- Contain business logic
- Access the database / call repositories directly
- Have functions >50 lines
- Write c.JSON(...) by hand for API responses
- Import another module's internals
```

**Template:**
```go
func (h *Handler) GetUser(c *gin.Context) {
    // 1. Open a trace span + parse input
    ctx, start := logger.LogStart(c.Request.Context(), "users.Handler.GetUser")
    id := c.Param("id")

    // 2. Call the service
    user, err := h.svc.GetUserByID(ctx, id)
    if err != nil {
        logger.LogFinish(ctx, "users.Handler.GetUser", err, start)
        utils.NotFound(c, err, "User not found")
        return
    }

    // 3. Respond via pkg/utils
    logger.LogFinish(ctx, "users.Handler.GetUser", nil, start)
    utils.Ok(c, user, "User retrieved successfully")
}
// Total: ~15 lines
```

### Service Layer
```go
✅ DO:
- Implement ALL business logic
- Define the `repository` interface it needs (consumer-defined, in service.go)
- Accept and propagate context.Context (LogStart / LogFinish)
- Orchestrate multiple repositories; return module sentinel errors
- Transform DTO ↔ model

❌ DON'T:
- Handle HTTP concerns (gin.Context) — except a lib that needs it (e.g. DataTables paging)
- Import the handler layer or another module's internals
- Exceed 400 lines per file (split into service_<topic>.go in the SAME package)
- Have functions >100 lines
```

**Template** (module-local types; `req`, `User`, `ErrEmailAlreadyExists` are in this module):
```go
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
        return nil, ErrEmailAlreadyExists // sentinel error — handler maps to 409
    }

    // 2. Apply business logic + persist
    hashed, err := s.hashPassword(req.Password)
    if err != nil {
        logger.LogFinish(ctx, "auth.Service.Register", err, start)
        return nil, fmt.Errorf("failed to process password: %w", err)
    }
    user := &User{Name: req.Name, Email: req.Email, Password: hashed}
    if err = s.userRepo.CreateUser(user); err != nil {
        logger.LogFinish(ctx, "auth.Service.Register", err, start)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    logger.Infof("user registered successfully: %s", user.Email)
    logger.LogFinish(ctx, "auth.Service.Register", nil, start)
    return &AuthResponse{ /* ... */ }, nil
}
// Total: ~30 lines
```

### Repository Layer
```go
✅ DO:
- CRUD operations only, on the INJECTED *gorm.DB (NewRepository(db); no globals)
- Database queries / transaction management
- Translate DB-specific errors (gorm.ErrRecordNotFound → nil, nil)
- Return the module's own models

❌ DON'T:
- Contain business logic / validate business rules
- Use the database.GetDB() global
- Call another module's repository
- Import the service or handler layer
```

**Template** (`User` is this module's model; `r.db` is the injected connection):
```go
func (r *Repository) FindByID(id uint) (*User, error) {
    var user User

    if err := r.db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // not found is not an error here; the service decides
        }
        return nil, fmt.Errorf("query failed: %w", err)
    }

    return &user, nil
}
// Total: ~12 lines
```

---

## 🔍 COMMON PATTERNS

> In the patterns below, `User`, `UpdateUserRequest`, and the `Err*` sentinels are all
> declared in the module's own package (no `models.`/`dto.` prefix).

### 1. List with Pagination
```go
func (s *Service) ListUsers(ctx context.Context, page, pageSize int) ([]*User, int64, error) {
    // Validate pagination
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    // Get data
    users, total, err := s.repo.List(page, pageSize)
    if err != nil {
        return nil, 0, fmt.Errorf("list users: %w", err)
    }

    return users, total, nil
}
```

### 2. Update Operations
```go
func (s *Service) UpdateUser(ctx context.Context, id uint, req *UpdateUserRequest) (*User, error) {
    // 1. Get existing
    user, err := s.repo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("get user: %w", err)
    }

    // 2. Check business rules
    if req.Email != user.Email {
        existing, _ := s.repo.GetUserByEmail(req.Email)
        if existing != nil {
            return nil, ErrEmailAlreadyExists // module sentinel error
        }
    }

    // 3. Update fields
    user.Name = req.Name
    user.Email = req.Email

    // 4. Save
    if err := s.repo.Update(user); err != nil {
        return nil, fmt.Errorf("update user: %w", err)
    }

    return user, nil
}
```

### 3. Delete Operations
```go
func (s *Service) DeleteUser(ctx context.Context, id uint) error {
    // 1. Check exists
    user, err := s.repo.FindByID(id)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }

    // 2. Check business rules (e.g., can't delete if it still has active transactions)
    hasTransactions, err := s.repo.HasActiveTransactions(id)
    if err != nil {
        return fmt.Errorf("check transactions: %w", err)
    }
    if hasTransactions {
        return ErrCannotDelete // module sentinel error
    }

    // 3. Delete
    if err := s.repo.Delete(user.ID); err != nil {
        return fmt.Errorf("delete user: %w", err)
    }

    logger.Infof("user deleted: ID=%d", id)
    return nil
}
```

### 4. Transaction Pattern (the tx lives in the REPOSITORY, on the injected *gorm.DB)
```go
// internal/modules/<name>/repository.go — the service orchestrates WHEN to call this;
// the repository owns the transaction mechanics (the service never touches *gorm.DB).
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

        // returning nil commits; any returned error rolls back automatically
        return nil
    })
}
```

---

## 🛡️ SECURITY CHECKLIST

### Before Committing

- [ ] No hardcoded passwords/secrets
- [ ] No SQL string concatenation
- [ ] All input validated (gin binding tags + c.ShouldBindJSON)
- [ ] Passwords hashed with bcrypt
- [ ] Protected routes behind the auth module's Middleware() guard
- [ ] Authorization checked where needed (using the authenticated user_id)
- [ ] Sensitive data sanitized in logs
- [ ] Rate limiting applied
- [ ] CORS configured properly

### Validation Pattern (gin binding tags, bound at the HTTP boundary)
```go
// internal/modules/<name>/dto.go — rules are gin binding tags (validated by ShouldBindJSON)
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// internal/modules/<name>/handler.go — bind + validate in one step;
// utils.BadRequest turns validator.ValidationErrors into a field→message map.
func (h *Handler) CreateUser(c *gin.Context) {
    var req CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    // req is validated; pass to the service
}
```

---

## 📝 TESTING PATTERNS

### Test File Structure (CO-LOCATED, white-box, in-package fake repo — no DB)

⚠️ Unit tests live next to the module as `internal/modules/<name>/*_test.go`, in the **same
package** (`package <module>`) so the fake can satisfy the module's **unexported** `repository`
interface. This is the canonical pattern; see `internal/modules/example/service_test.go`.

```go
// internal/modules/users/service_test.go
package users

import (
    "context"
    "errors"
    "testing"

    "github.com/stretchr/testify/assert"
)

// fakeRepo satisfies the unexported repository interface — no DB needed.
type fakeRepo struct {
    byEmail *User
    err     error
}

func (f *fakeRepo) GetUserByEmail(email string) (*User, error) { return f.byEmail, f.err }
func (f *fakeRepo) CreateUser(u *User) error                   { return f.err }

func TestCreateUser_Success(t *testing.T) {
    svc := NewService(&fakeRepo{}) // GetUserByEmail returns (nil, nil) → no duplicate

    user, err := svc.CreateUser(context.Background(), &CreateUserRequest{
        Name: "John Doe", Email: "john@example.com",
    })

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "John Doe", user.Name)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
    existing := &User{ID: 1, Email: "existing@example.com"}
    svc := NewService(&fakeRepo{byEmail: existing})

    user, err := svc.CreateUser(context.Background(), &CreateUserRequest{
        Name: "John Doe", Email: "existing@example.com",
    })

    assert.Nil(t, user)
    assert.True(t, errors.Is(err, ErrEmailAlreadyExists)) // module sentinel error
}
```

> Prefer in-package fakes over the shared `tests/mocks/` helpers for new code. Handler tests use
> `httptest` + a fake implementing the handler's `service` interface.

### Table-Driven Tests
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"missing domain", "user@", true},
        {"empty", "", true},
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

## 🚀 QUICK REFERENCE

### When Creating a New Module

The fastest path: `cp -r internal/modules/example internal/modules/<name>`, rename the package,
then adjust. See MODULE_GUIDE.md "How to add a new module".

1. Create `internal/modules/<name>/` (one package = one folder).
2. `model.go` — the GORM model(s); list them in `Module.Models()`.
3. `repository.go` — `Repository{db}` + `NewRepository(db)` (injected *gorm.DB, no globals).
4. `service.go` — `Service` + the `repository` interface it consumes; constructor `NewService`.
5. `handler.go` — `Handler` + the `service` interface it consumes; constructor `NewHandler`.
6. `module.go` — `New(db)`, `Name()`, `Models()`, `RegisterRoutes(api)`, public `API()`.
7. Co-located `service_test.go` — white-box `package <name>`, in-package fake repo (no DB).
8. Register it: add `<name>.New(db)` to `buildModules()` in `internal/bootstrap/modules.go`.
9. Document all exported items.
10. Run: `make test`.

### When Creating the Service File

1. Define the `repository` interface this service needs (in `service.go`, consumer-defined).
2. Implement `Service` struct + `NewService(repo)` constructor.
3. Implement methods (keep <100 lines each); accept/propagate `context.Context`.
4. Use `logger.LogStart`/`LogFinish` around each method; return module sentinel errors.
5. Add error handling to every function (wrap unexpected errors with `%w`).
6. Create co-located `service_test.go` with a fake repo; write at least 3 cases.

### When Creating the Repository File

1. `Repository{db *gorm.DB}` + `NewRepository(db)` — hold the INJECTED connection, no globals.
2. Use GORM / parameterized queries (no raw SQL string concatenation).
3. Return the module's own models.
4. Translate DB errors (`gorm.ErrRecordNotFound` → `nil, nil`).
5. Use transactions for multi-step operations (the tx lives here, not in the service).
6. The interface it satisfies is defined by the consuming service — not by the repository.

### When Creating the Handler File

1. Keep handler methods thin (<50 lines): bind request → call service → respond.
2. Define the `service` interface this handler needs (in `handler.go`, consumer-defined).
3. Respond via `pkg/utils` (`utils.Ok`, `utils.Created`, `utils.RespondWithAPIError`, …).
4. Protect routes with `m.Middleware()` (the auth module's JWT guard).
5. Don't put business logic here.

### When Refactoring Large Files

1. Identify logical boundaries
2. Extract to separate files:
   - `service.go` - main service logic
   - `service_helpers.go` - helper functions
   - `service_validators.go` - validation logic
   - `service_transformers.go` - data transformation
3. Update imports
4. Run tests to ensure nothing broke

---

## ❌ RED FLAGS (Stop and Refactor)

If you see ANY of these, STOP and refactor:

1. **File >300 lines** → Split now
2. **Function >100 lines** → Extract functions
3. **Error ignored (`_, _`)** → Handle it
4. **No tests** → Write tests
5. **No documentation** → Add comments
6. **panic() in service** → Return error instead
7. **Business logic in handler** → Move to service
8. **Database access in handler** → Use the repository
9. **Hardcoded secrets** → Move to .env
10. **SQL string concatenation** → Use GORM/parameterized queries
11. **Subfolder inside a module** → Split into more files in the SAME package
12. **Importing another module's internals** → Depend on its public interface, injected

---

## 💡 TIPS FOR AI AGENTS

### Estimating Line Count Before Writing

```
Simple CRUD:
- Handler: ~15 lines
- Service: ~30 lines
- Repository: ~12 lines

Complex operation:
- Handler: ~30 lines
- Service: ~80 lines (if exceeds, split into service_<topic>.go in the SAME package!)
- Repository: ~25 lines

If you estimate >100 lines for one function:
→ Plan to split into 3-5 smaller functions
```

### Naming Functions When Splitting

```go
// BEFORE: One giant function
func CreateTransactionFromApiLog() { } // 480 lines

// AFTER: Multiple focused functions
func CreateTransactionFromApiLog() { }      // 20 lines - orchestration
func parseRequestData() { }                 // 30 lines
func parseResponseData() { }                // 30 lines
func buildTransactionFromParsedData() { }   // 40 lines
func calculateTransactionFees() { }         // 25 lines
func determineTransactionStatus() { }       // 20 lines
func saveTransactionWithAudit() { }         // 30 lines
```

### Incremental Development

1. Write function signature + documentation FIRST
2. Write test cases SECOND
3. Implement logic THIRD
4. Refactor if needed FOURTH

This prevents writing too much code before realizing it's wrong.

---

## 📚 REFERENCE FILES

- Layout (source of truth): [`MODULE_GUIDE.md`](./MODULE_GUIDE.md)
- Full standards: [`CODING_STANDARDS.md`](./CODING_STANDARDS.md)
- Design patterns: [`DESIGN_PATTERNS.md`](./DESIGN_PATTERNS.md)
- Critical rules: [`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md), templates: [`AI_QUICK_REFERENCE.md`](./AI_QUICK_REFERENCE.md)
- Canonical reference module — copy it to start a new one:
  - `internal/modules/example/` (model, dto, repository, service, handler, module, service_test)
  - `internal/modules/auth/` (real module: JWT guard, exported Servicer, split service files)

---

**REMEMBER:**
- ✅ Small files = Easy to maintain
- ✅ Small functions = Easy to test
- ✅ Good tests = Confident refactoring
- ✅ Good docs = Future you will thank you

**Last updated:** 2026-06-10
**Enforcement:** MANDATORY for all commits
