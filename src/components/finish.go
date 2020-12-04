package com

import (
	"image/color"
	"time"
	sha "moonlander/src/shared"
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
			if !sha.LP.LapStartTime.IsZero() {
				duration := time.Now().Sub(sha.LP.LapStartTime)
				sha.LP.LapTimes = append(sha.LP.LapTimes, duration)
			}

			// set startTime of lap
			sha.LP.LapStartTime = time.Now()

			// increase currentLap
			sha.LP.CurrentLap++
			if sha.LP.CurrentLap > sha.LP.MaxLaps {
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
