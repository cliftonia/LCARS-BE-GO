package memory

import (
	"sync"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/google/uuid"
)

// UserRepository is an in-memory implementation of domain.UserRepository
type UserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

// NewUserRepository creates a new in-memory user repository with mock data
func NewUserRepository() *UserRepository {
	repo := &UserRepository{
		users: make(map[string]*domain.User),
	}

	// Add mock data
	mockUsers := []*domain.User{
		{
			ID:        "user-1",
			Name:      "John Doe",
			Email:     "john.doe@example.com",
			AvatarURL: stringPtr("https://i.pravatar.cc/150?img=1"),
			CreatedAt: time.Now().Add(-30 * 24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user-2",
			Name:      "Jane Smith",
			Email:     "jane.smith@example.com",
			AvatarURL: stringPtr("https://i.pravatar.cc/150?img=2"),
			CreatedAt: time.Now().Add(-20 * 24 * time.Hour),
			UpdatedAt: time.Now(),
		},
		{
			ID:        "user-3",
			Name:      "Bob Johnson",
			Email:     "bob.johnson@example.com",
			AvatarURL: nil,
			CreatedAt: time.Now().Add(-10 * 24 * time.Hour),
			UpdatedAt: time.Now(),
		},
	}

	for _, user := range mockUsers {
		repo.users[user.ID] = user
	}

	return repo
}

// GetByID retrieves a user by ID
func (r *UserRepository) GetByID(id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}

	return user, nil
}

// GetByEmail retrieves a user by email
func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, domain.ErrUserNotFound
}

// List retrieves a list of users with pagination
func (r *UserRepository) List(limit, offset int) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	users := make([]*domain.User, 0, len(r.users))
	for _, user := range r.users {
		users = append(users, user)
	}

	// Simple pagination
	start := offset
	if start > len(users) {
		return []*domain.User{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

// Count returns the total number of users
func (r *UserRepository) Count() (int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return len(r.users), nil
}

// Create creates a new user
func (r *UserRepository) Create(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	r.users[user.ID] = user
	return nil
}

// Update updates an existing user
func (r *UserRepository) Update(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}

	user.UpdatedAt = time.Now()
	r.users[user.ID] = user
	return nil
}

// Delete deletes a user by ID
func (r *UserRepository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}

	delete(r.users, id)
	return nil
}

func stringPtr(s string) *string {
	return &s
}
