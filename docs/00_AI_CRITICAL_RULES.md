# ⚠️ AI CRITICAL RULES - READ THIS FIRST

> **For AI Agents:** Read this BEFORE touching ANY code. These are NON-NEGOTIABLE rules.

> 🧭 **Structure:** Code lives in `internal/modules/<name>/` — package-by-feature, one vertical slice
> per module (`handler → service → repository → model`). See **[MODULE_GUIDE.md](./MODULE_GUIDE.md)**
> (the source of truth for the layout). The rules below — struct + constructor injection, consumer-defined
> interfaces, no standalone functions, error-as-value + APIError mapping, file-size limits — are the
> non-negotiables that apply *inside* that layout.

> 🔐 **Security posture:** Every change must align with **OWASP Top 10:2025** and **OWASP WSTG**.
> Before touching auth, input handling, PII, money movement, cryptography, or external I/O — review
> the applicable OWASP categories. The always-on hard rules live in the auto-loaded root
> **[`/CLAUDE.md`](../CLAUDE.md)**. When this service handles sensitive personal data or financial
> operations, that obligation is Tier 0.

---

## 🚨 TIER 0: ABSOLUTE RULES (NEVER VIOLATE)

### 1. Architecture Pattern (MANDATORY)

```go
// package auth (internal/modules/auth)

❌ WRONG - Standalone Functions on globals:
func Register(c *gin.Context) { }
func Login(c *gin.Context) { }

✅ CORRECT - Struct-based with constructor injection:
// handler.go — the HTTP layer is a Handler (not a "controller"); it depends on a
// consumer-defined `service` interface, never a concrete type.
type service interface {
    Register(ctx context.Context, req *RegisterRequest) (*AuthResponse, error)
    Login(ctx context.Context, req *LoginRequest) (*AuthResponse, error)
}

type Handler struct {
    svc service
}

func NewHandler(svc service) *Handler {
    return &Handler{svc: svc}
}

func (h *Handler) Register(c *gin.Context) { }
func (h *Handler) Login(c *gin.Context) { }
```

**Rule:** Handlers and Services MUST be structs with methods, wired via constructors. NO
standalone functions operating on package globals. Business logic lives on a `Service` method,
behind a constructor-injected, consumer-defined interface.

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
✅ CORRECT - Co-located with the module, white-box, fake repo (no DB):
internal/modules/example/service_test.go  (package example)

❌ WRONG - Stranded in a layer-named tests/ tree:
tests/unit/services/auth_service_test.go  (package services_test)
```

**Rule:** Unit-test a module **co-located** in `internal/modules/<name>/*_test.go`, in the
**same package** (white-box `package <module>`) so the test can implement the module's unexported
consumer-defined `repository` interface with an in-package fake (no DB). The `tests/` directory is
for integration tests and shared legacy mocks. See `internal/modules/example/service_test.go` —
this is the canonical template.

### 4. Dependency Injection (MANDATORY)

```go
❌ WRONG - Standalone func registered directly:
api.POST("/auth/register", auth.Register)   // free function on a global

✅ CORRECT - Constructor-based DI, assembled once in module.go:
// internal/modules/auth/module.go
func New(db *gorm.DB) *Module {
    svc := NewService(NewRepository(db)) // repo → service
    return &Module{svc: svc, handler: NewHandler(svc)} // service → handler
}

func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    g := api.Group("/auth")
    g.POST("/register", m.handler.Register)
}
```

**Rule:** Use constructor functions (`New*`) for dependency injection. Each layer's `New*` takes
its dependency as a consumer-defined interface; the module's `New(db)` is the single place that
assembles `repository → service → handler`.

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

### Routing (each module owns its routes)

There is no central router. A module declares its own routes in
`RegisterRoutes(api *gin.RouterGroup)`; `bootstrap` mounts every module under `/api/v1`.

```go
❌ WRONG - One giant central router with hand-wired controllers:
func RegisterRoutes(router *gin.Engine) {
    authRoutes := router.Group("/auth")
    {
        authRoutes.POST("/register", ...)
        authRoutes.POST("/login", ...)
    }
    // ... grows to 500+ lines as features pile up
}

✅ CORRECT - Each module registers its own routes:
// internal/modules/auth/module.go
func (m *Module) RegisterRoutes(api *gin.RouterGroup) {
    g := api.Group("/auth")
    g.POST("/register", m.handler.Register)
    g.POST("/login", m.handler.Login)

    protected := api.Group("")
    protected.Use(m.Middleware()) // JWT guard owned by the auth module
    protected.GET("/profile", m.handler.Profile)
}

// internal/bootstrap/modules.go — the single list of modules
func buildModules(db *gorm.DB) []Module {
    return []Module{
        auth.New(db),
        example.New(db),
        // newmodule.New(db),  ← adding a module is ONE line
    }
}
```

**Rules:**
- A module owns its routes in `RegisterRoutes(api *gin.RouterGroup)` (mounted under `/api/v1`).
- Adding a module is one line in `buildModules()` (`internal/bootstrap/modules.go`).
- System probes (`/health`, `/metrics`) are mounted at the **root** by the `health` system
  module via `RegisterSystem(r)`, not under `/api/v1`.

### Request Tracing with LogStart/LogFinish (MANDATORY in handlers)

```go
func (h *Handler) GetUser(c *gin.Context) {
    ctx, start := logger.LogStart(c.Request.Context(), "users.Handler.GetUser")
    defer func() { /* LogFinish called explicitly below */ }()

    id := c.Param("id")
    user, err := h.svc.GetUserByID(ctx, id)
    if err != nil {
        logger.LogFinish(ctx, "users.Handler.GetUser", err, start)
        utils.NotFound(c, err, "User not found")
        return
    }
    logger.LogFinish(ctx, "users.Handler.GetUser", nil, start)
    utils.Ok(c, user, "User retrieved successfully")
}
```

**Span name convention:** `<module>.<Type>.<Method>` — e.g., `auth.Handler.Login`, `auth.Service.Register`.

---

## 📁 File Structure Reference

See [MODULE_GUIDE.md](./MODULE_GUIDE.md) for the full tree. The essentials:

```
.
├── main.go                       → 3 lines: bootstrap.Run()
├── internal/
│   ├── bootstrap/                → the ONLY per-service wiring
│   │   ├── bootstrap.go          →   Run(): config → db → modules → migrate → serve
│   │   ├── modules.go            →   Module interface + buildModules() (the module list)
│   │   ├── server.go             →   gin engine + global middleware + route mounting
│   │   └── swagger.go            →   OpenAPI/Swagger UI (debug only)
│   ├── migrations/               → schema (migration.go: AutoMigrate; sql/ versioned)
│   └── modules/                  → ← business modules (work happens here)
│       ├── example/              →   THE reference module — copy it to make a new one
│       │   ├── model.go          →     GORM model(s) the module owns
│       │   ├── dto.go            →     request/response types (optional)
│       │   ├── repository.go     →     data access (holds injected *gorm.DB, no globals)
│       │   ├── service.go        →     business logic (defines the repository interface)
│       │   ├── handler.go        →     HTTP layer (defines the service interface)
│       │   ├── module.go         →     New(db), Name(), Models(), RegisterRoutes(), API()
│       │   └── service_test.go   →     co-located test (no DB — uses a fake repo)
│       ├── auth/                 →   JWT auth; exposes Middleware() and Auth() to others
│       └── health/               →   system module: /health, /metrics at ROOT
└── pkg/                          → ← cross-cutting kit; MUST never import internal/
    ├── config/  logger/  metrics/  types/  utils/  (utils/response.go → MUST use these)
    ├── middleware/               →   cors, request_id, request_log, metrics, rate_limit
    └── database/                 →   DbConnection(master, replica), GetDB()
```

> Unit tests are **co-located** in `internal/modules/<name>/*_test.go`. The `tests/` directory
> holds integration tests and shared legacy mocks only.

---

## ⚡ Quick Decision Tree

```
Writing a handler?
  → Struct + consumer-defined service interface? YES → Use response utils? YES → ✅
  → Standalone func / business logic in handler? ❌ STOP

Writing a service?
  → Struct + consumer-defined repository interface? YES
  → Co-located *_test.go with a fake repo? YES → ✅
  → No tests? ❌ STOP

Writing a repository?
  → Accepts injected *gorm.DB (not database.GetDB())? YES → ✅
  → Calling database.GetDB() directly? ❌ STOP

Returning response?
  → Using utils.Ok/Created/RespondWithAPIError? YES → ✅
  → Using c.JSON directly? ❌ STOP

Adding routes?
  → In the module's RegisterRoutes(api)? YES → ✅
  → Reaching for a central router? ❌ STOP

New module?
  → Added one line to buildModules()? YES → ✅
  → Wired it anywhere else? ❌ STOP

File approaching 250 lines?
  → Split into another file in the SAME module package? YES → ✅
  → Keep adding (or make a subfolder)? ❌ STOP
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

- Layout source of truth: [`MODULE_GUIDE.md`](./MODULE_GUIDE.md)
- Full standards: [`CODING_STANDARDS.md`](./CODING_STANDARDS.md) (read the sections marked CRITICAL)
- Design patterns: [`DESIGN_PATTERNS.md`](./DESIGN_PATTERNS.md) (read §1–§6)
- Quick templates: [`AI_QUICK_REFERENCE.md`](./AI_QUICK_REFERENCE.md)

**Critical sections in CODING_STANDARDS.md:**
- §3 Code Structure — handler → service → repository, consumer-defined interfaces
- §11.3 Response Format
- §11.4 Response Utilities (`pkg/utils`)

**Critical sections in DESIGN_PATTERNS.md:**
- §6 Implementation Patterns — Handler / Service / Repository / Module wiring
- §3.6 Dependency Injection

---

**Remember:** These are COMPANY STANDARDS. Violation = Code Rejected.
