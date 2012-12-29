// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package world

import (
	"bufio"
	"os"
	"testing"
)

func TestWrap(t *testing.T) {
	tests := [...]struct{ x, bound, w int }{
		{0, 500, 0},
		{-1, 500, 499},
		{-500, 500, 0},
		{500, 500, 0},
		{501, 500, 1},
	}
	for _, tst := range tests {
		if w := wrap(tst.x, tst.bound); w != tst.w {
			t.Error("Expected wrap(%d, %d)=%d, got %d", tst.x, tst.bound, tst.w, w)
		}
	}
}

// TestWriteRead tests writing a world and reading it back.
func TestWriteRead(t *testing.T) {
	var types []string
	for _, t := range Terrain {
		if t.Char != "" {
			types = append(types, t.Char)
		}
	}

	w := New(10, 10)
	for i := range w.locs {
		w.locs[i].Elevation = i % (MaxElevation + 1)
		w.locs[i].Terrain = &Terrain[types[i%len(types)][0]]
	}

	read, write, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe() error: %s", err)
	}

	if err := w.Write(write); err != nil {
		t.Fatalf("Failed to write the world: %s", err)
	}

	u, err := Read(bufio.NewReader(read))
	if err != nil {
		t.Fatalf("Failed to read the world: %s", err)
	}

	if w.W != u.W {
		t.Errorf("World widths don't match: %d and %d", w.W, u.W)
	}

	if w.H != u.H {
		t.Errorf("World heights don't match: %d and %d", w.H, u.H)
	}

	if len(w.locs) != len(u.locs) {
		t.Fatalf("Worlds have different numbers of locations: %d and %d",
			len(w.locs), len(u.locs))
	}

	for i, l0 := range w.locs {
		l1 := u.locs[i]
		if l0.Terrain != l1.Terrain {
			ch0, ch1 := "?", "?"
			if l0.Terrain != nil {
				ch0 = l0.Terrain.Char
			}
			if l1.Terrain != nil {
				ch1 = l1.Terrain.Char
			}
			t.Fatalf("Location %d: terrain mismatch: %p (%c) and %p (%c)", i,
				l0.Terrain, ch0, l1.Terrain, ch1)
		}
		if l0.Height() != l1.Height() {
			t.Fatalf("Location %d: height mismatch", i)
		}
	}
}
