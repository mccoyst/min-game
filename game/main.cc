#include <ui.hpp>
#include "world.hpp"

int main(){
	ui::Len width(640), height(480);
	std::unique_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));
	World w(stdin);

	return 0;
}
