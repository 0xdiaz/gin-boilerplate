package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/app/dto"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
	"github.com/0xdiaz/gin-boilerplate/pkg/config"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Common errors for auth service
var (
	ErrEmailAlreadyExists  = errors.New("email already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrUserNotFound        = errors.New("user not found")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrInvalidResetToken   = errors.New("invalid or expired reset token")
	ErrResetTokenExpired   = errors.New("reset token has expired")
)

// AuthService handles authentication-related business logic
type AuthService struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	mailer           EmailSender // optional; when set, ForgotPassword sends token via email instead of returning it
}

// NewAuthService creates a new AuthService instance. mailer may be nil (token returned in response for dev/testing).
func NewAuthService(userRepo repositories.UserRepository, refreshTokenRepo repositories.RefreshTokenRepository, mailer EmailSender) *AuthService {
	return &AuthService{userRepo: userRepo, refreshTokenRepo: refreshTokenRepo, mailer: mailer}
}

// Register creates a new user account with validation and password hashing.
//
// Returns ErrEmailAlreadyExists if email is already registered.
// Password is hashed using bcrypt before storage.
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (resp *dto.AuthResponse, err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.Register")

	existingUser, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		logger.LogFinish(ctx, "AuthService.Register", err, start)
		return nil, fmt.Errorf("failed to check email: %w", err)
	}

	if existingUser != nil {
		logger.Warnf("registration attempt with existing email: %s", req.Email)
		logger.LogFinish(ctx, "AuthService.Register", ErrEmailAlreadyExists, start)
		return nil, ErrEmailAlreadyExists
	}

	// Hash password
	hashedPassword, err := s.hashPassword(req.Password)
	if err != nil {
		logger.Errorf("failed to hash password: %v", err)
		logger.LogFinish(ctx, "AuthService.Register", err, start)
		return nil, fmt.Errorf("failed to process password: %w", err)
	}

	// Create user model
	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: hashedPassword,
	}

	// Save to database
	if err = s.userRepo.CreateUser(user); err != nil {
		logger.Errorf("failed to create user: %v", err)
		logger.LogFinish(ctx, "AuthService.Register", err, start)
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	logger.Infof("user registered successfully: %s", user.Email)

	response, err := s.issueAuthResponse(user)
	logger.LogFinish(ctx, "AuthService.Register", err, start)
	return response, err
}

// Login authenticates a user with email and password.
//
// Returns ErrInvalidCredentials if email or password is incorrect.
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (resp *dto.AuthResponse, err error) {
	ctx, start := logger.LogStart(ctx, "AuthService.Login")

	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		logger.Errorf("failed to get user: %v", err)
		logger.LogFinish(ctx, "AuthService.Login", err, start)
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	if user == nil {
		logger.Warnf("login attempt with non-existent email: %s", req.Email)
		logger.LogFinish(ctx, "AuthService.Login", ErrInvalidCredentials, start)
		return nil, ErrInvalidCredentials
	}

	// Verify password
	if err = s.verifyPassword(user.Password, req.Password); err != nil {
		logger.Warnf("login attempt with invalid password: %s", req.Email)
		logger.LogFinish(ctx, "AuthService.Login", ErrInvalidCredentials, start)
		return nil, ErrInvalidCredentials
	}

	logger.Infof("user logged in successfully: %s", user.Email)

	response, err := s.issueAuthResponse(user)
	logger.LogFinish(ctx, "AuthService.Login", err, start)
	return response, err
}

// issueAuthResponse generates an access token and a new refresh token rotation
// chain (new session/family) for user, and builds the resulting AuthResponse.
// Shared by Register and Login since both start a new session identically.
func (s *AuthService) issueAuthResponse(user *models.User) (*dto.AuthResponse, error) {
	accessToken, err := s.generateToken(user)
	if err != nil {
		logger.Errorf("failed to generate access token: %v", err)
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	familyID, err := s.newFamilyID()
	if err != nil {
		logger.Errorf("failed to generate token family id: %v", err)
		return nil, fmt.Errorf("failed to generate token family id: %w", err)
	}
	refreshToken, err := s.issueRefreshToken(user.ID, familyID)
	if err != nil {
		logger.Errorf("failed to issue refresh token: %v", err)
		return nil, fmt.Errorf("failed to issue refresh token: %w", err)
	}

	return &dto.AuthResponse{
		User: dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		TokenType:    "Bearer",
	}, nil
}

// ValidateToken validates a JWT token and returns the user ID.
//
// Returns error if token is invalid, expired, or malformed.
func (s *AuthService) ValidateToken(tokenString string) (uint, error) {
	secret := s.jwtSecret()
	if secret == "" {
		return 0, errors.New("JWT secret not configured")
	}

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	// Extract claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get user ID from claims
		userID, ok := claims["user_id"].(float64)
		if !ok {
			return 0, errors.New("invalid user_id in token")
		}
		return uint(userID), nil
	}

	return 0, errors.New("invalid token")
}

// GetProfile returns the current user's profile data.
//
// Returns ErrUserNotFound if userID doesn't exist (e.g. deleted after token was issued).
func (s *AuthService) GetProfile(ctx context.Context, userID uint) (*dto.UserResponse, error) {
	ctx, start := logger.LogStart(ctx, "AuthService.GetProfile")

	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		logger.Errorf("failed to get user for profile: %v", err)
		logger.LogFinish(ctx, "AuthService.GetProfile", err, start)
		return nil, fmt.Errorf("failed to get profile: %w", err)
	}
	if user == nil {
		logger.LogFinish(ctx, "AuthService.GetProfile", ErrUserNotFound, start)
		return nil, ErrUserNotFound
	}

	logger.LogFinish(ctx, "AuthService.GetProfile", nil, start)
	return &dto.UserResponse{ID: user.ID, Name: user.Name, Email: user.Email}, nil
}

// Private helper methods

// jwtSecret returns the JWT secret from config. Empty if config not loaded or not set.
func (s *AuthService) jwtSecret() string {
	if c := config.Get(); c != nil {
		return c.Server.JWTSecret
	}
	return ""
}

// hashPassword hashes a plain text password using bcrypt.
func (s *AuthService) hashPassword(password string) (string, error) {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedBytes), nil
}

// verifyPassword compares a hashed password with a plain text password.
func (s *AuthService) verifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// accessTokenTTL returns the configured JWT access token lifetime, defaulting to 15 minutes.
//
// Kept short because access tokens are stateless and cannot be revoked; logout/logout-all
// only revoke refresh tokens, so a short TTL bounds how long an already-issued access token
// keeps working after the user logs out.
func (s *AuthService) accessTokenTTL() time.Duration {
	if c := config.Get(); c != nil && c.Server.AccessTokenTTLMinutes > 0 {
		return time.Duration(c.Server.AccessTokenTTLMinutes) * time.Minute
	}
	return 15 * time.Minute
}

// generateToken creates a JWT token for authenticated user.
//
// Token contains user ID and email in claims.
// Expiry is configurable via accessTokenTTL() (default 15 minutes).
func (s *AuthService) generateToken(user *models.User) (string, error) {
	secret := s.jwtSecret()
	if secret == "" {
		return "", errors.New("JWT secret not configured")
	}

	// Create claims
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(s.accessTokenTTL()).Unix(),
		"iat":     time.Now().Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token with secret
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}

// generateRefreshToken creates a cryptographically secure random refresh token.
//
// Returns a 64-character hexadecimal string.
func (s *AuthService) generateRefreshToken() (string, error) {
	// Generate 32 random bytes (will be 64 chars in hex)
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random token: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// hashToken returns the SHA-256 hex digest of a raw token for at-rest storage.
//
// Plain SHA-256 (not bcrypt) is intentional: refresh tokens are 256-bit
// cryptographically random values, not low-entropy user secrets, so a fast
// hash is sufficient to make a stolen database dump useless without needing
// a deliberately slow KDF.
func hashToken(raw string) string {
	sum := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(sum[:])
}

// newFamilyID generates a new rotation-chain identifier, one per login session.
func (s *AuthService) newFamilyID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate family id: %w", err)
	}
	return hex.EncodeToString(bytes), nil
}

// refreshTokenTTL returns the configured refresh token lifetime, defaulting to 7 days.
func (s *AuthService) refreshTokenTTL() time.Duration {
	if c := config.Get(); c != nil && c.Server.RefreshTokenTTLDays > 0 {
		return time.Duration(c.Server.RefreshTokenTTLDays) * 24 * time.Hour
	}
	return 7 * 24 * time.Hour
}

// issueRefreshToken generates a new raw refresh token, persists its hash (with
// expiry) under familyID, and returns the raw token to send to the client.
// The raw value is never stored; only hashToken(raw) is written to the database.
func (s *AuthService) issueRefreshToken(userID uint, familyID string) (string, error) {
	raw, err := s.generateRefreshToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}

	record := &models.RefreshToken{
		UserID:    userID,
		TokenHash: hashToken(raw),
		FamilyID:  familyID,
		ExpiresAt: time.Now().Add(s.refreshTokenTTL()),
	}
	if err := s.refreshTokenRepo.Create(record); err != nil {
		return "", fmt.Errorf("failed to save refresh token: %w", err)
	}
	return raw, nil
}
