// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

type Camera struct {
	cam Point
}

func (c *Camera) Move(v Point) {
	c.cam.Add(v)
}

func (c *Camera) Center(v Point) {
	c.cam.X = ScreenDims.X.Div(F(2)).Sub(v.X)
	c.cam.Y = ScreenDims.Y.Div(F(2)).Sub(v.Y)
}

func (c *Camera) Pos() Point {
	return c.cam
}

func (c *Camera) Draw(p Point, ui Ui, img Img, shade float32) {
	img.Draw(ui, p.Add(c.cam), shade)
}
