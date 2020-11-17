package model

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// RoomSubscription information about users in a room
type RoomSubscription struct {
	UserID  string                         `json:"userId"`
	Name    string                         `json:"name"`
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

// RoomJoined
// RoomCreated
// RoomLeft
const (
	RoomJoined  RoomActionType = "joined"
	RoomCreated RoomActionType = "created"
	RoomLeft    RoomActionType = "left"
)

// RoomEventMessage information about users in a room
type RoomEventMessage struct {
	UserID string
	Name   string
	Action RoomActionType
}
