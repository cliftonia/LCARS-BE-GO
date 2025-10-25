// Package postgres provides PostgreSQL repository implementations.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
)

// MessageRepository implements domain.MessageRepository for PostgreSQL
type MessageRepository struct {
	db *sql.DB
}

// NewMessageRepository creates a new PostgreSQL message repository
func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

// Create creates a new message
func (r *MessageRepository) Create(message *domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		INSERT INTO messages (id, user_id, content, kind, is_read, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		message.ID,
		message.UserID,
		message.Content,
		message.Kind,
		message.IsRead,
		message.CreatedAt,
		message.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create message: %w", err)
	}

	return nil
}

// GetByID retrieves a message by ID
func (r *MessageRepository) GetByID(id string) (*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, content, kind, is_read, created_at, updated_at
		FROM messages
		WHERE id = $1
	`

	message := &domain.Message{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&message.ID,
		&message.UserID,
		&message.Content,
		&message.Kind,
		&message.IsRead,
		&message.CreatedAt,
		&message.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrMessageNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}

	return message, nil
}

// Update updates an existing message
func (r *MessageRepository) Update(message *domain.Message) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set updated_at to current time
	message.UpdatedAt = time.Now()

	query := `
		UPDATE messages
		SET content = $1, kind = $2, updated_at = $3
		WHERE id = $4
	`

	result, err := r.db.ExecContext(ctx, query,
		message.Content,
		message.Kind,
		message.UpdatedAt,
		message.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrMessageNotFound
	}

	return nil
}

// Delete deletes a message by ID
func (r *MessageRepository) Delete(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `DELETE FROM messages WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrMessageNotFound
	}

	return nil
}

// GetByUserID returns messages for a specific user with pagination
func (r *MessageRepository) GetByUserID(userID string, limit, offset int) ([]*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, content, kind, is_read, created_at, updated_at
		FROM messages
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages by user: %w", err)
	}
	defer func() { _ = rows.Close() }()

	messages := make([]*domain.Message, 0)
	for rows.Next() {
		message := &domain.Message{}
		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.Kind,
			&message.IsRead,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

// List returns a paginated list of all messages
func (r *MessageRepository) List(limit, offset int) ([]*domain.Message, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, content, kind, created_at, updated_at
		FROM messages
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list messages: %w", err)
	}
	defer func() { _ = rows.Close() }()

	messages := make([]*domain.Message, 0)
	for rows.Next() {
		message := &domain.Message{}
		err := rows.Scan(
			&message.ID,
			&message.UserID,
			&message.Content,
			&message.Kind,
			&message.CreatedAt,
			&message.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan message: %w", err)
		}
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating messages: %w", err)
	}

	return messages, nil
}

// CountByUserID returns the total number of messages for a specific user
func (r *MessageRepository) CountByUserID(userID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM messages WHERE user_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count messages by user: %w", err)
	}

	return count, nil
}

// GetUnreadCount returns the number of unread messages for a specific user
func (r *MessageRepository) GetUnreadCount(userID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `SELECT COUNT(*) FROM messages WHERE user_id = $1 AND is_read = false`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count unread messages: %w", err)
	}

	return count, nil
}

// MarkAsRead marks a message as read
func (r *MessageRepository) MarkAsRead(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
		UPDATE messages
		SET is_read = true, updated_at = $1
		WHERE id = $2
	`

	result, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to mark message as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrMessageNotFound
	}

	return nil
}
