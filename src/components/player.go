package com

import (
	"fmt"
	"image/color"
	"math"

	ass "moonlander/assets"
	sha "moonlander/src/shared"

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

// Controls stuff
type Controls struct {
	up, down, left, right, rr, rl bool
}

// NewPlayer constructor
func NewPlayer(id int, img *ebiten.Image, x, y, z int, v Vector, hx, hy, hw, hh int, c color.RGBA) Player {
	wImg, hImg := img.Size()
	p := Player{
		Object:   NewObject(id, img, x, y, z, v, hx, hy, hw, hh, true, c),
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
		hx: hx, hy: hy, hw: hw, hh: hh,
		hdiff: (hw - hh) / 2,
	}
	p.animU = NewAnimFromByte(ass.Up, 0, 0, 0, NewVector(0, 0), NewFrame(0, 0, 20, 48, 3, 5))
	p.animD = NewAnimFromByte(ass.Down, 0, 0, 0, NewVector(0, 0), NewFrame(0, 0, 10, 32, 3, 5))
	p.animL = NewAnimFromByte(ass.Left, 0, 0, 0, NewVector(0, 0), NewFrame(0, 0, 32, 10, 3, 5))
	p.animR = NewAnimFromByte(ass.Right, 0, 0, 0, NewVector(0, 0), NewFrame(0, 0, 32, 10, 3, 5))
	p.debug = false
	return p
}

// Draw Player
func (o *Player) Draw(screen *ebiten.Image) error {
	if o.Controls.left || o.Controls.rl {
		o.animL.Draw(screen)
	}
	if o.Controls.right || o.Controls.rr {
		o.animR.Draw(screen)
	}
	if o.Controls.up {
		o.animU.Draw(screen)
	}
	if o.Controls.down {
		o.animD.Draw(screen)
	}
	if o.img != nil {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-o.imgHW, -o.imgHH)
		op.GeoM.Rotate(o.z)
		op.GeoM.Translate(o.x+o.imgHW, o.y+o.imgHH)
		screen.DrawImage(o.img, op)
	}
	if o.debug {
		// draw hit shape (need to recreate image, because it changes shape)
		if o.rectImg != nil {
			o.rectImg, _ = ebiten.NewImage(o.rect.w, o.rect.h, ebiten.FilterNearest)
			o.rectImg.Fill(sha.Cyan25)
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(o.rect.x), float64(o.rect.y))
			screen.DrawImage(o.rectImg, op)
		}
	}

	// draw Player related info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("FPS: %0.2f\nx : %0.2f\ny : %0.2f\nz : %0.2f\nv : %0.4f",
		ebiten.CurrentTPS(), o.x, o.y, o.z, o.vector))

	return nil
}

// Update Player
func (o *Player) Update(screen *ebiten.Image) error {
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
	} else if o.z >= DPI {
		o.z = 0
	}

	// direction (based on rotation)
	zx := math.Sin(o.z)
	zy := math.Cos(o.z)

	//change hit shape, when rotating (morph between horizontal and vertical shape)
	hMP := math.Pow(zx, 2) //[0 ..1] horizontal
	vMP := math.Pow(zy, 2) //[0 ..1] vertical
	o.rect.w = int(float64(o.hh)*hMP + float64(o.hw)*vMP)
	o.rect.h = int(float64(o.hw)*hMP + float64(o.hh)*vMP)
	o.rx = o.hx + int(float64(o.hdiff)*hMP)
	o.ry = o.hy + int(float64(o.hdiff)*hMP)*-1

	// add velocity when pressing certan keys
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
		// when grounded, remove horizontal velocity and set ship facing up
		o.vector.x *= 0.96
		if o.z < PI {
			o.z -= 0.1 * (o.z / PI)
		} else {
			o.z += 0.1 * (DPI - o.z) / PI
		}
		// snap last part, because we ease out in rotation
		if o.z < 0.01 || o.z > PI*2-0.01 {
			o.z = 0
			o.vector.x = 0
		}
	} else {
		// gravity always pushes ship nose-up or nose-down
		// ship is in balance when flying perfectly horizontal
		if o.z > 0.01 && o.z < PI {
			if o.z < HPI {
				// go nose up
				o.z -= (sha.WP.Gravity / 6) * ((HPI - o.z) / HPI)
			} else {
				// go nose down
				o.z += (sha.WP.Gravity / 4) * ((HPI - o.z) / HPI) * -1
			}
		}
		if o.z < DPI-0.01 && o.z > PI {
			if o.z < (PI * 1.5) {
				// go nose down
				o.z -= (sha.WP.Gravity / 4) * ((PI/2*3 - o.z) / HPI)
			} else {
				// go nose up
				o.z += (sha.WP.Gravity / 6) * ((PI/2*3 - o.z) / HPI) * -1
			}
		}

		//add 'atmosphere' friction
		o.vector.x *= sha.WP.Friction * o.weight
		o.vector.y *= sha.WP.Friction * o.weight

		// add gravity
		o.vector.y += sha.WP.Gravity * o.weight
	}

	// update player position
	o.x += o.vector.x
	o.y += o.vector.y

	// update hit rect
	o.rect.setXY(int(o.x)+o.rx, int(o.y)+o.ry)

	// also update anim location + rotation, based on player rotation + location
	// do this after player postion is updated
	cx, cy := o.x+o.imgHW, o.y+o.imgHH
	if o.Controls.up {
		o.animU.x, o.animU.y = GetRotatedPoint(cx, cy, 0, +(o.imgHH + 16), o.z)
		o.animU.z = o.z
		o.animU.Update(screen)
	}
	if o.Controls.down {
		o.animD.z = o.z
		o.animD.x, o.animD.y = GetRotatedPoint(cx, cy, 0, -(o.imgHH + 16), o.z)
		o.animD.Update(screen)
	}
	if o.Controls.right || o.Controls.rr {
		o.animR.z = o.z
		o.animR.x, o.animR.y = GetRotatedPoint(cx, cy, -(o.imgHW + 8), -10, o.z)
		o.animR.Update(screen)
	}
	if o.Controls.left || o.Controls.rl {
		o.animL.z = o.z
		o.animL.x, o.animL.y = GetRotatedPoint(cx, cy, +(o.imgHW + 8), -10, o.z)
		o.animL.Update(screen)
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
				//fmt.Printf("%s hits %s on sides %+v\n", sha.ID[o.ID], sha.ID[t.ID], sides)
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
						if t.id == sha.IDWall && !o.grounded {
							o.y = float64(t.rect.y - o.rect.h - o.ry)
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
					if t.id == sha.IDCheckpoint {
						h.SetHit(o)
					} else if t.id == sha.IDFinish {
						h.SetHit(o)
					}
				}

			}
		}
	}
	return nil
}

func (o *Player) reset() {
	o.x = float64(sha.LP.PX)
	o.y = float64(sha.LP.PY)
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
