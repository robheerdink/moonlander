package com

import (
	"image/color"
	"math"

	sha "moonlander/src/shared"

	"github.com/hajimehoshi/ebiten"
)

// TestObject is used for testing collisions, its a controllable square
type TestObject struct {
	speed, imgHW, imgHH             float64
	hitX, hitY, hitW, hitH, hitDiff int
	collideObject                   *Object
	Object
}

// NewCollideTest constructor
func NewCollideTest(id, x, y, z int, v Vector, rx, ry, rw, rh int, c color.RGBA) TestObject {
	return TestObject{
		Object: NewObject(id, nil, x, y, z, v, rx, ry, rw, rh, true, c),
		speed:  0.8,
		imgHW:  float64(rw/2 + rx),
		imgHH:  float64(rh/2 + ry),
		// keep original hit rect values, to calc rotating hit rect
		hitX:    rx,
		hitY:    ry,
		hitW:    rw,
		hitH:    rh,
		hitDiff: (rw - rh) / 2,
	}
}

//Draw overrides Drawable interface (from default Sprite)
func (o *TestObject) Draw(screen *ebiten.Image) error {
	// draw vissual image
	if o.Img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-o.imgHW, -o.imgHH)
		op.GeoM.Rotate(o.R)
		op.GeoM.Translate(o.X+o.imgHW, o.Y+o.imgHH)
		screen.DrawImage(o.Img, op)
	}
	// draw hit rect
	if o.rectImg != nil {
		// need to recreate image, because it changes shape (expensive only use for debug)
		o.rectImg, _ = ebiten.NewImage(o.rect.w, o.rect.h, ebiten.FilterNearest)
		o.rectImg.Fill(sha.Cyan50)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.rect.x), float64(o.rect.y))
		screen.DrawImage(o.rectImg, op)
	}
	return nil
}

// Update ..
func (o *TestObject) Update(screen *ebiten.Image) error {
	o.removeHit()

	// slow down
	o.Vector.x *= 0.9
	o.Vector.y *= 0.9

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		o.Vector.y = o.speed * -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyG) {
		o.Vector.y = o.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyF) {
		o.Vector.x = o.speed * -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyH) {
		o.Vector.x = o.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyR) {
		o.R -= (o.speed * 2) * DegToRad
	}
	if ebiten.IsKeyPressed(ebiten.KeyY) {
		o.R += (o.speed * 2) * DegToRad
	}
	if math.Abs(o.R) > DPI {
		o.R = 0
	}

	zx := math.Sin(o.R)
	zy := math.Cos(o.R)

	nzx := math.Pow(zx, 2)
	nzy := math.Pow(zy, 2)
	o.rect.w = int(float64(o.hitH)*nzx + float64(o.hitW)*nzy)
	o.rect.h = int(float64(o.hitW)*nzx + float64(o.hitH)*nzy)
	o.rx = o.hitX + int(float64(o.hitDiff)*nzx)
	o.ry = o.hitY + int(float64(o.hitDiff)*nzx)*-1

	// update position
	o.X += o.Vector.x
	o.Y += o.Vector.y

	// update hit rect
	o.rect.setXY(int(o.X)+o.rx, int(o.Y)+o.ry)
	return nil
}

//Collide implements interface Collider, handles collission with ojects
func (o *TestObject) Collide(hitAbles []GameObject) error {

	for _, h := range hitAbles {
		t := h.GetObject()
		if &o.rect != &t.rect {
			hit, sides := CheckHit(o.GetObject(), t, true, true)
			if hit {
				// fmt.Printf("%s hits %s on sides %+v\n", sha.ID[o.ID], sha.ID[t.ID], sides)
				if t.solid {
					o.addHit(t)
					if sides.left {
						o.X = float64(t.rect.x-o.rect.w-o.rx) - 1
						o.Vector.x = 0
					}
					if sides.right {
						o.X = float64(t.rect.x+t.rect.w-o.rx) + 1
						o.Vector.x = 0
					}
					if sides.top {
						o.Y = float64(t.rect.y-o.rect.h-o.ry) - 1
						o.Vector.y = 0
					}
					if sides.bottom {
						o.Y = float64(t.rect.y+t.rect.h-o.ry) + 1
						o.Vector.y = 0
					}
				}
			}
		}
	}
	return nil
}

func (o *TestObject) addHit(obj *Object) {
	obj.Hit = true
	o.collideObject = obj
}

func (o *TestObject) removeHit() {
	if o.collideObject != nil {
		o.collideObject.Hit = false
		o.collideObject = nil
	}
}
