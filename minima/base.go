// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

type Base struct {
	Box geom.Rectangle
}

func NewBase(p geom.Point) Base {
	return Base{
		Box: geom.Rectangle{
			Min: p,
			Max: p.Add(geom.Pt(64, 64)),
		},
	}
}

func (b *Base) Draw(d ui.Drawer, cam ui.Camera) error {
	_, err := cam.Draw(d, ui.Sprite{
		Name:   "Base",
		Bounds: geom.Rect(0, 0, b.Box.Dx(), b.Box.Dy()),
		Shade:  1.0,
	}, b.Box.Min)

	return err
}
