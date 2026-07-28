# Tests Directory

This directory contains all test files following the project's testing standards.

## Running tests and coverage

From the project root:

- **Run unit tests:** `make test` or `go test ./tests/unit/... ./internal/... ./pkg/...`
- **Run with coverage:** `make test-coverage` or:
  ```bash
  go test -coverprofile=coverage.out ./tests/unit/... ./internal/... ./pkg/...
  go tool cover -func=coverage.out
  ```
- **Coverage gate (70%):** `make test-coverage-check` runs tests with coverage and fails if overall coverage is below 70%. The project standard is minimum 70% for the services layer (see docs/CODING_STANDARDS.md).
- **HTML coverage report:** `go tool cover -html=coverage.out` (opens in browser)

Mocks live under `tests/mocks/` (e.g. `MockUserRepository`, `MockAuthServicer`) for use in unit tests.

**Service coverage note:** `ExampleService` uses package-level repository calls (`repositories.Get`, `repositories.GetDataDatatables`) and is not unit-tested in isolation; it is covered by integration tests when the database is available.

## Directory Structure

```
tests/
├── unit/                  # Unit tests (isolated, mocked dependencies)
│   ├── controllers/       # Controller/handler tests
│   ├── services/          # Business logic tests
│   ├── repositories/      # Data access layer tests
│   └── utils/             # Utility function tests
├── integration/           # Integration tests (real DB, external services)
│   ├── api/              # End-to-end API tests
│   └── database/         # Database integration tests
└── fixtures/             # Test data (JSON, CSV, etc.)
```

## Test Organization Rules

### 1. Package Naming

All test files MUST use the `_test` suffix for package names:

```go
// ✅ CORRECT
package services_test  // External test (tests public API)
package controllers_test

// ❌ WRONG
package services  // Internal test (not allowed in tests/ folder)
```

### 2. File Naming

Test files MUST follow the pattern `{source_file}_test.go`:

```
auth_service.go       → auth_service_test.go
user_repository.go    → user_repository_test.go
validation_helper.go  → validation_helper_test.go
```

### 3. Test Types

#### Unit Tests (`tests/unit/`)
- Test individual functions/methods in isolation
- MUST mock all external dependencies (database, APIs, file system)
- Fast execution (< 100ms per test)
- No external service dependencies

Example:
```go
// tests/unit/services/auth_service_test.go
package services_test

import (
    "testing"
    "github.com/your-project/internal/app/services"
)

func TestAuthService_ValidateToken(t *testing.T) {
    service := services.NewAuthService()
    // Test with mocked dependencies
}
```

#### Integration Tests (`tests/integration/`)
- Test multiple components working together
- May use real database (test database)
- May call real external services (test environment)
- Slower execution

Example:
```go
// tests/integration/api/auth_api_test.go
package api_test

import (
    "testing"
    // Test with real database and HTTP server
)
```

#### Fixtures (`tests/fixtures/`)
- Reusable test data
- JSON, CSV, SQL files
- Shared across tests

Example:
```
fixtures/
├── users.json          # Sample user data
├── tokens.json         # Sample JWT tokens
└── transactions.sql    # SQL seed data
```

## Running Tests

```bash
# Run all tests
go test ./tests/...

# Run only unit tests
go test ./tests/unit/...

# Run only integration tests
go test ./tests/integration/...

# Run specific package tests
go test ./tests/unit/services/...

# Run with coverage
go test -cover ./tests/...

# Run with verbose output
go test -v ./tests/...

# Run integration tests (if using build tags)
go test -tags=integration ./tests/...
```

### Integration test notes

- **`tests/integration/api/health_test.go`** — Tests `GET /health` without requiring database. Expects 200 (healthy) or 503 (unhealthy). Run with: `go test ./tests/integration/api/...`
- **`tests/integration/database/connection_test.go`** — Tests database connectivity. Skipped unless `TEST_DB_MASTER_DSN` is set (e.g. `host=localhost user=postgres password=... dbname=test port=5432 sslmode=disable`). Run with: `go test ./tests/integration/database/...`

## Testing Standards

### Table-Driven Tests

MUST use table-driven tests for multiple scenarios:

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    bool
        wantErr bool
    }{
        {"valid input", "test@example.com", true, false},
        {"invalid input", "not-an-email", false, true},
        {"empty input", "", false, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Validate(tt.input)

            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }

            if got != tt.want {
                t.Errorf("got = %v, want %v", got, tt.want)
            }
        })
    }
}
```

### Test Coverage Goals

- **Services:** 70% minimum, 85% target
- **Repositories:** 70% minimum
- **Controllers:** 60% minimum (focus on validation and error handling)
- **Utils:** 80% minimum

### Mocking Guidelines

1. **Mock external dependencies**: Database, APIs, file system
2. **Use interfaces**: Makes mocking easier
3. **Keep mocks simple**: Don't over-complicate mock implementations
4. **Document mock behavior**: Add comments explaining what the mock does

Example:
```go
// Mock repository for testing
type MockUserRepo struct {
    GetUserFunc func(id uint) (*models.User, error)
}

func (m *MockUserRepo) GetUser(id uint) (*models.User, error) {
    return m.GetUserFunc(id)
}

// Usage in test
func TestGetUser(t *testing.T) {
    mockRepo := &MockUserRepo{
        GetUserFunc: func(id uint) (*models.User, error) {
            return &models.User{ID: id, Name: "Test"}, nil
        },
    }

    service := services.NewUserService(mockRepo)
    // Test with mocked repo
}
```

## Best Practices

1. ✅ Test files in `tests/` folder (not co-located with source)
2. ✅ Use `_test` package suffix
3. ✅ Mock all external dependencies for unit tests
4. ✅ Use table-driven tests
5. ✅ Test both happy path and error cases
6. ✅ Keep tests independent (no shared state)
7. ✅ Use meaningful test names
8. ✅ Add comments for complex test scenarios
9. ✅ Clean up resources in tests (use `t.Cleanup()`)
10. ✅ Keep fixtures in `tests/fixtures/`

## CI/CD Integration

Tests are automatically run in CI/CD pipeline. Ensure:

- All tests pass before merging
- Maintain minimum coverage requirements
- Integration tests use test database/environment
- No hardcoded secrets or credentials

## Resources

- [TESTING.md](../TESTING.md) - Comprehensive testing guide
- [CODING_STANDARDS.md](../docs/CODING_STANDARDS.md) - Testing standards section
- [Go Testing Package](https://pkg.go.dev/testing)
- [Testify Documentation](https://github.com/stretchr/testify)
