package handler

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/wile-ws/model"
)

var pingBroadcast = make(chan *model.Broadcast)

// DispatchPingMessage dispatch the ping messages
func DispatchPingMessage() {
	for {
		event := <-pingBroadcast
		val := event.Message.(model.PingMessage)

		origEnvelop := model.Envelop{
			Type: model.PongType,
			Message: model.PongMessage{
				Timecode: time.Now().UnixNano(),
			},
		}
		memberEnvelop := model.Envelop{
			Type:    model.PingType,
			Message: val,
		}

		// send to every client that is currently connected
		for ws, client := range rooms[val.RoomName].Clients {
			client.Socket.Mu.Lock()
			var msg []byte
			var err error
			if ws == event.Sender {
				msg, err = json.Marshal(origEnvelop)
				if err != nil {
					log.Fatalf("Unable to send message %v+", err)
					continue
				}
			} else {
				msg, err = json.Marshal(memberEnvelop)
				if err != nil {
					log.Fatalf("Unable to send message %v+", err)
					continue
				}
			}
			err = ws.WriteMessage(websocket.TextMessage, msg)
			client.Socket.Mu.Unlock()
			if err != nil {
				log.Printf("Websocket error: %s", err)
				ws.Close()
			}
		}
	}
}

// ReadPingMessage goroutine to read pings
func ReadPingMessage(msg *json.RawMessage, clientInfo *model.Socket) {
	broadcast := model.Broadcast{
		Sender: clientInfo.WebSocket,
	}
	var ping model.PingMessage
	err := json.Unmarshal(*msg, &ping)
	if err != nil {
		log.Printf("ping error: %v %s", err, string(*msg))
		return
	}
	// Ensure you're part of a room
	if room, ok := rooms[ping.RoomName]; ok {
		if _, present := room.Clients[clientInfo.WebSocket]; !present {
			log.Printf("%s not in room %v+", clientInfo.WebSocket.RemoteAddr(), room)
			return
		}
	} else {
		log.Printf("room %s doest no exist, %v", ping.RoomName, rooms)
	}
	broadcast.Message = ping
	pingBroadcast <- &broadcast
}

/*
Ping
when you send a ping message with a timecode
- you expect to a pong from the server
- members of the room also receive your ping

Goal
- detect lag between client and server
- detect timespan bitween members (and possibly ajust clock)
*/
