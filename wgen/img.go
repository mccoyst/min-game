package main

import (
	"minima/world"
	"image/color"
	"image"
	"image/png"
	"os"
)

const (
	BlkWidth = 1
	BlkHeight = 1
)

type worldImg world.World

// Bounds implements the Bounds() method of
// the image.Image interface.
func (w *worldImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.W * BlkWidth, w.H * BlkHeight)
}

// At implements the At() method of the
// image.Image interface.
func (w *worldImg) At(x, y int) color.Color {
	x /= BlkWidth
	y /= BlkHeight

	loc := (*world.World)(w).At(x, y)
	fact := float64(loc.Height + world.MaxHeight) / (2*world.MaxHeight)
	if fact > 1 {
		panic("Color factor is >1 in worldImg.At")
	}
	r, g, b, a := loc.Terrain.Color.RGBA()
	r = uint32(float64(r) * fact)
	g = uint32(float64(g) * fact)
	b = uint32(float64(b) * fact)
	return color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)}
}

// ColorModel implements the ColorModel() method
// of the image.Image interface.
func (w *worldImg) ColorModel() color.Model {
	return color.RGBAModel
}

// savePng saves the given world as a .png image.
func savePng(w *world.World) {
	out, err := os.Create("world.png")
	if err != nil {
		panic(err)
	}
	defer out.Close()
	png.Encode(out, (*worldImg)(w))
}
