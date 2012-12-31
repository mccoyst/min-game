// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

/* 
While the benchmarks in this package should prove useful, it should
benoted that they don't truely match the case that we see in the real
game.Here, we are distributing boids uniformly throughout the world,
but in thegame boids are really distributed in clusters. Some spacial
data structures,for example, will probably show better performance on
these benchmarks than will be seen in the actual game.
*/

package ai

import (
	"math/rand"
	"reflect"
	"testing"

	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/phys"
	"code.google.com/p/min-game/world"
)

func BenchmarkUpdateBoids100(b *testing.B) {
	updateN(100, b)
}

func BenchmarkUpdateBoids500(b *testing.B) {
	updateN(500, b)
}

func BenchmarkUpdateBoids1000(b *testing.B) {
	updateN(1000, b)
}

// UpdateN creates a random set of n boids within the benchmark world
// and benchmarks UpdateBoids by calling b.N times on the random boids.
func updateN(n int, b *testing.B) {
	b.StopTimer()
	w := benchWorld()
	p := benchPlayer()

	var bds boids
	boidVal := bds.Generate(rand.New(rand.NewSource(rand.Int63())), n)
	bds = boidVal.Interface().(boids)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		UpdateBoids(uint(i), bds, p, w)
	}
}

const (
	// WorldWidth and WorldHeight define the size of the world into which
	// randomly generated boids are placed.
	worldWidth = 500
	worldHeight
)

// BenchWorld returns a world of the given dimensions with all locations set
// to grassland.
func benchWorld() *world.World {
	w := world.New(worldWidth, worldHeight)
	for x := 0; x < worldWidth; x++ {
		for y := 0; y < worldHeight; y++ {
			w.At(x, y).Terrain = &world.Terrain[int('g')]
		}
	}
	return w
}

// BenchPlayer returns an random player location in the benchmark world.
func benchPlayer() *phys.Body {
	x := rand.Float64() * worldWidth * world.TileSize.X
	y := rand.Float64() * worldHeight * world.TileSize.Y
	return &phys.Body{
		Vel: geom.Point{rand.Float64(), rand.Float64()},
		Box: geom.Rect(x, y, x+boidWidth, y+boidHeight),
	}
}

type boids []Boid

var (
	// BoidWidth and BoidHeight define the size of randomly generated boids.
	boidWidth  = world.TileSize.X
	boidHeight = world.TileSize.Y
)

// Generate returns the given number of random boids, implementing the
// quick.Generator interface.
func (boids) Generate(r *rand.Rand, size int) reflect.Value {
	bs := make([]Boid, size)
	for i := range bs {
		x := r.Float64() * worldWidth * world.TileSize.X
		y := r.Float64() * worldHeight * world.TileSize.Y
		bs[i].Body = &phys.Body{
			Vel: geom.Point{r.Float64(), r.Float64()},
			Box: geom.Rect(x, y, x+boidWidth, y+boidHeight),
		}
	}
	return reflect.ValueOf(boids(bs))
}

func (b boids) Len() int {
	return len(b)
}

func (b boids) Boid(i int) Boid {
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
