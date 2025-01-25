package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gorilla/websocket"
)

var (
	ErrConnectionIsClose   = errors.New("connection is close")
	ErrOnClosingConnection = errors.New("error on closing connection")
)

type Service interface {
	JoinRoom() error
	LeaveRoom() error
	ReadMessageWS() ([]byte, error)
	SendMessageWS(string) error

	// get online users in room.
	GetUsers() ([]int, error)
}

type service struct {
	conn    *websocket.Conn
	wsURL   string
	httpURL string
}

func NewService(wsUrl, httpUrl string) Service {
	return &service{
		wsURL: wsUrl,
		httpURL: httpUrl,
	}
}

func (s *service) JoinRoom() error {
	conn, _, err := websocket.DefaultDialer.Dial(s.wsURL, nil)
	if err != nil {
		return fmt.Errorf("WebSocket Failed to connect: %v", err)
	}

	s.conn = conn
	return nil
}

func (s *service) LeaveRoom() error {
	if s.conn != nil {
		if err := s.conn.Close(); err != nil {
			return ErrOnClosingConnection
		}
		s.conn = nil
	}
	return nil
}

func (s *service) SendMessageWS(msg string) error {
	if s.conn != nil {
		err := s.conn.WriteMessage(websocket.TextMessage, []byte(msg))
		if err != nil {
			return fmt.Errorf("write error: %v", err)
		}
	}
	return nil
}

func (s *service) ReadMessageWS() ([]byte, error) {
	if s.conn != nil {
		_, resp, err := s.conn.ReadMessage()
		if err != nil {
			return nil, fmt.Errorf("read error: %v", err)
		}
		return resp, nil
	}
	return nil, nil
}

func (s *service) GetUsers() ([]int, error) {
	// Send a GET request
	resp, err := http.Get(s.httpURL)
	if err != nil {
		return nil, fmt.Errorf("error on sending get users request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d (%s)", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error on reading http response: %v", err)
	}

	var users []int
	if err := json.Unmarshal(body, &users); err != nil {
		return nil, fmt.Errorf("error on unmarshal response body; %v", err)
	}

	return users, nil
}
