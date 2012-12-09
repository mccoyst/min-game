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
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
)

// TileSize is the width and height of a tile in pixels.
const TileSize = 32

type Game struct {
	wo       *world.World
	cam      ui.Camera
	astro    *Player
	base     Base
	herbs    []animal.Herbivores
	treasure []Treasure
	// Keys is a bitmask of the currently pressed keys.
	keys ui.Button
}

// ReadGame reads an the items of an Game
// from the given reader and returns it, or an error if an error
// was encountered.
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
	e.astro = NewPlayer(e.wo, crashSite)
	e.base = NewBase(crashSite)

	for {
		var name string
		var size int64
		_, err = fmt.Fscanln(in, &name, &size)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		switch name {
		case "herbs":
			b, err := ioutil.ReadAll(&io.LimitedReader{in, size})
			if err != nil {
				return nil, err
			}
			var h animal.Herbivores
			if err = json.Unmarshal(b, &h); err != nil {
				return nil, err
			}
			e.herbs = append(e.herbs, h)

		default:
			panic("unknown input section: " + name)
		}
	}

	// for now
	e.treasure = []Treasure{Treasure{&item.Element{"Uranium"}, e.astro.body.Box.Add(geom.Pt(128, 128))}}

	e.CenterOnTile(e.wo.Tile(e.astro.body.Center()))

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
	for _, t := range e.treasure {
		if t.Item != nil {
			t.Draw(d, e.cam)
		}
	}
	e.astro.Draw(d, e.cam)
	for i := range e.herbs {
		e.herbs[i].Draw(d, e.cam)
	}

	e.astro.drawO2(d)

	if !*locInfo {
		return
	}
	d.SetFont("prstartk", 14)
	d.SetColor(White)
	sz := d.TextSize(e.astro.info)
	d.Draw(e.astro.info, geom.Pt(0, ScreenDims.Y-sz.Y))
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
		stk.Push(NewPauseScreen(ex.astro))
	case ui.Action:
		if ex.astro.body.Box.Overlaps(ex.base.Box) {
			stk.Push(NewBaseScreen(ex.astro, &ex.base))
		}
		for i, t := range ex.treasure {
			if t.Item != nil && ex.astro.body.Box.Overlaps(t.Box) {
				if ex.astro.PutPack(t.Item) {
					stk.Push(NewTreasureGet(t.Item.Name()))
					ex.treasure[i].Item = nil
					break
				} else {
					stk.Push(NewNormalMessage("You don't have room for that in your pack."))
					break
				}
			}
		}
	}
	return nil
}

func (e *Game) Update(stk *ui.ScreenStack) error {
	const speed = 4 // px

	if e.astro.o2 == 0 {
		et := e.astro.FindEtele()
		if et == nil || !et.Use() {
			stk.Push(NewGameOverScreen())
		} else {
			e.astro.body.Vel = geom.Pt(0, 0)
			dims := geom.Pt(e.astro.body.Box.Dx(), e.astro.body.Box.Dy())
			e.astro.body.Box.Min = e.base.Box.Min
			e.astro.body.Box.Max = e.base.Box.Min.Add(dims)
			e.astro.RefillO2()
		}
	}

	e.astro.body.Vel = geom.Pt(0, 0)
	if stk.Buttons&ui.Left != 0 {
		e.astro.body.Vel.X -= speed
	}
	if stk.Buttons&ui.Right != 0 {
		e.astro.body.Vel.X += speed
	}
	if stk.Buttons&ui.Down != 0 {
		e.astro.body.Vel.Y += speed
	}
	if stk.Buttons&ui.Up != 0 {
		e.astro.body.Vel.Y -= speed
	}
	e.astro.Move(e.wo)
	e.cam.Center(e.astro.body.Box.Center())

	for i := range e.herbs {
		ai.UpdateBoids(e.herbs[i], &e.astro.body, e.wo)
		e.herbs[i].Move(e.wo)
	}

	return nil
}

func randPoint(xmax, ymax float64) geom.Point {
	return geom.Pt(rand.Float64()*xmax, rand.Float64()*ymax)
}
