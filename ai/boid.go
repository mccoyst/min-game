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
	Boid(int) Boid
	BoidInfo() BoidInfo
}

// A Boid is a bird-like (cow-like, or fish-like) object.
type Boid struct {
	*phys.Body
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
}

// UpdateBoids updates the velocity of the boids.
func UpdateBoids(boids Boids, p *phys.Body, w *world.World) {
	info := boids.BoidInfo()
	local := localBoids(boids, w)
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
func localBoids(boids Boids, w *world.World) [][]Boid {
	localDist := boids.BoidInfo().LocalDist
	dd := localDist * localDist
	local := make([][]Boid, boids.Len())
	for i := range local {
		boid := boids.Boid(i)
		for j := i + 1; j < boids.Len(); j++ {
			b := boids.Boid(j)
			if boid.sqDist(b, w) > dd {
				continue
			}
			local[i] = append(local[i], b)
			local[j] = append(local[j], boid)
		}
	}
	return local
}

// MatchVel attempts to match the velocity of the local boids.
func (boid Boid) matchVel(local []Boid, i BoidInfo) {
	var avg geom.Point
	for _, b := range local {
		avg = avg.Add(b.Vel)
	}
	if len(local) == 0 {
		return
	}
	avg = avg.Div(float64(len(local))).Normalize().Mul(i.MatchBias)
	boid.Vel = boid.Vel.Add(avg)
}

// MoveCenter attempts to move the boid toward the center
// of its local flock mates.
func (boid Boid) moveCenter(local []Boid, i BoidInfo, w *world.World) {
	var avg, c geom.Point
	for _, b := range local {
		toCenter := w.Pixels.Sub(b.Box.Min, boid.Box.Min)
		c = c.Add(b.Box.Min)
		avg = avg.Add(toCenter)
	}
	if len(local) == 0 {
		return
	}
	n := float64(len(local))
	c = c.Div(n)
	if w.Pixels.SqDist(c, boid.Box.Min) < i.CenterDist*i.CenterDist {
		return
	}
	avg = avg.Div(n).Normalize().Mul(i.CenterBias)
	boid.Vel = boid.Vel.Add(avg)
}

// AvoidOthers attempts to avoid very close flock mates.
func (boid Boid) avoidOthers(local []Boid, i BoidInfo, w *world.World) {
	dd := i.AvoidDist * i.AvoidDist
	var a geom.Point
	for _, b := range local {
		if d := boid.sqDist(b, w); d > dd {
			continue
		}
		a = a.Add(avoidVec(boid.Center(), b.Center(), i.AvoidDist, w))
	}
	a = a.Mul(i.AvoidBias)
	boid.Vel = boid.Vel.Add(a)
}

// AvoidPlayer attempts to avoid the player.
func (boid Boid) avoidPlayer(p *phys.Body, i BoidInfo, w *world.World) {
	dd := i.PlayerDist * i.PlayerDist
	pt := p.Box.Min
	if p.Vel == geom.Pt(0, 0) || w.Pixels.SqDist(boid.Box.Min, pt) > dd {
		return
	}
	d := avoidVec(boid.Box.Min, pt, i.PlayerDist, w).Mul(i.PlayerBias)
	boid.Vel = boid.Vel.Add(d)
}

// AvoidTerrain attempts to avoid certain types of terrain.
func (boid Boid) avoidTerrain(i BoidInfo, w *world.World) {
	if i.AvoidTerrain == "" {
		return
	}

	var a geom.Point
	dd := i.TerrainDist * i.TerrainDist
	sz := geom.Pt(i.TerrainDist, i.TerrainDist)
	x0, y0 := w.Tile(boid.Box.Min.Sub(sz))
	x1, y1 := w.Tile(boid.Box.Min.Add(sz))

	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			r := w.At(x, y).Terrain.Char
			if strings.Index(i.AvoidTerrain, r) < 0 {
				continue
			}
			pt := geom.Pt(float64(x)*world.TileSize.X,
				float64(y)*world.TileSize.Y)
			if w.Pixels.SqDist(boid.Box.Min, pt) > dd {
				continue
			}
			a = a.Add(avoidVec(boid.Box.Min, pt, i.TerrainDist, w))
		}
	}
	a = a.Mul(i.TerrainBias)
	boid.Vel = boid.Vel.Add(a)
}

// ClampVel clamps the boid's velocity to have a magnitude of
// no more than max.
func (boid Boid) clampVel(max float64) {
	if boid.Vel.Len() > max {
		boid.Vel = boid.Vel.Normalize().Mul(max)
	}
}

// SqDist returns the squared distance between two boids.
func (b Boid) sqDist(o Boid, w *world.World) float64 {
	return w.Pixels.SqDist(b.Box.Min, o.Box.Min)
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
