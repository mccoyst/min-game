// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"io"

	"code.google.com/p/min-game/ui"
)

type TitleScreen struct {
}

func NewTitleScreen() *TitleScreen {
	return &TitleScreen{}
}

func (t *TitleScreen) Draw(w io.Writer) error {
	// w.SetColor(255, 255, 255, 255) oops, need a draw command for this
	_, err := w.Write([]byte(
`rectfill 64 64 23 76
img resrc/Astronaut.png 0 4 7 8 1.0 5 6
img resrc/Cow.png 128 128 16 16 1.0
`))
	return err
}

func (t *TitleScreen) Handle(stk *ScreenStack, e ui.Event) error {
	return nil
}

func (t *TitleScreen) Update(stk *ScreenStack) error {
	return nil
}
