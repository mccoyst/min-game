#include "ui.hpp"
#include "game.hpp"
#include "world.hpp"
#include "opengl.hpp"
#include "keyMap.hpp"
#include <SDL.h>
#include <SDL_image.h>
#include <SDL_ttf.h>
#include <cstdarg>
#include <cstddef>
#include <cassert>
#include <cstdio>

class SdlUi : public OpenGLUi {
	SDL_Surface *win;

public:
	SdlUi(Fixed w, Fixed h, const char *title);
	~SdlUi();
	virtual void Flip();
	virtual void Delay(unsigned long);
	virtual unsigned long Ticks();
	virtual bool PollEvent(Event&);
};

struct SdlImg : public OpenGLImg {
	Vec2 sz;
	SdlImg(SDL_Surface*);
	virtual ~SdlImg() { }
	virtual Vec2 Size() const { return sz; }
};

struct SdlFont : public Font {
	TTF_Font *font;
	char r, g, b;

	SdlFont(const char *, int, char, char, char);
	virtual ~SdlFont();
	virtual std::shared_ptr<Img> Render(const char*, ...);
};

SdlUi::SdlUi(Fixed w, Fixed h, const char *title) : OpenGLUi(w, h) {
	if (SDL_Init(SDL_INIT_VIDEO) == -1)
		throw Failure("Failed to initialized SDL video");

	SDL_GL_SetAttribute(SDL_GL_DOUBLEBUFFER, 1);

	win = SDL_SetVideoMode(w.whole(), h.whole(), 0, SDL_OPENGL);
	if (!win)
		throw Failure("Failed to set SDL video mode");

	fprintf(stderr, "Vendor: %s\nRenderer: %s\nVersion: %s\nShade Lang. Version: %s\n",
	glGetString(GL_VENDOR),
	glGetString(GL_RENDERER),
	glGetString(GL_VERSION),
	glGetString(GL_SHADING_LANGUAGE_VERSION));

	int imgflags = IMG_INIT_PNG;
	if ((IMG_Init(imgflags) & imgflags) != imgflags)
		throw Failure("Failed to initialize png support: %s", IMG_GetError());

	if (TTF_Init() == -1)
		throw Failure("Failed to initialize SDL_ttf: %s", TTF_GetError());

	InitOpenGL();
}

SdlUi::~SdlUi() {
	TTF_Quit();
	IMG_Quit();
	SDL_Quit();
}

void SdlUi::Flip() {
	SDL_GL_SwapBuffers();
}

void SdlUi::Delay(unsigned long msec) {
	SDL_Delay(msec);
}

unsigned long SdlUi::Ticks() {
	return SDL_GetTicks();
}

static bool getbutton(SDL_Event &sdle, Event &e) {
	switch (sdle.button.button) {
	case SDL_BUTTON_LEFT:
		e.button = Event::MouseLeft;
		break;
	case SDL_BUTTON_RIGHT:
		e.button = Event::MouseRight;
		break;
	case SDL_BUTTON_MIDDLE:
		e.button = Event::MouseCenter;
		break;
	default:
		return false;
	};
	return true;
}

static bool getkey(SDL_Event &sdle, Event &e) {
	switch (sdle.key.keysym.sym) {
	case SDLK_UP:
		e.button = KeyMap::KeyUpArrow;
		break;
	case SDLK_DOWN:
		e.button = KeyMap::KeyDownArrow;
		break;
	case SDLK_LEFT:
		e.button = KeyMap::KeyLeftArrow;
		break;
	case SDLK_RIGHT:
		e.button = KeyMap::KeyRightArrow;
		break;
	case SDLK_RSHIFT:
		e.button = KeyMap::KeyRShift;
		break;
	case SDLK_LSHIFT:
		e.button = KeyMap::KeyLShift;
		break;
	default:
		if (sdle.key.keysym.sym < 'a' || sdle.key.keysym.sym > 'z')
			return false;
		e.button = sdle.key.keysym.sym;
	}

	return true;
}

bool SdlUi::PollEvent(Event &e) {
	SDL_Event sdle;
	while (SDL_PollEvent(&sdle)) {
		switch (sdle.type) {
		case SDL_QUIT:
			e.type = Event::Closed;
			return true;

		case SDL_MOUSEBUTTONDOWN:
			e.type = Event::MouseDown;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;

		case SDL_MOUSEBUTTONUP:
			e.type = Event::MouseUp;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;


		case SDL_MOUSEMOTION:
			e.type = Event::MouseMoved;
			e.x = sdle.motion.x;
			e.y = sdle.motion.y;
			return true;

		case SDL_KEYUP:
			e.type = Event::KeyUp;
			if (!getkey(sdle, e))
				continue;
			return true;

		case SDL_KEYDOWN:
			e.type = Event::KeyDown;
			if (!getkey(sdle, e))
				continue;
			return true;

		default:
			// ignore
			break;
		}
	}
	return false;
}

SdlImg::SdlImg(SDL_Surface *surf) : sz(Fixed(surf->w), Fixed(surf->h)) {
	GLint pxSz = surf->format->BytesPerPixel;
	GLenum texFormat = GL_BGRA;
	switch (pxSz) {
	case 4:
		if (surf->format->Rmask == 0xFF)
			texFormat = GL_RGBA;
		break;
	case 3:
		if (surf->format->Rmask == 0xFF)
			texFormat = GL_RGB;
		else
			texFormat = GL_BGR;
		break;
	default:
		throw Failure("Bad image color type… apparently");
	}

	glGenTextures(1, &texid);
	glBindTexture(GL_TEXTURE_2D, texid);

	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);

	glTexImage2D(GL_TEXTURE_2D, 0, pxSz, surf->w, surf->h, 0,
		texFormat, GL_UNSIGNED_BYTE, surf->pixels);
}

SdlFont::SdlFont(const char *path, int sz, char _r, char _g, char _b)
		: r(_r), g(_g), b(_b) {
	font = TTF_OpenFont(path, sz);
	if (!font)
		throw Failure("Failed to load font %s: %s", path, TTF_GetError());
}

SdlFont::~SdlFont() {
	TTF_CloseFont(font);
}

std::shared_ptr<Img> SdlFont::Render(const char *fmt, ...) {
	char s[256];
	va_list ap;
	va_start(ap, fmt);
	vsnprintf(s, sizeof(s), fmt, ap);
	va_end(ap);

	SDL_Color c;
	c.r = r;
	c.g = g;
	c.b = b;
	SDL_Surface *surf = TTF_RenderUTF8_Blended(font, s, c);
	if (!surf)
		throw Failure("Failed to render text: %s", TTF_GetError());

	std::shared_ptr<Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}

std::shared_ptr<Ui> OpenWindow(Fixed w, Fixed h, const char *title) {
	return std::shared_ptr<Ui>(new SdlUi(w, h, title));
}

std::shared_ptr<Img> LoadImg(const char *path) {
	SDL_Surface *surf = IMG_Load(path);
	if (!surf)
		throw Failure("Failed to load image %s", path);
	std::shared_ptr<Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}

std::shared_ptr<Font> LoadFont(const char *path, int sz, char r, char g, char b) {
	return std::shared_ptr<Font>(new SdlFont(path, sz, r, g, b));
}
