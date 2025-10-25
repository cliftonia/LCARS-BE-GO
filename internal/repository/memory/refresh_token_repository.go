package memory

import (
	"sync"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
)

// RefreshTokenRepository implements domain.RefreshTokenRepository using in-memory storage
type RefreshTokenRepository struct {
	mu     sync.RWMutex
	tokens map[string]*domain.RefreshToken
}

// NewRefreshTokenRepository creates a new in-memory refresh token repository
func NewRefreshTokenRepository() *RefreshTokenRepository {
	return &RefreshTokenRepository{
		tokens: make(map[string]*domain.RefreshToken),
	}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.tokens[token.Token] = token
	return nil
}

// GetByToken retrieves a refresh token by its token string
func (r *RefreshTokenRepository) GetByToken(tokenStr string) (*domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	token, ok := r.tokens[tokenStr]
	if !ok {
		return nil, domain.ErrRefreshTokenNotFound
	}

	return token, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *RefreshTokenRepository) GetByUserID(userID string) ([]*domain.RefreshToken, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var tokens []*domain.RefreshToken
	for _, token := range r.tokens {
		if token.UserID == userID {
			tokens = append(tokens, token)
		}
	}

	return tokens, nil
}

// Revoke revokes a refresh token
func (r *RefreshTokenRepository) Revoke(tokenStr string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	token, ok := r.tokens[tokenStr]
	if !ok {
		return domain.ErrRefreshTokenNotFound
	}

	now := time.Now()
	token.RevokedAt = &now
	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllForUser(userID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for _, token := range r.tokens {
		if token.UserID == userID {
			token.RevokedAt = &now
		}
	}

	return nil
}

// DeleteExpired deletes all expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	for tokenStr, token := range r.tokens {
		if token.ExpiresAt.Before(now) {
			delete(r.tokens, tokenStr)
		}
	}

	return nil
}
