// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/uitil"
)

type PauseScreen struct {
	astro    *Player
	closing  bool
	selected int
	inPack   bool
}

func NewPauseScreen(astro *Player) *PauseScreen {
	return &PauseScreen{astro, false, 0, false}
}

func (p *PauseScreen) Transparent() bool {
	return true
}

func (p *PauseScreen) Draw(d ui.Drawer) {
	d.SetFont("prstartk", 16)

	origin := geom.Pt(32, 32)
	pad := 4.0

	suit := make([]Item, 0, len(p.astro.suit))
	for _, a := range p.astro.suit {
		suit = append(suit, a)
	}
	pt := p.drawInventory(d, "Suit: ", suit, !p.inPack, pad, origin)
	pt = p.drawInventory(d, "Pack: ", p.astro.pack, p.inPack, pad, geom.Pt(origin.X, pt.Y+pad))

	if p.inPack && p.astro.pack[p.selected] == nil {
		return
	}
	if !p.inPack && p.astro.suit[p.selected] == nil {
		return
	}

	descBounds := geom.Rectangle{
		Min: geom.Pt(origin.X, pt.Y+pad*2),
		Max: geom.Pt(ScreenDims.X-origin.X, ScreenDims.Y-origin.Y),
	}

	d.SetColor(White)
	d.Draw(descBounds, geom.Pt(0, 0))

	d.SetColor(Black)
	desc := ""
	if p.inPack {
		desc = p.astro.pack[p.selected].Desc()
	} else {
		desc = p.astro.suit[p.selected].Desc()
	}
	uitil.WordWrap(d, desc, descBounds.Rpad(pad))
}

func (p *PauseScreen) Handle(stk *ui.ScreenStack, e ui.Event) error {
	if p.closing {
		return nil
	}

	key, ok := e.(ui.Key)
	if !ok || !key.Down {
		return nil
	}

	switch key.Button {
	case ui.Menu:
		p.closing = true
	case ui.Left:
		p.selected--
		if p.selected < 0 {
			if p.inPack {
				p.selected = len(p.astro.pack) - 1
			} else {
				p.selected = len(p.astro.suit) - 1
			}
		}
	case ui.Right:
		p.selected++
		if p.inPack && p.selected == len(p.astro.pack) {
			p.selected = 0
		}
		if !p.inPack && p.selected == len(p.astro.suit) {
			p.selected = 0
		}
	case ui.Up, ui.Down:
		p.inPack = !p.inPack
		if p.inPack && p.selected >= len(p.astro.pack) {
			p.selected = len(p.astro.pack) - 1
		}
		if !p.inPack && p.selected >= len(p.astro.suit) {
			p.selected = len(p.astro.suit) - 1
		}
	}

	return nil
}

func (p *PauseScreen) Update(stk *ui.ScreenStack) error {
	if p.closing {
		stk.Pop()
		return nil
	}

	return nil
}

func (p *PauseScreen) drawInventory(d ui.Drawer, label string, items []Item, hilight bool, pad float64, origin geom.Point) geom.Point {
	size := d.TextSize(label)

	width := 32.0*float64(len(items)) + pad*float64(len(items)+3) + size.X
	height := 32 + pad*2
	bounds := geom.Rectangle{
		Min: origin,
		Max: origin.Add(geom.Pt(width, height)),
	}

	d.SetColor(White)
	d.Draw(bounds, geom.Pt(0, 0))

	pt := origin.Add(geom.Pt(pad, pad))

	d.SetColor(Black)
	d.Draw(label, pt)
	pt.X += size.X + pad

	for i, a := range items {
		if hilight && i == p.selected {
			d.SetColor(Black)
			d.Draw(geom.Rectangle{
				Min: pt.Sub(geom.Pt(2, 2)),
				Max: pt.Add(geom.Pt(34, 34)),
			}, geom.Pt(0, 0))
		}

		if a != nil {
			d.Draw(ui.Sprite{
				Name:   a.Name(),
				Bounds: geom.Rect(0, 0, 32, 32),
				Shade:  1.0,
			}, pt)
		}

		pt.X += 32.0 + pad
	}

	return bounds.Max
}
