#include "ui.hpp"
#include <SFML/Window.hpp>
#include <SFML/Graphics.hpp>
#include <cassert>
#include <ctime>
#include <cerrno>

namespace{

class SfmlUi : public ui::Ui{
	sf::RenderWindow win;
	sf::Clock ticks;
public:
	SfmlUi(Fixed w, Fixed h, const char *title);

	virtual void Flip();
	virtual void Clear();
	virtual void Delay(unsigned long);
	virtual unsigned long Ticks();
	virtual bool PollEvent(ui::Event&);

	virtual std::shared_ptr<ui::Img> LoadImg(const char *path);
	virtual void Draw(const Vec3&, std::shared_ptr<ui::Img> img);
	virtual void Shade(const Vec3&, const Vec3&, float);

};

struct SfmlImg : public ui::Img{
	sf::Image img;
	SfmlImg(const char *path);
};


SfmlUi::SfmlUi(Fixed w, Fixed h, const char *title)
	: Ui(w, h),
	win(sf::VideoMode(w.whole(), h.whole()), title){
	ticks.Reset();
	win.EnableKeyRepeat(false);
}

void SfmlUi::Flip(){
	win.Display();
}

void SfmlUi::Clear(){
	win.Clear();
}

void SfmlUi::Delay(unsigned long msec){
	struct timespec a, b;
	struct timespec *ts = &a, *rem = &b;
	ts->tv_sec = msec / 1000;
	ts->tv_nsec = (msec % 1000) * 1000000;
	int err = nanosleep(ts, rem);
	while (err == -1 && errno == EINTR) {
		struct timespec *tmp = ts;
		ts = rem;
		rem = tmp;
		err = nanosleep(ts, rem);
	}
	if (err == -1)
		abort();
}

unsigned long SfmlUi::Ticks() {
	return ticks.GetElapsedTime()*1000;
}

static bool getbutton(sf::Event &sfe, ui::Event &e) {
	switch (sfe.MouseButton.Button) {
	case sf::Mouse::Button::Left:
		e.button = ui::Event::MouseLeft;
		break;
	case sf::Mouse::Button::Right:
		e.button = ui::Event::MouseRight;
		break;
	case sf::Mouse::Button::Middle:
		e.button = ui::Event::MouseCenter;
		break;
	default:
		return false;
	}
	return true;
}

static bool getkey(sf::Event &sfe, ui::Event &e) {
	switch (sfe.Key.Code) {
	case sf::Key::Left:
		e.button = ui::Event::KeyLeftArrow;
		break;
	case sf::Key::Right:
		e.button = ui::Event::KeyRightArrow;
		break;
	case sf::Key::Up:
		e.button = ui::Event::KeyUpArrow;
		break;
	case sf::Key::Down:
		e.button = ui::Event::KeyDownArrow;
		break;
	case sf::Key::LShift:
		e.button = ui::Event::KeyLShift;
		break;
	case sf::Key::RShift:
		e.button = ui::Event::KeyRShift;
		break;
	default:
		if (sfe.Key.Code < 'a' || sfe.Key.Code > 'z')
			return false;
		e.button = sfe.Key.Code;
	}

	return true;
}

bool SfmlUi::PollEvent(ui::Event &e){
	sf::Event sfe;
	while (win.GetEvent(sfe)) {
		switch (sfe.Type) {
		case sf::Event::Closed:
			e.type = ui::Event::Closed;
			return true;
	
		case sf::Event::MouseButtonPressed:
			e.type = ui::Event::MouseDown;
			e.x = sfe.MouseButton.X;
			e.y = sfe.MouseButton.Y;
			if (!getbutton(sfe, e))
				continue;
			return true;
	
		case sf::Event::MouseButtonReleased:
			e.type = ui::Event::MouseUp;
			e.x = sfe.MouseButton.X;
			e.y = sfe.MouseButton.Y;
			if (!getbutton(sfe, e))
				continue;
			return true;
	
		case sf::Event::MouseMoved:
			e.type = ui::Event::MouseMoved;
			e.x = sfe.MouseMove.X;
			e.y = sfe.MouseMove.Y;
			return true;

		case sf::Event::KeyReleased:
			e.type = ui::Event::KeyUp;
			if (!getkey(sfe, e))
				continue;
			return true;

		case sf::Event::KeyPressed:
			e.type = ui::Event::KeyDown;
			if (!getkey(sfe, e))
				continue;
			return true;

		default:
			// ignore
			break;
		}
	}
	return false;
}

std::shared_ptr<ui::Img> SfmlUi::LoadImg(const char *path){
	return std::shared_ptr<ui::Img>(new SfmlImg(path));
}

void SfmlUi::Draw(const Vec3 &loc, std::shared_ptr<ui::Img> img){
	SfmlImg *s = dynamic_cast<SfmlImg*>(img.get());
	assert(s != nullptr);
	//TODO: optimize this
	sf::Vector2f pos(loc.x.whole(), loc.y.whole());
	sf::Sprite sprite(s->img, pos);
	sprite.SetBlendMode(sf::Blend::None);
	win.Draw(sprite);
}

void SfmlUi::Shade(const Vec3 &l, const Vec3 &sz, float f) {
	int x = l.x.whole(), y = l.y.whole();
	int w = sz.x.whole(), h = sz.y.whole();
	sf::Color c(0, 0, 0, 255*(1-f));
	sf::Shape s = sf::Shape::Rectangle(x, y, x+w, y+h, c);
	s.EnableOutline(false);
	win.Draw(s);
}

SfmlImg::SfmlImg(const char *path){
	//TODO: check return
	img.LoadFromFile(path);
	img.SetSmooth(false);
}

}

std::unique_ptr<ui::Ui> ui::OpenWindow(Fixed w, Fixed h, const char *title){
	return std::unique_ptr<ui::Ui>(new SfmlUi(w, h, title));
}
