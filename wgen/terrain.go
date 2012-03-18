package main

import (
	"minima/world"
	"code.google.com/p/eaburns/djsets"
	"math"
	"math/rand"
	"sort"
)

var(
	// mtnMin is the minimum value above which the
	// terrain is mountains.
	mntMin = int(math.Floor(world.MaxHeight*0.90))

	// waterMax is the maximum value below which
	// the terrain is water.
	waterMax = int(math.Floor(world.MaxHeight*0.45))
)

const (
	// lavaPercent is the percent of all lakes that become
	// lava pits instead.
	lavaPercent = 10
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
			case l.Height <=waterMax:
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
	x := w.W-1
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
	y := w.H-1
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
type Region struct{
	// size is the number of locations in this region.
	size int

	// terrain is the terrain of this region.
	terrain *world.TerrainType

	// set is this region's canonical set.
	set *djsets.Set
}

type RegionSlice []*Region

// Len implements the Len() method of sort.Interface
func (rs RegionSlice) Len() int {
	return len(rs)
}

// Less implements the Less() method of sort.Interface
func (rs RegionSlice) Less(i, j int) bool {
	return rs[i].size < rs[j].size
}

// Swap implements the Swap() method of sort.Interface
func (rs RegionSlice) Swap(i, j int) {
	rs[i], rs[j] = rs[j], rs[i]
}

// makeRegions returns a slice of Regions built from
// the connected components with the same terrain
// type.
func makeRegions(w *world.World, sets []djsets.Set) (regs RegionSlice) {
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
					size: 1,
					terrain: w.At(x, y).Terrain,
					set: set,
				}
				set.Aux = r
				regs = append(regs, r)
			}
		}
	}

	return
}

// findLakes returns all regions that are lakes
func findLakes(regs RegionSlice) (lakes RegionSlice) {
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
func addLava(w *world.World, lakes RegionSlice) {
	// all but the largest lake are candidates for lava
	if len(lakes) < 1 {
		return
	}
	sort.Sort(lakes)

	nlava := int(float64(len(lakes))*lavaPercent*0.01)
	if len(lakes) < nlava {
		nlava = len(lakes)
	}

	lava := make([]*Region,0, nlava)
	for i := 0; i < nlava; i++ {
		ind := rand.Intn(len(lakes))
		lakes[ind].terrain = &world.Terrain[int('l')]
		lava = append(lava, lakes[ind])
		lakes[ind], lakes = lakes[len(lakes)-1], lakes[:len(lakes)-1]
	}
}