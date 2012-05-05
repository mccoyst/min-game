package main

import (
	"math"
	"math/rand"
	"minima/world"
)

// doTerrain is the main routine for assigning a
// terrain value to each location.
func doTerrain(w *world.World) {
	start("Initializing terrain")
	initTerrain(w)
	finish()

	sz := float64(w.W * w.H)
	start("Adding oceans")
	ht := int(math.Ceil(world.MaxElevation * 0.2))
	addLiquid(w, 'w', int(sz*0.01), int(sz*0.4), int(sz*0.45), int(sz*0.6), ht)
	finish()

	start("Adding lakes")
	ht = int(math.Ceil(world.MaxElevation * 0.1))
	addLiquid(w, 'w', 7, int(sz*0.01), int(sz*0.05), int(sz*0.8), ht)
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
	const minMountain = world.MaxElevation * 0.75

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
func addLiquid(w *world.World, ch uint8, minSz, maxSz, minAmt, maxAmt, maxHt int) {
	nLiquid := 0
	tmap := makeTopoMap(w)

	mins := tmap.minima()
	for len(mins) > 0 && nLiquid < minAmt {
		i := rand.Intn(len(mins))
		min := mins[i]
		mins[i], mins = mins[len(mins)-1], mins[:len(mins)-1]

		if min.terrain.Char != 'g' {
			continue
		}

		amt := 1
		if maxHt > 1 {
			amt = rand.Intn(maxHt-1) + 1
		}
		hts := make([]int, amt)
		for i := range hts {
			hts[i] = min.height + i
		}
		for i := 0; i < len(hts)-1; i++ {
			j := rand.Intn(len(hts)-i) + i
			hts[i], hts[j] = hts[j], hts[i]
		}

		for _, ht := range hts {
			fl := tmap.flood(min, ht)
			sz := 0
			for _, c := range fl {
				if c.terrain != &world.Terrain[ch] {
					sz += c.size
				}
			}
			if sz > maxSz || nLiquid+sz > maxAmt {
				ht--
				continue
			} else if sz < minSz {
				break
			}
			for _, c := range fl {
				c.terrain = &world.Terrain[ch]
				c.depth += ht - c.height
				c.height = ht
			}
			nLiquid += sz
			break
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

type point struct {
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
	const minForrestFrac, maxForrestFrac = 0.08, 0.15
	frac := rand.Float64()*(maxForrestFrac-minForrestFrac) + minForrestFrac

	// get some seed locations.
	const seedFrac = 0.005
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
			if a.terrain.Char == 'g' && n+a.size < max {
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
