package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"time"

	"github.com/davecgh/go-spew/spew"
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

// ErrorBroadcast public error broadcast
var ErrorBroadcast = make(chan *model.Broadcast)

// DispatchErrorMessage dispatch error message to the socket
func DispatchErrorMessage() {
	for {
		event := <-ErrorBroadcast
		val := event.Message.(model.ErrorMessage)
		envelop := model.Envelop{
			Type:    model.ErrorType,
			Message: val,
		}

		// send to broadcast sender
		event.Sender.WriteJSON(envelop)
	}
}

// WSHandler handle room subscription
func WSHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	// register client
	newClient := model.Socket{
		ConnectionTime: time.Now(),
		WebSocket:      ws,
	}
	clients[ws] = &newClient
	fmt.Printf("new connection (%d) %v+\n", len(clients), ws.RemoteAddr())

	// Reader loop
	go func() {
		for {
			var msg json.RawMessage
			env := model.Envelop{
				Message: &msg,
			}
			err := ws.ReadJSON(&env)
			spew.Dump(env)
			if err != nil {
				log.Printf("read error: %s", reflect.TypeOf(err).String())
				broadcast := model.Broadcast{
					Sender: ws,
					Message: model.RoomEventMessage{
						Action: model.RoomLeft,
						UserID: "someone",
					},
				}
				roomLeftBroadcast <- &broadcast
				ws.Close()
				break
			}
			fmt.Printf("message %s from %v+\n", env.Type, ws.RemoteAddr())
			switch env.Type {
			case model.PingType:
				ReadPingMessage(env.Message.(*json.RawMessage), &newClient)
			case model.RoomType:
				ReadRoomMessage(env.Message.(*json.RawMessage), &newClient)
			}
		}
	}()

}
