package domain

import (
	"time"
)

// RefreshToken represents a refresh token in the system
type RefreshToken struct {
	ID        string     `json:"id" db:"id"`
	UserID    string     `json:"userId" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expiresAt" db:"expires_at"`
	CreatedAt time.Time  `json:"createdAt" db:"created_at"`
	RevokedAt *time.Time `json:"revokedAt,omitempty" db:"revoked_at"`
}

// IsExpired checks if the refresh token has expired
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsRevoked checks if the refresh token has been revoked
func (rt *RefreshToken) IsRevoked() bool {
	return rt.RevokedAt != nil
}

// IsValid checks if the refresh token is valid (not expired and not revoked)
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsExpired() && !rt.IsRevoked()
}

// RefreshTokenRepository defines the interface for refresh token data operations
type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	GetByUserID(userID string) ([]*RefreshToken, error)
	Revoke(token string) error
	RevokeAllForUser(userID string) error
	DeleteExpired() error
}
