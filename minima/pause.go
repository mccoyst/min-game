// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

type PauseScreen struct {
	closing bool
	size    float64
}

func NewPauseScreen() *PauseScreen {
	return &PauseScreen{size: 64}
}

func (p *PauseScreen) Transparent() bool {
	return true
}

func (p *PauseScreen) Draw(d ui.Drawer) error {
	d.SetColor(White)
	if err := d.SetFont("prstartk", p.size); err != nil {
		return err
	}

	txt := "PAUSED"
	sz := d.TextSize(txt)

	_, err := d.Draw(txt, geom.Pt(ScreenDims.X/2-sz.X/2, ScreenDims.Y/2-sz.Y))
	return err
}

func (p *PauseScreen) Handle(stk *ui.ScreenStack, e ui.Event) error {
	if p.closing {
		return nil
	}

	if key, ok := e.(ui.Key); ok && key.Down {
		p.closing = true
	}
	return nil
}

func (p *PauseScreen) Update(stk *ui.ScreenStack) error {
	if !p.closing {
		return nil
	}

	p.size -= 5
	if p.size < 12 {
		stk.Pop()
	}
	return nil
}
