package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/auth"
	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/cliftonbaggerman/subspace-backend/internal/repository/memory"
)

func TestAuthHandler_Register(t *testing.T) {
	userRepo := memory.NewUserRepository()
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	handler := NewAuthHandler(userRepo, jwtManager)

	t.Run("successful registration", func(t *testing.T) {
		reqBody := RegisterRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Register(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", w.Code)
		}

		var response AuthResponse
		_ = json.NewDecoder(w.Body).Decode(&response)

		if response.Token == "" {
			t.Error("expected token in response")
		}

		if response.User.Email != reqBody.Email {
			t.Errorf("expected email %s, got %s", reqBody.Email, response.User.Email)
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		reqBody := RegisterRequest{
			Name:     "Test User",
			Email:    "invalid-email",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Register(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})

	t.Run("password too short", func(t *testing.T) {
		reqBody := RegisterRequest{
			Name:     "Test User",
			Email:    "test2@example.com",
			Password: "short",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Register(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", w.Code)
		}
	})
}

func TestAuthHandler_Login(t *testing.T) {
	userRepo := memory.NewUserRepository()
	jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)
	handler := NewAuthHandler(userRepo, jwtManager)

	// Register a user first
	hashedPassword, _ := auth.HashPassword("password123")
	user := &domain.User{
		ID:       "test-user-1",
		Name:     "Test User",
		Email:    "login@example.com",
		Password: hashedPassword,
	}
	_ = userRepo.Create(user)

	t.Run("successful login", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "login@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var response AuthResponse
		_ = json.NewDecoder(w.Body).Decode(&response)

		if response.Token == "" {
			t.Error("expected token in response")
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "login@example.com",
			Password: "wrongpassword",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		reqBody := LoginRequest{
			Email:    "notfound@example.com",
			Password: "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		w := httptest.NewRecorder()

		handler.Login(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}
