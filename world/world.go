package world

const (
	// MaxHeight is the maximum height
	// value of any location in the world.
	MaxHeight = 8
)

// World is the main container for the world
// representation of minima.
type World struct{
	W, H int

	// locs is the grid of world locations.
	locs []Location
}

// A Location is a cell in the grid that
// represents the world
type Location struct{
	Terrain *TerrainType
	Height int
}

// Make returns a world of the given
// dimensions.
func Make(w, h int) World{
	const maxInt = int(^uint(0)>>1)
	if w <= 0 || h <= 0 {
		panic("World dimensions must be positive")
	}
	if maxInt/w < h {
		panic("The world dimensions are too big")
	}
	return World{
		W: w,
		H: h,
		locs: make([]Location, w*h),
	}
}

// At returns the location at the given x, y grid cell.
//
// Unlike AtCoord(), this roution does not wrap the
// x,y values around the boundaries of the grid.
func (w *World) At(x, y int) *Location{
	return &w.locs[x*w.H+y]
}

// AtCoord returns a pointer to the location at
// the given world coordinate.
func (w *World) AtCoord(x, y int) *Location{
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