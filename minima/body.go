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

	box1 := b.Box.Add(b.Vel)
	const MaxDh = 0.6
	dh := math.Min(math.Abs(meanHeight(box1, w)-meanHeight(b.Box, w)), MaxDh)
	step := 1.0 / 16
	slope := (step - 1.0) / MaxDh
	scale := dh*slope + 1.0
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
	wx := int(math.Floor(bc.X / TileSize))
	wy := int(math.Floor(bc.Y / TileSize))

	sum, area := 0.0, 0.0
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			l := w.At(wx+dx, wy+dy)
			x, y := float64(l.X*TileSize), float64(l.Y*TileSize)
			lbox := ui.Rect(x, y, x+TileSize, y+TileSize)
			is := IsectWorld(w, box, lbox)
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

func IsectWorld(w *world.World, a, b ui.Rectangle) ui.Rectangle {
	wSz := ui.Pt(float64(w.W)*TileSize, float64(w.H)*TileSize)

	if !Wraps(a, wSz) && !Wraps(b, wSz) {
		return a.Intersect(b)
	}

	if Wraps(b, wSz) {
		b = WrapMin(b, wSz)
	}

	dx := a.Dx()
	if a.Min.X < 0 {
		a.Min.X = wSz.X + math.Mod(a.Min.X, wSz.X)
	} else if a.Min.X >= wSz.X {
		a.Min.X = math.Mod(a.Min.X, wSz.X)
	}
	a.Max.X = a.Min.X + dx

	if a.Min.X >= b.Max.X || a.Max.X < b.Min.X {
		if a.Max.X < 0 {
			a.Max.X = wSz.X - math.Mod(a.Max.X, wSz.X)
		} else if a.Max.X >= wSz.X {
			a.Max.X = math.Mod(a.Max.X, wSz.X)
		}
		a.Min.X = a.Max.X - dx
	}
	if a.Min.X >= b.Max.X || a.Max.X < b.Min.X {
		return ui.Rectangle{}
	}

	dy := a.Dy()
	if a.Min.Y < 0 {
		a.Min.Y = wSz.Y + math.Mod(a.Min.Y, wSz.Y)
	} else if a.Min.Y >= wSz.Y {
		a.Min.Y = math.Mod(a.Min.Y, wSz.Y)
	}
	a.Max.Y = a.Min.Y + dy

	if a.Min.Y >= b.Max.Y || a.Max.Y < b.Min.Y {
		if a.Max.Y < 0 {
			a.Max.Y = wSz.Y - math.Mod(a.Max.Y, wSz.Y)
		} else if a.Max.Y >= wSz.Y {
			a.Max.Y = math.Mod(a.Max.Y, wSz.Y)
		}
		a.Min.Y = a.Max.Y - dy
	}
	return a.Intersect(b)
}

func Wraps(r ui.Rectangle, sz ui.Point) bool {
	if r.Min.X < 0 || r.Min.Y < 0 || r.Min.X >= sz.X || r.Min.Y >= sz.Y {
		return true
	}

	max := r.Min.Add(sz)
	return max.X < 0 || max.Y < 0 || max.X >= sz.X || max.Y >= sz.Y
}

func WrapMin(r ui.Rectangle, sz ui.Point) ui.Rectangle {
	if r.Min.X >= 0 && r.Min.Y >= 0 && r.Min.X < sz.X && r.Min.Y < sz.Y {
		return r
	}

	dx := r.Dx()
	dy := r.Dy()

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

	r.Max.X = r.Min.X + dx
	r.Max.Y = r.Min.Y + dy

	return r
}
