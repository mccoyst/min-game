// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ui

import (
	"code.google.com/p/min-game/geom"
)

type Camera struct {
	Pt   geom.Point
	Dims geom.Point
}

func (c *Camera) Move(v geom.Point) {
	c.Pt = c.Pt.Add(v)
}

func (c *Camera) Center(v geom.Point) {
	c.Pt.X = v.X - c.Dims.X/2.0
	c.Pt.Y = v.Y - c.Dims.Y/2.0
}

func (c *Camera) Draw(d Drawer, x interface{}, p geom.Point) (geom.Point, error) {
	return d.Draw(x, p.Sub(c.Pt))
}
