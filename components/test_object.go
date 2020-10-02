package comp

import (
	"image/color"
	"math"

	con "moonlander/constants"

	"github.com/hajimehoshi/ebiten"
)

// TestObject is used for testing collisions, its a controllable square
type TestObject struct {
	speed, imgHW, imgHH, rectHW, rectHH float64
	Object
	collideObject *Object
}

// NewCollideTest constructor
func NewCollideTest(id int, x, y, z float64, v Vector, rx, ry, rw, rh int, c color.RGBA) TestObject {
	return TestObject{
		Object: NewObjectNoImage(id, x, y, z, v, rx, ry, rw, rh, c),
		speed:  0.8,
		imgHW:  float64(rw/2 + rx),
		imgHH:  float64(rh/2 + ry),
		rectHW: float64((rw) / 2),
		rectHH: float64((rh) / 2),
	}
}

//Draw overrides Drawable interface (from default Sprite)
func (o *TestObject) Draw(screen *ebiten.Image) error {
	// draw vissual image
	if o.Img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-o.imgHW, -o.imgHW)
		op.GeoM.Rotate(o.Z)
		op.GeoM.Translate(o.X+o.imgHW, o.Y+o.imgHH)
		screen.DrawImage(o.Img, op)
	}
	// draw hit rect
	if o.RectImg != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.Rect.X), float64(o.Rect.Y))
		screen.DrawImage(o.RectImg, op)
	}
	return nil
}

// Update ..
func (o *TestObject) Update(screen *ebiten.Image) error {
	o.removeHit()

	// slow down
	o.Velocity.x *= 0.9
	o.Velocity.y *= 0.9

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		o.Velocity.y = o.speed * -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		o.Velocity.y = o.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		o.Velocity.x = o.speed * -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		o.Velocity.x = o.speed
	}
	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		o.Z -= (o.speed * 2) * con.DegToRad
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		o.Z += (o.speed * 2) * con.DegToRad
	}
	if math.Abs(o.Z) > math.Pi*2 {
		o.Z = 0
	}

	// update position
	o.X += o.Velocity.x
	o.Y += o.Velocity.y

	// update hit rect
	o.Rect.SetXY(int(o.X)+o.RX, int(o.Y)+o.RY)
	return nil
}

// GetObject implements interface Collider, so we can get the object from a Collider
func (o *TestObject) GetObject() *Object {
	return &o.Object
}

//Collide implements interface Collider, handles collission with ojects
func (o *TestObject) Collide(objects []*Object) error {
	for _, t := range objects {
		if &o.Rect != &t.Rect {
			hit, sides := CheckHit(o.GetObject(), t, true, true)
			if hit {
				o.addHit(t)
				// fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[t.ID], sides)
				if sides.left {
					o.X = float64(t.Rect.X-o.Rect.W-o.RX) - 1
					o.Velocity.x = 0
				}
				if sides.right {
					o.X = float64(t.Rect.X+t.Rect.W-o.RX) + 1
					o.Velocity.x = 0
				}
				if sides.top {
					o.Y = float64(t.Rect.Y-o.Rect.H-o.RY) - 1
					o.Velocity.y = 0
				}
				if sides.bottom {
					o.Y = float64(t.Rect.Y+t.Rect.H-o.RY) + 1
					o.Velocity.y = 0
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
