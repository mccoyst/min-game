package main

import (
	"bufio"
	"flag"
	"fmt"
	"math/rand"
	"minima/world"
	"os"
	"time"
)

const (
	// The parameters of the distribution over the
	// number of Gaussian2ds.
	meanGauss, stdevGauss = 150, 10

	// The parameters of the normal distribution
	// over mountain growths.
	meanGrowth, stdevGrowth = 0, world.MaxHeight / 5.0

	// Minimum and maximum Gaussian2d covariance of
	// random Gaussian2ds.
	covMin, covMax = -0.5, 0.5

	// Minimum and maximum standard deviation of
	// random Gaussian2ds as a factor of the
	// map width and height.
	sdevMin, sdevMax = 0.04, 0.1
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
	fmt.Println("seed", *seed)
	rand.Seed(*seed)

	fmt.Println("width", *width)
	fmt.Println("height", *height)
	w := initWorld(*width, *height)

	num := int(rand.NormFloat64()*stdevGauss + meanGauss)
	fmt.Println("num gaussians: ", num)
	for g := range gaussians(w, num) {
		grow(w, g)
	}

	fmt.Println("Saving image")
	savePng(w)

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

	sxmin, sxmax := sdevMin*float64(w.W), sdevMax*float64(w.W)
	symin, symax := sdevMin*float64(w.H), sdevMax*float64(w.H)
	sx := rand.Float64()*(sxmax-sxmin) + sxmin
	sy := rand.Float64()*(symax-symin) + symin

	ht := rand.NormFloat64()*stdevGrowth + meanGrowth
	cov := rand.Float64()*(covMax-covMin) + covMin

	return NewGaussian2d(mx, my, sx, sy, ht, cov)
}
