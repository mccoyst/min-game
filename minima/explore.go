// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

// TileSize is the width and height of a tile in pixels.
const TileSize = 32

type ExploreScreen struct {
	wo world.World
}

func NewExploreScreen(wo world.World) *ExploreScreen {
	return &ExploreScreen{wo}
}

func (e *ExploreScreen) Draw(d Drawer) error {
	w, h := int(ScreenDims.X/TileSize), int(ScreenDims.Y/TileSize)

	for y := -1; y <= h+1; y++ {
		for x := -1; x <= w; x++ {
			l := e.wo.At(x, y)
			err := e.drawCell(d, l, x*TileSize, y*TileSize)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *ExploreScreen) drawCell(d Drawer, l *world.Loc, x, y int) error {
	const minSh = 0.25
	const slope = (1 - minSh) / world.MaxElevation

	_, err := d.Draw(ui.Sprite{
		Name:   l.Terrain.Name,
		Bounds: ui.Rect(0, 0, TileSize, TileSize),
		Shade:  slope*float32(l.Elevation-l.Depth) + minSh,
	}, ui.Pt(float64(x), float64(y)))

	return err
}

func (ex *ExploreScreen) Handle(stk *ScreenStack, e ui.Event) error {
	return nil
}

func (t *ExploreScreen) Update(stk *ScreenStack) error {
	return nil
}
