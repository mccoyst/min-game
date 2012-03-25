#include "ui.hpp"
#include <cassert>

ui::Ui::~Ui(){
	Len n (22);
	Len m (20);
	assert(n + m == Len(42));
}

ui::Img::~Img(){
}
