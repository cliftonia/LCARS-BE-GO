package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/cliftonbaggerman/subspace-backend/internal/repository/memory"
	"github.com/gorilla/mux"
)

func TestMessageHandler_GetMessage(t *testing.T) {
	repo := memory.NewMessageRepository()
	handler := NewMessageHandler(repo)

	t.Run("get existing message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/messages/msg-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "msg-1"})
		w := httptest.NewRecorder()

		handler.GetMessage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var message domain.Message
		if err := json.NewDecoder(w.Body).Decode(&message); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if message.ID != "msg-1" {
			t.Errorf("expected message ID msg-1, got %s", message.ID)
		}
	})

	t.Run("get non-existent message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/messages/non-existent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "non-existent"})
		w := httptest.NewRecorder()

		handler.GetMessage(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestMessageHandler_GetUserMessages(t *testing.T) {
	repo := memory.NewMessageRepository()
	handler := NewMessageHandler(repo)

	t.Run("get messages for user with default pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-1/messages", nil)
		req = mux.SetURLVars(req, map[string]string{"userId": "user-1"})
		w := httptest.NewRecorder()

		handler.GetUserMessages(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["limit"].(float64) != 20 {
			t.Errorf("expected default limit 20, got %v", response["limit"])
		}
	})

	t.Run("get messages with custom pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-1/messages?limit=5&offset=0", nil)
		req = mux.SetURLVars(req, map[string]string{"userId": "user-1"})
		w := httptest.NewRecorder()

		handler.GetUserMessages(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["limit"].(float64) != 5 {
			t.Errorf("expected limit 5, got %v", response["limit"])
		}
	})

	t.Run("get messages with invalid pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-1/messages?limit=0", nil)
		req = mux.SetURLVars(req, map[string]string{"userId": "user-1"})
		w := httptest.NewRecorder()

		handler.GetUserMessages(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestMessageHandler_GetUnreadCount(t *testing.T) {
	repo := memory.NewMessageRepository()
	handler := NewMessageHandler(repo)

	t.Run("get unread count for user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-1/messages/unread-count", nil)
		req = mux.SetURLVars(req, map[string]string{"userId": "user-1"})
		w := httptest.NewRecorder()

		handler.GetUnreadCount(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]int
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if _, ok := response["unreadCount"]; !ok {
			t.Error("expected unreadCount in response")
		}
	})

	t.Run("get unread count for user with no messages", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/no-messages-user/messages/unread-count", nil)
		req = mux.SetURLVars(req, map[string]string{"userId": "no-messages-user"})
		w := httptest.NewRecorder()

		handler.GetUnreadCount(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]int
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["unreadCount"] != 0 {
			t.Errorf("expected unread count 0, got %d", response["unreadCount"])
		}
	})
}

func TestMessageHandler_CreateMessage(t *testing.T) {
	repo := memory.NewMessageRepository()
	handler := NewMessageHandler(repo)

	t.Run("create message successfully", func(t *testing.T) {
		message := domain.Message{
			UserID:  "user-1",
			Content: "Test message content",
			Kind:    domain.MessageKindInfo,
		}
		body, _ := json.Marshal(message)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateMessage(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var createdMessage domain.Message
		if err := json.NewDecoder(w.Body).Decode(&createdMessage); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if createdMessage.Content != "Test message content" {
			t.Errorf("expected content 'Test message content', got %s", createdMessage.Content)
		}
	})

	t.Run("create message without user ID", func(t *testing.T) {
		message := domain.Message{
			Content: "Test message content",
			Kind:    domain.MessageKindInfo,
		}
		body, _ := json.Marshal(message)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("create message with empty content", func(t *testing.T) {
		message := domain.Message{
			UserID:  "user-1",
			Content: "",
			Kind:    domain.MessageKindInfo,
		}
		body, _ := json.Marshal(message)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("create message with invalid kind", func(t *testing.T) {
		message := map[string]string{
			"userId":  "user-1",
			"content": "Test content",
			"kind":    "invalid-kind",
		}
		body, _ := json.Marshal(message)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("create message with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/messages", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateMessage(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestMessageHandler_MarkAsRead(t *testing.T) {
	repo := memory.NewMessageRepository()
	handler := NewMessageHandler(repo)

	t.Run("mark message as read successfully", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/messages/msg-1/read", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "msg-1"})
		w := httptest.NewRecorder()

		handler.MarkAsRead(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		// Verify message is marked as read
		message, _ := repo.GetByID("msg-1")
		if !message.IsRead {
			t.Error("expected message to be marked as read")
		}
	})

	t.Run("mark non-existent message as read", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/api/v1/messages/non-existent/read", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "non-existent"})
		w := httptest.NewRecorder()

		handler.MarkAsRead(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}

func TestMessageHandler_DeleteMessage(t *testing.T) {
	repo := memory.NewMessageRepository()
	handler := NewMessageHandler(repo)

	t.Run("delete message successfully", func(t *testing.T) {
		// First create a message to delete
		message := &domain.Message{
			ID:      "delete-me",
			UserID:  "user-1",
			Content: "Delete me",
			Kind:    domain.MessageKindInfo,
		}
		_ = repo.Create(message)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/messages/delete-me", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "delete-me"})
		w := httptest.NewRecorder()

		handler.DeleteMessage(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		// Verify message is deleted
		_, err := repo.GetByID("delete-me")
		if err == nil {
			t.Error("expected message to be deleted")
		}
	})

	t.Run("delete non-existent message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/messages/non-existent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "non-existent"})
		w := httptest.NewRecorder()

		handler.DeleteMessage(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
