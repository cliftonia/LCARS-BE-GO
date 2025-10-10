package memory

import (
	"testing"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
)

func TestUserRepository_GetByID(t *testing.T) {
	repo := NewUserRepository()

	t.Run("get existing user", func(t *testing.T) {
		user, err := repo.GetByID("user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.ID != "user-1" {
			t.Errorf("expected user ID 'user-1', got %s", user.ID)
		}
	})

	t.Run("get non-existent user", func(t *testing.T) {
		_, err := repo.GetByID("non-existent")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserRepository_GetByEmail(t *testing.T) {
	repo := NewUserRepository()

	t.Run("get existing user by email", func(t *testing.T) {
		user, err := repo.GetByEmail("john.doe@example.com")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if user.Email != "john.doe@example.com" {
			t.Errorf("expected email 'john.doe@example.com', got %s", user.Email)
		}
	})

	t.Run("get non-existent user by email", func(t *testing.T) {
		_, err := repo.GetByEmail("nonexistent@example.com")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserRepository_List(t *testing.T) {
	repo := NewUserRepository()

	t.Run("list all users", func(t *testing.T) {
		users, err := repo.List(10, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 3 {
			t.Errorf("expected 3 users, got %d", len(users))
		}
	})

	t.Run("list with pagination", func(t *testing.T) {
		users, err := repo.List(2, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 2 {
			t.Errorf("expected 2 users, got %d", len(users))
		}
	})

	t.Run("list with offset beyond data", func(t *testing.T) {
		users, err := repo.List(10, 100)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(users) != 0 {
			t.Errorf("expected 0 users, got %d", len(users))
		}
	})
}

func TestUserRepository_Count(t *testing.T) {
	repo := NewUserRepository()

	count, err := repo.Count()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if count != 3 {
		t.Errorf("expected count of 3, got %d", count)
	}
}

func TestUserRepository_Create(t *testing.T) {
	repo := NewUserRepository()

	newUser := &domain.User{
		Name:  "Test User",
		Email: "test@example.com",
	}

	err := repo.Create(newUser)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newUser.ID == "" {
		t.Error("expected ID to be generated")
	}

	if newUser.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	// Verify user was created
	count, _ := repo.Count()
	if count != 4 {
		t.Errorf("expected count of 4 after creation, got %d", count)
	}
}

func TestUserRepository_Update(t *testing.T) {
	repo := NewUserRepository()

	t.Run("update existing user", func(t *testing.T) {
		user, _ := repo.GetByID("user-1")
		user.Name = "Updated Name"

		err := repo.Update(user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updatedUser, _ := repo.GetByID("user-1")
		if updatedUser.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", updatedUser.Name)
		}
	})

	t.Run("update non-existent user", func(t *testing.T) {
		user := &domain.User{
			ID:    "non-existent",
			Name:  "Test",
			Email: "test@example.com",
		}

		err := repo.Update(user)
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestUserRepository_Delete(t *testing.T) {
	repo := NewUserRepository()

	t.Run("delete existing user", func(t *testing.T) {
		err := repo.Delete("user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.GetByID("user-1")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound after deletion, got %v", err)
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		err := repo.Delete("non-existent")
		if err != domain.ErrUserNotFound {
			t.Errorf("expected ErrUserNotFound, got %v", err)
		}
	})
}
