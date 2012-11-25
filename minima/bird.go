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

type Gull animal.Herbivore

var gullInfo animal.Info

func init() {
	var err error
	gullInfo, err = animal.LoadInfo("Gull")
	if err != nil {
		panic(err)
	}
}

func NewGull(p, v geom.Point) *Gull {
	return &Gull{
		Body: phys.Body{
			Box: geom.Rect(p.X, p.Y, p.X+TileSize, p.Y+TileSize),
			Vel: v,
		},
	}
}

func (g *Gull) Draw(d Drawer, cam ui.Camera) error {
	return (*animal.Herbivore)(g).Draw(&gullInfo, d, cam)
}

func (g *Gull) Move(w *world.World) {
	(*animal.Herbivore)(g).Move(&gullInfo, w)
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
