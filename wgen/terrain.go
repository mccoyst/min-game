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
	minMountain = world.MaxHeight * 0.90

	// minWaterFrac and maxWaterFrac define the minimum and
	// maximum amount of water that will be flooded into the
	// world.  Both are given as a fraction of the map size.
	minWaterFrac, maxWaterFrac = 0.40, 0.60

	// floodMaxHeight is the maximum amount of water to flood
	// into a minima given as fraction of the world.MaxHeight
	floodMaxHeight = 0.25
)

// doTerrain is the main routine for assigning a
// terrain value to each location.
func doTerrain(w *world.World) {
	initTerrain(w)
	addWater(w)
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
	fmt.Fprintln(os.Stderr, float64(waterSz)/float64(w.H*w.W)*100,
		"percent water")

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
