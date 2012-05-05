package main

import (
	"minima/world"
	"math/rand"
	"math"

)

const (
	// minMountain is the minimum value at and above which
	// terrain will initalize to mountain.
	minMountain = world.MaxElevation * 0.80

	// minWaterFrac and maxWaterFrac define the minimum and
	// maximum amount of water that will be flooded into the
	// world.  Both are given as a fraction of the map size.
	minWaterFrac, maxWaterFrac = 0.40, 0.60

	// minLavaFrac and maxLavaFrac define the minimum and
	// maximum amount of lava that will be flooded into the
	// world.  Both are given as a fraction of the map size.
	minLavaFrac, maxLavaFrac = 0.005, 0.01

	// floodMaxElevation is the maximum amount of water to flood
	// into a minima given as fraction of the world.MaxElevation
	floodMaxElevation = 0.25

	// minForrestFrac and maxForrestFrac give rough bounds on the
	// amount of forrest that can be added to the world.
	minForrestFrac, maxForrestFrac = 0.08, 0.15

	// seedFrac specifies the number of seed forrest.  This is given
	// as a fraction of the number of grass contours.
	seedFrac = 0.005
)

// doTerrain is the main routine for assigning a
// terrain value to each location.
func doTerrain(w *world.World) {
	start("Initializing terrain")
	initTerrain(w)
	finish()

	start("Adding water")
	addLiquid(w, 'w', minWaterFrac, maxWaterFrac)
	finish()

	start("Growing forrests")
	growTrees(w)
	finish()
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

// addLiquid adds some liquid (given by ch)  to the
// world by flooding some local minima to a random
// height.  The  percentage of the world that is flooded is
// based on the minFrac and maxFrac parameters.
func addLiquid(w *world.World, ch uint8, minFrac, maxFrac float64) {
	tmap := makeTopoMap(w)
	minNum := int(float64(w.W*w.H)*minFrac)
	maxNum := int(float64(w.W*w.H)*maxFrac)
	maxHeight := int(math.Floor(world.MaxElevation*floodMaxElevation))

	n := 0
	mins := tmap.minima()
	for len(mins) > 0 && n < minNum {
		i := rand.Intn(len(mins))
		min := mins[i]
		mins[i], mins = mins[len(mins)-1], mins[:len(mins)-1]

		if min.terrain.Char != 'g' {
			continue
		}
	
		amt := rand.Intn(maxHeight-1)+1
		ht := min.height + amt
	
		for ht > min.height {
			fl := tmap.flood(min, ht)
			sz := 0
			for _, d := range fl {
				sz += d.size
			}
			if n + sz <= maxNum {
				for _, d := range fl {
					d.terrain = &world.Terrain[ch]
					d.depth += ht - d.height
					d.height = ht
				}
				n += sz
			}
			ht--
		}
	}

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
	tmap := makeTopoMap(w)
	var grass []*contour
	for _, c := range tmap.conts {
		if c.terrain.Char == 'g' {
			grass = append(grass, c)
		}
	}

	// scramble
	for i := 0; i < len(grass)-1; i++ {
		j := rand.Intn(len(grass)-i) + i
		grass[i], grass[j] = grass[j], grass[i]
	}

	n := 0
	frac := rand.Float64() * (maxForrestFrac - minForrestFrac) + minForrestFrac

	// get some seed locations.
	seeds := grass[:int(float64(w.W)*float64(w.H)*frac*seedFrac)]
	for _, s := range seeds {
		s.terrain = &world.Terrain['f']
		n += s.size
	}

	min := int(frac * float64(w.W) * float64(w.H))
	max := int(maxForrestFrac * float64(w.W) * float64(w.H))
	for len(seeds) > 0 && n < min {
		i := rand.Intn(len(seeds))
		c := seeds[i]
		var adj []*contour
		for _, a := range c.adj {
			if a.terrain.Char == 'g' && n + a.size < max {
				adj = append(adj, a)
			}
		}
		if len(adj) == 0 {
			seeds[i], seeds = seeds[len(seeds)-1], seeds[:len(seeds)-1]
			continue
		}
		c = adj[rand.Intn(len(adj))]
		n += c.size
		c.terrain = &world.Terrain['f']
		seeds = append(seeds, c)
	}

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
