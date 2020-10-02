package comp

import (
	"image/color"

	con "moonlander/constants"

	"github.com/hajimehoshi/ebiten"
)

// Object is a Sprite with collision
type Object struct {
	Sprite
	HitShape
}

// NewObject creates new Object from source image, use hx, hy, w, h to position and size the hitbox,
func NewObject(id int, img *ebiten.Image, x, y, z float64, v Vector, rx, ry, rw, rh int, c color.RGBA) Object {
	rectImg, _ := ebiten.NewImage(int(rw), int(rh), ebiten.FilterNearest)
	rectImg.Fill(c)

	return Object{
		Sprite: NewSprite(id, img, x, y, z, v),
		HitShape: HitShape{
			RX:        rx,
			RY:        ry,
			Rect:      Rect{int(x) + rx, int(y) + ry, rw, rh},
			RectImg:   rectImg,
			RectColor: c,
		}}
}

// NewObjectNoImage uses no source image, an image will be draw from width, height, color.
func NewObjectNoImage(id int, x, y, z float64, v Vector, rx, ry, rw, rh int, c color.RGBA) Object {
	// visual image size
	// vissual size is rect dimensions + the offsets from rx, ry on both sides
	img, _ := ebiten.NewImage(int(rw+(rx*2)), int(rh+(ry*2)), ebiten.FilterNearest)
	img.Fill(c)

	// hit rect image (is centered in image)
	rectImg, _ := ebiten.NewImage(int(rw), int(rh), ebiten.FilterNearest)
	rectImg.Fill(con.Red)

	return Object{
		Sprite: NewSprite(id, img, x, y, z, v),
		HitShape: HitShape{
			RX:        rx,
			RY:        ry,
			Rect:      Rect{int(x) + rx, int(y) + ry, rw, rh},
			RectImg:   rectImg,
			RectColor: c,
		}}
}

// Draw Object
func (o *Object) Draw(screen *ebiten.Image) error {
	if o.Img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(o.X, o.Y)
		screen.DrawImage(o.Img, op)
	}
	if o.RectImg != nil {
		// only draw hit rect, when it gets a hit tag
		if o.Hit {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(o.Rect.X), float64(o.Rect.Y))
			screen.DrawImage(o.RectImg, op)
		}
	}
	return nil
}

// Update Object
func (o *Object) Update(screen *ebiten.Image) error {
	o.X += o.Velocity.x
	o.Y += o.Velocity.y
	o.Rect.SetXY(int(o.X)+o.RX, int(o.Y)+o.RY)
	return nil
}
