package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/gorilla/mux"
)

// UserHandler handles user-related HTTP requests
type UserHandler struct {
	repo domain.UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo domain.UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// GetUser handles GET /api/v1/users/{id}
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.repo.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// ListUsers handles GET /api/v1/users
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	limitStr := r.URL.Query().Get("limit")
	offsetStr := r.URL.Query().Get("offset")

	limit := 20 // default
	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}

	offset := 0 // default
	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	// Validate pagination
	if err := validatePagination(limit, offset); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	users, err := h.repo.List(limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve users")
		return
	}

	total, err := h.repo.Count()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve user count")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   users,
		"total":  total,
		"limit":  limit,
		"offset": offset,
	})
}

// CreateUser handles POST /api/v1/users
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate user input
	if err := validateUserName(user.Name); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateEmail(user.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.repo.Create(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create user")
		return
	}

	respondWithJSON(w, http.StatusCreated, user)
}

// UpdateUser handles PUT /api/v1/users/{id}
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	// Validate user input
	if err := validateUserName(user.Name); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := validateEmail(user.Email); err != nil {
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	user.ID = id

	if err := h.repo.Update(&user); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	respondWithJSON(w, http.StatusOK, user)
}

// DeleteUser handles DELETE /api/v1/users/{id}
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
