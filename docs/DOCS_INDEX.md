# 📑 Documentation Index - Quick Lookup Guide

> **For AI Agents:** Use Ctrl+F / Cmd+F to search keywords and the **§ section headings** below.
> The reference docs (CODING_STANDARDS, DESIGN_PATTERNS, OBSERVABILITY) are actively maintained, so
> any line numbers shown elsewhere are approximate — prefer the section heading / keyword.

---

## 🚨 MUST READ FIRST (5 minutes)

| File | Time | Purpose |
|------|------|---------|
| [MODULE_GUIDE.md](MODULE_GUIDE.md) | 3 min | **Project structure (source of truth).** Modular layout under `internal/modules/`; overrides old `internal/app`/`internal/domain` references in other docs. |
| [00_AI_CRITICAL_RULES.md](00_AI_CRITICAL_RULES.md) | 2 min | Non-negotiable rules (Tier 0-2) |
| [AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md) | 3 min | Code templates for all layers |

---

## 🔍 Quick Keyword Lookup

### Architecture & Structure
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **Struct-based Handlers** | DESIGN_PATTERNS.md | § Implementation Patterns | `struct`, `NewHandler`/`NewController`, `DI`, `dependency injection` |
| **Struct-based Services** | DESIGN_PATTERNS.md | § Implementation Patterns | `struct`, `NewService`, `business logic` |
| **Repository Pattern** | DESIGN_PATTERNS.md | § Implementation Patterns | `repository`, `CRUD`, `database`, `function-based` |
| **Clean Architecture** | DESIGN_PATTERNS.md | § Architecture | `layers`, `dependencies`, `separation of concerns` |
| **Directory Structure** | DESIGN_PATTERNS.md | § Directory Structure (modular layout: see MODULE_GUIDE) | `folder`, `organization`, `internal/modules/`, `pkg/` |

### Response & API
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **Standard Response Format** | CODING_STANDARDS.md | §11 API Design | `success`, `message`, `data`, `errors`, `JSON` |
| **Response Utilities** | CODING_STANDARDS.md | §11 API Design | `utils.Ok`, `utils.Created`, `utils.BadRequest` |
| **API Design** | CODING_STANDARDS.md | §11 API Design | `RESTful`, `endpoints`, `HTTP methods`, `status codes` |
| **Error Handling** | DESIGN_PATTERNS.md | § Error Handling | `error`, `logger.Errorf`, `fmt.Errorf`, `wrapping` |

### Routing & Middleware
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **Router Organization** | CODING_STANDARDS.md | §11 API Design | `routes`, `Register*Routes` (modular: `Module.RegisterRoutes`), `feature` |
| **Middleware** | DESIGN_PATTERNS.md | § Implementation Patterns | `gin.HandlerFunc`, `c.Next()`, `auth`, `rate limit` |
| **Rate Limiting** | CONFIGURATION.md | Optional vars, Example 3 | `RateLimitMiddleware`, `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`, per-IP, token bucket; config read in middleware |

### Testing
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **Test Organization** | CODING_STANDARDS.md | §7 Testing | `_test package`, `unit`, `integration` (modular: co-locate `*_test.go` per module — see MODULE_GUIDE) |
| **Test Patterns** | DESIGN_PATTERNS.md | § Testing | `table-driven`, `t.Run`, `setup`, `teardown` |
| **Mocking** | DESIGN_PATTERNS.md | § Testing | `mock`, `interface`, `testify/mock` |

### Database
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **GORM Models** | CODING_STANDARDS.md | §10 Database | `gorm.Model`, `tableName`, `migrations`, `soft delete` |
| **Repositories** | DESIGN_PATTERNS.md | § Implementation Patterns | `repository`, `function-based`, `CreateUser`, `GetByID` |
| **Transactions** | DESIGN_PATTERNS.md | § Data Flow | `Begin()`, `Commit()`, `Rollback()`, `atomic` |

### File Organization
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **DTO Placement** | CODING_STANDARDS.md | §1 File Organization | `dto`, per-module `dto.go`, `request`, `response` |
| **Constants** | CODING_STANDARDS.md | §1 File Organization | `constants`, module-local (no `pkg/enums`), `const`, `iota` |
| **Import Cycles** | CODING_STANDARDS.md | §1 File Organization | `import cycle`, `pkg →`, `internal →` |

### Observability
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **Health Checks** | OBSERVABILITY.md | § Health Checks | `/health`, `GET /health`, `database check`, `Kubernetes`, `Docker` |
| **Metrics** | OBSERVABILITY.md | § Metrics | `/metrics`, `request counters`, `uptime`, `error rate` |
| **Request Tracing** | OBSERVABILITY.md | § Request Tracing | `request ID`, `UUID`, `X-Request-ID`, `tracing`, `grep logs` |
| **Logger API (pkg/logger)** | OBSERVABILITY.md | § Logger API | `LogStart`, `LogFinish`, `FromContext`, `WithRequestID`, `spanName`, `request_id` |
| **Logging Standards** | CODING_STANDARDS.md | §8 Logging | `logger`, `Infof`, `Errorf`, `LogStart`, `LogFinish`, `request-scoped` |
| **Performance Impact** | OBSERVABILITY.md | § Performance Impact | `overhead`, `benchmark`, `negligible` |

### Configuration
| Topic | File | Section/Keyword | Keywords |
|-------|------|-------|----------|
| **Environment Detection** | CONFIGURATION.md | - | `APP_ENV`, `IsDevelopment()`, `IsProduction()`, `IsStaging()` |
| **Config Validation** | CONFIGURATION.md | - | `ValidateConfig()`, `startup`, `fail-fast`, `required keys` |
| **Secrets Management** | CONFIGURATION.md | - | `JWT_SECRET`, `openssl rand`, `32 characters` |
| **Environment Helpers** | CONFIGURATION.md | - | `config.IsDevelopment()`, `config.IsDebugEnabled()` |

### Security & Compliance 🔐
| Topic | File | Keywords |
|-------|------|----------|
| **Threat model / posture** | CODING_STANDARDS.md | `auth`, `trust boundary`, `JWT`, `CORS` |
| **Input validation** | CODING_STANDARDS.md | `binding`, `ShouldBindJSON`, `validated`, `OWASP injection` |
| **Secrets & config** | CONFIGURATION.md | `JWT_SECRET`, `HMAC`, `fail closed`, `env` |
| **Password security** | CODING_STANDARDS.md | `bcrypt`, `hashing`, `password` |
| **Rate limiting** | CONFIGURATION.md | `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`, `per-IP`, `token bucket` |

---

## 📚 Full Document Structure

### CODING_STANDARDS.md

> Use Ctrl+F on the section headings below — section line numbers are intentionally omitted because
> this doc is actively maintained and line numbers drift. Search the **§ heading** instead.

#### Critical Sections

| Section | What's Inside | Search Keywords |
|---------|---------------|-----------------|
| **§1. FILE ORGANIZATION** | File limits, directory structure, DTO/enum placement | `file size`, `300 lines`, `dto`, `enums`, `pkg/` |
| **§2. NAMING CONVENTIONS** | Variable, function, file naming rules | `camelCase`, `PascalCase`, `snake_case`, `naming` |
| **§3. CODE STRUCTURE** | Package structure, imports, grouping | `package`, `import`, `struct`, `interface` |
| **§4. FUNCTION GUIDELINES** | Function size, parameters, return values | `function`, `100 lines`, `parameters`, `return` |
| **§5. ERROR HANDLING** | Error wrapping, logging, recovery | `error`, `fmt.Errorf`, `%w`, `logger.Errorf` |
| **§6. DOCUMENTATION** | Comments, godoc, package docs | `comment`, `//`, `godoc`, `documentation` |
| **§7. TESTING** | Test organization, coverage, patterns | `test`, `_test`, `coverage`, `70%` |
| **§8. LOGGING** | Logger usage, levels, structured logging | `logger`, `Infof`, `Errorf`, `Warnf`, `Debugf` |
| **§9. SECURITY** | Input validation, SQL injection, XSS | `security`, `validation`, `sanitize`, `bcrypt` |
| **§10. DATABASE** | GORM models, migrations, queries | `gorm`, `model`, `migration`, `AutoMigrate` |
| **§11. API DESIGN** | RESTful, responses, utilities, routing | `API`, `REST`, `response`, `utils.Ok`, `routes` |
| **§12. CONFIGURATION** | Environment variables, config management | `config`, `env`, `.env`, `viper` |
| **§13. FORBIDDEN PRACTICES** | What NOT to do | `panic`, `global`, `god object`, `forbidden` |

#### Topic Quick Reference (search the keyword)

```
File Size Limits         §1  (MAX 300 lines, warning at 250)
DTO Placement            §1  (per-module dto.go)
Constants                §1  (module-local; no pkg/enums)
Import Cycles            §1  (how to avoid)
Code Structure & Imports §3
Error Handling           §5  (MUST wrap errors)
Testing                  §7  (70% coverage; modular: co-locate *_test.go — see MODULE_GUIDE)
Database & GORM          §10
Standard Response Format §11 (success/message/data/errors)
Response Utilities       §11 (utils.Ok, utils.Created, etc.)
Router Organization      §11 (modular: each module's RegisterRoutes; doc shows one-file-per-feature)
```

---

### DESIGN_PATTERNS.md

> Section line numbers are intentionally omitted (this doc is actively maintained and line numbers
> drift). Search the **§ heading**. Note: DESIGN_PATTERNS uses the older layered terminology
> (controllers/services in `internal/app`); for the current modular layout see MODULE_GUIDE.md —
> the HTTP layer is `handler.go` per module.

#### Critical Sections

| Section | What's Inside | Search Keywords |
|---------|---------------|-----------------|
| **§1. OVERVIEW** | Why these patterns, goals | `overview`, `principles`, `goals` |
| **§2. ARCHITECTURE** | Clean architecture, layers | `architecture`, `layers`, `dependencies` |
| **§3. CORE PATTERNS** | Design patterns used | `singleton`, `factory`, `repository`, `DI` |
| **§4. DIRECTORY STRUCTURE** | Project structure (modular: see MODULE_GUIDE) | `directory`, `folder`, `tree`, `structure` |
| **§5. LAYER RESPONSIBILITIES** | What each layer does | `handler`/`controller`, `service`, `repository`, `model` |
| **§6. IMPLEMENTATION PATTERNS** | How to implement (MOST CRITICAL) | `struct`, `handler`/`controller`, `service`, `repository` |
| **§7. REQUEST FLOW** | How requests are processed | `flow`, `request`, `middleware`, `handler` |
| **§8. DATA FLOW** | Data transformations | `DTO`, `model`, `mapping`, `transform` |
| **§9. ERROR HANDLING** | Error patterns | `error`, `recovery`, `logging` |
| **§10. TESTING** | Test patterns | `testing`, `mock`, `table-driven` |
| **§11. FEATURE GUIDE** | Complete feature implementation | `feature`, `step-by-step`, `example` |
| **§12. EXAMPLES** | Real code examples | `example`, `code`, `reference` |
| **§13. ANTI-PATTERNS** | What to avoid | `wrong`, `bad`, `anti-pattern`, `avoid` |

#### Topic Quick Reference (search the keyword)

```
Clean Architecture        §2  (layers, dependencies)
Core Design Patterns      §3  (Repository, DI, Factory)
Directory Structure       §4  (modular: see MODULE_GUIDE)
Layer Responsibilities    §5

🔥 MOST CRITICAL:
Handler/Controller & Service Pattern  §6  (STRUCT-BASED, MUST READ)
Repository Pattern        §6  (function-based CRUD)
Middleware Pattern        §6

Request Flow              §7  (complete lifecycle)
Data Flow & Transforms    §8
Error Handling Patterns   §9
Testing Patterns          §10 (table-driven, mocking)
Complete Feature Impl.    §11 (step-by-step)
Anti-Patterns             §13 (what NOT to do)
```

---

## 🎯 Common Tasks - Quick Navigation

> Line numbers below are kept only for files that did not move (`AI_QUICK_REFERENCE.md`,
> `00_AI_CRITICAL_RULES.md`). For `DESIGN_PATTERNS.md` / `CODING_STANDARDS.md` use the **§ section**
> (their line numbers drift). For the modular layout, start from MODULE_GUIDE.md.

### "I need to add a new module" (the modular way — start here)
1. Read: `MODULE_GUIDE.md` → "How to add a new module"
2. Copy: `internal/modules/example/` (the reference module)
3. Register: one line in `buildModules()` in `internal/bootstrap/modules.go`

### "I need to create a new handler" (HTTP layer of a module)
1. Read: `00_AI_CRITICAL_RULES.md` → Architecture Pattern
2. Reference: `internal/modules/example/handler.go`
3. Patterns: `DESIGN_PATTERNS.md` → § Implementation Patterns (uses old controller terminology)
4. Response Utils: `CODING_STANDARDS.md` → § 11 API Design

### "I need to create a new service" (business logic of a module)
1. Read: `00_AI_CRITICAL_RULES.md` → Architecture Pattern
2. Reference: `internal/modules/example/service.go`
3. Patterns: `DESIGN_PATTERNS.md` → § Implementation Patterns
4. Error Handling: `CODING_STANDARDS.md` → § 5 Error Handling

### "I need to add routes"
1. Mount them in the module's `Module.RegisterRoutes(api)` (modular layout) — see `internal/modules/auth/module.go`
2. Read: `00_AI_CRITICAL_RULES.md` → Router Organization
3. Details: `CODING_STANDARDS.md` → § 11 API Design (describes old one-file-per-feature routers)

### "I need to write tests"
1. Read: `MODULE_GUIDE.md` → Testing (co-locate `*_test.go`; fake repo, no DB)
2. Reference: `internal/modules/example/service_test.go`
3. Template: `AI_QUICK_REFERENCE.md` → Test Templates
4. Patterns: `DESIGN_PATTERNS.md` → § Testing

### "I need to handle errors"
1. Read: `00_AI_CRITICAL_RULES.md` → Error Handling
2. Details: `CODING_STANDARDS.md` → § 5 Error Handling
3. Patterns: `DESIGN_PATTERNS.md` → § Error Handling

### "I need to create database models"
1. Reference: `internal/modules/example/model.go`; list models in `Module.Models()`
2. Template: `AI_QUICK_REFERENCE.md` → Model Template
3. Details: `CODING_STANDARDS.md` → § 10 Database
4. Repository: `DESIGN_PATTERNS.md` → § Implementation Patterns (Repository)

---

## 🔑 Critical Reminders

### TIER 0 - NEVER VIOLATE
1. ✅ **Struct-based** handlers & services (NOT standalone functions)
2. ✅ **Use response utilities** (NOT c.JSON directly)
3. ✅ **Co-locate tests with the module** (`internal/modules/<name>/*_test.go`); `tests/` holds shared mocks + legacy/integration tests. (Older docs say "tests MUST live in `tests/`" — that is the pre-refactor rule; MODULE_GUIDE overrides it.)
4. ✅ **Dependency injection** via New* constructors

### TIER 1 - HARD LIMITS
- File size: MAX 300 lines
- Function size: MAX 100 lines
- Test coverage: MIN 70% for services

### TIER 2 - CRITICAL PATTERNS
- All responses use `pkg/utils` functions
- All errors logged with `logger.Errorf`
- Routing: each module mounts its own routes in `Module.RegisterRoutes(api)` (modular layout). Older docs describe one router file per feature (`{feature}_routes.go`) from the pre-refactor layout.

---

## 📖 Reading Order for New AI Agents

**Total Time: ~15 minutes for critical path**

1. **START HERE** (2 min): [00_AI_CRITICAL_RULES.md](00_AI_CRITICAL_RULES.md)
   - Absolute rules, decision trees, quick patterns

2. **TEMPLATES** (3 min): [AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md)
   - Copy-paste templates for all layers

3. **THIS INDEX** (2 min): [DOCS_INDEX.md](DOCS_INDEX.md)
   - Bookmark for quick lookups

4. **ON-DEMAND REFERENCE** (as needed):
   - [CODING_STANDARDS.md](CODING_STANDARDS.md) - Use Ctrl+F for specific topics
   - [DESIGN_PATTERNS.md](DESIGN_PATTERNS.md) - Deep dive when needed

---

## 🆘 Troubleshooting

### "I violated a pattern - where's the correct way?"
- Check: `00_AI_CRITICAL_RULES.md` → Decision Tree
- Reference: `AI_QUICK_REFERENCE.md` → Relevant template

### "Code review failed - what did I miss?"
- Check: `CODING_STANDARDS.md` → Code Review Checklist (search the heading)
- Verify: `00_AI_CRITICAL_RULES.md` → Tier 0 Rules

### "Import cycle error"
- Fix: `CODING_STANDARDS.md` → § 1 File Organization (Import Cycles)
- Pattern: internal → pkg (allowed), pkg → internal (forbidden)

### "Which response utility to use?"
- List: `CODING_STANDARDS.md` → § 11 API Design (Response Utilities)
- Quick: `00_AI_CRITICAL_RULES.md` → Response Utilities

---

**Last Updated:** 2026-06-10 (aligned to modular layout — see MODULE_GUIDE.md)
**Maintained By:** Boilerplate Team
