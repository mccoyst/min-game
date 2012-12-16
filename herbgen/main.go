// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	gomath "math"
	"math/rand"
	"os"
	"sort"
	"time"

	"code.google.com/p/min-game/animal"
	"code.google.com/p/min-game/geom"
	"code.google.com/p/min-game/math"
	"code.google.com/p/min-game/world"
)

var (
	name   = flag.String("name", "Gull", "Name of the herbivores to generate")
	num    = flag.Int("num", 25, "Number to generate")
	stdev  = flag.Float64("stdev", 20, "Standard deviation of high-probability regions")
	ngauss = flag.Int("ngauss", 5, "Number high-probability regions")
	seed   = flag.Int64("seed", time.Now().UnixNano(), "The random seed")
	draw   = flag.String("draw", "", "Draw the probability distribution to an image")
)

func main() {
	flag.Parse()
	rand.Seed(*seed)

	in := bufio.NewReader(os.Stdin)
	var err error
	w, err := world.Read(in)
	if err != nil {
		panic(err)
	}

	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	if err := w.Write(out); err != nil {
		panic(err)
	}

	herbs := placeHerbs(w)

	writeGame(in, out, herbs)
}

// PlaceHerbs places herbivores in the world.
func placeHerbs(w *world.World) animal.Herbivores {
	herbs, err := animal.MakeHerbivores(*name)
	if err != nil {
		panic(err)
	}

	ls := locs(w, herbs)
	ps := probs(w, ls)

	if *draw != "" {
		drawProbs(w, ls, ps)
	}

	for i := 1; i < len(ps); i++ {
		ps[i] += ps[i-1]
	}
	if gomath.Abs(ps[len(ps)-1]-1.0) > 0.0001 {
		panic(fmt.Sprintf("Probs don't sum to 1, they sum to %f", ps[len(ps)-1]))
	}
	ps[len(ps)-1] = 1 // Get rid of possible rounding issues.

	for n := 0; n < *num && len(ls) > 0; n++ {
		p := rand.Float64()
		i := sort.SearchFloat64s(ps, p)
		if i != 0 && (i >= len(ps) || ps[i] > p) {
			i--
		}
		vel := geom.Pt(rand.Float64(), rand.Float64()).Normalize()
		herbs.Spawn(ls[i].Point(), vel)
	}

	return herbs
}

// writeGame reads the game state, adds the herbivores, and writes it out.
func writeGame(in *bufio.Reader, out *bufio.Writer, herbs animal.Herbivores) {
	game := make(map[string]interface{})
	err := json.NewDecoder(in).Decode(&game)
	if err != nil && err != io.EOF {
		panic("Error reading JSON for " + *name + " " + err.Error())
	}
	if hs, ok := game["Herbivores"]; ok {
		game["Herbivores"] = append(hs.([]interface{}), herbs)
	} else {
		game["Herbivores"] = []interface{}{herbs}
	}

	b, err := json.MarshalIndent(game, "", "\t")
	if err != nil {
		panic(err)
	}
	if _, err := out.Write(b); err != nil {
		panic(err)
	}
}

// Locs returns the valid locations to place this herbivore type.
func locs(w *world.World, herbs animal.Herbivores) []*world.Loc {
	var typs []string
	maxAffinity := 0.0
	for t, a := range herbs.Info.Affinity {
		if a > maxAffinity {
			typs = []string{t}
			maxAffinity = a
		} else if a == maxAffinity {
			typs = append(typs, t)
		}
	}

	var locs []*world.Loc
	for _, t := range typs {
		ls := w.LocsWithType(t)
		for _, l := range ls {
			if l.Depth <= herbs.Info.BoidInfo.MaxDepth {
				locs = append(locs, l)
			}
		}
	}

	return locs
}

// Probs returns the probability corresponding to each location.
func probs(w *world.World, locs []*world.Loc) []float64 {
	wprobs := make([]float64, w.W*w.H)

	gauss := make([]*math.Gaussian2d, *ngauss)
	for i := range gauss {
		gauss[i] = randomGaussian2d(w)
	}

	const s = 2.0 // σ to compute prob around each gauss.
	for _, g := range gauss {
		xmin, xmax := int(g.Mx-s*g.Sx), int(g.Mx+s*g.Sx)
		ymin, ymax := int(g.My-s*g.Sy), int(g.My+s*g.Sy)
		for x := xmin; x < xmax; x++ {
			for y := ymin; y < ymax; y++ {
				i, j := w.Wrap(x, y)
				wprobs[i*w.H+j] += g.PDF(float64(x), float64(y))
			}
		}
	}

	sum := 0.0
	probs := make([]float64, len(locs))
	for i, l := range locs {
		probs[i] = wprobs[l.X*w.H+l.Y]
		sum += probs[i]
	}
	for i := range probs {
		probs[i] /= sum
	}
	return probs
}

// randomGaussian2d returns a random Gaussian2d.
func randomGaussian2d(w *world.World) *math.Gaussian2d {
	mx := rand.Float64() * float64(w.W)
	my := rand.Float64() * float64(w.H)
	ht := 1.0
	cov := 0.0
	return math.NewGaussian2d(mx, my, *stdev, *stdev, ht, cov)
}

// DrawProbs draws the world, with cells shaded lighter
// if they have a greater probability of containing an animal.
func drawProbs(w *world.World, locs []*world.Loc, probs []float64) {
	out, err := os.Create(*draw)
	if err != nil {
		panic(err)
	}
	defer out.Close()

	ps := make([]float64, w.W*w.H)
	mx := 0.0
	for i, l := range locs {
		ps[l.X*w.H+l.Y] = probs[i]
		if probs[i] > mx {
			mx = probs[i]
		}
	}
	png.Encode(out, &worldImg{w, ps, mx})
}

type worldImg struct {
	*world.World
	probs []float64
	mx    float64
}

// Bounds implements the Bounds() method of
// the image.Image interface.
func (w *worldImg) Bounds() image.Rectangle {
	return image.Rect(0, 0, w.W, w.H)
}

var (
	colors = []color.RGBA{
		'g': color.RGBA{109, 170, 44, 255},
		'm': color.RGBA{210, 125, 44, 255},
		'w': color.RGBA{109, 194, 202, 255},
		'l': color.RGBA{208, 70, 72, 255},
		'd': color.RGBA{218, 219, 94, 255},
		'f': color.RGBA{52, 101, 36, 255},
		'i': color.RGBA{222, 238, 214, 255},
	}
)

// At implements the At() method of the image.Image interface.
func (w *worldImg) At(x, y int) color.Color {
	p := w.probs[x*w.W+y]
	loc := w.World.At(x, y)
	min, max := 0.1, 1.0
	f := (p/w.mx)*(max-min) + min
	c := colors[loc.Terrain.Char[0]]
	return color.RGBA{
		R: uint8(float64(c.R) * f),
		G: uint8(float64(c.G) * f),
		B: uint8(float64(c.B) * f),
		A: c.A,
	}
}

// ColorModel implements the ColorModel() method
// of the image.Image interface.
func (w *worldImg) ColorModel() color.Model {
	return color.RGBAModel
}
