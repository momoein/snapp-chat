package main

import (
	"log"
	"net/url"
	"snappchat/api/handlers/cli"
	"snappchat/internal/client"

	"github.com/gorilla/websocket"
)

func main() {
	serverURL := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/ws"}

	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	app := client.NewClientApp(conn)

	cli.Run(app)
}
