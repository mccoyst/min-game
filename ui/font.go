// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
// Copyright 2012 The Plotinum Authors.

package ui

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"
)

const (
	// PxInch is the number of pixels in an inch.
	// 72 seems to make font sizes consistent with
	// SDL_ttf…
	pxInch = 72.0

	// PtInch is the number of Postscript points in an inch.
	ptInch = 72.0
)

// A Font is a size and style in which text can be
// rendered to an image.
type Font struct {
	// Size is the size of the font in points.
	size float64

	// Ttf is the truetype font handle.
	ttf *truetype.Font

	// Ctx is the freetype drawing context for this
	// font, size, and color.
	ctx *freetype.Context

	// Glyphs is a cache of pre-rendered glyphs.
	glyphs []image.Image
}

// NewFont returns a new Font loaded from a .ttf file.
// The default size is 12 points, and the default color
// is black.
func NewFont(path string) (*Font, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	ttf, err := freetype.ParseFont(bytes)
	if err != nil {
		return nil, err
	}

	const defaultSize = 12.
	ctx := freetype.NewContext()
	ctx.SetFont(ttf)
	ctx.SetFontSize(defaultSize)
	ctx.SetSrc(image.NewUniform(color.Black))
	ctx.SetDPI(pxInch)

	return &Font{size: defaultSize, ttf: ttf, ctx: ctx}, nil
}

// SetSize sets the font size.
func (f *Font) SetSize(sz float64) {
	f.ctx.SetFontSize(sz)
	f.size = sz
}

// SetColor sets the font color.
func (f *Font) SetColor(col color.Color) {
	f.ctx.SetSrc(image.NewUniform(col))
}

// FontExtents describes some size attributes of all text
// rendered in this font.
type FontExtents struct {

	// Height is the height of the font in pixels.
	// The height is size used to create the font, and it
	// can be used as a reasonable line spacing.
	Height int

	// Ascent is the distance in pixels from the baseline
	// (the y-value given as the font's drawing location
	// to the highest point of a rasterized glyph.
	Ascent int

	// Decent is the distance in pixels from the baseline
	// to the lowest point of a rasterized glyph.  The
	// descent extents below the baseline so it is always
	// a negative value.
	Descent int
}

// Extents returns the extents of a font.
func (f *Font) Extents() FontExtents {
	em := f.ttf.FUnitsPerEm()
	bounds := f.ttf.Bounds(em)
	scale := (f.size / ptInch * pxInch) / float64(em)
	return FontExtents{
		Height:  int(float64(bounds.YMax-bounds.YMin)*scale + 0.5),
		Ascent:  int(float64(bounds.YMax)*scale + 0.5),
		Descent: int(float64(bounds.YMin)*scale + 0.5),
	}
}

// Width returns width of a string in pixels.
func (f *Font) Width(s string) int {
	em := f.ttf.FUnitsPerEm()
	var width int32
	prev, hasPrev := truetype.Index(0), false
	for _, r := range s {
		index := f.ttf.Index(r)
		if hasPrev {
			width += f.ttf.Kerning(em, prev, index)
		}
		width += f.ttf.HMetric(em, index).AdvanceWidth
		prev, hasPrev = index, true
	}
	scale := (f.size / ptInch * pxInch) / float64(em)
	return int(float64(width)*scale + 0.5)
}

// Render returns an image of the given string.
func (f *Font) Render(s string) (image.Image, error) {
	w, h := f.Width(s), f.Extents().Height
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	em := f.ttf.FUnitsPerEm()
	scale := (f.size / ptInch * pxInch) / float64(em)

	var x int32
	prev, hasPrev := truetype.Index(0), false
	for _, r := range s {
		index := f.ttf.Index(r)
		if hasPrev {
			x += f.ttf.Kerning(em, prev, index)
		}

		g, err := f.glyph(r)
		if err != nil {
			return nil, err
		}
		b := g.Bounds().Add(image.Pt(int(float64(x)*scale), 0))
		draw.Draw(img, b, g, image.ZP, draw.Src)

		x += f.ttf.HMetric(em, index).AdvanceWidth
		prev, hasPrev = index, true
	}

	return img, nil
}

// Glyph returns an image.Image containing the glyph.
// If the glyph is in the cache then that is returned,
// otherwise the glyph is rendered, cached, and returned.
func (f *Font) glyph(r rune) (image.Image, error) {
	i := int(f.ttf.Index(r))
	if i < len(f.glyphs) && f.glyphs[i] != nil {
		return f.glyphs[i], nil
	}

	s := string([]rune{r})
	w := f.Width(s)
	ext := f.Extents()
	h := ext.Height

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	f.ctx.SetClip(img.Bounds())
	f.ctx.SetDst(img)

	pt := freetype.Pt(0, int(h+ext.Descent))
	if _, err := f.ctx.DrawString(s, pt); err != nil {
		return nil, err
	}

	if i >= len(f.glyphs) {
		gs := make([]image.Image, i+1)
		copy(gs, f.glyphs)
		f.glyphs = gs
	}
	f.glyphs[i] = img

	return img, nil
}
