// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
)

type TitleScreen struct {
	frame int
}

func NewTitleScreen() *TitleScreen {
	return &TitleScreen{}
}

func (t *TitleScreen) Draw(d Drawer) error {
	d.SetColor(255, 255, 255, 255)
	d.Draw(ui.Rect(64, 64, 64+23, 64+76), ui.Pt(0, 0))

	_, err := d.Draw(ui.Sprite{
		Name: "Astronaut",
		Bounds: ui.Rect(16, 0, 16+16, 0+16),
		Shade: 1.0,
	}, ui.Pt(64, 64))
	if err != nil {
		return err
	}

	_, err = d.Draw(ui.Sprite{
		Name: "Cow",
		Bounds: ui.Rect(0, 0, 16, 16),
		Shade: 1.0,
	}, ui.Pt(128, 128))

	d.SetColor(255, 255, 255, 255)
	_, err = d.Draw("Hello, this is minima", ui.Pt(0, 0))
	return err
}

func (t *TitleScreen) Handle(stk *ScreenStack, e ui.Event) error {
	return nil
}

func (t *TitleScreen) Update(stk *ScreenStack) error {
	return nil
}
