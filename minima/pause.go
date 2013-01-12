// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/uitil"
)

type PauseScreen struct {
	astro   *Player
	closing bool
}

func NewPauseScreen(astro *Player) *PauseScreen {
	return &PauseScreen{astro, false}
}

func (p *PauseScreen) Transparent() bool {
	return true
}

func (p *PauseScreen) Draw(d ui.Drawer) {
	d.SetFont(DialogFont, 16)

	origin := geom.Pt(32, 32)
	pad := 4.0

	pt := p.astro.suit.Draw("Suit: ", d, pad, origin, true)
	pt = p.astro.pack.Draw("Pack: ", d, pad, geom.Pt(origin.X, pt.Y+pad), true)

	if p.astro.pack.Selected >= 0 && p.astro.pack.Get(p.astro.pack.Selected) == nil {
		return
	}
	if p.astro.suit.Selected >= 0 && p.astro.suit.Get(p.astro.suit.Selected) == nil {
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
	if p.astro.pack.Selected >= 0 {
		desc = p.astro.pack.Get(p.astro.pack.Selected).Desc()
	} else {
		desc = p.astro.suit.Get(p.astro.suit.Selected).Desc()
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
	}

	HandleInvPair(&p.astro.pack, &p.astro.suit, key.Button)
	return nil
}

func (p *PauseScreen) Update(stk *ui.ScreenStack) error {
	if p.closing {
		stk.Pop()
		return nil
	}

	return nil
}
