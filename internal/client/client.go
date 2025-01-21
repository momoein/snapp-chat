package client

import (
	"fmt"
	"log"
	"net/url"

	"github.com/gorilla/websocket"
)

type ClientApp struct {
	conn *websocket.Conn
}

func NewClientApp(conn *websocket.Conn) *ClientApp {
	return &ClientApp{
		conn: conn,
	}
}

func ClientRun(serverURL url.URL) {
	log.Printf("Connecting to %s", serverURL.String())
	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()
}

func (a *ClientApp) SendMessageWS(msg string) error {
	err := a.conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		return fmt.Errorf("write error: %v", err)
	}
	return nil
}

func (a *ClientApp) ReadMessageWS() ([]byte, error) {
	_, resp, err := a.conn.ReadMessage()
	if err != nil {
		return nil, fmt.Errorf("read error: %v", err)
	}
	return resp, nil
}
