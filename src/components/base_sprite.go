package com

import (
	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
)

// Sprite is something which can be drawn and optionally updated
type Sprite struct {
	ID      int
	Img     *ebiten.Image
	X, Y, R float64
	Vector  Vector
}

// NewSprite creates a Sprite
// a Sprite can be drawn and optionally updated, sprites have no collision
func NewSprite(id int, img *ebiten.Image, x, y, z int, v Vector) Sprite {
	return Sprite{id, img, float64(x), float64(y), float64(z), v}
}

// GetID implements interface, returns an id
func (o *Sprite) GetID() int {
	return o.ID
}

// GetInfo implements interface, returns info for debugging
func (o *Sprite) GetInfo() (id int, name string, x, y, r float64, w, h int) {
	w, h = o.Img.Size()
	return o.ID, sha.Name[o.ID], o.X, o.Y, o.R, w, h
}

// Draw implements interface
func (o *Sprite) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.X, o.Y)
	screen.DrawImage(o.Img, op)
	return nil
}

// Update Sprite
func (o *Sprite) Update(screen *ebiten.Image) error {
	o.X += o.Vector.x
	o.Y += o.Vector.y
	return nil
}
