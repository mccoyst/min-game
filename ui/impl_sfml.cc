#include <ui.hpp>
#include <SFML/Window.hpp>
#include <SFML/Graphics.hpp>
#include <cassert>

namespace{

class SfmlUi : public ui::Ui{
	sf::RenderWindow win;
public:
	SfmlUi(ui::Len w, ui::Len h, const char *title);

	virtual void Flip();
	virtual void Clear();
	virtual void Delay(float);
	virtual bool PollEvent(ui::Event&);

	virtual std::shared_ptr<ui::Img> LoadImg(const char *path);
	virtual void Draw(const ui::Vec3&, std::shared_ptr<ui::Img> img);
};

struct SfmlImg : public ui::Img{
	sf::Image img;
	SfmlImg(const char *path);
};


SfmlUi::SfmlUi(ui::Len w, ui::Len h, const char *title)
	: Ui(w, h),
	win(sf::VideoMode(w.whole(), h.whole()), title){
}

void SfmlUi::Flip(){
	win.Display();
}

void SfmlUi::Clear(){
	win.Clear();
}

void SfmlUi::Delay(float sec){
	sf::Clock c;
	c.Reset();
	while (c.GetElapsedTime() < sec)
		;
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
			return true;
	
		case sf::Event::MouseButtonReleased:
			e.type = ui::Event::MouseUp;
			e.x = sfe.MouseButton.X;
			e.y = sfe.MouseButton.Y;
			return true;
	
		case sf::Event::MouseMoved:
			e.type = ui::Event::MouseMoved;
			e.x = sfe.MouseMove.X;
			e.y = sfe.MouseMove.Y;
			return true;
		}
	}
}

std::shared_ptr<ui::Img> SfmlUi::LoadImg(const char *path){
	return std::shared_ptr<ui::Img>(new SfmlImg(path));
}

void SfmlUi::Draw(const ui::Vec3 &loc, std::shared_ptr<ui::Img> img){
	SfmlImg *s = dynamic_cast<SfmlImg*>(img.get());
	assert(s != nullptr);
	//TODO: optimize this
	sf::Sprite sprite;
	sprite.SetImage(s->img);
	sprite.SetX(loc.x.whole());
	sprite.SetY(loc.y.whole());
	win.Draw(sprite);
}


SfmlImg::SfmlImg(const char *path){
	//TODO: check return
	img.LoadFromFile(path);
}

}

namespace ui {
std::unique_ptr<ui::Ui> OpenWindow(ui::Len w, ui::Len h, const char *title){
	return std::unique_ptr<ui::Ui>(new SfmlUi(w, h, title));
}
};