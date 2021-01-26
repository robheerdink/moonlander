package com

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Background is a sprite which can be scaled
type Background struct {
	Sprite
	scaleX, scaleY float64
}

// NewBackground constructor
func NewBackground(id int, imagePath string, x, y, z, w, h int, v Vector) Background {
	img, _, err := ebitenutil.NewImageFromFile(imagePath, ebiten.FilterDefault)
	if err != nil {
		log.Fatal(err)
	}
	wImg, hImg := img.Size()
	scaleX := float64(float64(w) / float64(wImg))
	scaleY := float64(float64(h) / float64(hImg))
	fmt.Println(w, h, wImg, hImg, scaleX, scaleY)
	return Background{
		Sprite: NewSprite(id, img, x, y, z, v),
		scaleX: scaleX,
		scaleY: scaleY,
	}
}

// Draw implements Drawer
func (o *Background) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.X, o.Y)
	op.GeoM.Scale(o.scaleX, o.scaleY)
	screen.DrawImage(o.Img, op)
	return nil
}
