// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "ui.hpp"
#include "game.hpp"
#include "io.hpp"
#include "screens.hpp"
#include "world.hpp"
#include <SDL_main.h>

using std::vector;

static void parseArgs(int, char*[]);

int main(int argc, char *argv[]) try{
	parseArgs(argc, argv);

	Ui win (ScreenDims.x, ScreenDims.y, "Minima");

	ScreenStack stk(win, NewTitleScreen());
	stk.Run();

	return 0;

}catch (const std::exception &f) {
	printf(cerr(), "Uncaught exception: \"%v\"\n", f.what());
	return 1;
}

static void parseArgs(int argc, char *argv[]) {
	extern bool worldOnStdin;
	extern bool drawHeights;
	vector<string> args (argv+1, argv+argc);

	for(auto &arg : args){
		if (arg == "-heights")
			drawHeights = true;
		if (arg == "-stdin")
			worldOnStdin = true;
	}
}

