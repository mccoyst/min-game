// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"strings"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

type PauseScreen struct {
	astro    *Player
	closing  bool
	selected int
}

func NewPauseScreen(astro *Player) *PauseScreen {
	return &PauseScreen{astro, false, 0}
}

func (p *PauseScreen) Transparent() bool {
	return true
}

func (p *PauseScreen) Draw(d ui.Drawer) error {
	if err := d.SetFont("prstartk", 16); err != nil {
		return err
	}

	origin := geom.Pt(32, 32)
	pad := 4.0
	suitLabel := "Suit: "
	suitSz := d.TextSize(suitLabel)

	width := 32.0*float64(p.astro.maxAugs) + pad*float64(p.astro.maxAugs+3) + suitSz.X
	height := 32 + pad*2
	suitBounds := geom.Rectangle{
		Min: origin,
		Max: origin.Add(geom.Pt(width, height)),
	}

	d.SetColor(Black)
	_, err := d.Draw(suitBounds, geom.Pt(0, 0))
	if err != nil {
		return err
	}

	pt := origin.Add(geom.Pt(pad, pad))

	d.SetColor(White)
	_, err = d.Draw(suitLabel, pt)
	if err != nil {
		return err
	}
	pt.X += suitSz.X + pad

	for i, a := range p.astro.suit {
		if i == p.selected {
			d.SetColor(Lemon)
			_, err = d.Draw(geom.Rectangle{
				Min: pt.Sub(geom.Pt(2, 2)),
				Max: pt.Add(geom.Pt(34, 34)),
			}, geom.Pt(0, 0))
			if err != nil {
				return err
			}
		}

		if a == nil {
			continue
		}

		_, err = d.Draw(ui.Sprite{
			Name:   a.Name(),
			Bounds: geom.Rect(0, 0, 32, 32),
			Shade:  1.0,
		}, pt)
		if err != nil {
			return err
		}
		pt.X += 32.0 + pad
	}

	if p.astro.suit[p.selected] == nil {
		return nil
	}

	descBounds := geom.Rectangle{
		Min: geom.Pt(origin.X, suitBounds.Max.Y+pad*2),
		Max: geom.Pt(ScreenDims.X-origin.X, ScreenDims.Y-origin.Y),
	}

	d.SetColor(Black)
	_, err = d.Draw(descBounds, geom.Pt(0, 0))
	if err != nil {
		return err
	}

	desc := p.astro.suit[p.selected].Desc()
	words := strings.Fields(desc)

	d.SetColor(White)

	left := descBounds.Min.X + pad
	wp := geom.Pt(left, descBounds.Min.Y+pad)
	for _, word := range words {
		word += " "
		wsz := d.TextSize(word)
		if wp.X+wsz.X > descBounds.Dx() {
			wp.Y += wsz.Y + pad
			wp.X = left
		}

		_, err = d.Draw(word, wp)
		if err != nil {
			return err
		}

		wp.X += wsz.X
	}

	return err
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
			p.selected = p.astro.maxAugs - 1
		}
	case ui.Right:
		p.selected++
		if p.selected == p.astro.maxAugs {
			p.selected = 0
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
