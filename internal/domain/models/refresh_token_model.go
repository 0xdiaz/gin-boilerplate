package models

import "time"

// RefreshToken represents one issued refresh token in a rotation chain.
//
// Raw tokens are never stored; TokenHash holds the SHA-256 hex digest.
// FamilyID links every token in a single rotation chain (one per login
// session), enabling reuse (theft) detection: if a token with a non-nil
// RevokedAt is presented again, the entire family must be revoked.
type RefreshToken struct {
	ID     uint `json:"id" gorm:"primaryKey"`
	UserID uint `json:"user_id" gorm:"not null;index"`

	TokenHash string `json:"-" gorm:"type:varchar(64);uniqueIndex;not null"`
	FamilyID  string `json:"-" gorm:"type:varchar(64);not null;index"`

	// ReplacedByHash is the TokenHash of the token that replaced this one via
	// rotation. Empty when this token was revoked directly (logout) instead
	// of rotated, or when it is still active.
	ReplacedByHash string `json:"-" gorm:"type:varchar(64)"`

	RevokedAt *time.Time `json:"-"`
	ExpiresAt time.Time  `json:"-" gorm:"not null"`
	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
}

// TableName specifies the database table name for RefreshToken model.
func (rt *RefreshToken) TableName() string {
	return "refresh_tokens"
}
