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
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	jwtManager       *auth.JWTManager
	appleVerifier    *auth.AppleAuthVerifier
	googleVerifier   *auth.GoogleAuthVerifier
}

// NewAuthHandler creates a new auth handler
func NewAuthHandler(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	jwtManager *auth.JWTManager,
) *AuthHandler {
	return &AuthHandler{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtManager:       jwtManager,
		appleVerifier:    auth.NewAppleAuthVerifier(),
		googleVerifier:   auth.NewGoogleAuthVerifier(),
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
	AccessToken  string       `json:"accessToken"`
	RefreshToken string       `json:"refreshToken"`
	User         *domain.User `json:"user"`
}

// AppleSignInRequest represents an Apple Sign In request
type AppleSignInRequest struct {
	UserID            string                 `json:"userId"`
	IdentityToken     string                 `json:"identityToken"`
	AuthorizationCode string                 `json:"authorizationCode"`
	Email             string                 `json:"email"`
	FullName          map[string]interface{} `json:"fullName"`
}

// GoogleSignInRequest represents a Google Sign In request
type GoogleSignInRequest struct {
	UserID      string `json:"userId"`
	IDToken     string `json:"idToken"`
	AccessToken string `json:"accessToken"`
	Email       string `json:"email"`
	FullName    string `json:"fullName"`
}

// RefreshTokenRequest represents a token refresh request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken"`
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

	// Generate access and refresh tokens
	accessToken, refreshToken, err := h.generateTokenPair(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Return tokens and user (password already excluded from JSON)
	respondWithJSON(w, http.StatusCreated, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
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

	// Generate access and refresh tokens
	accessToken, refreshToken, err := h.generateTokenPair(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	// Return tokens and user
	respondWithJSON(w, http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
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

// Refresh handles POST /api/v1/auth/refresh
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	// Validate refresh token
	if req.RefreshToken == "" {
		respondWithError(w, http.StatusBadRequest, "refresh token is required")
		return
	}

	// Get refresh token from database
	storedToken, err := h.refreshTokenRepo.GetByToken(req.RefreshToken)
	if err != nil {
		if err == domain.ErrRefreshTokenNotFound {
			respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to verify refresh token")
		return
	}

	// Validate token
	if !storedToken.IsValid() {
		if storedToken.IsExpired() {
			respondWithError(w, http.StatusUnauthorized, "refresh token has expired")
			return
		}
		if storedToken.IsRevoked() {
			respondWithError(w, http.StatusUnauthorized, "refresh token has been revoked")
			return
		}
		respondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	// Get user
	user, err := h.userRepo.GetByID(storedToken.UserID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	// Revoke old refresh token
	if err := h.refreshTokenRepo.Revoke(req.RefreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to revoke old token")
		return
	}

	// Generate new token pair
	accessToken, newRefreshToken, err := h.generateTokenPair(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User:         user,
	})
}

// AppleSignIn handles POST /api/v1/auth/apple
func (h *AuthHandler) AppleSignIn(w http.ResponseWriter, r *http.Request) {
	var req AppleSignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	// Validate required fields
	if req.IdentityToken == "" {
		respondWithError(w, http.StatusBadRequest, "identity token is required")
		return
	}

	// For development/testing: Allow mock tokens
	if req.IdentityToken == "mock-token" || req.IdentityToken == "mock-id-token" {
		// Create or get user for mock authentication
		user, err := h.getOrCreateAppleUser(req.UserID, req.Email, req.FullName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		// Generate access and refresh tokens
		accessToken, refreshToken, err := h.generateTokenPair(user.ID, user.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		respondWithJSON(w, http.StatusOK, AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		})
		return
	}

	// Production: Verify Apple identity token
	// Note: You'll need to configure your Apple app's client ID
	clientID := "com.subspace.app" // TODO: Move to config
	claims, err := h.appleVerifier.VerifyIdentityToken(req.IdentityToken, clientID)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Apple identity token")
		return
	}

	// Use verified email from token if not provided in request
	email := req.Email
	if email == "" {
		email = claims.Email
	}

	// Get or create user
	user, err := h.getOrCreateAppleUser(claims.Sub, email, req.FullName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := h.generateTokenPair(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// getOrCreateAppleUser gets an existing user or creates a new one for Apple Sign In
func (h *AuthHandler) getOrCreateAppleUser(_, email string, fullName map[string]interface{}) (*domain.User, error) {
	// Try to find existing user by email
	existingUser, err := h.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	// Build user name from fullName
	name := email // Default to email
	if fullName != nil {
		var nameParts []string
		if givenName, ok := fullName["givenName"].(string); ok && givenName != "" {
			nameParts = append(nameParts, givenName)
		}
		if familyName, ok := fullName["familyName"].(string); ok && familyName != "" {
			nameParts = append(nameParts, familyName)
		}
		if len(nameParts) > 0 {
			name = nameParts[0]
			if len(nameParts) > 1 {
				name = nameParts[0] + " " + nameParts[1]
			}
		}
	}

	// Create new user
	user := &domain.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  "", // No password for OAuth users
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// GoogleSignIn handles POST /api/v1/auth/google
func (h *AuthHandler) GoogleSignIn(w http.ResponseWriter, r *http.Request) {
	var req GoogleSignInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	// Validate required fields
	if req.IDToken == "" {
		respondWithError(w, http.StatusBadRequest, "id token is required")
		return
	}

	// For development/testing: Allow mock tokens
	if req.IDToken == "mock-id-token" || req.IDToken == "mock-token" {
		// Create or get user for mock authentication
		user, err := h.getOrCreateGoogleUser(req.Email, req.FullName)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to create user")
			return
		}

		// Generate access and refresh tokens
		accessToken, refreshToken, err := h.generateTokenPair(user.ID, user.Email)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
			return
		}

		respondWithJSON(w, http.StatusOK, AuthResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			User:         user,
		})
		return
	}

	// Production: Verify Google ID token
	tokenInfo, err := h.googleVerifier.VerifyIDToken(req.IDToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid Google ID token")
		return
	}

	// Use verified email and name from token
	email := tokenInfo.Email
	fullName := tokenInfo.Name
	if fullName == "" && req.FullName != "" {
		fullName = req.FullName
	}

	// Get or create user
	user, err := h.getOrCreateGoogleUser(email, fullName)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	// Generate access and refresh tokens
	accessToken, refreshToken, err := h.generateTokenPair(user.ID, user.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate tokens")
		return
	}

	respondWithJSON(w, http.StatusOK, AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         user,
	})
}

// getOrCreateGoogleUser gets an existing user or creates a new one for Google Sign In
func (h *AuthHandler) getOrCreateGoogleUser(email, fullName string) (*domain.User, error) {
	// Try to find existing user by email
	existingUser, err := h.userRepo.GetByEmail(email)
	if err == nil && existingUser != nil {
		return existingUser, nil
	}

	// Use email as name if fullName is empty
	name := fullName
	if name == "" {
		name = email
	}

	// Create new user
	user := &domain.User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		Password:  "", // No password for OAuth users
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// generateTokenPair creates both access and refresh tokens for a user
func (h *AuthHandler) generateTokenPair(userID, email string) (accessToken, refreshToken string, err error) {
	// Generate access token (JWT)
	accessToken, err = h.jwtManager.GenerateToken(userID, email)
	if err != nil {
		return "", "", err
	}

	// Generate refresh token (random string)
	refreshToken, err = auth.GenerateRefreshToken()
	if err != nil {
		return "", "", err
	}

	// Store refresh token in database
	refreshTokenModel := &domain.RefreshToken{
		ID:        uuid.New().String(),
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour), // 7 days
		CreatedAt: time.Now(),
	}

	if err := h.refreshTokenRepo.Create(refreshTokenModel); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
