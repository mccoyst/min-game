#include "ui.hpp"
#include "world.hpp"
#include "game.hpp"
#include <cstdio>
#include <SDL_main.h>

enum {
	// FrameMsec is the minimum frame time.
	FrameMsec = 20,

	// ScrollSpd is the amount to scroll per-frame
	// when an arrow key is held.
	ScrollSpd = 10,
};

static void loadingText(std::shared_ptr<ui::Ui>);

int main(int argc, char *argv[]) try{
	Fixed width(640), height(480);
	std::shared_ptr<ui::Ui> win(ui::OpenWindow(width, height, "Minima"));
	loadingText(win);

	// Must create the world *after* the window because
	// the world also loads some images.
	World world(stdin);

	bool drag = false;
	Fixed scrollx(0), scrolly(0), mul(1);
	int x0 = 0, y0 = 0;

	for ( ; ; ) {
		unsigned long t0 = win->Ticks();
		win->Clear();
		world.Draw(win);
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
				world.Scroll(Fixed(e.x - x0), -Fixed(e.y - y0));
				x0 = e.x;
				y0 = e.y;
				break;

			case ui::Event::KeyDown:
			case ui::Event::KeyUp:
				if (e.type == ui::Event::KeyDown)
					amt = Fixed(ScrollSpd);

				switch (e.button) {
				case ui::Event::KeyDownArrow:
					scrolly = amt;
					break;
				case ui::Event::KeyUpArrow:
					scrolly = -amt;
					break;
				case ui::Event::KeyLeftArrow:
					scrollx = amt;
					break;
				case ui::Event::KeyRightArrow:
					scrollx = -amt;
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
}catch (const Failure &f) {
	fputs(f.msg, stdout);
	fputc('\n', stdout);
	return 1;
}

static void loadingText(std::shared_ptr<ui::Ui> win) {
	std::shared_ptr<ui::Font> font = ui::LoadFont(
		"resrc/prstartk.ttf", 16, 255, 255, 255
	);
	std::shared_ptr<ui::Img> img = ui::RenderText(
		font, "Generating World"
	);
	win->Clear();
	win->Draw(Vec3(Fixed(0), Fixed(0)), img);
	win->Flip();
}