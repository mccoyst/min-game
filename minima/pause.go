// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"fmt"

	"github.com/mccoyst/min-game/geom"
	"github.com/mccoyst/min-game/item"
	"github.com/mccoyst/min-game/ui"
	"github.com/mccoyst/min-game/uitil"
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
	held := geom.Pt(pt.X+pad, origin.Y)
	packPt := geom.Pt(origin.X, pt.Y+pad)
	pt = p.astro.pack.Draw("Pack: ", d, pad, packPt, true)

	scrapPt := geom.Pt(pt.X+pad, packPt.Y)
	scrapText := fmt.Sprintf("Scrap: %d", p.astro.Scrap)
	scrapDims := d.TextSize(scrapText)
	scrapBounds := geom.Rectangle{
		Max: geom.Pt(scrapDims.X+2*pad, TileSize.Y+2*pad),
	}
	d.SetColor(Black)
	d.Draw(scrapBounds.Pad(pad), scrapPt)
	d.SetColor(White)
	d.Draw(scrapBounds, scrapPt)
	d.SetColor(Black)
	d.Draw(scrapText, scrapPt.Add(geom.Pt(pad, pad)))

	if p.astro.Held != nil {
		hinv := Inventory{[]*item.Item{p.astro.Held}, -1, true}
		hinv.Draw("Held: ", d, pad, held, true)
	}

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
	defer func() {
		// TODO(eaburns): Is there something more elegant than setting
		// every terrain type to it's best scale each time the player hands
		// something.
		for t, base := range baseScales {
			scales[t] = base
		}
		for _, i := range p.astro.suit.Items {
			if i == nil {
				continue
			}
			if b, ok := item.Bonus[i.Name]; ok && scales[b.Terrain] < b.Scale {
				scales[b.Terrain] = b.Scale
			}
		}
	}()

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
	case ui.Hands:
		a := p.astro
		if a.pack.Selected >= 0 {
			a.Held, a.pack.Items[a.pack.Selected] = a.pack.Items[a.pack.Selected], a.Held
		}
		if a.suit.Selected >= 0 {
			a.Held, a.suit.Items[a.suit.Selected] = a.suit.Items[a.suit.Selected], a.Held
		}
		return nil
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
