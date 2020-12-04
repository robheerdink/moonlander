package com

import "math"

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

// func GetFourRotatedPoints(cx, cy, ox, oy, rad float64) (x, y float64) {
// 	rx := cx + (ox * math.Cos(rad)) - (oy * math.Sin(rad))
// 	ry := cy + (ox * math.Sin(rad)) + (oy * math.Cos(rad))
// 	return rx, ry
// }
