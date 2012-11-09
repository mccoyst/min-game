// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package world

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"image/color"
	"strings"
	"os"
	"runtime"
)

const (
	// MaxElevation is the number of distinct
	// elevations, numbered 0..MaxElevation.
	MaxElevation = 19
)

// World is the main container for the world
// representation of minima.
type World struct {
	// W and H are the width and height of the world's
	// location grid.
	W, H int

	// locs is the grid of world locations.
	locs []Loc

	// X0 and Y0 are the start location.
	X0, Y0 int
}

// A Loc is a cell in the grid that represents the world
type Loc struct {
	// Terrain is the type of this locations terrain.
	Terrain *TerrainType

	// X and Y are the coordinates of this location
	X, Y int

	// Elevation is the elevation of the ground or the
	// surface of a body of water.
	Elevation int

	// Depth is the depth of water at this location.
	Depth int
}

// Height returns the elevation of the location minus
// its depth.  This is effectively the elevation of the ground
// or the elevation of the floor of a body of water.
func (l Loc) Height() int {
	return l.Elevation - l.Depth
}

// TerrainType holds information on a given type of terrain.
type TerrainType struct {
	// Char is the character representing this terrain type.
	Char  rune

	// Name is a human readable name of the terrain type.
	Name  string

	// Color is the color of the highest elevation when
	// locations of this type are drawn.
	Color color.RGBA
}

// Terrain is an array with the canonical terrain
// representations, indexed by the terrain type's
// unique character.
var Terrain = []TerrainType{
	'g': {'g', "grass", color.RGBA{0, 255, 0, 255}},
	'm': {'m', "mountain", color.RGBA{196, 128, 0, 255}},
	'w': {'w', "water", color.RGBA{0, 0, 255, 255}},
	'l': {'l', "lava", color.RGBA{255, 0, 0, 255}},
	'd': {'d', "desert", color.RGBA{255, 255, 0, 255}},
	'f': {'f', "forest", color.RGBA{0, 200, 128, 255}},
	'i': {'i', "glacier", color.RGBA{196, 196, 255, 255}},
}

// Make returns a world of the given dimensions.
func Make(w, h int) World {
	const maxInt = int(^uint(0) >> 1)
	if w <= 0 || h <= 0 {
		panic("World dimensions must be positive")
	}
	if maxInt/w < h {
		panic("The world dimensions are too big")
	}
	return World{
		W:    w,
		H:    h,
		locs: makeLocs(w, h),
	}
}

// makeLocs returns an array of locations with
// initialized X and Y fields.
func makeLocs(w, h int) []Loc {
	locs := make([]Loc, w*h)
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			l := &locs[x*h+y]
			l.X = x
			l.Y = y
		}
	}
	return locs
}

// At returns the location at the given world coordinate.
func (w *World) At(x, y int) *Loc {
	x, y = w.Wrap(x, y)
	return &w.locs[x*w.H+y]
}

// Wrap returns an x,y within the ranges 0–width-1 and
// 0–height-1.  This effectively maps world coordinates
// to a normalized point on a grid, making the world a
// torus shape.
func (w *World) Wrap(x, y int) (int, int) {
	return wrap(x, w.W), wrap(y, w.H)
}

// wrap returns the value of n wrapped around if it goes
// above bound-1 or below zero.
func wrap(n, bound int) int {
	if n >= 0 && n < bound {
		return n
	}

	if bound <= 0 {
		panic("Bad bound in wrap")
	}
	n %= bound
	if n < 0 {
		n = bound + n
		if n < 0 {
			panic("A value wrapped to a negative")
		}
	}
	return n
}

// LocsWithType returns a slice of pointers to all of the
// locations with any of the given types.
func (w *World) LocsWithType(types string) ([]*Loc) {
	var locs []*Loc
	for i, loc := range w.locs {
		if strings.ContainsRune(types, rune(loc.Terrain.Char)) {
			locs = append(locs, &w.locs[i])
		}
	}
	return locs
}

// Write writes a world.  If the writer is not buffered then
// it is wrapped in a buffered writer, so the caller does not
// need to worry about buffering writes.
func (w *World) Write(out io.Writer) error {
	var err error
	if _, err = fmt.Fprintln(out, "#", runtime.GOOS, runtime.GOARCH); err != nil {
		return err
	}
	if _, err = fmt.Fprintln(out, w.W, w.H); err != nil {
		return err
	}
	for _, l := range w.locs {
		if l.Terrain == nil {
			panic("Nil terrain")
		}
		if _, err = fmt.Fprintf(out, "%c %d %d\n", l.Terrain.Char, l.Elevation, l.Depth); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(out, w.X0, w.Y0)
	return err
}

// Read reads a world.  If an error is encountered then
// the error is returned as the second argument and
// the zero-world is returned as the first.
func Read(in *bufio.Reader) (World, error) {
	var err error
	var line string
	if line, err = readLine(in); err != nil {
		return World{}, err
	}
	var width, height int
	if _, err = fmt.Sscanln(line, &width, &height); err != nil {
		fmt.Fprintln(os.Stderr, "failed to scan", line)
		return World{}, err
	}

	w := Make(width, height)
	for i := range w.locs {
		var el, dp int
		var ch rune
		if line, err = readLine(in); err != nil {
			return World{}, err
		}
		if _, err = fmt.Sscanf(line, "%c %d %d", &ch, &el, &dp); err != nil {
			fmt.Fprintln(os.Stderr, "failed to scan", line)
			return World{}, err
		}
		if el < 0 || el > MaxElevation {
			err = fmt.Errorf("Location %d elevation %d is out of bounds", i, el)
			return World{}, err
		}
		if int(ch) >= len(Terrain) || Terrain[ch].Char == 0 {
			err = fmt.Errorf("Location %d invalid terrain: %c",
				i, ch)
			return World{}, err
		}
		w.locs[i].Terrain = &Terrain[int(ch)]
		w.locs[i].Elevation = el
		w.locs[i].Depth = dp
	}

	if line, err = readLine(in); err != nil {
		return World{}, err
	}
	_, err = fmt.Sscanln(line, &w.X0, &w.Y0)
	return w, err
}

// ReadLine returns the next non-comment line.  On error
// the empty string and error are returned.
func readLine(in *bufio.Reader) (string, error) {
	for {
		bytes, prefix, err := in.ReadLine()
		if prefix {
			err = errors.New("Line is too long")
		}
		if err != nil {
			return "", err
		}
		if bytes[0] != '#' {
			return string(bytes), nil
		}
	}
	panic("Unreachable")
}
