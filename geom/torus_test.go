// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package geom

import (
	"math"
	"testing"
	"testing/quick"
)

func TestNearWrap(t *testing.T) {
	tests := []struct {
		a, b, width, c float64
	}{
		{0, 0, 100, 0},
		{0, 100, 100, 0},
		{0, 99, 100, -1},
		{99, 100, 100, 100},
	}

	for _, test := range tests {
		if c := nearWrap(test.a, test.b, test.width); c != test.c {
			t.Errorf("expected nearWrap(%g, %g, %g) to be %g, not %g",
				test.a, test.b, test.width, test.c, c)
		}
	}
}

func TestTorusDist(t *testing.T) {
	torus := Torus{W: 100, H: 100}
	tests := []struct {
		a, b Point
		d    float64
	}{
		{Pt(0, 0), Pt(0, 0), 0},
		{Pt(100, 100), Pt(0, 0), 0},
		{Pt(-100, -100), Pt(100, 100), 0},
		{Pt(-100, -100), Pt(200, 200), 0},
		{Pt(-1, -1), Pt(0, 0), math.Sqrt(2)},
		{Pt(99, 99), Pt(99, 99), 0},
		{Pt(100, 100), Pt(100, 100), 0},
		{Pt(0, 0), Pt(100, 100), 0},
		{Pt(99, 99), Pt(100, 100), math.Sqrt(2)},
		{Pt(99, 99), Pt(0, 0), math.Sqrt(2)},
		{Pt(99, 99), Pt(10, 10), math.Sqrt(11*11 + 11*11)},
	}

	for _, test := range tests {
		if d := torus.Dist(test.a, test.b); d != test.d {
			t.Errorf("expected %s and %s to be %g appart, not %g",
				test.a, test.b, test.d, d)
		}
	}
}

func TestTorusSub(t *testing.T) {
	torus := Torus{W: 100, H: 100}
	tests := []struct {
		a, b, d Point
	}{
		{Pt(0, 0), Pt(0, 0), Pt(0, 0)},
		{Pt(0, 0), Pt(1, 1), Pt(-1, -1)},
		{Pt(50, 50), Pt(25, 25), Pt(25, 25)},
		{Pt(99, 99), Pt(0, 0), Pt(-1, -1)},
		{Pt(0, 0), Pt(99, 99), Pt(1, 1)},
	}

	for _, test := range tests {
		if d := torus.Sub(test.a, test.b); d != test.d {
			t.Errorf("expected %s - %s to be %g, not %g",
				test.a, test.b, test.d, d)
		}
	}
}

func TestAlign(t *testing.T) {
	torus := Torus{100, 100}
	tests := []struct {
		a, b   Rectangle
		balign Rectangle
	}{
		{Rect(0, 0, 1, 1), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},
		{Rect(0, 0, 1, 1), Rect(100, 100, 101, 101), Rect(0, 0, 1, 1)},
		{Rect(0, 0, 1, 1), Rect(-100, -100, -99, -99), Rect(0, 0, 1, 1)},
		{Rect(0, 0, 1, 1), Rect(100.5, 100.5, 101.5, 101.5), Rect(0.5, 0.5, 1.5, 1.5)},
	}
	for _, test := range tests {
		a, b := torus.align(test.a, test.b)
		if a != torus.NormRect(test.a) {
			t.Fatalf("%s is not normalized %s, wanted %s", a, test.a, torus.NormRect(test.a))
		}
		if b != test.balign {
			t.Errorf("Algining %s with %s: expected %s, got %s", test.b, test.a, b, test.balign)
		}
	}
}

func TestIsectTorus(t *testing.T) {
	torus := Torus{100, 100}
	z := Rectangle{}
	tests := []struct {
		a, b  Rectangle
		isect Rectangle
	}{
		{Rect(0, 0, 1, 1), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},
		{Rect(0, 0, 0.5, 0.5), Rect(0, 0, 1, 1), Rect(0, 0, 0.5, 0.5)},
		{Rect(0, 0, 1, 0.5), Rect(0, 0, 1, 1), Rect(0, 0, 1, 0.5)},
		{Rect(0, 0, 0.5, 1), Rect(0, 0, 1, 1), Rect(0, 0, 0.5, 1)},

		{Rect(100, 100, 101, 101), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},
		{Rect(-100, -100, -99, -99), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},
		{Rect(0, 0, 1, 1), Rect(100, 100, 101, 101), Rect(0, 0, 1, 1)},
		{Rect(-100, -100, -99, -99), Rect(100, 100, 101, 101), Rect(0, 0, 1, 1)},

		{Rect(0, -100, 1, -99), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},
		{Rect(0, 100, 1, 101), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},

		//
		{Rect(0, 0, 1, 1), Rect(0, -100, 1, -99), Rect(0, 0, 1, 1)},

		{Rect(0, 0, 1, 1), Rect(0, 100, 1, 101), Rect(0, 0, 1, 1)},
		{Rect(-100, 0, -99, 1), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},
		{Rect(100, 0, 101, 1), Rect(0, 0, 1, 1), Rect(0, 0, 1, 1)},

		//
		{Rect(0, 0, 1, 1), Rect(-100, 0, -99, 1), Rect(0, 0, 1, 1)},

		{Rect(0, 0, 1, 1), Rect(100, 0, 101, 1), Rect(0, 0, 1, 1)},

		{Rect(0, 2, 1, 3), Rect(0, 0, 1, 1), z},
		{Rect(2, 0, 3, 1), Rect(0, 0, 1, 1), z},
	}
	for _, test := range tests {
		isect := torus.Intersect(test.a, test.b)
		if isect == test.isect {
			continue
		}
		t.Errorf("%s and %s should isect by %s, not %s", test.a, test.b, test.isect, isect)
	}
}

func TestNormRect(t *testing.T) {
	torus := Torus{100, 100}
	tests := []struct {
		a, n Rectangle
	}{
		{Rect(-100, 0, -99, 1), Rect(0, 0, 1, 1)},
		{Rect(0, -100, 1, -99), Rect(0, 0, 1, 1)},
	}
	for _, test := range tests {
		n := torus.NormRect(test.a)
		if n == test.n {
			continue
		}
		t.Errorf("%s should normalize to %s, not %s", test.a, test.n, n)
	}
}

func TestNormRectQuick(t *testing.T) {
	torus := Torus{100, 100}
	f := func(x, y, w, h float64) bool {
		r := torus.NormRect(Rect(x, y, x+w, y+w))
		return r.Min.X >= 0 && r.Min.X < torus.W && r.Min.Y >= 0 && r.Min.Y <= torus.H
	}
	quick.Check(f, nil)
}
