package http

import (
	"fmt"
	"log"
	"net/http"
	"snappchat/api/http/handlers"
	"snappchat/api/service"
	app "snappchat/app/server"
)

func Run(app app.App) error {
	s := service.NewServerService(app.Service())

	api := "/api/v1"
	registerAPI(api, s)

	go s.Run()

	// Start the HTTP server
	httpPort := app.Config().Server.HttpPort
	log.Printf("Server started on http://localhost:%d", httpPort)
	return  http.ListenAndServe(fmt.Sprintf(":%d", httpPort), nil)
}

func registerAPI(api string, svc *service.ServerService) {
	handler := handlers.NewServerHandler(svc)

	http.HandleFunc("/", handler.Home)

	// WebSocket handler
	http.HandleFunc(api + "/ws", handler.AddClient)
	http.HandleFunc(api + "/ws/users", handler.GetUsers)
}
