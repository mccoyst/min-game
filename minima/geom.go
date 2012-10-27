// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

// A Point is an X, Y coordinate pair. The axes increase right and down.
type Point struct {
	X, Y Fxpt
}

// String returns a string representation of p like "(3.00,4.01)".
func (p Point) String() string {
	return "(" + p.X.String() + "," + p.Y.String() + ")"
}

// Add returns the vector p+q.
func (p Point) Add(q Point) Point {
	return Point{p.X.Add(q.X), p.Y.Add(q.Y)}
}

// Sub returns the vector p-q.
func (p Point) Sub(q Point) Point {
	return Point{p.X.Sub(q.X), p.Y.Sub(q.Y)}
}

// Eq returns whether p and q are equal.
func (p Point) Eq(q Point) bool {
	return p.X.Eq(q.X) && p.Y.Eq(q.Y)
}

// Pt is shorthand for Point{X, Y}.
func Pt(X, Y Fxpt) Point {
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
func (r Rectangle) Dx() Fxpt {
	return r.Max.X.Sub(r.Min.X)
}

// Dy returns r's height.
func (r Rectangle) Dy() Fxpt {
	return r.Max.Y.Sub(r.Min.Y)
}

// Size returns r's width and height.
func (r Rectangle) Size() Point {
	return Point{r.Dx(), r.Dy()}
}

// Add returns the rectangle r translated by p.
func (r Rectangle) Add(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X.Add(p.X), r.Min.Y.Add(p.Y)},
		Point{r.Max.X.Add(p.X), r.Max.Y.Add(p.Y)},
	}
}

// Sub returns the rectangle r translated by -p.
func (r Rectangle) Sub(p Point) Rectangle {
	return Rectangle{
		Point{r.Min.X.Sub(p.X), r.Min.Y.Sub(p.Y)},
		Point{r.Max.X.Sub(p.X), r.Max.Y.Sub(p.Y)},
	}
}

// Eq returns whether r and s are equal.
func (r Rectangle) Eq(s Rectangle) bool {
	return r.Min.X == s.Min.X && r.Min.Y == s.Min.Y &&
		r.Max.X == s.Max.X && r.Max.Y == s.Max.Y
}

// Overlaps returns whether r and s have a non-empty intersection.
func (r Rectangle) Overlaps(s Rectangle) bool {
	return r.Min.X.Lt(s.Max.X) && s.Min.X.Lt(r.Max.X) &&
		r.Min.Y.Lt(s.Max.Y) && s.Min.Y.Lt(r.Max.Y)
}

// Rect is shorthand for Rectangle{Pt(x0, y0), Pt(x1, y1)}.
func Rect(x0, y0, x1, y1 Fxpt) Rectangle {
	if x0.Gt(x1) {
		x0, x1 = x1, x0
	}
	if y0.Gt(y1) {
		y0, y1 = y1, y0
	}
	return Rectangle{Point{x0, y0}, Point{x1, y1}}
}
