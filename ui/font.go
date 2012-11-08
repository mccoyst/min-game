// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
// Copyright 2012 The Plotinum Authors.


package ui

import (
	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
	"image"
	"image/color"
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
}

// NewFont returns a new Font loaded from a .ttf file.
// The size parameter is specified in Postscript points.
func NewFont(path string, size float64, col color.Color) (*Font, error) {
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

	ctx := freetype.NewContext()
	ctx.SetFont(ttf)
	ctx.SetFontSize(size)
	ctx.SetSrc(image.NewUniform(col))
	ctx.SetDPI(pxInch)

	return &Font{ size: size, ttf: ttf, ctx: ctx }, nil
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
	scale := (f.size/ptInch * pxInch)/float64(em)
	return FontExtents{
		Height:  int(float64(bounds.YMax- bounds.YMin)*scale + 0.5),
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
	scale := (f.size/ptInch * pxInch)/float64(em)
	return int(float64(width)*scale + 0.5)
}

// Render returns an image of the given string.
func (f *Font) Render(s string) (image.Image, error) {
	w := f.Width(s)
	ext := f.Extents()
	h := ext.Height

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	f.ctx.SetClip(img.Bounds())
 	f.ctx.SetDst(img)

	_, err := f.ctx.DrawString(s, freetype.Pt(0, int(h+ext.Descent)))
	return img, err
}
