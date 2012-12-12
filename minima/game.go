// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/animal"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/item"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
	"encoding/json"
	"io"
	"math"
	"math/rand"
)

// TileSize is the width and height of a tile in pixels.
const TileSize = 32

type Game struct {
	wo         *world.World
	cam        ui.Camera
	base       Base
	Astro      *Player
	Herbivores []animal.Herbivores
	Treasure   []item.Treasure
}

// ReadGame returns a *Game, read from the given
// reader.
func ReadGame(r io.Reader) (*Game, error) {
	e := &Game{
		cam: ui.Camera{Dims: ScreenDims},
	}
	in := bufio.NewReader(r)
	var err error
	if e.wo, err = world.Read(in); err != nil {
		return nil, err
	}
	crashSite := geom.Pt(float64(e.wo.X0*TileSize), float64(e.wo.Y0*TileSize))
	e.Astro = NewPlayer(e.wo, crashSite)
	e.base = NewBase(crashSite)

	if err := json.NewDecoder(in).Decode(&e); err != nil {
		panic(err)
	}
	e.CenterOnTile(e.wo.Tile(e.Astro.body.Center()))
	return e, nil
}

func (e *Game) Transparent() bool {
	return false
}

// CenterOnTile centers the display on a given tile.
func (e *Game) CenterOnTile(x, y int) {
	e.cam.Center(geom.Pt(TileSize*float64(x)+TileSize/2,
		TileSize*float64(y)+TileSize/2))
}

func (e *Game) Draw(d ui.Drawer) {
	w, h := int(ScreenDims.X/TileSize), int(ScreenDims.Y/TileSize)
	x0, y0 := e.wo.Tile(e.cam.Pt)

	xoff0 := -math.Mod(e.cam.Pt.X, TileSize)
	if e.cam.Pt.X < 0 {
		xoff0 = -TileSize + xoff0
	}
	yoff0 := -math.Mod(e.cam.Pt.Y, TileSize)
	if e.cam.Pt.Y < 0 {
		yoff0 = -TileSize + yoff0
	}

	pt := geom.Pt(xoff0, yoff0)
	for x := x0; x <= x0+w; x++ {
		for y := y0; y <= y0+h; y++ {
			l := e.wo.At(x, y)
			drawCell(d, l, x, y, pt)
			pt.Y += TileSize
		}
		pt.Y = yoff0
		pt.X += TileSize
	}

	e.base.Draw(d, e.cam)
	for _, t := range e.Treasure {
		if t.Item == nil {
			continue
		}
		e.cam.Draw(d, ui.Sprite{
			Name:   "Present",
			Bounds: geom.Rect(0, 0, t.Box.Dx(), t.Box.Dy()),
			Shade:  1.0, //TODO: should shade with altitude
		}, t.Box.Min)
	}
	e.Astro.Draw(d, e.cam)
	for i := range e.Herbivores {
		e.Herbivores[i].Draw(d, e.cam)
	}

	e.Astro.drawO2(d)

	if !*debug {
		return
	}
	d.SetFont("prstartk", 14)
	d.SetColor(White)
	sz := d.TextSize(e.Astro.info)
	d.Draw(e.Astro.info, geom.Pt(0, ScreenDims.Y-sz.Y))
}

func drawCell(d ui.Drawer, l *world.Loc, x, y int, pt geom.Point) {
	const minSh = 0.15
	const slope = (1 - minSh) / world.MaxElevation

	d.Draw(ui.Sprite{
		Name:   l.Terrain.Name,
		Bounds: geom.Rect(0, 0, TileSize, TileSize),
		Shade:  slope*float32(l.Elevation-l.Depth) + minSh,
	}, pt)
}

func (ex *Game) Handle(stk *ui.ScreenStack, ev ui.Event) error {
	k, ok := ev.(ui.Key)
	if !ok || !k.Down {
		return nil
	}

	switch k.Button {
	case ui.Menu:
		stk.Push(NewPauseScreen(ex.Astro))
	case ui.Action:
		if ex.Astro.body.Box.Overlaps(ex.base.Box) {
			stk.Push(NewBaseScreen(ex.Astro, &ex.base))
		}
		for i, t := range ex.Treasure {
			if t.Item == nil || !ex.Astro.body.Box.Overlaps(t.Box) {
				continue
			}
			scr := NewNormalMessage("You don't have room for that in your pack.")
			if ex.Astro.PutPack(t.Item) {
				scr = NewNormalMessage("Bravo! You got the " + t.Item.Name + "!")
				ex.Treasure[i].Item = nil
			}
			stk.Push(scr)
			break
		}
	}
	return nil
}

func (e *Game) Update(stk *ui.ScreenStack) error {
	const speed = 4 // px

	if e.Astro.o2 == 0 {
		et := e.Astro.FindEtele()
		if et == nil || et.Uses == 0 {
			stk.Push(NewGameOverScreen())
		} else {
			et.Uses--
			e.Astro.body.Vel = geom.Pt(0, 0)
			dims := geom.Pt(e.Astro.body.Box.Dx(), e.Astro.body.Box.Dy())
			e.Astro.body.Box.Min = e.base.Box.Min
			e.Astro.body.Box.Max = e.base.Box.Min.Add(dims)
			e.Astro.RefillO2()
		}
	}

	e.Astro.body.Vel = geom.Pt(0, 0)
	if stk.Buttons&ui.Left != 0 {
		e.Astro.body.Vel.X -= speed
	}
	if stk.Buttons&ui.Right != 0 {
		e.Astro.body.Vel.X += speed
	}
	if stk.Buttons&ui.Down != 0 {
		e.Astro.body.Vel.Y += speed
	}
	if stk.Buttons&ui.Up != 0 {
		e.Astro.body.Vel.Y -= speed
	}
	e.Astro.Move(e.wo)
	e.cam.Center(e.Astro.body.Box.Center())

	for i := range e.Herbivores {
		ai.UpdateBoids(e.Herbivores[i], &e.Astro.body, e.wo)
		e.Herbivores[i].Move(e.wo)
	}

	return nil
}

func randPoint(xmax, ymax float64) geom.Point {
	return geom.Pt(rand.Float64()*xmax, rand.Float64()*ymax)
}
