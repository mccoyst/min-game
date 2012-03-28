#include "ui.hpp"
#include "game.hpp"
#include <SDL/SDL.h>
#include <SDL/SDL_opengl.h>
#include <SDL_image.h>
#include <SDL_ttf.h>
#include <cstdarg>
#include <cassert>

class SdlUi : public ui::Ui {
	SDL_Surface *win;
	unsigned long tick0;
public:
	SdlUi(Fixed w, Fixed h, const char *title);
	~SdlUi();
	virtual void Flip();
	virtual void Clear();
	virtual void Delay(unsigned long);
	virtual unsigned long Ticks();
	virtual bool PollEvent(ui::Event&);
	virtual void Draw(const Vec3&, std::shared_ptr<ui::Img>);
	virtual void Shade(const Vec3&, const Vec3&, float);
};

struct SdlImg : public ui::Img {
	GLuint texId;
	unsigned int w, h;

	SdlImg(SDL_Surface*);
	virtual ~SdlImg();
	virtual unsigned int Width() { return w; }
	virtual unsigned int Height() { return h; }
};

struct SdlFont : public ui::Font {
	TTF_Font *font;
	char r, g, b;

	SdlFont(const char *, int, char, char, char);
	virtual ~SdlFont();
};

SdlUi::SdlUi(Fixed w, Fixed h, const char *title) : Ui(w, h) {
	if (SDL_Init(SDL_INIT_VIDEO) == -1)
		throw Failure("Failed to initialized SDL video");
	tick0 = SDL_GetTicks();

	win = SDL_SetVideoMode(w.whole(), h.whole(), 0, SDL_OPENGL);
	if (!win)
		throw Failure("Failed to set SDL video mode");

	int imgflags = IMG_INIT_PNG;
	if ((IMG_Init(imgflags) & imgflags) != imgflags)
		throw Failure("Failed to initialize png support: %s", IMG_GetError());

	if (TTF_Init() == -1)
		throw Failure("Failed to initialize SDL_ttf: %s", TTF_GetError());

	glEnable(GL_TEXTURE_2D);
	glEnable(GL_BLEND);
	glBlendFunc(GL_SRC_ALPHA, GL_ONE);
	glMatrixMode(GL_PROJECTION);
	glLoadIdentity();
	glClearColor(0.0, 0.0, 0.0, 0.0);
	gluOrtho2D(0, w.whole(), 0, -h.whole());
	glMatrixMode(GL_MODELVIEW);
	glLoadIdentity();
	glTranslatef(0.0, -h.whole(), 0.0);
}

SdlUi::~SdlUi() {
	SDL_FreeSurface(win);
	TTF_Quit();
	IMG_Quit();
	SDL_Quit();
}

void SdlUi::Flip() {
	SDL_GL_SwapBuffers();
}

void SdlUi::Clear() {
	glClear(GL_COLOR_BUFFER_BIT);
}

void SdlUi::Delay(unsigned long msec) {
	SDL_Delay(msec);
}

unsigned long SdlUi::Ticks() {
	return tick0 - SDL_GetTicks();
}

static bool getbutton(SDL_Event &sdle, ui::Event &e) {
	switch (sdle.button.button) {
	case SDL_BUTTON_LEFT:
		e.button = ui::Event::MouseLeft;
		break;
	case SDL_BUTTON_RIGHT:
		e.button = ui::Event::MouseRight;
		break;
	case SDL_BUTTON_MIDDLE:
		e.button = ui::Event::MouseCenter;
		break;
	default:
		return false;
	};
	return true;
}

static bool getkey(SDL_Event &sdle, ui::Event &e) {
	switch (sdle.key.keysym.sym) {
	case SDLK_UP:
		e.button = ui::Event::KeyUpArrow;
		break;
	case SDLK_DOWN:
		e.button = ui::Event::KeyDownArrow;
		break;
	case SDLK_LEFT:
		e.button = ui::Event::KeyLeftArrow;
		break;
	case SDLK_RIGHT:
		e.button = ui::Event::KeyRightArrow;
		break;
	case SDLK_RSHIFT:
		e.button = ui::Event::KeyRShift;
		break;
	case SDLK_LSHIFT:
		e.button = ui::Event::KeyLShift;
		break;
	default:
		if (sdle.key.keysym.sym < 'a' || sdle.key.keysym.sym > 'z')
			return false;
		e.button = sdle.key.keysym.sym;
	}

	return true;
}

bool SdlUi::PollEvent(ui::Event &e) {
	SDL_Event sdle;
	while (SDL_PollEvent(&sdle)) {
		switch (sdle.type) {
		case SDL_QUIT:
			e.type = ui::Event::Closed;
			return true;

		case SDL_MOUSEBUTTONDOWN:
			e.type = ui::Event::MouseDown;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;

		case SDL_MOUSEBUTTONUP:
			e.type = ui::Event::MouseUp;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;


		case SDL_MOUSEMOTION:
			e.type = ui::Event::MouseMoved;
			e.x = sdle.motion.x;
			e.y = sdle.motion.y;
			return true;

		case SDL_KEYUP:
			e.type = ui::Event::KeyUp;
			if (!getkey(sdle, e))
				continue;
			return true;

		case SDL_KEYDOWN:
			e.type = ui::Event::KeyDown;
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

void SdlUi::Draw(const Vec3 &l, std::shared_ptr<ui::Img> _img) {
	SdlImg *img = static_cast<SdlImg*>(_img.get());
	float x = l.x.whole(), y = l.y.whole();

	glBindTexture(GL_TEXTURE_2D, img->texId);

	glBegin(GL_QUADS);
	glTexCoord2i(0, 0);
	glVertex3f(x, y, 0);
	glTexCoord2i(1, 0);
	glVertex3f(x+img->w, y, 0);
	glTexCoord2i(1, 1);
	glVertex3f(x+img->w, y+img->h, 0);
	glTexCoord2i(0, 1);
	glVertex3f(x, y+img->h, 0);
	glEnd();
}

void SdlUi::Shade(const Vec3 &l, const Vec3 &sz, float f) {
	float x = l.x.whole(), y = l.y.whole();
	float w = sz.x.whole(), h = sz.y.whole();

	if (f < 0)
		f = 0;
	if (f > 1)
		f = 1;
	glColor4f(0.5, 0.5, 0.5, f);

	glBegin(GL_QUADS);
	glVertex3f(x, y, 0);
	glVertex3f(x+w, y, 0);
	glVertex3f(x+w, y+h, 0);
	glVertex3f(x, y+h, 0);
	glEnd();
}

SdlImg::SdlImg(SDL_Surface *surf) : w(surf->w), h(surf->h) {
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

	glGenTextures(1, &texId);
	glBindTexture(GL_TEXTURE_2D, texId);
 
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
 
	glTexImage2D(GL_TEXTURE_2D, 0, pxSz, surf->w, surf->h, 0,
		texFormat, GL_UNSIGNED_BYTE, surf->pixels);
}

SdlImg::~SdlImg() {
	glDeleteTextures(1, &texId);
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

std::shared_ptr<ui::Ui> ui::OpenWindow(Fixed w, Fixed h, const char *title) {
	return std::shared_ptr<ui::Ui>(new SdlUi(w, h, title));
}

std::shared_ptr<ui::Img> ui::LoadImg(const char *path) {
	SDL_Surface *surf = IMG_Load(path);
	if (!surf)
		throw Failure("Failed to load image %s", path);
	std::shared_ptr<ui::Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}

std::shared_ptr<ui::Font> ui::LoadFont(const char *path, int sz, char r, char g, char b) {
	return std::shared_ptr<ui::Font>(new SdlFont(path, sz, r, g, b));
}

std::shared_ptr<ui::Img> ui::RenderText(std::shared_ptr<ui::Font> f, const char *fmt, ...) {
	SdlFont *font = dynamic_cast<SdlFont*>(f.get());
	assert (font);

	char s[256];
	va_list ap;
	va_start(ap, fmt);
	vsnprintf(s, sizeof(s), fmt, ap);
	va_end(ap);

	SDL_Color c;
	c.r = font->r;
	c.g = font->g;
	c.b = font->b;
	SDL_Surface *surf = TTF_RenderUTF8_Blended(font->font, s, c);
	if (!surf)
		throw Failure("Failed to render text: %s", TTF_GetError());

	std::shared_ptr<ui::Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}