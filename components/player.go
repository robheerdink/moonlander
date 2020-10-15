package comp

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math"

	images "moonlander/assets"
	con "moonlander/constants"

	"github.com/hajimehoshi/ebiten"
	"github.com/hajimehoshi/ebiten/ebitenutil"
)

// Player is a controllable object
type Player struct {
	animL, animR, animU, animD    Anim
	imgW, imgH, imgHW, imgHH      float64
	weight, thrust, retro, zSpeed float64
	hx, hy, hw, hh, hdiff         int
	grounded                      bool
	collideObject                 *Object
	Object
	Controls
}

// NewPlayer constructor
func NewPlayer(id int, img *ebiten.Image, x, y, z int, v Vector, rx, ry, rw, rh int, c color.RGBA) Player {
	wImg, hImg := img.Size()
	p := Player{
		Object:   NewObject(id, img, x, y, z, v, rx, ry, rw, rh, true, c),
		Controls: Controls{false, false, false, false, false, false},
		weight:   1,
		thrust:   0.06,
		retro:    0.02,
		zSpeed:   1.2,
		imgW:     float64(wImg),
		imgH:     float64(hImg),
		imgHW:    float64(wImg / 2),
		imgHH:    float64(hImg / 2),
		// keep original hit rect values, to calc rotating hit rect
		hx: rx, hy: ry, hw: rw, hh: rh,
		hdiff: (rw - rh) / 2,
	}

	image, _, _ := image.Decode(bytes.NewReader(images.Runner_png))
	imgStrip, _ := ebiten.NewImageFromImage(image, ebiten.FilterDefault)
	p.animL = NewAnim(imgStrip, 0, 0, 0, NewVector(0, 0), NewFrame(0, 32, 32, 32, 8, 5))
	p.animR = NewAnim(imgStrip, 0, 0, 0, NewVector(0, 0), NewFrame(0, 32, 32, 32, 8, 5))
	p.animU = NewAnim(imgStrip, 0, 0, 0, NewVector(0, 0), NewFrame(0, 32, 32, 32, 8, 5))
	p.animD = NewAnim(imgStrip, 0, 0, 0, NewVector(0, 0), NewFrame(0, 32, 32, 32, 8, 5))
	p.debug = true
	return p
}

// Draw Player
func (o *Player) Draw(screen *ebiten.Image) error {
	if o.img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-o.imgHW, -o.imgHH)
		op.GeoM.Rotate(o.z)
		op.GeoM.Translate(o.x+o.imgHW, o.y+o.imgHH)
		screen.DrawImage(o.img, op)
	}

	if o.Controls.left {
		o.animL.Draw(screen)
	}
	if o.Controls.right {
		o.animR.Draw(screen)
	}
	if o.Controls.up {
		o.animU.Draw(screen)
	}
	if o.Controls.down {
		o.animD.Draw(screen)
	}

	if o.debug {
		// draw hit shape (need to recreate image, because it changes shape)
		if o.rectImg != nil {
			o.rectImg, _ = ebiten.NewImage(o.rect.w, o.rect.h, ebiten.FilterNearest)
			o.rectImg.Fill(con.Cyan25)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(o.rect.x), float64(o.rect.y))
			screen.DrawImage(o.rectImg, op)
		}
		// draw Player related info
		ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nx : %0.2f\ny : %0.2f\nz : %0.2f\nv : %0.4f",
			ebiten.CurrentTPS(), o.x, o.y, o.z, o.vector))
	}

	return nil
}

// GetRotatedPoint transforms points
// cy,cy are the world coordinates of the center of an object.
// ox, oy is the relative offset for a point from the center of the object
// ox, oy will be different per point (top left, top right, etc)
// rad, rotation in radials
func GetRotatedPoint(cx, cy, ox, oy, rad float64) (x, y float64) {
	rx := cx + (ox * math.Cos(rad)) - (oy * math.Sin(rad))
	ry := cy + (ox * math.Sin(rad)) + (oy * math.Cos(rad))
	return rx, ry
}

// Update Player
func (o *Player) Update(screen *ebiten.Image) error {

	o.movePlayer()

	cx := o.x + o.imgHW
	cy := o.y + o.imgHH
	pad := 20.0

	// update animations
	if o.Controls.left {
		// gfx is center right
		o.animL.x, o.animL.y = GetRotatedPoint(cx, cy, +(o.imgHW + pad), 0, o.z)
		o.animL.z = o.z
		o.animL.Update(screen)
	}
	if o.Controls.right {
		// gfx is center left
		o.animR.x, o.animR.y = GetRotatedPoint(cx, cy, -(o.imgHW + pad), 0, o.z)
		o.animR.z = o.z
		o.animR.Update(screen)
	}
	if o.Controls.up {
		// gfx is center below
		o.animU.x, o.animU.y = GetRotatedPoint(cx, cy, 0, +(o.imgHH + pad), o.z)
		o.animU.z = o.z
		o.animU.Update(screen)
	}
	if o.Controls.down {
		// gfx is center top
		o.animD.x, o.animD.y = GetRotatedPoint(cx, cy, 0, -(o.imgHH + pad), o.z)
		o.animD.z = o.z
		o.animD.Update(screen)
	}
	return nil

}

// Collide implements interface Collider, handles collission with ojects
func (o *Player) Collide(hitAbles []HitAble) error {
	o.grounded = false

	for _, h := range hitAbles {
		t := h.GetObject()
		if &o.rect != &t.rect {
			hit, sides := CheckHit(o.GetObject(), t, true, true)
			if hit {
				//fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[t.ID], sides)
				o.addHit(t)

				// player hits somehting solid
				if t.solid {
					if sides.left {
						o.x = float64(t.rect.x-o.rect.w-o.rx) - 1
						o.vector.x = 0
					}
					if sides.right {
						o.x = float64(t.rect.x+t.rect.w-o.rx) + 1
						o.vector.x = 0
					}
					if sides.top {
						// if we hit a wall, set player to grounded and set on top block, without offset
						if t.id == con.IDWall && !o.grounded {
							o.y = float64(t.rect.y - o.rect.h - o.ry)
							o.vector.y = 0
							o.grounded = true
						} else {
							o.y = float64(t.rect.y-o.rect.h-o.ry) - 0
							o.vector.y = 0
						}
					}
					if sides.bottom {
						o.y = float64(t.rect.y+t.rect.h-o.ry) + 1
						o.vector.y = 0
					}
				} else {
					if t.id == con.IDCheckpoint {
						h.SetHit(o)
					} else if t.id == con.IDFinish {
						h.SetHit(o)
					}
				}

			}
		}
	}
	return nil
}

func (o *Player) movePlayer() {

	// keep track of all keys pressed
	o.Controls = Controls{false, false, false, false, false, false}
	if ebiten.IsKeyPressed(ebiten.KeyUp) {
		o.Controls.up = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyDown) && !o.grounded {
		o.Controls.down = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyLeft) && !o.grounded {
		o.Controls.left = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) && !o.grounded {
		o.Controls.right = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyZ) && !o.grounded {
		o.Controls.rl = true
	}
	if ebiten.IsKeyPressed(ebiten.KeyX) && !o.grounded {
		o.Controls.rr = true
	}

	// reset Player
	if ebiten.IsKeyPressed(ebiten.KeyBackspace) {
		o.reset()
	}

	// rotation
	if o.Controls.rl {
		o.z -= o.zSpeed * DegToRad
	}
	if o.Controls.rr {
		o.z += o.zSpeed * DegToRad
	}
	// convert radials always to be always positive between 0 - (2*Pi)
	if o.z < 0 {
		o.z += DPI
	} else if o.z > DPI {
		o.z = 0
	}

	// normalized direction (based on rotation)
	zx := math.Sin(o.z)
	zy := math.Cos(o.z)

	//normalzied possitve direction
	nzx := math.Pow(zx, 2)
	nzy := math.Pow(zy, 2)

	// change hit shape, when rotating (morph between horizontal and vertical shape)
	o.rect.w = int(float64(o.hh)*nzx + float64(o.hw)*nzy)
	o.rect.h = int(float64(o.hw)*nzx + float64(o.hh)*nzy)
	o.rx = o.hx + int(float64(o.hdiff)*nzx)
	o.ry = o.hy + int(float64(o.hdiff)*nzx)*-1

	// ship Velocity
	if o.Controls.up {
		o.vector.x -= o.thrust * zx * -1
		o.vector.y -= o.thrust * zy
	}
	if o.Controls.down {
		o.vector.x += o.retro * zx
		o.vector.y += o.retro * zy
	}
	if o.Controls.right {
		o.vector.x += o.retro * zy
		o.vector.y += o.retro * zx
	}
	if o.Controls.left {
		o.vector.x -= o.retro * zy
		o.vector.y -= o.retro * zx
	}

	if o.grounded {
		// remove horizontal velocity and set ship up right
		o.vector.x *= 0.96

		// turn spaceship uprigth (shortet direction + easout)
		if o.z < PI {
			o.z -= 0.1 * o.z / PI
		} else {
			o.z += 0.1 * (DPI - o.z) / PI
		}

		// snap upright
		if o.z < 0.01 || o.z > PI*2-0.02 {
			o.z = 0
			o.vector.x = 0
		}
	} else {
		// gravity always pushes ship nose-up or nose-down
		// ship is in balance when flying perfectly horizontal
		// ship is slightly more top heavy
		if o.z > 0.01 && o.z < PI {
			if o.z < HPI {
				// go nose up
				o.z -= (WP.Gravity / 6) * ((HPI - o.z) / HPI)
			} else {
				// go nose down
				o.z += (WP.Gravity / 4) * ((HPI - o.z) / HPI) * -1
			}
		}
		if o.z < DPI-0.01 && o.z > PI {
			if o.z < (PI * 1.5) {
				// go nose down
				o.z -= (WP.Gravity / 4) * ((PI/2*3 - o.z) / HPI)
			} else {
				// go nose up
				o.z += (WP.Gravity / 6) * ((PI/2*3 - o.z) / HPI) * -1
			}
		}

		//add 'atmosphere' friction
		o.vector.x *= WP.Friction * o.weight
		o.vector.y *= WP.Friction * o.weight

		// add gravity
		o.vector.y += WP.Gravity * o.weight
	}

	// add velocity to position
	o.x += o.vector.x
	o.y += o.vector.y

	// update hit rect
	o.rect.setXY(int(o.x)+o.rx, int(o.y)+o.ry)

}

func (o *Player) reset() {
	o.x = float64(LP.PX)
	o.y = float64(LP.PY)
	o.z = 0
	o.vector.x = 0
	o.vector.y = 0

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
