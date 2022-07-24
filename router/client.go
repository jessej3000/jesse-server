package router

import (
	"fmt"

	"github.com/gorilla/websocket"
)

// Message structure definition
type Message struct {
	Name string      `json:"name"`
	Data interface{} `json:"data"`
}

type MessagePacket struct {
	_message Message
	_client  *Client
}

type clientDetailPacket struct {
	Client *Client
	Detail ClientDetail
}

// ClientDetail struct definition
type ClientDetail struct {
	ID      string
	Name    string
	IP      string
	Session string
}

// Client object definition
type Client struct {
	Send   chan Message
	Socket *websocket.Conn
}

// Read incoming messages
func (C *Client) Read(manager *Manager) {
	var msg Message
	var msgPacket MessagePacket
	defer func() {
		manager.Disconnect <- C
		C.Socket.Close()
	}()

	for {
		if err := C.Socket.ReadJSON(&msg); err != nil {
			fmt.Println("Error reading from client.", err.Error())
			break
		}
		msgPacket._client = C
		msgPacket._message = msg
		fmt.Println("Event (" + msg.Name + ") received...")
		manager.InComingMessage <- msgPacket
	}
}

func (C *Client) Write(manager *Manager) {
	defer func() {
		manager.Disconnect <- C
		C.Socket.Close()
	}()

	for {
		msg := <-C.Send
		if err := C.Socket.WriteJSON(msg); err != nil {
			fmt.Println("Error writing to client...")
			break
		}
	}
}
