# 📑 Documentation Index - Quick Lookup Guide

> **For AI Agents:** Use Ctrl+F / Cmd+F to search keywords. All line numbers are clickable references.

---

## 🚨 MUST READ FIRST (5 minutes)

| File | Lines | Time | Purpose |
|------|-------|------|---------|
| [00_AI_CRITICAL_RULES.md](00_AI_CRITICAL_RULES.md) | 100 | 2 min | Non-negotiable rules (Tier 0-2) |
| [AI_QUICK_REFERENCE.md](AI_QUICK_REFERENCE.md) | 405 | 3 min | Code templates for all layers |

---

## 🔍 Quick Keyword Lookup

### Architecture & Structure
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **Struct-based Controllers** | DESIGN_PATTERNS.md | 900-1016 | `struct`, `NewController`, `DI`, `dependency injection` |
| **Struct-based Services** | DESIGN_PATTERNS.md | 900-1016 | `struct`, `NewService`, `business logic` |
| **Repository Pattern** | DESIGN_PATTERNS.md | 1017-1085 | `repository`, `CRUD`, `database`, `function-based` |
| **Clean Architecture** | DESIGN_PATTERNS.md | 99-153 | `layers`, `dependencies`, `separation of concerns` |
| **Directory Structure** | DESIGN_PATTERNS.md | 517-682 | `folder`, `organization`, `internal/`, `pkg/` |

### Response & API
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **Standard Response Format** | CODING_STANDARDS.md | 1494-1546 | `success`, `message`, `data`, `errors`, `JSON` |
| **Response Utilities** | CODING_STANDARDS.md | 1547-1711 | `utils.Ok`, `utils.Created`, `utils.BadRequest` |
| **API Design** | CODING_STANDARDS.md | 1421-1726 | `RESTful`, `endpoints`, `HTTP methods`, `status codes` |
| **Error Handling** | DESIGN_PATTERNS.md | 1485-1613 | `error`, `logger.Errorf`, `fmt.Errorf`, `wrapping` |

### Routing & Middleware
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **Router Organization** | CODING_STANDARDS.md | 1727-1952 | `routes`, `Register*Routes`, `index.go`, `feature` |
| **Middleware** | DESIGN_PATTERNS.md | 1086-1173 | `gin.HandlerFunc`, `c.Next()`, `auth`, `rate limit` |
| **Rate Limiting** | CONFIGURATION.md | Optional vars, Example 3 | `RateLimitMiddleware`, `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST`, per-IP, token bucket; config read in middleware |

### Testing
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **Test Organization** | CODING_STANDARDS.md | 863-1036 | `tests/`, `_test package`, `unit`, `integration` |
| **Test Patterns** | DESIGN_PATTERNS.md | 1614-1780 | `table-driven`, `t.Run`, `setup`, `teardown` |
| **Mocking** | DESIGN_PATTERNS.md | 1614-1780 | `mock`, `interface`, `testify/mock` |

### Database
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **GORM Models** | CODING_STANDARDS.md | 1264-1420 | `gorm.Model`, `tableName`, `migrations`, `soft delete` |
| **Repositories** | DESIGN_PATTERNS.md | 1017-1085 | `repository`, `function-based`, `CreateUser`, `GetByID` |
| **Transactions** | DESIGN_PATTERNS.md | 1363-1484 | `Begin()`, `Commit()`, `Rollback()`, `atomic` |

### File Organization
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **DTO Placement** | CODING_STANDARDS.md | 132-157 | `dto`, `internal/app/dto`, `request`, `response` |
| **Enum/Constants** | CODING_STANDARDS.md | 159-187 | `enums`, `pkg/enums`, `const`, `iota` |
| **Import Cycles** | CODING_STANDARDS.md | 188-217 | `import cycle`, `pkg →`, `internal →` |

### Observability
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **Health Checks** | OBSERVABILITY.md | 22-130 | `/health`, `GET /health`, `database check`, `Kubernetes`, `Docker` |
| **Metrics** | OBSERVABILITY.md | 132-243 | `/metrics`, `request counters`, `uptime`, `error rate` |
| **Request Tracing** | OBSERVABILITY.md | 245-336 | `request ID`, `UUID`, `X-Request-ID`, `tracing`, `grep logs` |
| **Logger API (pkg/logger)** | OBSERVABILITY.md | Logger API section | `LogStart`, `LogFinish`, `FromContext`, `WithRequestID`, `spanName`, `request_id` |
| **Logging Standards** | CODING_STANDARDS.md | §8 | `logger`, `Infof`, `Errorf`, `LogStart`, `LogFinish`, `request-scoped` |
| **Performance Impact** | OBSERVABILITY.md | 387-435 | `overhead`, `benchmark`, `0.12%`, `negligible` |

### Configuration
| Topic | File | Lines | Keywords |
|-------|------|-------|----------|
| **Environment Detection** | CONFIGURATION.md | - | `APP_ENV`, `IsDevelopment()`, `IsProduction()`, `IsStaging()` |
| **Config Validation** | CONFIGURATION.md | - | `ValidateConfig()`, `startup`, `fail-fast`, `required keys` |
| **Secrets Management** | CONFIGURATION.md | - | `JWT_SECRET`, `openssl rand`, `32 characters` |
| **Environment Helpers** | CONFIGURATION.md | - | `config.IsDevelopment()`, `config.IsDebugEnabled()` |

---

## 📚 Full Document Structure

### CODING_STANDARDS.md (2,226 lines)

#### Critical Sections

| Section | Lines | What's Inside | Search Keywords |
|---------|-------|---------------|-----------------|
| **1. FILE ORGANIZATION** | 54-218 | File limits, directory structure, DTO/enum placement | `file size`, `300 lines`, `dto`, `enums`, `pkg/` |
| **2. NAMING CONVENTIONS** | 219-328 | Variable, function, file naming rules | `camelCase`, `PascalCase`, `snake_case`, `naming` |
| **3. CODE STRUCTURE** | 329-435 | Package structure, imports, grouping | `package`, `import`, `struct`, `interface` |
| **4. FUNCTION GUIDELINES** | 436-561 | Function size, parameters, return values | `function`, `100 lines`, `parameters`, `return` |
| **5. ERROR HANDLING** | 562-711 | Error wrapping, logging, recovery | `error`, `fmt.Errorf`, `%w`, `logger.Errorf` |
| **6. DOCUMENTATION** | 712-862 | Comments, godoc, package docs | `comment`, `//`, `godoc`, `documentation` |
| **7. TESTING** | 863-1036 | Test organization, coverage, patterns | `test`, `tests/`, `_test`, `coverage`, `70%` |
| **8. LOGGING** | 1037-1125 | Logger usage, levels, structured logging | `logger`, `Infof`, `Errorf`, `Warnf`, `Debugf` |
| **9. SECURITY** | 1126-1263 | Input validation, SQL injection, XSS | `security`, `validation`, `sanitize`, `bcrypt` |
| **10. DATABASE** | 1264-1420 | GORM models, migrations, queries | `gorm`, `model`, `migration`, `AutoMigrate` |
| **11. API DESIGN** | 1421-1953 | RESTful, responses, utilities, routing | `API`, `REST`, `response`, `utils.Ok`, `routes` |
| **12. CONFIGURATION** | 1955-2035 | Environment variables, config management | `config`, `env`, `.env`, `viper` |
| **13. FORBIDDEN PRACTICES** | 2036-2118 | What NOT to do | `panic`, `global`, `god object`, `forbidden` |

#### Line Number Quick Reference

```
L56-67    File Size Limits (MAX 300 lines, warning at 250)
L132-157  DTO Placement (internal/app/dto)
L159-187  Enum Placement (pkg/enums)
L188-217  Import Cycles (how to avoid)
L329-435  Code Structure & Imports
L562-711  Error Handling (MUST wrap errors)
L863-1036 Testing (tests/ directory, 70% coverage)
L1264-420 Database & GORM
L1429-475 Standard Response Format (success/message/data/errors)
L1479-584 Response Utilities (utils.Ok, utils.Created, etc.)
L1727-952 Router Organization (one file per feature)
```

---

### DESIGN_PATTERNS.md (2,497 lines)

#### Critical Sections

| Section | Lines | What's Inside | Search Keywords |
|---------|-------|---------------|-----------------|
| **1. OVERVIEW** | 53-98 | Why these patterns, goals | `overview`, `principles`, `goals` |
| **2. ARCHITECTURE** | 99-153 | Clean architecture, layers | `architecture`, `layers`, `dependencies` |
| **3. CORE PATTERNS** | 154-516 | Design patterns used | `singleton`, `factory`, `repository`, `DI` |
| **4. DIRECTORY STRUCTURE** | 517-682 | Full project structure | `directory`, `folder`, `tree`, `structure` |
| **5. LAYER RESPONSIBILITIES** | 683-918 | What each layer does | `controller`, `service`, `repository`, `model` |
| **6. IMPLEMENTATION PATTERNS** | 919-1173 | How to implement (MOST CRITICAL) | `struct`, `controller`, `service`, `repository` |
| **7. REQUEST FLOW** | 1174-1362 | How requests are processed | `flow`, `request`, `middleware`, `handler` |
| **8. DATA FLOW** | 1363-1484 | Data transformations | `DTO`, `model`, `mapping`, `transform` |
| **9. ERROR HANDLING** | 1485-1613 | Error patterns | `error`, `recovery`, `logging` |
| **10. TESTING** | 1614-1780 | Test patterns | `testing`, `mock`, `table-driven` |
| **11. FEATURE GUIDE** | 1781-2196 | Complete feature implementation | `feature`, `step-by-step`, `example` |
| **12. EXAMPLES** | 2197-2300 | Real code examples | `example`, `code`, `reference` |
| **13. ANTI-PATTERNS** | 2301-2460 | What to avoid | `wrong`, `bad`, `anti-pattern`, `avoid` |

#### Line Number Quick Reference

```
L99-153   Clean Architecture (layers, dependencies)
L154-516  Core Design Patterns (Repository, DI, Factory)
L517-682  Directory Structure (full tree)
L683-918  Layer Responsibilities (detailed)

🔥 MOST CRITICAL:
L900-1016 Controller & Service Pattern (STRUCT-BASED, MUST READ)
L1017-85  Repository Pattern (function-based CRUD)
L1086-173 Middleware Pattern

L1174-362 Request Flow (complete lifecycle)
L1363-484 Data Flow & Transformations
L1485-613 Error Handling Patterns
L1614-780 Testing Patterns (table-driven, mocking)
L1781-196 Complete Feature Implementation (step-by-step)
L2301-460 Anti-Patterns (what NOT to do)
```

---

## 🎯 Common Tasks - Quick Navigation

### "I need to create a new controller"
1. Read: `00_AI_CRITICAL_RULES.md` → Architecture Pattern (L10-27)
2. Template: `AI_QUICK_REFERENCE.md` → Controller Template (L69-107)
3. Details: `DESIGN_PATTERNS.md` → Controller Pattern (L900-960)
4. Response Utils: `CODING_STANDARDS.md` → Response Utilities (L1547-711)

### "I need to create a new service"
1. Read: `00_AI_CRITICAL_RULES.md` → Architecture Pattern (L10-27)
2. Template: `AI_QUICK_REFERENCE.md` → Service Template (L110-180)
3. Details: `DESIGN_PATTERNS.md` → Service Pattern (L961-1016)
4. Error Handling: `CODING_STANDARDS.md` → Error Handling (L562-711)

### "I need to add routes"
1. Read: `00_AI_CRITICAL_RULES.md` → Router Organization (L143-178)
2. Template: `AI_QUICK_REFERENCE.md` → Router Template (L286-330)
3. Details: `CODING_STANDARDS.md` → Router Organization (L1727-952)

### "I need to write tests"
1. Read: `00_AI_CRITICAL_RULES.md` → Test Location (L56-66)
2. Template: `AI_QUICK_REFERENCE.md` → Test Templates (L333-398)
3. Details: `CODING_STANDARDS.md` → Testing Requirements (L863-1036)
4. Patterns: `DESIGN_PATTERNS.md` → Testing Patterns (L1614-780)

### "I need to handle errors"
1. Read: `00_AI_CRITICAL_RULES.md` → Error Handling (L113-128)
2. Details: `CODING_STANDARDS.md` → Error Handling (L562-711)
3. Patterns: `DESIGN_PATTERNS.md` → Error Patterns (L1485-613)

### "I need to create database models"
1. Template: `AI_QUICK_REFERENCE.md` → Model Template (L225-257)
2. Details: `CODING_STANDARDS.md` → Database Access (L1264-420)
3. Repository: `DESIGN_PATTERNS.md` → Repository Pattern (L1017-85)

---

## 🔑 Critical Reminders

### TIER 0 - NEVER VIOLATE
1. ✅ **Struct-based** controllers & services (NOT standalone functions)
2. ✅ **Use response utilities** (NOT c.JSON directly)
3. ✅ **Tests in tests/ directory** (NOT co-located)
4. ✅ **Dependency injection** via New* constructors

### TIER 1 - HARD LIMITS
- File size: MAX 300 lines
- Function size: MAX 100 lines
- Test coverage: MIN 70% for services

### TIER 2 - CRITICAL PATTERNS
- All responses use `pkg/utils` functions
- All errors logged with `logger.Errorf`
- Router: one file per feature (`{feature}_routes.go`)

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
- Check: `00_AI_CRITICAL_RULES.md` → Decision Tree (L219-241)
- Reference: `AI_QUICK_REFERENCE.md` → Relevant template

### "Code review failed - what did I miss?"
- Check: `CODING_STANDARDS.md` → Code Review Checklist (L2119-171)
- Verify: `00_AI_CRITICAL_RULES.md` → Tier 0 Rules (L8-81)

### "Import cycle error"
- Fix: `CODING_STANDARDS.md` → Import Cycles (L188-217)
- Pattern: internal → pkg (allowed), pkg → internal (forbidden)

### "Which response utility to use?"
- List: `CODING_STANDARDS.md` → Response Utilities (L1547-711)
- Quick: `00_AI_CRITICAL_RULES.md` → Response Utilities (L96-111)

---

**Last Updated:** 2025-11-09
**Maintained By:** Boilerplate Team
