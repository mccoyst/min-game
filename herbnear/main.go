// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/mccoyst/min-game/animal"
	"github.com/mccoyst/min-game/geom"
	"github.com/mccoyst/min-game/world"
)

var (
	name   = flag.String("name", "Gull", "Name of the herbivores to generate")
	num    = flag.Int("num", 25, "Number to generate")
	radius = flag.Int("radius", 8, "Radius (in tiles) around start location")
	seed   = flag.Int64("seed", time.Now().UnixNano(), "The random seed")
)

const TileSize = 32

func main() {
	flag.Parse()
	rand.Seed(*seed)

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
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		for tries := 0; tries < 1000; tries++ {
			x := rand.Float64()*(xmax-xmin) + xmin
			y := rand.Float64()*(ymax-ymin) + ymin

			herbs.Spawn(geom.Pt(x, y), vel)
			h := herbs.Herbs[len(herbs.Herbs)-1]

			loc := w.At(w.Tile(h.Body.Center()))
			tname := loc.Terrain.Char
			if loc.Depth <= herbs.Info.BoidInfo.MaxDepth && herbs.Info.Affinity[tname] == maxAffinity {
				break
			}

			// retry
			herbs.Herbs = herbs.Herbs[:len(herbs.Herbs)-1]
		}
	}

	for _, h := range herbs.Herbs {
		tx, ty := w.Tile(h.Body.Center())
		tname := w.At(tx, ty).Terrain.Char
		if *name == "Guppy" && tname != "w" {
			panic("Err")
		}
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
