# Stable Contracts

This document defines **contracts** that this service treats as stable. Changing them is considered a **breaking change** and should bump the **major version** (e.g. v1 → v2). Projects and integrations that depend on these can upgrade on their own schedule.

---

## 1. API response shape

All JSON responses use the same top-level structure.

**Success:**

- `success` (boolean) — always `true`
- `message` (string) — human-readable message
- `data` (object or null) — response payload; structure depends on endpoint
- `errors` (object or null) — typically `null` on success

**Error:**

- `success` (boolean) — always `false`
- `message` (string) — error message
- `data` (object or null) — typically `null`
- `errors` (object or null) — optional details (e.g. validation errors)

**Implementation:** [pkg/types/response.go](../pkg/types/response.go), [pkg/utils/response.go](../pkg/utils/response.go) (Ok, BadRequest, Unauthorized, etc.).

**Stability:** Field names (`success`, `message`, `data`, `errors`) and meaning are stable. Adding optional top-level fields in a backward-compatible way may be done in a minor release.

---

## 2. Authentication (protected routes)

- **Header:** `Authorization: Bearer <access_token>`
- **Token:** JWT access token from login or register (or refresh).
- **Missing/invalid:** Respond with `401 Unauthorized` and the standard error response shape above.

**Implementation:** JWT guard in [internal/app/middlewares/auth.go](../internal/app/middlewares/auth.go), applied as `middlewares.AuthMiddleware(authService)` (validation via `auth.AuthServicer.ValidateToken`). The protected route `GET /api/v1/profile` uses it.

**Stability:** Header name and `Bearer ` prefix are stable. Changing to another scheme (e.g. API key in header) is a breaking change.

---

## 3. API path prefix

- **Versioned API:** All versioned endpoints live under `/api/v1/` (e.g. `/api/v1/auth/login`, `/api/v1/profile`).
- **Unversioned:** `/health` and `/metrics` stay at root (no `/api/v1` prefix).

**Stability:** Path prefix `/api/v1` is stable. Introducing `/api/v2` is additive; removing or repurposing `/api/v1` is breaking.

---

## 4. Environment / config keys

The following env keys are part of the contract. Renaming or removing them is breaking; adding new keys is not.

| Key | Purpose |
|-----|---------|
| `APP_ENV` | development / staging / production |
| `JWT_SECRET` | JWT signing secret (min 32 chars) |
| `DEBUG` | Enable debug mode (e.g. Swagger) |
| `SERVER_HOST`, `SERVER_PORT` | HTTP server address |
| `MASTER_DB_*` (e.g. `MASTER_DB_HOST`, `MASTER_DB_NAME`, …) | Primary DB connection |
| `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` | Global API rate limit |

**Implementation:** [pkg/config](../pkg/config), [.env.example](../.env.example).

**Stability:** Key names and their role are stable. New optional keys may be added in minor releases.

---

## 5. Cross-service contracts (service-to-service)

> This section is a **template** for documenting service-to-service contracts when this boilerplate is used as the basis for a service that integrates with other services. Replace the placeholder content below with your actual integration contracts.

### 5a. Inter-service authentication

All service-to-service calls use a shared internal API key:

- **Header:** `X-Internal-API-Key: <key>` (plus optionally `Authorization: Bearer <key>` accepted)
- **Scope header:** `X-Client-ID: <tenant-id>` for tenant-scoped requests
- **Target:** mTLS is the long-term target for mutual authentication between services

Interim auth details and key management belong in your deployment secrets store (not in this file).

### 5b. Idempotency

For operations that must be safe to retry (e.g. payment callbacks, webhook deliveries):

- Use an `Idempotency-Key` or `X-Idempotency-Key` header (or `idempotency_key` body field).
- The server deduplicates on `(client_id, idempotency_key)`.
- Replaying a request with the same key returns the prior result (not a new operation).
- A conflicting payload on a duplicate key returns `409 Conflict`.

### 5c. Outbound webhook delivery

When this service must notify a downstream service of an async event:

- Enqueue a delivery record in a `webhook_deliveries` table within the same DB transaction that produces the event.
- A relay worker drains due rows and delivers via `POST {target}/api/v1/internal/<event>`.
- Use full-jitter exponential backoff, a `max_attempts` cap, and a dead-letter alarm.
- Delivery is at-least-once; the downstream service deduplicates on the event's idempotency key.

### 5d. Versioning of cross-service contracts

Cross-service contracts follow the same semver policy as the public API (§6 below). Breaking a service-to-service contract is a **major** version bump; adding new optional fields is **minor**.

---

## 6. Versioning policy

- **Semver:** This service uses semantic versioning (e.g. v1.2.0).
- **Breaking change:** Any change that violates the contracts above (or that requires callers/clients to change their code or config) → **major** version bump.
- **New features / backward-compatible changes** → **minor** version bump.
- **Bug fixes / docs** → **patch** version bump.

When a major version is released, a short **migration guide** (what changed, how to upgrade) should be provided.

---

## Why this matters

When many projects or teams use this boilerplate:

- **Stable contracts** let the core evolve (refactors, dependencies, internal structure) without forcing every consumer to change at once.
- **Clear versioning** lets consumers decide when to upgrade and what to expect (breaking vs non-breaking).
- **Documented contracts** (this file) give a single place to check "can I change this?" before doing a breaking change.
