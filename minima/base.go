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

func (s *BaseScreen) Draw(d ui.Drawer) error {
	origin := geom.Pt(32, 32)
	d.SetColor(Black)
	_, err := d.Draw(geom.Rectangle{
		Min: origin,
		Max: origin.Add(ScreenDims).Sub(origin.Mul(2)),
	}, geom.Pt(0, 0))
	if err != nil {
		return err
	}

	d.SetColor(White)
	if err := d.SetFont("prstartk", 16); err != nil {
		return err
	}

	_, err = d.Draw("Something will go here.", origin.Add(geom.Pt(16, 16)))
	return err
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
