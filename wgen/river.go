package main

import (
	"code.google.com/p/eaburns/perlin"
	"code.google.com/p/min-game/world"
	"container/heap"
	"math"
	"math/rand"
)

// addRivers adds rivers
// minSz gives the minimum river size and maxCnt is
// the maximum number of locations to add as rivers.
func addRivers(w *world.World, oceans []world.Coord, minSz, maxCnt int) {
	isOcean := make([]bool, w.W*w.H)
	for _, coord := range oceans {
		isOcean[coord.X*w.H+coord.Y] = true
	}

	sources := riverSources(w, isOcean)
	if len(sources) == 0 {
		return
	}

	noise := makeNoise(w)
	cnt := 0
	for cnt < maxCnt && len(sources) > 0 {
		i := rand.Intn(len(sources))
		src := sources[i]
		sources[i], sources = sources[len(sources)-1], sources[:len(sources)-1]

		river := riverLocs(w, noise, isOcean, minSz, src)
		if len(river) < minSz {
			continue
		}
		for _, node := range river {
			node.loc.Terrain = &world.Terrain[int('w')]
			if node.loc.Depth <= 0 {
				node.loc.Depth = 1
			}
			cnt++
		}
	}
}

// deltas is the Δx and Δy from a location to its neighbors.
var deltas = []struct{ dx, dy int }{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

// riverLocs returns a slice of coordinates that form a river.
func riverLocs(w *world.World, noise []float64, isOcean []bool, minSz int, src world.Coord) []*riverNode {
	init := rn(w, noise, src, nil)
	nodes := make([]*riverNode, w.W*w.H)
	nodes[src.X*w.H+src.Y] = init
	q := riverHeap{init}

	for len(q) > 0 {
		n := heap.Pop(&q).(*riverNode)
		if isOcean[n.X*w.H+n.Y] {
			if n.len() < minSz {
				continue
			}
			return n.path()
		}

		for _, d := range deltas {
			x, y := w.Wrap(n.X+d.dx, n.Y+d.dy)
			kid := nodes[x*w.H+y]
			if kid == nil {
				kid = rn(w, noise, world.Coord{x, y}, n)
				nodes[x*w.H+y] = kid
			} else if kid.g <= kid.edgecost+n.g {
				continue
			}
			if kid.loc.Elevation > n.loc.Elevation {
				continue
			}
			if kid.g > kid.edgecost+n.g {
				kid.g = kid.edgecost + n.g
				kid.parent = n
				if kid.pqind >= 0 {
					heap.Remove(&q, kid.pqind)
				}
			}
			heap.Push(&q, kid)
		}
	}

	return []*riverNode{}
}

type riverHeap []*riverNode

func (h *riverHeap) Push(x interface{}) {
	n := x.(*riverNode)
	n.pqind = len(*h)
	*h = append(*h, n)
}

func (h *riverHeap) Pop() interface{} {
	heap := *h
	n := heap[len(heap)-1]
	n.pqind = -1
	*h = heap[:len(heap)-1]
	return n
}

func (h riverHeap) Less(i, j int) bool {
	return h[i].g < h[j].g
}

func (h riverHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].pqind = i
	h[j].pqind = j
}

func (h riverHeap) Len() int {
	return len(h)
}

// A riverNode is a single location on a river.
type riverNode struct {
	world.Coord
	loc         *world.Loc
	parent      *riverNode
	g, edgecost float64
	pqind       int
}

// rn returns a new river node.
func rn(w *world.World, noise []float64, c world.Coord, p *riverNode) *riverNode {
	e := noise[c.X*w.H+c.Y]
	g := e
	if p != nil {
		g += p.g
	}
	return &riverNode{
		world.Coord: c,
		loc:         w.At(c.X, c.Y),
		parent:      p,
		edgecost:    e,
		g:           g,
		pqind:       -1,
	}
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
func riverSources(w *world.World, isOcean []bool) (sources []world.Coord) {
	min := minOcean(w, isOcean)
	for _, coord := range w.WithType("m") {
		mtn := w.At(coord.X, coord.Y)
		if mtn.Elevation >= min {
			sources = append(sources, coord)
		}
	}

	for _, coord := range w.WithType("w") {
		wtr := w.At(coord.X, coord.Y)
		if !isOcean[coord.X*w.H+coord.Y] && wtr.Elevation >= min {
			sources = append(sources, coord)
		}
	}

	return
}

// minWminOceanater returns the minimum ocean elevation in the world.
func minOcean(w *world.World, isOcean []bool) int {
	min := world.MaxElevation
	for _, coord := range w.WithType("w") {
		wtr := w.At(coord.X, coord.Y)
		if isOcean[coord.X*w.H+coord.Y] && wtr.Elevation < min {
			min = wtr.Elevation
		}
	}
	return min
}

// makeNoise makes a slice of normalized Perlin noise values.
func makeNoise(w *world.World) []float64 {
	noise := make([]float64, w.W*w.H)
	perlin := perlin.Make(0.8, 0.25, 2, rand.Int63(), nil)
	min, max := math.Inf(1), 0.0
	for i := range noise {
		x, y := i/w.H, i%w.H
		n := perlin(float64(x), float64(y))
		noise[i] = n
		if n < min {
			min = n
		}
		if n > max {
			max = n
		}
	}
	for i := range noise {
		noise[i] = 2 * (noise[i] - min) / (max - min)
		if noise[i] <= 0.01 {
			noise[i] = 0.01
		}
	}
	return noise
}
