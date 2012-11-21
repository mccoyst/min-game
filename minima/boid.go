package main

import (
	"math"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
)

// Boids is a list of boids that may flock together.
type Boids interface {
	Len() int
	Boid(int) Boid
	Info() BoidInfo
}

// A Boid is a bird-like (cow-like, or fish-like) object.
type Boid struct {
	*Body
}

// BoidInfo contains behavior information about boids.
type BoidInfo struct {
	// MaxVelocity is the fastest that a boid can move.
	MaxVelocity float64

	// LocalDist is the distance determining when two boids
	// are flocking together.
	LocalDist float64

	// AvoidDist is the distance at which boids will avoid
	// eachother.
	AvoidDist float64

	// PlayerDist is the distance at which boids will avoid
	// the player. 
	PlayerDist float64

	// These biases weight the three standard boid rules:
	//	Move toward the center of neighbors.
	//	Match the velocity of neighbors.
	//	Avoid very close neighbors.
	CenterBias, MatchBias, AvoidBias float64

	// PlayerBias is the weight applied to avoiding the player.
	PlayerBias float64
}

// UpdateBoids updates the velocity of the boids.
func UpdateBoids(boids Boids, p *Player, w *world.World) {
	info := boids.Info()
	local := localBoids(boids, w)
	for i, l := range local {
		boid := boids.Boid(i)
		boid.matchVel(l, info.MatchBias)
		boid.moveCenter(l, info.CenterBias, w)
		boid.avoidOthers(l, info.AvoidDist, info.AvoidBias, w)
		boid.avoidPlayer(p, info.PlayerDist, info.PlayerBias, w)
		boid.clampVel(info.MaxVelocity)
	}
}

// LocalBoids returns a slice containing the Boids that
// are local to the Boid with the corresponding index.
func localBoids(boids Boids, w *world.World) [][]Boid {
	localDist := boids.Info().LocalDist
	local := make([][]Boid, boids.Len())
	for i := range local {
		boid := boids.Boid(i)
		for j := i+1; j < boids.Len(); j++ {
			b := boids.Boid(j)
			if boid.dist(b, w) > localDist {
				continue
			}
			local[i] = append(local[i], b)
			local[j] = append(local[j], boid)
		}
	}
	return local
}

// MatchVel attempts to match the velocity of the local boids.
func (boid Boid) matchVel(local []Boid, bias float64) {
	var avg geom.Point
	for _, b := range local {
		avg = avg.Add(b.Vel)
	}
	if len(local) == 0 {
		return
	}
	avg = avg.Div(float64(len(local))).Normalize().Mul(bias)
	boid.Vel = boid.Vel.Add(avg)
}

// MoveCenter attempts to move the boid toward the center
// of its local flock mates.
func (boid Boid) moveCenter(local []Boid, bias float64, w *world.World) {
	var avg geom.Point
	for _, b := range local {
		toCenter := w.Pixels.Sub(b.Center(), boid.Center())
		avg = avg.Add(toCenter)
	}
	if len(local) == 0 {
		return
	}
	avg = avg.Div(float64(len(local))).Normalize().Mul(bias)
	boid.Vel = boid.Vel.Add(avg)
}

// AvoidOthers attempts to avoid very close flock mates.
func (boid Boid) avoidOthers(local []Boid, dist, bias float64, w *world.World) {
	var a geom.Point
	for _, b := range local {
		if d := boid.dist(b, w); d > dist {
			continue
		}
		a = a.Add(avoidVec(boid.Center(), b.Center(), dist, w))
	}
	a = a.Mul(bias)
	boid.Vel = boid.Vel.Add(a)
}

// AvoidPlayer attempts to avoid the player.
func (boid Boid) avoidPlayer(p *Player, dist, bias float64, w *world.World) {
	c := p.body.Center()
	if p.body.Vel == geom.Pt(0, 0) || w.Pixels.Dist(boid.Center(), c) > dist {
		return
	}
	d := avoidVec(boid.Center(), c, dist, w).Mul(bias)
	boid.Vel = boid.Vel.Add(d)
}

// ClampVel clamps the boid's velocity to have a magnitude of
// no more than max.
func (boid Boid) clampVel(max float64) {
	if boid.Vel.Len() > max {
		boid.Vel = boid.Vel.Normalize().Mul(max)
	}
}

// Dist returns the distance between two boids.
func (b Boid) dist(o Boid, w *world.World) float64 {
	return w.Pixels.Dist(b.Center(), o.Center())
}

// AvoidVec returns a vector to direct a away from b.
//
// The vector is biased such that it is stronger as the
// objects get closer.
func avoidVec(a, b geom.Point, dist float64, w *world.World) geom.Point {
	sqrt := math.Sqrt(dist)
	diff := w.Pixels.Sub(a, b)
	diff.X = math.Copysign(sqrt - math.Abs(diff.X), diff.X)
	diff.Y = math.Copysign(sqrt - math.Abs(diff.Y), diff.Y)
	return diff
}