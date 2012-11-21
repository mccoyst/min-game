// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
	"fmt"
)

type Player struct {
	wo   *world.World
	body Body

	// TileX and tileX give the coordinates of the player's current tile. 
	tileX, tileY int

	// Info is a string describing the player's current tile.  It is used
	// for debugging purposes.
	info string

	face, frame int
	ticks       int
}

var astroSheet SpriteSheet

var baseScales = map[rune]float64{
	'g': 1.0,
	'f': 0.85,
	'm': 0.5,
	'w': 0.1,
	'd': 0.75,
	'i': 0.4,
}

func init() {
	var err error
	astroSheet, err = LoadSpriteSheet("Astronaut")
	if err != nil {
		panic(err)
	}
}

func NewPlayer(wo *world.World, p geom.Point) *Player {
	return &Player{
		wo: wo,
		body: Body{
			Box: geom.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
		},
	}
}

func (p *Player) Move(w *world.World) {
	p.ticks++
	if p.ticks >= astroSheet.Tempo {
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

	p.body.Move(w, baseScales)

	if !*locInfo {
		return
	}
	tx, ty := w.Tile(p.body.Center())
	if tx == p.tileX && ty == p.tileY {
		return
	}
	p.tileX = tx
	p.tileY = ty
	p.info = fmt.Sprintf("%d,%d: %s", tx, ty, w.At(tx, ty).Terrain.Name)
}

func (p *Player) Draw(d Drawer, cam Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   astroSheet.Name,
		Bounds: astroSheet.Frame(p.face, p.frame),
		Shade:  1.0,
	}, p.body.Box.Min)
	return err
}
