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
	minMountain = world.MaxElevation * 0.80

	// minWaterFrac and maxWaterFrac define the minimum and
	// maximum amount of water that will be flooded into the
	// world.  Both are given as a fraction of the map size.
	minWaterFrac, maxWaterFrac = 0.40, 0.60

	// floodMaxElevation is the maximum amount of water to flood
	// into a minima given as fraction of the world.MaxElevation
	floodMaxElevation = 0.25

	// forrestFrac is the fraction of the world covered by forrest.
	minForrestFrac, maxForrestFrac = 0.2, 0.4

	// seedFrac is the fraction of the number of forrest tiles that
	// are seeds (randomly choosen grass land to conver to
	// forrest).  The remainder of the forrest tiles are 'grown' by
	// selecting a tile adjacent to a seed.
	seedFrac = 0.005
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
			if float64(l.Elevation) >= minMountain {
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
	maxHeight := int(math.Floor(world.MaxElevation*floodMaxElevation))

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
			loc := w.At(x, y)
			loc.Terrain = c.terrain
			loc.Elevation = c.height
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

	tmap := makeTopoMap(w)
	var grass []*contour
	for _, c := range tmap.conts {
		if c.terrain.Char == 'g' {
			grass = append(grass, c)
		}
	}

	// scramble
	for i := 0; i < len(grass); i++ {
		j := rand.Intn(len(grass))
		grass[i], grass[j] = grass[j], grass[i]
	}

	// build the seed locations.
	n := 0
	seeds := grass[:int(float64(w.W)*float64(w.H)*frac*seedFrac)]
	for _, s := range seeds {
		s.terrain = &world.Terrain['f']
		n += s.size
	}
	fmt.Fprintf(os.Stderr, "%.2f%% forrest\n", float64(n)/float64(w.H*w.W)*100)

	// blit the forrest to the map
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			c := tmap.getContour(x, y)
			loc := w.At(x, y)
			loc.Terrain = c.terrain
			loc.Elevation = c.height
			loc.Depth = c.depth
		}
	}
}
