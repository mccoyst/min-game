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

	virtual std::shared_ptr<ui::Img> LoadImg(const char *path);
	virtual void Draw(const ui::Vec3&, std::shared_ptr<ui::Img> img);
};

struct SfmlImg : public ui::Img{
	sf::Image img;
	SfmlImg(const char *path);
};


SfmlUi::SfmlUi(ui::Len w, ui::Len h, const char *title)
	: Ui(w, h),
	win(sf::VideoMode(w.value(), h.value()), title){
}

void SfmlUi::Flip(){
	win.Display();
}

void SfmlUi::Clear(){
	win.Clear();
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
	sprite.SetX(loc.x.value());
	sprite.SetY(loc.y.value());
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