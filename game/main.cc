#include "ui.hpp"
#include "world.hpp"
#include "game.hpp"
#include <cstdio>
#include <cstring>
#include <SDL_main.h>

enum {
	// FrameMsec is the minimum frame time.
	FrameMsec = 20,

	// ScrollSpd is the amount to scroll per-frame
	// when an arrow key is held.
	ScrollSpd = 10,
};

// drawHeights, when set to true makes the world draw the
// heigth of each tile on it.
bool drawHeights;

static void parseArgs(int, char*[]);
static void loadingText(std::shared_ptr<Ui>);

int main(int argc, char *argv[]) try{
	parseArgs(argc, argv);

	Fixed width(800), height(600);
	std::shared_ptr<Ui> win(OpenWindow(width, height, "Minima"));
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

		Event e;
		while (win->PollEvent(e)) {
			Fixed amt(0);
			switch (e.type) {
			case Event::Closed:
				goto out;

			case Event::MouseDown:
				scrollx = scrolly = Fixed(0);
				drag = true;
				x0 = e.x;
				y0 = e.y;
				break;

			case Event::MouseUp:
				drag = false;
				break;

			case Event::MouseMoved:
				if (!drag)
					break;
				world.Scroll(Fixed(e.x - x0), Fixed(y0 - e.y));
				x0 = e.x;
				y0 = e.y;
				break;

			case Event::KeyDown:
			case Event::KeyUp:
				if (e.type == Event::KeyDown)
					amt = Fixed(ScrollSpd);

				switch (e.button) {
				case Event::KeyDownArrow:
					scrolly = amt;
					break;
				case Event::KeyUpArrow:
					scrolly = -amt;
					break;
				case Event::KeyLeftArrow:
					scrollx = amt;
					break;
				case Event::KeyRightArrow:
					scrollx = -amt;
					break;
				case Event::KeyLShift:
				case Event::KeyRShift:
					if (e.type == Event::KeyDown)
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

static void parseArgs(int argc, char *argv[]) {
	for (int i = 1; i < argc; i++) {
		if (strcmp(argv[i], "-heights") == 0)
			drawHeights = true;
	}
}

static void loadingText(std::shared_ptr<Ui> win) {
	std::shared_ptr<Font> font = LoadFont(
		"resrc/prstartk.ttf", 16, 255, 255, 255
	);
	std::shared_ptr<Img> img = font->Render("Generating World");
	win->Clear();
	win->Draw(Vec3(Fixed(0), Fixed(0)), img);
	win->Flip();
}