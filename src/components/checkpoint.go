package com

import (
	"image/color"
)

// Checkpoint is <dunno yet>
type Checkpoint struct {
	Object
	done bool
}

// NewCheckpoint constructor
func NewCheckpoint(id, x, y, w, h int, c color.RGBA, done bool) Checkpoint {
	return Checkpoint{
		Object: NewObject(id, nil, x, y, 0, Vector{}, 0, 0, w, h, false, c),
		done:   done,
	}
}

// SetHit Override
func (o *Checkpoint) SetHit(collider Collider) {
	o.Hit = true
	o.done = true
}
