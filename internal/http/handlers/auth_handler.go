// Package handlers provides HTTP request handlers for the Subspace Backend API.
package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/cliftonbaggerman/subspace-backend/internal/auth"
	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/cliftonbaggerman/subspace-backend/internal/http/middleware"
	"github.com/google/uuid"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	userRepo   domain.UserRepository
	jwtManager *auth.JWTManager
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(userRepo domain.UserRepository, jwtManager *auth.JWTManager) *AuthHandler {
	return &AuthHandler{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

// RegisterRequest represents a user registration request
type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents a user login request
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents an authentication response
type AuthResponse struct {
	Token string       `json:"token"`
	User  *domain.User `json:"user"`
}

// Register handles POST /api/v1/auth/register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	// Validate input
	if err := validateUserName(req.Name); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateEmail(req.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if len(req.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "password must be at least 8 characters")
		return
	}

	// Check if email already exists
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err == nil && existingUser != nil {
		respondWithError(w, http.StatusConflict, "email already registered")
		return
	}

	// Hash password
	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to process registration")
		return
	}

	// Create user
	user := &domain.User{
		ID:        uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.Create(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token and user (password already excluded from JSON)
	respondWithJSON(w, http.StatusCreated, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Login handles POST /api/v1/auth/login
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	// Validate input
	if err := validateEmail(req.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.Password == "" {
		respondWithError(w, http.StatusBadRequest, "password is required")
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Verify password
	if err := auth.VerifyPassword(user.Password, req.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "invalid email or password")
		return
	}

	// Generate JWT token
	token, err := h.jwtManager.GenerateToken(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Return token and user
	respondWithJSON(w, http.StatusOK, AuthResponse{
		Token: token,
		User:  user,
	})
}

// Me handles GET /api/v1/auth/me
func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := middleware.GetUserIDFromContext(r.Context())
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "User ID not found in context")
		return
	}

	user, err := h.userRepo.GetByID(userID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}
