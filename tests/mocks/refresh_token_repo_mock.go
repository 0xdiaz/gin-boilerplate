package mocks

import (
	"sync"
	"time"

	"github.com/0xdiaz/gin-boilerplate/internal/domain/models"
	"github.com/0xdiaz/gin-boilerplate/internal/domain/repositories"
)

// MockRefreshTokenRepository is an in-memory RefreshTokenRepository for unit tests.
type MockRefreshTokenRepository struct {
	mu     sync.RWMutex
	nextID uint
	// key: token hash
	byHash map[string]*models.RefreshToken
}

// NewMockRefreshTokenRepository returns a new MockRefreshTokenRepository.
func NewMockRefreshTokenRepository() *MockRefreshTokenRepository {
	return &MockRefreshTokenRepository{byHash: make(map[string]*models.RefreshToken)}
}

// Ensure MockRefreshTokenRepository implements repositories.RefreshTokenRepository.
var _ repositories.RefreshTokenRepository = (*MockRefreshTokenRepository)(nil)

// Seed inserts a token record directly, for setting up "already existing" test state.
func (m *MockRefreshTokenRepository) Seed(token *models.RefreshToken) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byHash == nil {
		m.byHash = make(map[string]*models.RefreshToken)
	}
	if token.ID == 0 {
		m.nextID++
		token.ID = m.nextID
	}
	m.byHash[token.TokenHash] = token
}

func (m *MockRefreshTokenRepository) Create(token *models.RefreshToken) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.byHash == nil {
		m.byHash = make(map[string]*models.RefreshToken)
	}
	m.nextID++
	token.ID = m.nextID
	cp := *token
	m.byHash[token.TokenHash] = &cp
	return nil
}

func (m *MockRefreshTokenRepository) GetByTokenHash(hash string) (*models.RefreshToken, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if t, ok := m.byHash[hash]; ok {
		cp := *t
		return &cp, nil
	}
	return nil, nil
}

func (m *MockRefreshTokenRepository) MarkRotated(id uint, newTokenHash string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, t := range m.byHash {
		if t.ID == id {
			now := time.Now()
			t.RevokedAt = &now
			t.ReplacedByHash = newTokenHash
			return nil
		}
	}
	return nil
}

func (m *MockRefreshTokenRepository) RevokeFamily(familyID string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for _, t := range m.byHash {
		if t.FamilyID == familyID && t.RevokedAt == nil {
			t.RevokedAt = &now
		}
	}
	return nil
}

func (m *MockRefreshTokenRepository) RevokeAllForUser(userID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	now := time.Now()
	for _, t := range m.byHash {
		if t.UserID == userID && t.RevokedAt == nil {
			t.RevokedAt = &now
		}
	}
	return nil
}
