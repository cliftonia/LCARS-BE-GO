package domain

import (
	"time"
)

// User represents a user in the system
type User struct {
	ID        string    `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password_hash"` // Never include in JSON responses
	AvatarURL *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// UserRepository defines the interface for user data operations
type UserRepository interface {
	GetByID(id string) (*User, error)
	GetByEmail(email string) (*User, error)
	List(limit, offset int) ([]*User, error)
	Count() (int, error)
	Create(user *User) error
	Update(user *User) error
	Delete(id string) error
}
