#include "ui.hpp"
#include "game.hpp"	// Fatal
#include <SDL/SDL.h>
#include <SDL/SDL_image.h>

class SdlUi : public ui::Ui {
	SDL_Surface *win;
public:
	SdlUi(Fixed w, Fixed h, const char *title);
	~SdlUi();
	virtual void Flip();
	virtual void Clear();
	virtual void Delay(float);
	virtual bool PollEvent(ui::Event&);
	virtual std::shared_ptr<ui::Img> LoadImg(const char*);
	virtual void Draw(const Vec3&, std::shared_ptr<ui::Img>);
};

struct SdlImg : public ui::Img {
	SDL_Surface *surf;
	SdlImg(const char*);
	~SdlImg();
};

SdlUi::SdlUi(Fixed w, Fixed h, const char *title) : Ui(w, h) {
	if (SDL_Init(SDL_INIT_VIDEO) == -1)
		throw Failure("Failed to initialized SDL video");
	Uint32 flags = SDL_SWSURFACE | SDL_DOUBLEBUF;
	win = SDL_SetVideoMode(w.whole(), h.whole(), 0, flags);
	if (!win)
		throw Failure("Failed to set SDL video mode");
}

SdlUi::~SdlUi() {
	SDL_FreeSurface(win);
	SDL_Quit();
}

void SdlUi::Flip() {
	SDL_Flip(win);
}

void SdlUi::Clear() {
	Uint32 sc = SDL_MapRGB(win->format, 0, 0, 0);
	SDL_FillRect(win, NULL, sc);
}

void SdlUi::Delay(float sec) {
	SDL_Delay(sec * 1000);
}

static bool getbutton(SDL_Event *sdle, ui::Event &e) {
	switch (sdle->button.button) {
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
			if (!getbutton(&sdle, e))
				continue;
			return true;

		case SDL_MOUSEBUTTONUP:
			e.type = ui::Event::MouseUp;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(&sdle, e))
				continue;
			return true;


		case SDL_MOUSEMOTION:
			e.type = ui::Event::MouseMoved;
			e.x = sdle.motion.x;
			e.y = sdle.motion.y;
			return true;

		default:
			// ignore
			break;
		}
	}
	return false;
}

void SdlUi::Draw(const Vec3 &l, std::shared_ptr<ui::Img> img) {
	SDL_Rect dst;
	dst.x = l.x.whole();
	dst.y = l.y.whole();
	SDL_Surface *surf = static_cast<SdlImg*>(img.get())->surf;
	SDL_BlitSurface(surf, NULL, win, &dst);
}

SdlImg::SdlImg(const char *path) {
	surf = IMG_Load(path);
	if (!surf)
		throw Failure("Failed to load image %s", path);
}

SdlImg::~SdlImg() {
	SDL_FreeSurface(surf);
}

std::shared_ptr<ui::Img> SdlUi::LoadImg(const char *path) {
	return std::shared_ptr<ui::Img>(new SdlImg(path));
}

std::unique_ptr<ui::Ui> ui::OpenWindow(Fixed w, Fixed h, const char *title) {
	return std::unique_ptr<ui::Ui>(new SdlUi(w, h, title));
}