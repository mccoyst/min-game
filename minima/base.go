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
		Storage: Inventory{[]*item.Item{item.New(item.ETele)}, 0},
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

	inPack   bool
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
	return &BaseScreen{astro, base, false, false, 0}
}

func (s *BaseScreen) Transparent() bool {
	return true
}

func (s *BaseScreen) Draw(d ui.Drawer) {
	d.SetFont(DialogFont, 16)
	pt := DrawInventory(BaseInv{s, "Pack"}, d, pad, origin, true)
	DrawInventory(BaseInv{s, "Storage"}, d, pad, geom.Pt(origin.X, pt.Y+32+2*pad), false)
}

type BaseInv struct {
	s     *BaseScreen
	label string
}

func (b BaseInv) Label() string {
	return b.label
}

func (b BaseInv) Len() int {
	if b.label == "Pack" {
		return len(b.s.astro.pack)
	}
	return b.s.base.Storage.Len()
}

func (b BaseInv) Selected(n int) bool {
	if b.label == "Pack" && b.s.inPack || b.label == "Storage" && !b.s.inPack {
		return b.s.selected == n
	}
	return false
}

func (b BaseInv) Get(n int) *item.Item {
	if b.label == "Pack" {
		return b.s.astro.pack[n]
	}
	return b.s.base.Storage.Items[n]
}

func (b BaseInv) Set(n int, i *item.Item) {
	if b.label == "Pack" {
		b.s.astro.pack[n] = i
	} else {
		b.s.base.Storage.Items[n] = i
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
	case ui.Action:
		if s.inPack && s.astro.pack[s.selected] != nil {
			i := s.astro.pack[s.selected]
			s.astro.pack[s.selected] = nil
			s.base.PutStorage(i)
		}
		if !s.inPack && s.base.Storage.Items[s.selected] != nil {
			i := s.base.Storage.Items[s.selected]
			if s.astro.PutPack(i) {
				s.base.Storage.Items[s.selected] = nil
			}
		}
	case ui.Left:
		s.selected--
		if s.selected < 0 {
			if s.inPack {
				s.selected = len(s.astro.pack) - 1
			} else {
				s.selected = s.base.Storage.Len() - 1
			}
		}
	case ui.Right:
		s.selected++
		if s.inPack && s.selected == len(s.astro.pack) {
			s.selected = 0
		}
		if !s.inPack && s.selected == s.base.Storage.Len() {
			s.selected = 0
		}
	case ui.Up, ui.Down:
		s.inPack = !s.inPack
		if s.inPack && s.selected >= len(s.astro.pack) {
			s.selected = len(s.astro.pack) - 1
		}
		if !s.inPack && s.selected >= s.base.Storage.Len() {
			s.selected = s.base.Storage.Len() - 1
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
