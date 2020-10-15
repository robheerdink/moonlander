package comp

import (
	"image/color"

	con "moonlander/constants"

	"github.com/hajimehoshi/ebiten"
)

// Object is a Sprite with collision
type Object struct {
	id int
	Sprite
	HitShape
	debug bool
}

// NewObject creates new Object from source image, use hx, hy, w, h to position and size the hitbox,
func NewObject(id int, img *ebiten.Image, x, y, z int, v Vector, rx, ry, rw, rh int, solid bool, c color.RGBA) Object {

	// hit rect image (centered in image)
	rectImg, _ := ebiten.NewImage(rw, rh, ebiten.FilterNearest)
	rectImg.Fill(con.Red50)
	if img == nil {
		// create visual image, when bitmap is nil
		// the size of the image is hit rect + the offsets on both sides
		img, _ = ebiten.NewImage(rw+(rx*2), rh+(ry*2), ebiten.FilterNearest)
		img.Fill(c)
	}
	return Object{
		id:     id,
		Sprite: NewSprite(img, x, y, z, v),
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

// Draw implements Drawer
func (o *Object) Draw(screen *ebiten.Image) error {
	if o.img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(o.x, o.y)
		screen.DrawImage(o.img, op)
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
	o.x += o.vector.x
	o.y += o.vector.y
	o.rect.setXY(int(o.x)+o.rx, int(o.y)+o.ry)
	return nil
}

// SetHit implements HitAble
func (o *Object) SetHit(collider Collider) {
	//fmt.Println("set hit called")
}

// GetObject implements HitAble
func (o *Object) GetObject() *Object {
	return o
}

// GetRect returns the object hitshape rect
func (o *Object) GetRect() *Rect {
	return &o.rect
}

// GetSolid returns if a hitshape is solid
func (o *Object) GetSolid() bool {
	return o.solid
}
