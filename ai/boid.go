// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ai

import (
	"math"
	"strings"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/world"
)

// Boids is a list of boids that may flock together.
type Boids interface {
	Len() int
	Boid(int) *Boid
	BoidInfo() BoidInfo
}

const (
	// NThinkGroups is the number of think groups.
	NThinkGroups = 6 // At 60 fps, this gives us thinking at 10Hz.
)

// A Boid is a bird-like (cow-like, or fish-like) object.
type Boid struct {
	Body phys.Body

	// ThinkGroup is the number of the group with which this boid
	// considers its local neighbors when updating.
	ThinkGroup uint
}

// BoidInfo contains behavior information about boids.
// The XxxBias terms are fairly arbitrary weights that are
// applied to boid rules in order to prioritize them.
// Unfortunately, the only good way to set them seems to
// be painful trial-and-error.
type BoidInfo struct {
	// MaxVelocity is the fastest that a boid can move.
	MaxVelocity float64

	// LocalDist is the distance determining when two boids
	// are flocking together.
	LocalDist float64

	// MatchBias is multiplied by the velocity matching
	// velocity.
	MatchBias float64

	// CenterDist is the minimum distance a boid must be
	// away from the center before they try to center
	// themselves.
	CenterDist float64

	// CenterBias is multiplied by the centering velocity.
	CenterBias float64

	// AvoidDist is the distance at which boids will avoid
	// eachother.
	AvoidDist float64

	// AvoidBias is multiplied by the avoid velocity.
	AvoidBias float64

	// PlayerDist is the distance at which boids will avoid
	// the player. 
	PlayerDist float64

	// PlayerBias is the weight applied to avoiding the player.
	PlayerBias float64

	// TerrainDist is the distance at which terrain is avoided.
	TerrainDist float64

	// TerrainBias is the weight applied to avoiding terrain.
	TerrainBias float64

	// AvoidTerrain is a string of terrain types that are avoided.
	AvoidTerrain string

	// MaxDepth is the maximum water depth before this boid
	// attempts to avoid.
	MaxDepth int
}

// UpdateBoids updates the velocity of the boids.
func UpdateBoids(nframes uint, boids Boids, p *phys.Body, w *world.World) {
	info := boids.BoidInfo()
	local := localBoids(nframes, boids, w)
	for i, l := range local {
		boid := boids.Boid(i)
		boid.matchVel(l, info)
		boid.moveCenter(l, info, w)
		boid.avoidOthers(l, info, w)
		boid.avoidPlayer(p, info, w)
		boid.avoidTerrain(info, w)
		boid.clampVel(info.MaxVelocity)
	}
}

// LocalBoids returns a slice containing the Boids that
// are local to the Boid with the corresponding index.
func localBoids(nframes uint, boids Boids, w *world.World) [][]*Boid {
	g := newGrid(10, 10, w.Pixels)
	for i := 0; i < boids.Len(); i++ {
		b := boids.Boid(i)
		c := g.index(g.pt2Cell(b.Body.Box.Min))
		g.cells[c] = append(g.cells[c], b)
	}

	dist := boids.BoidInfo().LocalDist
	tGroup := nframes % NThinkGroups
	local := make([][]*Boid, boids.Len())
	for i := range local {
		boid := boids.Boid(i)
		if tGroup != boid.ThinkGroup {
			continue
		}
		local[i] = g.neighbors(boid, dist)
	}
	return local
}

// MatchVel attempts to match the velocity of the local boids.
func (boid *Boid) matchVel(local []*Boid, i BoidInfo) {
	var avg geom.Point
	for _, b := range local {
		if b == boid {
			continue
		}
		avg = avg.Add(b.Body.Vel)
	}
	if len(local) == 0 {
		return
	}
	avg = avg.Div(float64(len(local))).Normalize().Mul(i.MatchBias)
	boid.Body.Vel = boid.Body.Vel.Add(avg)
}

// MoveCenter attempts to move the boid toward the center
// of its local flock mates.
func (boid *Boid) moveCenter(local []*Boid, i BoidInfo, w *world.World) {
	var avg, c geom.Point
	for _, b := range local {
		if b == boid {
			continue
		}
		toCenter := w.Pixels.Sub(b.Body.Box.Min, boid.Body.Box.Min)
		c = c.Add(b.Body.Box.Min)
		avg = avg.Add(toCenter)
	}
	if len(local) == 0 {
		return
	}
	n := float64(len(local))
	c = c.Div(n)
	if w.Pixels.SqDist(c, boid.Body.Box.Min) < i.CenterDist*i.CenterDist {
		return
	}
	avg = avg.Div(n).Normalize().Mul(i.CenterBias)
	boid.Body.Vel = boid.Body.Vel.Add(avg)
}

// AvoidOthers attempts to avoid very close flock mates.
func (boid *Boid) avoidOthers(local []*Boid, i BoidInfo, w *world.World) {
	dd := i.AvoidDist * i.AvoidDist
	var a geom.Point
	for _, b := range local {
		if b == boid {
			continue
		}
		if d := boid.sqDist(b, w.Pixels); d > dd {
			continue
		}
		a = a.Add(avoidVec(boid.Body.Center(), b.Body.Center(), i.AvoidDist, w))
	}
	a = a.Mul(i.AvoidBias)
	boid.Body.Vel = boid.Body.Vel.Add(a)
}

// AvoidPlayer attempts to avoid the player.
func (boid *Boid) avoidPlayer(p *phys.Body, i BoidInfo, w *world.World) {
	dd := i.PlayerDist * i.PlayerDist
	pt := p.Box.Center()
	if p.Vel == geom.Pt(0, 0) || w.Pixels.SqDist(boid.Body.Box.Center(), pt) > dd {
		return
	}
	d := avoidVec(boid.Body.Box.Center(), pt, i.PlayerDist, w).Mul(i.PlayerBias)
	boid.Body.Vel = boid.Body.Vel.Add(d)
}

// AvoidTerrain attempts to avoid certain types of terrain.
func (boid *Boid) avoidTerrain(i BoidInfo, w *world.World) {
	if i.AvoidTerrain == "" {
		return
	}

	var a geom.Point
	dd := i.TerrainDist * i.TerrainDist
	sz := geom.Pt(i.TerrainDist, i.TerrainDist)
	x0, y0 := w.Tile(boid.Body.Box.Min.Sub(sz))
	x1, y1 := w.Tile(boid.Body.Box.Min.Add(sz))

	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			l := w.At(x, y)
			ch := l.Terrain.Char
			if l.Depth <= i.MaxDepth && strings.Index(i.AvoidTerrain, ch) < 0 {
				continue
			}
			pt := geom.Pt((float64(x)+0.5)*world.TileSize.X,
				(float64(y)+0.5)*world.TileSize.Y)
			if w.Pixels.SqDist(boid.Body.Box.Center(), pt) > dd {
				continue
			}
			a = a.Add(avoidVec(boid.Body.Box.Center(), pt, i.TerrainDist, w))
		}
	}
	a = a.Mul(i.TerrainBias)
	boid.Body.Vel = boid.Body.Vel.Add(a)
}

// ClampVel clamps the boid's velocity to have a magnitude of
// no more than max.
func (boid *Boid) clampVel(max float64) {
	if boid.Body.Vel.Len() > max {
		boid.Body.Vel = boid.Body.Vel.Normalize().Mul(max)
	}
}

// SqDist returns the squared distance between two boids.
func (b *Boid) sqDist(o *Boid, t geom.Torus) float64 {
	return t.SqDist(b.Body.Box.Min, o.Body.Box.Min)
}

// AvoidVec returns a vector to direct a away from b.
//
// The vector is biased such that it is stronger as the
// objects get closer.
func avoidVec(a, b geom.Point, dist float64, w *world.World) geom.Point {
	sqrt := math.Sqrt(dist)
	diff := w.Pixels.Sub(a, b)
	diff.X = math.Copysign(sqrt-math.Abs(diff.X), diff.X)
	diff.Y = math.Copysign(sqrt-math.Abs(diff.Y), diff.Y)
	return diff
}

type grid struct {
	w, h   int
	pixels geom.Torus
	cellSz geom.Point
	cells  [][]*Boid
}

func newGrid(w, h int, t geom.Torus) *grid {
	return &grid{
		w:      w,
		h:      h,
		pixels: t,
		cellSz: geom.Pt(t.W/float64(w), t.H/float64(h)),
		cells:  make([][]*Boid, w*h),
	}
}

// At returns the the cell for x, y
func (g *grid) index(x, y int) int {
	return wrap(x, g.w)*g.h + wrap(y, g.h)
}

// pt2Cell returns the cell that contains the given point.
func (g *grid) pt2Cell(p geom.Point) (int, int) {
	return int(math.Floor(p.X/g.cellSz.X + 0.5)),
		int(math.Floor(p.Y/g.cellSz.Y + 0.5))
}

// Neighbors returns the neighbors of the given boid at a specified radius.
func (g *grid) neighbors(boid *Boid, r float64) []*Boid {
	p := g.pixels.Norm(boid.Body.Box.Min)
	xmin, ymin := g.pt2Cell(p.Sub(geom.Pt(r, r)))
	xmax, ymax := g.pt2Cell(p.Add(geom.Pt(r, r)))

	var n []*Boid

	for x := xmin; x <= xmax; x++ {
		for y := ymin; y <= ymax; y++ {
			cell := g.cells[g.index(x, y)]
			for _, b := range cell {
				if boid.sqDist(b, g.pixels) <= r*r {
					n = append(n, b)
				}
			}
		}
	}

	return n
}

// wrap returns the value of n wrapped around if it goes
// above bound-1 or below zero.
func wrap(n, bound int) int {
	if n >= 0 && n < bound {
		return n
	}

	if bound <= 0 {
		panic("Bad bound in wrap")
	}
	n %= bound
	if n < 0 {
		n = bound + n
		if n < 0 {
			panic("A value wrapped to a negative")
		}
	}
	return n
}
