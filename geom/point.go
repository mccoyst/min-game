// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package geom

import (
	"fmt"
	"math"
)

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y float64
}

// String returns a string representation of p like "(3.00,4.01)".
func (p Point) String() string {
	return fmt.Sprintf("(%4.2f, %4.2f)", p.X, p.Y)
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X + q.X, p.Y + q.Y}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X - q.X, p.Y - q.Y}
}

// Mul returns the vector p*k.
func (p Point) Mul(k float64) Point {
	return Point{p.X * k, p.Y * k}
}

// Div returns the vector p/k.
func (p Point) Div(k float64) Point {
	return Point{p.X / k, p.Y / k}
}

// Len returns the magnitude of the vector defined by p.
func (p Point) Len() float64 {
	return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Dist returns the distance of two points.
func (p Point) Dist(q Point) float64 {
	dx, dy := p.X - q.X, p.Y - q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Normalize returns the vector normalized so that it has a magnitude of one.
func (p Point) Normalize() Point {
	return p.Div(p.Len())
}

// Eq returns whether p and q are equal.
func (p Point) Eq(q Point) bool {
	return p.X == q.X && p.Y == q.Y
}

// Pt is shorthand for Point{X, Y}.
func Pt(X, Y float64) Point {
	return Point{X, Y}
}
