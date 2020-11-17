package model

import "github.com/gorilla/websocket"

// Envelop read buffer of the websocket
type Envelop struct {
	Type    MessageType
	Message interface{}
}

// Broadcast wrap channel message
type Broadcast struct {
	Sender  *websocket.Conn
	Message interface{}
}

// MessageType enum
type MessageType string

// RoomType room incoming message
// PingType ping incoming message
const (
	RoomType MessageType = "room"
	PingType MessageType = "ping"
	PongType MessageType = "pong"
	NoneType MessageType = "none"
)

// Ping ping/pong message
type Ping struct {
	Timecode int64  `json:"timecode"`
	RoomName string `json:"room"`
}

// PongMessage return pong timecode
type PongMessage struct {
	Timecode int64 `json:"timecode"`
}
