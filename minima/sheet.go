// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"encoding/json"
	"os"

	"code.google.com/p/min-game/ui"
)

type SpriteSheet struct {
	Name                     string
	FrameSize                int
	Tempo                    int
	North, East, South, West int
}

func LoadSpriteSheet(s string) (SpriteSheet, error) {
	var sh SpriteSheet

	f, err := os.Open("resrc/" + s + ".sheet")
	if err != nil {
		return sh, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&sh)
	return sh, err
}

func (sh *SpriteSheet) Frame(face, frame int) ui.Rectangle {
	sz := float64(sh.FrameSize)
	x := float64(frame) * sz
	y := float64(face) * sz
	return ui.Rect(x, y, x+sz, y+sz)
}
