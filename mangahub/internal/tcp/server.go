package tcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Message is the framing unit sent over TCP
type Message struct {
	Type    string      `json:"type"`
	UserID  string      `json:"user_id,omitempty"`
	Payload interface{} `json:"payload"`
	Time    time.Time   `json:"time"`
}

// Server handles TCP connections for reading-progress sync
type Server struct {
	addr    string
	mu      sync.RWMutex
	clients map[net.Conn]string // conn -> userID
}

func NewServer(port string) *Server {
	return &Server{
		addr:    ":" + port,
		clients: make(map[net.Conn]string),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return fmt.Errorf("TCP listen failed: %w", err)
	}
	log.Printf("[TCP] listening on %s", s.addr)

	go func() {
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("[TCP] accept error: %v", err)
				continue
			}
			go s.handleConn(conn)
		}
	}()
	return nil
}

func (s *Server) handleConn(conn net.Conn) {
	defer func() {
		s.mu.Lock()
		delete(s.clients, conn)
		s.mu.Unlock()
		conn.Close()
		log.Printf("[TCP] client disconnected: %s", conn.RemoteAddr())
	}()

	log.Printf("[TCP] client connected: %s", conn.RemoteAddr())
	s.mu.Lock()
	s.clients[conn] = ""
	s.mu.Unlock()

	// Send welcome
	s.sendTo(conn, Message{Type: "welcome", Payload: "MangaHub TCP Sync Server", Time: time.Now()})

	scanner := bufio.NewScanner(conn)
	scanner.Buffer(make([]byte, 65536), 65536)

	for scanner.Scan() {
		line := scanner.Bytes()
		var msg Message
		if err := json.Unmarshal(line, &msg); err != nil {
			s.sendTo(conn, Message{Type: "error", Payload: "invalid JSON"})
			continue
		}

		switch msg.Type {
		case "identify":
			if uid, ok := msg.Payload.(string); ok {
				s.mu.Lock()
				s.clients[conn] = uid
				s.mu.Unlock()
				s.sendTo(conn, Message{Type: "identified", Payload: uid, Time: time.Now()})
				log.Printf("[TCP] user %s identified from %s", uid, conn.RemoteAddr())
			}

		case "progress_update":
			// Broadcast progress to all other clients
			userID := s.getClientUser(conn)
			broadcast := Message{
				Type:   "progress_broadcast",
				UserID: userID,
				Payload: msg.Payload,
				Time:   time.Now(),
			}
			s.BroadcastExcept(conn, broadcast)
			s.sendTo(conn, Message{Type: "ack", Payload: "progress synced", Time: time.Now()})

		case "ping":
			s.sendTo(conn, Message{Type: "pong", Payload: time.Now(), Time: time.Now()})

		default:
			s.sendTo(conn, Message{Type: "error", Payload: "unknown message type: " + msg.Type})
		}
	}
}

func (s *Server) sendTo(conn net.Conn, msg Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		return
	}
	data = append(data, '\n')
	conn.SetWriteDeadline(time.Now().Add(5 * time.Second))
	conn.Write(data)
}

func (s *Server) BroadcastExcept(exclude net.Conn, msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for conn := range s.clients {
		if conn != exclude {
			s.sendTo(conn, msg)
		}
	}
}

func (s *Server) Broadcast(msg Message) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for conn := range s.clients {
		s.sendTo(conn, msg)
	}
}

func (s *Server) ConnectedCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.clients)
}

func (s *Server) getClientUser(conn net.Conn) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.clients[conn]
}
