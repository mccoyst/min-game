// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Gull struct {
	body Body
}

var gullScales = map[rune]float64{
	'g': 1.0,
	'f': 1.0,
	'm': 1.0,
	'w': 1.0,
	'd': 1.0,
	'i': 0.1,
}

func NewGull(p, v ui.Point) *Gull {
	return &Gull{
		Body{
			Box: ui.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
			Vel: v,
		},
	}
}

func (g *Gull) Move(w *world.World) {
	g.body.Move(w, gullScales)
}

func (g *Gull) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name: "Bird0",
		Bounds: ui.Rect(0, 0, TileSize, TileSize),
		Shade: 1.0,
	}, g.body.Box.Min)
	return err
}
