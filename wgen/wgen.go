package main

import (
	"bufio"
	"flag"
	"math/rand"
	"minima/world"
	"os"
	"time"
)

const (
	// gaussFact is the number of Gaussians given as
	// a factor of the map size.
	gaussFact = 0.02

	// meanGroth and stdevGrowth are the parameters
	// of the normal distribution over mountain growths.
	meanGrowth, stdevGrowth = 0, world.MaxHeight / 6.0

	// conMin and covMax are the minimum and maximum
	// Gaussian2d covariance of random Gaussian2ds.
	covMin, covMax = -0.5, 0.5

	// stdevMin and stdevMax are the minimum and
	// maximum standard deviation of random Gaussians.
	sdevMin, sdevMax = 4, 10
)

var (
	width  = flag.Int("w", 100, "World width")
	height = flag.Int("h", 100, "World height")
	seed   = flag.Int64("seed", 0, "Random seed: 0 == use time")
)

func main() {
	flag.Parse()

	if *seed == 0 {
		*seed = int64(time.Now().Nanosecond())
	}
	rand.Seed(*seed)
	w := initWorld(*width, *height)

	num := int(float64(w.W*w.H)*gaussFact)
	for g := range gaussians(w, num) {
		grow(w, g)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if err := w.Write(out); err != nil {
		panic(err)
	}
}

// initWorld returns a newly initialized world
// with the given dimensions.
func initWorld(width, height int) *world.World {
	w := world.Make(width, height)
	for i := 0; i < w.W; i++ {
		for j := 0; j < w.H; j++ {
			l := w.AtCoord(i, j)
			l.Terrain = &world.Terrain['g']
			l.Height = world.MaxHeight / 2
		}
	}
	return &w
}

// grow generates a random height for the mean
// of the given Gaussian2d and grows the world
// around it.
func grow(w *world.World, g *Gaussian2d) {
	const s = 4.0 // standard deviations to grow around
	xmin, xmax := int(g.Mx-s*g.Sx), int(g.Mx+s*g.Sx)
	ymin, ymax := int(g.My-s*g.Sy), int(g.My+s*g.Sy)

	for x := xmin; x < xmax; x++ {
		for y := ymin; y < ymax; y++ {
			l := w.AtCoord(x, y)
			p := g.PDF(float64(x), float64(y))
			l.Height += int(p)
			if l.Height < 0 {
				l.Height = 0
			}
			if l.Height > world.MaxHeight {
				l.Height = world.MaxHeight
			}
		}
	}
}

// gaussians makes some random Gaussians and
// sends them out on the given channel.
func gaussians(w *world.World, num int) <-chan *Gaussian2d {
	ch := make(chan *Gaussian2d)
	go func() {
		for i := 0; i < num; i++ {
			ch <- randomGaussian2d(w)
		}
		close(ch)
	}()
	return ch
}

// randomGaussian2d creates a random Gaussian2d
// somewhere in the world and returns it.
func randomGaussian2d(w *world.World) *Gaussian2d {
	mx := rand.Float64() * float64(w.W)
	my := rand.Float64() * float64(w.H)

	sx := rand.Float64()*(sdevMax-sdevMin) + sdevMin
	sy := rand.Float64()*(sdevMax-sdevMin) + sdevMin

	ht := rand.NormFloat64()*stdevGrowth + meanGrowth
	cov := rand.Float64()*(covMax-covMin) + covMin

	return NewGaussian2d(mx, my, sx, sy, ht, cov)
}
