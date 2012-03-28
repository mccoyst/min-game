#include "ui.hpp"
#include "world.hpp"
#include <cstdio>
#include <google/profiler.h>

enum {
	// FrameMsec is the minimum frame time.
	FrameMsec = 20,

	// ScrollSpd is the amount to scroll per-frame
	// when an arrow key is held.
	ScrollSpd = 10,
};

int main(){
	Fixed width(640), height(480);
	std::unique_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));

	// Must create the world *after* the window because
	// the world also loads some images.
	World world(stdin);

	bool drag = false;
	Fixed scrollx(0), scrolly(0), mul(1);
	int x0 = 0, y0 = 0;

	for ( ; ; ) {
		unsigned long t0 = win->Ticks();
		win->Clear();
		world.Draw(*win);
		win->Flip();

		ui::Event e;
		while (win->PollEvent(e)) {
			Fixed amt(0);
			switch (e.type) {
			case ui::Event::Closed:
				goto out;

			case ui::Event::MouseDown:
				scrollx = scrolly = Fixed(0);
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

			case ui::Event::KeyDown:
			case ui::Event::KeyUp:
				if (e.type == ui::Event::KeyDown)
					amt = Fixed(ScrollSpd);

				switch (e.button) {
				case ui::Event::KeyDownArrow:
					scrolly = Fixed(0)-amt;
					break;
				case ui::Event::KeyUpArrow:
					scrolly = amt;
					break;
				case ui::Event::KeyLeftArrow:
					scrollx = amt;
					break;
				case ui::Event::KeyRightArrow:
					scrollx = Fixed(0)-amt;
					break;
				case ui::Event::KeyLShift:
				case ui::Event::KeyRShift:
					if (e.type == ui::Event::KeyDown)
						mul = Fixed(5);
					else
						mul = Fixed(1);
				default:
					// ignore
					break;
				}
				break;				
	
			default:
				// ignore
				break;
			}
		}

		world.Scroll(scrollx*mul, scrolly*mul);

		unsigned long t1 = win->Ticks();
		if (t0 + FrameMsec > t1)
			win->Delay(t0 + FrameMsec - t1);
	}

out:
	return 0;
}
