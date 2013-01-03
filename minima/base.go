// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/item"
	"code.google.com/p/min-game/ui"
)

type Base struct {
	Box geom.Rectangle

	Storage Inventory
}

func NewBase(p geom.Point) Base {
	return Base{
		Box: geom.Rectangle{
			Min: p,
			Max: p.Add(geom.Pt(64, 64)),
		},
		Storage: Inventory{[]*item.Item{item.New(item.ETele)}, 0, false},
	}
}

func (b *Base) Draw(d ui.Drawer, cam ui.Camera) {
	cam.Draw(d, ui.Sprite{
		Name:   "Base",
		Bounds: geom.Rect(0, 0, b.Box.Dx(), b.Box.Dy()),
		Shade:  1.0,
	}, b.Box.Min)
}

// PutStorage adds i to the base's storage.
func (b *Base) PutStorage(i *item.Item) {
	b.Storage.Put(i)
}

type BaseScreen struct {
	astro   *Player
	base    *Base
	closing bool
}

const pad = 4

var origin = geom.Pt(32, 32)
var bounds = geom.Rectangle{
	Min: origin,
	Max: origin.Add(geom.Pt(ScreenDims.X, ScreenDims.Y/2)).Sub(origin.Mul(2)),
}
var packBounds = bounds.Add(geom.Pt(0, bounds.Dy()+3*pad+32))

func NewBaseScreen(astro *Player, base *Base) *BaseScreen {
	return &BaseScreen{astro, base, false}
}

func (s *BaseScreen) Transparent() bool {
	return true
}

func (s *BaseScreen) Draw(d ui.Drawer) {
	d.SetFont(DialogFont, 16)
	pt := s.astro.pack.Draw("Pack", d, pad, origin, true)
	s.base.Storage.Draw("Storage", d, pad, geom.Pt(origin.X, pt.Y+32+2*pad), false)
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
	case ui.Action:
		if s.astro.pack.Selected >= 0 && s.astro.pack.Get(s.astro.pack.Selected) != nil {
			i := s.astro.pack.Get(s.astro.pack.Selected)
			s.astro.pack.Set(s.astro.pack.Selected, nil)
			s.base.PutStorage(i)
		}
		if s.base.Storage.Selected >= 0 && s.base.Storage.Get(s.base.Storage.Selected) != nil {
			i := s.base.Storage.Get(s.base.Storage.Selected)
			if s.astro.PutPack(i) {
				s.base.Storage.Set(s.base.Storage.Selected, nil)
			}
		}
	case ui.Left:
		if s.astro.pack.Selected >= 0 {
			s.astro.pack.Selected--
			if s.astro.pack.Selected < 0 {
				s.astro.pack.Selected = s.astro.pack.Len() - 1
			}
		} else {
			s.base.Storage.Selected--
			if s.base.Storage.Selected < 0 {
				s.astro.pack.Selected = s.base.Storage.Len() - 1
			}
		}
	case ui.Right:
		if s.astro.pack.Selected >= 0 {
			s.astro.pack.Selected++
			if s.astro.pack.Selected == s.astro.pack.Len() {
				s.astro.pack.Selected = 0
			}
		} else {
			s.base.Storage.Selected++
			if s.base.Storage.Selected == s.base.Storage.Len() {
				s.base.Storage.Selected = 0
			}
		}
	case ui.Up, ui.Down:
		if s.astro.pack.Selected >= 0 {
			s.base.Storage.Selected = s.astro.pack.Selected
			if s.base.Storage.Selected >= s.base.Storage.Len() {
				s.base.Storage.Selected = s.base.Storage.Len() - 1
			}
			s.astro.pack.Selected = -1
		} else {
			s.astro.pack.Selected = s.base.Storage.Selected
			if s.astro.pack.Selected >= s.astro.pack.Len() {
				s.astro.pack.Selected = s.astro.pack.Len() - 1
			}
			s.base.Storage.Selected = -1
		}
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
