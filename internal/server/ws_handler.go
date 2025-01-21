package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/nats-io/nats.go"
)

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

type Server struct {
	nc *nats.Conn

	hub *Hub
}

func NewServer(nc *nats.Conn, hub *Hub) *Server {
	return &Server{
		nc:  nc,
		hub: hub,
	}
}

func (s *Server) Run() error {
	defer s.nc.Close()

	go s.hub.run()

	// Subscribe to NATS for incoming messages
	subject := "chat.messages"
	_, err := s.nc.Subscribe(subject, func(m *nats.Msg) {
		log.Printf("Message received from NATS: %s", string(m.Data))
		s.hub.broadcast <- m.Data // Forward the message to WebSocket clients
	})
	if err != nil {
		return fmt.Errorf("error subscribing to nats: %v", err)
	}

	// WebSocket handler
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		s.ServeWs(subject, w, r)
	})

	// Start the HTTP server
	log.Println("Server started on http://localhost:8080")
	return fmt.Errorf("error starting server: %v", http.ListenAndServe(":8080", nil))
}

// serveWs handles websocket requests from the peer.
func (s *Server) ServeWs(subject string, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := &Client{hub: s.hub, conn: conn}
	s.hub.AddClient(client)
	log.Println("New WebSocket client connected")

	// Handle client messages
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from WebSocket: %v", err)
			s.hub.RemoveClient(client)
			break
		}
		log.Printf("Message received from client: %s", string(msg))

		// Publish client message to NATS
		err = s.nc.Publish(subject, msg)
		if err != nil {
			log.Printf("Error publishing to NATS: %v", err)
		}
	}
}
