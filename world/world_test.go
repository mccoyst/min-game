// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package world

import (
	"bufio"
	"errors"
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
		te := rand.Intn(len(types))
		w.locs[i].Terrain = &Terrain[types[te][0]]
	}

	u, err := writeRead(w)
	if err != nil {
		t.Error(err.Error())
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

	u, err := writeRead(w)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(w, u) {
		t.Error("Worlds don't match")
	}
}

// TestWriteReadRuns tests writing a world where some locations are the same.
func TestWriteReadRuns(t *testing.T) {
	var types []string
	for _, t := range Terrain {
		if t.Char != "" {
			types = append(types, t.Char)
		}
	}

	run := rand.Intn(4) + 1
	var el, de, te int

	w := New(10, 10)
	for i := range w.locs {
		if run == 0 {
			run = rand.Intn(4) + 1
			el = rand.Intn(MaxElevation)
			de = rand.Intn(el+1) - 1
			te = rand.Intn(len(types))
		}
		run--

		w.locs[i].Elevation = el
		w.locs[i].Depth = de
		w.locs[i].Terrain = &Terrain[types[te][0]]
	}

	u, err := writeRead(w)
	if err != nil {
		t.Error(err.Error())
	}

	if !reflect.DeepEqual(w, u) {
		t.Error("Worlds don't match")
	}
}

// WriteRead writes the given world, reads it, and returns what it read.
func writeRead(w *World) (*World, error) {
	read, write, err := os.Pipe()
	if err != nil {
		return nil, errors.New("Failed to create a pipe: " + err.Error())
	}
	if err := w.Write(write); err != nil {
		return nil, errors.New("Failed to write the world: " + err.Error())
	}
	u, err := Read(bufio.NewReader(read))
	if err != nil {
		err = errors.New("Failed to read the world: " + err.Error())
	}
	return u, err
}
