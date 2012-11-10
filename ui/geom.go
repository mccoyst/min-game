// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ui

import (
	"strconv"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y float64
}

// String returns a string representation of p like "(3.00,4.01)".
func (p Point) String() string {
	return "(" + strconv.FormatFloat(p.X, 'f', -1, 64) + "," + strconv.FormatFloat(p.Y, 'f', -1, 64) + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Eq returns whether p and q are equal.
func (p Point) Eq(q Point) bool {
	return p.X == q.X && p.Y == q.Y
}

// Pt is shorthand for Point{X, Y}.
func Pt(X, Y float64) Point {
	return Point{X, Y}
}

// A Rectangle contains the points with Min.X <= X < Max.X, Min.Y <= Y < Max.Y.
// It is well-formed if Min.X <= Max.X and likewise for Y. Points are always
// well-formed. A rectangle's methods always return well-formed outputs for
// well-formed inputs.
type Rectangle struct {
	Min, Max Point
}

// String returns a string representation of r like "(3,4)-(6,5)".
func (r Rectangle) String() string {
	return r.Min.String() + "-" + r.Max.String()
}

// Dx returns r's width.
func (r Rectangle) Dx() float64 {
	return r.Max.X - r.Min.X
}

// Dy returns r's height.
func (r Rectangle) Dy() float64 {
	return r.Max.Y - r.Min.Y
}

// Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Point{r.Dx(), r.Dy()}
}

// Add returns the rectangle r translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X + p.X, r.Min.Y + p.Y},
		Point{r.Max.X + p.X, r.Max.Y + p.Y},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle) Sub(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X - p.X, r.Min.Y - p.Y},
		Point{r.Max.X - p.X, r.Max.Y - p.Y},
	}
}

// Eq returns whether r and s are equal.
func (r Rectangle) Eq(s Rectangle) bool {
	return r.Min.X == s.Min.X && r.Min.Y == s.Min.Y &&
		r.Max.X == s.Max.X && r.Max.Y == s.Max.Y
}

// Overlaps returns whether r and s have a non-empty intersection.
func (r Rectangle) Overlaps(s Rectangle) bool {
	return r.Min.X < s.Max.X && s.Min.X < r.Max.X &&
		r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
}

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}.
func Rect(x0, y0, x1, y1 float64) Rectangle {
	if x0 > x1 {
		x0, x1 = x1, x0
	}
	if y0 > y1 {
		y0, y1 = y1, y0
	}
	return Rectangle{Point{x0, y0}, Point{x1, y1}}
}