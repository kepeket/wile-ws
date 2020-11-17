package model

// Training struct to unmarshall a training sent by room host
type Training struct {
	Name         string `json:"name"`
	Duration     int    `json:"duration"`
	Reps         string `json:"reps"`
	RepRate      string `json:"repRate"`
	TrainingType string `json:"trainingType"`
}

// TrainingTicker when changing excercise send this to sync devices
type TrainingTicker struct {
	Training Training `json:"training"`
	Timecode int32    `json:"timecode"`
}
