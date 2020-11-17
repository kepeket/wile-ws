package handler

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

var pingBroadcast = make(chan *model.Ping)

// DispatchPingMessage dispatch the ping messages
func DispatchPingMessage() {
	for {
		val := <-pingBroadcast
		log.Println("ping into room")
		mess := fmt.Sprintf("someone has joined the room")
		// send to every client that is currently connected
		for ws, client := range rooms[val.RoomName].Clients {
			client.Socket.Mu.Lock()
			err := ws.WriteMessage(websocket.TextMessage, []byte(mess))
			client.Socket.Mu.Unlock()
			if err != nil {
				log.Printf("Websocket error: %s", err)
				ws.Close()
			}
		}
	}
}

// ReadPingMessages goroutine to read pings
func ReadPingMessages(ws *websocket.Conn, clientInfo *model.Client) {
	go func() {
		for {
			var ping model.Ping
			clientInfo.Mu.Lock()
			fmt.Println("Read ping")
			err := ws.ReadJSON(&ping)
			clientInfo.Mu.Unlock()
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}

			if room, ok := rooms[ping.RoomName]; ok {
				if member, present := room.Clients[ws]; !present {
					log.Printf("%s not in room %v+", member.UserID, room)
				}
			}
			pingBroadcast <- &ping
		}
	}()
}
