package models

import (
	"time"
)

// User represents a user account in the system.
//
// This model is used for authentication and user management.
// Password field stores bcrypt hashed passwords only.
type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name" gorm:"type:varchar(255);not null"`
	Email    string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
	Password string `json:"-" gorm:"type:varchar(255);not null"` // Never expose in JSON

	// PasswordResetToken stores the SHA-256 hex digest of the reset token, never the raw value.
	PasswordResetToken  string     `json:"-" gorm:"type:varchar(255);index"`
	PasswordResetExpiry *time.Time `json:"-" gorm:"type:timestamp"`

	CreatedAt time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"` // Soft delete support
}

// TableName specifies the database table name for User model.
func (u *User) TableName() string {
	return "users"
}
