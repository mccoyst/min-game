// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/animal"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
	"math"
	"math/rand"
)

// TileSize is the width and height of a tile in pixels.
const TileSize = 32

type ExploreScreen struct {
	wo    *world.World
	cam   ui.Camera
	astro *Player
	base  Base
	anims animal.Animals
	// Keys is a bitmask of the currently pressed keys.
	keys ui.Button
}

func NewExploreScreen(wo *world.World, anims animal.Animals) *ExploreScreen {
	e := &ExploreScreen{
		wo:  wo,
		cam: ui.Camera{Dims: ScreenDims},
	}
	e.CenterOnTile(wo.X0, wo.Y0)
	crashSite := geom.Pt(float64(wo.X0*TileSize), float64(wo.Y0*TileSize))
	e.astro = NewPlayer(e.wo, crashSite)
	e.base = NewBase(crashSite)
	e.anims = anims
	return e
}

func (e *ExploreScreen) Transparent() bool {
	return false
}

// CenterOnTile centers the display on a given tile.
func (e *ExploreScreen) CenterOnTile(x, y int) {
	e.cam.Center(geom.Pt(TileSize*float64(x)+TileSize/2,
		TileSize*float64(y)+TileSize/2))
}

func (e *ExploreScreen) Draw(d ui.Drawer) error {
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
			err := drawCell(d, l, x, y, pt)
			if err != nil {
				return err
			}
			pt.Y += TileSize
		}
		pt.Y = yoff0
		pt.X += TileSize
	}

	if err := e.base.Draw(d, e.cam); err != nil {
		return err
	}

	if err := e.astro.Draw(d, e.cam); err != nil {
		return err
	}
	if err := e.anims.Draw(d, e.cam); err != nil {
		return err
	}

	if !*locInfo {
		return nil
	}
	d.SetFont("prstartk", 14)
	d.SetColor(White)
	if _, err := d.Draw(e.astro.info, geom.Pt(0, 0)); err != nil {
		return err
	}
	return nil
}

func drawCell(d ui.Drawer, l *world.Loc, x, y int, pt geom.Point) error {
	const minSh = 0.15
	const slope = (1 - minSh) / world.MaxElevation

	_, err := d.Draw(ui.Sprite{
		Name:   l.Terrain.Name,
		Bounds: geom.Rect(0, 0, TileSize, TileSize),
		Shade:  slope*float32(l.Elevation-l.Depth) + minSh,
	}, pt)

	return err
}

func (ex *ExploreScreen) Handle(stk *ui.ScreenStack, ev ui.Event) error {
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
	}
	return nil
}

func (e *ExploreScreen) Update(stk *ui.ScreenStack) error {
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

	e.anims.Move(&e.astro.body, e.wo)

	return nil
}

func randPoint(xmax, ymax float64) geom.Point {
	return geom.Pt(rand.Float64()*xmax, rand.Float64()*ymax)
}
