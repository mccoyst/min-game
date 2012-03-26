package world

import (
	"image"
	"image/color"
	"image/png"
	"os"
)

// SavePng saves the world to the given png file.
func (w *World) SavePng(file string, depth bool) error {
	out, err := os.Create(file)
	if err != nil {
		return err
	}
	defer out.Close()
	png.Encode(out, &worldImg{ w, depth })
	return nil
}

type worldImg struct {
	*World
	depth bool
}

// Bounds implements the Bounds() method of
// the image.Image interface.
func (w *worldImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.W, w.H)
}

// At implements the At() method of the
// image.Image interface.
func (w *worldImg) At(x, y int) color.Color {
	loc := w.World.At(x, y)
	ht := loc.Height
	if w.depth {
		ht -= loc.Depth
	}
	f := float64(ht+MaxHeight) / (2 * MaxHeight)
	if f > 1 {
		panic("Color factor is >1 in worldImg.At")
	}
	c := loc.Terrain.Color
	return color.RGBA{
		R: uint8(float64(c.R) * f),
		G: uint8(float64(c.G) * f),
		B: uint8(float64(c.B) * f),
		A: c.A,
	}
}

// ColorModel implements the ColorModel() method
// of the image.Image interface.
func (w *worldImg) ColorModel() color.Model {
	return color.RGBAModel
}
