// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package animal

import (
	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/geom"
)

type Gull struct {
	Herbivore
}

var gullInfo Info

func init() {
	var err error
	gullInfo, err = LoadInfo("Gull")
	if err != nil {
		panic(err)
	}
}

func NewGull(p, v geom.Point) *Gull {
	return &Gull{gullInfo.SpawnHerbivore(p, v)}
}

type Gulls []*Gull

func (gs Gulls) Len() int {
	return len(gs)
}

func (gs Gulls) Boid(n int) ai.Boid {
	return ai.Boid{&gs[n].Body}
}

func (Gulls) Info() ai.BoidInfo {
	return gullInfo.BoidInfo
}
