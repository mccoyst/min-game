// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"math"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Gull struct {
	body Body
	anim Anim
}

var gullSheet SpriteSheet
var gullScales = map[rune]float64{
	'g': 1.0,
	'f': 1.0,
	'm': 1.0,
	'w': 1.0,
	'd': 1.0,
	'i': 0.1,
}

func init() {
	var err error
	gullSheet, err = LoadSpriteSheet("Gull")
	if err != nil {
		panic(err)
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
	g.anim.Move(&gullSheet, g.body.Vel)
	g.body.Move(w, gullScales)
}

func (g *Gull) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   gullSheet.Name,
		Bounds: gullSheet.Frame(g.anim.face, g.anim.frame),
		Shade:  1.0,
	}, g.body.Box.Min)
	return err
}

type Gulls []*Gull

func (gs Gulls) Len() int {
	return len(gs)
}

func (gs Gulls) Boid(n int) Boid {
	return Boid{&gs[n].body}
}

func (Gulls) Info() BoidInfo {
	return BoidInfo{
		MaxVelocity: 2.0,
		LocalDist:   math.Pow(TileSize*10.0, 2),
		AvoidDist:   math.Pow(TileSize/2.0, 2),
		PlayerDist:  math.Pow(TileSize*3, 2),
		CenterBias:  0.05,
		MatchBias:   0.08,
		AvoidBias:   0.5,
		PlayerBias:  0.2,

		TerrainDist:  TileSize,
		TerrainBias:  0.5,
		AvoidTerrain: "i",
	}
}
