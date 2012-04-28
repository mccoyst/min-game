package main

import (
	"minima/world"
	"runtime/pprof"
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	// gaussFact is the number of Gaussians given as
	// a factor of the map size.
	gaussFact = 0.002

	// meanGroth and stdevGrowth are the parameters
	// of the normal distribution over mountain growths.
	meanGrowth, stdevGrowth = 0, world.MaxElevation * 0.15

	// conMin and maxCov are the minimum and maximum
	// Gaussian2d covariance of random Gaussian2ds.
	minCov, maxCov = -0.5, 0.5

	// stdevMin and stdevMax are the minimum and
	// maximum standard deviation of random Gaussians.
	minStdev, maxStdev = 10, 30
)

var (
	width      = flag.Int("w", 500, "World width")
	height     = flag.Int("h", 500, "World height")
	seed       = flag.Int64("seed", 0, "Random seed: 0 == use time")
	cpuprofile = flag.String("cprof", "", "Write cpu profile to file")
	memprofile = flag.String("mprof", "", "Write mem profile to file")
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			panic(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *seed == 0 {
		*seed = int64(time.Now().Nanosecond())
	}
	rand.Seed(*seed)
	fmt.Fprintln(os.Stderr, "seed", *seed)
	w := initWorld(*width, *height)

	num := int(float64(w.W*w.H) * gaussFact)
	for g := range gaussians(w, num) {
		grow(w, g)
	}
	clampHeights(w)
	doTerrain(w)
	placeStart(w)

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if err := w.Write(out); err != nil {
		panic(err)
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			panic(err)
		}
		pprof.WriteHeapProfile(f)
	}
}

// initWorld returns a newly initialized world
// with the given dimensions.
func initWorld(width, height int) *world.World {
	w := world.Make(width, height)
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			l := w.AtCoord(x, y)
			l.Elevation = world.MaxElevation/4 + 1
		}
	}
	return &w
}

// clampHeights ensures that all locations have a
// height that is within the allowable range.
func clampHeights(w *world.World) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			l := w.At(x, y)
			if l.Elevation < 0 {
				l.Elevation = 0
			}
			if l.Elevation > world.MaxElevation {
				l.Elevation = world.MaxElevation
			}
		}
	}

}

// grow generates a random height for the mean
// of the given Gaussian2d and grows the world
// around it.
func grow(w *world.World, g *Gaussian2d) {
	const s = 2.0 // standard deviations to grow around
	xmin, xmax := int(g.Mx-s*g.Sx), int(g.Mx+s*g.Sx)
	ymin, ymax := int(g.My-s*g.Sy), int(g.My+s*g.Sy)

	for x := xmin; x < xmax; x++ {
		for y := ymin; y < ymax; y++ {
			l := w.AtCoord(x, y)
			p := g.PDF(float64(x), float64(y))
			l.Elevation = l.Elevation + int(p)
		}
	}
}

// gaussians returns a channel upon which the
// given number of random Gaussians will be
// sent.
func gaussians(w *world.World, num int) <-chan *Gaussian2d {
	ch := make(chan *Gaussian2d, 100)
	go func() {
		for i := 0; i < num; i++ {
			ch <- randomGaussian2d(w)
		}
		close(ch)
	}()
	return ch
}

// randomGaussian2d returns a random Gaussian2d that
// is generated using the maxStdev, minStdev, stdevGrowth,
// meanGrowth, maxCov, and minCov constants.
func randomGaussian2d(w *world.World) *Gaussian2d {
	mx := rand.Float64() * float64(w.W)
	my := rand.Float64() * float64(w.H)

	sx := rand.Float64()*(maxStdev-minStdev) + minStdev
	sy := rand.Float64()*(maxStdev-minStdev) + minStdev

	ht := 0.0
	for int(ht) == 0 {
		ht = rand.NormFloat64()*stdevGrowth + meanGrowth
	}
	cov := rand.Float64()*(maxCov-minCov) + minCov

	return NewGaussian2d(mx, my, sx, sy, ht, cov)
}

// placeStart places the start location on a random grass tile.
func placeStart(w *world.World) {
	var grass []int
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			loc := w.At(x, y)
			if loc.Terrain == &world.Terrain[int('g')] {
				grass = append(grass, x*w.H + y)
			}
		}
	}

	if len(grass) == 0 {
		return
	}

	ind := rand.Intn(len(grass))
	w.X0 = grass[ind] / w.H
	w.Y0 = grass[ind] % w.H

	if w.At(w.X0, w.Y0).Terrain != &world.Terrain[int('g')] {
		panic("Start location is not grass")
	}
}
