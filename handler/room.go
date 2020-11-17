package handler

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

var rooms = make(map[string]model.RoomSubscription)
var roomCreatedBroadcast = make(chan *model.RoomCreatedMessage)
var roomJoinedBroadcast = make(chan *model.RoomJoinedMessage)

// DispatchRoomMessage dispatch the room join message
func DispatchRoomMessage() {
	for {
		select {
		case val := <-roomCreatedBroadcast:
			log.Printf("new room created %s", val.Name)
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

// ReadRoomMessages goroutine to read
func ReadRoomMessages(ws *websocket.Conn, clientInfo *model.Client) {
	go func() {
		for {
			var roomSub model.RoomSubscription
			clientInfo.Mu.Lock()
			fmt.Println("Read room")
			err := ws.ReadJSON(&roomSub)
			clientInfo.Mu.Unlock()
			if err != nil {
				log.Printf("error: %v", err)
				continue
			}
			var roomName = roomSub.Name
			if val, ok := rooms[roomName]; ok {
				notInRoom := true
				for _, info := range val.Clients {
					fmt.Printf("%s - %s\n", info.UserID, val.UserID)
					if info.UserID == roomSub.UserID {
						notInRoom = false
						break
					}
				}
				if notInRoom {
					val.Clients[ws] = model.RoomMember{
						UserID: val.UserID,
						Socket: clientInfo,
					}
					roomJoined := model.RoomJoinedMessage{
						UserID: val.UserID,
						Name:   roomName,
					}
					// Send the newly received message to the broadcast channel
					roomJoinedBroadcast <- &roomJoined
				}

			} else {
				// create new room
				roomSub.Clients = make(map[*websocket.Conn]model.RoomMember)
				roomSub.Clients[ws] = model.RoomMember{
					UserID: roomSub.UserID,
					Socket: clientInfo,
				}
				rooms[roomName] = roomSub
				message :=
					model.RoomCreatedMessage{
						Status: true,
						Name:   roomName,
					}
				roomCreatedBroadcast < &message
				log.Printf("new room created %s with %s\n", roomName, roomSub.UserID)
			}
		}
	}()
}
