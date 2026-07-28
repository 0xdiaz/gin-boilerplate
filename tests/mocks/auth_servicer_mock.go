package mocks

import (
	"context"
	"errors"

	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/internal/app/services/auth"
)

// MockAuthServicer is a test double for auth.AuthServicer.
// Set fields to control behaviour; nil means return zero value and nil error.
var _ auth.AuthServicer = (*MockAuthServicer)(nil)

// MockAuthServicer implements auth.AuthServicer for controller tests.
type MockAuthServicer struct {
	RegisterFunc       func(context.Context, *dto.RegisterRequest) (*dto.AuthResponse, error)
	LoginFunc          func(context.Context, *dto.LoginRequest) (*dto.AuthResponse, error)
	RefreshTokenFunc   func(context.Context, *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error)
	ForgotPasswordFunc func(context.Context, *dto.ForgotPasswordRequest) (string, error)
	ResetPasswordFunc  func(context.Context, *dto.ResetPasswordRequest) error
	ValidateTokenFunc  func(string) (uint, error)
	LogoutFunc         func(context.Context, *dto.LogoutRequest) error
	LogoutAllFunc      func(context.Context, uint) error
	GetProfileFunc     func(context.Context, uint) (*dto.UserResponse, error)
}

func (m *MockAuthServicer) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthServicer) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	if m.LoginFunc != nil {
		return m.LoginFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthServicer) RefreshToken(ctx context.Context, req *dto.RefreshTokenRequest) (*dto.RefreshTokenResponse, error) {
	if m.RefreshTokenFunc != nil {
		return m.RefreshTokenFunc(ctx, req)
	}
	return nil, nil
}

func (m *MockAuthServicer) ForgotPassword(ctx context.Context, req *dto.ForgotPasswordRequest) (string, error) {
	if m.ForgotPasswordFunc != nil {
		return m.ForgotPasswordFunc(ctx, req)
	}
	return "", nil
}

func (m *MockAuthServicer) ResetPassword(ctx context.Context, req *dto.ResetPasswordRequest) error {
	if m.ResetPasswordFunc != nil {
		return m.ResetPasswordFunc(ctx, req)
	}
	return nil
}

func (m *MockAuthServicer) ValidateToken(token string) (uint, error) {
	if m.ValidateTokenFunc != nil {
		return m.ValidateTokenFunc(token)
	}
	return 0, errors.New("invalid token")
}

func (m *MockAuthServicer) Logout(ctx context.Context, req *dto.LogoutRequest) error {
	if m.LogoutFunc != nil {
		return m.LogoutFunc(ctx, req)
	}
	return nil
}

func (m *MockAuthServicer) LogoutAll(ctx context.Context, userID uint) error {
	if m.LogoutAllFunc != nil {
		return m.LogoutAllFunc(ctx, userID)
	}
	return nil
}

func (m *MockAuthServicer) GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	if m.GetProfileFunc != nil {
		return m.GetProfileFunc(ctx, userID)
	}
	return nil, nil
}
