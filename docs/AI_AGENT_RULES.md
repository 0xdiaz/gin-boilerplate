# AI Agent Rules

> **CRITICAL**: These rules are MANDATORY for all AI agents working on this codebase.
> Violating these rules will result in rejected code.

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

### 6. DTO and Enum Placement - STRICT RULES
```
✅ ALL STRUCTS → internal/app/dto/
✅ ALL CONSTANTS → pkg/enums/
❌ EXCEPTION: Maps using DTOs → services/ (avoid import cycles)
```

**When creating any struct:**
```
✅ CORRECT: internal/app/dto/hub_dto.go
❌ WRONG: internal/app/services/hub_types.go
❌ WRONG: pkg/types/hub_types.go
```

**When creating any constant:**
```
✅ CORRECT: pkg/enums/hub_constants.go
❌ WRONG: internal/app/services/constants.go
❌ WRONG: internal/app/controllers/constants.go
```

**IMPORT CYCLE PREVENTION:**
```
✅ ALLOWED:
   internal/app/dto → pkg/enums
   internal/app/services → pkg/enums
   internal/app/services → internal/app/dto

❌ FORBIDDEN:
   pkg/enums → internal/app/dto         (import cycle!)
   pkg/enums → internal/app/services    (import cycle!)
```

**If you need a map using DTOs:**
```go
// ❌ WRONG - causes import cycle
// File: pkg/enums/notification_constants.go
import "internal/app/dto"
var NotificationRouteMapping = map[string]dto.NotificationRoute{...}

// ✅ CORRECT - keep in services
// File: internal/app/services/notification_router.go
import "internal/app/dto"
var NotificationRouteMapping = map[string]dto.NotificationRoute{...}
```

### 7. API Responses - MANDATORY UTILS
```
✅ MUST use pkg/utils response functions
❌ NEVER use c.JSON() directly in controllers
✅ ALWAYS import "pkg/utils" in controllers
```

**When sending success responses:**
```go
// ✅ CORRECT
import "github.com/.../pkg/utils"

func GetUser(c *gin.Context) {
    user, err := service.GetUserByID(id)
    if err != nil {
        utils.NotFound(c, err, "User not found")
        return
    }
    utils.Ok(c, user, "User retrieved successfully")
}

func CreateUser(c *gin.Context) {
    user, err := service.CreateUser(dto)
    if err != nil {
        utils.BadRequest(c, err, "Failed to create user")
        return
    }
    utils.Created(c, user, "User created successfully")
}

// ❌ WRONG - direct c.JSON()
func GetUser(c *gin.Context) {
    user, _ := service.GetUserByID(id)
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
   - What layer am I in? (controller/service/repository)
   - Am I following dependency direction?
   - Am I importing the right packages?

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

### Controller Layer
```go
✅ DO:
- Parse HTTP request (params, body, headers)
- Call service method
- Return HTTP response
- Handle HTTP status codes

❌ DON'T:
- Contain business logic
- Access database directly
- Have functions >50 lines
- Import repository packages
```

**Template:**
```go
func GetUser(c *gin.Context) {
    // 1. Parse input
    id := c.Param("id")

    // 2. Call service
    user, err := service.GetUserByID(id)
    if err != nil {
        utils.HandleErrors(c, http.StatusNotFound, nil, err.Error())
        return
    }

    // 3. Return response
    c.JSON(http.StatusOK, user)
}
// Total: ~15 lines
```

### Service Layer
```go
✅ DO:
- Implement ALL business logic
- Validate input
- Orchestrate multiple repositories
- Transform data between layers
- Handle business errors

❌ DON'T:
- Handle HTTP concerns (gin.Context)
- Import controller packages
- Exceed 400 lines per file
- Have functions >100 lines
```

**Template:**
```go
func (s *UserService) CreateUser(dto dto.CreateUserRequest) (*models.User, error) {
    // 1. Validate
    if err := s.validate(dto); err != nil {
        return nil, fmt.Errorf("validation failed: %w", err)
    }

    // 2. Check duplicates
    existing, _ := s.userRepo.FindByEmail(dto.Email)
    if existing != nil {
        return nil, types.ErrDuplicateEntry
    }

    // 3. Transform
    user := s.buildUser(dto)

    // 4. Save
    if err := s.userRepo.Create(user); err != nil {
        logger.Errorf("failed to create user: %v", err)
        return nil, fmt.Errorf("create user: %w", err)
    }

    logger.Infof("User created: ID=%d, Email=%s", user.ID, user.Email)
    return user, nil
}
// Total: ~30 lines
```

### Repository Layer
```go
✅ DO:
- CRUD operations only
- Database queries
- Transaction management
- Return domain models

❌ DON'T:
- Contain business logic
- Validate business rules
- Call other repositories
- Import service packages
```

**Template:**
```go
func (r *userRepository) FindByID(id uint) (*models.User, error) {
    var user models.User

    if err := r.db.First(&user, id).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, types.ErrNotFound
        }
        return nil, fmt.Errorf("query failed: %w", err)
    }

    return &user, nil
}
// Total: ~12 lines
```

---

## 🔍 COMMON PATTERNS

### 1. List with Pagination
```go
func (s *UserService) ListUsers(page, pageSize int) ([]*models.User, int64, error) {
    // Validate pagination
    if page < 1 {
        page = 1
    }
    if pageSize < 1 || pageSize > 100 {
        pageSize = 20
    }

    // Get data
    users, total, err := s.userRepo.List(page, pageSize)
    if err != nil {
        return nil, 0, fmt.Errorf("list users: %w", err)
    }

    return users, total, nil
}
```

### 2. Update Operations
```go
func (s *UserService) UpdateUser(id uint, dto dto.UpdateUserRequest) (*models.User, error) {
    // 1. Get existing
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return nil, fmt.Errorf("get user: %w", err)
    }

    // 2. Check business rules
    if dto.Email != user.Email {
        existing, _ := s.userRepo.FindByEmail(dto.Email)
        if existing != nil {
            return nil, types.ErrDuplicateEntry
        }
    }

    // 3. Update fields
    user.Name = dto.Name
    user.Email = dto.Email

    // 4. Save
    if err := s.userRepo.Update(user); err != nil {
        return nil, fmt.Errorf("update user: %w", err)
    }

    return user, nil
}
```

### 3. Delete Operations
```go
func (s *UserService) DeleteUser(id uint) error {
    // 1. Check exists
    user, err := s.userRepo.FindByID(id)
    if err != nil {
        return fmt.Errorf("get user: %w", err)
    }

    // 2. Check business rules (e.g., can't delete if has active transactions)
    hasTransactions, err := s.transactionRepo.HasActiveTransactions(id)
    if err != nil {
        return fmt.Errorf("check transactions: %w", err)
    }
    if hasTransactions {
        return types.ErrCannotDelete
    }

    // 3. Delete
    if err := s.userRepo.Delete(user.ID); err != nil {
        return fmt.Errorf("delete user: %w", err)
    }

    logger.Infof("User deleted: ID=%d", id)
    return nil
}
```

### 4. Transaction Pattern
```go
func (s *TransactionService) ProcessPayment(dto dto.PaymentRequest) error {
    return s.db.Transaction(func(tx *gorm.DB) error {
        // All operations in transaction
        // If any fails, all rollback

        // Step 1
        if err := s.createTransaction(tx, dto); err != nil {
            return err
        }

        // Step 2
        if err := s.updateBalance(tx, dto.UserID, dto.Amount); err != nil {
            return err
        }

        // Step 3
        if err := s.createAuditLog(tx, dto); err != nil {
            return err
        }

        return nil
    })
}
```

---

## 🛡️ SECURITY CHECKLIST

### Before Committing

- [ ] No hardcoded passwords/secrets
- [ ] No SQL string concatenation
- [ ] All input validated
- [ ] Passwords hashed with bcrypt
- [ ] Authentication checked in controllers
- [ ] Authorization checked where needed
- [ ] Sensitive data sanitized in logs
- [ ] Rate limiting applied
- [ ] CORS configured properly (never combine `Access-Control-Allow-Origin: *` with `Access-Control-Allow-Credentials`; use an explicit origin allowlist via `CORS_ALLOWED_ORIGINS`)

### Validation Pattern
```go
type CreateUserRequest struct {
    Name     string `json:"name" validate:"required,min=3,max=255"`
    Email    string `json:"email" validate:"required,email"`
    Password string `json:"password" validate:"required,min=8"`
    Age      int    `json:"age" validate:"gte=0,lte=150"`
}

func (s *UserService) validateCreateUser(dto CreateUserRequest) error {
    validate := validator.New()
    if err := validate.Struct(dto); err != nil {
        return fmt.Errorf("validation failed: %w", err)
    }
    return nil
}
```

---

## 📝 TESTING PATTERNS

### Test File Structure
```go
package services_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

// Mock repository
type MockUserRepository struct {
    mock.Mock
}

func (m *MockUserRepository) FindByID(id uint) (*models.User, error) {
    args := m.Called(id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*models.User), args.Error(1)
}

// Tests
func TestCreateUser_Success(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    dto := dto.CreateUserRequest{
        Name:  "John Doe",
        Email: "john@example.com",
    }

    mockRepo.On("FindByEmail", dto.Email).Return(nil, types.ErrNotFound)
    mockRepo.On("Create", mock.Anything).Return(nil)

    // Act
    user, err := service.CreateUser(dto)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, dto.Name, user.Name)
    mockRepo.AssertExpectations(t)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
    // Arrange
    mockRepo := new(MockUserRepository)
    service := NewUserService(mockRepo)

    dto := dto.CreateUserRequest{
        Name:  "John Doe",
        Email: "existing@example.com",
    }

    existingUser := &models.User{ID: 1, Email: dto.Email}
    mockRepo.On("FindByEmail", dto.Email).Return(existingUser, nil)

    // Act
    user, err := service.CreateUser(dto)

    // Assert
    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Equal(t, types.ErrDuplicateEntry, err)
}
```

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

### When Creating New Service

1. Create service file: `user_service.go`
2. Define interface (if needed)
3. Implement struct + constructor
4. Implement methods (keep <100 lines each)
5. Add error handling to every function
6. Add logging for important operations
7. Create test file: `user_service_test.go`
8. Write tests (minimum 3 cases)
9. Document all exported items
10. Run: `go test ./...`

### When Creating New Repository

1. Define interface in `/internal/domain/repositories/`
2. Implement in `/internal/adapters/database/`
3. Use GORM for queries (no raw SQL)
4. Return domain models
5. Handle GORM errors properly
6. Use transactions for multi-step operations

### When Creating New Controller

1. Keep functions thin (<50 lines)
2. Parse request → Call service → Return response
3. Handle HTTP status codes properly
4. Use middleware for auth/validation
5. Don't put business logic here

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
7. **Business logic in controller** → Move to service
8. **Database access in controller** → Use repository
9. **Hardcoded secrets** → Move to .env
10. **SQL string concatenation** → Use GORM/parameterized queries

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
- Service: ~80 lines (if exceeds, split!)
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

- Full standards: `CODING_STANDARDS.md`
- Architecture: See "Code Quality Analysis" in project docs
- Examples: Look at existing well-written files:
  - Small controllers: `internal/app/controllers/user_controller.go`
  - Good services: Check files <300 lines
  - Good repos: Most repository files

---

**REMEMBER:**
- ✅ Small files = Easy to maintain
- ✅ Small functions = Easy to test
- ✅ Good tests = Confident refactoring
- ✅ Good docs = Future you will thank you

**Last updated:** 2025-11-08
**Enforcement:** MANDATORY for all commits
