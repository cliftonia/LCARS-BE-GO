// Package memory provides in-memory repository implementations for testing.
package memory

import (
	"sync"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/google/uuid"
)

// MessageRepository is an in-memory implementation of domain.MessageRepository
type MessageRepository struct {
	mu       sync.RWMutex
	messages map[string]*domain.Message
}

// NewMessageRepository creates a new in-memory message repository with mock data
func NewMessageRepository() *MessageRepository {
	repo := &MessageRepository{
		messages: make(map[string]*domain.Message),
	}

	// Add mock data
	mockMessages := []*domain.Message{
		{
			ID:        "msg-1",
			UserID:    "user-1",
			Content:   "Welcome to Subspace!",
			Kind:      domain.MessageKindSuccess,
			IsRead:    false,
			CreatedAt: time.Now().Add(-2 * time.Hour),
			UpdatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        "msg-2",
			UserID:    "user-1",
			Content:   "Your profile has been updated successfully.",
			Kind:      domain.MessageKindInfo,
			IsRead:    true,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-30 * time.Minute),
		},
		{
			ID:        "msg-3",
			UserID:    "user-1",
			Content:   "New features are available. Check them out!",
			Kind:      domain.MessageKindInfo,
			IsRead:    false,
			CreatedAt: time.Now().Add(-30 * time.Minute),
			UpdatedAt: time.Now().Add(-30 * time.Minute),
		},
		{
			ID:        "msg-4",
			UserID:    "user-2",
			Content:   "Your password will expire in 7 days.",
			Kind:      domain.MessageKindWarning,
			IsRead:    false,
			CreatedAt: time.Now().Add(-1 * time.Hour),
			UpdatedAt: time.Now().Add(-1 * time.Hour),
		},
	}

	for _, message := range mockMessages {
		repo.messages[message.ID] = message
	}

	return repo
}

// GetByID retrieves a message by ID
func (r *MessageRepository) GetByID(id string) (*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	message, exists := r.messages[id]
	if !exists {
		return nil, domain.ErrMessageNotFound
	}

	return message, nil
}

// GetByUserID retrieves messages for a specific user with pagination
func (r *MessageRepository) GetByUserID(userID string, limit, offset int) ([]*domain.Message, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userMessages := make([]*domain.Message, 0)
	for _, message := range r.messages {
		if message.UserID == userID {
			userMessages = append(userMessages, message)
		}
	}

	// Simple pagination
	start := offset
	if start > len(userMessages) {
		return []*domain.Message{}, nil
	}

	end := start + limit
	if end > len(userMessages) {
		end = len(userMessages)
	}

	return userMessages[start:end], nil
}

// CountByUserID returns the total count of messages for a user
func (r *MessageRepository) CountByUserID(userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, message := range r.messages {
		if message.UserID == userID {
			count++
		}
	}

	return count, nil
}

// GetUnreadCount returns the count of unread messages for a user
func (r *MessageRepository) GetUnreadCount(userID string) (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	count := 0
	for _, message := range r.messages {
		if message.UserID == userID && !message.IsRead {
			count++
		}
	}

	return count, nil
}

// Create creates a new message
func (r *MessageRepository) Create(message *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if message.ID == "" {
		message.ID = uuid.New().String()
	}

	now := time.Now()
	message.CreatedAt = now
	message.UpdatedAt = now

	r.messages[message.ID] = message
	return nil
}

// Update updates an existing message
func (r *MessageRepository) Update(message *domain.Message) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.messages[message.ID]; !exists {
		return domain.ErrMessageNotFound
	}

	message.UpdatedAt = time.Now()
	r.messages[message.ID] = message
	return nil
}

// MarkAsRead marks a message as read
func (r *MessageRepository) MarkAsRead(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	message, exists := r.messages[id]
	if !exists {
		return domain.ErrMessageNotFound
	}

	message.IsRead = true
	message.UpdatedAt = time.Now()
	return nil
}

// Delete deletes a message by ID
func (r *MessageRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.messages[id]; !exists {
		return domain.ErrMessageNotFound
	}

	delete(r.messages, id)
	return nil
}
