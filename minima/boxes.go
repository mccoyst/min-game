// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/uitil"
)

func NewNormalMessage(msg string) *uitil.MessageBox {
	origin := geom.Pt(32, 32)
	dims := geom.Pt(ScreenDims.X, ScreenDims.Y/2)
	box := geom.Rectangle{
		Min: origin,
		Max: origin.Add(dims).Sub(origin.Mul(2)),
	}

	return &uitil.MessageBox{
		Text:   msg,
		Font:   "prstartk",
		Fontsz: 16,
		Fg:     Black,
		Bg:     White,
		Box:    box,
		Pad:    4.0,
	}
}
