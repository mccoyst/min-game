// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/item"
	"code.google.com/p/min-game/ui"
)

type Base struct {
	Box geom.Rectangle

	Storage []*item.Item
}

func NewBase(p geom.Point) Base {
	return Base{
		Box: geom.Rectangle{
			Min: p,
			Max: p.Add(geom.Pt(64, 64)),
		},
		Storage: []*item.Item{item.New(item.ETele)},
	}
}

func (b *Base) Draw(d ui.Drawer, cam ui.Camera) {
	cam.Draw(d, ui.Sprite{
		Name:   "Base",
		Bounds: geom.Rect(0, 0, b.Box.Dx(), b.Box.Dy()),
		Shade:  1.0,
	}, b.Box.Min)
}

type BaseScreen struct {
	astro   *Player
	base    *Base
	closing bool

	rowLen   int
	selected int
}

const pad = 4

var origin = geom.Pt(32, 32)
var bounds = geom.Rectangle{
	Min: origin,
	Max: origin.Add(geom.Pt(ScreenDims.X, ScreenDims.Y/2)).Sub(origin.Mul(2)),
}
var packBounds = bounds.Add(geom.Pt(0, bounds.Dy()+3*pad+32))

func NewBaseScreen(astro *Player, base *Base) *BaseScreen {
	r := int((bounds.Dx() - pad) / (TileSize + pad))
	return &BaseScreen{astro, base, false, r, 0}
}

func (s *BaseScreen) Transparent() bool {
	return true
}

func (s *BaseScreen) Draw(d ui.Drawer) {
	d.SetColor(Black)
	d.Draw(bounds.Pad(pad), geom.Pt(0, 0))

	d.SetColor(White)
	d.Draw(bounds, geom.Pt(0, 0))

	d.SetColor(Black)
	d.SetFont("prstartk", 16)
	pt := d.Draw("Storage", bounds.Min.Add(geom.Pt(pad, pad)))

	pt.X = bounds.Min.X + pad
	pt.Y = bounds.Min.Add(geom.Pt(pad, pad)).Y + pt.Y + pad
	for i, a := range s.base.Storage {
		if i == s.selected {
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

		pt.X += TileSize + pad
		if i >= s.rowLen {
			pt.Y += TileSize + pad
			pt.X = origin.X + pad
		}
	}

	d.SetColor(Black)
	d.Draw(packBounds.Pad(pad), geom.Pt(0, 0))

	d.SetColor(White)
	d.Draw(packBounds, geom.Pt(0, 0))

	d.SetColor(Black)
	d.SetFont("prstartk", 16)
	pt = d.Draw("Pack", packBounds.Min.Add(geom.Pt(pad, pad)))

	pt.X = packBounds.Min.X + pad
	pt.Y = packBounds.Min.Add(geom.Pt(pad, pad)).Y + pt.Y + pad
	for i, a := range s.astro.pack {
		if i == s.selected {
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

		pt.X += TileSize + pad
		if i >= s.rowLen {
			pt.Y += TileSize + pad
			pt.X = origin.X + pad
		}
	}
}

func (s *BaseScreen) Handle(stk *ui.ScreenStack, e ui.Event) error {
	if s.closing {
		return nil
	}

	key, ok := e.(ui.Key)
	if !ok || !key.Down {
		return nil
	}

	switch key.Button {
	case ui.Menu:
		s.closing = true
	case ui.Left:
	}
	return nil
}

func (s *BaseScreen) Update(stk *ui.ScreenStack) error {
	s.astro.RefillO2()

	if s.closing {
		stk.Pop()
	}
	return nil
}
