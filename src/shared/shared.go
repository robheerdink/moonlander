package shared

import "time"

// common used variables in components, values often set per level
var (
	LP LevelProperties
)

// LevelProperties are used to store info / progress off the level
type LevelProperties struct {
	Width        int
	Height       int
	Gravity      float64
	Friction     float64
	PlayerStartX int
	PlayerStartY int
	Lives        int
	CurrentLap   int
	MaxLaps      int
	LapTimes     []time.Duration
	LapStartTime time.Time
}
