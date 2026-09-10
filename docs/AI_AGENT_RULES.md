# AI Agent Rules

> **CRITICAL**: These rules are MANDATORY for all AI agents working on this codebase.
> Violating these rules will result in rejected code.

> 🧭 **Structure:** Code is organized **by technical layer** — HTTP in
> `internal/app/controllers`, business logic in `internal/app/services`, data access in
> `internal/domain/repositories`, GORM structs in `internal/domain/models`, wired in
> `internal/app/routers/index.go`; cross-cutting code in `pkg/`.
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
✅ DTOs live in one place → internal/app/dto/<name>_dto.go (package dto)
✅ Domain constants live with the model that owns them → internal/domain/models/<name>_model.go
✅ Sentinel errors live with the service that returns them → internal/app/services/<name>_service.go
✅ Truly cross-cutting shared types → pkg/types
❌ There is NO pkg/enums package — do not invent one
```

**When creating any DTO (request/response struct):**
```
✅ CORRECT: internal/app/dto/auth_dto.go     → RegisterRequest, LoginRequest, UserResponse
✅ CORRECT: internal/app/dto/auth_dto.go      → RegisterRequest, LoginRequest, AuthResponse
❌ WRONG:   internal/domain/models/user_model.go → (models are GORM structs, not API shapes)
❌ WRONG:   pkg/types/request.go               → (pkg/types is for shared response/error types only)
```

**When creating any constant or enum-like value:**
```
✅ CORRECT: internal/domain/models/refresh_token_model.go → status/type constants for that model
            (they live next to the model they describe, and must stay in sync with any
             CHECK constraint in the migration)
✅ CORRECT: internal/app/services/auth/auth_service.go → ErrUserNotFound
✅ CORRECT: internal/app/services/auth/auth_service.go → ErrInvalidCredentials, ErrUserNotFound
❌ WRONG:   pkg/enums/role.go                       → (no such package)
```

**IMPORT DIRECTION (no cycles):**
```
✅ ALLOWED:
   internal/... → pkg/utils, pkg/types, pkg/logger, pkg/config
   internal/app/controllers   → internal/app/services, internal/app/dto
   internal/app/services      → internal/domain/repositories, internal/domain/models, internal/app/dto
   internal/domain/repositories → internal/domain/models, internal/adapters/database
   internal/app/routers       → controllers, services, repositories, middlewares

❌ FORBIDDEN (import cycle / breaks layering):
   pkg/anything                 → internal/...
   internal/domain/repositories → internal/app/services
   internal/app/controllers     → internal/domain/repositories (skip a layer)
```

**Crossing features:** a service that needs another feature depends on that feature's
**repository interface** or service, injected through its `New*` constructor in
`internal/app/routers/index.go` — never by reaching for a package-level global. See
[DESIGN_PATTERNS.md](./DESIGN_PATTERNS.md) §3.6 (Dependency Injection).

### 7. API Responses - MANDATORY UTILS
```
✅ MUST use pkg/utils response functions
❌ NEVER use c.JSON() directly in handlers (for API responses)
✅ ALWAYS import "pkg/utils" in handlers
```

**When sending success responses:**
```go
// ✅ CORRECT — controller methods on a struct, using pkg/utils
import "github.com/0xdiaz/gin-boilerplate/pkg/utils"

func (ctrl *UserController) GetUser(c *gin.Context) {
    user, err := ctrl.service.GetUserByID(c.Request.Context(), id)
    if err != nil {
        utils.NotFound(c, err, "User not found")
        return
    }
    utils.Ok(c, user, "User retrieved successfully")
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
    var req dto.CreateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        utils.BadRequest(c, err, "Invalid request data")
        return
    }
    user, err := ctrl.service.CreateUser(c.Request.Context(), &req)
    if err != nil {
        utils.InternalServerError(c, err, "Failed to create user")
        return
    }
    utils.Created(c, user, "User created successfully")
}

// ❌ WRONG - direct c.JSON()
func (ctrl *UserController) GetUser(c *gin.Context) {
    user, _ := ctrl.service.GetUserByID(c.Request.Context(), id)
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
   - Which layer am I in? (controller / service / repository / model)
   - Am I following dependency direction?
   - Am I depending on a consumer-defined interface, not a concrete type?

### While Writing Code

1. **Follow the architecture**
   ```
   Controller (HTTP) → Service (Logic) → Repository (Data) → Model
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

### Controller Layer (the HTTP layer)
```go
✅ DO:
- Bind/parse the HTTP request (params, body, headers)
- Open/close a trace span (logger.LogStart / logger.LogFinish)
- Call the service it was given in its New*Controller constructor
- Map service sentinel errors to HTTP status via pkg/utils
- Respond via pkg/utils (utils.Ok, utils.Created, utils.BadRequest, …)

❌ DON'T:
- Contain business logic
- Access the database / call repositories directly
- Have functions >50 lines
- Write c.JSON(...) by hand for API responses
```

**Template** (`internal/app/controllers/auth_controller.go`):
```go
func (ctrl *AuthController) Profile(c *gin.Context) {
    // 1. Open a trace span + read input set by the middleware
    ctx, start := logger.LogStart(c.Request.Context(), "AuthController.Profile")

    userID := c.GetUint("user_id") // set by AuthMiddleware

    // 2. Call the service
    profile, err := ctrl.service.GetProfile(ctx, userID)
    if err != nil {
        if errors.Is(err, auth.ErrUserNotFound) {
            logger.LogFinish(ctx, "AuthController.Profile", err, start)
            utils.NotFound(c, err, "User not found")
            return
        }
        logger.Errorf("get profile failed: %v", err)
        logger.LogFinish(ctx, "AuthController.Profile", err, start)
        utils.InternalServerError(c, err, "Failed to retrieve profile")
        return
    }

    // 3. Respond via pkg/utils
    logger.LogFinish(ctx, "AuthController.Profile", nil, start)
    utils.Ok(c, profile, "Profile retrieved successfully")
}
// Total: ~20 lines
```

### Service Layer
```go
✅ DO:
- Implement ALL business logic
- Depend on repository INTERFACES, injected via New*Service
- Accept and propagate context.Context (LogStart / LogFinish)
- Orchestrate multiple repositories; return sentinel errors
- Transform model ↔ DTO

❌ DON'T:
- Handle HTTP concerns (gin.Context) — except a lib that needs it (e.g. DataTables paging)
- Touch database.DB directly
- Import the controller layer
- Exceed 300 lines per file (split into <name>_service_<topic>.go in the SAME package)
- Have functions >100 lines
```

**Template** (`internal/app/services/auth/auth_service.go`):
```go
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
    ctx, start := logger.LogStart(ctx, "AuthService.Register")

    // 1. Business constraint: email must be unique
    existing, err := s.userRepo.GetUserByEmail(req.Email)
    if err != nil {
        logger.LogFinish(ctx, "AuthService.Register", err, start)
        return nil, fmt.Errorf("failed to check email: %w", err)
    }
    if existing != nil {
        logger.LogFinish(ctx, "AuthService.Register", ErrEmailAlreadyExists, start)
        return nil, ErrEmailAlreadyExists // sentinel error — handler maps to 409
    }

    // 2. Apply business logic + persist
    hashed, err := s.hashPassword(req.Password)
    if err != nil {
        logger.LogFinish(ctx, "AuthService.Register", err, start)
        return nil, fmt.Errorf("failed to process password: %w", err)
    }
    user := &models.User{Name: req.Name, Email: req.Email, Password: hashed}
    if err = s.userRepo.CreateUser(user); err != nil {
        logger.LogFinish(ctx, "AuthService.Register", err, start)
        return nil, fmt.Errorf("failed to create user: %w", err)
    }

    logger.Infof("user registered successfully: %s", user.Email)
    logger.LogFinish(ctx, "AuthService.Register", nil, start)
    return &dto.AuthResponse{ /* ... */ }, nil
}
// Total: ~30 lines
```

### Repository Layer
```go
✅ DO:
- Declare an exported interface + unexported struct + New<Name>Repository()
- CRUD operations only, against the package-level database.DB handle
- Database queries / transaction management
- Translate DB-specific errors (gorm.ErrRecordNotFound → nil, nil)
- Log and wrap every error with %w
- Return models from internal/domain/models

❌ DON'T:
- Contain business logic / validate business rules
- Return gorm.ErrRecordNotFound to the caller
- Import the service or controller layer
```

**Template** (`internal/domain/repositories/user_repo.go`):
```go
func (r *userRepo) GetByID(id uint) (*models.User, error) {
    var user models.User

    if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil // not found is not an error here; the service decides
        }
        logger.Errorf("failed to get user by id: %v", err)
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return &user, nil
}
// Total: ~14 lines
```

---

## 🔍 COMMON PATTERNS

> In the patterns below, models come from `internal/domain/models`, request/response types
> from `internal/app/dto`, and the `Err*` sentinels are declared in the service that returns them.

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
            return nil, ErrEmailAlreadyExists // service sentinel error
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
        return ErrCannotDelete // service sentinel error
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
// internal/domain/repositories/<name>_repo.go — the service orchestrates WHEN to call this;
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
- [ ] Protected routes behind middlewares.AuthMiddleware(authService)
- [ ] Authorization checked where needed (using the authenticated user_id)
- [ ] Sensitive data sanitized in logs
- [ ] Rate limiting applied
- [ ] CORS configured properly

### Validation Pattern (gin binding tags, bound at the HTTP boundary)
```go
// internal/app/dto/<name>_dto.go — rules are gin binding tags (validated by ShouldBindJSON)
type CreateUserRequest struct {
    Name     string `json:"name" binding:"required,min=3,max=255"`
    Email    string `json:"email" binding:"required,email"`
    Password string `json:"password" binding:"required,min=8"`
    Age      int    `json:"age" binding:"gte=0,lte=150"`
}

// internal/app/controllers/<name>_controller.go — bind + validate in one step;
// utils.BadRequest turns validator.ValidationErrors into a field→message map.
func (ctrl *UserController) CreateUser(c *gin.Context) {
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

### Test File Structure (tests/ tree, black box, shared fake — no DB)

⚠️ Unit tests live in `tests/unit/<layer>/`, in package `<layer>_test`, driving the code
through its **exported** surface. The fake lives in `tests/mocks/` so every test can share it,
and asserts at compile time that it satisfies the real repository interface. This is the
canonical pattern; see `tests/mocks/user_repo_mock.go` and
`tests/unit/services/auth_service_test.go`.

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

// Compile-time proof the fake still matches the real interface.
var _ repositories.UserRepository = (*MockUserRepository)(nil)

func (m *MockUserRepository) GetByID(id uint) (*models.User, error) {
    if m.GetErr != nil {
        return nil, m.GetErr
    }
    m.mu.RLock()
    defer m.mu.RUnlock()

    user, ok := m.byID[id]
    if !ok {
        return nil, nil // matches the real repo: missing row is (nil, nil)
    }
    return user, nil
}
```

```go
// tests/unit/services/auth_service_test.go
package services_test

func TestAuthServiceGetProfile(t *testing.T) {
    userRepo := mocks.NewMockUserRepository()
    userRepo.AddUserByEmail("user@example.com", &models.User{ID: 1, Email: "user@example.com", Name: "User"})

    service := auth.NewAuthService(userRepo, mocks.NewMockRefreshTokenRepository(), nil)

    got, err := service.GetProfile(context.Background(), 1)

    require.NoError(t, err)
    assert.Equal(t, "user@example.com", got.Email)
}

func TestAuthServiceGet_NotFound(t *testing.T) {
    service := services.NewAuthService(mocks.NewMockUserRepository(), mocks.NewMockRefreshTokenRepository())

    got, err := service.Get(context.Background(), 999)

    assert.Nil(t, got)
    assert.True(t, errors.Is(err, services.ErrUserNotFound)) // service sentinel error
}
```

> When a fake needs to simulate failure, set its error field (`GetErr`, `ListErr`, …) rather
> than adding a new fake. Controller tests use `httptest` against the real controller with a
> service built on those fakes; see `tests/unit/controllers/`.

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

### When Creating a New Feature

The fastest path: copy the `auth` slice and rename. See MODULE_GUIDE.md
"How to add a new feature".

1. Migration — `internal/adapters/database/migrations/sql/NNNNNN_create_<table>.up.sql` plus
   the matching `.down.sql`.
2. `internal/domain/models/<name>_model.go` — the GORM model with a `TableName()` method.
3. `internal/domain/repositories/<name>_repo.go` — exported interface + unexported struct +
   `New<Name>Repository()`.
4. `internal/app/dto/<name>_dto.go` — request/response types with binding tags.
5. `internal/app/services/<name>_service.go` — `<Name>Service` holding repository interfaces +
   `New<Name>Service(...)`; declare sentinel errors here.
6. `internal/app/controllers/<name>_controller.go` — `<Name>Controller` holding the service +
   `New<Name>Controller(service)`.
7. `internal/app/routers/<name>_routes.go` — `Register<Name>Routes(group, service)`.
8. Wire it in `internal/app/routers/index.go`: build the repos and service, call
   `Register<Name>Routes(apiV1, ...)`.
9. `tests/mocks/<name>_repo_mock.go` + `tests/unit/services/<name>_service_test.go`.
10. Document all exported items, then run: `make test`.

### When Creating the Service File

1. Declare the sentinel errors this service can return (`var ErrUserNotFound = errors.New(...)`).
2. Implement the `<Name>Service` struct holding repository **interfaces** + a
   `New<Name>Service(...)` constructor.
3. Implement methods (keep <100 lines each); accept/propagate `context.Context`.
4. Use `logger.LogStart`/`LogFinish` around each exported method, span name `<Type>.<Method>`.
5. Add error handling to every function (wrap unexpected errors with `%w`).
6. Create `tests/unit/services/<name>_service_test.go` with fakes from `tests/mocks/`;
   write at least 3 cases.

### When Creating the Repository File

1. Exported `<Name>Repository` interface + unexported struct + `New<Name>Repository()`
   returning the interface.
2. Query through `database.DB` — repositories are the only layer allowed to.
3. Use GORM / parameterized queries (no raw SQL string concatenation).
4. Return models from `internal/domain/models`.
5. Translate DB errors (`gorm.ErrRecordNotFound` → `nil, nil`); log and wrap the rest with `%w`.
6. Use transactions for multi-step operations (the tx lives here, not in the service).
7. Add a matching fake in `tests/mocks/` with a `var _ repositories.<Name>Repository` assertion.

### When Creating the Controller File

1. Keep controller methods thin (<50 lines): bind request → call service → respond.
2. Hold the concrete `*services.<Name>Service`, injected via `New<Name>Controller`.
3. Map service sentinel errors with `errors.Is` before falling back to a 500.
4. Respond via `pkg/utils` (`utils.Ok`, `utils.Created`, `utils.RespondWithAPIError`, …).
5. Protect routes by mounting them behind `middlewares.AuthMiddleware(authService)` in
   `index.go`.
6. Don't put business logic here.

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
11. **Subfolder inside a layer** → Split into more files in the SAME package
12. **Skipping a layer** → A controller calls a service, never a repository

---

## 💡 TIPS FOR AI AGENTS

### Estimating Line Count Before Writing

```
Simple CRUD:
- Controller: ~15 lines
- Service: ~30 lines
- Repository: ~12 lines

Complex operation:
- Controller: ~30 lines
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
- Canonical reference slice — copy it to start a new feature:
  - `auth` — the only complete slice: `internal/domain/models/user_model.go`,
    `internal/domain/repositories/user_repo.go`, `internal/app/dto/auth_dto.go`,
    `internal/app/services/auth/auth_service.go`, `internal/app/controllers/auth_controller.go`,
    `internal/app/routers/auth_routes.go`, `tests/mocks/user_repo_mock.go`,
    `tests/unit/services/auth_service_test.go`. Its service has its own package because it
    outgrew one file (`auth_service.go` + `auth_service_tokens.go`); a smaller feature lives
    directly at `internal/app/services/<name>_service.go`.
  - `example` — the DataTables demo slice. Note its repository is a package-level function
    rather than the interface pattern above; follow `user_repo.go`, not `example_repo.go`.

---

**REMEMBER:**
- ✅ Small files = Easy to maintain
- ✅ Small functions = Easy to test
- ✅ Good tests = Confident refactoring
- ✅ Good docs = Future you will thank you

**Last updated:** 2026-06-10
**Enforcement:** MANDATORY for all commits
