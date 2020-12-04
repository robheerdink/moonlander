package com

import (
	"image/color"
)

// Wall is something you can smack in to
type Wall struct {
	Object
}

// NewWall constructor
func NewWall(id, x, y, w, h int, c color.RGBA) Wall {
	return Wall{Object: NewObject(id, nil, x, y, 0, Vector{}, 0, 0, w, h, true, c)}
}
