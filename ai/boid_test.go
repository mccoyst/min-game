// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ai

import (
	"math"
	"math/rand"
	"reflect"
	"testing"
	"testing/quick"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/world"
)

func TestGrid_Neighbors(t *testing.T) {
	err := quick.Check(func(boid *Boid, bs []*Boid, r float64) bool {
		r = math.Abs(r) * 320
		px := geom.Torus{500 * 32, 500 * 32}
		g := newGrid(10, 10, px)
		for _, b := range bs {
			c := g.index(g.pt2Cell(b.Body.Box.Min))
			g.cells[c] = append(g.cells[c], b)
		}

		good := make(map[*Boid]bool)
		for _, b := range near(boid, bs, px, r) {
			good[b] = true
		}

		n := make(map[*Boid]bool)
		for _, b := range g.neighbors(boid, r) {
			n[b] = true
		}

		for b := range good {
			if !n[b] {
				t.Errorf("Expected %v in range of %v", *b, *boid)
				return false
			}
		}

		for b := range n {
			if !good[b] {
				t.Errorf("Unxpected %v in range of %v", *b, *boid)
				return false
			}
		}
		return true
	}, nil)
	if err != nil {
		t.Error(err)
	}
}

func near(boid *Boid, bs []*Boid, px geom.Torus, r float64) []*Boid {
	var n []*Boid
	for _, b := range bs {
		if boid.sqDist(b, px) < r*r {
			n = append(n, b)
		}
	}
	return n
}

type boids []*Boid

var (
	// BoidWidth and BoidHeight define the size of randomly generated boids.
	boidWidth  = world.TileSize.X
	boidHeight = world.TileSize.Y
)

// Generate returns the given number of random boids, implementing the
// quick.Generator interface.
func (boids) Generate(r *rand.Rand, size int) reflect.Value {
	bs := make([]*Boid, size)
	for i := range bs {
		x := r.Float64() * worldWidth * world.TileSize.X
		y := r.Float64() * worldHeight * world.TileSize.Y
		bs[i] = &Boid{
			Body: phys.Body{
				Vel: geom.Point{r.Float64(), r.Float64()},
				Box: geom.Rect(x, y, x+boidWidth, y+boidHeight),
			},
		}
	}
	return reflect.ValueOf(boids(bs))
}

func (b boids) Len() int {
	return len(b)
}

func (b boids) Boid(i int) *Boid {
	return b[i]
}

func (boids) BoidInfo() BoidInfo {
	return BoidInfo{
		MaxVelocity:  0.5,
		LocalDist:    960,
		MatchBias:    0.0,
		CenterDist:   480,
		CenterBias:   0.005,
		AvoidDist:    48,
		AvoidBias:    0.001,
		PlayerDist:   64,
		PlayerBias:   0.02,
		TerrainDist:  35.2,
		TerrainBias:  0.0005,
		AvoidTerrain: "",
	}
}
