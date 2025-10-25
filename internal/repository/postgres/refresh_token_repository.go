// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/google/uuid"
)

// RefreshTokenRepository implements domain.RefreshTokenRepository for PostgreSQL
type RefreshTokenRepository struct {
	db *sql.DB
}

// NewRefreshTokenRepository creates a new PostgreSQL refresh token repository
func NewRefreshTokenRepository(db *sql.DB) *RefreshTokenRepository {
	return &RefreshTokenRepository{db: db}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	if token.ID == "" {
		token.ID = uuid.New().String()
	}
	if token.CreatedAt.IsZero() {
		token.CreatedAt = time.Now()
	}

	_, err := r.db.ExecContext(
		context.Background(),
		query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.CreatedAt,
	)

	return err
}

// GetByToken retrieves a refresh token by its token string
func (r *RefreshTokenRepository) GetByToken(tokenStr string) (*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE token = $1
	`

	var token domain.RefreshToken
	err := r.db.QueryRowContext(context.Background(), query, tokenStr).Scan(
		&token.ID,
		&token.UserID,
		&token.Token,
		&token.ExpiresAt,
		&token.CreatedAt,
		&token.RevokedAt,
	)

	if err == sql.ErrNoRows {
		return nil, domain.ErrRefreshTokenNotFound
	}
	if err != nil {
		return nil, err
	}

	return &token, nil
}

// GetByUserID retrieves all refresh tokens for a user
func (r *RefreshTokenRepository) GetByUserID(userID string) ([]*domain.RefreshToken, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at, revoked_at
		FROM refresh_tokens
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = rows.Close()
	}()

	var tokens []*domain.RefreshToken
	for rows.Next() {
		var token domain.RefreshToken
		if err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.Token,
			&token.ExpiresAt,
			&token.CreatedAt,
			&token.RevokedAt,
		); err != nil {
			return nil, err
		}
		tokens = append(tokens, &token)
	}

	return tokens, rows.Err()
}

// Revoke revokes a refresh token
func (r *RefreshTokenRepository) Revoke(tokenStr string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE token = $2 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(context.Background(), query, time.Now(), tokenStr)
	return err
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllForUser(userID string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked_at = $1
		WHERE user_id = $2 AND revoked_at IS NULL
	`

	_, err := r.db.ExecContext(context.Background(), query, time.Now(), userID)
	return err
}

// DeleteExpired deletes all expired refresh tokens
func (r *RefreshTokenRepository) DeleteExpired() error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < $1
	`

	_, err := r.db.ExecContext(context.Background(), query, time.Now())
	return err
}
