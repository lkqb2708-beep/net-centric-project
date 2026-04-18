package udp

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

// Notification is the UDP payload structure
type Notification struct {
	Type    string      `json:"type"`
	Title   string      `json:"title"`
	Body    string      `json:"body"`
	Payload interface{} `json:"payload,omitempty"`
	Time    time.Time   `json:"time"`
}

// Subscriber holds a UDP client's address and subscription topics
type Subscriber struct {
	Addr   *net.UDPAddr
	Topics map[string]bool
}

// Server handles UDP notification broadcasting
type Server struct {
	addr        string
	conn        *net.UDPConn
	mu          sync.RWMutex
	subscribers map[string]*Subscriber // key: addr.String()
}

func NewServer(port string) *Server {
	return &Server{
		addr:        ":" + port,
		subscribers: make(map[string]*Subscriber),
	}
}

func (s *Server) Start() error {
	addr, err := net.ResolveUDPAddr("udp", s.addr)
	if err != nil {
		return fmt.Errorf("UDP resolve addr failed: %w", err)
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return fmt.Errorf("UDP listen failed: %w", err)
	}
	s.conn = conn
	log.Printf("[UDP] listening on %s", s.addr)

	go s.readLoop()
	return nil
}

func (s *Server) readLoop() {
	buf := make([]byte, 4096)
	for {
		n, remoteAddr, err := s.conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("[UDP] read error: %v", err)
			continue
		}

		var msg map[string]interface{}
		if err := json.Unmarshal(buf[:n], &msg); err != nil {
			// Send error back
			s.sendTo(remoteAddr, Notification{
				Type:  "error",
				Title: "Invalid payload",
				Body:  "Expected JSON",
				Time:  time.Now(),
			})
			continue
		}

		msgType, _ := msg["type"].(string)
		switch msgType {
		case "subscribe":
			topics := []string{"chapter_release", "friend_activity", "system"}
			if t, ok := msg["topics"].([]interface{}); ok {
				topics = make([]string, 0, len(t))
				for _, v := range t {
					if s, ok := v.(string); ok {
						topics = append(topics, s)
					}
				}
			}
			s.mu.Lock()
			sub := &Subscriber{Addr: remoteAddr, Topics: make(map[string]bool)}
			for _, topic := range topics {
				sub.Topics[topic] = true
			}
			s.subscribers[remoteAddr.String()] = sub
			s.mu.Unlock()
			s.sendTo(remoteAddr, Notification{
				Type:  "subscribed",
				Title: "Subscribed",
				Body:  fmt.Sprintf("Subscribed to %v", topics),
				Time:  time.Now(),
			})
			log.Printf("[UDP] subscriber added: %s topics=%v", remoteAddr, topics)

		case "unsubscribe":
			s.mu.Lock()
			delete(s.subscribers, remoteAddr.String())
			s.mu.Unlock()
			s.sendTo(remoteAddr, Notification{
				Type: "unsubscribed", Title: "Unsubscribed", Time: time.Now(),
			})

		case "ping":
			s.sendTo(remoteAddr, Notification{
				Type: "pong", Title: "Pong", Time: time.Now(),
			})

		default:
			s.sendTo(remoteAddr, Notification{
				Type:  "error",
				Title: "Unknown type",
				Body:  "Unknown message type: " + msgType,
				Time:  time.Now(),
			})
		}
	}
}

func (s *Server) sendTo(addr *net.UDPAddr, n Notification) {
	data, err := json.Marshal(n)
	if err != nil {
		return
	}
	s.conn.WriteToUDP(data, addr)
}

// Broadcast sends a notification to all subscribers of a topic
func (s *Server) Broadcast(topic string, n Notification) {
	n.Time = time.Now()
	data, err := json.Marshal(n)
	if err != nil {
		return
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sub := range s.subscribers {
		if sub.Topics[topic] || sub.Topics["all"] {
			s.conn.WriteToUDP(data, sub.Addr)
		}
	}
	log.Printf("[UDP] broadcast topic=%s to %d subscribers", topic, len(s.subscribers))
}

// BroadcastChapterRelease notifies all subscribers of a new chapter
func (s *Server) BroadcastChapterRelease(mangaTitle string, chapter int) {
	s.Broadcast("chapter_release", Notification{
		Type:  "chapter_release",
		Title: "New Chapter Available!",
		Body:  fmt.Sprintf("%s - Chapter %d is now available", mangaTitle, chapter),
		Payload: map[string]interface{}{
			"manga_title": mangaTitle,
			"chapter":     chapter,
		},
	})
}

func (s *Server) SubscriberCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}
