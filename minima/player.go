// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"fmt"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/item"
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

	o2max   int
	o2      int
	o2ticks int

	suit Inventory
	pack Inventory
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
	if *debug {
		for i := range baseScales {
			baseScales[i] = 1.0
		}
	}
	return &Player{
		wo: wo,
		body: phys.Body{
			Box: geom.Rectangle{p, p.Add(TileSize)},
		},
		o2max: 50,
		o2:    50,
		suit:  Inventory{[]*item.Item{item.New(item.ETele), nil}, 0, true},
		pack:  Inventory{[]*item.Item{nil, nil, item.New(item.Uranium), nil}, -1, true},
	}
}

func (p *Player) Move(w *world.World) {
	p.o2ticks++
	if p.o2ticks > p.o2max && p.o2 > 0 {
		p.o2--
		p.o2ticks = 0
	}

	p.anim.Move(&astroSheet, p.body.Vel)
	p.body.Move(w, baseScales)

	if !*debug {
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

func (p *Player) Draw(d ui.Drawer, cam ui.Camera) {
	cam.Draw(d, ui.Sprite{
		Name:   astroSheet.Name,
		Bounds: astroSheet.Frame(p.anim.Face, p.anim.Frame),
		Shade:  1.0,
	}, p.body.Box.Min)
}

func (p *Player) RefillO2() {
	p.o2 = p.o2max
	p.o2ticks = 0
}

func (p *Player) drawO2(d ui.Drawer) {
	chunks := 10
	left := p.o2 / chunks
	chunk := geom.Rect(0, 0, 10, 10)

	dx, dy := 10.0, 10.0
	pt := geom.Pt(dx, dy)

	d.SetColor(Sky)
	i := 0
	for ; i < left; i++ {
		d.Draw(chunk, pt)
		pt.X += dx + 4
	}

	part := p.o2 % chunks
	if part != 0 {
		frac := float64(part) / float64(chunks)

		c := Sky
		c.R = uint8(float64(c.R) * frac)
		c.G = uint8(float64(c.G) * frac)
		c.B = uint8(float64(c.B) * frac)
		d.SetColor(c)

		d.Draw(chunk, pt)
	}
}

// FindEtele returns an E-Tele item with remaining uses from the player's suit or nil if such an item is not found.
func (p *Player) FindEtele() *item.Item {
	for _, i := range p.suit.Items {
		if i != nil && i.Name == item.ETele && i.Uses > 0 {
			return i
		}
	}
	return nil
}

// PutPack tries to add i to the player's backpack, and returns true iff successful.
func (p *Player) PutPack(i *item.Item) bool {
	return p.pack.Put(i)
}

// PutSuit tries to add i to the player's suit, and returns true iff successful.
func (p *Player) PutSuit(i *item.Item) bool {
	return p.suit.Put(i)
}
