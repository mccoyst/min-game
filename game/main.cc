#include "ui.hpp"
#include "world.hpp"
#include <cstdio>

enum {
	FrameMsec = 20,
};

int main(){
	World world(stdin);

	Fixed width(640), height(480);
	std::unique_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));

	bool drag = false;
	int x0 = 0, y0 = 0;

	for ( ; ; ) {
		unsigned long t0 = win->Ticks();
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
				world.Scroll(Fixed(e.x - x0), Fixed(e.y - y0));
				x0 = e.x;
				y0 = e.y;
				break;

			case ui::Event::Closed:
				goto out;
			}
		}
		unsigned long t1 = win->Ticks();
		if (t0 + FrameMsec > t1)
			win->Delay(t0 + FrameMsec - t1);
	}

out:
	return 0;
}
