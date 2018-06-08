// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
// Copyright 2012 The Plotinum Authors.

package ui

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"
	"os"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

const (
	// PxInch is the number of pixels in an inch.
	// 72 seems to make font sizes consistent with
	// SDL_ttf…
	pxInch = 72.0

	// PtInch is the number of Postscript points in an inch.
	ptInch = 72.0
)

type glyphKey struct {
	rune       rune
	size       float64
	r, g, b, a uint32
}

// A font is a size and style in which text can be
// rendered to an image.
type font struct {
	// Size is the size of the font in points.
	size float64

	// color is the current color.
	color color.Color

	// Ttf is the truetype font handle.
	ttf *truetype.Font

	// Ctx is the freetype drawing context for this
	// font, size, and color.
	ctx *freetype.Context

	// Glyphs is a cache of pre-rendered glyphs.
	glyphs map[glyphKey]image.Image
}

// NewFont returns a new font loaded from a .ttf file.
// The default size is 12 points, and the default color
// is black.
func newFont(path string) (*font, error) {
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
	defaultColor := color.RGBA{0, 0, 0, 255}
	ctx := freetype.NewContext()
	ctx.SetFont(ttf)
	ctx.SetDPI(pxInch)

	return &font{
		size:   defaultSize,
		color:  defaultColor,
		ttf:    ttf,
		ctx:    ctx,
		glyphs: make(map[glyphKey]image.Image),
	}, nil
}

// SetSize sets the font size.
func (f *font) setSize(sz float64) {
	f.size = sz
}

// SetColor sets the font color.
func (f *font) setColor(col color.Color) {
	f.color = col
}

// fontExtents describes some size attributes of all text
// rendered in this font.
type fontExtents struct {

	// Height is the height of the font in pixels.
	// The height is size used to create the font, and it
	// can be used as a reasonable line spacing.
	height int

	// Ascent is the distance in pixels from the baseline
	// (the y-value given as the font's drawing location
	// to the highest point of a rasterized glyph.
	ascent int

	// Decent is the distance in pixels from the baseline
	// to the lowest point of a rasterized glyph.  The
	// descent extents below the baseline so it is always
	// a negative value.
	descent int
}

// Extents returns the extents of a font.
func (f *font) extents() fontExtents {
	em := f.ttf.FUnitsPerEm()
	bounds := f.ttf.Bounds(fixed.I(int(em)))
	scale := (f.size / ptInch * pxInch) / float64(em)
	a := int(float64(bounds.Max.Y)*scale + 0.5)
	d := int(float64(bounds.Min.Y)*scale - 0.5)
	return fontExtents{
		height:  a - d,
		ascent:  a,
		descent: d,
	}
}

// Width returns width of a string in pixels.
func (f *font) width(s string) int {
	em := f.ttf.FUnitsPerEm()
	var width fixed.Int26_6
	prev, hasPrev := truetype.Index(0), false
	for _, r := range s {
		index := f.ttf.Index(r)
		if hasPrev {
			width += f.ttf.Kern(fixed.I(int(em)), prev, index)
		}
		width += f.ttf.HMetric(fixed.I(int(em)), index).AdvanceWidth
		prev, hasPrev = index, true
	}
	scale := (f.size / ptInch * pxInch) / float64(em)
	return int(float64(width)*scale + 0.5)
}

// Render returns an image of the given string.
func (f *font) render(s string) (image.Image, error) {
	w, h := f.width(s), f.extents().height
	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))

	em := f.ttf.FUnitsPerEm()
	scale := (f.size / ptInch * pxInch) / float64(em)

	var x fixed.Int26_6
	prev, hasPrev := truetype.Index(0), false
	for _, r := range s {
		index := f.ttf.Index(r)
		if hasPrev {
			x += f.ttf.Kern(fixed.I(int(em)), prev, index)
		}

		g, err := f.glyph(r)
		if err != nil {
			return nil, err
		}
		b := g.Bounds().Add(image.Pt(int(float64(x)*scale), 0))
		draw.Draw(img, b, g, image.ZP, draw.Src)

		x += f.ttf.HMetric(fixed.I(int(em)), index).AdvanceWidth
		prev, hasPrev = index, true
	}

	return img, nil
}

// Glyph returns an image.Image containing the glyph.
// If the glyph is in the cache then that is returned,
// otherwise the glyph is rendered, cached, and returned.
func (f *font) glyph(ru rune) (image.Image, error) {
	r, g, b, a := f.color.RGBA()
	key := glyphKey{
		rune: ru,
		size: f.size,
		r:    r,
		g:    g,
		b:    b,
		a:    a,
	}
	if img, ok := f.glyphs[key]; ok {
		return img, nil
	}
	s := string([]rune{ru})
	w := f.width(s)
	ext := f.extents()
	h := ext.height

	img := image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	f.ctx.SetFontSize(f.size)
	f.ctx.SetSrc(image.NewUniform(f.color))
	f.ctx.SetClip(img.Bounds())
	f.ctx.SetDst(img)

	pt := freetype.Pt(0, int(h+ext.descent))
	if _, err := f.ctx.DrawString(s, pt); err != nil {
		return nil, err
	}

	f.glyphs[key] = img

	return img, nil
}
