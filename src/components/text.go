package com

import (
	"image/color"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

// Text is a positional textfield, with text, font, and color
type Text struct {
	x, y  int
	text  string
	font  font.Face
	color color.RGBA
}

// NewText contructor
func NewText(x, y int, text string, font font.Face, color color.RGBA) Text {
	return Text{x, y, text, font, color}
}

// Draw it
func (o *Text) Draw(screen *ebiten.Image) error {
	text.Draw(screen, o.text, o.font, o.x, o.y, o.color)
	return nil
}
