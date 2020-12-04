package com

import (
	"bytes"
	"image"

	"github.com/hajimehoshi/ebiten"
)

// Frame has info about the animation (start position, dimensions of a frame, number of frames, delay for next frame)
type Frame struct {
	x, y, w, h, num, delay int
}

// NewFrame constructor
func NewFrame(x, y, w, h, num, delay int) Frame {
	return Frame{x, y, w, h, num, delay}
}

// Anim is an animated Sprite
type Anim struct {
	count int
	frame Frame
	Sprite
}

// NewAnim contructor
func NewAnim(img *ebiten.Image, x, y, z int, v Vector, frame Frame) Anim {
	return Anim{
		count:  0,
		frame:  frame,
		Sprite: NewSprite(img, x, y, z, v),
	}
}

// NewAnimFromByte contructor, will convert byte slice to Ebiten image format
func NewAnimFromByte(img []byte, x, y, z int, v Vector, frame Frame) Anim {
	animImg, _, _ := image.Decode(bytes.NewReader(img))
	imgStrip, _ := ebiten.NewImageFromImage(animImg, ebiten.FilterDefault)
	return NewAnim(imgStrip, x, y, z, v, frame)
}

// Update implements Updater
func (o *Anim) Update(screen *ebiten.Image) error {
	if o.count > o.frame.delay*o.frame.num {
		o.count = 0
	} else {
		o.count++
	}
	return nil
}

// Draw implements Drawer
func (o *Anim) Draw(screen *ebiten.Image) error {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(-float64(o.frame.w)/2, -float64(o.frame.h)/2)
	op.GeoM.Rotate(o.z)
	op.GeoM.Translate(o.x, o.y)

	i := (o.count / o.frame.delay) % o.frame.num
	sx, sy := o.frame.x+i*o.frame.w, o.frame.y

	screen.DrawImage(o.img.SubImage(image.Rect(sx, sy, sx+o.frame.w, sy+o.frame.h)).(*ebiten.Image), op)
	return nil
}
