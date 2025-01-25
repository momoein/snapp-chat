package service

import (
	"snappchat/internal/server"

	"github.com/gorilla/websocket"
)

type ServerService struct {
	svc server.Service
}

func NewServerService(svc server.Service) *ServerService {
	return &ServerService{svc: svc}
}

func (s *ServerService) AddClient(conn *websocket.Conn) {
	s.svc.AddClient(conn)
}

func (s *ServerService) Run() {
	s.svc.Run()
}

func (s *ServerService) GetUsers() []int {
	return s.svc.GetUsers()
}
