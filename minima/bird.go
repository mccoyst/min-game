// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"math"

	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
)

type Gull struct {
	body Body

	face, frame int
	ticks       int
}

var gullFrames [][]geom.Rectangle
var gullScales = map[rune]float64{
	'g': 1.0,
	'f': 1.0,
	'm': 1.0,
	'w': 1.0,
	'd': 1.0,
	'i': 0.1,
}

func init() {
	// TODO(mccoyst): Read this info from a file
	for y := 0; y < 4; y++ {
		gullFrames = append(gullFrames, make([]geom.Rectangle, 2))
		for x := 0; x < 2; x++ {
			gullFrames[y][x] = geom.Rect(float64(x*TileSize), float64(y*TileSize), float64(x*TileSize+TileSize), float64(y*TileSize+TileSize))
		}
	}
}

func NewGull(p, v geom.Point) *Gull {
	return &Gull{
		body: Body{
			Box: geom.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
			Vel: v,
		},
	}
}

func (g *Gull) Body() *Body {
	return &g.body
}

func (g *Gull) Move(w *world.World) {
	g.ticks++
	if g.ticks >= Tempo {
		g.frame++
		if g.frame >= 2 {
			g.frame = 0
		}
		g.ticks = 0
	}

	dx, dy := g.body.Vel.X, g.body.Vel.Y
	vertBiased := math.Abs(dy) > math.Abs(dx)

	// TODO(mccoyst): read from the same file, yadda yadda
	if dy > 0 && vertBiased {
		g.face = 0
	}
	if dy < 0 && vertBiased {
		g.face = 1
	}
	if dx > 0 && !vertBiased {
		g.face = 3
	}
	if dx < 0 && !vertBiased {
		g.face = 2
	}

	g.body.Move(w, gullScales)
}

func (g *Gull) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   "Bird0",
		Bounds: gullFrames[g.face][g.frame],
		Shade:  1.0,
	}, g.body.Box.Min)
	return err
}
