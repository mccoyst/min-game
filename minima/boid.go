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

	// Goal is some random goal velocity.
	goal ui.Point

	// nGoal is the number of frames to add the goal velocity
	// before using a new one.
	nGoal int

	boids []boid
}

type boid interface {
	Body() *Body
	Move(*world.World)
	Draw(Drawer, Camera) error
}

func (f *flock) Move(w *world.World) {
	f.update(w)
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
func (f *flock) update(w *world.World) {
	for _, b := range f.boids {
		local := f.localBoids(w, b)
		pt := b.Body().Box.Center()

		avoid := ui.Pt(0, 0)
		vel := ui.Pt(0, 0)
		center := ui.Pt(0, 0)
		for _, l := range local {
			body := l.Body()
			vel = vel.Add(body.Vel)
			center = center.Add(body.Box.Center())
			d := ptDist(pt, l.Body().Box.Center(), worldSize(w))
			if d <= f.avoidDist {
				avoid = avoid.Sub(l.Body().Box.Center().Sub(pt))
			}
		}

		n := float64(len(local))
		v := b.Body().Vel
		avoid = avoid.Mul(2)
		vel = v.Div(n).Sub(v).Div(8)
		center = center.Div(n).Sub(pt).Div(100)
		v = v.Add(avoid).Add(center).Add(vel)
		v = v.Add(f.goal.Div(2))
		b.Body().Vel = clampVel(v, f.maxSpeed)
	}

	f.nGoal--
	if f.nGoal <= 0 {
		f.goal = f.randVel()
		f.nGoal = 100
	}
}

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

func (f *flock) localBoids(w *world.World, bd boid) []boid {
	sz := worldSize(w)
	pt := bd.Body().Box.Center()

	var local []boid
	var nearest boid
	nearDist := math.Inf(1)
	for _, b := range f.boids {
		if bd == b {
			continue
		}
		d := ptDist(pt, b.Body().Box.Center(), sz)
		if d < nearDist {
			nearest = b
			nearDist = d
		}
		if d <= f.localDist {
			local = append(local, b)
		}
	}
	if len(local) == 0 && nearest != nil {
		return []boid{nearest}
	}
	return local
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
