package world

import (
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
