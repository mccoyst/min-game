// Small test program that draws an image of
// some text in a given font.

// +build ignore

package main

import (
	"flag"
	"image/color"
	"image/png"
	"os"

	"github.com/mccoyst/min-game/ui"
)

var (
	path = flag.String("f", "../resrc/prstartk.ttf", "The TTF file path")
	size = flag.Float64("s", 72.0, "The font size in points")
)

func main() {
	flag.Parse()

	font, err := ui.NewFont(*path)
	if err != nil {
		panic(err)
	}
	font.SetSize(*size)
	font.SetColor(color.Black)
	img, err := font.Render("Eloquent M")
	if err != nil {
		panic(err)
	}

	f, err := os.Create("img.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		panic(err)
	}
}
