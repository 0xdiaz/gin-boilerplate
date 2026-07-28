package seeders

import (
	"fmt"

	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
	"github.com/0xdiaz/gin-boilerplate/pkg/logger"
	"golang.org/x/crypto/bcrypt"
)

// demoPassword is the login password for every seeded demo user.
// Local development convenience only; never used outside IsDevelopment().
const demoPassword = "Password123!"

// demoUsers are inserted (upsert-by-email) so the API is immediately
// login-ready after `go run main.go` in a fresh development database.
var demoUsers = []struct {
	Name  string
	Email string
}{
	{"Demo Admin", "admin@example.local"},
	{"Demo User", "user@example.local"},
}

func seedDemoUsers() error {
	userRepo := repositories.NewUserRepository()

	hashed, err := bcrypt.GenerateFromPassword([]byte(demoPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash demo password: %w", err)
	}

	for _, u := range demoUsers {
		existing, err := userRepo.GetUserByEmail(u.Email)
		if err != nil {
			return err
		}
		if existing != nil {
			continue
		}

		if err := userRepo.CreateUser(&models.User{
			Name:     u.Name,
			Email:    u.Email,
			Password: string(hashed),
		}); err != nil {
			return err
		}
	}

	logger.Infof("Seeded %d demo users", len(demoUsers))
	return nil
}
