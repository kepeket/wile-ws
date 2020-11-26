package model

// WorkoutEventMessage information about users in a room
type WorkoutEventMessage struct {
	UserID           string            `json:"userId"`
	Name             string            `json:"name"`
	Action           WorkoutActionType `json:"action"`
	Training         Training          `json:"training"`
	Timestamp        int64             `json:"timestamp"`
	Countdown        int               `json:"countdown"`
	TrainingCount    int               `json:"trainingCount"`
	TrainingPosition int               `json:"trainingPos"`
}

// WorkoutActionType enum
type WorkoutActionType string

// WorkoutLobby wainting for devices to sync
// WorkoutStart receive start signal
// WorkoutStarted confirm workout started
// WorkoutStop end of workout
// WorkoutStopped confirm end of workout
// WorkoutTrainingStart new training
const (
	WorkoutTick          WorkoutActionType = "tick"
	WorkoutLobby         WorkoutActionType = "lobby"
	WorkoutReady         WorkoutActionType = "ready"
	WorkoutStart         WorkoutActionType = "start"
	WorkoutStarted       WorkoutActionType = "started"
	WorkoutPause         WorkoutActionType = "pause"
	WorkoutPaused        WorkoutActionType = "paused"
	WorkoutStop          WorkoutActionType = "stop"
	WorkoutStopped       WorkoutActionType = "stopped"
	WorkoutTrainingStart WorkoutActionType = "new_training"
)

// Training struct to unmarshall a training sent by room host
type Training struct {
	Name         string `json:"name"`
	Reps         int    `json:"reps"`
	RepRate      int    `json:"repRate"`
	TrainingType string `json:"trainingType"`
}
