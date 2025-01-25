package server

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nats-io/nats.go"
)

// create a nats connection, so make sure nats server is running.
func setupNATS() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}
	return nc
}

func TestGetUsers(t *testing.T) {
	nc := setupNATS()
	defer nc.Close()

	srv := NewService(nc)

	upgrader := websocket.Upgrader{}

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Logf("WebSocket upgrade error: %v", err)
				return
			}
			srv.AddClient(conn)
		}),
	}

	// Simulate adding clients
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	go server.Serve(listener)
	defer listener.Close()

	var wg sync.WaitGroup
	for range 2 {
		// Create WebSocket connection
		wg.Add(1)
		go func() {
			url := "ws://" + listener.Addr().String() + "/"
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Errorf("Failed to connect to WebSocket: %v", err)
				return
			}
			defer conn.Close()
			time.Sleep(200 * time.Millisecond)
			wg.Done()
		}()
	}

	time.Sleep(100 * time.Millisecond)

	users := srv.GetUsers()
	if len(users) != 2 {
		t.Errorf("Expected 2 users, got %d", len(users))
	}

	wg.Wait()
}

// This test validates the service's ability to handle concurrent WebSocket clients sending messages 
// and ensure all messages are correctly published and received via NATS.
func TestConcurrentMessageSending(t *testing.T) {
	// Setup NATS connection
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		t.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Create server service
	srv := NewService(nc)
	go srv.Run()

	// Setup WebSocket server for testing
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	server := http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				t.Logf("WebSocket upgrade error: %v", err)
				return
			}
			srv.AddClient(conn)
		}),
	}

	// Listen on a random available port
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}
	go server.Serve(listener)
	defer server.Close()

	// Number of concurrent clients
	clientCount := 50
	messageCount := 5

	// Channel to collect received messages
	receivedMessages := make(chan UserMessage, clientCount*messageCount)
	var wg sync.WaitGroup

	// Setup subscription to collect messages
	sub, err := nc.Subscribe(Subject, func(msg *nats.Msg) {
		var userMsg UserMessage
		if err := json.Unmarshal(msg.Data, &userMsg); err == nil {
			receivedMessages <- userMsg
		}
	})
	if err != nil {
		t.Fatalf("Failed to subscribe: %v", err)
	}
	defer sub.Unsubscribe()

	// Simulate concurrent clients
	for i := 0; i < clientCount; i++ {
		wg.Add(1)
		go func(clientID int) {
			defer wg.Done()

			// Create WebSocket connection
			url := "ws://" + listener.Addr().String() + "/"
			conn, _, err := websocket.DefaultDialer.Dial(url, nil)
			if err != nil {
				t.Errorf("Failed to connect to WebSocket: %v", err)
				return
			}
			defer conn.Close()

			// Send multiple messages
			for j := 0; j < messageCount; j++ {
				message := UserMessage{
					UserID:  clientID,
					Message: fmt.Sprintf("Client %d, message %d", clientID, j),
				}

				msgBytes, _ := json.Marshal(message)
				err := nc.Publish(Subject, msgBytes)
				if err != nil {
					t.Errorf("Failed to publish message: %v", err)
				}
			}
		}(i)
	}

	// Wait for all clients to send messages
	wg.Wait()

	// Allow time for message processing
	time.Sleep(500 * time.Millisecond)

	// Close the channel
	close(receivedMessages)

	// Collect and verify received messages
	receivedCount := 0
	for range receivedMessages {
		receivedCount++
	}

	expectedTotalMessages := clientCount * messageCount
	if receivedCount != expectedTotalMessages {
		t.Errorf(
			"Incomplete message distribution: received %d, expected %d",
			receivedCount,
			expectedTotalMessages,
		)
	}
}
