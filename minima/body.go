// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"fmt"
	"math"
	"os"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
)

type Body struct {
	Vel geom.Point
	Box geom.Rectangle
}

func (b *Body) Move(w *world.World, velScale map[rune]float64) {
	if b.Vel.X == 0 && b.Vel.Y == 0 {
		return
	}
	wx, wy := point2Tile(b.Box.Center())
	maxVel := velScale[w.At(wx, wy).Terrain.Char] * b.Vel.Len()
	b.Box = b.Box.Add(vecNorm(b.Vel, maxVel))
}

// VecNorm returns vec normalized to have the magnitude m.
func vecNorm(vec geom.Point, m float64) geom.Point {
	return vec.Mul(m / vec.Len())
}

// AvgElevation returns the average elevation of the world
// locations covered by a rectangle.
func avgElevation(box geom.Rectangle, w *world.World) float64 {
	sz := worldSize(w)
	wx, wy := point2Tile(box.Center())
	sum, area := 0.0, 0.0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			l := w.At(wx+dx, wy+dy)
			x, y := float64(l.X*TileSize), float64(l.Y*TileSize)
			lbox := geom.Rect(x, y, x+TileSize, y+TileSize)
			is := isectTorus(box, lbox, sz)
			sum += float64(l.Elevation) * is.Dx() * is.Dy()
			area += is.Dx() * is.Dy()
		}
	}
	if math.Abs(area-box.Dx()*box.Dy()) > 0.00001 {
		fmt.Printf("area=%g\n\n", area)
		os.Exit(1)
	}
	return sum / area
}

// Point2Tile returns the tile coordinate for a point
func point2Tile(p geom.Point) (int, int) {
	return int(math.Floor(p.X / TileSize)), int(math.Floor(p.Y / TileSize))
}

// WorldSize returns the size of the world in pixels.
func worldSize(w *world.World) geom.Point {
	return geom.Pt(float64(w.W)*TileSize, float64(w.H)*TileSize)
}

// isectTorus returns the intersection between two rectangles around
// a torus of the given dimensions.
func isectTorus(a, b geom.Rectangle, sz geom.Point) geom.Rectangle {
	b = normRect(b, sz)
	a, ok := align(a, b, sz)
	if !ok {
		return geom.Rectangle{}
	}
	return a.Intersect(b)
}

// Align attempts to align the a with b (which must be normalized)
// around a torus so that they overlap. If they overlap then the
// aligned version of a is returned with true, otherwise false is returned.
func align(a, b geom.Rectangle, sz geom.Point) (geom.Rectangle, bool) {
	var ok bool
	if a, ok = alignX(a, b, sz.X); !ok {
		return a, ok
	}
	return alignY(a, b, sz.Y)
}

func alignX(a, b geom.Rectangle, width float64) (geom.Rectangle, bool) {
	dx := a.Dx()

	a.Min.X = wrap(a.Min.X, width)
	a.Max.X = a.Min.X + dx
	if a.OverlapsX(b) {
		return a, true
	}

	a.Max.X = wrap(a.Max.X, width)
	a.Min.X = a.Max.X - dx
	return a, a.OverlapsX(b)
}

func alignY(a, b geom.Rectangle, height float64) (geom.Rectangle, bool) {
	dy := a.Dy()

	a.Min.Y = wrap(a.Min.Y, height)
	a.Max.Y = a.Min.Y + dy
	if a.OverlapsY(b) {
		return a, true
	}

	a.Max.Y = wrap(a.Max.Y, height)
	a.Min.Y = a.Max.Y - dy
	return a, a.OverlapsY(b)
}

// NormRect returns a normalized form of r that is wrapped around
// a torus with width and height given by worldSz such that its minimum
// point is within the rectangle (0,0), (worldSz.X-1,worldSz.Y-1),
func normRect(r geom.Rectangle, worldSz geom.Point) geom.Rectangle {
	if r.Min.X >= 0 && r.Min.Y >= 0 && r.Min.X < worldSz.X && r.Min.Y < worldSz.Y {
		return r
	}
	rectSz := r.Size()
	r.Min.X = wrap(r.Min.X, worldSz.X)
	r.Min.Y = wrap(r.Min.Y, worldSz.Y)
	r.Max = r.Min.Add(rectSz)
	return r
}

// Wrap returns x wrapped around at bound in both the positive
// and negative direction.
func wrap(x, bound float64) float64 {
	if x = math.Mod(x, bound); x < 0 {
		return bound + x
	}
	return x
}
