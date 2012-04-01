package main

import (
	"minima/world"
	"code.google.com/p/eaburns/djsets"
)

// topoMap are a set of connected components of
// equal heights.
type topoMap struct {
	sets []djsets.Set
	stride int
	conts []*contour
}

// A contour represents connected set of locations
// that have the same height.
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

// getContour returns the contour on which the given
// x,y point resides.
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
// when raising raising liquid to the given height
// in the given contour.
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
// successors of the given contour are not traversed on
// this branch of the search and will only be encountered
// if they are reached via another path along which foreach
// always returns true.
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

// findSets returns a slice of set, one for each
// location, that are unioned together such that
// sets with matching canonical representations
// all belong to the same contour.
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

// linkContours links adjacent contours.
func linkContours(w *world.World, m topoMap) {
	for x := 0; x < w.W-1; x++ {
		for y := 0; y < w.H-1; y++ {
			c := m.getContour(x, y)
			if right := m.getContour(x+1, y); c != right {
				addLink(c, right)
			}
			if down := m.getContour(x, y+1); c != down {
				addLink(c, down)
			}
			if diag := m.getContour(x+1, y+1); c != diag {
				addLink(c, diag)
			}
		}
	}

	// Right edge of the map (wraps to x == 0)
	x := w.W - 1
	for y := 0; y < w.H-1; y++ {
		c := m.getContour(x, y)
		if right := m.getContour(0, y); c != right {
			addLink(c, right)
		}
		if down := m.getContour(x, y); c != down {
			addLink(c, down)
		}
		if diag := m.getContour(0, y+1); c != diag {
			addLink(c, diag)
		}
	}

	// Bottom edge of the map (wraps to y==0)
	y := w.H - 1
	for x := 0; x < w.W-1; x++ {
		c := m.getContour(x, y)
		if right := m.getContour(x, y+1); c != right {
			addLink(c, right)
		}
		if down := m.getContour(x, 0); c != down {
			addLink(c, down)
		}
		if diag := m.getContour(x+1, 0); c != diag {
			addLink(c, diag)
		}
	}

	// Bottom left corner
	c := m.getContour(x, y)
	if right := m.getContour(0, y); c != right {
		addLink(c, right)
	}
	if down := m.getContour(x, 0); c != down {
		addLink(c, down)
	}
	if diag := m.getContour(0, 0); c != diag {
		addLink(c, diag)
	}
	return
}

// addLink adds a link between the two contours if
// one did not already exist.
func addLink(a, b *contour) {
	for _, l := range a.adj {
		if l == b {
			return
		}
	}

	a.adj = append(a.adj, b)
	b.adj = append(b.adj, a)
}
