// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"fmt"
	"math"
	"os"

	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Body struct {
	Vel ui.Point
	Box ui.Rectangle
}

func (b *Body) Move(w *world.World) {
	if b.Vel.X == 0 && b.Vel.Y == 0 {
		return
	}

	// Scale down the velocity based on the terrain.
	wx, wy := point2Tile(b.Box.Center())
	maxVel := velScale[w.At(wx, wy).Terrain.Char] * b.Vel.Len()

	const (
		maxDh  = 0.6
		minDh  = -maxDh
		minVel = 1.0 / 16.0
	)

	box1 := b.Box.Add(b.Vel)
	dh := avgElevation(box1, w) - avgElevation(b.Box, w)
	dh = math.Max(math.Min(dh, maxDh), minDh)

	slope := (minVel - maxVel) / maxDh
	b.Box = b.Box.Add(vecNorm(b.Vel, dh*slope+maxVel))
}

// VecNorm returns vec normalized to have the magnitude m.
func vecNorm(vec ui.Point, m float64) ui.Point {
	return vec.Mul(m / vec.Len())
}

var velScale = map[rune]float64{
	'g': 1.0,
	'f': 0.85,
	'm': 0.5,
	'w': 0.1,
	'd': 0.75,
	'i': 0.4,
}

// AvgElevation returns the average elevation of the world
// locations covered by a rectangle.
func avgElevation(box ui.Rectangle, w *world.World) float64 {
	sz := worldSize(w)
	wx, wy := point2Tile(box.Center())
	sum, area := 0.0, 0.0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			l := w.At(wx+dx, wy+dy)
			x, y := float64(l.X*TileSize), float64(l.Y*TileSize)
			lbox := ui.Rect(x, y, x+TileSize, y+TileSize)
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
func point2Tile(p ui.Point) (int, int) {
	return int(math.Floor(p.X / TileSize)), int(math.Floor(p.Y / TileSize))
}

// WorldSize returns the size of the world in pixels.
func worldSize(w *world.World) ui.Point {
	return ui.Pt(float64(w.W)*TileSize, float64(w.H)*TileSize)
}

// isectTorus returns the intersection between two rectangles around
// a torus of the given dimensions.
func isectTorus(a, b ui.Rectangle, sz ui.Point) ui.Rectangle {
	b = normRect(b, sz)
	a, ok := align(a, b, sz)
	if !ok {
		return ui.Rectangle{}
	}
	return a.Intersect(b)
}

// Align attempts to align the a with b (which must be normalized)
// around a torus so that they overlap. If they overlap then the
// aligned version of a is returned with true, otherwise false is returned.
func align(a, b ui.Rectangle, sz ui.Point) (ui.Rectangle, bool) {
	var ok bool
	if a, ok = alignX(a, b, sz.X); !ok {
		return a, ok
	}
	return alignY(a, b, sz.Y)
}

func alignX(a, b ui.Rectangle, width float64) (ui.Rectangle, bool) {
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

func alignY(a, b ui.Rectangle, height float64) (ui.Rectangle, bool) {
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
func normRect(r ui.Rectangle, worldSz ui.Point) ui.Rectangle {
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
