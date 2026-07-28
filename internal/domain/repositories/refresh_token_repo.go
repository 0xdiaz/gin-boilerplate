package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"gorm.io/gorm"
)

// RefreshTokenRepository defines data access for the refresh token rotation chain.
type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	GetByTokenHash(hash string) (*models.RefreshToken, error)
	// MarkRotated revokes the token identified by id. newTokenHash is recorded
	// as ReplacedByHash when non-empty (rotation); left empty for a terminal
	// revoke (logout).
	MarkRotated(id uint, newTokenHash string) error
	// RevokeFamily revokes every still-active token sharing familyID.
	// Used for reuse (theft) detection: a rotated-out token being replayed
	// means the whole chain must be killed.
	RevokeFamily(familyID string) error
	// RevokeAllForUser revokes every still-active token for userID (logout-all-devices).
	RevokeAllForUser(userID uint) error
}

type refreshTokenRepo struct{}

// NewRefreshTokenRepository returns a new RefreshTokenRepository implementation.
func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepo{}
}

func (r *refreshTokenRepo) Create(token *models.RefreshToken) error {
	if err := database.DB.Create(token).Error; err != nil {
		logger.Errorf("failed to create refresh token: %v", err)
		return fmt.Errorf("failed to create refresh token: %w", err)
	}
	return nil
}

func (r *refreshTokenRepo) GetByTokenHash(hash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	err := database.DB.Where("token_hash = ?", hash).First(&token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get refresh token by hash: %v", err)
		return nil, fmt.Errorf("failed to get refresh token by hash: %w", err)
	}
	return &token, nil
}

func (r *refreshTokenRepo) MarkRotated(id uint, newTokenHash string) error {
	updates := map[string]interface{}{
		"revoked_at":       time.Now(),
		"replaced_by_hash": newTokenHash,
	}
	if err := database.DB.Model(&models.RefreshToken{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		logger.Errorf("failed to mark refresh token rotated: %v", err)
		return fmt.Errorf("failed to mark refresh token rotated: %w", err)
	}
	return nil
}

func (r *refreshTokenRepo) RevokeFamily(familyID string) error {
	err := database.DB.Model(&models.RefreshToken{}).
		Where("family_id = ? AND revoked_at IS NULL", familyID).
		Update("revoked_at", time.Now()).Error
	if err != nil {
		logger.Errorf("failed to revoke refresh token family: %v", err)
		return fmt.Errorf("failed to revoke refresh token family: %w", err)
	}
	return nil
}

func (r *refreshTokenRepo) RevokeAllForUser(userID uint) error {
	err := database.DB.Model(&models.RefreshToken{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", time.Now()).Error
	if err != nil {
		logger.Errorf("failed to revoke all refresh tokens for user: %v", err)
		return fmt.Errorf("failed to revoke all refresh tokens for user: %w", err)
	}
	return nil
}
