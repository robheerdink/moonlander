package comp

import (
	"image/color"
	"time"
)

// Finish is <dunno yet>
type Finish struct {
	Object
	checkpoints []*Checkpoint
	finished    bool
}

// NewFinish constructor
func NewFinish(id, x, y, w, h int, c color.RGBA, checkpoints []*Checkpoint) Finish {
	return Finish{
		Object:      NewObject(id, nil, x, y, 0, Vector{}, 0, 0, w, h, false, c),
		checkpoints: checkpoints,
	}
}

// SetHit Override
func (o *Finish) SetHit(collider Collider) {

	if !o.finished {

		// check if we passed all checkpoints
		allHit := true
		for _, cp := range o.checkpoints {
			if !cp.done {
				allHit = false
				break
			}
		}

		// valid lap
		if allHit {
			// save duration's of laps
			if !LP.LapStartTime.IsZero() {
				duration := time.Now().Sub(LP.LapStartTime)
				LP.LapTimes = append(LP.LapTimes, duration)
			}

			// set startTime of lap
			LP.LapStartTime = time.Now()

			// increase currentLap
			LP.CurrentLap++
			if LP.CurrentLap > LP.MaxLaps {
				o.finished = true
			}

			// reset checkpoint
			for _, cp := range o.checkpoints {
				cp.done = false
				cp.Hit = false
			}
		}
	}

}
