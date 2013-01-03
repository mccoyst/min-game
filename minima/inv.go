// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/item"
	"code.google.com/p/min-game/ui"
)

type Inventory struct {
	Items    []*item.Item
	Selected int
}

func (i *Inventory) Len() int {
	return len(i.Items)
}

func (i *Inventory) Get(n int) *item.Item {
	return i.Items[n]
}

func (i *Inventory) Draw(label string, d ui.Drawer, pad float64, origin geom.Point, fit bool) geom.Point {
	size := d.TextSize(label)

	width, height := 0.0, 0.0
	if fit {
		width = TileSize*float64(i.Len()) + pad*float64(i.Len()+3) + size.X
		height = TileSize + pad*2
	} else {
		width = ScreenDims.X - TileSize - origin.X
		height = ScreenDims.Y - TileSize - origin.Y
	}

	bounds := geom.Rectangle{
		Min: origin,
		Max: origin.Add(geom.Pt(width, height)),
	}

	d.SetColor(Black)
	d.Draw(bounds.Pad(pad), geom.Pt(0, 0))

	d.SetColor(White)
	d.Draw(bounds, geom.Pt(0, 0))

	pt := origin.Add(geom.Pt(pad, pad))

	d.SetColor(Black)
	d.Draw(label, pt)
	if fit {
		pt.X += size.X + pad
	} else {
		pt.Y += size.Y + pad
	}

	for j := 0; j < i.Len(); j++ {
		a := i.Get(j)
		if i.Selected == j {
			d.SetColor(Black)
			d.Draw(geom.Rectangle{
				Min: pt.Sub(geom.Pt(2, 2)),
				Max: pt.Add(geom.Pt(34, 34)),
			}, geom.Pt(0, 0))
		}

		if a != nil {
			d.Draw(ui.Sprite{
				Name:   a.Name,
				Bounds: geom.Rect(0, 0, 32, 32),
				Shade:  1.0,
			}, pt)
		}

		pt.X += 32.0 + pad
		if pt.X >= bounds.Dx() {
			pt.X = bounds.Min.X
			pt.Y += TileSize + pad
		}
	}

	return bounds.Max
}
