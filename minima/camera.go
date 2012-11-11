// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"code.google.com/p/min-game/ui"
)

type Camera struct {
	pt ui.Point
}

func (c *Camera) Move(v ui.Point) {
	c.pt = c.pt.Add(v)
}

func (c *Camera) Center(v ui.Point) {
	c.pt.X = v.X - ScreenDims.X/2.0
	c.pt.Y = v.Y - ScreenDims.Y/2.0
}

func (c *Camera) Pos() ui.Point {
	return c.pt
}

func (c *Camera) Draw(d Drawer, x interface{}, p ui.Point) (ui.Point, error) {
	return d.Draw(x, p.Sub(c.pt))
}
