package server

import (
	"github.com/gorilla/websocket"
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn
	id   int
}

func NewClient(conn *websocket.Conn, id int) *Client {
	return &Client{
		conn: conn,
		id:   id,
	}
}

func (c *Client) ID() int {
	return c.id
}

type UserMessage struct {
	UserID  int    `json:"userId"`
	Message string `json:"message"`
}
