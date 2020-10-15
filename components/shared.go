package comp

import (
	"image/color"
	"math"
	"time"

	"github.com/hajimehoshi/ebiten"
)

// common used variables in components, values change between levels
var (
	WP WorldProperties
	LP LevelProperties
)

// common used constants in components
const (
	PI       float64 = math.Pi
	DPI      float64 = math.Pi * 2
	HPI      float64 = math.Pi / 2
	RadToDeg float64 = 180 / math.Pi
	DegToRad float64 = math.Pi / 180
)

// WorldProperties are properties of the world or level
type WorldProperties struct {
	Gravity     float64
	Friction    float64
	LevelWidth  int
	LevelHeight int
}

// LevelProperties are used to store info / progress off the level
type LevelProperties struct {
	Lives        int
	CurrentLap   int
	MaxLaps      int
	LapTimes     []time.Duration
	LapStartTime time.Time
	PX, PY       int
}

// Drawer can be drawn every frame
type Drawer interface {
	Draw(screen *ebiten.Image) error
}

// Updater can be updated every frame
type Updater interface {
	Update(screen *ebiten.Image) error
}

// HitAble something that can be hit / collided with
type HitAble interface {
	SetHit(collider Collider)
	GetObject() *Object
}

// Collider checks collisions with HitAble's
type Collider interface {
	Collide(hitList []HitAble) error
	GetObject() *Object
}

// Vector used for direction of objects
type Vector struct {
	x, y float64
}

// NewVector creates a Vector
func NewVector(x, y float64) Vector {
	return Vector{x, y}
}

//Rect as format x,y,w,h
type Rect struct {
	x, y, w, h int
}

//NewRect creates a Rect
func NewRect(x, y, w, h int) Rect{
	return Rect{x,y,w,h}
}

// setXY position of Rect
func (r *Rect) setXY(x, y int) {
	r.x = x
	r.y = y
}

// HitShape is used for custom hit area for objects
type HitShape struct {
	rx, ry    int
	rect      Rect
	rectImg   *ebiten.Image
	rectColor color.RGBA
	Hit       bool
	solid     bool
}

// Controls stuff (TODO should not be part of object)
type Controls struct {
	up, down, left, right, rr, rl bool
}

// Sides has a boolean for each side, to indicate collision on specific side('s)
type Sides struct {
	left, right, top, bottom bool
}

// hitHV returs true if there is a horizontal and vertical hit
func (s *Sides) hitHV() bool {
	return (s.left || s.right) && (s.top || s.bottom)
}

// CheckHit checks for a hit between two objects and
// optionally returns the hit side('s) from which we hit the targets
func CheckHit(o *Object, target *Object, resolveSides bool, resolveHV bool) (bool, Sides) {
	r, t := &o.rect, &target.rect
	var sides Sides
	if CheckOverlap(r, t) {
		if resolveSides {
			checkSides(&sides, r, t)
			if resolveHV {
				resolveHitHV(&sides, getIntersetRect(r, t))
				//fmt.Printf("%s hits %s on sides %+v\n", con.ID[o.ID], con.ID[target.ID], sides)
			}
		}
		return true, sides
	}
	return false, sides
}

// CheckOverlap checks for a hit between two rects
func CheckOverlap(r *Rect, t *Rect) bool {
	if !(r.x > t.x+t.w || r.x+r.w < t.x || r.y > t.y+t.h || r.y+r.h < t.y) {
		return true
	}
	return false
}

// checkSides returns which side('s) we hit a target (assumes an overlap)
func checkSides(sides *Sides, r *Rect, t *Rect) {
	// dont remove equals signs, else e.g if the collider is bigger then the target
	// a hit bottom could return only hit Left, right as true
	// most likely rounding between position float and integers of the hit rect
	aLeft, aRight := r.x, r.x+r.w
	aTop, aBottom := r.y, r.y+r.h
	bLeft, bRight := t.x, t.x+t.w
	bTop, bBottom := t.y, t.y+t.h
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

// resolveHitHV resolves collisions in case we have a horizontal and vertical collision
// with the same object, based on the overlap dimensions it ignores the horizontal or vertical hit.
// so the collider only react to the biggest hit
func resolveHitHV(sides *Sides, ir *Rect) {
	//fmt.Println("ResolveHitHV, intersect rect: ", ir)
	if sides.hitHV() {
		if ir.w > ir.h {
			sides.left, sides.right = false, false
		} else {
			sides.top, sides.bottom = false, false
		}
	}
}

// getIntersetRect returns the intersection rect between two rectanges
func getIntersetRect(r *Rect, t *Rect) *Rect {
	x1 := int(math.Max(float64(r.x), float64(t.x)))
	y1 := int(math.Max(float64(r.y), float64(t.y)))
	x2 := int(math.Min(float64(r.x+r.w), float64(t.x+t.w)))
	y2 := int(math.Min(float64(r.y+r.h), float64(t.y+t.h)))
	return &Rect{x1, y1, x2 - x1, y2 - y1}
}
