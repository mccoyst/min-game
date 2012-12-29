// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package world

import (
	"bufio"
	"math/rand"
	"os"
	"reflect"
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
		w.locs[i].Elevation = rand.Intn(MaxElevation-1) + 1
		w.locs[i].Depth = rand.Intn(w.locs[i].Elevation)
		t := rand.Intn(len(types))
		w.locs[i].Terrain = &Terrain[types[t][0]]
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

	if !reflect.DeepEqual(w, u) {
		t.Error("Worlds don't match")
	}
}

// TestWriteReadSame tests writing a world where all locations are
// identical and reading it back.
func TestWriteReadSame(t *testing.T) {
	w := New(10, 10)
	for i := range w.locs {
		w.locs[i].Elevation = 1
		w.locs[i].Depth = 0
		w.locs[i].Terrain = &Terrain[int('g')]
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

	if !reflect.DeepEqual(w, u) {
		t.Error("Worlds don't match")
	}
}
