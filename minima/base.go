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
}

func NewBaseScreen(astro *Player, base *Base) *BaseScreen {
	return &BaseScreen{astro, base, false}
}

func (s *BaseScreen) Transparent() bool {
	return true
}

func (s *BaseScreen) Draw(d ui.Drawer) {
	origin := geom.Pt(32, 32)
	d.SetColor(White)
	d.Draw(geom.Rectangle{
		Min: origin,
		Max: origin.Add(ScreenDims).Sub(origin.Mul(2)),
	}, geom.Pt(0, 0))

	d.SetColor(Black)
	d.SetFont("prstartk", 16)
	d.Draw("Something will go here.", origin.Add(geom.Pt(16, 16)))
}

func (s *BaseScreen) Handle(stk *ui.ScreenStack, e ui.Event) error {
	if s.closing {
		return nil
	}

	if key, ok := e.(ui.Key); ok && key.Down {
		s.closing = true
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
