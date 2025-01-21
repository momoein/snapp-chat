package main

import (
	"log"
	"snappchat/internal/server"
	// gorillachat "snappchat/sample/gorilla-caht"

	"github.com/nats-io/nats.go"
)

func main() {
	// gorillachat.Run()

	var (
		NATSServer = "nats://localhost:4222" // NATS server address
	)

	nc, err := nats.Connect(NATSServer)
	if err != nil {
		log.Fatal("Error on connecting to nats: ", err)
	}

	hub := server.NewHub()

	s := server.NewServer(nc, hub)

	log.Fatal(s.Run())
}
