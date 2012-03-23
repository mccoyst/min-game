#include <ui.hpp>
#include "world.hpp"

int main(){
	World world(stdin);

	ui::Len width(640), height(480);
	std::unique_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));

	for (unsigned int i = 0; i < 16*10; i++) {
		win->Clear();
		world.Draw(*win);
		win->Flip();
		world.xoff += 1;
		world.yoff += 1;
		win->Delay(0.02);
	}

	return 0;
}
