// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Player struct {
	wo   *world.World
	body Body
}

func NewPlayer(wo *world.World, p ui.Point) *Player {
	return &Player{
		wo: wo,
		body: Body{
			Box: ui.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
		},
	}
}

func (p *Player) Move(w *world.World) {
	p.body.Move(w)
}

func (p *Player) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   "Astronaut",
		Bounds: ui.Rect(0, 0, TileSize, TileSize),
		Shade:  1.0,
	}, p.body.Box.Min)
	return err
}
