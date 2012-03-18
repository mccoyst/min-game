package main

import (
	"bufio"
	"flag"
	"image"
	"image/color"
	"image/png"
	"minima/world"
	"os"
)

var (
	outFile = flag.String("o", "world.png", "The output file")
	echo    = flag.Bool("e", false, "Echo the world to standard output")
	blkSize = flag.Int("s", 1, "Size of each block in pixels")
)

func main() {
	flag.Parse()

	w, err := world.Read(bufio.NewReader(os.Stdin))
	if err != nil {
		panic(err)
	}

	imgout, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	defer imgout.Close()
	png.Encode(imgout, (*worldImg)(&w))

	if *echo {
		out := bufio.NewWriter(os.Stdout)
		defer out.Flush()
		if err := w.Write(out); err != nil {
			panic(err)
		}
	}
}

type worldImg world.World

// Bounds implements the Bounds() method of
// the image.Image interface.
func (w *worldImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.W*(*blkSize), w.H*(*blkSize))
}

// At implements the At() method of the
// image.Image interface.
func (w *worldImg) At(x, y int) color.Color {
	x /= *blkSize
	y /= *blkSize
	loc := (*world.World)(w).At(x, y)
	f := float64(loc.Height+world.MaxHeight) / (2 * world.MaxHeight)
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
