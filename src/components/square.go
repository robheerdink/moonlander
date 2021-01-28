package com

import (
	"image/color"
)

// Square is <dunno yet>
type Square struct {
	Object
}

// NewSquare constructor
func NewSquare(id, x, y, z int, v Vector, rx, ry, rw, rh int, c color.RGBA) Square {
	return Square{Object: NewObject(id, nil, x, y, z, v, rx, ry, rw, rh, true, c)}
}

// Collide implements interface Collider
func (o *Square) Collide(hitAbles []GameObject) error {

	for _, h := range hitAbles {
		t := h.GetObject()
		if &o.rect != &t.rect {
			hit, sides := CheckHit(o.GetObject(), t, true, true)
			if hit {
				// fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[t.ID], sides)
				if t.solid {
					if sides.left {
						o.X = float64(t.rect.x-o.rect.w-o.rx) - 1
						o.Vector.x *= -1
					}
					if sides.right {
						o.X = float64(t.rect.x+t.rect.w-o.rx) + 1
						o.Vector.x *= -1
					}
					if sides.top {
						o.Y = float64(t.rect.y-o.rect.h-o.ry) - 1
						o.Vector.y *= -1
					}
					if sides.bottom {
						o.Y = float64(t.rect.y+t.rect.h-o.ry) + 1
						o.Vector.y *= -1
					}
				}
			}
		}
	}
	return nil
}
