package realtime

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 8192
)

// Event is a typed WebSocket message
type Event struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
	RoomID  string      `json:"room_id,omitempty"`
	UserID  string      `json:"user_id,omitempty"`
}

// Client represents a single WebSocket connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	roomID string
}

// Hub manages all active WebSocket connections
type Hub struct {
	mu      sync.RWMutex
	clients map[*Client]bool
	rooms   map[string]map[*Client]bool // roomID -> clients
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Client]bool),
		rooms:   make(map[string]map[*Client]bool),
	}
}

// Register adds a client to the hub
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[client] = true
	log.Printf("[WS] client connected: user=%s room=%s total=%d", client.userID, client.roomID, len(h.clients))
}

// Unregister removes a client from the hub
func (h *Hub) Unregister(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		close(client.send)
		if client.roomID != "" {
			if room, ok := h.rooms[client.roomID]; ok {
				delete(room, client)
			}
		}
	}
	log.Printf("[WS] client disconnected: user=%s total=%d", client.userID, len(h.clients))
}

// JoinRoom adds client to a room
func (h *Hub) JoinRoom(client *Client, roomID string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	// Leave old room
	if client.roomID != "" {
		if room, ok := h.rooms[client.roomID]; ok {
			delete(room, client)
		}
	}
	client.roomID = roomID
	if _, ok := h.rooms[roomID]; !ok {
		h.rooms[roomID] = make(map[*Client]bool)
	}
	h.rooms[roomID][client] = true
}

// BroadcastToRoom sends an event to all clients in a room
func (h *Hub) BroadcastToRoom(roomID string, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	if room, ok := h.rooms[roomID]; ok {
		for client := range room {
			select {
			case client.send <- data:
			default:
				// Slow client; skip
			}
		}
	}
}

// BroadcastToUser sends an event to a specific user (all their connections)
func (h *Hub) BroadcastToUser(userID string, event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		if client.userID == userID {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}

// BroadcastAll sends to every connected client
func (h *Hub) BroadcastAll(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for client := range h.clients {
		select {
		case client.send <- data:
		default:
		}
	}
}

// ConnectedUsers returns count of connected users
func (h *Hub) ConnectedUsers() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// NewClient creates and runs a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	client := &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
	hub.Register(client)
	go client.readPump()
	go client.writePump()
	return client
}

func (c *Client) readPump() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// Handle incoming client events (join room, etc.)
		var evt Event
		if err := json.Unmarshal(message, &evt); err == nil {
			switch evt.Type {
			case "join_room":
				if roomID, ok := evt.Payload.(string); ok {
					c.hub.JoinRoom(c, roomID)
				} else if m, ok := evt.Payload.(map[string]interface{}); ok {
					if rid, ok := m["room_id"].(string); ok {
						c.hub.JoinRoom(c, rid)
					}
				}
			case "ping":
				data, _ := json.Marshal(Event{Type: "pong", Payload: time.Now()})
				c.send <- data
			}
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)
			// Drain queued messages
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte("\n"))
				w.Write(<-c.send)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
