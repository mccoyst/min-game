package main

import (
	"bufio"
	"flag"
	"math/rand"
	"minima/world"
	"os"
	"time"
	"math"
)

const (
	// gaussFact is the number of Gaussians given as
	// a factor of the map size.
	gaussFact = 0.002

	// meanGroth and stdevGrowth are the parameters
	// of the normal distribution over mountain growths.
	meanGrowth, stdevGrowth = 0, 10

	// conMin and covMax are the minimum and maximum
	// Gaussian2d covariance of random Gaussian2ds.
	covMin, covMax = -0.5, 0.5

	// stdevMin and stdevMax are the minimum and
	// maximum standard deviation of random Gaussians.
	sdevMin, sdevMax = 10, 30
)

var(
	// mtnMin is the minimum value above which the
	// terrain is mountains.
	mntMin = int(math.Floor(world.MaxHeight*0.90))

	// waterMax is the maximum value below which
	// the terrain is water.
	waterMax = int(math.Floor(world.MaxHeight*0.45))
)

var (
	width  = flag.Int("w", 500, "World width")
	height = flag.Int("h", 500, "World height")
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
	doTerrain(w)

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
			l.Height = (world.MaxHeight / 2) +1
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
			l.Height = l.Height + int(p)
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

	ht := 0.0
	for int(ht) == 0 {
		ht = rand.NormFloat64()*stdevGrowth + meanGrowth
	}
	cov := rand.Float64()*(covMax-covMin) + covMin

	return NewGaussian2d(mx, my, sx, sy, ht, cov)
}

// doTerrain clamps the heights of each cell and
// assigns their terrains.
func doTerrain(w *world.World) {
	for x := 0; x < w.W; x++ {
		for y := 0; y < w.H; y++ {
			l := w.At(x, y)
			if l.Height < 0 {
				l.Height = 0
			}
			if l.Height > world.MaxHeight {
				l.Height = world.MaxHeight
			}
			switch {
			case l.Height >= mntMin:
				l.Terrain = &world.Terrain['m']
			case l.Height <=waterMax:
				l.Terrain = &world.Terrain['w']
			}
		}
	}
}