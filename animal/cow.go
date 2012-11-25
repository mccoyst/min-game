// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
)

type Cow struct {
	Herbivore
}

var cowInfo Info

func init() {
	var err error
	cowInfo, err = LoadInfo("Cow")
	if err != nil {
		panic(err)
	}
}

func NewCow(p, v geom.Point) *Cow {
	return &Cow{
		Herbivore{
			Body: phys.Body{
				Box: geom.Rect(p.X, p.Y, p.X+32, p.Y+32),
				Vel: v,
			},
			info: &cowInfo,
		},
	}
}

type Cows []*Cow

func (cs Cows) Len() int {
	return len(cs)
}

func (cs Cows) Boid(n int) ai.Boid {
	return ai.Boid{&cs[n].Body}
}

func (Cows) Info() ai.BoidInfo {
	return cowInfo.BoidInfo
}
