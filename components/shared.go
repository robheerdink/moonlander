package comp

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten"
)

var (
	// WP will be set on level creation, Player need to be able to access World Properties
	WP WorldProperties
)

// WorldProperties are properties of the world or level
type WorldProperties struct {
	Gravity     float64
	Friction    float64
	LevelWidth  int
	LevelHeight int
}

// Drawer ..
type Drawer interface {
	Draw(screen *ebiten.Image) error
}

// Updater ..
type Updater interface {
	Update(screen *ebiten.Image) error
}

// Collider ..
type Collider interface {
	Collide(objectList []*Object) error
	GetObject() *Object
}

// Vector used for direction of objects
type Vector struct {
	x, y float64
}

// NewVector creates a Vector with direction and speed
func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

//Rect as format x,y,w,h
type Rect struct {
	X, Y, W, H int
}

// SetXY position of Rect
func (r *Rect) SetXY(x, y int) {
	r.X = x
	r.Y = y
}

// HitShape is used for custom hit area for objects
type HitShape struct {
	RX, RY    int
	Rect      Rect
	RectImg   *ebiten.Image
	RectColor color.RGBA
	Hit       bool
}

// Controls stuff (TODO should not be part of object)
type Controls struct {
	up, down, left, right, rr, rl bool
}

// Sides has a boolean for each side, to indicate collision on specific side('s)
type Sides struct {
	left, right, top, bottom bool
}

// HitHV returs true if there is a horizontal and vertical hit
func (s *Sides) HitHV() bool {
	return (s.left || s.right) && (s.top || s.bottom)
}

// CheckHit checks for a hit between two objects and
// optionally returns the hit side('s) from which we hit the targets
func CheckHit(o *Object, target *Object, checkSides bool, resolveHV bool) (bool, Sides) {
	r, t := &o.Rect, &target.Rect
	var sides Sides
	if CheckOverlap(r, t) {
		if checkSides {
			CheckSides(&sides, r, t)
			if resolveHV {
				ResolveHitHV(&sides, GetIntersetRect(r, t))
				//fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[target.ID], sides)
			}
		}
		return true, sides
	}
	return false, sides
}

// CheckOverlap checks for a hit between two rects
func CheckOverlap(r *Rect, t *Rect) bool {
	if !(r.X > t.X+t.W || r.X+r.W < t.X || r.Y > t.Y+t.H || r.Y+r.H < t.Y) {
		return true
	}
	return false
}

// CheckSides returns which side('s) we hit a target (assumes an overlap)
func CheckSides(sides *Sides, r *Rect, t *Rect) {
	// dont remove equals signs, else e.g if the collider is bigger then the target
	// a hit bottom could return only hit Left, right as true
	// most likely rounding between position float and integers of the hit rect
	aLeft, aRight := r.X, r.X+r.W
	aTop, aBottom := r.Y, r.Y+r.H
	bLeft, bRight := t.X, t.X+t.W
	bTop, bBottom := t.Y, t.Y+t.H
	if aRight >= bLeft && aLeft < bLeft {
		sides.left = true
	}
	if aLeft <= bRight && aRight > bRight {
		sides.right = true
	}
	if aBottom >= bTop && aTop < bTop {
		sides.top = true
	}
	if aTop <= bBottom && aBottom > bBottom {
		sides.bottom = true
	}
}

// ResolveHitHV resolves collisions in case we have a horizontal and vertical collision
// with the same object, based on the overlap dimensions it ignores the horizontal or vertical hit.
// so the collider only react to the biggest hit
func ResolveHitHV(sides *Sides, ir *Rect) {
	//fmt.Println("ResolveHitHV, intersect rect: ", ir)
	if sides.HitHV() {
		if ir.W > ir.H {
			sides.left, sides.right = false, false
		} else {
			sides.top, sides.bottom = false, false
		}
	}
}

// GetIntersetRect returns the intersection rect between two rectanges
func GetIntersetRect(r *Rect, t *Rect) *Rect {
	x1 := int(math.Max(float64(r.X), float64(t.X)))
	y1 := int(math.Max(float64(r.Y), float64(t.Y)))
	x2 := int(math.Min(float64(r.X+r.W), float64(t.X+t.W)))
	y2 := int(math.Min(float64(r.Y+r.H), float64(t.Y+t.H)))
	return &Rect{x1, y1, x2 - x1, y2 - y1}
}
