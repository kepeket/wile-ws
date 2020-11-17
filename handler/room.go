package handler

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

var rooms = make(map[string]model.RoomSubscription)
var roomCreatedBroadcast = make(chan *model.RoomEventMessage)
var roomJoinedBroadcast = make(chan *model.RoomEventMessage)
var roomLeftBroadcast = make(chan *model.RoomEventMessage)

// DispatchRoomMessage dispatch the room join message
func DispatchRoomMessage() {
	for {
		select {
		case val := <-roomCreatedBroadcast:
			log.Printf("new room created %s", val.Name)
			var msg []byte
			envelop := model.Envelop{
				Type:    model.RoomType,
				Message: val,
			}
			msg, err := json.Marshal(envelop)
			if err != nil {
				log.Fatalf("Unable to send message %v+", err)
				continue
			}

			// send to every client that is currently connected
			for ws, client := range rooms[val.Name].Clients {
				client.Socket.Mu.Lock()
				err := ws.WriteMessage(websocket.TextMessage, msg)
				client.Socket.Mu.Unlock()
				if err != nil {
					log.Printf("Websocket error: %s", err)
					ws.Close()
				}
			}
		case val := <-roomJoinedBroadcast:
			log.Println("new room message received")
			mess := fmt.Sprintf("someone has joined the room")
			// send to every client that is currently connected
			for ws, client := range rooms[val.Name].Clients {
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
}

// ReadRoomMessage goroutine to read
func ReadRoomMessage(msg *[]byte, clientInfo *model.Socket) {
	var roomSub model.RoomSubscription
	err := json.Unmarshal(*msg, &roomSub)
	if err != nil {
		log.Printf("room error: %v [%s]", err, msg)
		return
	}
	var roomName = roomSub.Name
	if val, ok := rooms[roomName]; ok {
		notInRoom := true
		for _, info := range val.Clients {
			if info.UserID == roomSub.UserID {
				notInRoom = false
				break
			}
		}
		if notInRoom {
			val.Clients[clientInfo.WebSocket] = model.RoomMember{
				UserID: val.UserID,
				Socket: clientInfo,
				Host:   false,
			}
			event := model.RoomEventMessage{
				UserID: val.UserID,
				Name:   roomName,
				Action: model.RoomJoined,
			}
			// Send new joiner message
			roomJoinedBroadcast <- &event
		}

	} else {
		// create new room
		roomSub.Clients = make(map[*websocket.Conn]model.RoomMember)
		roomSub.Clients[clientInfo.WebSocket] = model.RoomMember{
			UserID: roomSub.UserID,
			Socket: clientInfo,
			Host:   true,
		}
		rooms[roomName] = roomSub
		event :=
			model.RoomEventMessage{
				Action: model.RoomCreated,
				Name:   roomName,
				UserID: roomSub.UserID,
			}
		roomCreatedBroadcast <- &event
		log.Printf("new room created %s by %s\n", roomName, roomSub.UserID)
	}
}
