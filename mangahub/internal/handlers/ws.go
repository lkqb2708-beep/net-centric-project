package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"mangahub/internal/middleware"
	"mangahub/internal/realtime"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (restrict in production)
	},
}

type WSHandler struct {
	hub *realtime.Hub
}

func NewWSHandler(hub *realtime.Hub) *WSHandler {
	return &WSHandler{hub: hub}
}

// GET /ws  — upgrades to WebSocket
func (h *WSHandler) Handle(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("[WS] upgrade failed: %v", err)
		return
	}

	client := realtime.NewClient(h.hub, conn, userID)
	_ = client

	// Send welcome event
	h.hub.BroadcastToUser(userID, realtime.Event{
		Type: "connected",
		Payload: map[string]interface{}{
			"user_id": userID,
			"message": "Welcome to MangaHub!",
		},
	})
}
