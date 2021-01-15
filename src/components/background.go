package com

import (
	"github.com/hajimehoshi/ebiten"
)

// Background is a sprite which can be scaled
type Background struct {
	Sprite
	scaleX, scaleY float64
}

// NewBackground constructor
func NewBackground(img *ebiten.Image, x, y, z, w, h int, v Vector) Background {
	wImg, hImg := img.Size()
	scaleX := float64(w / wImg)
	scaleY := float64(h / hImg)
	return Background{
		Sprite: NewSprite(img, x, y, z, v),
		scaleX: scaleX,
		scaleY: scaleY,
	}
}

// Draw implements Drawer
func (o *Background) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.x, o.y)
	op.GeoM.Scale(o.scaleX, o.scaleY)
	screen.DrawImage(o.img, op)
	return nil
}
