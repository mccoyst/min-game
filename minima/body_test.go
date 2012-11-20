package main

import (
	"code.google.com/p/min-game/ui"
	"testing"
	"testing/quick"
)

type alignTest struct {
	a, b    ui.Rectangle
	aligned ui.Rectangle
	ok      bool
}

func TestAlign(t *testing.T) {
	tests := []alignTest{
		{ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(100, 100, 101, 101), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(-100, -100, -99, -99), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(100.5, 100.5, 101.5, 101.5), ui.Rect(0, 0, 1, 1), ui.Rect(0.5, 0.5, 1.5, 1.5), true},
		{ui.Rect(-1, -1, 0, 0), ui.Rect(99, 99, 100, 100), ui.Rect(99, 99, 100, 100), true},
		{ui.Rect(100, 0, 101, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(-100, 0, -99, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(100.5, 0.5, 101.5, 1.5), ui.Rect(0, 0, 1, 1), ui.Rect(0.5, 0.5, 1.5, 1.5), true},
		{ui.Rect(-1, 99, 0, 100), ui.Rect(99, 99, 100, 100), ui.Rect(99, 99, 100, 100), true},
		{ui.Rect(0, 100, 1, 101), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(0, -100, 1, -99), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), true},
		{ui.Rect(0.5, 100.5, 1.5, 101.5), ui.Rect(0, 0, 1, 1), ui.Rect(0.5, 0.5, 1.5, 1.5), true},
		{ui.Rect(99, -1, 100, 0), ui.Rect(99, 99, 100, 100), ui.Rect(99, 99, 100, 100), true},

		{ui.Rect(0, 0, 1, 1), ui.Rect(99, 99, 100, 100), ui.Rect(0, 0, 0, 0), false},
		{ui.Rect(0, 0, 1, 1), ui.Rect(99, 99, 100, 100), ui.Rect(0, 0, 0, 0), false},
	}
	for _, test := range tests {
		aligned, ok := align(test.a, test.b, ui.Pt(100, 100))
		switch {
		case test.ok && !ok:
			t.Errorf("%s and %s didn't align", test.a, test.b)

		case test.ok && aligned != test.aligned:
			t.Errorf("%s and %s should align to %v, not %v", test.a, test.b, test.aligned, aligned)
		}

	}
}

type isectTest struct {
	a, b  ui.Rectangle
	isect ui.Rectangle
}

func TestIsectTorus(t *testing.T) {
	z := ui.Rectangle{}
	tests := []isectTest{
		{ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(0, 0, 0.5, 0.5), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 0.5, 0.5)},
		{ui.Rect(0, 0, 1, 0.5), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 0.5)},
		{ui.Rect(0, 0, 0.5, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 0.5, 1)},

		{ui.Rect(100, 100, 101, 101), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(-100, -100, -99, -99), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(0, 0, 1, 1), ui.Rect(100, 100, 101, 101), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(-100, -100, -99, -99), ui.Rect(100, 100, 101, 101), ui.Rect(0, 0, 1, 1)},

		{ui.Rect(0, -100, 1, -99), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(0, 100, 1, 101), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},

		//
		{ui.Rect(0, 0, 1, 1), ui.Rect(0, -100, 1, -99), ui.Rect(0, 0, 1, 1)},

		{ui.Rect(0, 0, 1, 1), ui.Rect(0, 100, 1, 101), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(-100, 0, -99, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(100, 0, 101, 1), ui.Rect(0, 0, 1, 1), ui.Rect(0, 0, 1, 1)},

		//
		{ui.Rect(0, 0, 1, 1), ui.Rect(-100, 0, -99, 1), ui.Rect(0, 0, 1, 1)},

		{ui.Rect(0, 0, 1, 1), ui.Rect(100, 0, 101, 1), ui.Rect(0, 0, 1, 1)},

		{ui.Rect(0, 2, 1, 3), ui.Rect(0, 0, 1, 1), z},
		{ui.Rect(2, 0, 3, 1), ui.Rect(0, 0, 1, 1), z},
	}
	for _, test := range tests {
		isect := isectTorus(test.a, test.b, ui.Pt(100, 100))
		if isect == test.isect {
			continue
		}
		t.Errorf("%s and %s should isect by %s, not %s", test.a, test.b, test.isect, isect)
	}
}

func TestNormRect(t *testing.T) {
	tests := []struct {
		a, n ui.Rectangle
	}{
		{ui.Rect(-100, 0, -99, 1), ui.Rect(0, 0, 1, 1)},
		{ui.Rect(0, -100, 1, -99), ui.Rect(0, 0, 1, 1)},
	}
	for _, test := range tests {
		n := normRect(test.a, ui.Pt(100, 100))
		if n == test.n {
			continue
		}
		t.Errorf("%s should normalize to %s, not %s", test.a, test.n, n)
	}
}

func TestNormRectQuick(t *testing.T) {
	f := func(x, y, w, h float64) bool {
		r := normRect(ui.Rect(x, y, x+w, y+w), ui.Pt(100, 100))
		return r.Min.X >= 0 && r.Min.X < 100 && r.Min.Y >= 0 && r.Min.Y <= 100
	}
	quick.Check(f, nil)
}
