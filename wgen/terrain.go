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

	comps := makeComponents(w, false, func(a, b *world.Loc) bool{
		return a.Terrain == b.Terrain
	})

	var lakes []*component
	for _, c := range comps.comps {
		if c.loc.Terrain.Char == 'w' {
			lakes = append(lakes, c)
		}
	}

	maxLava := int(float64(len(lakes)) * lavaFact)
	maxSz := int(float64(w.W*w.H) * lavaMaxFact)
	for i := 0; i < maxLava && len(lakes) > 0; i++ {
		ind := rand.Intn(len(lakes))
		c := lakes[ind]
		lakes[ind], lakes = lakes[len(lakes)-1], lakes[:len(lakes)-1]

		if c.size <= maxSz {
			c.loc.Terrain = &world.Terrain[int('l')]
		}
	}

	finalizeTerrain(w, comps)
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

// finalizeTerrain sets the terrain for each location based
// on the final terrain of its component.
func finalizeTerrain(w *world.World, comps components) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			c := comps.find(x, y)
			w.At(x, y).Terrain = c.loc.Terrain
		}
	}
}

// components are a set of connected components.
type components struct {
	sets []djsets.Set
	stride int
	comps []*component
}

// A component is some set of locations that can be represented
// by a canonical location.
type component struct {
	// size is the number of locations in this group.
	size int

	// minHt and maxHt are the extreme
	// heights of the region.
	minHt, maxHt int

	// loc is the canonical location for this group.
	// All other locations have something similar
	// to the canonical location.
	loc *world.Loc

	// set is the set for the canonical location.
	set *djsets.Set
}

// sameComponent tests if two adjacent locations
// should fall within the same component.
type sameComponent func(a, b *world.Loc)bool

// find returns the component for the given location.
func (c components) find(x, y int) *component {
	return c.sets[x*c.stride+y].Find().Aux.(*component)
}

// makeComponents returns a components containing all 
// connected components for which p evaluates to true
// on adjacent locations.
//
// If d is true then diagonals are considered adjacent.
func makeComponents(w *world.World, d bool, p sameComponent) components {
	sets := findSets(w, d, p)
	return components{
		comps: findComps(w, sets),
		stride: w.H,
		sets: sets,
	}
}

// findSets returns a slice of djsets.Sets containing all
// connected components where p evaluates to true for
// adjacent of locations.
//
// If d is false then Diagonals are not considered adjacent.
func findSets(w *world.World, d bool, p sameComponent) (sets []djsets.Set) {
	sets = make([]djsets.Set, w.W*w.H)

	for x := 0; x < w.W-1; x++ {
		for y := 0; y < w.H-1; y++ {
			loc := w.At(x, y)
			set := &sets[x*w.H+y]
			if right := w.At(x+1, y); p(loc, right) {
				set.Union(&sets[(x+1)*w.H+y])
			}
			if down := w.At(x, y+1); p(loc, down) {
				set.Union(&sets[x*w.H+y+1])
			}
			if diag := w.At(x+1, y+1); d && p(loc, diag) {
				set.Union(&sets[(x+1)*w.H+y+1])
			}
		}
	}

	// Right edge of the map (wraps to x == 0)
	x := w.W - 1
	for y := 0; y < w.H-1; y++ {
		loc := w.At(x, y)
		set := &sets[x*w.H+y]
		if right := w.At(0, y); p(loc, right) {
			set.Union(&sets[y])
		}
		if down := w.At(x, y+1); p(loc, down) {
			set.Union(&sets[x*w.H+y+1])
		}
		if diag := w.At(0, y+1); d && p(loc, diag) {
			set.Union(&sets[y+1])
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
		if diag := w.At(x+1, 0); d && p(loc, diag) {
			set.Union(&sets[(x+1)*w.H])
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
	if diag := w.At(0, 0); d && p(loc, diag) {
		set.Union(&sets[0])
	}
	return
}

// findComps returns a slice of all components,
// where each component has one canonical location.
func findComps(w *world.World, sets []djsets.Set) (comps []*component) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			switch set := sets[x*w.H+y].Find(); {
			case set.Aux != nil:
				c := set.Aux.(*component)
				l := w.At(x, y)
				if l.Height < c.minHt {
					c.minHt = l.Height
				}
				if l.Height > c.maxHt {
					c.maxHt = l.Height
				}
				c.size++
			default:
				l := w.At(x, y)
				c := &component{
					size:    1,
					minHt: l.Height,
					maxHt: l.Height,
					loc: l,
					set:     set,
				}
				set.Aux = c
				comps = append(comps, c)
			}
		}
	}
	return
}

