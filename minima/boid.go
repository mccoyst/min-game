package main

import (
	"math"
	"math/rand"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/world"
)

type flock struct {
	// LocalDist is the distance defining the local neighborhood of each boid.
	localDist float64

	// AvoidDist is the distance at which boids begin to avoid eachother.
	avoidDist float64

	// MaxSpeed is the fastest than any boid can move.
	maxSpeed float64

	boids []boid
}

type boid interface {
	Body() *Body
	Move(*world.World)
	Draw(Drawer, Camera) error
}

func (f *flock) Move(p *Player, w *world.World) {
	f.update(p, w)
	for _, b := range f.boids {
		b.Move(w)
	}
}

func (f *flock) Draw(d Drawer, cam Camera) error {
	for _, b := range f.boids {
		if err := b.Draw(d, cam); err != nil {
			return err
		}
	}
	return nil
}

// Update updates the velocity of the boids in the flock.
func (f *flock) update(p *Player, w *world.World) {

	for _, b := range f.boids {
		f.moveWith(b, w)
		f.moveCloser(b, w)
		f.moveAway(b, w)
		f.avoidPlayer(b, p, w)
		b.Body().Vel = clampVel(b.Body().Vel, f.maxSpeed)
	}
}

func (f *flock) avoidPlayer(cur boid, p *Player, w *world.World) {
	if p.body.Vel == geom.Pt(0, 0) || w.Pixels.Dist(cur.Body().Box.Center(), p.body.Box.Center()) > TileSize*3 {
		return
	}
	d := w.Pixels.Sub(p.body.Box.Center(), cur.Body().Box.Center())
	cur.Body().Vel = cur.Body().Vel.Sub(d.Div(5))
}

func (f *flock) moveAway(cur boid, w *world.World) {
	var dist geom.Point
	for _, b := range f.boids {
		
		if b == cur || w.Pixels.Dist(b.Body().Box.Center(), cur.Body().Box.Center()) > f.avoidDist {
			continue
		}
		diff := w.Pixels.Sub(cur.Body().Box.Center(), b.Body().Box.Center())
		sqrt := math.Sqrt(f.avoidDist)
		if diff.X >= 0 {
			diff.X = sqrt - diff.X
		} else {
			diff.X = -sqrt - diff.X
		}
		if diff.Y >= 0 {
			diff.Y = sqrt - diff.Y
		} else {
			diff.Y = -sqrt - diff.Y
		}
		dist = dist.Add(diff)
	}

	cur.Body().Vel = cur.Body().Vel.Sub(dist.Div(2))
}

func (f *flock) moveCloser(cur boid, w *world.World) {
	var avg geom.Point
	var n float64
	for _, b := range f.boids {
		
		if b == cur || w.Pixels.Dist(b.Body().Box.Center(), cur.Body().Box.Center()) > f.localDist {
			continue
		}
		n++
		avg = avg.Add(w.Pixels.Sub(cur.Body().Box.Center(), b.Body().Box.Center()))
	}
	if n == 0 {
		return
	}
	avg = avg.Div(n)
	avg = vecNorm(avg, 0.05)
	cur.Body().Vel = cur.Body().Vel.Sub(avg)
	cur.Body().Vel = clampVel(cur.Body().Vel, f.maxSpeed)
}

func (f *flock) moveWith(cur boid, w *world.World) {
	var avg geom.Point
	var n float64
	for _, b := range f.boids {
		
		if b == cur || w.Pixels.Dist(b.Body().Box.Center(), cur.Body().Box.Center()) > f.localDist {
			continue
		}
		n++
		avg = avg.Add(b.Body().Vel)
	}
	if n == 0 {
		return
	}
	avg = avg.Div(n)
	avg = vecNorm(avg, 0.08)
	cur.Body().Vel = cur.Body().Vel.Add(avg)
	cur.Body().Vel = clampVel(cur.Body().Vel, f.maxSpeed)
}

// RandVel returns a random velocity within the speed limit
// of the flock.
func (f *flock) randVel() geom.Point {
	x := rand.Float64()*2 - 1
	y := rand.Float64()*2 - 1
	speed := rand.Float64() * f.maxSpeed
	return vecNorm(geom.Pt(x, y), speed)
}

// ClampVel returns v, clamped so that its magnitude is no more
// than a maximum value.
func clampVel(v geom.Point, max float64) geom.Point {
	if v.Len() > max {
		return vecNorm(v, max)
	} else if v.Len() < -max {
		return vecNorm(v, -max)
	}
	return v
}
