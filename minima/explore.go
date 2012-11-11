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
	wo *world.World

	// Point is the pixel of the world drawn in the upper
	// left of the screen.  Each tile is TileSize by TileSize
	// pixels, so for example, to have the full 1,1 tile in
	// the upper left:
	// 	Point=ui.Pt(TileSize, TileSize)
	// To have the center of tile 0,0 in the upper right:
	//	Point=ui.Pt(TileSize/2, TileSize/2)
	ui.Point
}

func NewExploreScreen(wo *world.World) *ExploreScreen {
	e := &ExploreScreen{wo, ui.Pt(0, 0)}
	e.CenterOnTile(wo.X0, wo.Y0)
	return e
}

// CenterOnTile centers the display on a given tile.
func (e *ExploreScreen) CenterOnTile(x, y int) {
	e.X = TileSize*float64(x) + TileSize/2 - ScreenDims.X/2
	e.Y = TileSize*float64(y) + TileSize/2 - ScreenDims.Y/2
}

func (e *ExploreScreen) Draw(d Drawer) error {
	w, h := int(ScreenDims.X/TileSize), int(ScreenDims.Y/TileSize)
	x0 := int(e.X / TileSize)
	xoff0 := -math.Mod(e.X, TileSize)
	if e.X < 0 {
		x0 -= 1
		xoff0 = -TileSize + xoff0
	}
	y0 := int(e.Y / TileSize)
	yoff0 := -math.Mod(e.Y, TileSize)
	if e.Y < 0 {
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

	return nil
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
			ex.Point = ex.Add(ui.Pt(speed, 0))
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Right:
			ex.Point = ex.Sub(ui.Pt(speed, 0))
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Down:
			ex.Point = ex.Add(ui.Pt(0, speed))
		case k.Down && ui.DefaultKeymap[k.Code] == ui.Up:
			ex.Point = ex.Sub(ui.Pt(0, speed))
		}
	}
	return nil
}

func (t *ExploreScreen) Update(stk *ScreenStack) error {
	return nil
}
