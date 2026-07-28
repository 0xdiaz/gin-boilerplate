package auth

import (
	"context"

	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
)

// AuthServicer is the interface used by AuthController and AuthMiddleware.
// *AuthService implements this interface.
type AuthServicer interface {
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (string, error)
	ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error
	ValidateToken(token string) (uint, error)
	// GetProfile returns the current user's profile data. Returns ErrUserNotFound if userID doesn't exist.
	GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error)
	// Logout revokes the single refresh token supplied (this device/session only).
	Logout(ctx context.Context, req *dto.LogoutRequest) error
	// LogoutAll revokes every active refresh token for userID (all devices/sessions).
	LogoutAll(ctx context.Context, userID uint) error
}
