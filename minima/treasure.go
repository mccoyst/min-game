// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/uitil"
)

type Treasure struct {
	Item Item
	Box  geom.Rectangle
}

func (t *Treasure) Draw(d ui.Drawer, cam ui.Camera) {
	cam.Draw(d, ui.Sprite{
		Name:   "Present",
		Bounds: geom.Rect(0, 0, t.Box.Dx(), t.Box.Dy()),
		Shade:  1.0, //TODO: should shade with altitude
	}, t.Box.Min)
}

func NewTreasureGet(name string) *uitil.MessageBox {
	return NewNormalMessage("Bravo! You got the " + name + "!")
}
