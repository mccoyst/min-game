// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
)

type Phys struct {
	Box ui.Rectangle
	Vel ui.Point
}

func (p *Phys) Move() {
	p.Box = p.Box.Add(p.Vel)
}

func (p *Phys) Overlaps(b *Phys) bool {
	return p.Box.Overlaps(b.Box)
}
