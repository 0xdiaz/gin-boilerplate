# 📚 Documentation Guide

Welcome to the project documentation. This guide helps you navigate all documentation files.

> 🧭 **Architecture: layered service.** Code is organized **by technical layer** —
> `internal/app/{controllers,dto,middlewares,routers,services}` and
> `internal/domain/{models,repositories}`, with the database adapter in
> `internal/adapters/database` and cross-cutting code in `pkg/`. A feature is one file per
> layer, sharing a name prefix (`auth_controller.go`, `auth_service.go`, `user_repo.go`).
> **[`MODULE_GUIDE.md`](./MODULE_GUIDE.md) is the source of truth** for the layout.

---

## 🚨 FOR AI AGENTS - START HERE

### Reading Order (MANDATORY)

0. **[`MODULE_GUIDE.md`](./MODULE_GUIDE.md)** 🧭 **SOURCE OF TRUTH** (read first for structure)
   - How code is organized: one folder per layer under `internal/app/` and `internal/domain/`
   - The dependency direction: `controller → service → repository → model`
   - How to add a feature (copy the `auth` reference slice)

1. **[`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md)** ⚠️ **START HERE for rules** (100 lines, 2 min)
   - Absolute non-negotiable rules
   - Struct-based patterns (MANDATORY)
   - Response utilities (MANDATORY)
   - Test location rules (MANDATORY)
   - **READ THIS FIRST OR YOUR CODE WILL BE REJECTED**

2. **[`AI_QUICK_REFERENCE.md`](./AI_QUICK_REFERENCE.md)** (405 lines, 5 min)
   - Quick templates for controllers, services, repositories
   - The 5 Commandments (file/function size limits)
   - Testing checklist
   - Common patterns

3. **[`DOCS_INDEX.md`](./DOCS_INDEX.md)** ⭐ **NEW - BOOKMARK THIS** (comprehensive index)
   - Quick keyword lookup table
   - Line number references for all topics
   - Task-based navigation (e.g., "I need to create a controller")
   - Common task quick links
   - Full document structure overview
   - **USE THIS FOR QUICK LOOKUPS DURING WORK**

4. **Use as Reference (Ctrl+F + line numbers):**
   - **[`CODING_STANDARDS.md`](./CODING_STANDARDS.md)** - Enhanced TOC with line numbers & keywords
   - **[`DESIGN_PATTERNS.md`](./DESIGN_PATTERNS.md)** - Enhanced TOC with line numbers & keywords

---

## 📖 Documentation Files

### For AI Agents

| File | Size | Purpose | When to Read |
|------|------|---------|--------------|
| **MODULE_GUIDE.md** 🧭 | ~180 lines | **Project structure (source of truth):** layered layout, how to add a feature | **FIRST - for structure** |
| **00_AI_CRITICAL_RULES.md** | 100 lines | Non-negotiable rules | **FIRST - for rules** |
| **AI_AGENT_RULES.md** | ~700 lines | Mandatory rules for AI (file/function size, testing, docs, errors) | **IMPORTANT** |
| **AI_QUICK_REFERENCE.md** | 405 lines | Templates & checklists | Before writing code |
| **DOCS_INDEX.md** ⭐ | ~500 lines | Master index with line refs | **BOOKMARK - use during work** |
| **CODING_STANDARDS.md** | 2,200+ lines | Complete coding standards | Reference (use Ctrl+F + line numbers) |
| **DESIGN_PATTERNS.md** | 2,600+ lines | Architecture patterns | Reference (use Ctrl+F + line numbers) |

### For Developers

| File | Purpose |
|------|---------|
| **MODULE_GUIDE.md** 🧭 | Source of truth for project structure: the layered layout, how to add a feature, wiring |
| **CODING_STANDARDS.md** | Comprehensive coding standards, naming conventions, best practices |
| **DESIGN_PATTERNS.md** | Architecture patterns, layer responsibilities, implementation guides |
| **CONFIGURATION.md** | Environment configuration, validation, secrets management |
| **OBSERVABILITY.md** | Health checks, metrics, request tracing, monitoring guide |
| **AUTHENTICATION.md** | Auth flows, JWT, password reset |
| **CONTRACTS.md** | Stable contracts (response, auth, env); versioning policy — avoid breaking changes |
| **AI_QUICK_REFERENCE.md** | Quick templates and decision trees |
| **00_AI_CRITICAL_RULES.md** | Quick reference for critical rules |
| **AI_AGENT_RULES.md** | Mandatory rules for AI agents (file size, testing, docs, errors) |
| **OPENTELEMETRY_TRACING_ANALYSIS.md** | Optional OpenTelemetry analysis; request_id is sufficient for monolith |
| **CONTROLLER_COMPLIANCE_AUDIT.md** | Controller compliance checklist |
| **SERVICE_COMPLIANCE_AUDIT.md** | Service compliance checklist |

---

## 🎯 Quick Navigation

### I Want To...

**Add a new feature (start here):**
1. Read: `MODULE_GUIDE.md` → "How to add a new feature"
2. Copy the reference slice: `auth` (model, repo, dto, service, controller, routes, mock, test)
3. Wire it: build the repo + service in `internal/app/routers/index.go`, then call
   `Register<Name>Routes(apiV1, ...)`

**Write a new controller (the HTTP layer):**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1 — struct-based + DI)
2. Reference: `internal/app/controllers/auth_controller.go`
3. Patterns: `AI_QUICK_REFERENCE.md` → Layer Templates; `DESIGN_PATTERNS.md` §6.1

**Write a new service (business logic):**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1 — struct-based + DI)
2. Reference: `internal/app/services/auth/auth_service.go`
3. Patterns: `AI_QUICK_REFERENCE.md` → Templates section; `DESIGN_PATTERNS.md` (layered terminology, see MODULE_GUIDE for the modular mapping)

**Return a response:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 2)
2. Utils: `pkg/utils/response.go`
3. Details: `CODING_STANDARDS.md` → § 11 API Design (Response Utilities)

**Write tests:**
1. Read: `MODULE_GUIDE.md` → Testing (co-locate `*_test.go` with the module; fake repo, no DB)
2. Reference: `tests/unit/services/auth_service_test.go`
3. Guide: `tests/README.md` (shared mocks + legacy/integration tests)

**Handle errors:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 2)
2. Details: `CODING_STANDARDS.md` → Error Handling section

**Use dependency injection:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 4)
2. Details: `DESIGN_PATTERNS.md` → § Core Patterns (Dependency Injection); DI here: each layer's `New*` constructor, assembled in `routers/index.go`

**Setup observability (health checks, metrics):**
1. Guide: `OBSERVABILITY.md`
2. Endpoints: GET /health, GET /metrics
3. Implementation: Middleware-based, <1% overhead

**Configure environment:**
1. Guide: `CONFIGURATION.md`
2. Validation: Startup fail-fast
3. Environment helpers: IsDevelopment(), IsProduction()

**View or update API documentation (Swagger):**
1. When `DEBUG=true`, open `/swagger/` in the browser (e.g. http://localhost:8000/swagger/).
2. The source of truth is [api/openapi.yaml](../api/openapi.yaml); edit that file to change the docs (no annotations in controllers).

**Check what is stable (avoid breaking consumers):**
1. Read: [CONTRACTS.md](./CONTRACTS.md) — response shape, auth header, path prefix, env keys, versioning policy.

---

## 🚫 Common Mistakes (Don't Do This)

### Mistake #1: Not Reading Critical Rules
```
❌ Skipping 00_AI_CRITICAL_RULES.md
✅ Reading it first (takes 3 minutes)
```

**Result of skipping:** Code rejected, need to refactor everything.

### Mistake #2: Using Standalone Functions
```go
❌ func Register(ctx *gin.Context) { }  // Rejected
✅ func (ctrl *AuthController) Register(c *gin.Context) { }
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 0, Rule 1

### Mistake #3: Direct c.JSON() Calls
```go
❌ c.JSON(200, gin.H{"data": user})  // Rejected
✅ utils.Ok(c, user, "Success")
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 0, Rule 2

### Mistake #4: Tests in the Wrong Place
```
✅ tests/unit/services/auth_service_test.go     // unit test, package services_test
✅ tests/mocks/user_repo_mock.go                // the shared fake it drives
❌ internal/app/services/auth_service_test.go    // tests do not live beside the code here
```

**Why:** unit tests live in the root `tests/` tree, in package `<layer>_test`, so they exercise
the code through its exported surface exactly as production callers do. The fakes live in
`tests/mocks/` so one interface change surfaces in one place — each carries a
`var _ repositories.UserRepository = (*MockUserRepository)(nil)` assertion. See
`MODULE_GUIDE.md` → Testing and `tests/unit/services/auth_service_test.go`.

### Mistake #5: Exceeding File Size
```
❌ File with 400 lines  // Rejected
✅ Split into multiple focused files (max 300 lines)
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 1

---

## 📊 Documentation Statistics

- **Total Documentation:** ~5,000+ lines (after cleanup of redundant docs)
- **Critical Rules:** 100 lines
- **Must Read Before Coding:** 505 lines
- **Quick Lookup Index:** DOCS_INDEX.md
- **Reference Material:** CODING_STANDARDS, DESIGN_PATTERNS, CONFIGURATION, OBSERVABILITY, etc.

**Efficiency Tip:**
1. Read critical rules + quick ref first (~10 min)
2. Bookmark DOCS_INDEX.md for quick lookups
3. Use enhanced TOCs in full docs with Ctrl+F + line numbers

---

## 🔄 Document Updates

**Last Updated:** 2026-02-03

**Recent Changes:**
- **Docs realigned with the code:** the doc set previously described a package-by-feature layout (`internal/modules/<name>/`, `internal/bootstrap/`, co-located tests) that was never implemented. Every doc now describes the layered layout that actually exists. See `MODULE_GUIDE.md` (source of truth), whose "Planned direction" section records the unbuilt modular design.
- API versioning (`/api/v1`), global rate limit, single config source, request_id/LogStart/LogFinish logging, pluggable EmailSender (see "Recent changes" above).
- `AI_AGENT_RULES.md` is kept (important for AI agents).
- Added "New developer onboarding" path in this file.
- `OPENTELEMETRY_TRACING_ANALYSIS.md` — optional OpenTelemetry analysis; request_id is sufficient for monolith.

---

## ✅ Checklist Before First Code Contribution

```
□ Read MODULE_GUIDE.md (project structure — source of truth)
□ Read 00_AI_CRITICAL_RULES.md (100 lines)
□ Read AI_QUICK_REFERENCE.md (405 lines)
□ Understand the layered layout (internal/app/* and internal/domain/*)
□ Understand struct-based pattern requirement
□ Understand response utilities requirement
□ Understand where tests live (tests/unit/<layer>/, fakes in tests/mocks/)
□ Know file size limits (300 lines max)
□ Know function size limits (100 lines max)
```

**Time Required:** 15-20 minutes
**Time Saved:** Hours of refactoring

---

## 📞 Questions?

If documentation is unclear:
1. Check `00_AI_CRITICAL_RULES.md` first
2. Search in `CODING_STANDARDS.md` or `DESIGN_PATTERNS.md`
3. Look for examples in existing code
4. Ask the team

**For AI Agents:** If you're unsure, ASK. Don't guess and violate critical rules.
