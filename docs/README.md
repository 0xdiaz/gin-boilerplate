# 📚 Documentation Guide

Welcome to the project documentation. This guide helps you navigate all documentation files.

---

## 🚨 FOR AI AGENTS - START HERE

### Reading Order (MANDATORY)

1. **[`00_AI_CRITICAL_RULES.md`](./00_AI_CRITICAL_RULES.md)** ⚠️ **START HERE** (100 lines, 2 min)
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
| **00_AI_CRITICAL_RULES.md** | 100 lines | Non-negotiable rules | **FIRST - ALWAYS** |
| **AI_AGENT_RULES.md** | ~700 lines | Mandatory rules for AI (file/function size, testing, docs, errors) | **IMPORTANT** |
| **AI_QUICK_REFERENCE.md** | 405 lines | Templates & checklists | Before writing code |
| **DOCS_INDEX.md** ⭐ | ~500 lines | Master index with line refs | **BOOKMARK - use during work** |
| **CODING_STANDARDS.md** | 2,200+ lines | Complete coding standards | Reference (use Ctrl+F + line numbers) |
| **DESIGN_PATTERNS.md** | 2,600+ lines | Architecture patterns | Reference (use Ctrl+F + line numbers) |

### For Developers

| File | Purpose |
|------|---------|
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

**Write a new controller:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1)
2. Template: `AI_QUICK_REFERENCE.md` → Templates section
3. Details: `DESIGN_PATTERNS.md` lines 900-948

**Write a new service:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 1)
2. Template: `AI_QUICK_REFERENCE.md` → Templates section
3. Details: `DESIGN_PATTERNS.md` lines 984-1016

**Return a response:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 2)
2. Utils: `pkg/utils/response.go`
3. Details: `CODING_STANDARDS.md` lines 1479-1584

**Write tests:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 3)
2. Guide: `tests/README.md`
3. Details: `CODING_STANDARDS.md` lines 840-1009

**Handle errors:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 2)
2. Details: `CODING_STANDARDS.md` → Error Handling section

**Use dependency injection:**
1. Read: `00_AI_CRITICAL_RULES.md` (Tier 0, Rule 4)
2. Details: `DESIGN_PATTERNS.md` lines 439-492

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

### Mistake #4: Co-located Tests
```
❌ internal/app/services/auth_service_test.go  // Rejected
✅ tests/unit/services/auth_service_test.go
```

**Why:** See `00_AI_CRITICAL_RULES.md` Tier 0, Rule 3

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
- API versioning (`/api/v1`), global rate limit, single config source, request_id/LogStart/LogFinish logging, pluggable EmailSender (see "Recent changes" above).
- `AI_AGENT_RULES.md` is kept (important for AI agents).
- Added "New developer onboarding" path in this file.
- `OPENTELEMETRY_TRACING_ANALYSIS.md` — optional OpenTelemetry analysis; request_id is sufficient for monolith.

---

## ✅ Checklist Before First Code Contribution

```
□ Read 00_AI_CRITICAL_RULES.md (100 lines)
□ Read AI_QUICK_REFERENCE.md (405 lines)
□ Understand struct-based pattern requirement
□ Understand response utilities requirement
□ Understand tests/ directory requirement
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
