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

	bounds := img.Bounds()
	newTex := func(format C.Uint32) *C.SDL_Texture {
		return C.SDL_CreateTexture(
			ui.rend,
			format,
			C.SDL_TEXTUREACCESS_STATIC,
			C.int(bounds.Dx()),
			C.int(bounds.Dy()))
	}

	rgba := asRgba(img, path)
	tex := newTex(C.sdl_rgba_fmt(C.int(isLE)))
	e := C.SDL_UpdateTexture(tex, nil, unsafe.Pointer(&rgba.Pix[0]), C.int(rgba.Stride))
	if e != 0 {
		C.SDL_DestroyTexture(tex)
		return nil, sdlError()
	}
	C.SDL_SetTextureBlendMode(tex, C.SDL_BLENDMODE_BLEND)
	si := &sdlImg{tex}
	ui.imgCache[path] = si
	return si, nil
}

// BUG(mccoyst): asRgba assumes the image bounds starts at (0,0).
func asRgba(img image.Image, name string) *image.RGBA {
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
// Fonts are named like "FontName:Pts", where the
// number given by Pts is the size in PostScript points.
type Text struct {
	Font string
	string
}

// Sprite represents an image, the portion of
// the image to be rendered, and its shading.
type Sprite struct {
	Name string
	Bounds image.Rectangle
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
		return Pt(float64(d.Bounds.Dx()), float64(d.Bounds.Dy())),
			drawImg(ui, d.Name, int(p.X), int(p.Y), d.Bounds.Min.X, d.Bounds.Min.Y, d.Bounds.Dx(), d.Bounds.Dy(), d.Shade)
	}
	panic("That's not a thing to draw")
}

func fillRect(ui *Ui, x, y, w, h int) {
	C.SDL_RenderFillRect(ui.rend, &C.SDL_Rect{C.int(x), C.int(y), C.int(w), C.int(h)})
}

func drawImg(ui *Ui, name string, x, y, subx, suby, w, h int, shade float32) error {
	img, err := loadImg(ui, "resrc/" + name + ".png")
	if err != nil {
		return err
	}
	C.SDL_RenderCopy(ui.rend, img.tex,
		&C.SDL_Rect{C.int(subx), C.int(suby), C.int(w), C.int(h)},
		&C.SDL_Rect{C.int(x), C.int(y), C.int(w), C.int(h)})
	return nil
}