package com

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

// constants only used by components
const (
	PI       float64 = math.Pi
	DPI      float64 = math.Pi * 2
	HPI      float64 = math.Pi / 2
	RadToDeg float64 = 180 / math.Pi
	DegToRad float64 = math.Pi / 180
)

// GameObject interface for all game objects
type GameObject interface {
	GetID() (id int)
	GetInfo() (id int, name string, x, y, r float64, w, h int)
	Draw(screen *ebiten.Image) error
	Update(screen *ebiten.Image) error
	SetHit(collider GameObject)
	GetObject() *Object
	Collide(hitList []GameObject) error
}

// Vector used for direction of objects
type Vector struct {
	x, y float64
}

// NewVector creates a Vector
func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

//Rect as format x,y,w,h
type Rect struct {
	x, y, w, h int
}

//NewRect creates a Rect
func NewRect(x, y, w, h int) Rect {
	return Rect{x, y, w, h}
}

// setXY position of Rect
func (r *Rect) setXY(x, y int) {
	r.x = x
	r.y = y
}

// HitShape is used for custom hit area for objects
type HitShape struct {
	rx, ry    int
	rect      Rect
	rectImg   *ebiten.Image
	rectColor color.RGBA
	Hit       bool
	solid     bool
}
