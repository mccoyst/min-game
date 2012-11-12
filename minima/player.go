// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Player struct {
	wo   *world.World
	body Body
	face, frame int
	ticks int
}

const Tempo = 40
var frames [][]ui.Rectangle

func init() {
	// TODO(mccoyst): Read this info from a file
	for y := 0; y < 4; y++ {
		frames = append(frames, make([]ui.Rectangle, 2))
		for x := 0; x < 2; x++ {
			frames[y][x] = ui.Rect(float64(x*TileSize), float64(y*TileSize), float64(x*TileSize + TileSize), float64(y*TileSize + TileSize))
		}
	}
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
	p.ticks++
	if p.ticks >= Tempo {
		p.frame++
		if p.frame >= 2 {
			p.frame = 0
		}
		p.ticks = 0
	}

	// TODO(mccoyst): read from the same file, yadda yadda
	if p.body.Vel.Y > 0 {
		p.face = 0
	}
	if p.body.Vel.Y < 0 {
		p.face = 3
	}
	if p.body.Vel.X > 0 {
		p.face = 2
	}
	if p.body.Vel.X < 0 {
		p.face = 1
	}

	p.body.Move(w)
}

func (p *Player) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   "Astronaut",
		Bounds: frames[p.face][p.frame],
		Shade:  1.0,
	}, p.body.Box.Min)
	return err
}
