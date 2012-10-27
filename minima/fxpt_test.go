// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package main

import "testing"

func TestFp(t *testing.T) {
	a, b := Fp(0, 0), Fxpt{}
	if !a.Eq(b) {
		t.Error(a, "!=", b)
	}

	a, b = Fp(1, 0), Fxpt{1 << FixedPoint}
	if !a.Eq(b) {
		t.Error(a, "!=", b)
	}
}

func TestString(t *testing.T) {
	a := Fp(7, 8)
	if a.String() != "7.08" {
		t.Error(a, "!=", "7.08")
	}
}

func TestAdd(t *testing.T) {
	type test struct{ a, b, c Fxpt }
	tests := []test{
		test{Fp(1, 0), Fp(1, 0), Fp(2, 0)},
		test{Fp(0, 8), Fp(0, 8), Fp(1, 0)},
		test{Fp(0, 2), Fp(0, 3), Fp(0, 5)},
		test{Fp(0, 15), Fp(0, 1), Fp(1, 0)},
		test{Fp(0, 15), Fp(0, 9), Fp(1, 8)},
		test{Fp(-1, 0), Fp(1, 0), Fp(0, 0)},
	}

	for _, tt := range tests {
		if !tt.a.Add(tt.b).Eq(tt.c) {
			t.Error(tt.a, "+", tt.b, "==", tt.a.Add(tt.b), "!=", tt.c)
		}
	}
}

func TestSub(t *testing.T) {
	type test struct{ a, b, c Fxpt }
	tests := []test{
		test{Fp(1, 0), Fp(1, 0), Fp(0, 0)},
		test{Fp(1, 0), Fp(0, 0), Fp(1, 0)},
		test{Fp(0, 0), Fp(1, 0), Fp(-1, 0)},
		test{Fp(1, 0), Fp(2, 0), Fp(-1, 0)},
		test{Fp(0, 8), Fp(0, 8), Fp(0, 0)},
		test{Fp(1, 0), Fp(0, 1), Fp(0, 15)},
	}

	for _, tt := range tests {
		if !tt.a.Sub(tt.b).Eq(tt.c) {
			t.Error(tt.a, "-", tt.b, "==", tt.a.Sub(tt.b), "!=", tt.c)
		}
	}
}

func TestMul(t *testing.T) {
	type test struct{ a, b, c Fxpt }
	tests := []test{
		test{Fp(1, 0), Fp(0, 0), Fp(0, 0)},
		test{Fp(1, 0), Fp(1, 0), Fp(1, 0)},
		test{Fp(1, 0), Fp(-1, 0), Fp(-1, 0)},
		test{Fp(1, 0), Fp(2, 0), Fp(2, 0)},
	}

	for _, tt := range tests {
		if !tt.a.Mul(tt.b).Eq(tt.c) {
			t.Error(tt.a, "*", tt.b, "==", tt.a.Mul(tt.b), "!=", tt.c)
		}
	}
}

func TestDiv(t *testing.T) {
	type test struct{ a, b, c Fxpt }
	tests := []test{
		test{Fp(1, 0), Fp(1, 0), Fp(1, 0)},
		test{Fp(2, 0), Fp(1, 0), Fp(2, 0)},
		test{Fp(1, 0), Fp(2, 0), Fp(0, 1<<(FixedPoint-1))},
	}

	for _, tt := range tests {
		if !tt.a.Div(tt.b).Eq(tt.c) {
			t.Error(tt.a, "÷", tt.b, "==", tt.a.Div(tt.b), "!=", tt.c)
		}
	}
}

func TestRem(t *testing.T) {
	type test struct{ a, b, c Fxpt }
	tests := []test{
		test{Fp(100, 0), Fp(100, 0), Fp(0, 0)},
		test{Fp(100, 0), Fp(50, 0), Fp(0, 0)},
		test{Fp(100, 0), Fp(25, 0), Fp(0, 0)},
		test{Fp(100, 0), Fp(99, 0), Fp(1, 0)},
		test{Fp(100, 0), Fp(51, 0), Fp(49, 0)},
	}

	for _, tt := range tests {
		if !tt.a.Rem(tt.b).Eq(tt.c) {
			t.Error(tt.a, "%", tt.b, "==", tt.a.Rem(tt.b), "!=", tt.c)
		}
	}
}
