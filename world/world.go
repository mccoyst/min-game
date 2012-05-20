// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
package world

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
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

// A Loc is a cell in the grid that
// represents the world
type Loc struct {
	Terrain *TerrainType
	Elevation, Depth  int
}

// Height returns the height of the location, which is
// its elevation minus its depth.
func (l Loc) Height() int {
	return l.Elevation - l.Depth
}

// Make returns a world of the given
// dimensions.
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

// At returns the location at the given x, y grid cell.
//
// Unlike AtCoord(), this roution does not wrap the
// x,y values around the boundaries of the grid.
func (w *World) At(x, y int) *Loc {
	return &w.locs[x*w.H+y]
}

// AtCoord returns a pointer to the location at
// the given world coordinate.
func (w *World) AtCoord(x, y int) *Loc {
	return &w.locs[w.CoordToIndex(x, y)]
}

// WrapCord returns the x,y point for a coordinate
// after wrapping it around the world.
func (w *World) WrapCoord(x, y int) (int, int) {
	return wrap(x, w.W), wrap(y, w.H)
}

// CoordToIndex returns the array index that
// corresponds to the given x,y world coordinate.
func (w *World) CoordToIndex(x, y int) int {
	x = wrap(x, w.W)
	y = wrap(y, w.H)
	return x*w.H + y
}

// wrap returns the value of n wrapped
// around if it goes above bound-1 or
// below zero.
func wrap(n, bound int) int {
	// probably quicker to do this test for the
	// common case than to bother using %
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

// Write writes the world to the given io.Writer.
func (w *World) Write(out io.Writer) (err error) {
	fmt.Fprintln(out, "#", runtime.GOOS, runtime.GOARCH)

	if _, err = fmt.Fprintln(out, w.W, w.H); err != nil {
		return
	}

	for _, l := range w.locs {
		if l.Terrain == nil {
			panic("Nil terrain")
		}
		_, err = fmt.Fprintf(out, "%c %d %d\n", l.Terrain.Char, l.Elevation, l.Depth)
		if err != nil {
			return
		}
	}

	fmt.Fprintln(out, w.X0, w.Y0)

	return
}

// Read reads the world from the given io.Reader
// and returns it.  If an error is encountered then
// the error is returned as the second argument and
// the zero-world is returned as the first.
func Read(in *bufio.Reader) (_ World, err error) {
	line, err := readLine(in)
	if err != nil {
		return
	}
	var width, height int
	_, err = fmt.Sscanln(line, &width, &height)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to scan", line)
		return
	}

	w := Make(width, height)
	for i := range w.locs {
		var el, dp int
		var ch uint8
		line, err = readLine(in)
		if err != nil {
			return
		}
		_, err = fmt.Sscanf(line, "%c %d %d", &ch, &el, &dp)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to scan", line)
			return
		}
		if el < 0 || el > MaxElevation {
			err = fmt.Errorf("Location %d elevation %d is out of bounds", i, el)
			return
		}
		if int(ch) >= len(Terrain) || Terrain[int(ch)].Char == uint8(0) {
			err = fmt.Errorf("Location %d invalid terrain: %c",
				i, ch)
			return
		}
		w.locs[i].Terrain = &Terrain[int(ch)]
		w.locs[i].Elevation = el
		w.locs[i].Depth = dp
	}

	line, err = readLine(in)
	if err != nil {
		return
	}
	_, err = fmt.Sscanln(line, &w.X0, &w.Y0)
	return w, err
}

func readLine(in *bufio.Reader) (string, error) {
	for {
		r, _, err := in.ReadRune()
		if err != nil {
			return "", err
		}
		in.UnreadRune()
	
		if !unicode.IsSpace(r) && r != '#' {
			bytes, prefix, err := in.ReadLine()
			if prefix {
				err = fmt.Errorf("Line is too long")
			}
			return string(bytes), err
		}
		_, prefix, err := in.ReadLine()
		for prefix && err == nil {
			_, prefix, err = in.ReadLine()
		}
		if err != nil {
			return "", err
		}
	}
	panic("Unreachable")
}
