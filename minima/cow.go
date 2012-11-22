// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Cow struct {
	body Body
	anim Anim
}

var cowSheet SpriteSheet
var cowScales = map[rune]float64{
	'g': 1.0,
	'f': 0.1,
	'm': 0.3,
	'w': 0.1,
	'd': 0.5,
	'i': 0.1,
}

func init() {
	var err error
	cowSheet, err = LoadSpriteSheet("Cow")
	if err != nil {
		panic(err)
	}
}

func NewCow(p, v geom.Point) *Cow {
	return &Cow{
		body: Body{
			Box: geom.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
			Vel: v,
		},
	}
}

func (c *Cow) Body() *Body {
	return &c.body
}

func (c *Cow) Move(w *world.World) {
	c.anim.Move(&cowSheet, c.body.Vel)
	c.body.Move(w, cowScales)
}

func (c *Cow) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   cowSheet.Name,
		Bounds: cowSheet.Frame(c.anim.face, c.anim.frame),
		Shade:  1.0,
	}, c.body.Box.Min)
	return err
}

type Cows []*Cow

func (cs Cows) Len() int {
	return len(cs)
}

func (cs Cows) Boid(n int) Boid {
	return Boid{&cs[n].body}
}

func (Cows) Info() BoidInfo {
	return BoidInfo{
		MaxVelocity: 0.5,
		LocalDist:   TileSize * 30.0,

		MatchBias: 0.0,

		CenterDist: TileSize * 15,
		CenterBias: 0.005,

		AvoidDist: TileSize * 1.5,
		AvoidBias: 0.001,

		PlayerDist: TileSize * 2,
		PlayerBias: 0.02,

		TerrainDist:  TileSize * 1.1,
		TerrainBias:  0.0005,
		AvoidTerrain: "fmwdi",
	}
}
