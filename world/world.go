// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
package world

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	Terrain          *TerrainType
	Elevation, Depth int
}

// Height returns the height of the location, which is
// its elevation minus its depth.
func (l Loc) Height() int {
	return l.Elevation - l.Depth
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
		locs: make([]Loc, w*h),
	}
}

// At returns the location at the given world coordinate.
func (w *World) At(x, y int) *Loc {
	return &w.locs[w.Index(x, y)]
}

// Index returns an array index for a world coordinate.
func (w *World) Index(x, y int) int {
	x, y = w.Wrap(x, y)
	return x*w.H + y
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
		var ch uint8
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
		if int(ch) >= len(Terrain) || Terrain[int(ch)].Char == uint8(0) {
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
