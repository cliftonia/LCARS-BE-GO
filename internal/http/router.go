// Package http provides HTTP server routing and configuration.
package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/auth"
	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/handlers"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// RouterConfig holds dependencies for creating routes
type RouterConfig struct {
	UserHandler    *handlers.UserHandler
	MessageHandler *handlers.MessageHandler
	AuthHandler    *handlers.AuthHandler
	UserRepo       domain.UserRepository
	MessageRepo    domain.MessageRepository
	JWTManager     *auth.JWTManager
	CORSOrigins    []string
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(config RouterConfig) http.Handler {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)
	r.Use(middleware.RateLimit(100, 10)) // 100 requests/minute, burst of 10

	// API version prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	r.HandleFunc("/health", healthCheckHandler(config.UserRepo, config.MessageRepo)).Methods(http.MethodGet)

	// Auth routes (public)
	api.HandleFunc("/auth/register", config.AuthHandler.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", config.AuthHandler.Login).Methods(http.MethodPost)
	api.HandleFunc("/auth/apple", config.AuthHandler.AppleSignIn).Methods(http.MethodPost)
	api.HandleFunc("/auth/google", config.AuthHandler.GoogleSignIn).Methods(http.MethodPost)
	api.HandleFunc("/auth/refresh", config.AuthHandler.Refresh).Methods(http.MethodPost)

	// Protected auth routes
	authAPI := api.PathPrefix("/auth").Subrouter()
	authAPI.Use(middleware.Auth(config.JWTManager))
	authAPI.HandleFunc("/me", config.AuthHandler.Me).Methods(http.MethodGet)

	// User routes (protected)
	usersAPI := api.PathPrefix("/users").Subrouter()
	usersAPI.Use(middleware.Auth(config.JWTManager))
	usersAPI.HandleFunc("", config.UserHandler.ListUsers).Methods(http.MethodGet)
	usersAPI.HandleFunc("", config.UserHandler.CreateUser).Methods(http.MethodPost)
	usersAPI.HandleFunc("/{id}", config.UserHandler.GetUser).Methods(http.MethodGet)
	usersAPI.HandleFunc("/{id}", config.UserHandler.UpdateUser).Methods(http.MethodPut)
	usersAPI.HandleFunc("/{id}", config.UserHandler.DeleteUser).Methods(http.MethodDelete)

	// Message routes (protected)
	messagesAPI := api.PathPrefix("/messages").Subrouter()
	messagesAPI.Use(middleware.Auth(config.JWTManager))
	messagesAPI.HandleFunc("", config.MessageHandler.CreateMessage).Methods(http.MethodPost)
	messagesAPI.HandleFunc("/{id}", config.MessageHandler.GetMessage).Methods(http.MethodGet)
	messagesAPI.HandleFunc("/{id}", config.MessageHandler.DeleteMessage).Methods(http.MethodDelete)
	messagesAPI.HandleFunc("/{id}/read", config.MessageHandler.MarkAsRead).Methods(http.MethodPatch)

	// User-specific message routes (protected)
	usersAPI.HandleFunc("/{userId}/messages", config.MessageHandler.GetUserMessages).Methods(http.MethodGet)
	usersAPI.HandleFunc("/{userId}/messages/unread-count", config.MessageHandler.GetUnreadCount).Methods(http.MethodGet)

	// Configure CORS
	c := cors.New(cors.Options{
		AllowedOrigins:   config.CORSOrigins,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	})

	return c.Handler(r)
}

// HealthStatus represents the health status response
type HealthStatus struct {
	Status    string            `json:"status"`
	Service   string            `json:"service"`
	Timestamp time.Time         `json:"timestamp"`
	Checks    map[string]string `json:"checks"`
}

// healthCheckHandler returns a health check handler that checks dependencies
func healthCheckHandler(userRepo domain.UserRepository, messageRepo domain.MessageRepository) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		checks := make(map[string]string)
		allHealthy := true

		// Check user repository
		if _, err := userRepo.Count(); err != nil {
			checks["user_repository"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks["user_repository"] = "healthy"
		}

		// Check message repository
		if _, err := messageRepo.GetUnreadCount("health-check"); err != nil {
			checks["message_repository"] = "unhealthy: " + err.Error()
			allHealthy = false
		} else {
			checks["message_repository"] = "healthy"
		}

		status := "healthy"
		statusCode := http.StatusOK
		if !allHealthy {
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
		}

		response := HealthStatus{
			Status:    status,
			Service:   "subspace-backend",
			Timestamp: time.Now(),
			Checks:    checks,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		_ = json.NewEncoder(w).Encode(response)
	}
}
