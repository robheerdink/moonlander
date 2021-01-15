package com

import (
	"github.com/hajimehoshi/ebiten"
)

// Sprite is something which can be drawn and optionally updated
type Sprite struct {
	img     *ebiten.Image
	x, y, z float64
	vector  Vector
}

// NewSprite creates a Sprite, a Sprite can be drawn and optionally updated, sprites have no collision
func NewSprite(img *ebiten.Image, x, y, z int, v Vector) Sprite {
	return Sprite{img, float64(x), float64(y), float64(z), v}
}

// Draw implements Drawer
func (o *Sprite) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.x, o.y)
	screen.DrawImage(o.img, op)
	return nil
}

// GetImageInfo implements Drawer
func (o *Sprite) GetImageInfo() (x, y, z float64, img *ebiten.Image) {
	return  o.x, o.y, o.z, o.img
}

// Update Sprite
func (o *Sprite) Update(screen *ebiten.Image) error {
	o.x += o.vector.x
	o.y += o.vector.y
	return nil
}
