// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/geom"
)

type Camera struct {
	pt geom.Point
}

func (c *Camera) Move(v geom.Point) {
	c.pt = c.pt.Add(v)
}

func (c *Camera) Center(v geom.Point) {
	c.pt.X = v.X - ScreenDims.X/2.0
	c.pt.Y = v.Y - ScreenDims.Y/2.0
}

func (c *Camera) Pos() geom.Point {
	return c.pt
}

func (c *Camera) Draw(d Drawer, x interface{}, p geom.Point) (geom.Point, error) {
	return d.Draw(x, p.Sub(c.pt))
}
