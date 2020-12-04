package shared

import "time"

// common used variables in components, values often set per level
var (
	WP WorldProperties
	LP LevelProperties
)

// WorldProperties are properties of the world or level
type WorldProperties struct {
	Gravity     float64
	Friction    float64
	LevelWidth  int
	LevelHeight int
}

// LevelProperties are used to store info / progress off the level
type LevelProperties struct {
	Lives        int
	CurrentLap   int
	MaxLaps      int
	LapTimes     []time.Duration
	LapStartTime time.Time
	PX, PY       int
}
