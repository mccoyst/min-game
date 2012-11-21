// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"encoding/json"
	"math"
	"os"

	"code.google.com/p/min-game/geom"
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

func (sh *SpriteSheet) Frame(face, frame int) geom.Rectangle {
	sz := float64(sh.FrameSize)
	x := float64(frame) * sz
	y := float64(face) * sz
	return geom.Rect(x, y, x+sz, y+sz)
}

type Anim struct {
	face, frame int
	ticks       int
}

func (a *Anim) Move(sh *SpriteSheet, vel geom.Point) {
	a.ticks++
	if a.ticks >= sh.Tempo {
		a.frame++
		if a.frame >= 2 {
			a.frame = 0
		}
		a.ticks = 0
	}

	dx, dy := vel.X, vel.Y
	vertBiased := math.Abs(dy) > math.Abs(dx)

	if dy > 0 && vertBiased {
		a.face = sh.South
	}
	if dy < 0 && vertBiased {
		a.face = sh.North
	}
	if dx > 0 && !vertBiased {
		a.face = sh.East
	}
	if dx < 0 && !vertBiased {
		a.face = sh.West
	}
}
