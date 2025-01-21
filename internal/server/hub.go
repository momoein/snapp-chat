package server

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	lock sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		broadcast: make(chan []byte),
		clients:   make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		message := <-h.broadcast
		h.lock.Lock()
		for client := range h.clients {
			err := client.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Printf("error writing to client: %v", err)
				client.conn.Close()
				delete(h.clients, client)
			}
		}
		h.lock.Unlock()
	}
}

func (h *Hub) AddClient(client *Client) {
	h.lock.Lock()
	h.clients[client] = true
	h.lock.Unlock()
}

func (h *Hub) RemoveClient(client *Client) {
	h.lock.Lock()
	delete(h.clients, client)
	h.lock.Unlock()
}
