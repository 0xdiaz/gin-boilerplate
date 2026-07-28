# Configuration Management Guide

This document explains how to configure the application for different environments.

---

## 📋 Table of Contents

1. [Quick Start](#quick-start)
2. [Environment Detection (APP_ENV)](#environment-detection-app_env)
3. [Environment Variables](#environment-variables)
4. [Configuration Files](#configuration-files)
5. [Validation](#validation)
6. [Security Best Practices](#security-best-practices)
7. [Environment-Specific Setup](#environment-specific-setup)
8. [Environment-Based Features](#environment-based-features)
9. [Troubleshooting](#troubleshooting)

---

## 🚀 Quick Start

### 1. Copy Example Configuration

```bash
cp .env.example .env
```

### 2. Generate Secure Secrets

```bash
# Generate JWT_SECRET
openssl rand -base64 32
```

### 3. Update .env File

```env
# Replace with generated secret
JWT_SECRET=<your-generated-jwt-secret-here>

# Update database credentials
MASTER_DB_NAME=your_actual_db_name
MASTER_DB_USER=your_actual_db_user
MASTER_DB_PASSWORD=your_actual_db_password
MASTER_DB_HOST=localhost
```

### 4. Run the Application

```bash
go run main.go
```

The application will validate all required configuration on startup and **fail fast** if any required values are missing or using insecure defaults.

**How config is loaded:** Environment variables are read once at startup in `config.SetupConfig()` (from `.env` via Viper). Application code in `internal/` reads configuration only through `config.Get()`, which returns the loaded `Configuration` struct (Server, Database, etc.). Do not read env or Viper directly in app code; use `config.Get()` instead.

---

## 🌍 Environment Detection (APP_ENV)

The application uses `APP_ENV` to determine which environment it's running in. This controls feature flags, logging behavior, and other environment-specific configurations.

### Supported Environments

| Environment | APP_ENV Value | Use Case |
|-------------|---------------|----------|
| **Development** | `development` | Local development, debugging |
| **Staging** | `staging` | Pre-production testing |
| **Production** | `production` | Live production environment |

### Setting Environment

```env
# .env file
APP_ENV=development  # or staging, or production
```

If `APP_ENV` is not set, it defaults to `development`.

### Environment Helper Functions

Use these helper functions in your code to conditionally execute logic:

```go
import "github.com/0xdiaz/gin-boilerplate/pkg/config"

// Get current environment
env := config.GetEnvironment()  // Returns: "development", "staging", or "production"

// Check specific environment
if config.IsDevelopment() {
    // Only runs in development
    logger.Debugf("Detailed debug info here")
}

if config.IsProduction() {
    // Only runs in production
    initMetrics()
    enableStrictSecurity()
}

if config.IsStaging() {
    // Only runs in staging
    enableTestingFeatures()
}

// Check debug mode (auto-enabled in development)
if config.IsDebugEnabled() {
    // Debug mode active
}
```

### Real-World Examples

**Example 1: Enable SQL Logging Only in Development**
```go
// internal/adapters/database/database.go
func DbConnection(masterDSN, replicaDSN string) error {
    logMode := config.IsDevelopment() // true in dev, false in prod

    masterDB, err := gorm.Open(postgres.Open(masterDSN), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
        // Only log SQL in development
        DisableForeignKeyConstraintWhenMigrating: config.IsDevelopment(),
    })

    if config.IsDevelopment() {
        masterDB = masterDB.Debug() // Enable SQL logging
    }

    return nil
}
```

**Example 2: Seed Test Data Only in Development**
```go
// migrations/migrate.go
func Migrate() {
    // Run migrations
    database.MasterDB.AutoMigrate(&models.User{})

    // Seed test data only in development
    if config.IsDevelopment() {
        seedTestData()
    }
}

func seedTestData() {
    logger.Infof("Seeding test data for development...")
    // Create test users, etc.
}
```

**Example 3: Rate Limits From Environment**

Rate limit is read inside `RateLimitMiddleware()` from `RATE_LIMIT_RPS` and `RATE_LIMIT_BURST` (see `internal/app/middlewares/rate_limit.go`). If unset or ≤0, defaults (100 rps, 200 burst) are used. Set these in each environment's `.env` (e.g. lower in production, higher in development).

**Example 4: Enable Profiling in Non-Production**
```go
// cmd/main.go or middleware
if !config.IsProduction() {
    // Enable pprof profiling endpoints
    router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
}
```

**Example 5: Environment-Specific Error Messages**
```go
// pkg/utils/response.go
func InternalServerError(c *gin.Context, err error, message string) {
    if config.IsDevelopment() {
        // Show detailed error in development
        HandleErrors(c, http.StatusInternalServerError, err, fmt.Sprintf("%s: %v", message, err))
    } else {
        // Hide details in production
        HandleErrors(c, http.StatusInternalServerError, nil, message)
    }
}
```

### Environment Constants

Available constants for comparison:

```go
import "github.com/0xdiaz/gin-boilerplate/pkg/config"

config.EnvDevelopment  // "development"
config.EnvStaging      // "staging"
config.EnvProduction   // "production"

// Usage
if config.GetEnvironment() == config.EnvProduction {
    // Production-specific code
}
```

---

## 🔑 Environment Variables

### Required Variables

The following environment variables **MUST** be set. The application will not start if any are missing:

| Variable | Description | Example | Notes |
|----------|-------------|---------|-------|
| `JWT_SECRET` | JWT signing secret | Generated via openssl | Min 32 chars, MUST be changed from example |
| `SERVER_HOST` | Server bind address | `0.0.0.0` | Use `0.0.0.0` to bind all interfaces |
| `SERVER_PORT` | Server port | `8000` | Any available port |
| `MASTER_DB_NAME` | Master database name | `my_app_db` | Primary database for writes |
| `MASTER_DB_USER` | Master database user | `postgres` | User with write permissions |
| `MASTER_DB_PASSWORD` | Master database password | `secure_password` | Strong password required |
| `MASTER_DB_HOST` | Master database host | `localhost` | Database server address |
| `MASTER_DB_PORT` | Master database port | `5432` | PostgreSQL default: 5432 |

### Optional Variables

| Variable | Description | Default | Notes |
|----------|-------------|---------|-------|
| `APP_ENV` | Application environment | `development` | Values: `development`, `staging`, `production` |
| `DEBUG` | Debug mode | Auto (true in dev) | Set to `True` only in development |
| `ALLOWED_HOSTS` | Reserved for future Host-header validation | `0.0.0.0` | Comma-separated list. **Not** used for CORS or reverse proxy trust — see `CORS_ALLOWED_ORIGINS` and `TRUSTED_PROXIES` below. |
| `CORS_ALLOWED_ORIGINS` | Frontend origins allowed to make cross-origin requests | Localhost dev origins (see below) | Comma-separated, full URL with scheme (`https://app.example.com`), no wildcard. **Required in production** — startup fails if unset. |
| `TRUSTED_PROXIES` | Reverse proxy IPs/CIDRs trusted for `X-Forwarded-For` | `127.0.0.1,::1` in development | Comma-separated. Not related to CORS. Empty (trust none) by default outside development. |
| `SERVER_TIMEZONE` | Server timezone | `UTC` | Must be valid IANA timezone (e.g. UTC, Asia/Jakarta). Default applied in main.go when unset. |
| `RATE_LIMIT_RPS` | Auth rate limit (requests per second per IP) | `100` | Applied to `/auth` routes only. Set to 0 or omit to use default. |
| `RATE_LIMIT_BURST` | Auth rate limit burst size | `200` | Max tokens in bucket. Set to 0 or omit to use default. |
| `MASTER_DB_LOG_MODE` | Enable DB query logging | `True` | Set to `False` in production |
| `MASTER_SSL_MODE` | Database SSL mode | `disable` | Use `require` in production |
| `SERVER_SHUTDOWN_TIMEOUT` | Graceful shutdown timeout (seconds) | `10` | Max time to wait for in-flight requests before exit |
| `REFRESH_TOKEN_TTL_DAYS` | Refresh token lifetime (days) | `7` | Token is hashed at rest and rotated on every use; see [AUTHENTICATION.md](AUTHENTICATION.md) |

### CORS Configuration

`CORS_ALLOWED_ORIGINS`, `TRUSTED_PROXIES`, and `ALLOWED_HOSTS` are three separate concepts that are easy to confuse:

| Variable | Purpose | Consumed by |
|----------|---------|-------------|
| `CORS_ALLOWED_ORIGINS` | Which frontend origins may call this API from a browser | `middlewares.CORSMiddleware()` |
| `TRUSTED_PROXIES` | Which upstream proxy IPs are trusted to set `X-Forwarded-For` | `gin.Engine.SetTrustedProxies()` in `routers.SetupRoute()` |
| `ALLOWED_HOSTS` | Reserved for future `Host` header validation | Not currently enforced |

**Behavior:**

- `CORSMiddleware` never sets `Access-Control-Allow-Origin: *`. It echoes back the exact request `Origin` only when that origin is present in `CORS_ALLOWED_ORIGINS`, and sets `Vary: Origin` accordingly.
- `Access-Control-Allow-Credentials` is never set, since this API authenticates via a Bearer JWT in the `Authorization` header, not cookies.
- In development (`APP_ENV=development` or unset), unset `CORS_ALLOWED_ORIGINS`/`TRUSTED_PROXIES` fall back to safe localhost defaults.
- In staging/production, an empty `CORS_ALLOWED_ORIGINS` fails startup validation; an empty `TRUSTED_PROXIES` means no proxy is trusted (Gin uses the raw remote address).
- Entries must start with `http://` or `https://`; a bare `*` is rejected at startup.

**Verify with curl:**

```bash
# Allowed origin — should return the origin back in Access-Control-Allow-Origin
curl -i -X OPTIONS http://localhost:8000/api/v1/auth/login \
  -H "Origin: http://localhost:3000" \
  -H "Access-Control-Request-Method: POST" \
  -H "Access-Control-Request-Headers: Content-Type,Authorization"

# Disallowed origin — Access-Control-Allow-Origin should be absent
curl -i -X OPTIONS http://localhost:8000/api/v1/auth/login \
  -H "Origin: http://evil.com" \
  -H "Access-Control-Request-Method: POST"
```

### Graceful Shutdown

The application handles `SIGTERM` and `SIGINT` (e.g. Ctrl+C) for graceful shutdown: it stops accepting new requests, waits for in-flight requests to complete (up to `SERVER_SHUTDOWN_TIMEOUT` seconds), then closes the database connection and exits. Optional env `SERVER_SHUTDOWN_TIMEOUT` (integer, seconds) defaults to 10 if unset or zero.

### Replica Database (Optional)

For read scaling with database replicas:

| Variable | Description | Notes |
|----------|-------------|-------|
| `REPLICA_DB_NAME` | Replica database name | Can be same as master in development |
| `REPLICA_DB_USER` | Replica database user | Read-only user recommended |
| `REPLICA_DB_PASSWORD` | Replica database password | |
| `REPLICA_DB_HOST` | Replica database host | |
| `REPLICA_DB_PORT` | Replica database port | |
| `REPLICA_SSL_MODE` | Replica SSL mode | Use `require` in production |

---

## 📁 Configuration Files

### .env (Main Configuration)

- **Location:** Root directory
- **Purpose:** Main configuration file loaded on startup
- **Git:** `.gitignore`d - NEVER commit this file
- **Contains:** Actual secrets and environment-specific values

### .env.example (Template)

- **Location:** Root directory
- **Purpose:** Template showing all required variables
- **Git:** Committed to repository
- **Contains:** Placeholder values and documentation

### config/ Package

Configuration is managed through the `pkg/config` package:

```go
// config.go - Main config loader with validation
func SetupConfig() error

// db.go - Database configuration
func DbConfiguration() (masterDSN, replicaDSN string)

// server.go - Server configuration
func ServerConfig() string
```

---

## ✅ Validation

### Automatic Validation on Startup

The application automatically validates configuration when starting:

```go
// main.go
if err := config.SetupConfig(); err != nil {
    logger.Fatalf("config SetupConfig() error: %s", err)
}
```

### What Gets Validated

1. **Required Keys:** All required environment variables must be present
2. **Secret Security:** Secrets cannot use example/default values
3. **Secret Length:** Secrets must be at least 32 characters
4. **Value Presence:** No empty values for required fields

### Validation Error Examples

```bash
# Missing required key
Error: missing required config keys: JWT_SECRET, MASTER_DB_PASSWORD

# Using example value
Error: JWT_SECRET must be changed from example value and be at least 32 characters long.
Generate with: openssl rand -base64 32

# Empty value
Error: SERVER_PORT cannot be empty
```

---

## 🔒 Security Best Practices

### 1. Generate Strong Secrets

**❌ DON'T:**
```env
JWT_SECRET=mysecret
```

**✅ DO:**
```bash
# Generate cryptographically secure secrets
openssl rand -base64 32
# Output: K7gNU3sdo+OL0wNhqoVWhr3g6s1xYv72ol/pe/Unols=
```

### 2. Never Commit Secrets

**❌ DON'T:**
- Commit `.env` file to git
- Put real secrets in `.env.example`
- Hardcode secrets in source code

**✅ DO:**
- Keep `.env` in `.gitignore`
- Use placeholder values in `.env.example`
- Load secrets from environment variables

### 3. Use Different Secrets Per Environment

**❌ DON'T:**
```env
# Same secret everywhere
Development: JWT_SECRET=abc123
Production:  JWT_SECRET=abc123  # BAD!
```

**✅ DO:**
```env
# Different secrets per environment
Development: JWT_SECRET=dev_secret_K7gNU3sdo+OL0wNhqoVWhr3g
Production:  JWT_SECRET=prod_secret_X9mPQ7tzo+RL2xOisqXYks4h
```

### 4. Rotate Secrets Regularly

**Best Practice:**
- Rotate secrets every 90 days
- Rotate immediately if compromised
- Keep old secrets for brief transition period

### 5. Secure Secret Storage

**For Production:**
- **AWS:** Use AWS Secrets Manager or Parameter Store
- **GCP:** Use Google Secret Manager
- **Azure:** Use Azure Key Vault
- **Kubernetes:** Use Kubernetes Secrets
- **HashiCorp:** Use Vault

**Example with AWS Secrets Manager:**
```bash
# Store secret
aws secretsmanager create-secret \
    --name prod/myapp/jwt-secret \
    --secret-string "your-secret-here"

# Retrieve in application
# (requires AWS SDK integration)
```

---

## 🌍 Environment-Specific Setup

### Development Environment

```env
# .env.development (if implementing multi-env support)
DEBUG=True
SERVER_HOST=localhost
SERVER_PORT=8000

MASTER_DB_HOST=localhost
MASTER_DB_PORT=5432
MASTER_DB_NAME=myapp_dev
MASTER_DB_LOG_MODE=True
MASTER_SSL_MODE=disable

# Use same DB for replica in dev
REPLICA_DB_HOST=localhost
REPLICA_DB_NAME=myapp_dev
```

### Staging Environment

```env
# .env.staging
DEBUG=False
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

MASTER_DB_HOST=staging-db.internal.company.com
MASTER_DB_PORT=5432
MASTER_DB_NAME=myapp_staging
MASTER_DB_LOG_MODE=False
MASTER_SSL_MODE=require

REPLICA_DB_HOST=staging-db-replica.internal.company.com
```

### Production Environment

```env
# .env.production
DEBUG=False
SERVER_HOST=0.0.0.0
SERVER_PORT=8000

MASTER_DB_HOST=prod-db.internal.company.com
MASTER_DB_PORT=5432
MASTER_DB_NAME=myapp_production
MASTER_DB_LOG_MODE=False
MASTER_SSL_MODE=require

REPLICA_DB_HOST=prod-db-replica.internal.company.com
```

### Loading Environment-Specific Config (Future Enhancement)

```go
// Example: Load based on APP_ENV
env := os.Getenv("APP_ENV")
if env == "" {
    env = "development"
}

configFile := fmt.Sprintf(".env.%s", env)
viper.SetConfigFile(configFile)
```

---

## 🐛 Troubleshooting

### Application Won't Start

**Error:** `missing required config keys: JWT_SECRET`

**Solution:**
1. Ensure `.env` file exists
2. Check all required variables are present
3. Verify no empty values

```bash
# Check .env file
cat .env | grep JWT_SECRET
```

---

### "JWT_SECRET must be changed from example value"

**Error:** `JWT_SECRET must be changed from example value and be at least 32 characters long`

**Solution:**
```bash
# Generate new secret
openssl rand -base64 32

# Update .env
JWT_SECRET=<paste-generated-secret-here>
```

---

### Database Connection Failed

**Error:** `database DbConnection error: connection refused`

**Possible Causes:**
1. Wrong `MASTER_DB_HOST` or `MASTER_DB_PORT`
2. Database not running
3. Incorrect credentials
4. Firewall blocking connection

**Solution:**
```bash
# Test database connection
psql -h localhost -U your_user -d your_db

# Check if PostgreSQL is running
pg_isready -h localhost -p 5432
```

---

### "JWT secret not configured"

**Error:** `JWT secret not configured`

**Solution:**
Ensure `JWT_SECRET` is set in `.env`:
```env
JWT_SECRET=your-generated-jwt-secret-min-32-chars
```

---

## 8️⃣ Environment-Based Features

Use environment detection to conditionally enable/disable features and adjust behavior.

### Basic Usage

```go
import "github.com/0xdiaz/gin-boilerplate/pkg/config"

// Get current environment
env := config.GetEnvironment()
fmt.Println("Running in:", env) // development, staging, or production

// Check specific environments
if config.IsDevelopment() {
    fmt.Println("Development mode active")
}

if config.IsProduction() {
    fmt.Println("Production mode active")
}
```

### Real-World Examples

**1. Conditional Database Logging**

```go
// internal/adapters/database/database.go
func DbConnection(masterDSN, replicaDSN string) error {
    if config.IsDevelopment() {
        masterDB = masterDB.Debug() // Enable SQL logging only in dev
    }
    return nil
}
```

**2. Environment-Specific Seeding**

```go
// migrations/migrate.go
func Migrate() {
    database.MasterDB.AutoMigrate(&models.User{})

    if config.IsDevelopment() {
        seedTestData() // Only seed in development
    }
}
```

**3. Different Rate Limits Per Environment**

```go
// Rate limit is read from RATE_LIMIT_RPS / RATE_LIMIT_BURST in .env by RateLimitMiddleware().
// Set different values per environment (e.g. RATE_LIMIT_RPS=10 in production, 100 in development).
// Apply once on the apiV1 group (index.go) — sub-groups like authRoutes inherit it automatically.
// Calling RateLimitMiddleware() again on a sub-group reuses the same singleton limiter/key,
// which halves the effective burst for that sub-group instead of adding protection.
apiV1.Use(middlewares.RateLimitMiddleware())
```

**4. Development-Only Debug Endpoints**

```go
if config.IsDevelopment() {
    router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
}
```

**5. Environment-Specific Error Messages**

```go
func InternalServerError(c *gin.Context, err error, message string) {
    if config.IsDevelopment() {
        // Show details in dev
        HandleErrors(c, 500, err, fmt.Sprintf("%s: %v", message, err))
    } else {
        // Hide details in production
        HandleErrors(c, 500, nil, message)
    }
}
```

**6. Feature Flags**

```go
type FeatureFlags struct {
    EmailVerificationRequired bool
    TwoFactorAuthEnabled      bool
    BetaFeatures              bool
}

func GetFeatureFlags() FeatureFlags {
    if config.IsProduction() {
        return FeatureFlags{
            EmailVerificationRequired: true,
            TwoFactorAuthEnabled:      true,
            BetaFeatures:              false,
        }
    }
    // Relaxed for development
    return FeatureFlags{
        EmailVerificationRequired: false,
        TwoFactorAuthEnabled:      false,
        BetaFeatures:              true,
    }
}
```

**7. Environment-Specific Logging**

```go
func InitLogger() {
    if config.IsProduction() {
        logrus.SetLevel(logrus.WarnLevel) // Only warnings in prod
        logrus.SetFormatter(&logrus.JSONFormatter{})
    } else {
        logrus.SetLevel(logrus.DebugLevel) // Everything in dev
        logrus.SetFormatter(&logrus.TextFormatter{})
    }
}
```

**8. Conditional Middleware**

```go
// Only enable metrics in staging and production
if !config.IsDevelopment() {
    router.Use(middlewares.MetricsMiddleware())
}

// Security headers only in production
if config.IsProduction() {
    router.Use(middlewares.SecurityHeadersMiddleware())
}
```

### Security Best Practices

**Don't expose environment in responses:**

```go
// ❌ BAD
router.GET("/status", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "environment": config.GetEnvironment(), // ❌ Security risk!
    })
})

// ✅ GOOD
router.GET("/status", func(c *gin.Context) {
    c.JSON(200, gin.H{"status": "ok"})
})
```

**Environment-based security:**

```go
if config.IsProduction() {
    minPasswordLength = 12 // Strict in production
} else {
    minPasswordLength = 6  // Relaxed for testing
}
```

---

## 🔗 Related Documentation

- [README.md](../README.md) - Project overview and quick start
- [CODING_STANDARDS.md](CODING_STANDARDS.md) - Section 12: Configuration
- [Security Guidelines](CODING_STANDARDS.md#9-security-guidelines)

---

## 📝 Configuration Checklist

Before deploying to any environment:

```
□ Copied .env.example to .env
□ Generated unique JWT_SECRET (min 32 chars)
□ Updated all database credentials
□ Set DEBUG=False for production
□ Enabled SSL for database (MASTER_SSL_MODE=require)
□ Verified all required variables are set
□ Tested application starts without errors
□ Never committed .env file to git
□ Documented where production secrets are stored
```

---

**Last Updated:** 2025-11-09
**Maintainer:** Boilerplate Team
