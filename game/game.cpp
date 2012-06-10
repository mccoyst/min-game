// Copyright © 2012 the Minima Authors under the MIT license. See AUTHORS for the list of authors.
#include "game.hpp"
#include "world.hpp"
#include "ui.hpp"
#include "entities.hpp"

ExploreScreen::ExploreScreen(World &w)
	: world(w),
	view((ScreenDims.x/World::TileW).whole() + 2,
		(ScreenDims.y/World::TileH).whole() + 3,
		World::TileW.whole(),
		World::TileH.whole(),
		LoadImg("resrc/tiles.png")),
	astroimg(LoadImg("resrc/Astronaut.png")){
	astro.reset(new Astro{astroimg.get()});
}

ExploreScreen::~ExploreScreen() { }

void ExploreScreen::Update(ScreenStack&) {
	astro->Move();
}

void ExploreScreen::Draw(Ui &win) {
	win.CenterCam(astro->Box().min);
	win.Clear();
	world.Draw(win, view);
	astro->Draw(win);
	win.Flip();
}

void ExploreScreen::Handle(ScreenStack &stk, Event &e) {
	if(e.type != Event::KeyDown && e.type != Event::KeyUp)
		return;

	Fixed speed;
	if(e.type == Event::KeyDown)
		speed = Fixed{2};

	switch (e.button) {
	case Event::DownArrow:
		astro->vel.y = speed;
		break;
	case Event::UpArrow:
		astro->vel.y = -speed;
		break;
	case Event::LeftArrow:
		astro->vel.x = speed;
		break;
	case Event::RightArrow:
		astro->vel.x = -speed;
		break;
	case Event::Action:
		if(e.type == Event::KeyDown) stk.Pop();
		break;
	case Event::None:
	default:
		// ignore
		break;
	}

}
