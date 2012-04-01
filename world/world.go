package world

import (
	"fmt"
	"io"
)

const (
	// NHeights is the number of distinct
	// heights, numbered 0..MaxHeight.
	MaxHeight = 19
)

// World is the main container for the world
// representation of minima.
type World struct {
	W, H int

	// locs is the grid of world locations.
	locs []Loc
}

// A Loc is a cell in the grid that
// represents the world
type Loc struct {
	Terrain *TerrainType
	Height, Depth  int
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
	if _, err = fmt.Fprintln(out, w.W, w.H); err != nil {
		return
	}

	for _, l := range w.locs {
		if l.Terrain == nil {
			panic("Nil terrain")
		}
		if _, err = fmt.Fprintf(out, " %c %d %d", l.Terrain.Char, l.Height, l.Depth); err != nil {
			return
		}
	}

	fmt.Fprintln(out)

	return
}

// Read reads the world from the given io.Reader
// and returns it.  If an error is encountered then
// the error is returned as the second argument and
// the zero-world is returned as the first.
func Read(in io.Reader) (_ World, err error) {
	var width, height, n int
	if n, err = fmt.Fscanln(in, &width, &height); n != 2 || err != nil {
		if err == nil {
			err = fmt.Errorf("Failed to scan width and height")
		}
		return
	}

	w := Make(width, height)
	for i := range w.locs {
		var ht, dp int
		var ch uint8
		if n, err = fmt.Fscanf(in, " %c %d %d", &ch, &ht, &dp); n != 3 || err != nil {
			if err == nil {
				err = fmt.Errorf("Failed to scan location %d", i)
			}
			return
		}
		if ht < 0 || ht > MaxHeight {
			err = fmt.Errorf("Location %d height %d is out of bounds", i, ht)
			return
		}
		if int(ch) >= len(Terrain) || Terrain[int(ch)].Char == uint8(0) {
			err = fmt.Errorf("Location %d invalid terrain: %c",
				i, ch)
			return
		}
		w.locs[i].Terrain = &Terrain[int(ch)]
		w.locs[i].Height = ht
		w.locs[i].Depth = dp
	}

	return w, nil
}
