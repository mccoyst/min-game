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
static void loadingText(std::shared_ptr<Ui>, std::shared_ptr<Font>);

int main(int argc, char *argv[]) try{
	parseArgs(argc, argv);

	Fixed width(800), height(600);
	auto win = std::make_shared<Ui>(width, height, "Minima");

	auto font = LoadFont("resrc/prstartk.ttf", 12, 255, 255, 255);
	loadingText(win, font);

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

static void loadingText(std::shared_ptr<Ui> win, std::shared_ptr<Font> font) {
	auto img = font->Render("Generating World");
	win->Clear();
	win->Draw(Vec2(Fixed(0), Fixed(0)), img);
	win->Flip();
}
