// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ai"
	"code.google.com/p/min-game/animal"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/ui"
	"code.google.com/p/min-game/world"
)

type Cow animal.Herbivore

var cowInfo animal.Info

func init() {
	var err error
	cowInfo, err = animal.LoadInfo("Cow")
	if err != nil {
		panic(err)
	}
}

func NewCow(p, v geom.Point) *Cow {
	return &Cow{
		Body: phys.Body{
			Box: geom.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
			Vel: v,
		},
	}
}

func (c *Cow) Draw(d Drawer, cam ui.Camera) error {
	return (*animal.Herbivore)(c).Draw(&cowInfo, d, cam)
}

func (c *Cow) Move(w *world.World) {
	(*animal.Herbivore)(c).Move(&cowInfo, w)
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
