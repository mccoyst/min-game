// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/item"
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
	d.SetFont(DialogFont, 16)

	origin := geom.Pt(32, 32)
	pad := 4.0

	pt := p.drawInventory(d, "Suit: ", pad, origin)
	pt = p.drawInventory(d, "Pack: ", pad, geom.Pt(origin.X, pt.Y+pad))

	if p.inPack && p.astro.pack[p.selected] == nil {
		return
	}
	if !p.inPack && p.astro.suit.Items[p.selected] == nil {
		return
	}

	descBounds := geom.Rectangle{
		Min: geom.Pt(origin.X, pt.Y+pad),
		Max: geom.Pt(ScreenDims.X-origin.X, ScreenDims.Y-origin.Y),
	}

	d.SetColor(Black)
	d.Draw(descBounds.Pad(pad), geom.Pt(0, 0))
	d.SetColor(White)
	d.Draw(descBounds, geom.Pt(0, 0))

	d.SetColor(Black)
	desc := ""
	if p.inPack {
		desc = p.astro.pack[p.selected].Desc()
	} else {
		desc = p.astro.suit.Items[p.selected].Desc()
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
	case ui.Action:
		if p.inPack && p.astro.pack[p.selected] != nil {
			if p.astro.PutSuit(p.astro.pack[p.selected]) {
				p.astro.pack[p.selected] = nil
			}
		}
		if !p.inPack && p.astro.suit.Items[p.selected] != nil {
			if p.astro.PutPack(p.astro.suit.Items[p.selected]) {
				p.astro.suit.Items[p.selected] = nil
			}
		}
	case ui.Left:
		p.selected--
		if p.selected < 0 {
			if p.inPack {
				p.selected = len(p.astro.pack) - 1
			} else {
				p.selected = p.astro.suit.Len() - 1
			}
		}
	case ui.Right:
		p.selected++
		if p.inPack && p.selected == len(p.astro.pack) {
			p.selected = 0
		}
		if !p.inPack && p.selected == p.astro.suit.Len() {
			p.selected = 0
		}
	case ui.Up, ui.Down:
		p.inPack = !p.inPack
		if p.inPack && p.selected >= len(p.astro.pack) {
			p.selected = len(p.astro.pack) - 1
		}
		if !p.inPack && p.selected >= p.astro.suit.Len() {
			p.selected = p.astro.suit.Len() - 1
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

func (p *PauseScreen) drawInventory(d ui.Drawer, label string, pad float64, origin geom.Point) geom.Point {
	return DrawInventory(PauseInv{p, label}, d, pad, origin, true)
}

type PauseInv struct {
	p     *PauseScreen
	label string
}

func (p PauseInv) Label() string {
	return p.label
}

func (p PauseInv) Len() int {
	if p.label == "Pack: " {
		return len(p.p.astro.pack)
	}
	return p.p.astro.suit.Len()
}

func (p PauseInv) Selected(i int) bool {
	if p.label == "Pack: " && p.p.inPack || p.label == "Suit: " && !p.p.inPack {
		return i == p.p.selected
	}
	return false
}

func (p PauseInv) Get(i int) *item.Item {
	if p.label == "Pack: " {
		return p.p.astro.pack[i]
	}
	return p.p.astro.suit.Get(i)
}

func (p PauseInv) Set(i int, n *item.Item) {
	if p.label == "Pack: " {
		p.p.astro.pack[i] = n
	} else {
		p.p.astro.suit.Items[i] = n
	}
}
