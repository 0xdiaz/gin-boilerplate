# ⚠️ AI CRITICAL RULES - READ THIS FIRST

> **For AI Agents:** Read this BEFORE touching ANY code. These are NON-NEGOTIABLE rules.

---

## 🚨 TIER 0: ABSOLUTE RULES (NEVER VIOLATE)

### 1. Architecture Pattern (MANDATORY)

```go
❌ WRONG - Standalone Functions:
func Register(ctx *gin.Context) { }
func Login(ctx *gin.Context) { }

✅ CORRECT - Struct-based with DI:
type AuthController struct {
    service *services.AuthService
}

func NewAuthController(service *services.AuthService) *AuthController {
    return &AuthController{service: service}
}

func (ctrl *AuthController) Register(c *gin.Context) { }
func (ctrl *AuthController) Login(c *gin.Context) { }
```

**Rule:** Controllers and Services MUST be structs with methods. NO standalone functions.

### 2. Response Format (MANDATORY)

```go
❌ WRONG - Direct gin.H:
c.JSON(200, gin.H{"status": 200, "data": user})

✅ CORRECT - Use Response Utilities:
utils.Ok(c, user, "User retrieved successfully")
utils.Created(c, user, "User created successfully")
utils.BadRequest(c, err, "Invalid input")
utils.Unauthorized(c, err, "Invalid credentials")
```

**Standard Format:**
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {...},
  "errors": null
}
```

**Rule:** ALL responses MUST use `pkg/utils/response.go` utilities. NO direct c.JSON().

### 3. Test Location (MANDATORY)

```
❌ WRONG - Co-located:
internal/app/services/auth_service_test.go

✅ CORRECT - In tests/ directory:
tests/unit/services/auth_service_test.go (package services_test)
```

**Rule:** ALL tests in `tests/` directory with `_test` package suffix.

### 4. Dependency Injection (MANDATORY)

```go
❌ WRONG - Direct instantiation:
authRoutes.POST("/register", controllers.Register)

✅ CORRECT - Constructor-based DI:
authService := services.NewAuthService()
authController := controllers.NewAuthController(authService)
authRoutes.POST("/register", authController.Register)
```

**Rule:** Use constructor functions (New*) for dependency injection.

---

## 🔥 TIER 1: HARD LIMITS (EXCEED = REJECT CODE)

```
File Size:     MAX 300 lines  (warning at 250)
Function Size: MAX 100 lines  (warning at 80)
Test Coverage: MIN 70% for services
```

---

## 📍 TIER 2: CRITICAL PATTERNS

### Response Utilities (pkg/utils/response.go)

```go
// Success responses
utils.Ok(c, data, message)              // 200
utils.Created(c, data, message)          // 201
utils.NoContent(c)                       // 204

// Error responses
utils.BadRequest(c, err, message)        // 400
utils.Unauthorized(c, err, message)      // 401
utils.Forbidden(c, err, message)         // 403
utils.NotFound(c, err, message)          // 404
utils.Conflict(c, err, message)          // 409
utils.InternalServerError(c, err, msg)   // 500
```

### Error Handling

```go
❌ WRONG:
_, _ = someFunc()  // Ignored error
if err != nil {
    return
}

✅ CORRECT:
result, err := someFunc()
if err != nil {
    logger.Errorf("operation failed: %v", err)
    return fmt.Errorf("failed to do X: %w", err)
}
```

### Logging

```go
❌ WRONG:
log.Printf("User created")
fmt.Println("Error:", err)

✅ CORRECT:
logger.Infof("user created: ID=%d, Email=%s", user.ID, user.Email)
logger.Errorf("failed to create user: %v", err)
logger.Warnf("approaching rate limit: %d/%d", current, limit)
```

### Router Organization

```go
❌ WRONG - All routes in index.go:
// internal/app/routers/index.go (500+ lines)
func RegisterRoutes(router *gin.Engine) {
    authRoutes := router.Group("/auth")
    {
        authRoutes.POST("/register", ...)
        authRoutes.POST("/login", ...)
    }
    userRoutes := router.Group("/users")
    {
        // ... 50+ routes ...
    }
    // ... becomes 500+ lines
}

✅ CORRECT - Separate files by feature:
// internal/app/routers/auth_routes.go (40 lines)
func RegisterAuthRoutes(router *gin.Engine, authService *services.AuthService) {
    authController := controllers.NewAuthController(authService)
    authRoutes := router.Group("/auth")
    {
        authRoutes.POST("/register", authController.Register)
        authRoutes.POST("/login", authController.Login)
    }
}

// internal/app/routers/index.go (50 lines)
func RegisterRoutes(router *gin.Engine) {
    authService := services.NewAuthService()
    RegisterAuthRoutes(router, authService)
    RegisterUserRoutes(router, userService, authService)
}
```

**Rules:**
- One file per controller/feature: `{feature}_routes.go`
- Function naming: `Register{Feature}Routes()`
- Max 100 lines per route file
- Main `index.go` only calls Register functions

---

## 📁 File Structure Reference

```
project/
├── internal/
│   ├── app/
│   │   ├── controllers/       → Struct-based, use response utils
│   │   ├── services/          → Struct-based, business logic
│   │   ├── dto/              → Request/Response structs
│   │   ├── middlewares/      → Gin middleware functions
│   │   └── routers/          → Route registration (ONE FILE PER FEATURE)
│   │       ├── index.go      → Main router (calls all Register functions)
│   │       ├── auth_routes.go    → Auth routes only
│   │       ├── user_routes.go    → User routes only
│   │       └── product_routes.go → Product routes only
│   └── domain/
│       ├── models/           → GORM entities
│       └── repositories/     → Function-based CRUD
├── pkg/
│   └── utils/
│       └── response.go       → MUST use these utilities
└── tests/                    → ALL tests here
    ├── unit/
    │   ├── controllers/
    │   ├── services/
    │   └── repositories/
    └── integration/
```

---

## ⚡ Quick Decision Tree

```
Writing a controller?
  → Struct-based? YES → Use response utils? YES → ✅
  → Standalone func? ❌ STOP

Writing a service?
  → Struct-based? YES → Has tests in tests/? YES → ✅
  → No tests? ❌ STOP

Returning response?
  → Using utils.Ok/Created/etc? YES → ✅
  → Using c.JSON directly? ❌ STOP

Adding routes?
  → Separate {feature}_routes.go file? YES → ✅
  → All in index.go? ❌ STOP

File approaching 250 lines?
  → Split now? YES → ✅
  → Keep adding? ❌ STOP
```

---

## Git: Do not commit non-essential .md files

**Rule:** Markdown files that are **local, analysis-only, or temporary** must not be committed or pushed.

**Do not commit (examples):**
- `PROJECT_ANALYSIS.md` — project analysis/score (local only)
- Draft docs, personal notes, or .md used only for internal reference and not part of the shared project

**Do commit:** All files in `docs/` that are part of the project standard (CODING_STANDARDS, DESIGN_PATTERNS, OBSERVABILITY, CONFIGURATION, AI rules, README, .env.example, etc.).

Non-essential files are listed in `.gitignore` (e.g. `PROJECT_ANALYSIS.md`). Before committing, ensure no new analysis/local .md files are staged.

---

## 📚 For More Details

- Full standards: `CODING_STANDARDS.md` (read sections marked CRITICAL)
- Design patterns: `DESIGN_PATTERNS.md` (read sections 1-4)
- Quick templates: `AI_QUICK_REFERENCE.md`

**Critical sections in CODING_STANDARDS.md:**
- Lines 900-1100: Struct-based patterns
- Lines 1429-1475: Response format
- Lines 1479-1584: Response utilities

**Critical sections in DESIGN_PATTERNS.md:**
- Lines 900-1016: Controller & Service patterns
- Lines 439-492: Dependency injection

---

**Remember:** These are COMPANY STANDARDS. Violation = Code Rejected.
