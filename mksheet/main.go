// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
// Usage: mksheet src0 src1 src2 … dest
// All source images are assumed to have the same dimensions.
package main

import (
	"flag"
	"image"
	"image/draw"
	"image/png"
	"os"
)

var width = flag.Int("width", -1, "Max width (in tiles) of a row. A value < 1 produces one line.")

func main() {
	flag.Parse()
	tiles := flag.Args()

	if len(tiles) < 2 {
		os.Stderr.WriteString("I need some input tile files ☹\n")
		os.Exit(1)
	}

	var w, h int
	img0 := readImg(tiles[0])
	dims0 := img0.Bounds()
	var destrect image.Rectangle
	if *width < 1 {
		destrect = image.Rect(0, 0, dims0.Dx()*len(tiles), dims0.Dy())
	} else {
		w = min(*width, len(tiles))
		h = len(tiles)/w
		destrect = image.Rect(0, 0, dims0.Dx()*w, dims0.Dy()*h)
	}
	dest := image.NewRGBA(destrect)

	dp := image.Rect(0, 0, dims0.Dx(), dims0.Dy())
	draw.Draw(dest, dp, img0, dims0.Min, draw.Src)

	for i := 1; i < len(tiles)-1; i++ {
		img := readImg(tiles[i])
		x, y := i, 0
		if *width > 0 {
			x = i % w
			y = i / w
		}
		r := dp.Add(image.Pt(x*dims0.Dx(), y*dims0.Dy()))
		draw.Draw(dest, r, img, img.Bounds().Min, draw.Src)
	}

	outname := tiles[len(tiles)-1]
	out, err := os.Create(outname)
	if err != nil {
		os.Stderr.WriteString("Error creating \""+outname+"\": "+err.Error())
		os.Exit(1)
	}
	defer out.Close()

	err = png.Encode(out, dest)
	if err != nil {
		os.Stderr.WriteString("Error encoding \""+outname+"\": "+err.Error())
		os.Exit(1)
	}
}

func readImg(file string) image.Image {
	tile, err := os.Open(file)
	if err != nil {
		os.Stderr.WriteString("Error opening \""+file+"\": "+err.Error())
		os.Exit(1)
	}
	defer tile.Close()

	img, _, err := image.Decode(tile)
	if err != nil {
		os.Stderr.WriteString("Error decoding \""+file+"\": "+err.Error())
		os.Exit(1)
	}
	return img
}

func min(a, b int) int {
	if a < b { return a }
	return b
}
