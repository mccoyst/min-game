// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"code.google.com/p/min-game/world"
	"flag"
	"image"
	"image/color"
	"image/png"
	"os"
)

var (
	outFile = flag.String("o", "world.png", "The output file")
	echo    = flag.Bool("e", false, "Echo the world to standard output")
	depth   = flag.Bool("d", true, "Draw water depth")
)

func main() {
	flag.Parse()

	w, err := world.Read(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}

	out, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	png.Encode(out, &worldImg{w, *depth})

	if *echo {
		out := bufio.NewWriter(os.Stdout)
		defer out.Flush()
		if err := w.Write(out); err != nil {
			panic(err)
		}
	}
}

type worldImg struct {
	*world.World
	depth bool
}

// Bounds implements the Bounds() method of
// the image.Image interface.
func (w *worldImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.W, w.H)
}

var (
	colors = []color.RGBA{
		'g': color.RGBA{109, 170, 44, 255},
		'm': color.RGBA{210, 125, 44, 255},
		'w': color.RGBA{109, 194, 202, 255},
		'l': color.RGBA{208, 70, 72, 255},
		'd': color.RGBA{218, 219, 94, 255},
		'f': color.RGBA{52, 101, 36, 255},
		'i': color.RGBA{222, 238, 214, 255},
	}
)

// At implements the At() method of the
// image.Image interface.
func (w *worldImg) At(x, y int) color.Color {
	loc := w.World.At(x, y)
	el := loc.Elevation
	if w.depth {
		el -= loc.Depth
	}
	f := float64(el+world.MaxElevation) / (2 * world.MaxElevation)
	if f > 1 {
		panic("Color factor is >1 in worldImg.At")
	}
	c := colors[loc.Terrain.Char]
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
