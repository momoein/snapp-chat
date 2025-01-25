package server

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

const (
	Subject = "chat.messages"
)

type Service interface {

	// it add new client and
	// listen for new message from client websocket and publish the message.
	AddClient(conn *websocket.Conn)

	// Subscribe to NATS for incoming messages
	// and start broadcast messages to clients
	Run()

	// get online users in room.
	GetUsers() []int
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

// it add new client and
// listen for new message from client websocket and publish the message.
func (s *service) AddClient(conn *websocket.Conn) {
	defer conn.Close()

	client := NewClient(conn, time.Now().UTC().Nanosecond())
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

		userMsg := UserMessage{
			UserID: client.ID(),
			Message: string(msg),
		}

		userMsgByte, err := json.Marshal(userMsg)
		if err != nil {
			log.Printf("error on marshal user message: %v", err)
			continue
		}
		
		// Publish client message to NATS
		err = s.nc.Publish(Subject, userMsgByte)
		if err != nil {
			log.Printf("Error publishing to NATS: %v", err)
		}
	}
}

func (s *service) GetUsers() []int {
	var users []int
	for client := range s.hub.clients {
		users = append(users, client.ID())
	}
	return users
}
