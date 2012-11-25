// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"fmt"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/sprite"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Player struct {
	wo   *world.World
	body phys.Body

	// TileX and tileX give the coordinates of the player's current tile. 
	tileX, tileY int

	// Info is a string describing the player's current tile.  It is used
	// for debugging purposes.
	info string

	anim sprite.Anim
}

var astroSheet sprite.Sheet

var baseScales = map[string]float64{
	"g": 1.0,
	"f": 0.85,
	"m": 0.5,
	"w": 0.1,
	"d": 0.75,
	"i": 0.4,
}

func init() {
	var err error
	astroSheet, err = sprite.LoadSheet("Astronaut")
	if err != nil {
		panic(err)
	}
}

func NewPlayer(wo *world.World, p geom.Point) *Player {
	return &Player{
		wo: wo,
		body: phys.Body{
			Box: geom.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
		},
	}
}

func (p *Player) Move(w *world.World) {
	p.anim.Move(&astroSheet, p.body.Vel)
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

func (p *Player) Draw(d Drawer, cam ui.Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   astroSheet.Name,
		Bounds: astroSheet.Frame(p.anim.Face, p.anim.Frame),
		Shade:  1.0,
	}, p.body.Box.Min)
	return err
}
