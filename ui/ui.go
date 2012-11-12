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
	Action
)

var ButtonNames = map[Button]string{
	Unknown: "Unknown",
	Left:    "Left",
	Right:   "Right",
	Up:      "Up",
	Down:    "Down",
	Action:  "Action",
}

func (b Button) String() string {
	if s, ok := ButtonNames[b]; ok {
		return s
	}
	return ButtonNames[Unknown]
}

var (
	DefaultKeymap = map[KeyCode]Button{
		KeyCode(C.SDLK_s): Left,
		KeyCode(C.SDLK_f): Right,
		KeyCode(C.SDLK_e): Up,
		KeyCode(C.SDLK_d): Down,
		KeyCode(C.SDLK_j): Action,
	}

	DvorakKeymap = map[KeyCode]Button{
		KeyCode(C.SDLK_o):      Left,
		KeyCode(C.SDLK_u):      Right,
		KeyCode(C.SDLK_PERIOD): Up,
		KeyCode(C.SDLK_e):      Down,
		KeyCode(C.SDLK_h):      Action,
	}
)

// SDL-specific:

type Ui struct {
	win  *C.SDL_Window
	rend *C.SDL_Renderer

	// Font is the current font.
	font *font

	// Color is the current color.
	color color.Color

	// NFrames is the number of frames drawn.
	nFrames uint64

	imgCache  map[string]*sdlImg
	fontCache map[string]*font
	txtCache  map[textKey]*cachedText
}

type textKey struct {
	txt        string
	size       float64
	r, g, b, a uint32
}

type cachedText struct {
	img   *sdlImg
	frame uint64
	rect  Rectangle
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
		C.SDL_WINDOW_SHOWN|C.SDL_WINDOW_OPENGL)
	if win == nil {
		return nil, sdlError()
	}

	rend := C.SDL_CreateRenderer(win, -1, C.SDL_RENDERER_ACCELERATED)
	if rend == nil {
		return nil, sdlError()
	}

	ui := &Ui{
		win:       win,
		rend:      rend,
		imgCache:  make(map[string]*sdlImg),
		fontCache: make(map[string]*font),
		txtCache:  make(map[textKey]*cachedText),
	}
	if err := ui.SetFont("prstartk", 12); err != nil {
		return nil, err
	}
	ui.SetColor(color.Black)
	return ui, nil
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
	for k, c := range ui.txtCache {
		if c.frame < ui.nFrames {
			delete(ui.txtCache, k)
			c.img.Close()
		}
	}
	ui.nFrames++
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

// Sprite represents an image, the portion of
// the image to be rendered, and its shading.
type Sprite struct {
	Name   string
	Bounds Rectangle
	Shade  float32
}

// SetColor sets the current drawing color.
func (ui *Ui) SetColor(col color.Color) {
	r, g, b, a := col.RGBA()
	r8 := uint8(float64(r) / 0xFFFF * 255)
	g8 := uint8(float64(g) / 0xFFFF * 255)
	b8 := uint8(float64(b) / 0xFFFF * 255)
	a8 := uint8(float64(a) / 0xFFFF * 255)
	C.SDL_SetRenderDrawColor(ui.rend,
		C.Uint8(r8), C.Uint8(g8), C.Uint8(b8), C.Uint8(a8))
	ui.color = col
	ui.font.setColor(ui.color)
}

// SetFont sets the current font face and size.
func (ui *Ui) SetFont(name string, sz float64) error {
	var ok bool
	if ui.font, ok = ui.fontCache[name]; !ok {
		var err error
		if ui.font, err = newFont("resrc/" + name + ".ttf"); err != nil {
			return err
		}
		ui.fontCache[name] = ui.font
	}
	ui.font.setSize(sz)
	ui.font.setColor(ui.color)
	return nil
}

// TextSize returns the size of the text when rendered in the current font.
func (ui *Ui) TextSize(txt string) Point {
	w := ui.font.width(txt)
	h := ui.font.extents().height
	return Pt(float64(w), float64(h))
}

/*
Draw queues a rendering of x and returns the dimensions of what
will be rendered, or an error. Draw supports the following types:

	string
		The given string is drawn at p in the current font, in the
		current color.

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
	case string:
		return drawText(ui, d, p)
	case image.Image:
		return drawImage(ui, d, p)
	}
	panic("That's not a thing to draw")
}

func fillRect(ui *Ui, x, y, w, h int) {
	C.SDL_RenderFillRect(ui.rend, &C.SDL_Rect{C.int(x), C.int(y), C.int(w), C.int(h)})
}

func drawSprite(ui *Ui, s Sprite, p Point) error {
	img, err := loadImg(ui, "resrc/"+s.Name+".png")
	if err != nil {
		return err
	}
	img.Draw(ui, s, p)
	return nil
}

func (img *sdlImg) Draw(ui *Ui, s Sprite, p Point) {
	if s.Shade < 1.0 {
		sh := C.Uint8(s.Shade * 255)
		C.SDL_SetTextureColorMod(img.tex, sh, sh, sh)
		defer C.SDL_SetTextureColorMod(img.tex, 255, 255, 255)
	}
	C.SDL_RenderCopy(ui.rend, img.tex,
		&C.SDL_Rect{C.int(s.Bounds.Min.X), C.int(s.Bounds.Min.Y), C.int(s.Bounds.Dx()), C.int(s.Bounds.Dy())},
		&C.SDL_Rect{C.int(p.X), C.int(p.Y), C.int(s.Bounds.Dx()), C.int(s.Bounds.Dy())})
}

// DrawText draws the string to the ui at the given point, 
// using the ui's current font, and current color.
func drawText(ui *Ui, txt string, p Point) (Point, error) {
	r, g, b, a := ui.color.RGBA()
	key := textKey{
		txt:  txt,
		size: ui.font.size,
		r:    r,
		g:    g,
		b:    b,
		a:    a,
	}
	var img *sdlImg
	c, ok := ui.txtCache[key]
	if ok {
		c.frame = ui.nFrames
		img = c.img
	} else {
		i, err := ui.font.render(txt)
		if err != nil {
			return Point{}, err
		}
		img, err = newSdlImage(ui, i, "")
		if err != nil {
			return Point{}, err
		}
		c = &cachedText{
			img,
			ui.nFrames,
			toRect(i.Bounds()),
		}
		ui.txtCache[key] = c
	}
	img.Draw(ui, Sprite{Bounds: c.rect, Shade: 1.0}, p)
	return Pt(float64(c.rect.Dx()), float64(c.rect.Dy())), nil
}

// DrawImage draws an image to the UI at the given point.
func drawImage(ui *Ui, i image.Image, p Point) (Point, error) {
	s, err := newSdlImage(ui, i, "")
	if err != nil {
		return Point{}, err
	}
	defer s.Close()

	s.Draw(ui, Sprite{Bounds: toRect(i.Bounds()), Shade: 1.0}, p)
	return Pt(float64(i.Bounds().Dx()), float64(i.Bounds().Dy())), nil
}

// BUG(mccoyst): barf
func toRect(r image.Rectangle) Rectangle {
	return Rect(float64(r.Min.X), float64(r.Min.Y),
		float64(r.Max.X), float64(r.Max.Y))
}
