package handlers

import (
	"net/http"

	"github.com/google/uuid"
	"mangahub/internal/middleware"
	"mangahub/internal/repositories"
)

type NotificationHandler struct {
	notifs *repositories.NotificationRepository
}

func NewNotificationHandler(notifs *repositories.NotificationRepository) *NotificationHandler {
	return &NotificationHandler{notifs: notifs}
}

// GET /api/notifications
func (h *NotificationHandler) List(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	unreadOnly := r.URL.Query().Get("unread") == "true"

	notifs, err := h.notifs.GetByUser(r.Context(), uid, unreadOnly)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to fetch notifications")
		return
	}
	count, _ := h.notifs.CountUnread(r.Context(), uid)
	respondOK(w, map[string]interface{}{
		"notifications": notifs,
		"unread_count":  count,
	})
}

// PUT /api/notifications/{id}/read
func (h *NotificationHandler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)

	// Extract notification id from path using chi
	idStr := r.PathValue("id")
	if idStr == "" {
		respondError(w, http.StatusBadRequest, "missing notification ID")
		return
	}
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid notification ID")
		return
	}
	if err := h.notifs.MarkRead(r.Context(), id, uid); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark read")
		return
	}
	respondOK(w, map[string]string{"message": "marked as read"})
}

// PUT /api/notifications/read-all
func (h *NotificationHandler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	uid, _ := uuid.Parse(userID)
	if err := h.notifs.MarkAllRead(r.Context(), uid); err != nil {
		respondError(w, http.StatusInternalServerError, "failed to mark all read")
		return
	}
	respondOK(w, map[string]string{"message": "all marked as read"})
}
