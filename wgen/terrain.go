package main

import (
	"minima/world"
	"math/rand"
	"math"
	"fmt"
	"os"

)

const (
	// minMountain is the minimum value at and above which
	// terrain will initalize to mountain.
	minMountain = world.MaxHeight * 0.80

	// minWaterFrac and maxWaterFrac define the minimum and
	// maximum amount of water that will be flooded into the
	// world.  Both are given as a fraction of the map size.
	minWaterFrac, maxWaterFrac = 0.40, 0.60

	// floodMaxHeight is the maximum amount of water to flood
	// into a minima given as fraction of the world.MaxHeight
	floodMaxHeight = 0.25

	// forrestFrac is the fraction of the world covered by forrest.
	minForrestFrac, maxForrestFrac = 0.1, 0.2

	// seedFrac is the fraction of the number of forrest tiles that
	// are seeds (randomly choosen grass land to conver to
	// forrest).  The remainder of the forrest tiles are 'grown' by
	// selecting a tile adjacent to a seed.
	seedFrac = 0.0005
)

// doTerrain is the main routine for assigning a
// terrain value to each location.
func doTerrain(w *world.World) {
	initTerrain(w)
	addWater(w)
	growTrees(w)
}

// initTerrain initializes the world's terrain.
//
// Currently terrain is all initialized to grass
// land unless it is above a certain threshold
// in which case it is made a mountain.
func initTerrain(w *world.World) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			l := w.At(x, y)
			if float64(l.Height) >= minMountain {
				l.Terrain = &world.Terrain['m']
			} else {
				l.Terrain = &world.Terrain['g']
			}
		}
	}
}

// addWater adds water to the world by flooding
// some local minima to a random height.  The
// percentage of the world that is flooded is based
// by the minWaterFrac and maxWaterFrac constants.
func addWater(w *world.World) {
	tmap := makeTopoMap(w)
	minWater := int(float64(w.W*w.H)*minWaterFrac)
	maxWater := int(float64(w.W*w.H)*maxWaterFrac)
	maxHeight := int(math.Floor(world.MaxHeight*floodMaxHeight))

	waterSz := 0
	mins := tmap.minima()
	for len(mins) > 0 && waterSz < minWater {
		i := rand.Intn(len(mins))
		min := mins[i]
		mins[i], mins = mins[len(mins)-1], mins[:len(mins)-1]

		if min.terrain == &world.Terrain['w'] {
			continue
		}
	
		amt := rand.Intn(maxHeight-1)+1
		ht := min.height + amt
	
		for ht > min.height {
			fl := tmap.flood(min, ht)
			sz := 0
			for _, d := range fl {
				if d.terrain != &world.Terrain['w'] {
					sz += d.size
				}
			}
			if waterSz + sz <= maxWater {
				for _, d := range fl {
					d.terrain = &world.Terrain['w']
					d.depth += ht - d.height
					d.height = ht
				}
				waterSz += sz
			}
			ht--
		}
	}
	fmt.Fprintf(os.Stderr, "%.2f%% water\n", float64(waterSz)/float64(w.H*w.W)*100)

	// blit the water to the map
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			c := tmap.getContour(x, y)
			if c.terrain == nil {
				continue
			}
			loc := w.At(x, y)
			loc.Terrain = c.terrain
			loc.Height = c.height
			loc.Depth = c.depth
		}
	}
}

type point struct{
	x, y int
}

// growTrees changes forrest tiles into grass tiles.
func growTrees(w *world.World) {
	frac := rand.Float64()*(maxForrestFrac-minForrestFrac) + minForrestFrac
	totalForrest := int(float64(w.W)*float64(w.H)*frac)

	// build the seed locations.
	grass := grassLocs(w)
	scramble(grass)
	seeds := grass[:int(float64(totalForrest)*seedFrac)]
	for _, s := range seeds {
		w.At(s.x, s.y).Terrain = &world.Terrain['f']
	}

	n := len(seeds)
	for n < totalForrest && len(seeds) > 0 {
		i := rand.Intn(len(seeds))
		pt := seeds[i]

		var adj []point
		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if (x == 0 && y == 0) ||
					pt.x + x < 0 || pt.x + x >= w.W ||
					pt.y + y < 0 || pt.y + y >= w.H ||
					w.At(pt.x+x, pt.y+y).Terrain.Char != 'g' {
					continue
				}
				adj = append(adj, point{pt.x+x, pt.y+y})
			}
		}

		if len(adj) == 0 {
			seeds[i], seeds = grass[len(seeds)-1], grass[:len(seeds)-1]
			continue
		}

		pt = adj[rand.Intn(len(adj))]
		w.At(pt.x, pt.y).Terrain = &world.Terrain['f']
		seeds = append(seeds, pt)
		n++
	}
	fmt.Fprintf(os.Stderr, "%.2f%% forrest\n", float64(n)/float64(w.H*w.W)*100)
}

// grassLocs returns a slice of all locations that are grass land.
func grassLocs(w *world.World) (grass []point) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			l := w.At(x, y)
			if l.Terrain.Char == 'g' {
				grass = append(grass, point{x, y})
			}
		}
	}
	return
}

// scramble scrambles the given array.
func scramble(a []point) {
	for i := 0; i < len(a); i++ {
		j := rand.Intn(len(a))
		a[i], a[j] = a[j], a[i]
	}
}