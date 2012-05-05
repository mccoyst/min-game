package main

import (
	"minima/world"
	"math/rand"
	"math"
)

// doTerrain is the main routine for assigning a
// terrain value to each location.
func doTerrain(w *world.World) {
	start("Initializing terrain")
	initTerrain(w)
	finish()

	start("Adding water")
	addLiquid(w, 'w', 0.5, 0.6)
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
func addLiquid(w *world.World, ch uint8, minFrac, maxFrac float64) {
	tmap := makeTopoMap(w)
	minNum := int(float64(w.W*w.H)*minFrac)
	maxNum := int(float64(w.W*w.H)*maxFrac)
	maxHeight := int(math.Ceil(world.MaxElevation*0.2))

	n := 0
	mins := tmap.minima()
	for len(mins) > 0 && n < minNum {
		i := rand.Intn(len(mins))
		min := mins[i]
		mins[i], mins = mins[len(mins)-1], mins[:len(mins)-1]

		if min.terrain.Char != 'g' {
			continue
		}

		amt := 1
		if (maxHeight > 1) {
			amt = rand.Intn(maxHeight-1)+1
		}
		ht := min.height + amt
	
		for ht > min.height {
			fl := tmap.flood(min, ht)
			sz := 0
			for _, c := range fl {
				if c.terrain != &world.Terrain[ch] {
					sz += c.size
				}
			}
			if n + sz > maxNum {
				ht--
				continue
			}
			for _, c := range fl {
				c.terrain = &world.Terrain[ch]
				c.depth += ht - c.height
				c.height = ht
			}
			n += sz
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
	const minForrestFrac, maxForrestFrac = 0.08, 0.15
	frac := rand.Float64() * (maxForrestFrac - minForrestFrac) + minForrestFrac

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
