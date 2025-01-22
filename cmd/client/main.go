package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"snappchat/api/cli"
	config "snappchat/config/client"
	"snappchat/internal/client"

	"github.com/gorilla/websocket"
)

var configPath = flag.String("config", "client_config.json", "configuration file")

func main() {
	flag.Parse()

	if v := os.Getenv("CONFIG_PATH"); len(v) > 0 {
		*configPath = v
	}
	cfg := config.MustReadConfig(*configPath)

	serverURL := url.URL{
		Scheme: "ws", 
		Host: fmt.Sprintf("%s:%d", cfg.WebsocketAddr.Host, cfg.WebsocketAddr.Port), 
		Path: cfg.WebsocketAddr.Path,
	}

	conn, _, err := websocket.DefaultDialer.Dial(serverURL.String(), nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer conn.Close()

	app := client.NewClientApp(conn)

	cli.Run(app)
}
