package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

// Rooms list to broadcast message to
var Rooms = make(map[string]model.RoomEventMessage)
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
			envelop := model.Envelop{
				Type:    model.RoomType,
				Message: val,
			}

			// send to every client that is currently connected
			for ws, client := range Rooms[val.Name].Clients {
				client.Socket.Mu.Lock()
				err := ws.WriteJSON(envelop)
				client.Socket.Mu.Unlock()
				if err != nil {
					log.Printf("Websocket error: %s", err)
					ws.Close()
				}
			}
		case event := <-roomJoinedBroadcast:
			val := event.Message.(model.RoomEventMessage)
			log.Printf("%s joined the room %v+", val.UserID, val)
			envelop := model.Envelop{
				Type:    model.RoomType,
				Message: val,
			}

			// send to every client that is currently connected
			for ws, client := range Rooms[val.Name].Clients {
				client.Socket.Mu.Lock()
				err := ws.WriteJSON(envelop)
				client.Socket.Mu.Unlock()
				if err != nil {
					log.Printf("Websocket error: %s", err)
					ws.Close()
				}
			}
		case event := <-roomLeftBroadcast:
			val := event.Message.(model.RoomEventMessage)
			log.Println("left room message received")

			for roomname, room := range Rooms {
				for ws, client := range room.Clients {
					if ws == event.Sender {
						val.UserID = client.UserID
						envelop := model.Envelop{
							Type:    model.RoomType,
							Message: val,
						}

						// that was one of its room, lets notify everyone
						for ws2 := range room.Clients {
							client.Socket.Mu.Lock()
							err := ws2.WriteJSON(envelop)
							client.Socket.Mu.Unlock()
							if err != nil {
								log.Printf("Websocket error: %s", err)
								ws2.Close()
							}
						}
						delete(room.Clients, event.Sender)
						// do not break in case multiple people connected to the same socket
						// break
					}
				}
				// if room is empty, lets delete it to avoid memory leaks
				if len(room.Clients) == 0 {
					delete(Rooms, roomname)
				}
			}
		}
	}
}

// ReadRoomMessage goroutine to read
func ReadRoomMessage(msg *json.RawMessage, clientInfo *model.Socket) {
	var roomSub model.RoomEventMessage
	err := json.Unmarshal(*msg, &roomSub)
	if err != nil {
		log.Printf("room error: %v [%s]", err, msg)
		return
	}
	broadcast := model.Broadcast{
		Sender: clientInfo.WebSocket,
	}
	switch roomSub.Action {
	case model.RoomJoin:
		err = join(roomSub, &broadcast, clientInfo)
	case model.RoomLeave:
		err = leave(roomSub, &broadcast, clientInfo)
	case model.RoomCreate:
		err = create(roomSub, &broadcast, clientInfo)
	}
	if err != nil {
		broadcast.Message = model.ErrorMessage{
			Message: err.Error(),
		}
		ErrorBroadcast <- &broadcast
	}
}

func create(r model.RoomEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if _, ok := Rooms[r.Name]; !ok {
		Rooms[r.Name] = model.RoomEventMessage{
			Name:    r.Name,
			UserID:  r.UserID,
			Clients: make(map[*websocket.Conn]model.RoomMember),
		}
		Rooms[r.Name].Clients[broadcast.Sender] = model.RoomMember{
			UserID: r.UserID,
			Socket: client,
			Host:   true,
		}
		broadcast.Message =
			model.RoomEventMessage{
				Action: model.RoomCreated,
				Name:   r.Name,
				UserID: r.UserID,
			}
		roomCreatedBroadcast <- broadcast
		return nil
	}
	return errors.New("room_already_exists")

}

func join(r model.RoomEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		notInRoom := true
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				notInRoom = false
				break
			}
		}
		if notInRoom {
			val.Clients[broadcast.Sender] = model.RoomMember{
				UserID: r.UserID,
				Socket: client,
				Host:   false,
			}
			broadcast.Message = model.RoomEventMessage{
				UserID: r.UserID,
				Name:   r.Name,
				Action: model.RoomJoined,
			}
			roomJoinedBroadcast <- broadcast
			return nil
		}
		return errors.New("already_in_room")
	}
	return errors.New("room_does_not_exists")
}

func leave(r model.RoomEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				val.Clients[broadcast.Sender] = model.RoomMember{
					UserID: r.UserID,
					Socket: client,
					Host:   false,
				}
				broadcast.Message = model.RoomEventMessage{
					Action: model.RoomLeft,
					UserID: r.UserID,
				}
				roomLeftBroadcast <- broadcast
				return nil
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}
