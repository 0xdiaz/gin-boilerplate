package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
)

// RefreshToken generates new access and refresh tokens using a valid refresh token.
//
// The presented token is rotated: it is revoked and replaced by a newly
// issued token in the same family. If a token that was already rotated (or
// revoked) is presented again, this is treated as reuse/theft and the entire
// token family is revoked, forcing re-login on every device in that session.
//
// Returns ErrInvalidRefreshToken if the refresh token is invalid, revoked, or expired.
func (s *AuthService) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (resp *dto.RefreshTokenResponse, err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.RefreshToken")

	tokenRecord, err := s.refreshTokenRepo.GetByTokenHash(hashToken(req.RefreshToken))
	if err != nil {
		logger.Errorf("failed to get refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}

	if tokenRecord == nil {
		logger.Warnf("refresh token attempt with unknown token")
		logger.LogFinish(ctx, "AuthService.RefreshToken", ErrInvalidRefreshToken, start)
		return nil, ErrInvalidRefreshToken
	}

	if tokenRecord.RevokedAt != nil {
		logger.Warnf("refresh token reuse detected for user %d, revoking family %s", tokenRecord.UserID, tokenRecord.FamilyID)
		if revokeErr := s.refreshTokenRepo.RevokeFamily(tokenRecord.FamilyID); revokeErr != nil {
			logger.Errorf("failed to revoke token family after reuse detection: %v", revokeErr)
		}
		logger.LogFinish(ctx, "AuthService.RefreshToken", ErrInvalidRefreshToken, start)
		return nil, ErrInvalidRefreshToken
	}

	if time.Now().After(tokenRecord.ExpiresAt) {
		logger.Warnf("refresh token attempt with expired token for user %d", tokenRecord.UserID)
		logger.LogFinish(ctx, "AuthService.RefreshToken", ErrInvalidRefreshToken, start)
		return nil, ErrInvalidRefreshToken
	}

	user, err := s.userRepo.GetUserByID(tokenRecord.UserID)
	if err != nil {
		logger.Errorf("failed to get user for refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to validate refresh token: %w", err)
	}
	if user == nil {
		logger.LogFinish(ctx, "AuthService.RefreshToken", ErrInvalidRefreshToken, start)
		return nil, ErrInvalidRefreshToken
	}

	// Generate new access token
	accessToken, err := s.generateToken(user)
	if err != nil {
		logger.Errorf("failed to generate new access token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to generate new access token: %w", err)
	}

	// Rotate: issue a new refresh token in the same family, then revoke the old one
	newRefreshToken, err := s.issueRefreshToken(user.ID, tokenRecord.FamilyID)
	if err != nil {
		logger.Errorf("failed to issue rotated refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to issue rotated refresh token: %w", err)
	}
	if err := s.refreshTokenRepo.MarkRotated(tokenRecord.ID, hashToken(newRefreshToken)); err != nil {
		logger.Errorf("failed to revoke rotated refresh token: %v", err)
		logger.LogFinish(ctx, "AuthService.RefreshToken", err, start)
		return nil, fmt.Errorf("failed to revoke rotated refresh token: %w", err)
	}

	logger.Infof("refresh token successful for user: %s", user.Email)

	// Build response
	response := &dto.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		TokenType:    "Bearer",
	}

	logger.LogFinish(ctx, "AuthService.RefreshToken", nil, start)
	return response, nil
}

// Logout revokes the given refresh token (this device/session only).
//
// Always returns nil for an unknown or already-revoked token so callers
// cannot use this endpoint to probe whether a token exists.
func (s *AuthService) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	ctx, start := logger.LogStart(ctx, "AuthService.Logout")

	tokenRecord, err := s.refreshTokenRepo.GetByTokenHash(hashToken(req.RefreshToken))
	if err != nil {
		logger.Errorf("failed to get refresh token for logout: %v", err)
		logger.LogFinish(ctx, "AuthService.Logout", err, start)
		return fmt.Errorf("failed to process logout: %w", err)
	}
	if tokenRecord == nil || tokenRecord.RevokedAt != nil {
		logger.LogFinish(ctx, "AuthService.Logout", nil, start)
		return nil
	}

	if err := s.refreshTokenRepo.MarkRotated(tokenRecord.ID, ""); err != nil {
		logger.Errorf("failed to revoke refresh token on logout: %v", err)
		logger.LogFinish(ctx, "AuthService.Logout", err, start)
		return fmt.Errorf("failed to process logout: %w", err)
	}

	logger.Infof("logout successful for user: %d", tokenRecord.UserID)
	logger.LogFinish(ctx, "AuthService.Logout", nil, start)
	return nil
}

// LogoutAll revokes every active refresh token for userID (all devices/sessions).
func (s *AuthService) LogoutAll(ctx context.Context, userID uint) error {
	ctx, start := logger.LogStart(ctx, "AuthService.LogoutAll")

	if err := s.refreshTokenRepo.RevokeAllForUser(userID); err != nil {
		logger.Errorf("failed to revoke all refresh tokens: %v", err)
		logger.LogFinish(ctx, "AuthService.LogoutAll", err, start)
		return fmt.Errorf("failed to process logout-all: %w", err)
	}

	logger.Infof("logout-all successful for user: %d", userID)
	logger.LogFinish(ctx, "AuthService.LogoutAll", nil, start)
	return nil
}

// ForgotPassword generates a password reset token and returns it.
//
// In production, this token should be sent via email instead of returned in response.
// Returns ErrUserNotFound if email doesn't exist.
func (s *AuthService) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (token string, err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.ForgotPassword")

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		logger.Errorf("failed to get user by email: %v", err)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
		return "", fmt.Errorf("failed to process request: %w", err)
	}

	if user == nil {
		logger.Warnf("password reset attempt for non-existent email: %s", req.Email)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", ErrUserNotFound, start)
		return "", ErrUserNotFound
	}

	// Generate password reset token
	resetToken, err := s.generatePasswordResetToken()
	if err != nil {
		logger.Errorf("failed to generate reset token: %v", err)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Set token expiry (15 minutes from now). Only the hash is persisted;
	// the raw token is never stored, mirroring refresh token handling.
	expiry := time.Now().Add(15 * time.Minute)
	user.PasswordResetToken = hashToken(resetToken)
	user.PasswordResetExpiry = &expiry

	// Save to database
	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Errorf("failed to save reset token: %v", err)
		logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
		return "", fmt.Errorf("failed to save reset token: %w", err)
	}

	logger.Infof("password reset token generated for user: %s", user.Email)

	if s.mailer != nil {
		if err := s.mailer.SendPasswordResetEmail(user.Email, resetToken); err != nil {
			logger.Errorf("failed to send password reset email: %v", err)
			logger.LogFinish(ctx, "AuthService.ForgotPassword", err, start)
			return "", fmt.Errorf("failed to send reset email: %w", err)
		}
		logger.LogFinish(ctx, "AuthService.ForgotPassword", nil, start)
		return "", nil
	}
	logger.LogFinish(ctx, "AuthService.ForgotPassword", nil, start)
	return resetToken, nil
}

// ResetPassword resets user password using a valid reset token.
//
// Returns ErrInvalidResetToken if token is invalid.
// Returns ErrResetTokenExpired if token has expired.
func (s *AuthService) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) (err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.ResetPassword")

	user, err := s.userRepo.GetUserByPasswordResetToken(hashToken(req.Token))
	if err != nil {
		logger.Errorf("failed to get user by reset token: %v", err)
		logger.LogFinish(ctx, "AuthService.ResetPassword", err, start)
		return fmt.Errorf("failed to validate reset token: %w", err)
	}

	if user == nil {
		logger.Warnf("password reset attempt with invalid token")
		logger.LogFinish(ctx, "AuthService.ResetPassword", ErrInvalidResetToken, start)
		return ErrInvalidResetToken
	}

	// Check if token has expired
	if user.PasswordResetExpiry == nil || time.Now().After(*user.PasswordResetExpiry) {
		logger.Warnf("password reset attempt with expired token for user: %s", user.Email)
		logger.LogFinish(ctx, "AuthService.ResetPassword", ErrResetTokenExpired, start)
		return ErrResetTokenExpired
	}

	// Hash new password
	hashedPassword, err := s.hashPassword(req.NewPassword)
	if err != nil {
		logger.Errorf("failed to hash new password: %v", err)
		logger.LogFinish(ctx, "AuthService.ResetPassword", err, start)
		return fmt.Errorf("failed to process password: %w", err)
	}

	// Update password and clear reset token
	user.Password = hashedPassword
	user.PasswordResetToken = ""
	user.PasswordResetExpiry = nil

	// Save to database
	if err := s.userRepo.UpdateUser(user); err != nil {
		logger.Errorf("failed to update password: %v", err)
		logger.LogFinish(ctx, "AuthService.ResetPassword", err, start)
		return fmt.Errorf("failed to update password: %w", err)
	}

	logger.Infof("password reset successful for user: %s", user.Email)

	logger.LogFinish(ctx, "AuthService.ResetPassword", nil, start)
	return nil
}

// generatePasswordResetToken creates a cryptographically secure random reset token.
//
// Returns a 64-character hexadecimal string.
func (s *AuthService) generatePasswordResetToken() (string, error) {
	// Generate 32 random bytes (will be 64 chars in hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}
