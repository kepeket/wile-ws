package handler

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Clients list of subscribed clients
var clients = make(map[*websocket.Conn]*model.Socket)

// WSHandler handle room subscription
func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new connection")

	// register client
	newClient := model.Socket{
		ConnectionTime: time.Now(),
		WebSocket:      ws,
	}
	clients[ws] = &newClient

	envelop := model.Envelop{}

	// Reader loop
	go func() {
		for {
			err := ws.ReadJSON(&envelop)
			if err != nil {
				log.Printf("read error: %v", err)
				ws.Close()

				continue
			}
			switch envelop.Type {
			case model.PingType:
				ReadPingMessage(envelop.Message, &newClient)
			case model.RoomType:
				ReadRoomMessage(envelop.Message, &newClient)
			}
		}
	}()

}
