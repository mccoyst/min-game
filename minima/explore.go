// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ai"
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
	gulls animal.Gulls
	cows  animal.Cows

	// Keys is a bitmask of the currently pressed keys.
	keys ui.Button
}

func NewExploreScreen(wo *world.World) *ExploreScreen {
	e := &ExploreScreen{
		wo:  wo,
		cam: ui.Camera{Dims: ScreenDims},
	}
	e.CenterOnTile(wo.X0, wo.Y0)
	e.astro = NewPlayer(e.wo, geom.Pt(float64(wo.X0*TileSize), float64(wo.Y0*TileSize)))

	xmin, xmax := float64(wo.X0-8)*TileSize, float64(wo.X0+8)*TileSize
	ymin, ymax := float64(wo.Y0-8)*TileSize, float64(wo.Y0+8)*TileSize
	for i := 0; i < 25; i++ {
		x := rand.Float64()*(xmax-xmin) + xmin
		y := rand.Float64()*(ymax-ymin) + ymin
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		e.gulls = append(e.gulls, animal.NewGull(geom.Pt(x, y), vel))
	}

	for i := 0; i < 25; i++ {
		var x, y float64
		for i := 0; i < 1000; i++ {
			x = rand.Float64()*(xmax-xmin) + xmin
			y = rand.Float64()*(ymax-ymin) + ymin
			pt := geom.Pt(x, y)
			tx, ty := e.wo.Tile(pt.Add(world.TileSize.Div(2)))
			if e.wo.At(tx, ty).Terrain.Char == "g" {
				break
			}
		}
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		e.cows = append(e.cows, animal.NewCow(geom.Pt(x, y), vel))
	}
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

func (e *ExploreScreen) Draw(d Drawer) error {
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

	if err := e.astro.Draw(d, e.cam); err != nil {
		return err
	}
	for _, g := range e.gulls {
		if err := g.Draw(d, e.cam); err != nil {
			return err
		}
	}
	for _, c := range e.cows {
		if err := c.Draw(d, e.cam); err != nil {
			return err
		}
	}

	if !*locInfo {
		return nil
	}
	if err := d.SetFont("prstartk", 14); err != nil {
		return err
	}
	d.SetColor(White)
	if _, err := d.Draw(e.astro.info, geom.Pt(0, 0)); err != nil {
		return err
	}
	return nil
}

func drawCell(d Drawer, l *world.Loc, x, y int, pt geom.Point) error {
	const minSh = 0.15
	const slope = (1 - minSh) / world.MaxElevation

	_, err := d.Draw(ui.Sprite{
		Name:   l.Terrain.Name,
		Bounds: geom.Rect(0, 0, TileSize, TileSize),
		Shade:  slope*float32(l.Elevation-l.Depth) + minSh,
	}, pt)

	return err
}

func (ex *ExploreScreen) Handle(stk *ScreenStack, ev ui.Event) error {
	if k, ok := ev.(ui.Key); ok && k.Down && k.Button == ui.Menu {
		stk.Push(NewPauseScreen())
	}
	return nil
}

func (e *ExploreScreen) Update(stk *ScreenStack) error {
	const speed = 4 // px

	e.astro.body.Vel = geom.Pt(0, 0)
	if stk.buttons&ui.Left != 0 {
		e.astro.body.Vel.X -= speed
	}
	if stk.buttons&ui.Right != 0 {
		e.astro.body.Vel.X += speed
	}
	if stk.buttons&ui.Down != 0 {
		e.astro.body.Vel.Y += speed
	}
	if stk.buttons&ui.Up != 0 {
		e.astro.body.Vel.Y -= speed
	}
	e.astro.Move(e.wo)
	e.cam.Center(e.astro.body.Box.Center())

	ai.UpdateBoids(e.gulls, &e.astro.body, e.wo)
	for _, g := range e.gulls {
		g.Move(e.wo)
	}

	ai.UpdateBoids(e.cows, &e.astro.body, e.wo)
	for _, c := range e.cows {
		c.Move(e.wo)
	}

	return nil
}

func randPoint(xmax, ymax float64) geom.Point {
	return geom.Pt(rand.Float64()*xmax, rand.Float64()*ymax)
}
