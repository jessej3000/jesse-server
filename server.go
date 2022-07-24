package main

import (
	"fmt"
	"net/http"
	"server/controller"
	"server/router"

	"github.com/gorilla/websocket"
)

type Server struct {
	Host string
	Port string
}

// var manager *router.Manager

// Creates and initialize server instance
// Params
// host : string - host address
// port : string
// manager : *router.Manager)
// Returns
// *Server : pointer to Server object
// *router.Router : pointer to Router
func NewServer(host string, port string) (*Server, *router.Router) {

	messageRouter := router.NewMessageRouter()
	messageRouter.Handle("msg_in", controller.MessageIn)
	messageRouter.Handle("identify", controller.Identify)

	return &Server{
		Host: host,
		Port: port,
	}, messageRouter
}

// ListenAndServe
// func (S *Server) ListenAndServe(ctx context.Context) {
func (S *Server) ListenAndServe() {
	http.HandleFunc("/", upgradeToWebsocket)
	// s := &http.Server{
	// 	Addr:           S.Host + ":" + S.Port,
	// 	Handler:        nil,
	// 	ReadTimeout:    30000,
	// 	WriteTimeout:   30000,
	// 	MaxHeaderBytes: 1 << 20,
	// }
	// s.ListenAndServe()

	http.ListenAndServe(S.Host+":"+S.Port, nil)
}

// Setup response header
func setupResponse(w *http.ResponseWriter, r *http.Request) {
	// (*w).Header().Set("Access-Control-Allow-Origin", "*") // "https://virtualjesse.herokuapp.com")
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:9000")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

// Handles http request and upgrade to websocket
func upgradeToWebsocket(w http.ResponseWriter, r *http.Request) {
	// setupResponse(&w, r)
	fmt.Println("Conntection found")
	socket, error := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(w, r, nil)
	if error != nil {
		fmt.Println("Upgrade to Websocket Error", error)
		return
	}

	client := &router.Client{Send: make(chan router.Message), Socket: socket}

	manager.Connect <- client

	go client.Read(&manager)
	go client.Write(&manager)

}
