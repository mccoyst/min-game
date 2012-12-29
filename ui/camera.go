// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ui

import (
	"code.google.com/p/min-game/geom"
)

type Camera struct {
	Torus geom.Torus
	Pt    geom.Point
	Dims  geom.Point
}

func (c *Camera) Move(v geom.Point) {
	c.Pt = c.Pt.Add(v)
}

func (c *Camera) Center(v geom.Point) {
	c.Pt.X = v.X - c.Dims.X/2.0
	c.Pt.Y = v.Y - c.Dims.Y/2.0
}

func (c *Camera) Draw(d Drawer, s Sprite, p geom.Point) geom.Point {
	screen := geom.Rect(0, 0, c.Dims.X, c.Dims.Y)
	p = p.Sub(c.Pt)
	box := geom.Rect(p.X, p.Y, p.X+s.Bounds.Dx(), p.Y+s.Bounds.Dy())
	_, box = c.Torus.AlignRects(screen, box)
	return d.Draw(s, box.Min)
}
