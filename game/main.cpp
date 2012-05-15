// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
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

static void parseArgs(int, char*[]);
static void loadingText(Ui &, Font*);

int main(int argc, char *argv[]) try{
	parseArgs(argc, argv);

	Fixed width(800), height(600);
	Ui win (width, height, "Minima");

	auto font = std::unique_ptr<Font>(LoadFont("resrc/prstartk.ttf", 12, 255, 255, 255));
	loadingText(win, font.get());

	// Must create the world *after* the window because
	// the world also loads some images.
	World world(stdin);

	ScreenStack stk(win,
		std::shared_ptr<Screen>(new ExploreScreen(win, world)));
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
	}
}

static void loadingText(Ui &win, Font *font) {
	auto img = std::unique_ptr<Img>(font->Render("Generating World"));
	win.Clear();
	win.Draw(Vec2(Fixed(0), Fixed(0)), img.get());
	win.Flip();
}
