#include <ui.hpp>
#include "world.hpp"
#include <cstdio>

int main(){
	World world(stdin);

	ui::Len width(640), height(480);
	std::unique_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));

	bool drag = false;
	int x0 = 0, y0 = 0;

	for ( ; ; ) {
		win->Clear();
		world.Draw(*win);
		win->Flip();

		ui::Event e;
		while (win->PollEvent(e)) {
			switch (e.type) {
			case ui::Event::MouseDown:
				drag = true;
				x0 = e.x;
				y0 = e.y;
				break;

			case ui::Event::MouseUp:
				drag = false;
				break;

			case ui::Event::MouseMoved:
				if (!drag)
					break;
				world.Scroll(e.x - x0, e.y - y0);
				x0 = e.x;
				y0 = e.y;
				break;

			case ui::Event::Closed:
				goto out;
			}
		}
		win->Delay(0.02);
	}

out:
	return 0;
}
