package main

import (
	"code.google.com/p/eaburns/djsets"
	"math"
	"math/rand"
	"minima/world"
)

var (
	// mtnMin is the minimum value above which the
	// terrain is mountains.
	mntMin = int(math.Floor(world.MaxHeight * 0.90))

	// waterMax is the maximum value below which
	// the terrain is water.
	waterMax = int(math.Floor(world.MaxHeight * 0.45))
)

const (
	// lavaFact is the factor of all lakes that will
	// be converted to lava pits.
	lavaFact = 0.10

	// lavaMaxFact is the upper limit, given as a factor
	// of the map area, on the size of a lava pool.
	lavaMaxFact = 0.05
)

// doTerrain clamps the heights of each cell and
// assigns their terrains.
func doTerrain(w *world.World) {
	initTerrain(w)

	sets := makeSets(w)
	regs := makeRegions(w, sets)
	addLava(w, findLakes(regs))

	finalizeTerrain(w, sets)
}

// initTerrain initializes the world's terrain by
// setting it to water or mountain depending
// on its height.
func initTerrain(w *world.World) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			l := w.At(x, y)
			if l.Height < 0 {
				l.Height = 0
			}
			if l.Height > world.MaxHeight {
				l.Height = world.MaxHeight
			}
			switch {
			case l.Height >= mntMin:
				l.Terrain = &world.Terrain['m']
			case l.Height <= waterMax:
				l.Terrain = &world.Terrain['w']
			}
		}
	}
}

// makeSets returns a slice of djsets.Sets that are
// unioned together in regions.
func makeSets(w *world.World) (sets []djsets.Set) {
	sets = make([]djsets.Set, w.W*w.H)

	for x := 0; x < w.W-1; x++ {
		for y := 0; y < w.H-1; y++ {
			loc := w.At(x, y)
			set := &sets[x*w.H+y]
			if right := w.At(x+1, y); right.Terrain == loc.Terrain {
				set.Union(&sets[(x+1)*w.H+y])
			}
			if down := w.At(x, y+1); down.Terrain == loc.Terrain {
				set.Union(&sets[x*w.H+y+1])
			}
		}
	}

	// Right edge of the map (wraps to x == 0)
	x := w.W - 1
	for y := 0; y < w.H-1; y++ {
		loc := w.At(x, y)
		set := &sets[x*w.H+y]
		if right := w.At(0, y); right.Terrain == loc.Terrain {
			set.Union(&sets[y])
		}
		if down := w.At(x, y+1); down.Terrain == loc.Terrain {
			set.Union(&sets[x*w.H+y+1])
		}
	}

	// Bottom edge of the map (wraps to y==0)
	y := w.H - 1
	for x := 0; x < w.W-1; x++ {
		loc := w.At(x, y)
		set := &sets[x*w.H+y]
		if right := w.At(x+1, y); right.Terrain == loc.Terrain {
			set.Union(&sets[x*w.H+y])
		}
		if down := w.At(x, 0); down.Terrain == loc.Terrain {
			set.Union(&sets[x*w.H])
		}
	}

	// Bottom left corner
	loc := w.At(x, y)
	set := &sets[x*w.H+y]
	if right := w.At(0, y); right.Terrain == loc.Terrain {
		set.Union(&sets[y])
	}
	if down := w.At(x, 0); down.Terrain == loc.Terrain {
		set.Union(&sets[x*w.H])
	}
	return
}

// A Region is a connected component of the world
// that has the same terrain type.
type Region struct {
	// size is the number of locations in this region.
	size int

	// terrain is the terrain of this region.
	terrain *world.TerrainType

	// set is this region's canonical set.
	set *djsets.Set
}

// makeRegions returns a slice of Regions built from
// the connected components with the same terrain
// type.
func makeRegions(w *world.World, sets []djsets.Set) (regs []*Region) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			switch set := sets[x*w.H+y].Find(); {
			case set.Aux != nil:
				r := set.Aux.(*Region)
				if r.terrain != w.At(x, y).Terrain {
					panic("A region has multiple terrains")
				}
				r.size++
			default:
				r := &Region{
					size:    1,
					terrain: w.At(x, y).Terrain,
					set:     set,
				}
				set.Aux = r
				regs = append(regs, r)
			}
		}
	}

	return
}

// findLakes returns all regions that are lakes
func findLakes(regs []*Region) (lakes []*Region) {
	for _, r := range regs {
		if r.terrain.Char == 'w' {
			lakes = append(lakes, r)
		}
	}
	return
}

// finalizeTerrain sets the terrain for each location based
// on the final terrain of its region.
func finalizeTerrain(w *world.World, sets []djsets.Set) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			set := sets[x*w.H+y].Find()
			w.At(x, y).Terrain = set.Aux.(*Region).terrain
		}
	}
}

// addLava randomly selects some lakes and makes them
// into lava.
func addLava(w *world.World, lakes []*Region) {
	maxLava := int(float64(len(lakes)) * lavaFact)
	maxSz := int(float64(w.W*w.H) * lavaMaxFact)

	for i := 0; i < maxLava && len(lakes) > 0; i++ {
		ind := rand.Intn(len(lakes))
		l := lakes[ind]
		lakes[ind], lakes = lakes[len(lakes)-1], lakes[:len(lakes)-1]

		if l.size <= maxSz {
			l.terrain = &world.Terrain[int('l')]
		}
	}
}
