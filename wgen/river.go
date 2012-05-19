package main

import (
	"minima/world"
	"math/rand"
)

// addRivers adds rivers
// minSz gives the minimum river size and maxCnt is
// the maximum number of locations to add as rivers.
func addRivers(w *world.World, oceans []*world.Loc, minSz, maxCnt int) {
	isOcean := make(map[*world.Loc]bool, w.W*w.H)
	for _, l := range oceans {
		isOcean[l] = true
	}

	sources := riverSources(w, isOcean)
	if len(sources) == 0 {
		return
	}

	cnt := 0
	for cnt < maxCnt {
		i := rand.Intn(len(sources))
		src := sources[i]
		sources[i], sources = sources[len(sources)-1], sources[:len(sources)-1]

		river := riverLocs(w, isOcean, minSz, src)
		if len(river) < minSz {
			continue
		}
		for _, node := range river {
			node.loc.Terrain = &world.Terrain[int('w')]
			cnt++
		}
	}
}

// riverLocs returns a slice of coordinates that form a river.
func riverLocs(w *world.World, isOcean map[*world.Loc]bool, minSz int, src world.Coord) []*riverNode {
	init := &riverNode{
		world.Coord: src,
		loc: w.At(src.X, src.Y),
		cost: 0,
		onq: true,
	}
	nodes := make(map[*world.Loc]*riverNode, w.W*w.H)
	nodes[init.loc] = init
	q := []*riverNode{ init }

	for len(q) > 0 {
		n := q[0]
		q = q[1:]
		n.onq = false;

		if isOcean[n.loc] && n.len() >= minSz {
			return n.path()
		}

		for i, d := range deltas {
			x, y := n.X + d.dx, n.Y + d.dy
			kidloc := w.AtCoord(x, y)
			if n.edgecosts[i] == 0 {
				n.edgecosts[i] = float64(rand.Intn(5)+1)
				if kidloc.Elevation < n.loc.Elevation {
					n.edgecosts[i] *= 0.1
				}
			}
			cost := n.edgecosts[i]

			kid, ok := nodes[kidloc]
			if !ok {
				kid = &riverNode {
					world.Coord: world.Coord{ x, y },
					loc: kidloc,
					parent: n,
					cost: n.cost + cost,
					onq: true,
				}
				nodes[kidloc] = kid
				q = append(q, kid)
				continue
			}
			if kidloc.Elevation > n.loc.Elevation || kid.cost <= n.cost + cost {
				continue
			}

			kid.cost = n.cost + cost
			kid.parent = n
			if !kid.onq {
				kid.onq = true
				q = append(q, kid)
			}
		}
	}
	return []*riverNode{}
}

var (
	// deltas is the Δx and Δy from a location to its neighbors.
	deltas = []struct{ dx, dy int } {
		{ 1, 0 },
		{ -1, 0 },
		{ 0, 1 },
		{ 0, -1 },
	}
)

// A riverNode is a single location on a river.
type riverNode struct {
	world.Coord
	loc *world.Loc
	parent *riverNode
	cost float64
	onq bool
	edgecosts [4]float64
}

// len returns the length of the path from this node back
// to the source
func (n *riverNode) len() int {
	l := 1
	for n.parent != nil {
		l++
		n = n.parent
	}
	return l
}

// path returns the path from this node back to the source
func (n *riverNode) path() (path []*riverNode) {
	for n.parent != nil {
		path = append(path, n)
		n = n.parent
	}
	path = append(path, n)
	return
}

// riverSources returns a scrambled list of all possible
// source locations for a river.
func riverSources(w *world.World, isOcean map[*world.Loc]bool) (sources []world.Coord) {
	min := minOcean(w, isOcean)
	for _, coord := range w.WithType("m") {
		mtn := w.At(coord.X, coord.Y)
		if mtn.Elevation >= min {
			sources = append(sources, coord)
		}
	}

	for _, coord := range w.WithType("w") {
		wtr := w.At(coord.X, coord.Y)
		if !isOcean[wtr] && wtr.Elevation >= min {
			sources = append(sources, coord)
		}
	}

	return
}

// minWminOceanater returns the minimum ocean elevation in the world.
func minOcean(w *world.World, isOcean map[*world.Loc]bool) int {
	min := world.MaxElevation
	for _, coord := range w.WithType("w") {
		water := w.At(coord.X, coord.Y)
		if isOcean[water] && water.Elevation < min {
			min = water.Elevation
		}
	}
	return min
}
