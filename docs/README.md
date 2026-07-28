# 📚 Documentation Guide

Welcome to the project documentation. This guide helps you navigate all documentation files.

> 🧭 **Architecture: modular service.** Code is organized **by business module** under
> `internal/modules/<name>/` (each owning `model.go`, `dto.go`, `repository.go`, `service.go`,
> `handler.go`, `module.go`). Cross-cutting code lives in `pkg/`; per-service wiring in
> `internal/bootstrap/`. **[`MODULE_GUIDE.md`](./MODULE_GUIDE.md) is the source of truth** and
> overrides any older doc (DESIGN_PATTERNS, CODING_STANDARDS, AI_AGENT_RULES) that still describes
> the previous `internal/app` / `internal/domain` layered layout.

---

## 🚨 FOR AI AGENTS - START HERE

### Reading Order (MANDATORY)

0. **[`MODULE_GUIDE.md`](./MODULE_GUIDE.md)** 🧭 **SOURCE OF TRUTH** (read first for structure)
   - How code is organized: one folder per business module under `internal/modules/`
   - The vertical slice: `handler → service → repository → model`
   - How to add a module (copy the `example` reference module)
   - Overrides the layered-layout descriptions in the deeper reference docs

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
| **MODULE_GUIDE.md** 🧭 | ~100 lines | **Project structure (source of truth):** modular layout under `internal/modules/`; overrides layered-layout docs | **FIRST - for structure** |
| **00_AI_CRITICAL_RULES.md** | 100 lines | Non-negotiable rules | **FIRST - for rules** |
| **AI_AGENT_RULES.md** | ~700 lines | Mandatory rules for AI (file/function size, testing, docs, errors) | **IMPORTANT** |
| **AI_QUICK_REFERENCE.md** | 405 lines | Templates & checklists | Before writing code |
| **DOCS_INDEX.md** ⭐ | ~500 lines | Master index with line refs | **BOOKMARK - use during work** |
| **CODING_STANDARDS.md** | 2,200+ lines | Complete coding standards | Reference (use Ctrl+F + line numbers) |
| **DESIGN_PATTERNS.md** | 2,600+ lines | Architecture patterns | Reference (use Ctrl+F + line numbers) |

### For Developers

| File | Purpose |
|------|---------|
| **MODULE_GUIDE.md** 🧭 | Source of truth for project structure: modular layout (`internal/modules/<name>/`), how to add a module, cross-module communication |
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

**Add a new module (the modular way — start here):**
1. Read: `MODULE_GUIDE.md` → "How to add a new module"
2. Copy the reference module: `internal/modules/example/`
3. Register it: one line in `buildModules()` in `internal/bootstrap/modules.go`

**Write a new handler (HTTP layer of a module):**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1 — struct-based + DI)
2. Reference: `internal/modules/example/handler.go`
3. Patterns: `AI_QUICK_REFERENCE.md` → Templates section; `DESIGN_PATTERNS.md` (layered terminology, see MODULE_GUIDE for the modular mapping)

**Write a new service (business logic of a module):**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1 — struct-based + DI)
2. Reference: `internal/modules/example/service.go`
3. Patterns: `AI_QUICK_REFERENCE.md` → Templates section; `DESIGN_PATTERNS.md` (layered terminology, see MODULE_GUIDE for the modular mapping)

**Return a response:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 2)
2. Utils: `pkg/utils/response.go`
3. Details: `CODING_STANDARDS.md` → § 11 API Design (Response Utilities)

**Write tests:**
1. Read: `MODULE_GUIDE.md` → Testing (co-locate `*_test.go` with the module; fake repo, no DB)
2. Reference: `internal/modules/example/service_test.go`
3. Guide: `tests/README.md` (shared mocks + legacy/integration tests)

**Handle errors:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 2)
2. Details: `CODING_STANDARDS.md` → Error Handling section

**Use dependency injection:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 4)
2. Details: `DESIGN_PATTERNS.md` → § Core Patterns (Dependency Injection); modular DI: each module's `New(db)` constructor

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
✅ internal/modules/auth/service_test.go        // Co-located with the module (preferred)
✅ tests/unit/...                                // Legacy/shared tests still live here
```

**Why:** Under the modular layout, **tests are co-located** with the module they cover
(`internal/modules/<name>/*_test.go`) — see `MODULE_GUIDE.md` → Testing and
`internal/modules/example/service_test.go`. Older docs that say "tests MUST live in `tests/`
and co-located tests are rejected" describe the pre-refactor layered layout; `MODULE_GUIDE.md`
overrides them. New modules should prefer in-package fakes and co-located tests; `tests/` remains
for shared mocks and legacy/integration tests.

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
- **Modular refactor:** code is now organized by business module under `internal/modules/<name>/` (was `internal/app` / `internal/domain` layered layout). See `MODULE_GUIDE.md` (source of truth). Tests are now co-located with their module. The `CONTROLLER_COMPLIANCE_AUDIT.md` and `SERVICE_COMPLIANCE_AUDIT.md` docs are historical (they audited the pre-refactor layout).
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
□ Understand the modular layout (one folder per module under internal/modules/)
□ Understand struct-based pattern requirement
□ Understand response utilities requirement
□ Understand co-located tests (internal/modules/<name>/*_test.go)
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
