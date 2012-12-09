// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"code.google.com/p/min-game/animal"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/json"
	"code.google.com/p/min-game/world"
	"io"
	"math/rand"
	"os"
)

const TileSize = 32

func main() {
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	w, err := world.Read(bufio.NewReader(io.TeeReader(os.Stdin, out)))
	if err != nil {
		panic(err)
	}

	gulls, err := animal.MakeHerbivores("Gull")
	if err != nil {
		panic(err)
	}
	xmin, xmax := float64(w.X0-8)*TileSize, float64(w.X0+8)*TileSize
	ymin, ymax := float64(w.Y0-8)*TileSize, float64(w.Y0+8)*TileSize
	for i := 0; i < 25; i++ {
		x := rand.Float64()*(xmax-xmin) + xmin
		y := rand.Float64()*(ymax-ymin) + ymin
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		gulls.Spawn(geom.Pt(x, y), vel)
	}

	if err = json.Encode(out, "herbs", gulls); err != nil {
		panic(err)
	}

	cows, err := animal.MakeHerbivores("Cow")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 25; i++ {
		var x, y float64
		for i := 0; i < 1000; i++ {
			x = rand.Float64()*(xmax-xmin) + xmin
			y = rand.Float64()*(ymax-ymin) + ymin
			pt := geom.Pt(x, y)
			tx, ty := w.Tile(pt.Add(world.TileSize.Div(2)))
			if w.At(tx, ty).Terrain.Char == "g" {
				break
			}
		}
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		cows.Spawn(geom.Pt(x, y), vel)
	}

	if err = json.Encode(out, "herbs", cows); err != nil {
		panic(err)
	}
}
