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

func (c *Camera) Rect() geom.Rectangle {
	return geom.Rect(c.Pt.X, c.Pt.Y, c.Pt.X+c.Dims.X, c.Pt.Y+c.Dims.Y)
}

func (c *Camera) Draw(d Drawer, s Sprite, p geom.Point) geom.Point {
	sz := s.Bounds.Size()
	p = p.Sub(c.Pt)
	rect := geom.Rect(p.X, p.Y, p.X+sz.X, p.Y+sz.Y)
	_, rect = c.Torus.AlignRects(c.Rect(), rect)
	return d.Draw(s, rect.Min)
}
