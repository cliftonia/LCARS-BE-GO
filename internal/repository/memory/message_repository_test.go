package memory

import (
	"testing"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
)

func TestMessageRepository_GetByID(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("get existing message", func(t *testing.T) {
		message, err := repo.GetByID("msg-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if message.ID != "msg-1" {
			t.Errorf("expected message ID 'msg-1', got %s", message.ID)
		}
	})

	t.Run("get non-existent message", func(t *testing.T) {
		_, err := repo.GetByID("non-existent")
		if err != domain.ErrMessageNotFound {
			t.Errorf("expected ErrMessageNotFound, got %v", err)
		}
	})
}

func TestMessageRepository_GetByUserID(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("get messages for user with messages", func(t *testing.T) {
		messages, err := repo.GetByUserID("user-1", 10, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(messages) != 3 {
			t.Errorf("expected 3 messages for user-1, got %d", len(messages))
		}
	})

	t.Run("get messages with pagination", func(t *testing.T) {
		messages, err := repo.GetByUserID("user-1", 2, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(messages) != 2 {
			t.Errorf("expected 2 messages, got %d", len(messages))
		}
	})

	t.Run("get messages for user with no messages", func(t *testing.T) {
		messages, err := repo.GetByUserID("user-999", 10, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(messages) != 0 {
			t.Errorf("expected 0 messages, got %d", len(messages))
		}
	})
}

func TestMessageRepository_CountByUserID(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("count messages for user-1", func(t *testing.T) {
		count, err := repo.CountByUserID("user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 3 {
			t.Errorf("expected count of 3, got %d", count)
		}
	})

	t.Run("count messages for user with no messages", func(t *testing.T) {
		count, err := repo.CountByUserID("user-999")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected count of 0, got %d", count)
		}
	})
}

func TestMessageRepository_GetUnreadCount(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("get unread count for user-1", func(t *testing.T) {
		count, err := repo.GetUnreadCount("user-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 2 {
			t.Errorf("expected 2 unread messages, got %d", count)
		}
	})

	t.Run("get unread count for user with no unread messages", func(t *testing.T) {
		count, err := repo.GetUnreadCount("user-999")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if count != 0 {
			t.Errorf("expected 0 unread messages, got %d", count)
		}
	})
}

func TestMessageRepository_Create(t *testing.T) {
	repo := NewMessageRepository()

	newMessage := &domain.Message{
		UserID:  "user-1",
		Content: "Test message",
		Kind:    domain.MessageKindInfo,
		IsRead:  false,
	}

	err := repo.Create(newMessage)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if newMessage.ID == "" {
		t.Error("expected ID to be generated")
	}

	if newMessage.CreatedAt.IsZero() {
		t.Error("expected CreatedAt to be set")
	}

	// Verify message was created
	count, _ := repo.CountByUserID("user-1")
	if count != 4 {
		t.Errorf("expected count of 4 after creation, got %d", count)
	}
}

func TestMessageRepository_Update(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("update existing message", func(t *testing.T) {
		message, _ := repo.GetByID("msg-1")
		message.Content = "Updated content"

		err := repo.Update(message)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updatedMessage, _ := repo.GetByID("msg-1")
		if updatedMessage.Content != "Updated content" {
			t.Errorf("expected content 'Updated content', got %s", updatedMessage.Content)
		}
	})

	t.Run("update non-existent message", func(t *testing.T) {
		message := &domain.Message{
			ID:      "non-existent",
			UserID:  "user-1",
			Content: "Test",
			Kind:    domain.MessageKindInfo,
		}

		err := repo.Update(message)
		if err != domain.ErrMessageNotFound {
			t.Errorf("expected ErrMessageNotFound, got %v", err)
		}
	})
}

func TestMessageRepository_MarkAsRead(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("mark existing unread message as read", func(t *testing.T) {
		err := repo.MarkAsRead("msg-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		message, _ := repo.GetByID("msg-1")
		if !message.IsRead {
			t.Error("expected message to be marked as read")
		}
	})

	t.Run("mark non-existent message as read", func(t *testing.T) {
		err := repo.MarkAsRead("non-existent")
		if err != domain.ErrMessageNotFound {
			t.Errorf("expected ErrMessageNotFound, got %v", err)
		}
	})
}

func TestMessageRepository_Delete(t *testing.T) {
	repo := NewMessageRepository()

	t.Run("delete existing message", func(t *testing.T) {
		err := repo.Delete("msg-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = repo.GetByID("msg-1")
		if err != domain.ErrMessageNotFound {
			t.Errorf("expected ErrMessageNotFound after deletion, got %v", err)
		}
	})

	t.Run("delete non-existent message", func(t *testing.T) {
		err := repo.Delete("non-existent")
		if err != domain.ErrMessageNotFound {
			t.Errorf("expected ErrMessageNotFound, got %v", err)
		}
	})
}
