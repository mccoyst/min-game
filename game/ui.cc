#include "ui.hpp"
#include <cassert>

ui::Ui::~Ui(){
	Fixed n (22);
	Fixed m (20);
	assert(n + m == Fixed(42));
}

ui::Img::~Img(){
}

ui::Font::~Font(){
}