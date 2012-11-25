// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package sprite

import (
	"encoding/json"
	"math"
	"os"

	"code.google.com/p/min-game/geom"
)

type Sheet struct {
	Name                     string
	FrameSize                int
	Tempo                    int
	North, East, South, West int
}

func LoadSheet(s string) (Sheet, error) {
	var sh Sheet

	f, err := os.Open("resrc/" + s + ".sheet")
	if err != nil {
		return sh, err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	err = dec.Decode(&sh)
	return sh, err
}

func (sh *Sheet) Frame(Face, Frame int) geom.Rectangle {
	sz := float64(sh.FrameSize)
	x := float64(Frame) * sz
	y := float64(Face) * sz
	return geom.Rect(x, y, x+sz, y+sz)
}

type Anim struct {
	Face, Frame int
	Ticks       int
}

func (a *Anim) Move(sh *Sheet, vel geom.Point) {
	a.Ticks++
	if a.Ticks >= sh.Tempo {
		a.Frame++
		if a.Frame >= 2 {
			a.Frame = 0
		}
		a.Ticks = 0
	}

	dx, dy := vel.X, vel.Y
	vertBiased := math.Abs(dy) > math.Abs(dx)

	if dy > 0 && vertBiased {
		a.Face = sh.South
	}
	if dy < 0 && vertBiased {
		a.Face = sh.North
	}
	if dx > 0 && !vertBiased {
		a.Face = sh.East
	}
	if dx < 0 && !vertBiased {
		a.Face = sh.West
	}
}
