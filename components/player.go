package comp

import (
	"fmt"
	"image/color"
	"math"

	con "moonlander/constants"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Player is a controllable object
type Player struct {
	weight, mainThrust, retroThrust, rotateSpeed, imgHW, imgHH, rectHW, rectHH float64
	Object
	Controls
	collideObject *Object
}

// NewPlayer constructor
func NewPlayer(id int, img *ebiten.Image, x, y, z float64, v Vector, rx, ry, rw, rh int, c color.RGBA) Player {
	wImg, hImg := img.Size()
	return Player{
		Object:      NewObject(id, img, x, y, z, v, rx, ry, rw, rh, c),
		Controls:    Controls{false, false, false, false, false, false},
		weight:      1,
		mainThrust:  0.06,
		retroThrust: 0.02,
		rotateSpeed: 1,
		imgHW:       float64(wImg / 2),
		imgHH:       float64(hImg / 2),
		rectHW:      float64(rw / 2),
		rectHH:      float64(rh / 2),
	}
}

// Draw Player
func (o *Player) Draw(screen *ebiten.Image) error {
	if o.RectImg != nil {
		// dont rotate hit, because hit rotation is not calculated in collision
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(o.Rect.X), float64(o.Rect.Y))
		screen.DrawImage(o.RectImg, op)
	}
	if o.Img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-o.imgHW, -o.imgHH)
		op.GeoM.Rotate(o.Z)
		op.GeoM.Translate(o.X+o.imgHW, o.Y+o.imgHH)
		screen.DrawImage(o.Img, op)
	}
	// draw Player related info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nx : %0.2f\ny : %0.2f\nz : %0.2f\nv : %0.4f",
		ebiten.CurrentTPS(), o.X, o.Y, o.Z, o.Velocity))
	return nil
}

// Update Player
func (o *Player) Update(screen *ebiten.Image) error {
	o.movePlayer()
	return nil
}

// GetObject implements interface Collider, so we can get the object from a Collider
func (o *Player) GetObject() *Object {
	return &o.Object
}

// Collide implements interface Collider, handles collission with ojects
func (o *Player) Collide(objects []*Object) error {

	for _, t := range objects {
		if &o.Rect != &t.Rect {
			hit, sides := CheckHit(o.GetObject(), t, true, true)
			if hit {
				fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[t.ID], sides)
				o.addHit(t)
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

func (o *Player) movePlayer() {
	// update hit rect
	o.Rect.SetXY(int(o.X)+o.RX, int(o.Y)+o.RY)

	// keep track of all keys pressed
	o.Controls = Controls{false, false, false, false, false, false}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		o.Controls.up = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) {
		o.Controls.down = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		o.Controls.left = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		o.Controls.right = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) {
		o.Controls.rl = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) {
		o.Controls.rr = true
	}

	// reset Player
	if ebiten.IsKeyPressed(ebiten.KeyC) {
		o.reset()
	}

	// rotation
	if o.Controls.rl {
		o.Z -= o.rotateSpeed * con.DegToRad
	}
	if o.Controls.rr {
		o.Z += o.rotateSpeed * con.DegToRad
	}
	if math.Abs(o.Z) > math.Pi*2 {
		o.Z = 0
	}

	// normalized direction (based on rotation)
	zx := math.Sin(o.Z)
	zy := math.Cos(o.Z)

	// ship Velocity
	if o.Controls.up {
		o.Velocity.x -= o.mainThrust * zx * -1
		o.Velocity.y -= o.mainThrust * zy
	}
	if o.Controls.down {
		o.Velocity.x += o.retroThrust * zx
		o.Velocity.y += o.retroThrust * zy
	}
	if o.Controls.right {
		o.Velocity.x += o.retroThrust * zy
		o.Velocity.y += o.retroThrust * zx
	}
	if o.Controls.left {
		o.Velocity.x -= o.retroThrust * zy
		o.Velocity.y -= o.retroThrust * zx
	}

	//add 'atmosphere' friction
	o.Velocity.x *= WP.Friction * o.weight
	o.Velocity.y *= WP.Friction * o.weight

	// add gravity
	o.Velocity.y += WP.Gravity * o.weight

	// velocity to position
	o.X += o.Velocity.x
	o.Y += o.Velocity.y
}

func (o *Player) reset() {
	o.X = float64(WP.LevelWidth / 2)
	o.Y = float64(WP.LevelHeight/2) + 100
	o.Velocity.x = 0
	o.Velocity.y = 0
}

func (o *Player) addHit(obj *Object) {
	obj.Hit = true
	o.collideObject = obj
}

func (o *Player) removeHit() {
	if o.collideObject != nil {
		o.collideObject.Hit = false
		o.collideObject = nil
	}
}
