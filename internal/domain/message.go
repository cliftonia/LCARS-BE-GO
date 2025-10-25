package domain

import (
	"time"
)

// MessageKind represents the type of message
type MessageKind string

// MessageKind constants define the types of messages
const (
	MessageKindInfo    MessageKind = "info"
	MessageKindWarning MessageKind = "warning"
	MessageKindError   MessageKind = "error"
	MessageKindSuccess MessageKind = "success"
)

// Message represents a message in the system
type Message struct {
	ID        string      `json:"id" db:"id"`
	UserID    string      `json:"userId" db:"user_id"`
	Content   string      `json:"content" db:"content"`
	Kind      MessageKind `json:"kind" db:"kind"`
	IsRead    bool        `json:"isRead" db:"is_read"`
	CreatedAt time.Time   `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time   `json:"updatedAt" db:"updated_at"`
}

// MessageRepository defines the interface for message data operations
type MessageRepository interface {
	GetByID(id string) (*Message, error)
	GetByUserID(userID string, limit, offset int) ([]*Message, error)
	CountByUserID(userID string) (int, error)
	GetUnreadCount(userID string) (int, error)
	Create(message *Message) error
	Update(message *Message) error
	MarkAsRead(id string) error
	Delete(id string) error
}
