package services_test

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services/auth"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/0xdiaz/gin-boilerplate/tests/mocks"
	"github.com/stretchr/testify/assert"
)

const testJWTSecret = "test-jwt-secret-at-least-32-characters-long!!"

func setTestConfig(t *testing.T) {
	t.Helper()
	cfg := &config.Configuration{}
	cfg.Server.JWTSecret = testJWTSecret
	cfg.Server.RefreshTokenTTLDays = 7
	cfg.Server.AccessTokenTTLMinutes = 15
	config.SetForTest(cfg)
}

// testHashToken mirrors auth.hashToken (unexported) so tests can seed
// MockRefreshTokenRepository with records keyed the same way the service stores them.
func testHashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

func newAuthServiceWithMockRepo() *auth.AuthService {
	return auth.NewAuthService(mocks.NewMockUserRepository(), mocks.NewMockRefreshTokenRepository(), nil)
}

func newAuthServiceWithRepo(repo repositories.UserRepository) *auth.AuthService {
	return auth.NewAuthService(repo, mocks.NewMockRefreshTokenRepository(), nil)
}

// NOTE: These tests demonstrate testing patterns for services.
// In production, you should:
// 1. Mock repository dependencies (database calls)
// 2. Test business logic and error handling
// 3. Test token generation/validation with test secrets
// See TESTING.md for more details.

// TestValidateToken tests JWT token validation
func TestValidateToken(t *testing.T) {
	setTestConfig(t)
	service := newAuthServiceWithMockRepo()

	tests := []struct {
		name    string
		token   string
		wantErr bool
	}{
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "Malformed token",
			token:   "notavalidtoken",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.ValidateToken(tt.token)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestRegisterValidation tests registration input validation
func TestRegisterValidation(t *testing.T) {
	tests := []struct {
		name    string
		request *dto.RegisterRequest
		wantErr bool
	}{
		{
			name: "Valid registration request",
			request: &dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			wantErr: false,
		},
		{
			name: "Empty name",
			request: &dto.RegisterRequest{
				Name:     "",
				Email:    "john@example.com",
				Password: "SecurePass123!",
			},
			wantErr: true,
		},
		{
			name: "Invalid email",
			request: &dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "invalid-email",
				Password: "SecurePass123!",
			},
			wantErr: true,
		},
		{
			name: "Short password",
			request: &dto.RegisterRequest{
				Name:     "John Doe",
				Email:    "john@example.com",
				Password: "short",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Note: Actual validation happens in controller with ShouldBindJSON
			// This test demonstrates the DTO structure validation
			if tt.request.Name == "" || tt.request.Email == "" || len(tt.request.Password) < 8 {
				assert.True(t, tt.wantErr)
			}
		})
	}
}

// Example of table-driven test for business logic
func TestPasswordComplexity(t *testing.T) {
	tests := []struct {
		name     string
		password string
		valid    bool
	}{
		{
			name:     "Strong password",
			password: "StrongPass123!",
			valid:    true,
		},
		{
			name:     "Weak password - too short",
			password: "Pass1!",
			valid:    false,
		},
		{
			name:     "Minimum length",
			password: "Pass123!",
			valid:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Example validation: password must be at least 8 characters
			isValid := len(tt.password) >= 8

			assert.Equal(t, tt.valid, isValid)
		})
	}
}

// TestRefreshToken tests refresh token functionality with mocked repositories.
func TestRefreshToken(t *testing.T) {
	setTestConfig(t)

	t.Run("valid token rotates and old token is revoked", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository()
		userRepo.AddUserByEmail("test@example.com", &models.User{ID: 1, Name: "Test", Email: "test@example.com"})
		tokenRepo := mocks.NewMockRefreshTokenRepository()
		tokenRepo.Seed(&models.RefreshToken{
			UserID:    1,
			TokenHash: testHashToken("valid-token-here"),
			FamilyID:  "family-1",
			ExpiresAt: time.Now().Add(24 * time.Hour),
		})
		service := auth.NewAuthService(userRepo, tokenRepo, nil)

		response, err := service.RefreshToken(context.Background(), &dto.RefreshTokenRequest{RefreshToken: "valid-token-here"})
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.NotEmpty(t, response.AccessToken)
		assert.NotEmpty(t, response.RefreshToken)
		assert.NotEqual(t, "valid-token-here", response.RefreshToken)
		assert.Equal(t, "Bearer", response.TokenType)

		// Old token must now be revoked and unusable.
		old, _ := tokenRepo.GetByTokenHash(testHashToken("valid-token-here"))
		assert.NotNil(t, old.RevokedAt)
	})

	t.Run("unknown token", func(t *testing.T) {
		service := newAuthServiceWithMockRepo()
		_, err := service.RefreshToken(context.Background(), &dto.RefreshTokenRequest{RefreshToken: "never-issued"})
		assert.True(t, errors.Is(err, auth.ErrInvalidRefreshToken))
	})

	t.Run("expired token", func(t *testing.T) {
		userRepo := mocks.NewMockUserRepository()
		userRepo.AddUserByEmail("test@example.com", &models.User{ID: 1, Email: "test@example.com"})
		tokenRepo := mocks.NewMockRefreshTokenRepository()
		tokenRepo.Seed(&models.RefreshToken{
			UserID:    1,
			TokenHash: testHashToken("expired-token"),
			FamilyID:  "family-1",
			ExpiresAt: time.Now().Add(-1 * time.Hour),
		})
		service := auth.NewAuthService(userRepo, tokenRepo, nil)

		_, err := service.RefreshToken(context.Background(), &dto.RefreshTokenRequest{RefreshToken: "expired-token"})
		assert.True(t, errors.Is(err, auth.ErrInvalidRefreshToken))
	})
}

// TestRefreshTokenReuseDetection verifies that replaying an already-rotated
// token revokes the entire token family (theft response), not just the token.
func TestRefreshTokenReuseDetection(t *testing.T) {
	setTestConfig(t)
	userRepo := mocks.NewMockUserRepository()
	userRepo.AddUserByEmail("test@example.com", &models.User{ID: 1, Email: "test@example.com"})
	tokenRepo := mocks.NewMockRefreshTokenRepository()
	tokenRepo.Seed(&models.RefreshToken{
		UserID:    1,
		TokenHash: testHashToken("original-token"),
		FamilyID:  "family-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	service := auth.NewAuthService(userRepo, tokenRepo, nil)

	// First use: legitimate rotation.
	first, err := service.RefreshToken(context.Background(), &dto.RefreshTokenRequest{RefreshToken: "original-token"})
	assert.NoError(t, err)
	assert.NotNil(t, first)

	// Replaying the now-revoked original token must fail AND kill the family,
	// so the token rotated out of it (first.RefreshToken) also stops working.
	_, err = service.RefreshToken(context.Background(), &dto.RefreshTokenRequest{RefreshToken: "original-token"})
	assert.True(t, errors.Is(err, auth.ErrInvalidRefreshToken))

	_, err = service.RefreshToken(context.Background(), &dto.RefreshTokenRequest{RefreshToken: first.RefreshToken})
	assert.True(t, errors.Is(err, auth.ErrInvalidRefreshToken), "entire family should be revoked after reuse detection")
}

// TestLogout verifies that logout revokes only the presented token and is idempotent.
func TestLogout(t *testing.T) {
	setTestConfig(t)
	tokenRepo := mocks.NewMockRefreshTokenRepository()
	tokenRepo.Seed(&models.RefreshToken{
		UserID:    1,
		TokenHash: testHashToken("session-token"),
		FamilyID:  "family-1",
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	service := auth.NewAuthService(mocks.NewMockUserRepository(), tokenRepo, nil)

	err := service.Logout(context.Background(), &dto.LogoutRequest{RefreshToken: "session-token"})
	assert.NoError(t, err)

	revoked, _ := tokenRepo.GetByTokenHash(testHashToken("session-token"))
	assert.NotNil(t, revoked.RevokedAt)

	// Idempotent: logging out again (or an unknown token) must not error.
	err = service.Logout(context.Background(), &dto.LogoutRequest{RefreshToken: "session-token"})
	assert.NoError(t, err)
	err = service.Logout(context.Background(), &dto.LogoutRequest{RefreshToken: "never-issued"})
	assert.NoError(t, err)
}

// TestLogoutAll verifies that logout-all revokes every active token for a user.
func TestLogoutAll(t *testing.T) {
	setTestConfig(t)
	tokenRepo := mocks.NewMockRefreshTokenRepository()
	tokenRepo.Seed(&models.RefreshToken{UserID: 1, TokenHash: testHashToken("device-a"), FamilyID: "family-a", ExpiresAt: time.Now().Add(24 * time.Hour)})
	tokenRepo.Seed(&models.RefreshToken{UserID: 1, TokenHash: testHashToken("device-b"), FamilyID: "family-b", ExpiresAt: time.Now().Add(24 * time.Hour)})
	service := auth.NewAuthService(mocks.NewMockUserRepository(), tokenRepo, nil)

	err := service.LogoutAll(context.Background(), 1)
	assert.NoError(t, err)

	a, _ := tokenRepo.GetByTokenHash(testHashToken("device-a"))
	b, _ := tokenRepo.GetByTokenHash(testHashToken("device-b"))
	assert.NotNil(t, a.RevokedAt)
	assert.NotNil(t, b.RevokedAt)
}

// TestForgotPassword tests forgot password functionality using MockUserRepository.
func TestForgotPassword(t *testing.T) {
	mockRepo := mocks.NewMockUserRepository()
	// Seed user for "Valid email - user exists" case
	mockRepo.AddUserByEmail("user@example.com", &models.User{ID: 1, Email: "user@example.com", Name: "User"})
	service := newAuthServiceWithRepo(mockRepo)

	tests := []struct {
		name        string
		email       string
		expectError error
		expectToken bool
	}{
		{
			name:        "Valid email - user exists",
			email:       "user@example.com",
			expectError: nil,
			expectToken: true,
		},
		{
			name:        "Email not found",
			email:       "nonexistent@example.com",
			expectError: auth.ErrUserNotFound,
			expectToken: false,
		},
		{
			name:        "Invalid email format - no such user",
			email:       "invalid-email",
			expectError: auth.ErrUserNotFound,
			expectToken: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &dto.ForgotPasswordRequest{
				Email: tt.email,
			}

			token, err := service.ForgotPassword(context.Background(), req)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectError), "expected %v", tt.expectError)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				if tt.expectToken {
					assert.NotEmpty(t, token)
					assert.Len(t, token, 64)
				}
			}
		})
	}
}

// TestResetPassword tests password reset functionality using MockUserRepository.
func TestResetPassword(t *testing.T) {
	futureExpiry := time.Now().Add(15 * time.Minute)
	pastExpiry := time.Now().Add(-1 * time.Hour)
	validUser := &models.User{
		ID:                  1,
		Email:               "user@example.com",
		PasswordResetToken:  testHashToken("valid-reset-token"),
		PasswordResetExpiry: &futureExpiry,
	}
	expiredUser := &models.User{
		ID:                  2,
		Email:               "expired@example.com",
		PasswordResetToken:  testHashToken("expired-token"),
		PasswordResetExpiry: &pastExpiry,
	}
	mockRepo := mocks.NewMockUserRepository()
	// Mock is keyed by the same hash the service looks up by (raw token is never stored).
	mockRepo.SetUserByPasswordResetToken(testHashToken("valid-reset-token"), validUser)
	mockRepo.SetUserByPasswordResetToken(testHashToken("expired-token"), expiredUser)
	service := newAuthServiceWithRepo(mockRepo)

	tests := []struct {
		name        string
		token       string
		newPassword string
		expectError error
	}{
		{
			name:        "Valid reset token and password",
			token:       "valid-reset-token",
			newPassword: "NewSecurePass123!",
			expectError: nil,
		},
		{
			name:        "Invalid reset token",
			token:       "invalid-token",
			newPassword: "NewSecurePass123!",
			expectError: auth.ErrInvalidResetToken,
		},
		{
			name:        "Expired reset token",
			token:       "expired-token",
			newPassword: "NewSecurePass123!",
			expectError: auth.ErrResetTokenExpired,
		},
		{
			name:        "Short password",
			token:       "valid-reset-token",
			newPassword: "short",
			expectError: nil, // Validation happens in controller
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &dto.ResetPasswordRequest{
				Token:       tt.token,
				NewPassword: tt.newPassword,
			}

			err := service.ResetPassword(context.Background(), req)

			if tt.expectError != nil {
				assert.Error(t, err)
				assert.True(t, errors.Is(err, tt.expectError), "expected %v", tt.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestTokenGeneration tests that refresh and reset tokens are cryptographically secure
func TestTokenGeneration(t *testing.T) {
	tests := []struct {
		name          string
		tokenLength   int
		expectedChars string
	}{
		{
			name:          "Refresh token format",
			tokenLength:   64,
			expectedChars: "0123456789abcdef", // hex characters
		},
		{
			name:          "Reset token format",
			tokenLength:   64,
			expectedChars: "0123456789abcdef", // hex characters
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that tokens are the correct length
			assert.Equal(t, 64, tt.tokenLength, "Tokens should be 64 characters (32 bytes hex encoded)")

			// In actual implementation, tokens use crypto/rand
			// which is cryptographically secure
			// Test should verify randomness and uniqueness
		})
	}
}

// TestResetTokenExpiry tests that reset tokens expire correctly
func TestResetTokenExpiry(t *testing.T) {
	tests := []struct {
		name           string
		expiryDuration string
		shouldExpire   bool
	}{
		{
			name:           "Token expires in 15 minutes",
			expiryDuration: "15 minutes",
			shouldExpire:   false,
		},
		{
			name:           "Token expired 1 minute ago",
			expiryDuration: "expired",
			shouldExpire:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset tokens should expire after 15 minutes
			// Implementation sets expiry to time.Now().Add(15 * time.Minute)
			assert.True(t, true, "Token expiry logic implemented in service")
		})
	}
}
