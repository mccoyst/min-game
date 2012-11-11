// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
	"math"
)

// TileSize is the width and height of a tile in pixels.
const TileSize = 32

type ExploreScreen struct {
	wo    *world.World
	cam   Camera
	astro *Player
}

func NewExploreScreen(wo *world.World) *ExploreScreen {
	e := &ExploreScreen{wo: wo}
	e.CenterOnTile(wo.X0, wo.Y0)
	e.astro = NewPlayer(e.wo, ui.Pt(float64(wo.X0*TileSize), float64(wo.Y0*TileSize)))
	return e
}

// CenterOnTile centers the display on a given tile.
func (e *ExploreScreen) CenterOnTile(x, y int) {
	e.cam.Center(ui.Pt(TileSize*float64(x)+TileSize/2,
		TileSize*float64(y)+TileSize/2))
}

func (e *ExploreScreen) Draw(d Drawer) error {
	e.cam.Center(e.astro.box.Center())

	w, h := int(ScreenDims.X/TileSize), int(ScreenDims.Y/TileSize)
	x0 := int(e.cam.pt.X / TileSize)
	xoff0 := -math.Mod(e.cam.pt.X, TileSize)
	if e.cam.pt.X < 0 {
		x0 -= 1
		xoff0 = -TileSize + xoff0
	}
	y0 := int(e.cam.pt.Y / TileSize)
	yoff0 := -math.Mod(e.cam.pt.Y, TileSize)
	if e.cam.pt.Y < 0 {
		y0 -= 1
		yoff0 = -TileSize + yoff0
	}
	pt := ui.Pt(xoff0, yoff0)
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

	return e.astro.Draw(d, e.cam)
}

func drawCell(d Drawer, l *world.Loc, x, y int, pt ui.Point) error {
	const minSh = 0.25
	const slope = (1 - minSh) / world.MaxElevation

	_, err := d.Draw(ui.Sprite{
		Name:   l.Terrain.Name,
		Bounds: ui.Rect(0, 0, TileSize, TileSize),
		Shade:  slope*float32(l.Elevation-l.Depth) + minSh,
	}, pt)

	return err
}

func (ex *ExploreScreen) Handle(stk *ScreenStack, ev ui.Event) error {
	switch k := ev.(type) {
	case ui.Key:
		const speed = 5 // px
		switch {
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Left:
			ex.astro.Move(ui.Pt(speed, 0))
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Right:
			ex.astro.Move(ui.Pt(-speed, 0))
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Down:
			ex.astro.Move(ui.Pt(0, speed))
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Up:
			ex.astro.Move(ui.Pt(0, -speed))
		}
	}
	return nil
}

func (t *ExploreScreen) Update(stk *ScreenStack) error {
	return nil
}
