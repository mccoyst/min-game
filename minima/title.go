// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"io"

	"code.google.com/p/min-game/ui"
)

type TitleScreen struct {
	frame int
}

func NewTitleScreen() *TitleScreen {
	return &TitleScreen{}
}

func (t *TitleScreen) Draw(w io.Writer) error {
	_, err := w.Write([]byte(
`color 255 255 255
rectfill 64 64 23 76
img Astronaut 64 64 16 16 1.0 16 0
img Cow 128 128 16 16 1.0
`))
	return err
}

func (t *TitleScreen) Handle(stk *ScreenStack, e ui.Event) error {
	return nil
}

func (t *TitleScreen) Update(stk *ScreenStack) error {
	return nil
}
