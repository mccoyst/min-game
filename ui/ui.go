// © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.

package ui

/*
#include <SDL.h>

#cgo darwin CFLAGS: -I/Users/steve/Code/src/sdl2/include
#cgo darwin LDFLAGS: -framework SDL2

#cgo linux CFLAGS: -I/usr/local/include/SDL2
#cgo linux LDFLAGS: -L/usr/local/lib -lSDL2

static Uint32 sdl_event_type(SDL_Event *e){
	return e->type;
}

static Uint32 sdl_rgba_fmt(int isLE){
	// SDL is doing some stupid byte-order-specific garbage.
	if(isLE)
		return SDL_PIXELFORMAT_ABGR8888;
	return SDL_PIXELFORMAT_RGBA8888;
}
*/
import "C"

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"unsafe"
)

type Event interface{}
type Quit struct{}
type KeyCode C.SDL_Keycode
type Key struct {
	Down   bool
	Repeat bool
	Code   KeyCode
}

func (k KeyCode) String() string {
	return C.GoString(C.SDL_GetKeyName(C.SDL_Keycode(k)))
}

type Button int

const (
	Unknown Button = iota
	Left
	Right
	Up
	Down
	Shoot
	Bomb
)

var ButtonNames = map[Button]string{
	Unknown: "Unknown",
	Left: "Left",
	Right: "Right",
	Up: "Up",
	Down: "Down",
	Shoot: "Shoot",
	Bomb: "Bomb",
}

func (b Button) String() string {
	if s, ok := ButtonNames[b]; ok {
		return s
	}
	return ButtonNames[Unknown]
}

var DefaultKeymap = map[KeyCode]Button{
	KeyCode(C.SDLK_s): Left,
	KeyCode(C.SDLK_f): Right,
	KeyCode(C.SDLK_e): Up,
	KeyCode(C.SDLK_d): Down,
	KeyCode(C.SDLK_j): Shoot,
	KeyCode(C.SDLK_k): Bomb,
}

// SDL-specific:

type Ui struct {
	win  *C.SDL_Window
	rend *C.SDL_Renderer

	imgCache map[string]*sdlImg
	fontCache map[string]*Font
}

func New(title string, w, h int) (*Ui, error) {
	e := C.SDL_Init(C.SDL_INIT_EVERYTHING)
	if e != 0 {
		return nil, sdlError()
	}

	t := C.CString(title)
	defer C.free(unsafe.Pointer(t))
	win := C.SDL_CreateWindow(
		t,
		C.SDL_WINDOWPOS_UNDEFINED,
		C.SDL_WINDOWPOS_UNDEFINED,
		C.int(w),
		C.int(h),
		C.SDL_WINDOW_SHOWN)
	if win == nil {
		return nil, sdlError()
	}

	rend := C.SDL_CreateRenderer(win, -1, C.SDL_RENDERER_PRESENTVSYNC)
	if rend == nil {
		return nil, sdlError()
	}

	return &Ui{win: win, rend: rend, imgCache: make(map[string]*sdlImg)}, nil
}

func (ui *Ui) Close() {
	C.SDL_DestroyRenderer(ui.rend)
	C.SDL_DestroyWindow(ui.win)
	C.SDL_Quit()
}

func (ui *Ui) PollEvent() Event {
	var e C.SDL_Event
	if C.SDL_PollEvent(&e) == 0 {
		return nil
	}

	switch C.sdl_event_type(&e) {
	case C.SDL_QUIT:
		return Quit{}
	case C.SDL_KEYDOWN, C.SDL_KEYUP:
		k := (*C.SDL_KeyboardEvent)(unsafe.Pointer(&e))
		return Key{
			Down:   k._type == C.SDL_KEYDOWN,
			Repeat: k.repeat != 0,
			Code:   KeyCode(k.keysym.sym),
		}
	}

	return nil
}

func (ui *Ui) Clear() {
	C.SDL_RenderClear(ui.rend)
}

func (ui *Ui) Sync() error {
	C.SDL_RenderPresent(ui.rend)
	return nil
}

type sdlImg struct {
	tex *C.SDL_Texture
}

func (s *sdlImg) Close() {
	C.SDL_DestroyTexture(s.tex)
}

func loadImg(ui *Ui, path string) (*sdlImg, error) {
	if img, ok := ui.imgCache[path]; ok {
		return img, nil
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}

	return newSdlImage(ui, img, path)
}

func newSdlImage(ui *Ui, img image.Image, path string) (*sdlImg, error) {
	bounds := img.Bounds()
	newTex := func(format C.Uint32) *C.SDL_Texture {
		return C.SDL_CreateTexture(
			ui.rend,
			format,
			C.SDL_TEXTUREACCESS_STATIC,
			C.int(bounds.Dx()),
			C.int(bounds.Dy()))
	}

	rgba := asRgba(img)
	tex := newTex(C.sdl_rgba_fmt(C.int(isLE)))
	e := C.SDL_UpdateTexture(tex, nil, unsafe.Pointer(&rgba.Pix[0]), C.int(rgba.Stride))
	if e != 0 {
		C.SDL_DestroyTexture(tex)
		return nil, sdlError()
	}
	C.SDL_SetTextureBlendMode(tex, C.SDL_BLENDMODE_BLEND)
	si := &sdlImg{tex}
	if path != "" {
		ui.imgCache[path] = si
	}
	return si, nil
}

// BUG(mccoyst): asRgba assumes the image bounds starts at (0,0).
func asRgba(img image.Image) *image.RGBA {
	if rgba, ok := img.(*image.RGBA); ok {
		return rgba
	}

	bounds := img.Bounds()
	rgba := image.NewRGBA(bounds)
	rm := rgba.ColorModel()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, rm.Convert(img.At(x, y)))
		}
	}
	return rgba
}

func sdlError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}

// Text pairs a string of text with a font.
// Pts is given in PostScript points.
type Text struct {
	Font string
	Pts float64
	string
}

// Sprite represents an image, the portion of
// the image to be rendered, and its shading.
type Sprite struct {
	Name string
	Bounds Rectangle
	Shade float32
}

func (ui *Ui) SetColor(r, g, b, a uint8) {
	C.SDL_SetRenderDrawColor(ui.rend,
		C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a))
}

/*
Draw queues a rendering of x and returns the dimensions of what
will be rendered, or an error. Draw supports the following types:

	string
		The given string is drawn at p in the default font, in the
		current color, on a background of the opposite color.

	Text
		The given text is drawn at p, using the font given in the text.

	Rectangle
		The given rectangle is filled at offset p, in the current color.

	Sprite
		The given image is drawn at offset p.

	image.Image
		The given image is drawn at offset p.
*/
func (ui *Ui) Draw(i interface{}, p Point) (Point, error) {
	switch d := i.(type) {
	case Rectangle:
		loc := d.Min.Add(p)
		fillRect(ui, int(loc.X), int(loc.Y), int(d.Dx()), int(d.Dy()))
		return d.Size(), nil
	case Sprite:
		return d.Bounds.Size(), drawSprite(ui, d, p)
	case Text:
		return drawText(ui, d, p)
	case string:
		return drawText(ui, Text{ "prstartk", 16.0, d }, p)
	}
	panic("That's not a thing to draw")
}

func fillRect(ui *Ui, x, y, w, h int) {
	C.SDL_RenderFillRect(ui.rend, &C.SDL_Rect{C.int(x), C.int(y), C.int(w), C.int(h)})
}

func drawSprite(ui *Ui, s Sprite, p Point) error {
	img, err := loadImg(ui, "resrc/" + s.Name + ".png")
	if err != nil {
		return err
	}
	img.Draw(ui, s, p)
	return nil
}

func (img *sdlImg) Draw(ui *Ui, s Sprite, p Point) {
	C.SDL_RenderCopy(ui.rend, img.tex,
		&C.SDL_Rect{C.int(s.Bounds.Min.X), C.int(s.Bounds.Min.Y), C.int(s.Bounds.Dx()), C.int(s.Bounds.Dy())},
		&C.SDL_Rect{C.int(p.X), C.int(p.Y), C.int(s.Bounds.Dx()), C.int(s.Bounds.Dy())})
}

func drawText(ui *Ui, t Text, p Point) (Point, error) {
	font, ok := ui.fontCache[t.Font]
	if !ok {
		var err error
		if font, err = NewFont("resrc/" + t.Font + ".ttf"); err != nil {
			return Point{}, err
		}
	}

	var r, g, b, a C.Uint8
	C.SDL_GetRenderDrawColor(ui.rend, &r, &g, &b, &a)
	font.SetColor(color.RGBA{ uint8(r), uint8(g), uint8(b), uint8(a) })
	font.SetSize(t.Pts)

	img, err := font.Render(t.string)
	if err != nil {
		return Point{}, err
	}

	return Pt(float64(img.Bounds().Dx()), float64(img.Bounds().Dy())),
		drawImage(ui, img, p)
}

func drawImage(ui *Ui, i image.Image, p Point) error {
	s, err := newSdlImage(ui, i, "")
	if err != nil {
		return err
	}
	defer s.Close()

	s.Draw(ui, Sprite{ Bounds: toRect(i.Bounds()), Shade: 1.0 }, p)
	return nil
}

// BUG(mccoyst): barf
func toRect(r image.Rectangle) Rectangle {
	return Rect(float64(r.Min.X), float64(r.Min.Y),
		float64(r.Max.X), float64(r.Max.Y))
}
