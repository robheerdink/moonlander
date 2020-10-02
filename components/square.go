package comp

import (
	"image/color"
)

// Square is <dunno yet>
type Square struct {
	Object
}

// NewSquare constructor
func NewSquare(id int, x, y, z float64, v Vector, rx, ry, rw, rh int, c color.RGBA) Square {
	return Square{Object: NewObjectNoImage(id, x, y, z, v, rx, ry, rw, rh, c)}
}

// GetObject implements interface Collider, so we can get the object from a Collider
func (o *Square) GetObject() *Object {
	return &o.Object
}

// Collide implements interface Collider, handles collission with ojects
func (o *Square) Collide(objects []*Object) error {

	for _, t := range objects {
		if &o.Rect != &t.Rect {
			hit, sides := CheckHit(o.GetObject(), t, true, true)
			if hit {
				// fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[t.ID], sides)
				if sides.left {
					o.X = float64(t.Rect.X-o.Rect.W-o.RX) - 1
					o.Velocity.x *= -1
				}
				if sides.right {
					o.X = float64(t.Rect.X+t.Rect.W-o.RX) + 1
					o.Velocity.x *= -1
				}
				if sides.top {
					o.Y = float64(t.Rect.Y-o.Rect.H-o.RY) - 1
					o.Velocity.y *= -1
				}
				if sides.bottom {
					o.Y = float64(t.Rect.Y+t.Rect.H-o.RY) + 1
					o.Velocity.y *= -1
				}
			}
		}
	}
	return nil
}
