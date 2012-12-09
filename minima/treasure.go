// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/uitil"
)

type Treasure struct {
	Item Item
	Box  geom.Rectangle
}

func (t *Treasure) Draw(d ui.Drawer, cam ui.Camera) {
	cam.Draw(d, ui.Sprite{
		Name:   "Present",
		Bounds: geom.Rect(0, 0, t.Box.Dx(), t.Box.Dy()),
		Shade:  1.0, //TODO: should shade with altitude
	}, t.Box.Min)
}

type TreasureGet struct {
	name    string
	closing bool
}

func NewTreasureGet(name string) *TreasureGet {
	return &TreasureGet{name, false}
}

func (g *TreasureGet) Transparent() bool {
	return true
}

func (g *TreasureGet) Draw(d ui.Drawer) {
	origin := geom.Pt(32, 32)
	dims := geom.Pt(ScreenDims.X, ScreenDims.Y/2)
	d.SetFont("prstartk", 16)
	d.SetColor(White)
	box := geom.Rectangle{
		Min: origin,
		Max: origin.Add(dims).Sub(origin.Mul(2)),
	}
	d.Draw(box, geom.Pt(0, 0))

	d.SetColor(Black)
	uitil.WordWrap(d, "Bravo! You got the "+g.name+"!", box.Rpad(4))
}

func (g *TreasureGet) Handle(stk *ui.ScreenStack, e ui.Event) error {
	if g.closing {
		return nil
	}

	key, ok := e.(ui.Key)
	if !ok || !key.Down {
		return nil
	}

	g.closing = true
	return nil
}

func (g *TreasureGet) Update(stk *ui.ScreenStack) error {
	if g.closing {
		stk.Pop()
		return nil
	}

	return nil
}
