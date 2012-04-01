#include "ui.hpp"
#include <cassert>

Ui::~Ui(){
	Fixed n (22);
	Fixed m (20);
	assert(n + m == Fixed(42));
}

Img::~Img(){
}

Font::~Font(){
}