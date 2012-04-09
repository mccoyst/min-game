#include "ui.hpp"
#include "world.hpp"
#include "game.hpp"
#include "geom.hpp"
#include "screen.hpp"
#include <cstdio>
#include <cstring>
#include <cassert>
#include <SDL_main.h>

// drawHeights, when set to true makes the world draw the
// heigth of each tile on it.
bool drawHeights;

// drawFps, when set to true, draws the frames per-second on
// the screen.
bool drawFps;

static void parseArgs(int, char*[]);
static void loadingText(std::shared_ptr<Ui>, std::shared_ptr<Font>);

int main(int argc, char *argv[]) try{
	parseArgs(argc, argv);

	Fixed width(800), height(600);
	auto win(OpenWindow(width, height, "Minima"));

	auto font = LoadFont("resrc/prstartk.ttf", 12, 255, 255, 255);
	loadingText(win, font);

	// Must create the world *after* the window because
	// the world also loads some images.
	World world(stdin);

	ScreenStack stk(win);
	stk.SetDrawFps(drawFps);
	stk.Push(std::shared_ptr<Screen>(new ExploreScreen(win, world)));
	stk.Run();

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
	auto img = font->Render("Generating World");
	win->Clear();
	win->Draw(Vec2(Fixed(0), Fixed(0)), img);
	win->Flip();
}
