package model

import "github.com/gorilla/websocket"

// Envelop read buffer of the websocket
type Envelop struct {
	Type    MessageType `json:"type"`
	Message interface{} `json:"message"`
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
	RoomType  MessageType = "room"
	PingType  MessageType = "ping"
	PongType  MessageType = "pong"
	NoneType  MessageType = "none"
	ErrorType MessageType = "error"
)

// Ping ping/pong message
type PingMessage struct {
	Timecode int64  `json:"timecode"`
	RoomName string `json:"room"`
	UserID   string `json:"userId"`
}

// PongMessage return pong timecode
type PongMessage struct {
	Timecode int64 `json:"timecode"`
}

// ErrorMessage ship error message to client
type ErrorMessage struct {
	Message string `json:"message"`
}
