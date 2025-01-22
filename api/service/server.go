package service

import (
	"log"
	"net/http"
	"snappchat/internal/server"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (you can customize this for security)
	},
}


type ServerService struct {
	svc server.Service
}

func NewServerService(svc server.Service) *ServerService {
	return &ServerService{svc: svc}
}

func (s *ServerService) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	s.svc.AddClient(conn)
}

func (s *ServerService) AddClient(conn *websocket.Conn) {
	s.svc.AddClient(conn)
}


func (s *ServerService) Run() {
	s.svc.Run()
}
