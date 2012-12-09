// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"code.google.com/p/min-game/animal"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
	"encoding/json"
	"flag"
	"io"
	"math/rand"
	"os"
)

var (
	name   = flag.String("name", "Gull", "Name of the herbivores to generate")
	num    = flag.Int("num", 25, "Number to generate")
	radius = flag.Int("radius", 8, "Radius (in tiles) around start location")
)

const TileSize = 32

func main() {
	flag.Parse()

	in := bufio.NewReader(os.Stdin)
	w, err := world.Read(in)
	if err != nil {
		panic(err)
	}

	xmin := float64(w.X0-*radius) * TileSize
	xmax := float64(w.X0+*radius) * TileSize
	ymin := float64(w.Y0-*radius) * TileSize
	ymax := float64(w.Y0+*radius) * TileSize

	herbs, err := animal.MakeHerbivores(*name)
	if err != nil {
		panic(err)
	}

	maxAffinity := 0.0
	for _, a := range herbs.Info.Affinity {
		if a > maxAffinity {
			maxAffinity = a
		}
	}

	for i := 0; i < *num; i++ {
		var x, y float64
		for i := 0; i < 1000; i++ {
			x = rand.Float64()*(xmax-xmin) + xmin
			y = rand.Float64()*(ymax-ymin) + ymin
			pt := geom.Pt(x, y)
			tx, ty := w.Tile(pt.Add(world.TileSize.Div(2)))
			tname := w.At(tx, ty).Terrain.Char
			if herbs.Info.Affinity[tname] == maxAffinity {
				break
			}
		}
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		herbs.Spawn(geom.Pt(x, y), vel)
	}

	game := make(map[string]interface{})
	if err = json.NewDecoder(in).Decode(&game); err != nil && err != io.EOF {
		panic("Error reading JSON for " + *name + " " + err.Error())
	}
	if hs, ok := game["Herbivores"]; ok {
		game["Herbivores"] = append(hs.([]interface{}), herbs)
	} else {
		game["Herbivores"] = []interface{}{herbs}
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	if err := w.Write(out); err != nil {
		panic(err)
	}
	b, err := json.MarshalIndent(game, "", "\t")
	if err != nil {
		panic(err)
	}
	if _, err := out.Write(b); err != nil {
		panic(err)
	}
}
