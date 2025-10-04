package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cliftonbaggerman/subspace-backend/internal/domain"
	"github.com/gorilla/mux"
)

// MessageHandler handles message-related HTTP requests
type MessageHandler struct {
	repo domain.MessageRepository
}

// NewMessageHandler creates a new message handler
func NewMessageHandler(repo domain.MessageRepository) *MessageHandler {
	return &MessageHandler{repo: repo}
}

// GetMessage handles GET /api/v1/messages/{id}
func (h *MessageHandler) GetMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	message, err := h.repo.GetByID(id)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Message not found")
		return
	}

	respondWithJSON(w, http.StatusOK, message)
}

// GetUserMessages handles GET /api/v1/users/{userId}/messages
func (h *MessageHandler) GetUserMessages(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

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

	messages, err := h.repo.GetByUserID(userID, limit, offset)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve messages")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]interface{}{
		"data":   messages,
		"limit":  limit,
		"offset": offset,
	})
}

// GetUnreadCount handles GET /api/v1/users/{userId}/messages/unread-count
func (h *MessageHandler) GetUnreadCount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userId"]

	count, err := h.repo.GetUnreadCount(userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to retrieve unread count")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]int{"unreadCount": count})
}

// CreateMessage handles POST /api/v1/messages
func (h *MessageHandler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	var message domain.Message
	if err := json.NewDecoder(r.Body).Decode(&message); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := h.repo.Create(&message); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create message")
		return
	}

	respondWithJSON(w, http.StatusCreated, message)
}

// MarkAsRead handles PATCH /api/v1/messages/{id}/read
func (h *MessageHandler) MarkAsRead(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.MarkAsRead(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to mark message as read")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Message marked as read"})
}

// DeleteMessage handles DELETE /api/v1/messages/{id}
func (h *MessageHandler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	if err := h.repo.Delete(id); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to delete message")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "Message deleted successfully"})
}
