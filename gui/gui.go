package gui

import (
	//"fmt"
	//"image"

	"image/color"
	"log"

	"github.com/golang/freetype/truetype"
	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/text"
	"golang.org/x/image/font"
)

var (
	fontNormal font.Face
	fontBig    font.Face
	mouse      pointer
)

type pointer struct {
	x, y    int
	pressed bool
}

type clickable interface {
	getName() string
	getActive() bool
	setActive(bool)
	checkHit(x, y int) clickable
	draw(screen *ebiten.Image)
}

type button struct {
	name, text         string
	x, y, w, h         int
	active             bool
	img                *ebiten.Image
	face               font.Face
	btnColor, txtColor color.RGBA
}

func init() {
	tt, err := truetype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	fontNormal = truetype.NewFace(tt, &truetype.Options{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	fontBig = truetype.NewFace(tt, &truetype.Options{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
}

func newButton(name, text string, x, y, w, h int, face font.Face, btnColor, txtColor color.RGBA) button {
	img, _ := ebiten.NewImage(w, h, ebiten.FilterNearest)
	img.Fill(btnColor)
	return button{
		x: x, y: y, w: w, h: h,
		name:     name,
		text:     text,
		active:   true,
		img:      img,
		face:     face,
		btnColor: btnColor,
		txtColor: txtColor,
	}
}

func (b *button) getName() string {
	return b.name
}

func (b *button) getActive() bool {
	return b.active
}

func (b *button) setActive(value bool) {
	b.active = value
}

func (b *button) checkHit(x, y int) clickable {
	if x > b.x && x < b.x+b.w && y > b.y && y < b.y+b.h {
		return b
	}
	return nil
}

func (b *button) draw(screen *ebiten.Image) {
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Translate(float64(b.x), float64(b.y))
	// draw button as image
	screen.DrawImage(b.img, opts)
	// draw text as image (center in button, text y is from position of char '.' in the top line)
	// because of this position multi-line centering is difficult
	textSize := text.BoundString(b.face, b.text)
	textWidth := textSize.Max.X - textSize.Min.X
	textHeight := textSize.Max.Y - textSize.Min.Y
	textX := (b.w - textWidth) / 2
	textY := (b.h + textHeight) / 2
	text.Draw(screen, b.text, fontNormal, b.x+textX, b.y+textY, b.txtColor)
}

// Check hits on group of buttons, if hit set the button active and all other in the group as inactive
func checkHits(x, y int, clickables []clickable) clickable {
	var clicked clickable
	for _, c := range clickables {
		clicked = c.checkHit(mouse.x, mouse.y)
		if clicked != nil && clicked.getActive() {
			for _, o := range clickables {
				o.setActive(true)
			}
			clicked.setActive(false)
			return clicked
		}
	}
	return nil
}

func (p *pointer) update() {
	p.pressed = false
	p.x, p.y = ebiten.CursorPosition()
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		p.pressed = true
	}
}
