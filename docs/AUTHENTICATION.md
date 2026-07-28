# Authentication System Documentation

**Version:** 2.0
**Last Updated:** 2025-11-09
**Features:** JWT Authentication, Refresh Token, Password Reset

---

## Overview

Authentication is a self-contained module at `internal/modules/auth/`. It owns its full vertical slice (`handler → service → repository → model`) and exposes a small public surface to the rest of the app via `module.go` (see [MODULE_GUIDE.md](./MODULE_GUIDE.md)):

- `auth.New(db)` — build the module.
- `Module.RegisterRoutes(api)` — mount auth routes under `/api/v1`.
- `Module.Middleware()` — the JWT guard, handed to any module that needs to protect routes.
- `Module.Auth()` — the `auth.Servicer` contract (e.g. `ValidateToken`) for other modules.

This module provides a complete authentication system with the following features:

- ✅ User Registration
- ✅ User Login
- ✅ JWT Access Token (24 hours expiry)
- ✅ Refresh Token Mechanism
- ✅ Password Reset Flow
- ✅ Token Rotation (security best practice)

**Sensitive data:** Logs must not contain passwords or full tokens. The application masks or omits these; see [OBSERVABILITY.md](./OBSERVABILITY.md) for logging and masking details.

---

## Authentication Flow

### 1. Registration Flow

```
User → POST /api/v1/auth/register → auth.Handler → auth.Service → auth.Repository → Database
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
- Rate limiting (default 100 req/s per IP, applied to all `/api/v1`)
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
- Token rotation: Old refresh token invalidated immediately
- New refresh token generated for each refresh
- Reduces risk of token theft/replay attacks
- Refresh token stored in database (can be revoked)

**Token Lifecycle:**
- Access Token: 24 hours expiry (configurable)
- Refresh Token: No expiry, but rotated on each use

**Error responses:**
- `400 Bad Request` — Missing or invalid refresh_token in body.
- `401 Unauthorized` — Invalid or expired refresh token (e.g. not found or already rotated).

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
    "message": "If the email exists, a password reset link has been sent",
    "token": "c5d7e9f1a3b5c7d9e1f3a5b7c9d1e3f5..."  // Only when no mailer is configured (dev/testing)
  }
}
```

**Security Features:**
- Reset token: 64-char cryptographically secure hex string
- Token expiry: 15 minutes
- Generic success message (don't reveal if email exists)
- Rate limiting applied
- Token stored in database with expiry timestamp

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

In production, the reset token should be sent by email instead of returned in the API response. The module supports this via a pluggable **EmailSender** interface:

- **Interface:** `auth.EmailSender` (in `internal/modules/auth/mailer.go`) with a single method `SendPasswordResetEmail(to, resetToken string) error`. Implement this in your project (e.g. SMTP, SendGrid, SES).
- **Wire-up:** Pass your implementation into `auth.NewWithMailer(db, mailer)` (which calls `NewService(NewRepository(db), mailer)`). When `mailer` is non-nil, `ForgotPassword` sends the token via email and returns an empty string; the handler then omits the token from the response. When `mailer` is nil (default), the token is returned in the response for development and testing.
- **Default:** `auth.New(db)` calls `NewWithMailer(db, nil)`, so by default the token is returned in the response. To enable production-style behaviour, implement `EmailSender` and swap `auth.New(db)` for `auth.NewWithMailer(db, yourMailer)` in `buildModules()` (`internal/bootstrap/modules.go`).

---

## Implementation Details

### Database Schema

**User Model Fields** (`internal/modules/auth/model.go`):
```go
type User struct {
    ID       uint   `json:"id" gorm:"primaryKey"`
    Name     string `json:"name" gorm:"type:varchar(255);not null"`
    Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
    Password string `json:"-" gorm:"type:varchar(255);not null"` // Never expose in JSON

    // Refresh token for JWT token refresh mechanism
    RefreshToken string `json:"-" gorm:"type:varchar(500);index"`

    // Password reset token and expiry for forgot password flow
    PasswordResetToken  string     `json:"-" gorm:"type:varchar(255);index"`
    PasswordResetExpiry *time.Time `json:"-" gorm:"type:timestamp"`

    CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
    DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Soft delete support
}

// TableName specifies the database table name.
func (u *User) TableName() string { return "users" }
```

**Database Indexes:**
- `email`: Unique index for fast lookup and uniqueness
- `refresh_token`: Index for fast refresh token validation
- `password_reset_token`: Index for fast reset token validation

---

### Token Security

#### Access Token (JWT)
- **Algorithm:** HS256 (HMAC with SHA-256)
- **Claims:**
  - `user_id`: User's database ID
  - `email`: User's email
  - `exp`: Expiry timestamp (24 hours)
  - `iat`: Issued at timestamp
- **Secret:** Environment variable `JWT_SECRET` (min 32 chars)
- **Storage:** Client-side only (LocalStorage/Memory)

#### Refresh Token
- **Type:** Cryptographically secure random hex string
- **Length:** 64 characters (32 bytes)
- **Generation:** `crypto/rand` package
- **Storage:** Database (can be revoked)
- **Rotation:** New token generated on each refresh

#### Password Reset Token
- **Type:** Cryptographically secure random hex string
- **Length:** 64 characters (32 bytes)
- **Generation:** `crypto/rand` package
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
- Token rotation on refresh
- Refresh tokens stored in database (revocable)
- Access tokens with expiry

✅ **Rate Limiting:**
- Applied to all `/api/v1` routes (auth routes are mounted under `/api/v1`, so they inherit it)
- Prevents brute force attacks
- IP-based limiting (default 100 req/s, burst 200)

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

🔐 **Token Blacklisting:**
- Implement token blacklist for logout
- Use Redis for fast blacklist lookup
- Clear expired tokens periodically

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
| Invalid refresh token | 401 Unauthorized | "Invalid or expired refresh token" |
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

Per the modular layout, prefer co-locating tests with the module: `internal/modules/auth/*_test.go`, unit-testing `Service` against a fake `Repository` (no DB needed) — see `internal/modules/example/service_test.go` for the recipe. Legacy auth tests still live under `tests/unit/services/auth_service_test.go` and `tests/unit/controllers/auth_controller_test.go`.

**Test Coverage:**
- ✅ RefreshToken functionality
- ✅ ForgotPassword functionality
- ✅ ResetPassword functionality
- ✅ Token generation security
- ✅ Token expiry validation

**Running Tests:**
```bash
make test   # runs ./tests/unit/... ./internal/... ./pkg/...
```

---

## Configuration

### Environment Variables

Required variables in `.env`:

```bash
# JWT Configuration
JWT_SECRET=your-secret-key-min-32-characters  # Min 32 chars required

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

The `users` table is created/updated automatically on startup. `bootstrap.Run()` collects every module's models from `Module.Models()` (auth declares `&User{}`) and passes them to `migrations.Run(db, models)`:

```bash
# Using GORM AutoMigrate (development)
go run main.go  # Automatically migrates owned models on startup

# Using versioned SQL (production)
# Add versioned SQL under internal/migrations/sql/ (run via golang-migrate)
```

### Fields Added

New fields added to `users` table:
- `refresh_token` (varchar 500)
- `password_reset_token` (varchar 255)
- `password_reset_expiry` (timestamp)

---

## API Reference

### Summary

| Endpoint | Method | Auth Required | Description |
|----------|--------|---------------|-------------|
| `/api/v1/auth/register` | POST | No | Register new user |
| `/api/v1/auth/login` | POST | No | Authenticate user |
| `/api/v1/auth/refresh` | POST | No | Refresh access token |
| `/api/v1/auth/forgot-password` | POST | No | Request password reset |
| `/api/v1/auth/reset-password` | POST | No | Complete password reset |

### Rate Limiting

Rate limiting is applied globally to the `/api/v1` group in `bootstrap.buildEngine` (`internal/bootstrap/server.go`), so auth endpoints are covered. It is per client IP (token bucket) and reads limits from config inside the middleware (`pkg/middleware/rate_limit.go`):
- **Env vars:** `RATE_LIMIT_RPS`, `RATE_LIMIT_BURST` (see [CONFIGURATION.md](CONFIGURATION.md))
- **Defaults:** 100 requests per second, burst 200
- **Response when exceeded:** 429 Too Many Requests

---

## Changelog

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
- Module layout (source of truth): [docs/MODULE_GUIDE.md](MODULE_GUIDE.md)
- Stable API/config contracts: [docs/CONTRACTS.md](CONTRACTS.md)
- See main README: [README.md](../README.md)

---

**Built with ❤️ following enterprise-grade security practices**
