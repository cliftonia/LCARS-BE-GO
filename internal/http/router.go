package http

import (
	"net/http"

	"github.com/cliftonbaggerman/subspace-backend/internal/http/handlers"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/middleware"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// RouterConfig holds dependencies for creating routes
type RouterConfig struct {
	UserHandler    *handlers.UserHandler
	MessageHandler *handlers.MessageHandler
	CORSOrigins    []string
}

// NewRouter creates a new HTTP router with all routes configured
func NewRouter(config RouterConfig) http.Handler {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recovery)

	// API version prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	r.HandleFunc("/health", healthCheck).Methods(http.MethodGet)

	// User routes
	api.HandleFunc("/users", config.UserHandler.ListUsers).Methods(http.MethodGet)
	api.HandleFunc("/users", config.UserHandler.CreateUser).Methods(http.MethodPost)
	api.HandleFunc("/users/{id}", config.UserHandler.GetUser).Methods(http.MethodGet)
	api.HandleFunc("/users/{id}", config.UserHandler.UpdateUser).Methods(http.MethodPut)
	api.HandleFunc("/users/{id}", config.UserHandler.DeleteUser).Methods(http.MethodDelete)

	// Message routes
	api.HandleFunc("/messages", config.MessageHandler.CreateMessage).Methods(http.MethodPost)
	api.HandleFunc("/messages/{id}", config.MessageHandler.GetMessage).Methods(http.MethodGet)
	api.HandleFunc("/messages/{id}", config.MessageHandler.DeleteMessage).Methods(http.MethodDelete)
	api.HandleFunc("/messages/{id}/read", config.MessageHandler.MarkAsRead).Methods(http.MethodPatch)

	// User-specific message routes
	api.HandleFunc("/users/{userId}/messages", config.MessageHandler.GetUserMessages).Methods(http.MethodGet)
	api.HandleFunc("/users/{userId}/messages/unread-count", config.MessageHandler.GetUnreadCount).Methods(http.MethodGet)

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

// healthCheck returns a simple health status
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy","service":"subspace-backend"}`))
}
