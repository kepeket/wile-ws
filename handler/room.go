package handler

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

var rooms = make(map[string]model.RoomSubscription)
var roomCreatedBroadcast = make(chan *model.Broadcast)
var roomJoinedBroadcast = make(chan *model.Broadcast)
var roomLeftBroadcast = make(chan *model.Broadcast)

// DispatchRoomMessage dispatch the room join message
func DispatchRoomMessage() {
	for {
		select {
		case event := <-roomCreatedBroadcast:
			val := event.Message.(model.RoomEventMessage)
			log.Printf("new room created %s", val.Name)
			var msg []byte
			envelop := model.Envelop{
				Type:    model.RoomType,
				Message: val,
			}
			msg, err := json.Marshal(envelop)
			if err != nil {
				log.Fatalf("Unable to prepare message %v+", err)
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
		case event := <-roomJoinedBroadcast:
			val := event.Message.(model.RoomEventMessage)
			log.Printf("%s joined the room %v+", val.UserID, val)
			var msg []byte
			envelop := model.Envelop{
				Type:    model.RoomType,
				Message: val,
			}
			msg, err := json.Marshal(envelop)
			if err != nil {
				log.Fatalf("Unable to prepare message %v+", err)
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
		case event := <-roomLeftBroadcast:
			val := event.Message.(model.RoomEventMessage)
			log.Println("left room message received")
			var msg []byte
			delete(rooms[val.Name].Clients, event.Sender)
			envelop := model.Envelop{
				Type:    model.RoomType,
				Message: val,
			}
			msg, err := json.Marshal(envelop)
			if err != nil {
				log.Fatalf("Unable to prepare message %v+", err)
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
		}
	}
}

// ReadRoomMessage goroutine to read
func ReadRoomMessage(msg *json.RawMessage, clientInfo *model.Socket) {
	var roomSub model.RoomSubscription
	err := json.Unmarshal(*msg, &roomSub)
	if err != nil {
		log.Printf("room error: %v [%s]", err, msg)
		return
	}
	broadcast := model.Broadcast{
		Sender: clientInfo.WebSocket,
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
				UserID: roomSub.UserID,
				Socket: clientInfo,
				Host:   false,
			}
			broadcast.Message = model.RoomEventMessage{
				UserID: roomSub.UserID,
				Name:   roomName,
				Action: model.RoomJoined,
			}
			// Send new joiner message
			roomJoinedBroadcast <- &broadcast
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
		broadcast.Message =
			model.RoomEventMessage{
				Action: model.RoomCreated,
				Name:   roomName,
				UserID: roomSub.UserID,
			}
		roomCreatedBroadcast <- &broadcast
		log.Printf("new room created %s by %s\n", roomName, roomSub.UserID)
	}
}
