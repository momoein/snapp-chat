package client

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"encoding/json"
)

func TestJoinRoom(t *testing.T) {
	s := NewService("ws://localhost:8080/api/v1/ws", "http://localhost:8080/api/v1ws/users")
	err := s.JoinRoom()
	if err != nil {
		t.Errorf("JoinRoom() error = %v", err)
	}
	defer s.LeaveRoom()
}

func TestSendMessage(t *testing.T) {
	s := NewService("ws://localhost:8080/api/v1/ws", "http://localhost:8080/api/v1/ws/users")
	err := s.JoinRoom()
	if err != nil {
		t.Fatalf("Failed to join room: %v", err)
	}
	defer s.LeaveRoom()

	err = s.SendMessageWS("test message")
	if err != nil {
		t.Errorf("SendMessageWS() error = %v", err)
	}
}

func TestGetUsers(t *testing.T) {
	// Mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		users := []int{1, 2, 3}
		json.NewEncoder(w).Encode(users)
	}))
	defer server.Close()

	s := NewService("ws://localhost:8080/ws", server.URL)
	
	users, err := s.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers() error = %v", err)
	}
	
	if len(users) != 3 {
		t.Errorf("Expected 3 users, got %d", len(users))
	}
}