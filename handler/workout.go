package handler

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/wile-ws/model"
)

var workoutBroadcast = make(chan *model.Broadcast)

// var workoutLobbyBroadcast = make(chan *model.Broadcast)
// var workoutStartedBroadcast = make(chan *model.Broadcast)
// var workoutStoppedBroadcast = make(chan *model.Broadcast)
// var workoutTrainingStartBroadcast = make(chan *model.Broadcast)

// DispatchWorkoutMessage dispatch workout message to room
func DispatchWorkoutMessage() {
	for {
		select {
		// case event := <-workoutLobbyBroadcast:
		// case event := <-workoutStartedBroadcast:
		// case event := <-workoutStoppedBroadcast:
		// case event := <-workoutTrainingStartBroadcast:
		case event := <-workoutBroadcast:
			val := event.Message.(model.WorkoutEventMessage)
			log.Printf("workout about to begin in room %s", val.Name)
			envelop := model.Envelop{
				Type:    model.WorkoutType,
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
		}
	}
}

// ReadWorkoutMessage goroutine to read workout
func ReadWorkoutMessage(msg *json.RawMessage, clientInfo *model.Socket) {
	var event model.WorkoutEventMessage
	err := json.Unmarshal(*msg, &event)
	if err != nil {
		log.Printf("room error: %v [%s]", err, msg)
		return
	}
	broadcast := model.Broadcast{
		Sender: clientInfo.WebSocket,
	}
	switch event.Action {
	case model.WorkoutTick:
		err = tick(event, &broadcast, clientInfo)
	case model.WorkoutLobby:
		err = lobby(event, &broadcast, clientInfo)
	case model.WorkoutReady:
		err = ready(event, &broadcast, clientInfo)
	case model.WorkoutStart:
		err = start(event, &broadcast, clientInfo)
	case model.WorkoutStop:
		err = stop(event, &broadcast, clientInfo)
	case model.WorkoutTrainingStart:
		err = trainingStart(event, &broadcast, clientInfo)
	case model.WorkoutPause:
		err = pause(event, &broadcast, clientInfo)
	}
	if err != nil {
		broadcast.Message = model.ErrorMessage{
			Message: err.Error(),
		}
		ErrorBroadcast <- &broadcast
	}
}

func tick(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				broadcast.Message = r
				workoutBroadcast <- broadcast
				return nil
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}

func lobby(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				if info.Host {
					broadcast.Message = r
					workoutBroadcast <- broadcast
					return nil
				}
				return errors.New("not_owner_of_room")
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}

func pause(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				if info.Host {
					r.Action = model.WorkoutPaused
					broadcast.Message = r
					workoutBroadcast <- broadcast
					return nil
				}
				return errors.New("not_owner_of_room")
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}

func ready(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				broadcast.Message = r
				workoutBroadcast <- broadcast
				return nil
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")

}

func start(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				if info.Host {
					r.Action = model.WorkoutStarted
					broadcast.Message = r
					workoutBroadcast <- broadcast
					return nil
				}
				return errors.New("not_owner_of_room")
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}

func stop(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				if info.Host == true {
					r.Action = model.WorkoutStopped
					broadcast.Message = r
					workoutBroadcast <- broadcast
					return nil
				}
				return errors.New("not_owner_of_room")
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}

func trainingStart(r model.WorkoutEventMessage, broadcast *model.Broadcast, client *model.Socket) error {
	if val, ok := Rooms[r.Name]; ok {
		for _, info := range val.Clients {
			if info.UserID == r.UserID {
				if info.Host == true {
					r.Action = model.WorkoutTrainingStart
					broadcast.Message = r
					workoutBroadcast <- broadcast
					return nil
				}
				return errors.New("not_owner_of_room")
			}
		}
		return errors.New("not_in_room")
	}
	return errors.New("room_does_not_exists")
}
