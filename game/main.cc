#include "ui.hpp"
#include "world.hpp"
#include "game.hpp"
#include <cstdio>
#include <cstring>
#include <cassert>
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

// drawFps, when set to true, draws the frames per-second on
// the screen.
bool drawFps;

static void parseArgs(int, char*[]);
static void loadingText(std::shared_ptr<Ui>, std::shared_ptr<Font>);
static void doFps(std::shared_ptr<Ui>, std::shared_ptr<Font>, unsigned long);

int main(int argc, char *argv[]) try{
	parseArgs(argc, argv);

	Fixed width(800), height(600);
	std::shared_ptr<Ui> win(OpenWindow(width, height, "Minima"));

	std::shared_ptr<Font> font = LoadFont("resrc/prstartk.ttf", 12, 255, 255, 255);
	loadingText(win, font);

	// Must create the world *after* the window because
	// the world also loads some images.
	World world(stdin);
	world.Center(win, world.x0, world.y0);

	std::shared_ptr<Img> guy = LoadImg("resrc/Astronaut.png");
	Vec2 guyloc(Fixed(world.x0) * World::TileW, Fixed(world.y0) * World::TileH);

	bool running = true;
	bool drag = false;
	Fixed scrollx(0), scrolly(0), mul(1);
	int x0 = 0, y0 = 0;

	// Compute the mean frame time.
	double meanTime = 0;
	unsigned long nFrames = 0;
	unsigned long t0 = win->Ticks();
	unsigned long t1 = t0;

	while(running){
		unsigned long frameTime = t1 - t0;
		t0 = t1;

		win->Clear();
		world.Draw(win);
		win->Draw(guyloc + world.Offset(), guy);
		doFps(win, font, frameTime);
		win->Flip();

		Event e;
		while (win->PollEvent(e)) {
			Fixed amt(0);
			switch (e.type) {
			case Event::Closed:
				running = false;

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

		t1 = win->Ticks();
		nFrames++;
		meanTime = meanTime + (t1-t0 - meanTime)/nFrames;
		if (t0 + FrameMsec > t1)
			win->Delay(t0 + FrameMsec - t1);
	}

	printf("%lu frames\n", nFrames);
	printf("Mean frame time: %g msec\n", meanTime);

	return 0;
}catch (const Failure &f) {
	fputs(f.msg, stderr);
	fputc('\n', stderr);
	return 1;
}

static void parseArgs(int argc, char *argv[]) {
	for (int i = 1; i < argc; i++) {
		if (strcmp(argv[i], "-heights") == 0)
			drawHeights = true;
		else if (strcmp(argv[i], "-fps") == 0)
			drawFps = true;
	}
}

static void loadingText(std::shared_ptr<Ui> win, std::shared_ptr<Font> font) {
	std::shared_ptr<Img> img = font->Render("Generating World");
	win->Clear();
	win->Draw(Vec2(Fixed(0), Fixed(0)), img);
	win->Flip();
}

static void doFps(std::shared_ptr<Ui> win, std::shared_ptr<Font> font,
		unsigned long msec) {
	if (!drawFps)
		return;
	unsigned long rate = 1.0 / (msec / 1000.0);
	win->Draw(Vec2(Fixed(0), Fixed(0)), font->Render("%lu fps", rate));
}