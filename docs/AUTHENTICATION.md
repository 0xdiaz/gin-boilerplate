# Authentication System Documentation

**Version:** 2.0
**Last Updated:** 2025-11-09
**Features:** JWT Authentication, Refresh Token, Password Reset

---

## Overview

This boilerplate provides a complete authentication system with the following features:

- ✅ User Registration
- ✅ User Login
- ✅ JWT Access Token (configurable expiry via `ACCESS_TOKEN_TTL_MINUTES`, default 15 minutes)
- ✅ Refresh Token Mechanism (hashed at rest, per-token expiry, rotation, reuse/theft detection)
- ✅ Logout (single session) and Logout-All (all devices)
- ✅ Password Reset Flow
- ✅ Token Rotation (security best practice)

**Sensitive data:** Logs must not contain passwords or full tokens. The application masks or omits these; see [OBSERVABILITY.md](./OBSERVABILITY.md) for logging and masking details.

---

## Authentication Flow

### 1. Registration Flow

```
User → POST /api/v1/auth/register → AuthController → AuthService → Repository → Database
                                    ↓
                    Generate Access Token & Refresh Token
                                    ↓
                              Return Response
```

**Endpoint:** `POST /api/v1/auth/register`

**Request:**
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "user": {
      "id": 1,
      "name": "John Doe",
      "email": "john@example.com"
    },
    "access_token": "eyJhbGciOiJIUzI1NiIs...",
    "refresh_token": "a3d5e8f9b2c1d4e6f7a8b9c0d1e2f3a4...",
    "token_type": "Bearer"
  }
}
```

**Error responses:**
- `400 Bad Request` — Validation failed (e.g. missing/invalid name, email, or password; password &lt; 8 chars).
- `409 Conflict` — Email already registered (e.g. `"email already exists"`).

**Security Features:**
- Password hashed with bcrypt (cost 10)
- Email uniqueness validation
- Input validation (name min 3 chars, password min 8 chars)

---

### 2. Login Flow

**Endpoint:** `POST /api/v1/auth/login`

**Request:**
```json
{
  "email": "john@example.com",
  "password": "SecurePass123!"
}
```

**Response:** Same as registration response (user, access_token, refresh_token, token_type).

**Error responses:**
- `400 Bad Request` — Validation failed (e.g. missing email or password, invalid JSON).
- `401 Unauthorized` — Invalid credentials (generic message; does not reveal whether email exists).

**Security Features:**
- Rate limiting (100 req/s per IP)
- Password verification with bcrypt
- Generic error messages (don't reveal if email exists)
- Refresh token rotation on each login

---

### 3. Refresh Token Flow

**Purpose:** Obtain new access token without re-authentication

**Endpoint:** `POST /api/v1/auth/refresh`

**Request:**
```json
{
  "refresh_token": "a3d5e8f9b2c1d4e6f7a8b9c0d1e2f3a4..."
}
```

**Response:**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIs...",  // New access token
    "refresh_token": "b4e6f8a0c2d4e6f8a0b2c4d6e8f0a2b4...",  // New refresh token
    "token_type": "Bearer"
  }
}
```

**Security Features:**
- **Hashed at rest:** Only the SHA-256 hash of the refresh token is stored; the raw value is never persisted.
- **Rotation:** Every refresh revokes the presented token and issues a new one in the same `family_id` chain.
- **Reuse (theft) detection:** If a token that was already rotated (or revoked) is presented again, the entire `family_id` is revoked — every device/session sharing that login is forced to re-authenticate. This is the signal that a refresh token was stolen and replayed.
- **Per-token expiry:** Configurable via `REFRESH_TOKEN_TTL_DAYS` (default 7 days); expired tokens are rejected the same as invalid ones.

**Token Lifecycle:**
- Access Token: `ACCESS_TOKEN_TTL_MINUTES` (default 15 minutes). Kept short because it is a stateless JWT and cannot be revoked — logout/logout-all only revoke refresh tokens, so this TTL bounds how long an already-issued access token keeps working after the user logs out.
- Refresh Token: `REFRESH_TOKEN_TTL_DAYS` (default 7 days), rotated on each use, single-use (old token immediately invalid)

**Error responses:**
- `400 Bad Request` — Missing or invalid refresh_token in body.
- `401 Unauthorized` — Invalid, expired, reused, or already-rotated refresh token (generic message; does not distinguish cause).

---

### 3b. Logout Flow (single session)

**Purpose:** Revoke one refresh token (this device only) without affecting other sessions.

**Endpoint:** `POST /api/v1/auth/logout`

**Request:**
```json
{
  "refresh_token": "a3d5e8f9b2c1d4e6f7a8b9c0d1e2f3a4..."
}
```

**Response:** `200 OK` with `"Logout successful"`, always — even if the token was already revoked or never existed, so this endpoint cannot be used to probe token validity.

---

### 3c. Logout-All Flow (all sessions)

**Purpose:** Revoke every refresh token for the authenticated user (e.g. "sign out everywhere" after a suspected compromise).

**Endpoint:** `POST /api/v1/logout-all` (requires `Authorization: Bearer <access_token>`)

**Response:** `200 OK` with `"Logged out from all devices successfully"`.

---

### 4. Forgot Password Flow

**Purpose:** Initiate password reset for forgotten passwords

**Endpoint:** `POST /api/v1/auth/forgot-password`

**Request:**
```json
{
  "email": "john@example.com"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset initiated",
  "data": {
    "message": "Password reset instructions sent to email",
    "token": "c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5..."  // Only in dev mode
  }
}
```

**Security Features:**
- Reset token: 64-char cryptographically secure hex string
- Token expiry: 15 minutes
- Generic success message (don't reveal if email exists)
- Rate limiting applied
- Only `SHA-256(raw)` is stored in `users.password_reset_token`; the raw token is never persisted, so a database leak alone does not yield a usable token

**Error responses:**
- `400 Bad Request` — Missing or invalid email.
- `404 Not Found` or generic success — In production, a generic success message is returned regardless of whether the email exists (don't reveal if email is registered).

**Production Consideration:**
- Token should be sent via email, not in response
- Include link to password reset page: `https://yourapp.com/reset-password?token={token}`
- Consider SMS verification for sensitive applications

---

### 5. Reset Password Flow

**Purpose:** Complete password reset using valid token

**Endpoint:** `POST /api/v1/auth/reset-password`

**Request:**
```json
{
  "token": "c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5...",
  "new_password": "NewSecurePass456!"
}
```

**Response:**
```json
{
  "success": true,
  "message": "Password reset successfully",
  "data": null
}
```

**Security Features:**
- Token validation (exists & not expired)
- Password hashed with bcrypt
- Token and expiry cleared after successful reset
- New password validation (min 8 chars)

**Error Responses:**
- Invalid token: `400 Bad Request - "Invalid reset token"`
- Expired token: `400 Bad Request - "Reset token has expired"`

---

### Password reset and email

In production, the reset token should be sent by email instead of returned in the API response. The boilerplate supports this via a pluggable **EmailSender** interface:

- **Interface:** `auth.EmailSender` with a single method `SendPasswordResetEmail(to, resetToken string) error`. Implement this in your project (e.g. SMTP, SendGrid, SES).
- **Wire-up:** Pass your implementation into `auth.NewAuthService(userRepo, refreshTokenRepo, mailer)`. When `mailer` is non-nil, `ForgotPassword` sends the token via email and does not return it in the response. When `mailer` is nil (default), the token is returned in the response for development and testing.
- **Default:** The boilerplate wires `NewAuthService(userRepo, refreshTokenRepo, nil)`, so by default the token is returned in the response. To enable production-style behaviour, implement `EmailSender` and inject it when constructing the auth service in your routes.

---

## Implementation Details

### Database Schema

**User Model Fields:**
```go
type User struct {
    ID                   uint       `gorm:"primaryKey"`
    Name                 string     `gorm:"type:varchar(255);not null"`
    Email                string     `gorm:"type:varchar(255);uniqueIndex;not null"`
    Password             string     `gorm:"type:varchar(255);not null"`

    // Password reset mechanism
    PasswordResetToken   string     `gorm:"type:varchar(255);index"`
    PasswordResetExpiry  *time.Time `gorm:"type:timestamp"`

    CreatedAt            time.Time
    UpdatedAt            time.Time
    DeletedAt            *time.Time `gorm:"index"`
}
```

**RefreshToken Model Fields** (dedicated table, not on `User`):
```go
type RefreshToken struct {
    ID              uint       `gorm:"primaryKey"`
    UserID          uint       `gorm:"not null;index"`
    TokenHash       string     `gorm:"type:varchar(64);uniqueIndex;not null"` // SHA-256 hex; raw token never stored
    FamilyID        string     `gorm:"type:varchar(64);not null;index"`       // rotation chain id, one per login session
    ReplacedByHash  string     `gorm:"type:varchar(64)"`                      // set on rotation, links to successor token
    RevokedAt       *time.Time
    ExpiresAt       time.Time  `gorm:"not null"`
    CreatedAt       time.Time
}
```

**Database Indexes:**
- `users.email`: Unique index for fast lookup and uniqueness
- `users.password_reset_token`: Index for fast reset token validation
- `refresh_tokens.token_hash`: Unique index, used to look up a token on refresh/logout
- `refresh_tokens.family_id`: Index, used to revoke a whole rotation chain on reuse detection or explicit revoke

---

### Token Security

#### Access Token (JWT)
- **Algorithm:** HS256 (HMAC with SHA-256)
- **Claims:**
  - `user_id`: User's database ID
  - `email`: User's email
  - `exp`: Expiry timestamp (`ACCESS_TOKEN_TTL_MINUTES`, default 15 minutes)
  - `iat`: Issued at timestamp
- **Secret:** Environment variable `JWT_SECRET` (min 32 chars)
- **Storage:** Client-side only (LocalStorage/Memory)

#### Refresh Token
- **Type:** Cryptographically secure random hex string
- **Length:** 64 characters (32 bytes)
- **Generation:** `crypto/rand` package
- **Storage:** Only `SHA-256(raw)` is stored in `refresh_tokens.token_hash`; the raw value returned to the client is never persisted, so a database leak alone does not yield usable tokens
- **Expiry:** `REFRESH_TOKEN_TTL_DAYS` (default 7 days) from issuance
- **Rotation:** New token generated on each refresh, in the same `family_id`; the old token is marked revoked with `replaced_by_hash` pointing at the new one
- **Reuse detection:** Presenting a revoked/already-rotated token revokes its entire `family_id`, invalidating every token in that login session

#### Password Reset Token
- **Type:** Cryptographically secure random hex string
- **Length:** 64 characters (32 bytes)
- **Generation:** `crypto/rand` package
- **Storage:** Only `SHA-256(raw)` is stored in `users.password_reset_token`; the raw value sent to the client (or emailed) is never persisted
- **Expiry:** 15 minutes from generation
- **Single Use:** Cleared after successful password reset

---

## Security Best Practices

### Implemented

✅ **Password Security:**
- Bcrypt hashing (cost 10)
- Minimum password length (8 chars)
- Password never exposed in JSON responses

✅ **Token Security:**
- Cryptographically secure token generation
- Refresh tokens hashed at rest (SHA-256); raw value never stored
- Token rotation on refresh, with reuse/theft detection (revokes whole session family)
- Explicit revocation: logout (single session) and logout-all (all devices)
- Per-token expiry (`REFRESH_TOKEN_TTL_DAYS`)
- Access tokens with expiry

✅ **Rate Limiting:**
- Applied to all auth endpoints
- Prevents brute force attacks
- IP-based limiting (100 req/s, burst 200)

✅ **Generic Error Messages:**
- Don't reveal if email exists (forgot password)
- Same error for invalid email/password
- Prevents user enumeration attacks

✅ **SQL Injection Prevention:**
- GORM parameterized queries
- Input validation with go-playground/validator

### Recommendations for Production

🔐 **Multi-Factor Authentication (MFA):**
- Add TOTP/SMS verification
- Require for sensitive operations

🔐 **Email Service Integration:**
- Send reset tokens via email (not in response)
- Use templates for professional emails
- Track email delivery status

🔐 **Account Security:**
- Login attempt tracking
- Account lockout after failed attempts
- Suspicious activity detection

🔐 **HTTPS Only:**
- Enforce HTTPS in production
- Use secure cookie flags
- HSTS headers

---

## Error Handling

### Common Errors

| Error | HTTP Status | Message |
|-------|-------------|---------|
| Email already exists | 409 Conflict | "Email already exists" |
| Invalid credentials | 401 Unauthorized | "Invalid email or password" |
| Invalid/expired/reused refresh token | 401 Unauthorized | "Invalid or expired refresh token" |
| Invalid reset token | 400 Bad Request | "Invalid reset token" |
| Expired reset token | 400 Bad Request | "Reset token has expired" |
| Validation error | 400 Bad Request | Specific validation message |

### Response Format

All errors follow the standard response format:

```json
{
  "success": false,
  "message": "Error message here",
  "data": null,
  "errors": [
    {
      "field": "email",
      "message": "Email is required"
    }
  ]
}
```

---

## Testing

### Unit Tests

Tests are located in `tests/unit/services/auth_service_test.go`

**Test Coverage:**
- ✅ RefreshToken functionality
- ✅ ForgotPassword functionality
- ✅ ResetPassword functionality
- ✅ Token generation security
- ✅ Token expiry validation

**Running Tests:**
```bash
go test ./tests/unit/services/... -v
```

**Note:** Most tests are skipped pending mock implementation. See test files for detailed test scenarios.

---

## Configuration

### Environment Variables

Required variables in `.env`:

```bash
# JWT Configuration
JWT_SECRET=your-secret-key-min-32-characters  # Min 32 chars required

# Refresh token lifetime in days (default 7)
REFRESH_TOKEN_TTL_DAYS=7

# Database
MASTER_DB_HOST=localhost
MASTER_DB_PORT=5432
MASTER_DB_NAME=your_database
MASTER_DB_USER=your_user
MASTER_DB_PASSWORD=your_password

# Server
SERVER_HOST=localhost
SERVER_PORT=8000
DEBUG=true
```

### Validation

Configuration is validated on startup:
- Secrets must be min 32 characters
- Cannot use example/default values
- All required variables must be present

---

## Migration

### Database Migration

Run migrations to add new fields:

```bash
go run main.go  # Applies pending SQL migrations automatically on startup (fail-fast)

# Add new schema changes as versioned SQL files instead:
make migrate-create NAME=add_something
```

### Schema Changes

- `users` table: `password_reset_token` (varchar 255), `password_reset_expiry` (timestamp). `refresh_token` was removed (moved to its own table, see below).
- `refresh_tokens` table (added in `000004_create_refresh_tokens_table`): `token_hash`, `family_id`, `replaced_by_hash`, `revoked_at`, `expires_at` — see [Database Schema](#database-schema) above.

---

## API Reference

### Summary

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/auth/register` | POST | No | Register new user |
| `/api/v1/auth/login` | POST | No | Authenticate user |
| `/api/v1/auth/refresh` | POST | No | Refresh access token (rotates refresh token) |
| `/api/v1/auth/logout` | POST | No (refresh token in body) | Revoke one refresh token (this session) |
| `/api/v1/logout-all` | POST | Yes | Revoke all refresh tokens for the user (all devices) |
| `/api/v1/auth/forgot-password` | POST | No | Request password reset |
| `/api/v1/auth/reset-password` | POST | No | Complete password reset |

### Rate Limiting

All auth endpoints use rate limiting (per IP, token bucket). Limits are read from config inside the middleware:
- **Env vars:** `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` (see [CONFIGURATION.md](CONFIGURATION.md))
- **Defaults:** 100 requests per second, burst 200
- **Response when exceeded:** 429 Too Many Requests

---

## Changelog

### Version 2.1 (2026-07-08)

**Added:**
- ✅ Refresh tokens moved to a dedicated `refresh_tokens` table, hashed (SHA-256) at rest
- ✅ Per-token expiry (`REFRESH_TOKEN_TTL_DAYS`)
- ✅ Rotation-family (`family_id`) reuse/theft detection — replaying a rotated token revokes the whole session
- ✅ `POST /api/v1/auth/logout` (single session) and `POST /api/v1/logout-all` (all devices)

**Fixed:**
- 🔒 Refresh tokens were previously stored in plaintext on `users.refresh_token` with no expiry; a database leak gave permanent session hijack. Now hashed, time-limited, and revocable.

### Version 2.0 (2025-11-09)

**Added:**
- ✅ Refresh token mechanism
- ✅ Token rotation on refresh
- ✅ Password reset flow (forgot/reset)
- ✅ Cryptographically secure tokens
- ✅ Token expiry management
- ✅ Repository methods for token operations
- ✅ Comprehensive unit tests
- ✅ Updated documentation

**Security Improvements:**
- Token rotation prevents replay attacks
- Time-limited reset tokens (15 min)
- Generic error messages prevent enumeration
- Refresh tokens stored in database (revocable)

### Version 1.0 (Initial)

**Features:**
- Basic JWT authentication
- User registration
- User login
- Protected routes
- Password hashing with bcrypt

---

## Support

For questions or issues:
- Check coding standards: [docs/CODING_STANDARDS.md](CODING_STANDARDS.md)
- Review design patterns: [docs/DESIGN_PATTERNS.md](DESIGN_PATTERNS.md)
- See main README: [README.md](../README.md)

---

**Built with ❤️ following enterprise-grade security practices**
