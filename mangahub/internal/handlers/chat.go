package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"mangahub/internal/middleware"
	"mangahub/internal/models"
	"mangahub/internal/realtime"
	"mangahub/internal/repositories"
)

type ChatHandler struct {
	chat *repositories.ChatRepository
	hub  *realtime.Hub
}

func NewChatHandler(chat *repositories.ChatRepository, hub *realtime.Hub) *ChatHandler {
	return &ChatHandler{chat: chat, hub: hub}
}

// GET /api/chat/rooms
func (h *ChatHandler) GetRooms(w http.ResponseWriter, r *http.Request) {
	rooms, err := h.chat.GetRooms(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch chat rooms")
		return
	}
	if rooms == nil {
		rooms = []*models.ChatRoom{}
	}
	respondOK(w, rooms)
}

// GET /api/chat/rooms/{id}/messages
func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room ID")
		return
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	if limit < 1 || limit > 100 {
		limit = 50
	}

	msgs, err := h.chat.GetMessages(r.Context(), roomID, limit, offset)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch messages")
		return
	}
	if msgs == nil {
		msgs = []*models.ChatMessage{}
	}
	respondOK(w, msgs)
}

// POST /api/chat/rooms/{id}/messages
func (h *ChatHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	roomID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room ID")
		return
	}

	var req models.SendMessageRequest
	if err := decodeJSON(r, &req); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request")
		return
	}
	if req.Content == "" {
		respondError(w, http.StatusBadRequest, "content is required")
		return
	}

	msg := &models.ChatMessage{
		RoomID:    roomID,
		UserID:    uid,
		Content:   req.Content,
		CreatedAt: time.Now(),
	}
	if err := h.chat.SaveMessage(r.Context(), msg); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to save message")
		return
	}

	// Broadcast via WebSocket
	h.hub.BroadcastToRoom(roomID.String(), realtime.Event{
		Type: "chat_message",
		Payload: map[string]interface{}{
			"id":         msg.ID,
			"room_id":    roomID,
			"user_id":    uid,
			"username":   middleware.GetUsername(r),
			"content":    msg.Content,
			"created_at": msg.CreatedAt,
		},
	})

	respondCreated(w, msg)
}
