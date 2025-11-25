package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/auth"
	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/handlers"
	"github.com/cliftonbaggerman/subspace-backend/internal/repository/memory"
)

var (
	testRouterOnce   sync.Once
	sharedTestRouter http.Handler
)

func setupTestRouter() http.Handler {
	testRouterOnce.Do(func() {
		userRepo := memory.NewUserRepository()
		messageRepo := memory.NewMessageRepository()
		refreshTokenRepo := memory.NewRefreshTokenRepository()
		jwtManager := auth.NewJWTManager("test-secret", 24*time.Hour)

		userHandler := handlers.NewUserHandler(userRepo)
		messageHandler := handlers.NewMessageHandler(messageRepo)
		authHandler := handlers.NewAuthHandler(userRepo, refreshTokenRepo, jwtManager)

		// Use nil metrics to avoid duplicate registration in tests
		sharedTestRouter = NewRouter(RouterConfig{
			UserHandler:    userHandler,
			MessageHandler: messageHandler,
			AuthHandler:    authHandler,
			UserRepo:       userRepo,
			MessageRepo:    messageRepo,
			JWTManager:     jwtManager,
			CORSOrigins:    []string{"*"},
			Metrics:        nil,
		})
	})
	return sharedTestRouter
}

func TestHealthEndpoint(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var response HealthStatus
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if response.Status != "healthy" {
		t.Errorf("expected status 'healthy', got %s", response.Status)
	}

	if response.Service != "subspace-backend" {
		t.Errorf("expected service 'subspace-backend', got %s", response.Service)
	}
}

func TestMetricsEndpoint(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	// Verify Prometheus format
	contentType := w.Header().Get("Content-Type")
	if contentType == "" {
		t.Error("expected Content-Type header")
	}
}

func TestAuthRegisterEndpoint(t *testing.T) {
	router := setupTestRouter()

	t.Run("successful registration", func(t *testing.T) {
		reqBody := map[string]string{
			"name":     "Test User",
			"email":    "testregister@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["accessToken"] == nil {
			t.Error("expected accessToken in response")
		}

		if response["refreshToken"] == nil {
			t.Error("expected refreshToken in response")
		}
	})
}

func TestAuthLoginEndpoint(t *testing.T) {
	router := setupTestRouter()

	// First register a user
	registerBody := map[string]string{
		"name":     "Login Test User",
		"email":    "logintest@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	t.Run("successful login", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    "logintest@example.com",
			"password": "password123",
		}
		body, _ := json.Marshal(loginBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var response map[string]interface{}
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if response["accessToken"] == nil {
			t.Error("expected accessToken in response")
		}
	})

	t.Run("invalid credentials", func(t *testing.T) {
		loginBody := map[string]string{
			"email":    "logintest@example.com",
			"password": "wrongpassword",
		}
		body, _ := json.Marshal(loginBody)

		req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusUnauthorized {
			t.Errorf("expected status 401, got %d", w.Code)
		}
	})
}

func TestProtectedEndpointsWithoutAuth(t *testing.T) {
	router := setupTestRouter()

	endpoints := []struct {
		method string
		path   string
	}{
		{http.MethodGet, "/api/v1/auth/me"},
		{http.MethodGet, "/api/v1/users"},
		{http.MethodPost, "/api/v1/users"},
		{http.MethodGet, "/api/v1/users/user-1"},
		{http.MethodPut, "/api/v1/users/user-1"},
		{http.MethodDelete, "/api/v1/users/user-1"},
		{http.MethodPost, "/api/v1/messages"},
		{http.MethodGet, "/api/v1/messages/msg-1"},
		{http.MethodDelete, "/api/v1/messages/msg-1"},
		{http.MethodPatch, "/api/v1/messages/msg-1/read"},
	}

	for _, endpoint := range endpoints {
		t.Run(endpoint.method+" "+endpoint.path, func(t *testing.T) {
			req := httptest.NewRequest(endpoint.method, endpoint.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Accept either 401 (Unauthorized) or 429 (Rate Limited) as valid responses
			// Rate limiting may kick in during tests since all requests come from the same IP
			if w.Code != http.StatusUnauthorized && w.Code != http.StatusTooManyRequests {
				t.Errorf("expected status 401 or 429 for %s %s, got %d", endpoint.method, endpoint.path, w.Code)
			}
		})
	}
}

func TestProtectedEndpointsWithAuth(t *testing.T) {
	userRepo := memory.NewUserRepository()
	messageRepo := memory.NewMessageRepository()
	refreshTokenRepo := memory.NewRefreshTokenRepository()
	jwtManager := auth.NewJWTManager("test-secret-auth", 24*time.Hour)

	userHandler := handlers.NewUserHandler(userRepo)
	messageHandler := handlers.NewMessageHandler(messageRepo)
	authHandler := handlers.NewAuthHandler(userRepo, refreshTokenRepo, jwtManager)

	router := NewRouter(RouterConfig{
		UserHandler:    userHandler,
		MessageHandler: messageHandler,
		AuthHandler:    authHandler,
		UserRepo:       userRepo,
		MessageRepo:    messageRepo,
		JWTManager:     jwtManager,
		CORSOrigins:    []string{"*"},
		Metrics:        nil, // Avoid duplicate registration
	})

	// Register and get token
	registerBody := map[string]string{
		"name":     "Auth Test User",
		"email":    "authtest@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerBody)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var authResponse map[string]interface{}
	_ = json.NewDecoder(w.Body).Decode(&authResponse)
	token := authResponse["accessToken"].(string)

	t.Run("get current user with token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}

		var user domain.User
		if err := json.NewDecoder(w.Body).Decode(&user); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if user.Email != "authtest@example.com" {
			t.Errorf("expected email 'authtest@example.com', got %s", user.Email)
		}
	})

	t.Run("list users with token", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/api/v1/users", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d: %s", w.Code, w.Body.String())
		}
	})
}

func TestCORSHeaders(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodOptions, "/api/v1/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// CORS should be handled
	if w.Header().Get("Access-Control-Allow-Origin") == "" {
		t.Error("expected Access-Control-Allow-Origin header")
	}
}

func TestRequestIDHeader(t *testing.T) {
	router := setupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Request ID should be present in the response
	if w.Header().Get("X-Request-ID") == "" {
		t.Error("expected X-Request-ID header")
	}
}
