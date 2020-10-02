package comp

import (
	"github.com/hajimehoshi/ebiten"
)

// Sprite is something which can be drawn and optionally updated
type Sprite struct {
	ID       int
	Img      *ebiten.Image
	X, Y, Z  float64
	Velocity Vector
}

// NewSprite creates a Sprite, a Sprite can be drawn and optionally updated, sprites have no collision
func NewSprite(id int, img *ebiten.Image, x, y, z float64, v Vector) Sprite {
	return Sprite{id, img, x, y, z, v}
}

// Draw Sprite
func (o *Sprite) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(o.X, o.Y)
	screen.DrawImage(o.Img, op)
	return nil
}

// Update Sprite
func (o *Sprite) Update(screen *ebiten.Image) error {
	o.X += o.Velocity.x
	o.Y += o.Velocity.y
	return nil
}
