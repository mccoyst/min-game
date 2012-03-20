package main

import (
	"code.google.com/p/eaburns/djsets"
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
	minWaterFrac, maxWaterFrac = 0.50, 0.75

	// maxFloodFrac is the maximum amount of water that a
	// single flood can add to the world given as a fraction of
	// the world size.
	maxFloodFrac = 0.15

	// floodMaxHeight is the maximum amount of water to flood
	// into a minima given as fraction of the world.MaxHeight
	floodMaxHeight = 0.25
)

// doTerrain clamps the heights of each cell and
// assigns their terrains.
func doTerrain(w *world.World) {
	initTerrain(w)

	tmap := makeTopoMap(w)

	addWater(w, tmap)

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
			if float64(l.Height) >= minMountain {
				l.Terrain = &world.Terrain['m']
			}
		}
	}
}

// addWater adds water to the world by flooding some local
// minima with water.
func addWater(w *world.World, tmap topoMap) {
	minWater := int(float64(w.W*w.H)*minWaterFrac)
	maxWater := int(float64(w.W*w.H)*maxWaterFrac)
	maxHeight := int(math.Floor(world.MaxHeight*floodMaxHeight))
	maxFlood := int(float64(w.W*w.H)*maxFloodFrac)

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
			if sz <= maxFlood && waterSz + sz <= maxWater {
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
	fmt.Fprintln(os.Stderr, float64(waterSz)/float64(w.H*w.W)*100, "percent water")
}

// topoMap are a set of connected components.
type topoMap struct {
	sets []djsets.Set
	stride int
	conts []*contour
}

// A contour is a connected set of locations that are
// of the same height.
type contour struct {
	// id is the unique small int that names this contour.
	id int

	// size is the number of locations in this group.
	size int

	// terrain, if non-nil, specifies the terrain type
	// for all locations of this contour.
	terrain *world.TerrainType

	// height and depth are the height and depth values
	// for all locations of this contour.
	height, depth int

	// adj is the list of adjacent countours.
	adj []*contour

	// set is the set for the canonical location.
	set *djsets.Set
}

// topoMap returns a topological map of the world.
func makeTopoMap(w *world.World) topoMap {
	sets := findSets(w)
	m := topoMap {
		stride: w.H,
		sets: sets,
		conts: findContours(w, sets),
	}
	linkContours(w, m)
	return m
}

// find returns the contour on which the given x,y point resides.
func (m topoMap) getContour(x, y int) *contour {
	return m.sets[x*m.stride+y].Find().Aux.(*contour)
}

// minima returns a slice of all contours that are local minima.
func (m topoMap) minima() (mins []*contour) {
	for _, c := range m.conts {
		min := true
		for _, a := range c.adj {
			if a.height < c.height {
				min = false
				break
			}
		}
		if min {
			mins = append(mins, c)
		}
	}
	return
}

// flood returns all of the contours that would flood
// when raising the water to the given height from the
// receiver.
func (t topoMap) flood(c *contour, ht int) (fl []*contour) {
	t.walk(c, func (c *contour) bool {
		if c.height > ht {
			return false
		}
		fl = append(fl, c)
		return true
	})
	return
}

// walk traverses the contours out from the initial
// in a depth-first order, calling foreach on each newly
// visited contour.  If foreach returns false then the
// successors of the given contour are not traversed
// unless they are reached via another path.
func (t topoMap) walk(init *contour, foreach func(*contour)bool) {
	seen := make([]bool, len(t.conts))
	var stack []*contour

	seen[init.id] = true
	stack = append(stack, init)
	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if !foreach(n) {
			continue
		}
		for _, kid := range n.adj {
			if !seen[kid.id] {
				seen[kid.id] = true
				stack = append(stack, kid)
			}
		}
	}
}

// findSets returns a slice of the set structures for each location,
// unioned into contours.
func findSets(w *world.World) (sets []djsets.Set) {
	sets = make([]djsets.Set, w.W*w.H)

	for x := 0; x < w.W-1; x++ {
		for y := 0; y < w.H-1; y++ {
			loc := w.At(x, y)
			set := &sets[x*w.H+y]
			if right := w.At(x+1, y); loc.Height == right.Height {
				set.Union(&sets[(x+1)*w.H+y])
			}
			if down := w.At(x, y+1); loc.Height == down.Height {
				set.Union(&sets[x*w.H+y+1])
			}
			if diag := w.At(x+1, y+1); loc.Height == diag.Height {
				set.Union(&sets[(x+1)*w.H+y+1])
			}
		}
	}

	// Right edge of the map (wraps to x == 0)
	x := w.W - 1
	for y := 0; y < w.H-1; y++ {
		loc := w.At(x, y)
		set := &sets[x*w.H+y]
		if right := w.At(0, y); loc.Height == right.Height {
			set.Union(&sets[y])
		}
		if down := w.At(x, y+1); loc.Height == down.Height {
			set.Union(&sets[x*w.H+y+1])
		}
		if diag := w.At(0, y+1); loc.Height == diag.Height {
			set.Union(&sets[y+1])
		}
	}

	// Bottom edge of the map (wraps to y==0)
	y := w.H - 1
	for x := 0; x < w.W-1; x++ {
		loc := w.At(x, y)
		set := &sets[x*w.H+y]
		if right := w.At(x+1, y); loc.Height == right.Height {
			set.Union(&sets[x*w.H+y])
		}
		if down := w.At(x, 0); loc.Height == down.Height {
			set.Union(&sets[x*w.H])
		}
		if diag := w.At(x+1, 0); loc.Height == diag.Height {
			set.Union(&sets[(x+1)*w.H])
		}
	}

	// Bottom left corner
	loc := w.At(x, y)
	set := &sets[x*w.H+y]
	if right := w.At(0, y); loc.Height == right.Height {
		set.Union(&sets[y])
	}
	if down := w.At(x, 0); loc.Height == down.Height {
		set.Union(&sets[x*w.H])
	}
	if diag := w.At(0, 0); loc.Height == diag.Height {
		set.Union(&sets[0])
	}
	return
}

// findContours returns a slice of all contours.
func findContours(w *world.World, sets []djsets.Set) (comps []*contour) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			loc := w.At(x, y)
			switch set := sets[x*w.H+y].Find(); {
			case set.Aux != nil:
				c := set.Aux.(*contour)
				c.size++
			default:
				c := &contour{
					id: len(comps),
					size:    1,
					height: loc.Height,
					depth: loc.Depth,
					set:     set,
				}
				set.Aux = c
				comps = append(comps, c)
			}
		}
	}
	return
}

// linkContours links the topoMap into a graph.
func linkContours(w *world.World, m topoMap) {
	for x := 0; x < w.W-1; x++ {
		for y := 0; y < w.H-1; y++ {
			c := m.getContour(x, y)
			if right := m.getContour(x+1, y); c != right {
				link(c, right)
			}
			if down := m.getContour(x, y+1); c != down {
				link(c, down)
			}
			if diag := m.getContour(x+1, y+1); c != diag {
				link(c, diag)
			}
		}
	}

	// Right edge of the map (wraps to x == 0)
	x := w.W - 1
	for y := 0; y < w.H-1; y++ {
		c := m.getContour(x, y)
		if right := m.getContour(0, y); c != right {
			link(c, right)
		}
		if down := m.getContour(x, y); c != down {
			link(c, down)
		}
		if diag := m.getContour(0, y+1); c != diag {
			link(c, diag)
		}
	}

	// Bottom edge of the map (wraps to y==0)
	y := w.H - 1
	for x := 0; x < w.W-1; x++ {
		c := m.getContour(x, y)
		if right := m.getContour(x, y+1); c != right {
			link(c, right)
		}
		if down := m.getContour(x, 0); c != down {
			link(c, down)
		}
		if diag := m.getContour(x+1, 0); c != diag {
			link(c, diag)
		}
	}

	// Bottom left corner
	c := m.getContour(x, y)
	if right := m.getContour(0, y); c != right {
		link(c, right)
	}
	if down := m.getContour(x, 0); c != down {
		link(c, down)
	}
	if diag := m.getContour(0, 0); c != diag {
		link(c, diag)
	}
	return
}

// link adds a link between the two contours if
// one did not already exist.
func link(a, b *contour) {
	for _, l := range a.adj {
		if l == b {
			return
		}
	}

	a.adj = append(a.adj, b)
	b.adj = append(b.adj, a)
}