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

func TestUserHandler_GetUser(t *testing.T) {
	repo := memory.NewUserRepository()
	handler := NewUserHandler(repo)

	t.Run("get existing user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/user-1", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "user-1"})
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var user domain.User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user.ID != "user-1" {
			t.Errorf("expected user ID user-1, got %s", user.ID)
		}
	})

	t.Run("get non-existent user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users/non-existent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "non-existent"})
		w := httptest.NewRecorder()

		handler.GetUser(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", w.Code)
		}
	})
}

func TestUserHandler_ListUsers(t *testing.T) {
	repo := memory.NewUserRepository()
	handler := NewUserHandler(repo)

	t.Run("list users with default pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

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

		if response["offset"].(float64) != 0 {
			t.Errorf("expected default offset 0, got %v", response["offset"])
		}
	})

	t.Run("list users with custom pagination", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=10&offset=1", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["limit"].(float64) != 10 {
			t.Errorf("expected limit 10, got %v", response["limit"])
		}

		if response["offset"].(float64) != 1 {
			t.Errorf("expected offset 1, got %v", response["offset"])
		}
	})

	t.Run("list users with invalid limit", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=0", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("list users with limit exceeding max", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users?limit=200", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	repo := memory.NewUserRepository()
	handler := NewUserHandler(repo)

	t.Run("create user successfully", func(t *testing.T) {
		user := domain.User{
			Name:  "New User",
			Email: "newuser@example.com",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var createdUser domain.User
		if err := json.NewDecoder(w.Body).Decode(&createdUser); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if createdUser.Name != "New User" {
			t.Errorf("expected name 'New User', got %s", createdUser.Name)
		}

		if createdUser.Email != "newuser@example.com" {
			t.Errorf("expected email 'newuser@example.com', got %s", createdUser.Email)
		}
	})

	t.Run("create user with invalid email", func(t *testing.T) {
		user := domain.User{
			Name:  "Invalid Email User",
			Email: "invalid-email",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("create user with empty name", func(t *testing.T) {
		user := domain.User{
			Name:  "",
			Email: "valid@example.com",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("create user with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/v1/users", bytes.NewBufferString("invalid json"))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.CreateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestUserHandler_UpdateUser(t *testing.T) {
	repo := memory.NewUserRepository()
	handler := NewUserHandler(repo)

	t.Run("update user successfully", func(t *testing.T) {
		user := domain.User{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "user-1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var updatedUser domain.User
		if err := json.NewDecoder(w.Body).Decode(&updatedUser); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if updatedUser.Name != "Updated Name" {
			t.Errorf("expected name 'Updated Name', got %s", updatedUser.Name)
		}
	})

	t.Run("update user with invalid email", func(t *testing.T) {
		user := domain.User{
			Name:  "Valid Name",
			Email: "invalid-email",
		}
		body, _ := json.Marshal(user)

		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-1", bytes.NewBuffer(body))
		req = mux.SetURLVars(req, map[string]string{"id": "user-1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("update user with invalid JSON", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/api/v1/users/user-1", bytes.NewBufferString("invalid json"))
		req = mux.SetURLVars(req, map[string]string{"id": "user-1"})
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.UpdateUser(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestUserHandler_DeleteUser(t *testing.T) {
	repo := memory.NewUserRepository()
	handler := NewUserHandler(repo)

	t.Run("delete user successfully", func(t *testing.T) {
		// First create a user to delete
		user := &domain.User{
			ID:    "delete-me",
			Name:  "Delete Me",
			Email: "deleteme@example.com",
		}
		_ = repo.Create(user)

		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/delete-me", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "delete-me"})
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		// Verify user is deleted
		_, err := repo.GetByID("delete-me")
		if err == nil {
			t.Error("expected user to be deleted")
		}
	})

	t.Run("delete non-existent user", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/api/v1/users/non-existent", nil)
		req = mux.SetURLVars(req, map[string]string{"id": "non-existent"})
		w := httptest.NewRecorder()

		handler.DeleteUser(w, req)

		// Memory repo returns error for non-existent users
		if w.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", w.Code)
		}
	})
}
