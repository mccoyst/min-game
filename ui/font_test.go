package ui

import (
	"testing"
	"image/color"
	"math/rand"
)

func BenchmarkFontRender1(b *testing.B) {
	render(b, 1)
}

func BenchmarkFontRender2(b *testing.B) {
	render(b, 2)
}

func BenchmarkFontRender4(b *testing.B) {
	render(b, 4)
}

func BenchmarkFontRender8(b *testing.B) {
	render(b, 8)
}

func BenchmarkFontRender16(b *testing.B) {
	render(b, 16)
}

func BenchmarkFontRender32(b *testing.B) {
	render(b, 32)
}

// Render the given number of random ASCII characters.
func render(b *testing.B, cnt int) {
	b.StopTimer()
	f, err := NewFont("../resrc/prstartk.ttf", 12, color.Black)
	if err != nil {
		b.Fatal(err.Error())
	}


	const NStrs = 100
	strs := make([]string, NStrs)
	bytes := make([]byte, cnt)
	for i := range strs {
		for j := range bytes {
			// ASCII 33—126 are the printable characters.
			bytes[j] = byte(rand.Int31n(126-33 + 1) + 33)
		}
		strs[i] = string(bytes)
	}

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		if _, err = f.Render(strs[i%len(strs)]); err != nil {
			b.Fatal(err.Error())
		}
	}

}