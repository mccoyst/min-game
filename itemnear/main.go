// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"code.google.com/p/min-game/item"
	"code.google.com/p/min-game/world"
	"encoding/json"
	"flag"
	"io"
	"math/rand"
	"os"
	"time"
)

var (
	name   = flag.String("name", "Uranium", "Name of the item to generate")
	num    = flag.Int("num", 1, "Number to generate")
	radius = flag.Int("radius", 4, "Radius (in tiles) around start location")
	seed   = flag.Int64("seed", time.Now().UnixNano(), "The random seed")
)

func main() {
	flag.Parse()
	rand.Seed(*seed)

	in := bufio.NewReader(os.Stdin)
	w, err := world.Read(in)
	if err != nil {
		panic(err)
	}

	r := *radius
	var items []interface{}
	for i := 0; i < *num; i++ {
		x := float64(w.X0+(rand.Intn(2*r)-r)) * world.TileSize.X
		y := float64(w.Y0+(rand.Intn(2*r)-r)) * world.TileSize.Y
		it := item.NewTreasure(x, y, item.New(*name))
		if it == nil {
			panic("Unknown item name: " + *name)
		}
		items = append(items, it)
	}

	game := make(map[string]interface{})
	if err = json.NewDecoder(in).Decode(&game); err != nil && err != io.EOF {
		panic("Error reading JSON for " + *name + " " + err.Error())
	}

	if treasure, ok := game["Treasure"]; ok {
		items = append(treasure.([]interface{}), items...)
	}
	game["Treasure"] = items

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
