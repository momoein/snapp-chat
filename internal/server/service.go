package server

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

const (
	Subject = "chat.messages"
)

type Service interface {
	
	// read message from client websocket and publish message.
	AddClient(conn *websocket.Conn)

	// Subscribe to NATS for incoming messages
	// and start broadcast messages to clients
	Run()
}

type service struct {
	nc  *nats.Conn
	hub *Hub
}

func NewService(conn *nats.Conn) Service {
	return &service{
		nc:  conn,
		hub: NewHub(),
	}
}

// Subscribe to NATS for incoming messages
// and start broadcast messages to clients
func (s *service) Run() {
	_, err := s.nc.Subscribe(Subject, func(m *nats.Msg) {
		log.Printf("Message received from NATS: %s", string(m.Data))
		s.hub.broadcast <- m.Data // Forward the message to WebSocket clients
	})
	if err != nil {
		log.Printf("error subscribing to nats: %v\n", err)
		return
	}
	defer s.nc.Close()

	s.hub.startBroadcast()
}

// read message from client websocket and publish message.
func (s *service) AddClient(conn *websocket.Conn) {
	defer conn.Close()

	client := &Client{conn: conn}
	s.hub.AddClient(client)
	log.Println("New WebSocket client connected")

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading from WebSocket: %v", err)
			s.hub.RemoveClient(client)
			break
		}
		log.Printf("Message received from client: %s", string(msg))

		// Publish client message to NATS
		err = s.nc.Publish(Subject, msg)
		if err != nil {
			log.Printf("Error publishing to NATS: %v", err)
		}
	}
}
