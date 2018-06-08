// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package phys

import (
	"github.com/mccoyst/min-game/geom"
	"github.com/mccoyst/min-game/world"
)

type Body struct {
	Vel geom.Point
	Box geom.Rectangle
}

func (b *Body) Move(w *world.World, velScale map[string]float64) {
	if b.Vel.X == 0 && b.Vel.Y == 0 {
		return
	}
	wx, wy := w.Tile(b.Center())
	m := velScale[w.At(wx, wy).Terrain.Char] * b.Vel.Len()
	maxVel := geom.Pt(m, m)
	b.Box = b.Box.Add(b.Vel.Normalize().Mul(maxVel))
	b.Box = w.Pixels.NormRect(b.Box)
}

func (b *Body) Center() geom.Point {
	return b.Box.Center()
}
