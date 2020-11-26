package model

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RoomSubscription information about users in a room
type RoomEventMessage struct {
	UserID  string                         `json:"userId"`
	Name    string                         `json:"name"`
	Action  RoomActionType                 `json:"action"`
	Clients map[*websocket.Conn]RoomMember `json:"-"`
}

// RoomMember identify a client in a room
type RoomMember struct {
	UserID string
	Socket *Socket
	Host   bool
}

// Socket maintain information and mutex of connection
type Socket struct {
	Mu             sync.Mutex
	ConnectionTime time.Time
	WebSocket      *websocket.Conn
}

// RoomActionType enum
type RoomActionType string

// RoomJoin ask to join
// RoomJoined notify someone joined
// RoomCreate to create a room
// RoomCreated notify it was created
// RoomLeft notify someone left
// RoomLeave ask to leave
const (
	RoomJoin    RoomActionType = "join"
	RoomJoined  RoomActionType = "joined"
	RoomCreate  RoomActionType = "create"
	RoomCreated RoomActionType = "created"
	RoomLeft    RoomActionType = "left"
	RoomLeave   RoomActionType = "leave"
)
