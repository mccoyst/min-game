// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"math"

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

// Point2Tile returns the tile coordinate for a point
func point2Tile(p geom.Point) (int, int) {
	return int(math.Floor(p.X / TileSize)), int(math.Floor(p.Y / TileSize))
}
