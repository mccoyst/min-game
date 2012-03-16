package main

import (
	"minima/world"
	"math/rand"
	"time"
	"flag"
	"fmt"
)

const (
	// The mean of the distribution over the
	// number of Gaussian2ds.
	meanGauss = 200

	// Standard deviation of the distribution
	// over the number of Gaussian2ds.
	stdevGauss = 20

	// The rate parameter of the exponential
	// distribution of mountain heights.
	rateHeight = 1.1

	// Minimum and maximum Gaussian2d covariance of
	// random Gaussian2ds.
	covMin, covMax = -0.5, 0.5

	// Minimum and maximum standard deviation of
	// random Gaussian2ds as a factor of the
	// map width and height.
	sdevMin, sdevMax = 0.04, 0.07
)

var (
	width = flag.Int("w", 100, "World width")
	height = flag.Int("h", 100, "World height")
	seed = flag.Int64("seed", 0, "Random seed: 0 == use time")
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
	w := world.Make(*width, *height)
	for i := 0; i < w.W; i++ {
		for j := 0; j < w.H; j++ {
			l := w.At(i, j)
			l.Terrain = &world.Terrain['g']
		}
	}

	ch := make(chan *Gaussian2d)
	go genGaussian2ds(&w, ch)
	for g := range ch {
		grow(&w, g)
	}

	fmt.Println("Saving image")
	savePng(&w)
}

// grow generates a random height for the mean
// of the given Gaussian2d and grows the world
// around it.
func grow(w *world.World, g *Gaussian2d) {
	const s = 4.0	// standard deviations to grow around
	xmin, xmax := int(g.Mx - s*g.Sx), int(g.Mx + s*g.Sx)
	ymin, ymax := int(g.My - s*g.Sy), int(g.My + s*g.Sy)

	ht := rndHt()
	peek := g.PDF(g.Mx, g.My)
	mul := float64(ht) / peek

	for x := xmin; x < xmax; x++ {
		for y := ymin; y < ymax; y++ {
			l := w.At(x, y)
			p := g.PDF(float64(x), float64(y)) * mul
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

// rndHt returns a randomly choosen height for a
// mountain.  Heights are choosen from an exponential
// distribution.
func rndHt() int {
	ht := rand.ExpFloat64()/1.1 + 1

	if ht > world.MaxHeight {
		ht = world.MaxHeight
	}

	return int(ht)
}

// genGaussian2ds generates a bunch of random
// Gaussian2ds and sends them to the given channel.
func genGaussian2ds(w *world.World, ch chan<-*Gaussian2d) {
	num := int(rand.NormFloat64() * stdevGauss + meanGauss)
	fmt.Println("num gaussians: ", num)

	gs := make([]*Gaussian2d,0, num)
	for i := 0; i < num; i++ {
		tries := 0
retry:
		g := randomGaussian2d(w)
		for j := range gs {
			if tries < 1000 && near(g, gs[j]) {
				tries++
				goto retry
			}
		}

		gs = append(gs, g)
		ch <- g
	}
	close(ch)
}

// randomGaussian2d creates a random Gaussian2d
// somewhere in the world and returns it.
func randomGaussian2d(w *world.World) *Gaussian2d {
	mx := rand.Float64() * float64(w.W)
	my := rand.Float64() * float64(w.H)

	sxmin, sxmax := sdevMin*float64(w.W), sdevMax*float64(w.W)
	symin, symax := sdevMin*float64(w.H), sdevMax*float64(w.H)
	sx := rand.Float64() * (sxmax-sxmin) + sxmin
	sy := rand.Float64() * (symax-symin) + symin

	cov := rand.Float64() * (covMax - covMin) + covMin

	return NewGaussian2d(mx, my, sx, sy, cov)
}

// near returns true if the two Gassians are too close.
func near(a, b *Gaussian2d) bool {
	const s = 2	// standard deviations to grow around
	axmin, axmax := int(a.Mx - s*a.Sx), int(a.Mx + s*a.Sx)
	bxmin, bxmax := int(b.Mx - s*b.Sx), int(b.Mx + s*b.Sx)
	if axmin > bxmax || bxmin > axmax {
		return false
	}

	aymin, aymax := int(a.My - s*a.Sy), int(a.My + s*a.Sy)
	bymin, bymax := int(b.My - s*b.Sy), int(b.My + s*b.Sy)
	if aymin > bymax || bymin > aymax {
		return false
	}
	return true
}

