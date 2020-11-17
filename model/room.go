package model

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client maintain information and mutex of connection
type Client struct {
	Mu             sync.Mutex
	ConnectionTime time.Time
}

// RoomSubscription information about users in a room
type RoomSubscription struct {
	UserID  string                         `json:"userId"`
	Name    string                         `json:"name"`
	Clients map[*websocket.Conn]RoomMember `json:"-"`
}

// RoomJoinedMessage information about users in a room
type RoomJoinedMessage struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

// RoomCreatedMessage information about users in a room
type RoomCreatedMessage struct {
	Status bool   `json:"status"`
	Name   string `json:"name"`
}

// RoomMember identify a client in a room
type RoomMember struct {
	UserID string  `json:"userId"`
	Socket *Client `json:"-"`
}

// Ping ping/pong message
type Ping struct {
	Timecode int32  `json:"timecode"`
	RoomName string `json:"room"`
}
