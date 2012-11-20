package main

import (
	"math"
	"math/rand"

	"code.google.com/p/min-game/ui"
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
	if p.body.Vel == ui.Pt(0, 0) || ptDist(cur.Body().Box.Center(), p.body.Box.Center(), worldSize(w)) > TileSize*3 {
		return
	}
	d := p.body.Box.Center().Sub(cur.Body().Box.Center())
	cur.Body().Vel = cur.Body().Vel.Sub(d.Div(5))
}

func (f *flock) moveAway(cur boid, w *world.World) {
	sz := worldSize(w)
	var dist ui.Point
	for _, b := range f.boids {
		
		if b == cur || ptDist(b.Body().Box.Center(), cur.Body().Box.Center(), sz) > f.avoidDist {
			continue
		}
		diff := cur.Body().Box.Center().Sub(b.Body().Box.Center())
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
	sz := worldSize(w)
	var avg ui.Point
	var n float64
	for _, b := range f.boids {
		
		if b == cur || ptDist(b.Body().Box.Center(), cur.Body().Box.Center(), sz) > f.localDist {
			continue
		}
		n++
		avg = avg.Add(cur.Body().Box.Center().Sub(b.Body().Box.Center()))
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
	sz := worldSize(w)
	var avg ui.Point
	var n float64
	for _, b := range f.boids {
		
		if b == cur || ptDist(b.Body().Box.Center(), cur.Body().Box.Center(), sz) > f.localDist {
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
func (f *flock) randVel() ui.Point {
	x := rand.Float64()*2 - 1
	y := rand.Float64()*2 - 1
	speed := rand.Float64() * f.maxSpeed
	return vecNorm(ui.Pt(x, y), speed)
}

// ClampVel returns v, clamped so that its magnitude is no more
// than a maximum value.
func clampVel(v ui.Point, max float64) ui.Point {
	if v.Len() > max {
		return vecNorm(v, max)
	} else if v.Len() < -max {
		return vecNorm(v, -max)
	}
	return v
}

// PtDist returns the distance of two points on a torus.
func ptDist(a, b, sz ui.Point) float64 {
	dx := dist(a.X, b.X, sz.X)
	dy := dist(a.Y, b.Y, sz.Y)
	return math.Sqrt(dx*dx + dy*dy)
}

// Dist returns the distance between two values, wrapped at width.
func dist(a, b, width float64) float64 {
	min, max := wrap(a, width), wrap(b, width)
	if min > max {
		min, max = max, min
	}
	return math.Min(max-min, min+width-max)
}
