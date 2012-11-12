// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"math"

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

	box1 := b.Box
	box1.Add(b.Vel)
	const MaxDh = 0.6
	dh := math.Min(math.Abs(meanHeight(box1, w) - meanHeight(b.Box, w)), MaxDh)
	step := 1.0/16
	slope := (step - 1.0) / MaxDh
	scale := dh * slope + 1.0
	bc := b.Box.Center()
	wx := int(bc.X / TileSize)
	wy := int(bc.Y / TileSize)
	v := b.Vel.Mul(scale).Mul(velScale[w.At(wx, wy).Terrain.Char])
	if b.Vel.X < 0 && v.X == 0 {
		v.X = -step
	}
	if b.Vel.X > 0 && v.X == 0 {
		v.X = step
	}
	if b.Vel.Y < 0 && v.Y == 0 {
		v.Y = -step
	}
	if b.Vel.Y > 0 && v.Y == 0 {
		v.Y = step
	}
	b.Box = b.Box.Add(v)
}

var velScale = map[rune]float64{
	'g': 1.0,
	'f': 0.85,
	'm': 0.5,
	'w': 0.1,
	'd': 0.75,
	'i': 0.4,
}

func meanHeight(box ui.Rectangle, w *world.World) float64 {
	bc := box.Center()
	wx := int(bc.X / TileSize)
	wy := int(bc.Y / TileSize)
	l := w.At(wx, wy)
	locs := make([]struct{ loc *world.Loc; box ui.Rectangle }, 9)
	i := 0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			locs[i].loc = w.At(l.X + dx, l.Y + dy)
			x, y := float64(l.X*TileSize), float64(l.Y*TileSize)
			locs[i].box = ui.Rect(x, y, x + TileSize, y + TileSize)
			i++
		}
	}

	zrect := ui.Rectangle{}
	sum, area := 0.0, 0.0
	for i := 0; i < 9; i++ {
		lbox := locs[i].box
		is := IsectWorld(w, box, lbox)
		if is.Eq(zrect) {
			continue
		}
		sum += float64(locs[i].loc.Elevation) * is.Dx() * is.Dy()
		area += is.Dx() * is.Dy()
	}

	return sum / area
}

func IsectWorld(w *world.World, a, b ui.Rectangle) ui.Rectangle {
	size := ui.Pt(float64(w.W)*TileSize, float64(w.H)*TileSize)

	if !Wraps(a, size) && !Wraps(b, size) {
		return a.Intersect(b)
	}

	zrect := ui.Rectangle{}

	if Wraps(b, size) {
		b = WrapMin(b, size)
		is := a.Intersect(b)
		if !is.Eq(zrect) {
			return is
		}

		b = WrapMax(b, size)
		return a.Intersect(b)
	}

	a = WrapMin(a, size)
	is := a.Intersect(b)
	if !is.Eq(zrect) {
		return is
	}

	a = WrapMax(a, size)
	return a.Intersect(b)
}

func Wraps(r ui.Rectangle, sz ui.Point) bool {
	if r.Min.X < 0 || r.Min.Y < 0 || r.Min.X >= sz.X || r.Min.Y >= sz.Y {
		return true
	}

	max := r.Min.Add(sz)
	return max.X < 0 && max.Y < 0 && max.X >= sz.X && max.Y >= sz.Y
}

func WrapMin(r ui.Rectangle, sz ui.Point) ui.Rectangle {
	if r.Min.X >= 0 && r.Min.Y >= 0 && r.Min.X < sz.X && r.Min.Y < sz.Y {
		return r
	}

	if r.Min.X < 0 {
		r.Min.X = sz.X + math.Mod(r.Min.X, sz.X)
	} else if r.Min.X >= sz.X {
		r.Min.X = math.Mod(r.Min.X, sz.X)
	}

	if r.Min.Y < 0 {
		r.Min.Y = sz.Y + math.Mod(r.Min.Y, sz.Y)
	} else if r.Min.Y >= sz.Y {
		r.Min.Y = math.Mod(r.Min.Y, sz.Y)
	}

	return r
}

func WrapMax(r ui.Rectangle, sz ui.Point) ui.Rectangle {
	max := r.Max

	if max.X >= 0 && max.Y >= 0 && max.X < sz.X && max.Y < sz.Y {
		return r
	}

	if max.X < 0 {
		max.X = sz.X + math.Mod(max.X, sz.X)
	} else if max.X >= sz.X {
		max.X = math.Mod(max.X, sz.X)
	}

	if max.Y < 0 {
		max.Y = sz.Y + math.Mod(max.Y, sz.Y)
	} else if max.Y >= sz.Y {
		max.Y = math.Mod(max.Y, sz.Y)
	}

	r.Min = max.Sub(sz)
	return r
}
