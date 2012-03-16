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

// Index returns the array index of the
// given x and y coordinates of the world.
// Coordinates beyond the bounds of the
// array wrap around, i.e., this always
// returns a valid index into an array
// that matches the size of the world.
//
// This function is intended indexing into
// arrays of auxiliary data for each world
// location.
func (w *World) Index(x, y int) int {
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

// At returns a pointer to the location at
// the given x, y coordinate.
func (w *World) At(x, y int) *Location{
	return &w.locs[w.Index(x, y)]
}
