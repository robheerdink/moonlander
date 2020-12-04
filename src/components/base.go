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

// Drawer can be drawn every frame
type Drawer interface {
	// Draw in screen
	Draw(screen *ebiten.Image) error
	// Draw in world
	GetImageInfo() (x, y, z float64, img *ebiten.Image)
}

// Updater can be updated every frame
type Updater interface {
	Update(screen *ebiten.Image) error
}

// HitAble something that can be hit / collided with
type HitAble interface {
	SetHit(collider Collider)
	GetObject() *Object
}

// Collider checks collisions with HitAble's
type Collider interface {
	Collide(hitList []HitAble) error
	GetObject() *Object
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