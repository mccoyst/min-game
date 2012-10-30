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

func (b Button) String() string {
	switch b {
	case Left:
		return "Left"
	case Right:
		return "Right"
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Shoot:
		return "Shoot"
	case Bomb:
		return "Bomb"
	}
	return "?Unkown Button?"
}

var DefaultKeymap = map[KeyCode]Button{
	KeyCode(C.SDLK_s): Left,
	KeyCode(C.SDLK_f): Right,
	KeyCode(C.SDLK_e): Up,
	KeyCode(C.SDLK_d): Down,
	KeyCode(C.SDLK_j): Shoot,
	KeyCode(C.SDLK_k): Bomb,
}

type Ui interface {
	Quit()
	Clear()
	LoadImg(string) (Img, error)
	FillRect(x, y, w, h int)
	SetColor(r, g, b, a uint8)
	Show()
	PollEvent() Event
}

// SDL-specific:

type sdl struct {
	win  *C.SDL_Window
	rend *C.SDL_Renderer
}

func NewUi(title string, w, h int) (Ui, error) {
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

	return &sdl{win: win, rend: rend}, nil
}

func (ui *sdl) Quit() {
	C.SDL_DestroyRenderer(ui.rend)
	C.SDL_DestroyWindow(ui.win)
	C.SDL_Quit()
}

func (ui *sdl) PollEvent() Event {
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

func (ui *sdl) SetColor(r, g, b, a uint8) {
	C.SDL_SetRenderDrawColor(ui.rend,
		C.Uint8(r), C.Uint8(g), C.Uint8(b), C.Uint8(a))
}

func (ui *sdl) Clear() {
	C.SDL_RenderClear(ui.rend)
}

func (ui *sdl) FillRect(x, y, w, h int) {
	C.SDL_RenderFillRect(ui.rend, &C.SDL_Rect{C.int(x), C.int(y), C.int(w), C.int(h)})
}

func (ui *sdl) Show() {
	C.SDL_RenderPresent(ui.rend)
}

type Img interface {
	Draw(Ui, Point, float32)
}

type sdlImg struct {
	tex *C.SDL_Texture
}

func (ui *sdl) LoadImg(path string) (Img, error) {
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

	rgba := RGBA(img, path)
	tex := newTex(C.sdl_rgba_fmt(C.int(isLE)))
	e := C.SDL_UpdateTexture(tex, nil, unsafe.Pointer(&rgba.Pix[0]), C.int(rgba.Stride))
	if e != 0 {
		C.SDL_DestroyTexture(tex)
		return nil, sdlError()
	}
	C.SDL_SetTextureBlendMode(tex, C.SDL_BLENDMODE_BLEND)
	return sdlImg{tex}, nil
}

// BUG(mccoyst): RGBA assumes the image bounds starts at (0,0).
func RGBA(img image.Image, name string) *image.RGBA {
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

func (img sdlImg) Draw(ui Ui, p Point, shade float32) {
	sdl := ui.(*sdl)
	var format C.Uint32
	var w, h, access C.int
	C.SDL_QueryTexture(img.tex, &format, &access, &w, &h)
	C.SDL_RenderCopy(sdl.rend, img.tex, nil, &C.SDL_Rect{
		C.int(p.X), C.int(p.Y), w, h})
}

func sdlError() error {
	return errors.New(C.GoString(C.SDL_GetError()))
}
