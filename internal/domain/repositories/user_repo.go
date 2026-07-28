package repositories

import (
	"errors"
	"fmt"

	"github.com/0xdiaz/gin-boilerplate/internal/adapters/database"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"gorm.io/gorm"
)

// UserRepository defines data access for user entity (used by auth and others).
type UserRepository interface {
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CreateUser(user *models.User) error
	UpdateUser(user *models.User) error
	GetUserByPasswordResetToken(token string) (*models.User, error)
}

// userRepo is the default implementation of UserRepository.
type userRepo struct{}

// NewUserRepository returns a new UserRepository implementation.
func NewUserRepository() UserRepository {
	return &userRepo{}
}

func (r *userRepo) CreateUser(user *models.User) error {
	if err := database.DB.Create(user).Error; err != nil {
		logger.Errorf("failed to create user: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}
	return nil
}

func (r *userRepo) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get user by email: %v", err)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}
	return &user, nil
}

func (r *userRepo) UpdateUser(user *models.User) error {
	if err := database.DB.Save(user).Error; err != nil {
		logger.Errorf("failed to update user: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *userRepo) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := database.DB.First(&user, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get user by id: %v", err)
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}
	return &user, nil
}

func (r *userRepo) GetUserByPasswordResetToken(token string) (*models.User, error) {
	var user models.User
	err := database.DB.Where("password_reset_token = ?", token).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		logger.Errorf("failed to get user by password reset token: %v", err)
		return nil, fmt.Errorf("failed to get user by password reset token: %w", err)
	}
	return &user, nil
}
