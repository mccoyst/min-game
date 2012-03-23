#include <ui.hpp>
#include "world.hpp"

int main(){
	World world(stdin);

	ui::Len width(40), height(30);
	std::unique_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));

	win->Clear();
	win->Flip();

	world.Draw(*win);
	win->Flip();

	sleep(2);

	return 0;
}
