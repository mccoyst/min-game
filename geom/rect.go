// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package geom

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

// Center returns r's center point.
func (r Rectangle) Center() Point {
	return Point{r.Min.X + r.Dx()/2, r.Min.Y + r.Dy()/2}
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
	return r.OverlapsX(s) && r.OverlapsY(s)
}

// OverlapsX returns whether r and s have a non-empty intersection in the x direction.
func (r Rectangle) OverlapsX(s Rectangle) bool {
	return r.Min.X < s.Max.X && s.Min.X < r.Max.X
}

// OverlapsY returns whether r and s have a non-empty intersection in the y direction.
func (r Rectangle) OverlapsY(s Rectangle) bool {
	return r.Min.Y < s.Max.Y && s.Min.Y < r.Max.Y
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

// Intersect returns the largest rectangle contained by both r and s. If the
// two rectangles do not overlap then the zero rectangle will be returned.
func (r Rectangle) Intersect(s Rectangle) Rectangle {
	if r.Min.X < s.Min.X {
		r.Min.X = s.Min.X
	}
	if r.Min.Y < s.Min.Y {
		r.Min.Y = s.Min.Y
	}
	if r.Max.X > s.Max.X {
		r.Max.X = s.Max.X
	}
	if r.Max.Y > s.Max.Y {
		r.Max.Y = s.Max.Y
	}
	if r.Min.X > r.Max.X || r.Min.Y > r.Max.Y {
		return Rectangle{}
	}
	return r
}

// Rpad returns a rectangle with p removed from the left, right, top, and bottom of r.
func (r Rectangle) Rpad(p float64) Rectangle {
	r.Min.X += p
	r.Min.Y += p
	r.Max.X -= p
	r.Max.Y -= p
	return r
}
