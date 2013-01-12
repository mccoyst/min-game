// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package geom

import (
	"math"
)

// A Torus represents a torus with the given width and height.
type Torus struct {
	W, H float64
}

// SqDist returns the squared distance of two points on the torus.
func (t Torus) SqDist(a, b Point) float64 {
	a, b = t.Norm(a), t.Norm(b)
	dx := math.Abs(a.X - b.X)
	dx = math.Min(t.W-dx, dx)
	dy := math.Abs(a.Y - b.Y)
	dy = math.Min(t.H-dy, dy)
	return dx*dx + dy*dy
}

// Dist returns the distance of two points on the torus.
func (t Torus) Dist(a, b Point) float64 {
	return math.Sqrt(t.SqDist(a, b))
}

// Sub returns vector a - b respecting the torus.  For example,
// with a 100 by 100 torus, (99,99) - (0,0) points toward
// the negative x and y directions since (0,0) == (100,100).
func (t Torus) Sub(a, b Point) Point {
	a, b = t.normPair(a, b)
	return a.Sub(b)
}

// Norm returns a point that is equivalent to a on the torus, but
// is within (0,0)-(W-1, H-1).
func (t Torus) Norm(a Point) Point {
	return Pt(wrap(a.X, t.W), wrap(a.Y, t.H))
}

// NormPair returns a pair of points where the first is equavalent
// to a on the torus, but within (0,0)-(W-1, H-1), and the second
// is the point equivalent to b that is nearest a.
func (t Torus) normPair(a, b Point) (Point, Point) {
	a = t.Norm(a)
	b.X = nearWrap(a.X, b.X, t.W)
	b.Y = nearWrap(a.Y, b.Y, t.H)
	return a, b
}

// NormRect returns a rectangle equivalent to r on the torus,
// but with its minimum point normalized.
func (t Torus) NormRect(r Rectangle) Rectangle {
	sz := r.Size()
	r.Min = t.Norm(r.Min)
	r.Max = r.Min.Add(sz)
	return r
}

// Intersect returns the intersection of a and b on the torus.
//
// Both a and b must be smaller than the size of the torus.
func (t Torus) Intersect(a, b Rectangle) Rectangle {
	a, b = t.AlignRects(a, b)
	return a.Intersect(b)
}

// Overlaps returns true if the two rectangles intersect on the torus.
//
// Both a and b must be smaller than the size of the torus.
func (t Torus) Overlaps(a, b Rectangle) bool {
	return t.Intersect(a, b) != Rectangle{}
}

// AlignRects returns a pair of rectangles.  The first is the normalized equivalent
// of a.  The second is equivalent to b on the torus.  If there is b equivalent
// that overlaps a on the torus, then it is returned, otherwise it is any b
// equivalent.
//
// Both a and b are assumed to be smaller than the width and
// height of the torus.
func (t Torus) AlignRects(a, b Rectangle) (Rectangle, Rectangle) {
	if a.Dx() >= t.W || a.Dy() >= t.H {
		panic("Torus.align, rectangle a is too big")
	} else if b.Dx() >= t.W || b.Dy() >= t.H {
		panic("Torus.align, rectangle b is too big")
	}

	a = t.NormRect(a)

	sz := b.Size()
	b.Min.X = wrap(b.Min.X, t.W)
	b.Max.X = b.Min.X + sz.X
	if !b.OverlapsX(a) {
		b.Max.X = wrap(b.Max.X, t.W)
		b.Min.X = b.Max.X - sz.X
	}

	b.Min.Y = wrap(b.Min.Y, t.H)
	b.Max.Y = b.Min.Y + sz.Y
	if !b.OverlapsY(a) {
		b.Max.Y = wrap(b.Max.Y, t.H)
		b.Min.Y = b.Max.Y - sz.Y
	}
	return a, b
}

// NearWrap returns the value equivalent to b modulo the width
// that is nearest to a, which is assumed to be within 0-width-1.
func nearWrap(a, b, width float64) float64 {
	if a < 0 || a >= width {
		panic("geom.nearWrap: a is out of range")
	}

	b = wrap(b, width)
	n := b
	if math.Abs(b+width-a) < math.Abs(n-a) {
		n = b + width
	}
	if math.Abs(b-width-a) < math.Abs(n-a) {
		n = b - width
	}
	return n
}

// Wrap returns x wrapped around at bound in both the positive
// and negative direction.
func wrap(x, bound float64) float64 {
	if x >= 0 && x < bound {
		return x
	}
	if x = math.Mod(x, bound); x < 0 {
		return bound + x
	}
	return x
}
