package com

import (
	"image/color"

	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
)

// Object is a Sprite with collision
type Object struct {
	Sprite
	HitShape
	debug bool
}

// NewObject creates new Object from source image, use hx, hy, w, h to position and size the hitbox,
func NewObject(id int, img *ebiten.Image, x, y, z int, v Vector, rx, ry, rw, rh int, solid bool, c color.RGBA) Object {

	// hit rect image (centered in image)
	rectImg, _ := ebiten.NewImage(rw, rh, ebiten.FilterNearest)
	rectImg.Fill(sha.Red50)
	if img == nil {
		// create visual image, when bitmap is nil
		// the size of the image is hit rect + the offsets on both sides
		img, _ = ebiten.NewImage(rw+(rx*2), rh+(ry*2), ebiten.FilterNearest)
		img.Fill(c)
	}
	return Object{
		Sprite: NewSprite(id, img, x, y, z, v),
		HitShape: HitShape{
			rx:        rx,
			ry:        ry,
			rect:      Rect{x + rx, y + ry, rw, rh},
			rectImg:   rectImg,
			rectColor: c,
			solid:     solid,
		},
		debug: false,
	}
}

// GetID implements interface, returns an id
func (o *Object) GetID() int {
	return o.ID
}

// GetInfo implements interface, returns info for debugging
func (o *Object) GetInfo() (id int, name string, x, y, r float64, w, h int) {
	w, h = o.Img.Size()
	return o.ID, sha.Name[o.ID], o.X, o.Y, o.R, w, h
}

// GetObject implements interface
func (o *Object) GetObject() *Object {
	return o
}

// SetHit implements interface
func (o *Object) SetHit(collider Collider) {
	//fmt.Println("set hit called")
}

// Draw i implements interface
func (o *Object) Draw(screen *ebiten.Image) error {
	if o.Img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(o.X, o.Y)
		screen.DrawImage(o.Img, op)
	}
	if o.rectImg != nil {
		// only draw hit rect, when it gets a hit tag
		if o.Hit {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(o.rect.x), float64(o.rect.y))
			screen.DrawImage(o.rectImg, op)
		}
	}
	return nil
}

// Update implements Updater
func (o *Object) Update(screen *ebiten.Image) error {
	o.X += o.Vector.x
	o.Y += o.Vector.y
	o.rect.setXY(int(o.X)+o.rx, int(o.Y)+o.ry)
	return nil
}

// GetRect returns the object hitshape rect
func (o *Object) GetRect() *Rect {
	return &o.rect
}

// GetSolid returns if a hitshape is solid
func (o *Object) GetSolid() bool {
	return o.solid
}
