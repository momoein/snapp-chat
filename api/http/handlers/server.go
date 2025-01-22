package handlers

import (
	"log"
	"net/http"
	"snappchat/api/service"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins (you can customize this for security)
	},
}

type ServerHandler struct {
	svc *service.ServerService
}

func NewServerHandler(s *service.ServerService) *ServerHandler {
	return &ServerHandler{
		svc: s,
	}
}

func (h *ServerHandler) Home(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}


func (h *ServerHandler) wsUpgrade(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (h *ServerHandler) AddClient(w http.ResponseWriter, r *http.Request) {
	conn, err := h.wsUpgrade(w, r)
	if err != nil {
		log.Println(err)
		return
	}
	h.svc.AddClient(conn)
}
