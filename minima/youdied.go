// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
)

type GameOverScreen struct {
}

func NewGameOverScreen() *GameOverScreen {
	return &GameOverScreen{}
}

func (t *GameOverScreen) Transparent() bool {
	return false
}

func (t *GameOverScreen) Draw(d ui.Drawer) error {
	d.SetColor(Red)
	_, err := d.Draw(geom.Rect(0, 0, ScreenDims.X, ScreenDims.Y), geom.Pt(0, 0))
	if err != nil {
		return err
	}

	d.SetColor(Black)
	if err := d.SetFont("bit_outline", 96); err != nil {
		return err
	}
	text := "You Died"
	textSz := d.TextSize(text)
	textPos := geom.Pt(ScreenDims.X/2-textSz.X/2,
		ScreenDims.Y/2-textSz.Y)
	wh, err := d.Draw(text, textPos)
	if err != nil {
		return err
	}

	if err := d.SetFont("prstartk", 12); err != nil {
		return err
	}
	flavor := "…and there was nothing…"
	flavorSz := d.TextSize(flavor)
	flavorPos := geom.Pt(ScreenDims.X/2-flavorSz.X/2, textPos.Y+wh.Y+flavorSz.Y)
	_, err = d.Draw(flavor, flavorPos)

	return err
}

func (t *GameOverScreen) Handle(stk *ui.ScreenStack, e ui.Event) error {
	return nil
}

func (t *GameOverScreen) Update(stk *ui.ScreenStack) error {
	return nil
}
