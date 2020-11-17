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
var clients = make(map[*websocket.Conn]*model.Client)

// WSHandler handle room subscription
func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("new connection")

	// register client
	newClient := model.Client{
		ConnectionTime: time.Now(),
	}
	clients[ws] = &newClient

	ReadRoomMessages(ws, &newClient)
	ReadPingMessages(ws, &newClient)

}
