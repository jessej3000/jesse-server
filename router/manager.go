package router

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// Manager object definition
type Manager struct {
	Clients         map[*Client]ClientDetail
	InComingMessage chan MessagePacket
	Connect         chan *Client
	Disconnect      chan *Client
	ModifyDetail    chan clientDetailPacket
}

// Start starts Manager
func (C *Manager) Start(r *Router) {
	fmt.Println("Chat server started...")
	for {
		select {
		case connectedClient := <-C.Connect:
			// Check if admin
			//		Check if domain payment ok
			// If client
			//		Check if client payment ok
			uuid, _ := GUID()
			remoteAddressAsID := connectedClient.Socket.RemoteAddr().String()
			// if strings.Contains(remoteAddressAsID, "127.0.0.1") {
			C.Clients[connectedClient] = ClientDetail{
				ID:      "",
				Name:    "",
				IP:      remoteAddressAsID,
				Session: uuid,
			}
			fmt.Println("New Client Connected:", remoteAddressAsID)
			// } else {
			// 	fmt.Println("Invalid connection attempt from: ", remoteAddressAsID)
			// 	connectedClient.Socket.Close()
			// }
		case disconnectedClient := <-C.Disconnect:
			if detail, ok := C.Clients[disconnectedClient]; ok {
				fmt.Println(
					detail.ID,
					detail.Name, "disconnected...")
				close(disconnectedClient.Send)
				delete(C.Clients, disconnectedClient)
			}
		case inComingMessage := <-C.InComingMessage:
			fmt.Println(
				C.Clients[inComingMessage._client].ID,
				C.Clients[inComingMessage._client].Name,
				"sent a message...")
			if Handler, found := r.FindHandler(inComingMessage._message.Name); found {
				fmt.Println(inComingMessage._message.Name, "Handler found")
				go Handler(inComingMessage._client, inComingMessage._message.Data, C)
			} else {
				fmt.Println(inComingMessage._message.Name, "Handler not found")
			}
		}
	}
}

// MonitorDetails monitor any changes
func (C *Manager) MonitorDetails() {
	for {
		detailsUpdate := <-C.ModifyDetail
		C.Clients[detailsUpdate.Client] = detailsUpdate.Detail
	}
}

// GUID generates GUID
func GUID() (string, string) {
	id := uuid.New()
	uid := strings.Replace(id.String(), "-", "", -1)
	return uid, uid[0:4]
}
